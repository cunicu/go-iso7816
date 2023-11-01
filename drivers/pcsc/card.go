// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pcsc

import (
	"fmt"

	"github.com/ebfe/scard"

	iso "cunicu.li/go-iso7816"
)

var _ iso.PCSCCard = (*Card)(nil)

// Card implements the iso7816.PCSCCard interface
// via github.com/ebfe/scard.
type Card struct {
	*scard.Card
}

// NewCard creates a new card by connecting via the PC/SC API.
func NewCard(ctx *scard.Context, reader string) (*iso.Card, error) {
	sc, err := ctx.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to reader: %w", err)
	}

	return iso.NewCard(&Card{sc}), nil
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
