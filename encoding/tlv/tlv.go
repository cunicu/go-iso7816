// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// Package tlv implements various version (ASN.1 BER, Simple, Compact) TLV encoding using in ISO 7816-4.
package tlv

import (
	"encoding"
	"errors"
)

var errInvalidLength = errors.New("invalid length")

type TagValue struct {
	Tag        Tag
	Value      []byte
	Children   []TagValue
	SkipLength bool
}

func New(t Tag, values ...any) (tv TagValue) {
	tv.Tag = t

	for _, value := range values {
		switch v := value.(type) {
		case encoding.BinaryMarshaler:
			data, err := v.MarshalBinary()
			if err != nil {
				panic("failed to marshal")
			}

			tv.Value = append(tv.Value, data...)

		case []byte:
			tv.Value = append(tv.Value, v...)

		case string:
			tv.Value = append(tv.Value, []byte(v)...)

		case TagValue:
			tv.Children = append(tv.Children, v)
		}
	}

	return tv
}
