// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816/encoding/tlv"
)

func TestTagBER(t *testing.T) {
	require := require.New(t)

	cases := []struct {
		tag         tlv.Tag
		class       tlv.Class
		constructed bool
	}{
		{tlv.NewBERTag(0x20, tlv.ClassUniversal), tlv.ClassUniversal, false},
		{tlv.NewBERTag(0x200, tlv.ClassContext), tlv.ClassContext, false},
		{tlv.NewBERTag(0x20000, tlv.ClassPrivate), tlv.ClassPrivate, false},
		{0x21, tlv.ClassUniversal, true},
		{0x01, tlv.ClassUniversal, false},
		{0x41, tlv.ClassApplication, false},
		{0x81, tlv.ClassContext, false},
		{0xC1, tlv.ClassPrivate, false},
	}

	for _, c := range cases {
		require.Equal(c.class, c.tag.Class())
		require.Equal(c.constructed, c.tag.IsConstructed())
	}
}

func TestEncodeBER(t *testing.T) {
	require := require.New(t)

	long1 := make([]byte, 0x80)
	long2 := make([]byte, 0x100)

	expected := append([]byte{}, 0x01, 0x04, 0x10, 0x11, 0x12, 0x13)
	expected = append(expected, 0x02, 0x81, 0x80)
	expected = append(expected, long1...)
	expected = append(expected, 0x03, 0x82, 0x01, 0o0)
	expected = append(expected, long2...)
	expected = append(expected, 0x08, 0x03, 0x20, 0x21, 0x22)
	expected = append(expected, 0x1f, 0x82, 0x00, 0x01, 0x20)

	buf, err := tlv.EncodeBER(
		tlv.New(0x1, []byte{0x10, 0x11, 0x12, 0x13}),
		tlv.New(0x2, long1),
		tlv.New(0x3, long2),
		tlv.New(0x8, []byte{0x20, 0x21, 0x22}),
		tlv.New(tlv.NewBERTag(0x100, tlv.ClassUniversal), []byte{0x20}),
	)
	require.NoError(err)
	require.Equal(expected, buf)

	tvs, err := tlv.DecodeBER(buf)
	require.NoError(err)
	require.Equal(tlv.TagValues{
		tlv.New(0x1, []byte{0x10, 0x11, 0x12, 0x13}),
		tlv.New(0x2, long1),
		tlv.New(0x3, long2),
		tlv.New(0x8, []byte{0x20, 0x21, 0x22}),
		tlv.New(tlv.NewBERTag(0x100, tlv.ClassUniversal), []byte{0x20}),
	}, tvs)
}

func TestEncodeNestedBER(t *testing.T) {
	require := require.New(t)

	expected := []byte{0x21, 0x09, 0x02, 0x02, 0x03, 0x04, 0x03, 0x03, 0x05, 0x06, 0x07}
	expected = append(expected, 0x21, 0x09, 0x02, 0x02, 0x03, 0x04, 0x03, 0x03, 0x05, 0x06, 0x07)
	expected = append(expected, 0x3F, 0x33, 0x09, 0x02, 0x02, 0x03, 0x04, 0x03, 0x03, 0x05, 0x06, 0x07)

	buf, err := tlv.EncodeBER(
		tlv.New(0x01,
			tlv.New(0x2, []byte{3, 4}),
			tlv.New(0x3, []byte{5, 6, 7}),
		),
		tlv.New(0x21, // Test manually setting the constructed bit
			tlv.New(0x2, []byte{3, 4}),
			tlv.New(0x3, []byte{5, 6, 7}),
		),
		tlv.New(tlv.NewBERTag(0x33, tlv.ClassUniversal),
			tlv.New(0x2, []byte{3, 4}),
			tlv.New(0x3, []byte{5, 6, 7}),
		),
	)
	require.NoError(err)
	require.Equal(expected, buf)

	tvs, err := tlv.DecodeBER(buf)
	require.NoError(err)
	require.True(tvs.Equal(tlv.TagValues{
		tlv.New(0x21,
			tlv.New(0x2, []byte{3, 4}),
			tlv.New(0x3, []byte{5, 6, 7}),
		),
		tlv.New(0x21,
			tlv.New(0x2, []byte{3, 4}),
			tlv.New(0x3, []byte{5, 6, 7}),
		),
		tlv.New(tlv.NewBERTag(0x33, tlv.ClassContext),
			tlv.New(0x2, []byte{3, 4}),
			tlv.New(0x3, []byte{5, 6, 7}),
		),
	}))
}

func TestBERError(t *testing.T) {
	require := require.New(t)

	_, err := tlv.DecodeBER([]byte{0x11})
	require.Error(err)
}

func TestDecodeBEREmpty(t *testing.T) {
	require := require.New(t)

	tvs, err := tlv.DecodeBER(nil)
	require.NoError(err)
	require.Empty(tvs)
}

func TestEncodeBEREmpty(t *testing.T) {
	require := require.New(t)

	buf, err := tlv.EncodeBER()
	require.NoError(err)
	require.Empty(buf)
}

func FuzzBER(f *testing.F) {
	f.Fuzz(func(_ *testing.T, buf []byte) {
		tlv.DecodeBER(buf) //nolint:errcheck
	})
}
