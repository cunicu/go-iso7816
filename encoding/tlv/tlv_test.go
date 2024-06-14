// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package tlv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816/encoding/tlv"
)

func TestTagValues(t *testing.T) {
	require := require.New(t)

	var tvs tlv.TagValues

	tvs.Put(tlv.New(0x01, "Hallo"))
	require.Len(tvs, 1)

	tvs.Put(tlv.New(0x01, "Welt"))
	require.Len(tvs, 2)

	h, _, ok := tvs.Get(0x01)
	require.True(ok)
	require.Equal([]byte("Hallo"), h)

	tv, ok := tvs.Pop(0x01)
	require.True(ok)
	require.Equal([]byte("Hallo"), tv.Value)
	require.Len(tvs, 1)

	tv, ok = tvs.Pop(0x01)
	require.True(ok)
	require.Equal([]byte("Welt"), tv.Value)
	require.Len(tvs, 0)

	_, ok = tvs.Pop(0x01)
	require.False(ok)
}

func TestTagValuesGetChild(t *testing.T) {
	require := require.New(t)

	tvs := tlv.TagValues{
		tlv.New(0x06, 1, 2, 3),
		tlv.New(0x01,
			tlv.New(0x05),
			tlv.New(0x02,
				tlv.New(0x03, "Hallo"),
				tlv.New(0x04, "Welt"),
			),
		),
		tlv.New(0x07, []byte{0xA, 0xB, 0xC}),
	}

	v, _, ok := tvs.GetChild(0x01, 0x02, 0x04)
	require.True(ok)
	require.Equal([]byte("Welt"), v)

	_, c, ok := tvs.GetChild(0x01, 0x02)
	require.True(ok)
	require.Len(c, 2)
}

func TestTagValuesAll(t *testing.T) {
	require := require.New(t)

	tvs := tlv.TagValues{
		tlv.New(0x01),
		tlv.New(0x02),
		tlv.New(0x03),
		tlv.New(0x04),
		tlv.New(0x02),
		tlv.New(0x02),
		tlv.New(0x04),
		tlv.New(0x04),
	}

	tvs2a := tvs.GetAll(0x02)
	require.Len(tvs2a, 3)

	tvs2b := tvs.GetAll(0x04)
	require.Len(tvs2b, 3)

	tvs3 := tvs.PopAll(0x02)
	require.Len(tvs3, 3)

	tvs4 := tvs.GetAll(0x02)
	require.Len(tvs4, 0)

	r := tvs.DeleteAll(0x04)
	require.Equal(r, 3)
	require.Len(tvs, 2)
}
