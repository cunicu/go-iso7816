package feitan

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	iso "cunicu.li/go-iso7816"
)

// transmit sends an APDU to the card and returns the response.
// This is a custom version of iso7816.Card.Transmit()
// as FEITIAN has a broken ISO-7816-4 implementation in their tokens.
func transmit(c *iso.Card, cmd *iso.CAPDU) ([]byte, error) {
	cmdBuf := []byte{
		cmd.Cla,
		byte(cmd.Ins),
		cmd.P1,
		cmd.P2,
	}

	if cmd.Data == nil {
		cmdBuf = append(cmdBuf, byte(cmd.Ne))
	}

	if lc := len(cmd.Data); lc <= 0xff {
		cmdBuf = append(cmdBuf, byte(lc))
	} else {
		panic("unsupported command length")
	}

	cmdBuf = append(cmdBuf, cmd.Data...)

	return c.Transmit(cmdBuf)
}

func printable(b []byte) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII || r < 31 {
			return -1
		}
		return r
	}, string(b))
}

func GetSerialNumber(c *iso.Card) (string, error) {
	resp, err := transmit(c, &iso.CAPDU{
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

func GetAppletVersion(c *iso.Card) (string, error) {
	if _, err := c.Select(iso.AidFIDO); err != nil {
		return "", err
	}

	resp, err := transmit(c, &iso.CAPDU{
		Cla:  128,
		Ins:  226,
		P1:   128,
		P2:   0,
		Data: []byte{0xDF, 0xFF, 0x02, 0x80, 0x01},
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%02x%02x", resp[8], resp[9]), nil
}
