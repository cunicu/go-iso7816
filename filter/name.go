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
	return func(name string, c *iso.Card) (bool, error) {
		return name == nameExpected, nil
	}
}

// HasName matches the name of the smart card reader
// against the provided regular expression.
func HasNameRegex(regex string) Filter {
	re := regexp.MustCompile(regex)
	return func(name string, c *iso.Card) (bool, error) {
		return re.MatchString(name), nil
	}
}

// IsYubikey checks if the smart card is a YubiKey
// based on the name of the smart card reader.
func IsYubiKey(n string, c *iso.Card) (bool, error) {
	return HasNameRegex("(?i)YubiKey")(n, c)
}

// IsNikrokey checks if the smart card is a Nitrokey
// based on the name of the smart card reader.
func IsNitrokey(n string, c *iso.Card) (bool, error) {
	return HasNameRegex("(?i)Nitrokey")(n, c)
}

// IsNikrokey3 checks if the smart card is a Nitrokey 3
// based on the name of the smart card reader.
func IsNitrokey3(n string, c *iso.Card) (bool, error) {
	return HasNameRegex("Nitrokey 3")(n, c)
}
