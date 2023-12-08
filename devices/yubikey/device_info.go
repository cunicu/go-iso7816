// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package yubikey

import (
	"encoding/binary"
	"fmt"
	"time"

	iso "cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/encoding/tlv"
)

const (
	TagCapsSupportedUSB tlv.Tag = 0x01
	TagSerialNumber     tlv.Tag = 0x02
	TagCapsEnabledUSB   tlv.Tag = 0x03
	TagFormFactor       tlv.Tag = 0x04
	TagFirmwareVersion  tlv.Tag = 0x05
	TagAutoEjectTimeout tlv.Tag = 0x06
	TagChalRespTimeout  tlv.Tag = 0x07
	TagDeviceFlags      tlv.Tag = 0x08
	TagAppVersions      tlv.Tag = 0x09
	TagConfigLock       tlv.Tag = 0x0a
	TagUnlock           tlv.Tag = 0x0b
	TagReboot           tlv.Tag = 0x0c
	TagCapsSupportedNFC tlv.Tag = 0x0d
	TagCapsEnabledNFC   tlv.Tag = 0x0e
)

type DeviceFlag byte

const (
	DeviceFlagRemoteWakeup DeviceFlag = 0x40
	DeviceFlagEject        DeviceFlag = 0x80
)

type Capability int

const (
	CapOTP     Capability = 0x01
	CapU2F     Capability = 0x02
	CapFIDO2   Capability = 0x200
	CapOATH    Capability = 0x20
	CapPIV     Capability = 0x10
	CapOpenPGP Capability = 0x08
	CapHSMAUTH Capability = 0x100
)

type FormFactor byte

const (
	FormFactorUnknown       FormFactor = 0x00
	FormFactorUSBAKeychain  FormFactor = 0x01
	FormFactorUSBANano      FormFactor = 0x02
	FormFactorUSBCKeychain  FormFactor = 0x03
	FormFactorUSBCNano      FormFactor = 0x04
	FormFactorUSBCLightning FormFactor = 0x05
	FormFactorUSBABio       FormFactor = 0x06
	FormFactorUSBCBio       FormFactor = 0x07
)

type DeviceInfo struct {
	Flags            DeviceFlag
	CapsSupportedUSB Capability
	CapsEnabledUSB   Capability
	CapsSupportedNFC Capability
	CapsEnabledNFC   Capability
	SerialNumber     uint32
	FirmwareVersion  iso.Version
	FormFactor       FormFactor
	AutoEjectTimeout time.Duration
	ChalRespTimeout  time.Duration
	IsLocked         bool
	IsSky            bool
	IsFIPS           bool
}

// nolint: gocognit
func (di *DeviceInfo) Unmarshal(b []byte) error {
	tvs, err := tlv.DecodeSimple(b[1:])
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	readCap := func(tv tlv.TagValue) (Capability, error) {
		switch len(tv.Value) {
		case 1: // For Yubikey 4.x
			return Capability(tv.Value[0]), nil
		case 2:
			return Capability(binary.BigEndian.Uint16(tv.Value)), nil
		default:
			return 0, ErrInvalidResponseLength
		}
	}

	for _, tv := range tvs {
		switch tv.Tag {
		case TagCapsSupportedUSB:
			if di.CapsSupportedNFC, err = readCap(tv); err != nil {
				return fmt.Errorf("%w: CapsSupportedUSB", err)
			}

		case TagCapsEnabledUSB:
			if di.CapsEnabledUSB, err = readCap(tv); err != nil {
				return fmt.Errorf("%w: CapsEnabledUSB", err)
			}

		case TagCapsSupportedNFC:
			if len(tv.Value) != 2 {
				return fmt.Errorf("%w: CapsSupportedNFC", ErrInvalidResponseLength)
			}

			di.CapsSupportedNFC = Capability(binary.BigEndian.Uint16(tv.Value))

		case TagCapsEnabledNFC:
			if len(tv.Value) != 2 {
				return fmt.Errorf("%w: CapsEnabledNFC", ErrInvalidResponseLength)
			}

			di.CapsEnabledNFC = Capability(binary.BigEndian.Uint16(tv.Value))

		case TagSerialNumber:
			if len(tv.Value) != 4 {
				return fmt.Errorf("%w: SerialNumber", ErrInvalidResponseLength)
			}

			di.SerialNumber = binary.BigEndian.Uint32(tv.Value)

		case TagFormFactor:
			if len(tv.Value) != 1 {
				return fmt.Errorf("%w: FormFactor", ErrInvalidResponseLength)
			}

			di.FormFactor = FormFactor(tv.Value[0] & 0xf)
			di.IsFIPS = tv.Value[0]&0x80 != 0
			di.IsSky = tv.Value[0]&0x40 != 0

		case TagFirmwareVersion:
			if len(tv.Value) != 3 {
				return fmt.Errorf("%w: FirmwareVersion", ErrInvalidResponseLength)
			}

			di.FirmwareVersion = iso.Version{
				Major: int(tv.Value[0]),
				Minor: int(tv.Value[1]),
				Patch: int(tv.Value[2]),
			}

		case TagAutoEjectTimeout:
			if len(tv.Value) != 2 {
				return fmt.Errorf("%w: AutoEjectTimeout", ErrInvalidResponseLength)
			}

			di.AutoEjectTimeout = time.Second * time.Duration(binary.BigEndian.Uint16(tv.Value))

		case TagChalRespTimeout:
			if len(tv.Value) != 1 {
				return fmt.Errorf("%w: ChalRespTimeout", ErrInvalidResponseLength)
			}

			di.ChalRespTimeout = time.Second * time.Duration(tv.Value[0])

		case TagDeviceFlags:
			if len(tv.Value) != 1 {
				return fmt.Errorf("%w: DeviceFlags", ErrInvalidResponseLength)
			}

			di.Flags = DeviceFlag(tv.Value[0])

		case TagAppVersions:
			// TODO

		case TagConfigLock:
			if len(tv.Value) != 1 {
				return fmt.Errorf("%w: ConfigLock", ErrInvalidResponseLength)
			}

			di.IsLocked = tv.Value[0] != 0
		}
	}

	return nil
}

// GetDeviceInfo returns device information about the YubiKey token.
func GetDeviceInfo(card *iso.Card) (*DeviceInfo, error) {
	resp, err := card.Send(&iso.CAPDU{
		Ins: 0x1D,
		P1:  0x00,
		P2:  0x00,
	})
	if err != nil {
		return nil, err
	}

	di := &DeviceInfo{}
	if err := di.Unmarshal(resp); err != nil {
		return nil, err
	}

	return di, nil
}
