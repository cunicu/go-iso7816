// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
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

func withCard(t *testing.T, aid []byte, cb func(t *testing.T, card *yk.Card)) {
	test.WithCard(t, filter.IsYubiKey, func(t *testing.T, card *iso7816.Card) {
		require := require.New(t)

		_, err := card.Select(aid)
		require.NoError(err)

		cb(t, &yk.Card{card})
	})
}

func TestDeviceInfo(t *testing.T) {
	withCard(t, iso7816.AidYubicoManagement, func(t *testing.T, card *yk.Card) {
		require := require.New(t)

		di, err := card.DeviceInfo()
		require.NoError(err)

		t.Logf("Device Info: %+#v", di)
	})
}

func TestStatus(t *testing.T) {
	withCard(t, iso7816.AidYubicoOTP, func(t *testing.T, card *yk.Card) {
		require := require.New(t)

		sts, err := card.Status()
		require.NoError(err)

		t.Logf("Status: %+#v", sts)
	})
}

func TestSerialNumber(t *testing.T) {
	withCard(t, iso7816.AidYubicoOTP, func(t *testing.T, card *yk.Card) {
		require := require.New(t)

		sno, err := card.SerialNumber()
		require.NoError(err)

		t.Logf("Serial Number: %d", sno)
	})
}

func TestFIPSMode(t *testing.T) {
	withCard(t, iso7816.AidYubicoOTP, func(t *testing.T, card *yk.Card) {
		require := require.New(t)

		fm, err := card.FIPSMode()
		require.NoError(err)

		t.Logf("FIPS Mode: %t", fm)
	})
}
