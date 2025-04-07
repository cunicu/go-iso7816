// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv

func DecodeCompact(buf []byte) (tvs TagValues, err error) {
	for len(buf) > 0 {
		if len(buf) < 1 {
			return nil, errInvalidLength
		}

		l := int(buf[0] & 0xF)

		if len(buf) < l+1 {
			return nil, errInvalidLength
		}

		tvs = append(tvs, TagValue{
			Tag:   Tag(buf[0] >> 4),
			Value: buf[1 : 1+l],
		})
		buf = buf[1+l:]
	}

	return tvs, nil
}

func EncodeCompact(tvs ...TagValue) (buf []byte, err error) {
	for _, tv := range tvs {
		if tv.Tag > 0xF {
			return nil, ErrTagToBig
		}

		if len(tv.Value) > 0xF {
			return nil, ErrValueToLarge
		}

		buf = append(buf, byte((int(tv.Tag)<<4)|len(tv.Value))) //nolint:gosec
		buf = append(buf, tv.Value...)
	}

	return buf, nil
}
