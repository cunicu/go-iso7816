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
	return withApplet(iso.AidYubicoOTP, func(card *iso.Card) (bool, error) {
		if sts, err := GetStatus(card); err != nil {
			return false, err
		} else if !sts.Version.Less(v) {
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
func HasOTP(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapOTP)(reader, card)
}

// HasU2F is a filter which checks if the YubiKey has the U2F
// applet enabled.
func HasU2F(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapU2F)(reader, card)
}

// HasFIDO2 is a filter which checks if the YubiKey has the FIDO2
// applet enabled.
func HasFIDO2(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapFIDO2)(reader, card)
}

// HasOATH is a filter which checks if the YubiKey has the OATH
// applet enabled.
func HasOATH(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapOATH)(reader, card)
}

// HasPIV is a filter which checks if the YubiKey has the PIV
// applet enabled.
func HasPIV(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapPIV)(reader, card)
}

// HasOpenPGP is a filter which checks if the YubiKey has the OpenPGP
// applet enabled.
func HasOpenPGP(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapOpenPGP)(reader, card)
}

// HasHSMAuth is a filter which checks if the YubiKey has the HSM authentication
// applet enabled.
func HasHSMAuth(reader string, card *iso.Card) (bool, error) {
	return hasCapabilityEnabled(CapOpenPGP)(reader, card)
}

func hasCapabilityEnabled(c Capability) filter.Filter {
	return withDeviceInfo(func(di *DeviceInfo) bool {
		return (di.CapsEnabledUSB|di.CapsEnabledNFC)&c != 0
	})
}

func withDeviceInfo(cb func(di *DeviceInfo) bool) filter.Filter {
	return withApplet(iso.AidYubicoManagement, func(card *iso.Card) (bool, error) {
		di, err := GetDeviceInfo(card)
		if err != nil {
			return false, fmt.Errorf("failed to get device information: %w", err)
		}

		return cb(di), nil
	})
}

func withApplet(aid []byte, cb func(card *iso.Card) (bool, error)) filter.Filter {
	return func(reader string, card *iso.Card) (bool, error) {
		// Matching against the name first saves us from connecting to the card
		if match, err := filter.IsYubiKey(reader, card); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}

		if card == nil {
			return false, filter.ErrOpen
		}

		if _, err := card.Select(aid); err != nil {
			return false, nil //nolint:nilerr
		}

		return cb(card)
	}
}
