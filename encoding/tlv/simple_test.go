// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816/encoding/tlv"
)

func TestEncodeSimple(t *testing.T) {
	require := require.New(t)

	long := make([]byte, 0xFF)

	expected := append([]byte{}, 0x01, 0x04, 0x10, 0x11, 0x12, 0x13, 0x02, 0xFF, 0x00, 0xFF)
	expected = append(expected, long...)
	expected = append(expected, 0x08, 0x03, 0x20, 0x21, 0x22)

	tvsIn := tlv.TagValues{
		tlv.New(1, []byte{0x10, 0x11, 0x12, 0x13}),
		tlv.New(2, long),
		tlv.New(8, []byte{0x20, 0x21, 0x22}),
	}
	buf, err := tlv.EncodeSimple(tvsIn...)
	require.NoError(err)
	require.Equal(expected, buf)

	tvs, err := tlv.DecodeSimple(buf)
	require.NoError(err)
	require.Equal(tvsIn, tvs)
}

func TestDecodeSimpleError(t *testing.T) {
	require := require.New(t)

	_, err := tlv.DecodeSimple([]byte{0x11})
	require.Error(err)
}

func TestEncodeSimpleError(t *testing.T) {
	require := require.New(t)

	_, err := tlv.EncodeSimple(tlv.New(0x100))
	require.ErrorIs(err, tlv.ErrTagToBig)

	tooBig := make([]byte, 0x10000)

	_, err = tlv.EncodeSimple(tlv.New(0x01, tooBig))
	require.ErrorIs(err, tlv.ErrValueToLarge)
}

func TestDecodeSimpleEmpty(t *testing.T) {
	require := require.New(t)

	tvs, err := tlv.DecodeSimple(nil)
	require.NoError(err)
	require.Empty(tvs)
}

func TestEncodeSimpleEmpty(t *testing.T) {
	require := require.New(t)

	buf, err := tlv.EncodeSimple()
	require.NoError(err)
	require.Empty(buf)
}

func FuzzSimple(f *testing.F) {
	f.Fuzz(func(_ *testing.T, buf []byte) {
		tlv.DecodeSimple(buf) //nolint:errcheck
	})
}
