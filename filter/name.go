// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package filter

import (
	"regexp"

	iso "cunicu.li/go-iso7816"
)

// HasName compares the name of the smart card reader
// with the provided name.
func HasName(nameExpected string) Filter {
	return func(card iso.PCSCCard) (bool, error) {
		if card, ok := card.(iso.ReaderCard); ok {
			return card.Reader() == nameExpected, nil
		}

		return false, nil
	}
}

// HasName matches the name of the smart card reader
// against the provided regular expression.
func HasNameRegex(regex string) Filter {
	re := regexp.MustCompile(regex)
	return func(card iso.PCSCCard) (bool, error) {
		if card, ok := card.(iso.ReaderCard); ok {
			return re.MatchString(card.Reader()), nil
		}

		return false, nil
	}
}

// IsYubiKey checks if the smart card is a YubiKey
// based on the name of the smart card reader.
func IsYubiKey(card iso.PCSCCard) (bool, error) {
	return HasNameRegex("(?i)YubiKey")(card)
}

// IsNikrokey checks if the smart card is a Nitrokey
// based on the name of the smart card reader.
func IsNitrokey(card iso.PCSCCard) (bool, error) {
	return HasNameRegex("(?i)Nitrokey")(card)
}

// IsNikrokey3 checks if the smart card is a Nitrokey 3
// based on the name of the smart card reader.
func IsNitrokey3(card iso.PCSCCard) (bool, error) {
	return HasNameRegex("Nitrokey 3")(card)
}

// IsFeitian checks if the smart card is a FEITIAN key
// based on the name of the smart card reader.
func IsFeitian(card iso.PCSCCard) (bool, error) {
	return HasNameRegex("^FT")(card)
}
