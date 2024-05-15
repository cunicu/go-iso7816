// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// Package nitrokey implements basic support for getting status and details about Nitrokey 3 tokens.
package nitrokey

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ebfe/scard"

	iso "cunicu.li/go-iso7816"
)

var ErrInvalidLength = errors.New("invalid length")

const (
	LenRandom  = 57
	LenVersion = 4
)

const (
	// https://github.com/Nitrokey/admin-app/blob/main/src/admin.rs
	InsUpdate             iso.Instruction = 0x51
	InsReboot             iso.Instruction = 0x53
	InsRNG                iso.Instruction = 0x60
	InsGetFirmwareVersion iso.Instruction = 0x61
	InsGetUUID            iso.Instruction = 0x62
	InsLocked             iso.Instruction = 0x63

	InsAdminStatus          byte = 0x80
	InsAdminTestSE050       byte = 0x81
	InsAdminGetConfig       byte = 0x82
	InsAdminSetConfig       byte = 0x83
	InsAdminFactoryReset    byte = 0x84
	InsAdminFactoryResetApp byte = 0x85
)

// https://github.com/Nitrokey/pynitrokey/blob/781d4b9e3e9fc3cfc297611d31e7e43643547ac8/pynitrokey/nk3/admin_app.py#L20
type InitStatus byte

const (
	InitStatusNFCError           InitStatus = 0b0001
	InitStatusInternalFlashError InitStatus = 0b0010
	InitStatusExternalFlashError InitStatus = 0b0100
	InitStatusMigrationError     InitStatus = 0b1000
)

// https://github.com/Nitrokey/pynitrokey/blob/781d4b9e3e9fc3cfc297611d31e7e43643547ac8/pynitrokey/nk3/admin_app.py#L41
type Variant byte

const (
	VariantUSBIP Variant = 0
	VariantLPC55 Variant = 1
	VariantNRF52 Variant = 2
)

// https://github.com/Nitrokey/pynitrokey/blob/781d4b9e3e9fc3cfc297611d31e7e43643547ac8/pynitrokey/nk3/admin_app.py#L77
type DeviceStatus struct {
	InitStatus InitStatus
	IfsBlocks  byte
	EfsBlocks  uint16
	Variant    Variant
}

func (ds *DeviceStatus) Unmarshal(b []byte) error {
	if len(b) < 1 {
		return ErrInvalidLength
	}

	if len(b) >= 4 {
		ds.IfsBlocks = b[1]
		ds.EfsBlocks = binary.BigEndian.Uint16(b[2:])
	}

	if len(b) >= 5 {
		ds.Variant = Variant(b[4])
	}

	return nil
}

// GetDeviceStatus returns the device status of the Nitrokey 3 token.
func GetDeviceStatus(card *iso.Card) (*DeviceStatus, error) {
	resp, err := card.Send(&iso.CAPDU{
		Ins:  iso.Instruction(InsAdminStatus),
		P1:   0x00,
		P2:   0x00,
		Data: []byte{InsAdminStatus},
		Ne:   0x05,
	})
	if err != nil {
		return nil, err
	}

	ds := &DeviceStatus{}
	if err := ds.Unmarshal(resp); err != nil {
		return nil, err
	}

	return ds, nil
}

// GetUUID returns the UUID of the Nitrokey 3 token.
func GetUUID(card *iso.Card) ([]byte, error) {
	resp, err := card.Send(&iso.CAPDU{
		Ins: InsGetUUID,
		P1:  0x00,
		P2:  0x00,
		Ne:  0x10,
	})
	if err != nil {
		return nil, err
	}

	if len(resp) != 16 {
		return nil, iso.ErrWrongLength
	}

	return resp, nil
}

// GetFirmwareVersion returns the firmware version of the Nitrokey 3 token.
func GetFirmwareVersion(card *iso.Card) (*iso.Version, error) {
	resp, err := card.Send(&iso.CAPDU{
		Ins:  InsGetFirmwareVersion,
		P1:   0x00,
		P2:   0x00,
		Data: []byte{0x01},
		Ne:   LenVersion,
	})
	if err != nil {
		return nil, err
	}

	if len(resp) != 4 {
		return nil, ErrInvalidLength
	}

	version := binary.BigEndian.Uint32(resp)

	// This is the reverse of the calculation in runners/lpc55/build.rs (CARGO_PKG_VERSION):
	// https://github.com/Nitrokey/nitrokey-3-firmware/blob/main/runners/lpc55/build.rs#L131
	return &iso.Version{
		Major: int(version >> 22),
		Minor: int((version >> 6) & ((1 << 16) - 1)),
		Patch: int(version & ((1 << 6) - 1)),
	}, nil
}

// GetRandom returns LenRandom bytes of random bytes.
func GetRandom(card *iso.Card) ([]byte, error) {
	return card.Send(&iso.CAPDU{
		Ins: InsRNG,
		Ne:  LenRandom,
	})
}

// Reboot resets the token.
func Reboot(card *iso.Card) error {
	_, err := card.Send(&iso.CAPDU{Ins: InsReboot})
	if !errors.Is(err, scard.ErrReaderUnavailable) {
		return fmt.Errorf("unexpected error: %w", err)
	}

	// Rebooting, detaches the card reader briefly which
	// requires re-enumerating/reconnecting.
	if pcscCard := card.Base(); pcscCard != nil {
		if rcard, ok := pcscCard.(iso.ReconnectableCard); ok {
			if err := rcard.Reconnect(false); err != nil {
				return err
			}
		}
	}

	return nil
}

// IsLocked checks if the bootloader is locked.
// Locked bootloaders can only be updated via officially signed firmware images.
func IsLocked(card *iso.Card) (bool, error) {
	resp, err := card.Send(&iso.CAPDU{
		Ins: InsLocked,
		Ne:  1,
	})
	if err != nil {
		return false, err
	} else if len(resp) != 1 {
		return false, iso.ErrWrongLength
	}

	return resp[0] == 1, nil
}

func Metadata(card *iso.Card) (meta map[string]any) {
	if _, err := card.Select(iso.AidSolokeysAdmin); err != nil {
		return nil
	}

	meta = map[string]any{}

	if v, err := GetFirmwareVersion(card); err == nil {
		meta["version"] = v
	}

	if id, err := GetUUID(card); err == nil {
		meta["uuid"] = id
	}

	return meta
}
