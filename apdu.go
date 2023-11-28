// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-FileCopyrightText: 2020 skythen
// SPDX-License-Identifier: Apache-2.0

package iso7816

import (
	"errors"
	"fmt"
)

var errInvalidLength = errors.New("invalid length")

const (
	LenHeader          = 4 // LenHeader defines the length of the header of an APDU.
	LenLCStandard      = 1 // LenLCStandard defines the length of the LC of a standard length APDU.
	LenLCExtended      = 3 // LenLCExtended defines the length of the LC of an extended length APDU.
	LenResponseTrailer = 2 // LenResponseTrailer defines the length of the trailer of a RAPDU.

	MaxLenCommandDataStandard  = (1 << 8) - 1  // MaxLenCommandDataStandard defines the maximum command data length of a standard length CAPDU.
	MaxLenResponseDataStandard = (1 << 8)      // MaxLenResponseDataStandard defines the maximum response data length of a standard length RAPDU.
	MaxLenCommandDataExtended  = (1 << 16) - 1 // MaxLenCommandDataExtended defines the maximum command data length of an extended length CAPDU.
	MaxLenResponseDataExtended = (1 << 16)     // MaxLenResponseDataExtended defines the maximum response data length of an extended length RAPDU.
)

type CAPDU struct {
	Cla    byte        // Cla is the class byte.
	Ins    Instruction // Ins is the instruction byte.
	P1, P2 byte        // P1, P2 is the p1, p2 byte.
	Data   []byte      // Data is the data field.
	Ne     int         // Ne is the total number of expected response data byte (not LE encoded).
}

func (c *CAPDU) Bytes() ([]byte, error) {
	if len(c.Data) > MaxLenCommandDataExtended {
		return nil, fmt.Errorf("%w:CAPDU data length %d exceeds maximum allowed length of %d",
			errInvalidLength, len(c.Data), MaxLenCommandDataExtended)
	}

	if c.Ne > MaxLenResponseDataExtended {
		return nil, fmt.Errorf("%w: expected length ne %d exceeds maximum allowed length of %d",
			errInvalidLength, len(c.Data), MaxLenResponseDataExtended)
	}

	switch {
	case len(c.Data) == 0 && c.Ne == 0: // Case 1: Cla | Ins | P1 | P2
		return []byte{c.Cla, byte(c.Ins), c.P1, c.P2}, nil

	case len(c.Data) == 0 && c.Ne > 0: // Case 2
		// Extended format: Cla | Ins | P1 | P2 | Le (extended)
		if c.Ne > MaxLenResponseDataStandard {
			le := make([]byte, LenLCExtended) // First byte is zero byte, so LE length is equal to LC length

			if c.Ne == MaxLenResponseDataExtended {
				le[1] = 0x00
				le[2] = 0x00
			} else {
				le[1] = (byte)((c.Ne >> 8) & 0xFF)
				le[2] = (byte)(c.Ne & 0xFF)
			}

			result := make([]byte, 0, LenHeader+LenLCExtended)
			result = append(result, c.Cla, byte(c.Ins), c.P1, c.P2)
			result = append(result, le...)

			return result, nil
		}

		// Standard format: Cla | Ins | P1 | P2 | Le
		result := make([]byte, 0, LenHeader+LenLCStandard)
		result = append(result, c.Cla, byte(c.Ins), c.P1, c.P2)

		if c.Ne == MaxLenResponseDataStandard {
			result = append(result, 0x00)
		} else {
			result = append(result, byte(c.Ne))
		}

		return result, nil

	case len(c.Data) != 0 && c.Ne == 0: // Case 3
		// Extended format: Cla | Ins | P1 | P2 | Lc (extended) | Data
		if len(c.Data) > MaxLenCommandDataStandard {
			lc := make([]byte, LenLCExtended)
			lc[1] = (byte)((len(c.Data) >> 8) & 0xFF)
			lc[2] = (byte)(len(c.Data) & 0xFF)

			result := make([]byte, 0, LenHeader+LenLCExtended+len(c.Data))
			result = append(result, c.Cla, byte(c.Ins), c.P1, c.P2)
			result = append(result, lc...)
			result = append(result, c.Data...)

			return result, nil
		}

		// Standard format: Cla | Ins | P1 | P2 | Lc | Data
		result := make([]byte, 0, LenHeader+1+len(c.Data))
		result = append(result, c.Cla, byte(c.Ins), c.P1, c.P2, byte(len(c.Data)))
		result = append(result, c.Data...)

		return result, nil

	case c.Ne > MaxLenResponseDataStandard || len(c.Data) > MaxLenCommandDataStandard: // Case 4: Cla | Ins | P1 | P2 | Lc (extended) | Data | Le (extended)
		lc := make([]byte, LenLCExtended) // First byte is zero byte
		lc[1] = (byte)((len(c.Data) >> 8) & 0xFF)
		lc[2] = (byte)(len(c.Data) & 0xFF)

		le := make([]byte, 2)

		if c.Ne == MaxLenResponseDataExtended {
			le[0] = 0x00
			le[1] = 0x00
		} else {
			le[0] = (byte)((c.Ne >> 8) & 0xFF)
			le[1] = (byte)(c.Ne & 0xFF)
		}

		result := make([]byte, 0, LenHeader+LenLCExtended+len(c.Data)+len(le))
		result = append(result, c.Cla, byte(c.Ins), c.P1, c.P2)
		result = append(result, lc...)
		result = append(result, c.Data...)
		result = append(result, le...)

		return result, nil

	default: // Standard format: Cla | Ins | P1 | P2 | Lc | Data | Ne
		result := make([]byte, 0, LenHeader+LenLCStandard+len(c.Data)+1)
		result = append(result, c.Cla, byte(c.Ins), c.P1, c.P2, byte(len(c.Data)))
		result = append(result, c.Data...)
		result = append(result, byte(c.Ne))

		return result, nil
	}
}

type RAPDU struct {
	Data     []byte
	SW1, SW2 byte
}

// ParseRAPDU parses a Response APDU and returns a RAPDU.
func ParseRAPDU(b []byte) (*RAPDU, error) {
	if len(b) < LenResponseTrailer || len(b) > 65538 {
		return nil, fmt.Errorf("%w: a RAPDU must consist of at least 2 byte and maximum of 65538 byte, got %d", errInvalidLength, len(b))
	}

	if len(b) == LenResponseTrailer {
		return &RAPDU{SW1: b[0], SW2: b[1]}, nil
	}

	return &RAPDU{Data: b[:len(b)-LenResponseTrailer], SW1: b[len(b)-2], SW2: b[len(b)-1]}, nil
}

// IsSuccess returns true if the RAPDU indicates the successful execution of a command ('0x61xx' or '0x9000'), otherwise false.
func (r *RAPDU) IsSuccess() bool {
	return r.SW1 == 0x61 || r.SW1 == 0x90 && r.SW2 == 0x00
}

// IsWarning returns true if the RAPDU indicates the execution of a command with a warning ('0x62xx' or '0x63xx'), otherwise false.
func (r *RAPDU) IsWarning() bool {
	return r.SW1 == 0x62 || r.SW1 == 0x63
}

// IsError returns true if the RAPDU indicates an error during the execution of a command ('0x64xx', '0x65xx' or from '0x67xx' to 0x6Fxx'), otherwise false.
func (r *RAPDU) IsError() bool {
	return (r.SW1 == 0x64 || r.SW1 == 0x65) || (r.SW1 >= 0x67 && r.SW1 <= 0x6F)
}

func (r *RAPDU) Code() Code {
	return Code{r.SW1, r.SW2}
}
