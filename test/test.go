// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ebfe/scard"
	"github.com/stretchr/testify/require"

	iso "cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/drivers/pcsc"
	"cunicu.li/go-iso7816/filter"
)

var (
	ErrFailedToCastCard  = errors.New("failed to cast card")
	ErrCardNotResettable = errors.New("card not resettable")
)

// DangerousWipeRealCard is a global flag which should be used to safe guard any
// real smart cards from being accidentally wiped.
// nolint: gochecknoglobals
var DangerousWipeRealCard = os.Getenv("TEST_DANGEROUS_WIPE_REAL_CARD") != ""

// WithCard opens the first card matching the filter flt and wraps it with the
// TraceCard and MockCard wrappers to trace commands for debugging as well
// as record a transcript of commands exchanged during the test.
// Finally the provided test callback is invoked.
// The user is expected to implement any test code within the provided callback.
// Thanks to wrapping the card into a mock object, tests can also be executed
// without the presence of a real card (e.g. in a CI environment).
// In this case, previously recorded command transcripts are used to emulate &
// assert the communication with the card.
func WithCard(t *testing.T, flt filter.Filter, cb func(t *testing.T, card *iso.Card)) {
	require := require.New(t)

	var realCard iso.PCSCCard

	if DangerousWipeRealCard {
		ctx, err := scard.EstablishContext()
		require.NoError(err)

		defer func() {
			err = ctx.Release()
			require.NoError(err)
		}()

		if realCard, err = pcsc.OpenFirstCard(ctx, flt, true); errors.Is(err, pcsc.ErrNoCardFound) {
			t.Log("Warning: no real cards found. Using mocked card instead!")
		} else if err != nil {
			t.Fatalf("failed to open card: %s", err)
		} else {
			defer func() {
				err := realCard.Close()
				require.NoError(err)
			}()
		}
	} else {
		t.Log("Warning: Running with mocked smart card. Set env var TEST_DANGEROUS_WIPE_REAL_CARD=1 to test against a real smart card.")
	}

	withMock := func(t *testing.T) {
		mockedCard, err := NewMockCard(t, realCard)
		require.NoError(err)

		tracedCard := NewTraceCard(mockedCard, slog.Default())
		isoCard := iso.NewCard(tracedCard)

		cb(t, isoCard)

		err = mockedCard.Close()
		require.NoError(err)
	}

	mockDir := filepath.Join("mockdata", t.Name())
	if fi, err := os.Stat(mockDir); err == nil {
		require.True(fi.IsDir(), "Mockdata directory must be a directory")
	} else if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(mockDir, 0o755)
		require.NoError(err)
	}

	testName := os.Getenv("TEST_RUN_NAME")
	if testName == "" {
		testName = "latest"
	}

	if realCard != nil {
		t.Run(testName, withMock)
	} else {
		entries, err := os.ReadDir(mockDir)
		require.NoError(err)

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// subMockDir := filepath.Join(mockDir, entry.Name())
			t.Run(entry.Name(), withMock)
		}
	}
}

// PCSCCard unwraps a card which has been passed to the callback of WithCard().
// It can be useful to assert that a test is currently running with a real card.
func PCSCCard(card *iso.Card) iso.PCSCCard {
	prev := card.PCSCCard

	for {
		switch card := prev.(type) {
		case *iso.Card:
			prev = card.PCSCCard
		case *MockCard:
			prev = card.next
		case *TraceCard:
			prev = card.next
		case *pcsc.Card:
			return card
		default:
			return nil
		}
	}
}

// ResetCard is a helper to reset a test card
func ResetCard(card *iso.Card) error {
	pcscCard := PCSCCard(card)
	if pcscCard == nil {
		return ErrFailedToCastCard
	}

	resetCard, ok := pcscCard.(iso.ResettableCard)
	if !ok {
		return ErrCardNotResettable
	}

	return resetCard.Reset()
}
