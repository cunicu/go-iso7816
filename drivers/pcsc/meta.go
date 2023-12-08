package pcsc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ebfe/scard"

	"cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/devices/nitrokey"
	"cunicu.li/go-iso7816/devices/yubikey"
)

// See: https://learn.microsoft.com/en-us/windows/win32/api/winscard/nf-winscard-scardgetattrib
func (c *Card) Metadata() (meta map[string]string) {
	protos := map[scard.Protocol]string{
		scard.ProtocolUndefined: "undefined",
		scard.ProtocolT0:        "t0",
		scard.ProtocolT1:        "t1",
		scard.ProtocolAny:       "any",
	}

	states := map[scard.State]string{
		scard.Unknown:    "unknown",
		scard.Absent:     "absent",
		scard.Present:    "present",
		scard.Swallowed:  "swallowed",
		scard.Powered:    "powered",
		scard.Negotiable: "negotiable",
		scard.Specific:   "specific",
	}

	meta = map[string]string{}

	if sts, err := c.Status(); err == nil {
		meta["status.reader"] = sts.Reader
		meta["status.atr"] = hex.EncodeToString(sts.Atr)
		meta["status.active_protocol"] = bitsToString(sts.ActiveProtocol, protos)

		if sts.State != 0 {
			meta["status.state"] = bitsToString(sts.State, states)
		}
	}

	if data, err := c.GetAttrib(scard.AttrVendorName); err == nil {
		meta["attr.name.vendor"] = string(data)
	}

	if data, err := c.GetAttrib(scard.AttrDeviceSystemName); err == nil {
		meta["attr.name.system"] = string(data)
	}

	if data, err := c.GetAttrib(scard.AttrDeviceFriendlyName); err == nil {
		meta["attr.name.friendly"] = string(data)
	}

	if data, err := c.GetAttrib(scard.AttrVendorIfdSerialNo); err == nil {
		meta["attr.ifd.serial"] = string(data)
	}

	if data, err := c.GetAttrib(scard.AttrVendorIfdType); err == nil {
		meta["attr.ifd.type"] = string(data)
	}

	if data, err := c.GetAttrib(scard.AttrVendorIfdVersion); err == nil && len(data) == 4 {
		idata := binary.NativeEndian.Uint32(data)
		v := iso7816.Version{
			Major: int(idata>>24) & 0xff,
			Minor: int(idata>>16) & 0xff,
			Patch: int(idata & 0xffff),
		}

		meta["attr.ifd.version"] = v.String()
	}

	// https://ludovicrousseau.blogspot.com/2020/04/scardattrchannelid-and-usb-devices.html
	if data, err := c.GetAttrib(scard.AttrChannelId); err == nil && len(data) == 4 {
		idata := binary.NativeEndian.Uint32(data)

		chType := idata >> 16
		chNum := idata & 0xffff

		switch chType {
		case 0x01: // Serial I/O; chNum is a port number.
			meta["attr.channel.type"] = "serial"
			meta["attr.channel.serial.port"] = fmt.Sprint(chNum)

		case 0x02: // Parallel I/O; chNum is a port number.
			meta["attr.channel.type"] = "parallel"
			meta["attr.channel.parallel.port"] = fmt.Sprint(chNum)

		case 0x04: // PS/2 keyboard port; chNum is zero.
			meta["attr.channel.type"] = "ps/2"

		case 0x08: // SCSI; chNum is SCSI ID number.
			meta["attr.channel.type"] = "scsi"
			meta["attr.channel.scsi.id"] = fmt.Sprint(chNum)

		case 0x10: // IDE; chNum is device number.
			meta["attr.channel.type"] = "ide"
			meta["attr.channel.ide_dev_id"] = fmt.Sprint(chNum)

		case 0x20: // USB; chNum is device number.
			bus := (chNum & 0xFF00) >> 8
			addr := chNum & 0xFF

			meta["attr.channel.type"] = "usb"
			meta["attr.channel.usb.bus"] = fmt.Sprint(bus)
			meta["attr.channel.usb.addr"] = fmt.Sprint(addr)

		// Vendor-defined interface with y in the range zero through 15; chNum is vendor defined.
		case 0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff:
			meta["attr.channel.type"] = fmt.Sprintf("vendor:%2x", chType)
			meta["attr.channel.vendor"] = fmt.Sprint(chNum)
		}
	}

	ic := iso7816.NewCard(c)

	metadatas := map[string]func(*iso7816.Card) map[string]any{
		"yubikey":  yubikey.Metadata,
		"nitrokey": nitrokey.Metadata,
	}

	for prefix, metadata := range metadatas {
		for key, value := range metadata(ic) {
			var strValue string

			switch value := value.(type) {
			case []byte:
				strValue = hex.EncodeToString(value)
			default:
				strValue = fmt.Sprint(value)
			}

			meta[prefix+"."+key] = strValue
		}
	}

	return meta
}

func bitsToString[E ~uint32](v E, m map[E]string) string {
	vs := []string{}
	for m, s := range m {
		if v&m != 0 {
			vs = append(vs, s)
		}
	}
	return strings.Join(vs, ",")
}
