package feitan_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/devices/feitan"
	"cunicu.li/go-iso7816/filter"
	"cunicu.li/go-iso7816/test"
)

func TestGetSerialNumber(t *testing.T) {
	test.WithCard(t, filter.IsFeitian, func(t *testing.T, c *iso7816.Card) {
		require := require.New(t)

		// _, err := c.Select(iso7816.AidFeitianOTP)
		// require.NoError(err)

		sno, err := feitan.GetSerialNumber(c)
		require.NoError(err)

		t.Logf("Serial Number: %s", sno)
	})
}

func TestGetAppletVersion(t *testing.T) {
	test.WithCard(t, filter.IsFeitian, func(t *testing.T, c *iso7816.Card) {
		require := require.New(t)

		// _, err := c.Select(iso7816.AidFeitianOTP)
		// require.NoError(err)

		sno, err := feitan.GetAppletVersion(c)
		require.NoError(err)

		t.Logf("Applet version: %s", sno)
	})
}
