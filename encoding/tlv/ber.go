// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv

import "errors"

var ErrNotConstructed = errors.New("tag is not constructed but contains children")

type Class byte

const (
	ClassUniversal   Class = 0b00
	ClassApplication Class = 0b01
	ClassContext     Class = 0b10
	ClassPrivate     Class = 0b11
)

func (t Tag) Class() Class {
	return Class((t & 0xff) >> 6)
}

func (t Tag) IsConstructed() bool {
	return t&(1<<5) != 0
}

func (tv TagValue) MarshalBER() (buf []byte, err error) {
	var cbuf []byte

	for _, child := range tv.Children {
		ccb, err := child.MarshalBER()
		if err != nil {
			return nil, err
		}

		cbuf = append(cbuf, ccb...)
	}

	tb, err := encodeTagBER(tv.Tag)
	if err != nil {
		panic("tag too large")
	}

	lb, err := encodeLengthBER(len(cbuf) + len(tv.Value))
	if err != nil {
		panic("tag too large")
	}

	buf = append(buf, tb...)
	buf = append(buf, lb...)
	if cbuf != nil {
		buf = append(buf, cbuf...)
	} else {
		buf = append(buf, tv.Value...)
	}

	return buf, nil
}

func (tv *TagValue) UnmarshalBER(buf []byte) ([]byte, error) {
	var err error

	t, buf, err := decodeTagBER(buf)
	if err != nil {
		return nil, err
	}

	l, buf, err := decodeLengthBER(buf)
	if err != nil {
		return nil, err
	}

	if len(buf) < l {
		return nil, errInvalidLength
	}

	tv.Tag = t
	tv.Value = buf[:l]

	if tv.Tag.IsConstructed() {
		if tv.Children, err = DecodeBER(tv.Value); err != nil {
			return nil, err
		}
	}

	return buf[l:], nil
}

func EncodeBER(tvs ...TagValue) (buf []byte, err error) {
	for _, tv := range tvs {
		vb, err := tv.MarshalBER()
		if err != nil {
			return nil, err
		}

		buf = append(buf, vb...)
	}

	return buf, nil
}

func DecodeBER(buf []byte) (tvs TagValues, err error) {
	for len(buf) > 0 {
		var tv TagValue
		if buf, err = tv.UnmarshalBER(buf); err != nil {
			return nil, err
		}

		tvs = append(tvs, tv)
	}

	return tvs, nil
}

// decodeTagBER decodes an ASN.1 BER-TLV encoded tag field.
// See: ISO 7816-4 Section 5.2.2.1 BER-TLV tag fields
func decodeTagBER(buf []byte) (t Tag, rbuf []byte, err error) {
	if len(buf) < 1 {
		return 0, nil, errInvalidLength
	}

	if buf[0]&0x1f != 0x1f {
		return Tag(buf[0]), buf[1:], nil
	}

	if len(buf) < 2 {
		return 0, nil, errInvalidLength
	} else if buf[1]&0x80 == 0 && buf[1]&0x7f > 30 {
		return Tag(uint32(buf[0])<<8 | uint32(buf[1])), buf[2:], nil
	}

	if len(buf) < 3 {
		return 0, nil, errInvalidLength
	} else if buf[1]&0x80 == 0x80 && buf[1]&0x7f != 0 {
		return Tag(uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])), buf[3:], nil
	}

	return 0, nil, nil
}

// decodeLengthBER decodes an  ASN.1 BER-TLV encoded length field.
// See: ISO 7816-4 Section 5.2.2.2 BER-TLV length fields
func decodeLengthBER(buf []byte) (int, []byte, error) {
	if len(buf) < 1 {
		return 0, nil, errInvalidLength
	}

	// Short form
	if buf[0] < 0x80 {
		return int(buf[0]), buf[1:], nil
	}

	// Long form
	n := int(buf[0] - 0x80)
	if n > 4 || len(buf) < n+1 {
		return -1, nil, errInvalidLength
	}

	l := 0
	for i := 1; i <= n; i++ {
		l <<= 8
		l |= int(buf[i])
	}

	return l, buf[n+1:], nil
}

func encodeTagBER(t Tag) (buf []byte, err error) {
	switch {
	case t>>8 == 0:
		buf = []byte{byte(t >> 0)}
	case t>>16 == 0:
		buf = []byte{byte(t >> 8), byte(t >> 0)}
	case t>>24 == 0:
		buf = []byte{byte(t >> 16), byte(t >> 8), byte(t >> 0)}
	case t>>32 == 0:
		buf = []byte{byte(t >> 24), byte(t >> 16), byte(t >> 8), byte(t >> 0)}
	default:
		return nil, errInvalidLength
	}

	return buf, nil
}

func encodeLengthBER(l int) ([]byte, error) {
	switch {
	case l < 0x80:
		return []byte{byte(l)}, nil
	case l>>8 == 0:
		return []byte{0x81, byte(l)}, nil
	case l>>16 == 0:
		return []byte{0x82, byte(l >> 8), byte(l)}, nil
	case l>>24 == 0:
		return []byte{0x83, byte(l >> 16), byte(l >> 8), byte(l)}, nil
	case l>>32 == 0:
		return []byte{0x84, byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}, nil
	default:
		return nil, errInvalidLength
	}
}
