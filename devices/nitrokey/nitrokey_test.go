// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package nitrokey_test

import (
	"testing"

	"cunicu.li/go-iso7816"
	iso "cunicu.li/go-iso7816"
	nk "cunicu.li/go-iso7816/devices/nitrokey"
	"cunicu.li/go-iso7816/filter"
	"cunicu.li/go-iso7816/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetUUID(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, c *iso7816.Card) {
		require := require.New(t)

		_, err := c.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		uuidBuf, err := nk.GetUUID(c)
		require.NoError(err)

		uid, err := uuid.FromBytes(uuidBuf)
		require.NoError(err)

		t.Logf("UUID: %s", uid)
	})
}

func TestGetFirmwareVersion(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, c *iso7816.Card) {
		require := require.New(t)

		_, err := c.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		ver, err := nk.GetFirmwareVersion(c)
		require.NoError(err)

		t.Logf("Version: %+#v", ver)
	})
}

func TestGetStatus(t *testing.T) {
	test.WithCard(t, filter.IsNitrokey3, func(t *testing.T, c *iso7816.Card) {
		require := require.New(t)

		_, err := c.Select(iso.AidSolokeysAdmin)
		require.NoError(err)

		ds, err := nk.GetDeviceStatus(c)
		require.NoError(err)

		t.Logf("Status: %+#v", ds)
	})
}
