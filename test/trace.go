// SPDX-FileCopyrightText: 2023-2024 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/hex"
	"log/slog"
	"time"

	iso "cunicu.li/go-iso7816"
)

var _ iso.PCSCCard = (*TraceCard)(nil)

// TraceCard is a wrapper around iso7816.PCSCCard
// which logs a exchanged commands (APDUs) to a log/slog
// logger.
type TraceCard struct {
	iso.PCSCCard
	logger *slog.Logger
}

// NewTraceCard wraps a iso7816.PCSCCard into a TraceCard
// which logs a exchanged commands (APDUs) to a log/slog
// logger.
func NewTraceCard(next iso.PCSCCard, logger *slog.Logger) *TraceCard {
	if logger == nil {
		logger = slog.Default()
	}

	return &TraceCard{
		PCSCCard: next,
		logger:   logger,
	}
}

func (c *TraceCard) Transmit(cmd []byte) ([]byte, error) {
	c.logger.Info("Send ->",
		slog.Any("cmd", hex.EncodeToString(cmd)),
		slog.Int("len", len(cmd)))

	start := time.Now()

	resp, err := c.PCSCCard.Transmit(cmd)

	end := time.Now()

	args := []any{
		slog.Duration("after", end.Sub(start)),
	}

	if err == nil {
		args = append(args,
			slog.Int("len", len(resp)),
			slog.String("resp", hex.EncodeToString(resp)))
		c.logger.Info("Recv <-", args...)
	} else {
		args = append(args, slog.Any("error", err))
		c.logger.Error("Recv <-", args...)
	}

	return resp, err
}

func (c *TraceCard) BeginTransaction() error {
	c.logger.Info("BeginTransaction")

	return c.PCSCCard.BeginTransaction()
}

func (c *TraceCard) EndTransaction() error {
	c.logger.Info("EndTransaction")

	return c.PCSCCard.EndTransaction()
}

func (c *TraceCard) Close() error {
	return nil
}
