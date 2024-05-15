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

func TestGetUUID(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		uuidBuf, err := nk.GetUUID(card)
		require.NoError(err)

		uid, err := uuid.FromBytes(uuidBuf)
		require.NoError(err)

		t.Logf("UUID: %s", uid)
	})
}

func TestGetRandom(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		rand, err := nk.GetRandom(card)
		require.NoError(err)
		require.Len(rand, nk.LenRandom)

		t.Logf("Random: %s", hex.EncodeToString(rand))
	})
}

func TestReboot(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		err = nk.Reboot(card)
		require.NoError(err)
	})
}

func TestIsLocked(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		locked, err := nk.IsLocked(card)
		require.NoError(err)
		require.True(locked)
	})
}

func TestGetFirmwareVersion(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		ver, err := nk.GetFirmwareVersion(card)
		require.NoError(err)

		t.Logf("Version: %+#v", ver)
	})
}

func TestGetStatus(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, card *iso.Card) {
		require := require.New(t)

		_, err := card.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		ds, err := nk.GetDeviceStatus(card)
		require.NoError(err)

		t.Logf("Status: %+#v", ds)
	})
}
