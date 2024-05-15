// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package feitian_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/devices/feitian"
	"cunicu.li/go-iso7816/filter"
	"cunicu.li/go-iso7816/test"
)

func withCard(t *testing.T, cb func(t *testing.T, c *feitian.Card)) {
	test.WithCard(t, filter.IsFeitian, func(t *testing.T, c *iso7816.Card) {
		cb(t, &feitian.Card{c})
	})
}

func TestSerialNumber(t *testing.T) {
	withCard(t, func(t *testing.T, c *feitian.Card) {
		require := require.New(t)

		sno, err := c.SerialNumber()
		require.NoError(err)
		require.NotEmpty(sno)

		t.Logf("Serial Number: %s", sno)
	})
}

func TestCosVersion(t *testing.T) {
	withCard(t, func(t *testing.T, c *feitian.Card) {
		require := require.New(t)

		v, err := c.COSVersion()
		require.NoError(err)
		require.NotEmpty(v)

		t.Logf("COS version: %s", v)
	})
}
