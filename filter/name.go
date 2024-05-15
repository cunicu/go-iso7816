// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package filter

import (
	"regexp"

	iso "cunicu.li/go-iso7816"
)

// HasName compares the name of the smart card reader
// with the provided name.
func HasName(nameExpected string) Filter {
	return func(reader string, _ *iso.Card) (bool, error) {
		return reader == nameExpected, nil
	}
}

// HasName matches the name of the smart card reader
// against the provided regular expression.
func HasNameRegex(regex string) Filter {
	re := regexp.MustCompile(regex)
	return func(reader string, _ *iso.Card) (bool, error) {
		return re.MatchString(reader), nil
	}
}

// IsYubikey checks if the smart card is a YubiKey
// based on the name of the smart card reader.
func IsYubiKey(reader string, card *iso.Card) (bool, error) {
	return HasNameRegex("(?i)YubiKey")(reader, card)
}

// IsNikrokey checks if the smart card is a Nitrokey
// based on the name of the smart card reader.
func IsNitrokey(reader string, card *iso.Card) (bool, error) {
	return HasNameRegex("(?i)Nitrokey")(reader, card)
}

// IsNikrokey3 checks if the smart card is a Nitrokey 3
// based on the name of the smart card reader.
func IsNitrokey3(reader string, card *iso.Card) (bool, error) {
	return HasNameRegex("Nitrokey 3")(reader, card)
}
