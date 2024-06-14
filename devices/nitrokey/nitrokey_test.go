// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package nitrokey_test

import (
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	iso "cunicu.li/go-iso7816"
	nk "cunicu.li/go-iso7816/devices/nitrokey"
	"cunicu.li/go-iso7816/filter"
	"cunicu.li/go-iso7816/test"
)

func withCard(t *testing.T, cb func(t *testing.T, c *nk.Card)) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		cb(t, &nk.Card{card})
	})
}

func TestGetUUID(t *testing.T) {
	withCard(t, func(t *testing.T, card *nk.Card) {
		require := require.New(t)

		uuidBuf, err := card.UUID()
		require.NoError(err)

		uid, err := uuid.FromBytes(uuidBuf)
		require.NoError(err)

		t.Logf("UUID: %s", uid)
	})
}

func TestRandom(t *testing.T) {
	withCard(t, func(t *testing.T, card *nk.Card) {
		require := require.New(t)

		rand, err := card.Random()
		require.NoError(err)
		require.Len(rand, nk.LenRandom)

		t.Logf("Random: %s", hex.EncodeToString(rand))
	})
}

func TestReboot(t *testing.T) {
	withCard(t, func(t *testing.T, card *nk.Card) {
		require := require.New(t)

		err := card.Reboot()
		require.NoError(err)
	})
}

func TestIsLocked(t *testing.T) {
	withCard(t, func(t *testing.T, card *nk.Card) {
		require := require.New(t)

		locked, err := card.IsLocked()
		require.NoError(err)
		require.True(locked)
	})
}

func TestGetFirmwareVersion(t *testing.T) {
	withCard(t, func(t *testing.T, card *nk.Card) {
		require := require.New(t)

		ver, err := card.FirmwareVersion()
		require.NoError(err)

		t.Logf("Version: %+#v", ver)
	})
}

func TestGetStatus(t *testing.T) {
	withCard(t, func(t *testing.T, card *nk.Card) {
		require := require.New(t)

		ds, err := card.DeviceStatus()
		require.NoError(err)

		t.Logf("Status: %+#v", ds)
	})
}
