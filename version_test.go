// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package iso7816_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	iso "cunicu.li/go-iso7816"
)

func TestParseVersion(t *testing.T) {
	tests := map[string]iso.Version{
		"1":     {1, -1, -1},
		"1.0":   {1, 0, -1},
		"1.0.0": {1, 0, 0},
		"":      {-1, -1, -1},
	}
	for strVer, expVer := range tests {
		t.Run(strVer, func(t *testing.T) {
			ver, err := iso.ParseVersion(strVer)
			require.NoError(t, err)
			require.Equal(t, expVer, ver)
		})
	}
}

func TestParseVersionError(t *testing.T) {
	tests := []string{
		"-1",
		"a",
		"a.1",
		"1.2.3.4",
	}
	for _, strVer := range tests {
		t.Run(strVer, func(t *testing.T) {
			_, err := iso.ParseVersion(strVer)
			require.Error(t, err)
		})
	}
}

func TestVersionLess(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		less bool
	}{
		{"1.2.4", "1.2.4", false},
		{"1.2.4", "1.2.3", true},
		{"1.2.3", "1.2.4", false},
		{"1", "1", false},
		{"1", "2", false},
		{"1", "0", true},
		{"0", "1", false},
	}
	for _, test := range tests {
		a, err := iso.ParseVersion(test.a)
		require.NoError(t, err)

		b, err := iso.ParseVersion(test.b)
		require.NoError(t, err)

		require.True(t, a.Less(b) == test.less, "%s > %s == %t", test.a, test.b, test.less)
	}
}
