// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816/encoding/tlv"
)

func TestEncodeCompact(t *testing.T) {
	require := require.New(t)

	buf, err := tlv.EncodeCompact(
		tlv.New(1, []byte{0x10, 0x11, 0x12, 0x13}),
		tlv.New(8, []byte{0x20, 0x21, 0x22}),
	)
	require.NoError(err)
	require.Equal([]byte{0x14, 0x10, 0x11, 0x12, 0x13, 0x83, 0x20, 0x21, 0x22}, buf)

	tvs, err := tlv.DecodeCompact(buf)
	require.NoError(err)
	require.Equal([]tlv.TagValue{
		{
			Tag:   1,
			Value: []byte{0x10, 0x11, 0x12, 0x13},
		},
		{
			Tag:   8,
			Value: []byte{0x20, 0x21, 0x22},
		},
	}, tvs)
}

func TestDecodeCompactError(t *testing.T) {
	require := require.New(t)

	_, err := tlv.DecodeCompact([]byte{0x11})
	require.Error(err)
}

func TestEncodeCompactError(t *testing.T) {
	require := require.New(t)

	_, err := tlv.EncodeCompact(tlv.New(0x10))
	require.ErrorIs(err, tlv.ErrTagToBig)

	tooBig := make([]byte, 0x10)

	_, err = tlv.EncodeCompact(tlv.New(0x01, tooBig))
	require.ErrorIs(err, tlv.ErrValueToLarge)
}

func TestDecodeCompactEmpty(t *testing.T) {
	require := require.New(t)

	tvs, err := tlv.DecodeCompact(nil)
	require.NoError(err)
	require.Empty(tvs)
}

func TestEncodeCompactEmpty(t *testing.T) {
	require := require.New(t)

	buf, err := tlv.EncodeCompact()
	require.NoError(err)
	require.Empty(buf)
}

func FuzzCompact(f *testing.F) {
	f.Fuzz(func(t *testing.T, buf []byte) {
		tlv.DecodeCompact(buf) //nolint:errcheck
	})
}
