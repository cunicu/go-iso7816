// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pcsc

import (
	"bytes"
	"fmt"

	"github.com/ebfe/scard"

	"cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/filter"
)

// HasAttribute uses the PCSC API SCardGetAttrib to
// check if the card has given attribute
func HasAttribute(attr scard.Attrib, value []byte) filter.Filter {
	return func(name string, c *iso7816.Card) (bool, error) {
		if c == nil {
			return false, filter.ErrOpen
		}

		sc, ok := c.PCSCCard.(*Card)
		if !ok {
			return false, nil
		}

		val, err := sc.GetAttrib(attr)
		if err != nil {
			return false, fmt.Errorf("failed to get attribute: %w", err)
		}

		return bytes.Equal(val, value), nil
	}
}