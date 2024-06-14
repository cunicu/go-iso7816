// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// Package yubikey implements basic support for getting status and details about YubiKey tokens.
package yubikey

import (
	"errors"

	iso "cunicu.li/go-iso7816"
)

var ErrInvalidResponseLength = errors.New("invalid response length")

const (
	// https://docs.yubico.com/yesdk/users-manual/application-otp/otp-commands.html
	InsOTP        iso.Instruction = 0x01 // Most commands of the OTP applet use this value
	InsReadStatus iso.Instruction = 0x03
)

type Card struct {
	*iso.Card
}
