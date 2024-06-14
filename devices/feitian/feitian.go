// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package feitian

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	iso "cunicu.li/go-iso7816"
)

type Card struct {
	*iso.Card
}

func NewCard(card iso.PCSCCard) *Card {
	isoCard := iso.NewCard(card)
	return &Card{isoCard}
}

// Transmit sends an APDU to the card and returns the response.
// This is a custom version of iso7816.Card.Transmit()
// as FEITIAN has a broken ISO-7816-4 implementation in their tokens.
func (c *Card) Transmit(cmd *iso.CAPDU) ([]byte, error) {
	cmdBuf := []byte{
		cmd.Cla,
		byte(cmd.Ins),
		cmd.P1,
		cmd.P2,
	}

	if cmd.Data == nil {
		cmdBuf = append(cmdBuf, byte(cmd.Ne))
	}

	if lc := len(cmd.Data); lc <= 0xFF {
		cmdBuf = append(cmdBuf, byte(lc))
	} else {
		panic("unsupported command length")
	}

	cmdBuf = append(cmdBuf, cmd.Data...)

	return c.Card.Transmit(cmdBuf)
}

func printable(b []byte) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII || r < 31 {
			return -1
		}
		return r
	}, string(b))
}

func (c *Card) SerialNumber() (string, error) {
	resp, err := c.Transmit(&iso.CAPDU{
		Cla: 128,
		Ins: 227,
		P1:  3,
		P2:  0,
	})
	if err != nil {
		return "", err
	}

	if len(resp) == 8 {
		return hex.EncodeToString(resp), nil
	}

	return printable(resp), nil
}

func (c *Card) COSVersion() (string, error) {
	if resp, err := c.Transmit(&iso.CAPDU{
		Cla: 128,
		Ins: 227,
		P1:  0,
		P2:  0,
		Ne:  3,
	}); err == nil {
		if len(resp) >= 3 {
			if resp[0]>>4 == 0 {
				return fmt.Sprintf("%x%x%02X", resp[0], resp[1]&15, resp[2]), nil
			}
			return fmt.Sprintf("%x%x%02X", resp[0]>>4, resp[0]&15, resp[1]), nil
		}
	}

	if _, err := c.Transmit(&iso.CAPDU{
		Cla: 0,
		Ins: 164,
		P1:  4,
		P2:  0,
	}); err != nil {
		return "", err
	}

	if resp, err := c.Transmit(&iso.CAPDU{
		Cla: 0,
		Ins: 202,
		P1:  159,
		P2:  127,
	}); err == nil {
		if len(resp) > 20 {
			if resp[8] == 1 {
				return fmt.Sprintf("%02X%02X", resp[8], resp[9]), nil
			} else if resp[8] == 0 {
				return fmt.Sprintf("%x%x%02X", resp[9], resp[10], resp[11]), nil
			}

			// if len(self.serial_number) < 8 {
			// 	sno := hex.EncodeToString(resp[12:16])
			// }
		}
	}

	return "", nil
}
