// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package filter

import (
	"errors"
	"fmt"

	iso "cunicu.li/go-iso7816"
)

// ErrOpen is a sentinel error returned by a Filter indicating
// that the filter requires a connection to the card for checking
// its predicate.
// This is usually the case if the predicate needs to exchange
// APDUs with the card rather than simply checking the readers
// name
var ErrOpen = errors.New("open card for detailed filtering")

// Filter is a predicate which evaluates
// whether a given reader/card matches
// a given condition.
type Filter func(card iso.PCSCCard) (bool, error)

// Any matches any card
func Any(_ iso.PCSCCard) (bool, error) {
	return true, nil
}

// None matches no card
func None(_ iso.PCSCCard) (bool, error) {
	return false, nil
}

// HasApplet matches card which can select an applet
// with the given application identifier (AID).
func HasApplet(aid []byte) Filter {
	return func(card iso.PCSCCard) (bool, error) {
		if card == nil {
			return false, ErrOpen
		}

		isoCard := iso.NewCard(card)
		tx, err := isoCard.NewTransaction()
		if err != nil {
			return false, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Close()

		if _, err := tx.Select(aid); err != nil {
			return false, nil //nolint:nilerr
		}

		return true, nil
	}
}
