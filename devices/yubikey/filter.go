// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package yubikey

import (
	"fmt"

	iso "cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/filter"
)

func HasVersionStr(s string) filter.Filter {
	if v, err := iso.ParseVersion(s); err == nil {
		return HasVersion(v)
	}

	return filter.None
}

// HasVersion checks that the card has a firmware version equal or higher
// than the given one.
func HasVersion(v iso.Version) filter.Filter {
	return withApplet(iso.AidYubicoOTP, func(card iso.PCSCCard) (bool, error) {
		ykCard := Card{iso.NewCard(card)}
		if sts, err := ykCard.Status(); err != nil {
			return false, err
		} else if v.Less(sts.Version) {
			return false, nil
		}

		return true, nil
	})
}

func IsSerialNumber(sno uint32) filter.Filter {
	return withDeviceInfo(func(di *DeviceInfo) bool {
		return di.SerialNumber == sno
	})
}

// HasFormFactor returns a filter which checks if the YubiKey
// has a given form factor.
func HasFormFactor(ff FormFactor) filter.Filter {
	return withDeviceInfo(func(di *DeviceInfo) bool {
		return di.FormFactor == ff
	})
}

//nolint:gochecknoglobals
var (
	IsFIPS = withDeviceInfo(func(di *DeviceInfo) bool {
		return di.IsFIPS
	})
	IsLocked = withDeviceInfo(func(di *DeviceInfo) bool {
		return di.IsLocked
	})
)

// HasOTP is a filter which checks if the YubiKey has the OTP
// applet enabled.
func HasOTP(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapOTP)(card)
}

// HasU2F is a filter which checks if the YubiKey has the U2F
// applet enabled.
func HasU2F(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapU2F)(card)
}

// HasFIDO2 is a filter which checks if the YubiKey has the FIDO2
// applet enabled.
func HasFIDO2(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapFIDO2)(card)
}

// HasOATH is a filter which checks if the YubiKey has the OATH
// applet enabled.
func HasOATH(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapOATH)(card)
}

// HasPIV is a filter which checks if the YubiKey has the PIV
// applet enabled.
func HasPIV(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapPIV)(card)
}

// HasOpenPGP is a filter which checks if the YubiKey has the OpenPGP
// applet enabled.
func HasOpenPGP(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapOpenPGP)(card)
}

// HasHSMAuth is a filter which checks if the YubiKey has the HSM authentication
// applet enabled.
func HasHSMAuth(card iso.PCSCCard) (bool, error) {
	return hasCapabilityEnabled(CapOpenPGP)(card)
}

func hasCapabilityEnabled(c Capability) filter.Filter {
	return withDeviceInfo(func(di *DeviceInfo) bool {
		return (di.CapsEnabledUSB|di.CapsEnabledNFC)&c != 0
	})
}

func withDeviceInfo(cb func(di *DeviceInfo) bool) filter.Filter {
	return withApplet(iso.AidYubicoManagement, func(card iso.PCSCCard) (bool, error) {
		ykCard := Card{iso.NewCard(card)}
		di, err := ykCard.DeviceInfo()
		if err != nil {
			return false, fmt.Errorf("failed to get device information: %w", err)
		}

		return cb(di), nil
	})
}

func withApplet(aid []byte, cb func(card iso.PCSCCard) (bool, error)) filter.Filter {
	return func(card iso.PCSCCard) (bool, error) {
		// Matching against the name first saves us from connecting to the card
		if match, err := filter.IsYubiKey(card); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}

		if card == nil {
			return false, filter.ErrOpen
		}

		isoCard := iso.NewCard(card)
		if _, err := isoCard.Select(aid); err != nil {
			return false, nil //nolint:nilerr
		}

		ykCard := &Card{iso.NewCard(card)}

		return cb(ykCard)
	}
}
