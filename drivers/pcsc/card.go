// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pcsc

import (
	"errors"
	"fmt"
	"time"

	"github.com/ebfe/scard"

	iso "cunicu.li/go-iso7816"
)

var _ iso.PCSCCard = (*Card)(nil)

// Card implements the iso7816.PCSCCard interface
// via github.com/ebfe/scard.
type Card struct {
	*scard.Card

	ctx    *scard.Context
	reader string
	mode   scard.ShareMode
}

// NewCard creates a new card by connecting via the PC/SC API.
func NewCard(ctx *scard.Context, reader string, shared bool) (*iso.Card, error) {
	mode := scard.ShareExclusive
	if shared {
		mode = scard.ShareShared
	}

	sc, err := ctx.Connect(reader, mode, scard.ProtocolAny)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to reader: %w", err)
	}

	return iso.NewCard(&Card{sc, ctx, reader, mode}), nil
}

func (c *Card) Base() iso.PCSCCard {
	return c
}

// Transmit wraps SCardTransmit.
func (c *Card) Transmit(cmd []byte) ([]byte, error) {
	return c.Card.Transmit(cmd)
}

// BeginTransaction wraps SCardBeginTransaction.
func (c *Card) BeginTransaction() error {
	return c.Card.BeginTransaction()
}

// EndTransaction wraps SCardEndTransaction.
func (c *Card) EndTransaction() error {
	return c.Card.EndTransaction(scard.LeaveCard)
}

// Close disconnects and resets the card
func (c *Card) Close() error {
	return c.Card.Disconnect(scard.ResetCard)
}

// Reset reconnects and resets the card
func (c *Card) Reconnect(reset bool) (err error) {
	if reset {
		return c.Card.Reconnect(scard.ShareShared, scard.ProtocolT1, scard.ResetCard)
	}

	for {
		if c.Card, err = c.ctx.Connect(c.reader, c.mode, scard.ProtocolAny); err == nil {
			return nil
		} else if errors.Is(err, scard.ErrUnknownReader) || errors.Is(err, scard.ErrNoSmartcard) {
			time.Sleep(100 * time.Millisecond)
		} else {
			return err
		}
	}
}
