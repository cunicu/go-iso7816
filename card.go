// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816

import (
	"fmt"
)

type ReconnectableCard interface {
	Reconnect(reset bool) error
}

type ReaderCard interface {
	Reader() string
}

type MetadataCard interface {
	Metadata() map[string]string
}

type PCSCCard interface {
	Transmit([]byte) ([]byte, error)
	BeginTransaction() error
	EndTransaction() error
	Close() error
	Base() PCSCCard
}

type Card struct {
	PCSCCard

	InsGetRemaining Instruction
}

func NewCard(c PCSCCard) *Card {
	return &Card{
		PCSCCard: c,

		// Some applets like Yubico's OATH applet use a different
		// command for fetching remaining data
		InsGetRemaining: InsGetResponse,
	}
}

func (c *Card) Select(aid []byte) (respBuf []byte, err error) {
	return c.Send(&CAPDU{
		Ins:  InsSelect,
		P1:   0x04,
		P2:   0x00,
		Data: aid,
		Ne:   MaxLenRespDataStandard,
	})
}

// If checks that the card meets the provided filter condition.
func (c *Card) If(flt func(card PCSCCard) (bool, error)) (bool, error) {
	return flt(c)
}

// Send sends a command APDU to the card
// nolint: unparam
func (c *Card) Send(cmd *CAPDU) (respBuf []byte, err error) {
	cmdBuf, err := cmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize CAPDU: %w", err)
	}

	for {
		r, err := c.Transmit(cmdBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to transmit CAPDU: %w", err)
		}

		resp, err := ParseRAPDU(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RAPDU: %w", err)
		}

		respCode := resp.Code()
		respBuf = append(respBuf, resp.Data...)

		switch {
		case respCode.HasMore():
			cmdBuf = []byte{0x00, byte(c.InsGetRemaining), 0x00, 0x00}

		case respCode.IsSuccess():
			return respBuf, nil

		default:
			return nil, respCode
		}
	}
}

type Transaction struct {
	*Card
}

func (c *Card) NewTransaction() (*Transaction, error) {
	if err := c.BeginTransaction(); err != nil {
		return nil, err
	}

	return &Transaction{c}, nil
}

func (tx *Transaction) Close() error {
	return tx.EndTransaction()
}
