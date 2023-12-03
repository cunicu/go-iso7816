// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// Package tlv implements various version (ASN.1 BER, Simple, Compact) TLV encoding using in ISO 7816-4.
package tlv

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"slices"
)

var errInvalidLength = errors.New("invalid length")

type TagValue struct {
	Tag        Tag
	Value      []byte
	Children   TagValues
	SkipLength bool
}

func (tv TagValue) Equal(w TagValue) bool {
	if !tv.Tag.IsConstructed() {
		return bytes.Equal(tv.Value, w.Value)
	}

	return tv.Children.Equal(w.Children)
}

func New(t Tag, values ...any) (tv TagValue) {
	tv.Tag = t

	for _, value := range values {
		tv.Append(value)
	}

	return tv
}

func (tv *TagValue) Append(value any) {
	switch v := value.(type) {
	case encoding.BinaryMarshaler:
		data, err := v.MarshalBinary()
		if err != nil {
			panic("failed to marshal")
		}

		tv.Value = append(tv.Value, data...)

	case byte:
		tv.Value = append(tv.Value, v)

	case []byte:
		tv.Value = append(tv.Value, v...)

	case string:
		tv.Value = append(tv.Value, []byte(v)...)

	case uint16:
		tv.Value = binary.BigEndian.AppendUint16(tv.Value, v)
	case uint32:
		tv.Value = binary.BigEndian.AppendUint32(tv.Value, v)
	case uint64:
		tv.Value = binary.BigEndian.AppendUint64(tv.Value, v)

	case TagValue:
		tv.Children = append(tv.Children, v)
	case TagValues:
		tv.Children = append(tv.Children, v...)
	}
}

type TagValues []TagValue

func (tvs TagValues) Get(tag Tag) ([]byte, TagValues, bool) {
	for _, tv := range tvs {
		if tv.Tag == tag {
			return tv.Value, tv.Children, true
		}
	}

	return nil, nil, false
}

func (tvs *TagValues) Put(tv TagValue) {
	*tvs = append(*tvs, tv)
}

func (tvs *TagValues) Pop(tag Tag) (TagValue, bool) {
	for i, tv := range *tvs {
		if tv.Tag == tag {
			*tvs = slices.Delete(*tvs, i, i+1)
			return tv, true
		}
	}

	return TagValue{}, false
}

func (tvs TagValues) GetChild(tag Tag, subs ...Tag) ([]byte, TagValues, bool) {
	value, children, ok := tvs.Get(tag)
	if ok {
		if len(subs) > 0 && len(children) > 0 {
			return children.GetChild(subs[0], subs[1:]...)
		}

		return value, children, true
	}

	return nil, nil, false
}

func (tvs *TagValues) GetAll(tag Tag) (s TagValues) {
	for _, tv := range *tvs {
		if tv.Tag == tag {
			s = append(s, tv)
		}
	}

	return s
}

func (tvs *TagValues) DeleteAll(tag Tag) (removed int) {
	var n, r int
	for _, tv := range *tvs {
		if tv.Tag != tag {
			(*tvs)[n] = tv
			n++
		} else {
			r++
		}
	}
	(*tvs) = (*tvs)[:n]

	return r
}

func (tvs *TagValues) PopAll(tag Tag) (s TagValues) {
	var n int
	for _, tv := range *tvs {
		if tv.Tag != tag {
			(*tvs)[n] = tv
			n++
		} else {
			s = append(s, tv)
		}
	}
	(*tvs) = (*tvs)[:n]

	return s
}

func (tvs TagValues) Equal(w TagValues) bool {
	if len(tvs) != len(w) {
		return false
	}

	for i, vc := range tvs {
		wc := w[i]

		if !vc.Equal(wc) {
			return false
		}
	}

	return true
}
