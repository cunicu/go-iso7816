// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pcsc

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ebfe/scard"

	iso "cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/filter"
)

var ErrNoCardFound = errors.New("no card found")

// OpenCards opens up to cnt cards which match the provided filter flt.
func OpenCards(ctx *scard.Context, cnt int, flt filter.Filter, shared bool) (cards []iso.PCSCCard, err error) {
	readers, err := ctx.ListReaders()
	if err != nil {
		return nil, fmt.Errorf("failed to list readers: %w", err)
	}

	// Make the list of returned cards deterministic
	slices.Sort(readers)

	for _, reader := range readers {
		var card *iso.Card
		for i := 0; i < 2; i++ {
			if match, err := flt(reader, card); errors.Is(err, filter.ErrOpen) || match {
				if card, err = NewCard(ctx, reader, shared); err != nil {
					return nil, fmt.Errorf("failed to connect to card: %w", err)
				}

				if match {
					cards = append(cards, card.PCSCCard)
					break
				}
			} else if err != nil {
				return nil, err
			}
		}

		if cnt >= 0 && len(cards) >= cnt {
			break
		}
	}

	return cards, nil
}

// OpenFirstCard opens the first card which matches the filter flt
// or returns ErrNoCardFound if none was found.
func OpenFirstCard(ctx *scard.Context, flt filter.Filter, shared bool) (iso.PCSCCard, error) {
	cards, err := OpenCards(ctx, 1, flt, shared)
	if err != nil {
		return nil, err
	} else if len(cards) != 1 {
		return nil, ErrNoCardFound
	}

	return cards[0], nil
}
