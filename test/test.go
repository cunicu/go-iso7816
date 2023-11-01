// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/ebfe/scard"
	"github.com/stretchr/testify/require"

	iso "cunicu.li/go-iso7816"
	"cunicu.li/go-iso7816/drivers/pcsc"
	"cunicu.li/go-iso7816/filter"
)

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

	ctx, err := scard.EstablishContext()
	require.NoError(err)

	realCard, err := pcsc.OpenFirstCard(ctx, flt)
	if errors.Is(err, pcsc.ErrNoCardFound) {
		t.Log("Warn: no real cards found. Using mocked card instead!")
	}

	mockedCard, err := NewMockCard(t, realCard)
	require.NoError(err)

	tracedCard := NewTraceCard(mockedCard, slog.Default())
	isoCard := iso.NewCard(tracedCard)

	cb(t, isoCard)

	err = ctx.Release()
	require.NoError(err)

	err = mockedCard.Close()
	require.NoError(err)
}
