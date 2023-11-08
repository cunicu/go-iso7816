// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv

func EncodeSimple(tvs ...TagValue) (buf []byte, err error) {
	for _, tv := range tvs {
		if tv.Tag == 0 || tv.Tag == 0xFF {
			return nil, ErrInvalidTag
		} else if tv.Tag > 0xff {
			return nil, ErrTagToBig
		}

		if tv.SkipLength {
			buf = append(buf, byte(tv.Tag))
		} else if l := len(tv.Value); l < 0xff {
			buf = append(buf, byte(tv.Tag), byte(l))
		} else if l <= 0xffff {
			buf = append(buf, byte(tv.Tag), 0xff, byte(l>>8), byte(l))
		} else {
			return nil, ErrValueToLarge
		}

		buf = append(buf, tv.Value...)
	}

	return buf, nil
}

func DecodeSimple(buf []byte) (tvs TagValues, err error) {
	for len(buf) > 0 {
		if len(buf) < 2 {
			return nil, errInvalidLength
		}

		var o, l int
		if buf[1] != 0xff {
			o = 2
			l = int(buf[1])
		} else {
			if len(buf) < 4 {
				return nil, errInvalidLength
			}

			o = 4
			l = int(buf[2])<<8 + int(buf[3])

			if len(buf) < 4+l {
				return nil, errInvalidLength
			}
		}

		if len(buf) < o+l {
			return nil, errInvalidLength
		}

		tvs = append(tvs, TagValue{
			Tag:   Tag(buf[0]),
			Value: buf[o : o+l],
		})
		buf = buf[o+l:]
	}

	return tvs, nil
}
