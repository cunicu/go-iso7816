// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package yubikey_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816"
	yk "cunicu.li/go-iso7816/devices/yubikey"
	"cunicu.li/go-iso7816/filter"
	"cunicu.li/go-iso7816/test"
)

func TestGetDeviceInfo(t *testing.T) {
	test.WithCard(t, filter.IsYubiKey, func(t *testing.T, card *iso7816.Card) {
		require := require.New(t)

		_, err := card.Select(iso7816.AidYubicoManagement)
		require.NoError(err)

		di, err := yk.GetDeviceInfo(card)
		require.NoError(err)

		t.Logf("Device Info: %+#v", di)
	})
}

func TestGetStatus(t *testing.T) {
	test.WithCard(t, filter.IsYubiKey, func(t *testing.T, card *iso7816.Card) {
		require := require.New(t)

		_, err := card.Select(iso7816.AidYubicoOTP)
		require.NoError(err)

		sts, err := yk.GetStatus(card)
		require.NoError(err)

		t.Logf("Status: %+#v", sts)
	})
}

func TestGetSerialNumber(t *testing.T) {
	test.WithCard(t, filter.IsYubiKey, func(t *testing.T, card *iso7816.Card) {
		require := require.New(t)

		_, err := card.Select(iso7816.AidYubicoOTP)
		require.NoError(err)

		sno, err := yk.GetSerialNumber(card)
		require.NoError(err)

		t.Logf("Serial Number: %d", sno)
	})
}

func TestGetFIPSMode(t *testing.T) {
	test.WithCard(t, filter.IsYubiKey, func(t *testing.T, card *iso7816.Card) {
		require := require.New(t)

		_, err := card.Select(iso7816.AidYubicoOTP)
		require.NoError(err)

		fm, err := yk.GetFIPSMode(card)
		require.NoError(err)

		t.Logf("FIPS Mode: %t", fm)
	})
}
