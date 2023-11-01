// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	iso "cunicu.li/go-iso7816"
)

// TraceCard is a wrapper around iso7816.PCSCCard
// which logs a exchanged commands (APDUs) to a log/slog
// logger.
type TraceCard struct {
	logger *slog.Logger
	next   iso.PCSCCard
}

// NewTraceCard wraps a iso7816.PCSCCard into a TraceCard
// which logs a exchanged commands (APDUs) to a log/slog
// logger.
func NewTraceCard(next iso.PCSCCard, logger *slog.Logger) *TraceCard {
	if logger == nil {
		logger = slog.Default()
	}

	return &TraceCard{
		logger: logger,
		next:   next,
	}
}

func (c *TraceCard) Transmit(cmd []byte) ([]byte, error) {
	c.logger.Info("Send ->",
		slog.Any("cmd", hex.EncodeToString(cmd)),
		slog.Int("len", len(cmd)))

	start := time.Now()

	resp, err := c.next.Transmit(cmd)

	end := time.Now()

	args := []any{
		slog.String("resp", hex.EncodeToString(resp)),
		slog.Int("len", len(resp)),
		slog.Duration("after", end.Sub(start)),
	}

	if err == nil {
		c.logger.Info("Recv <-", args...)
	} else {
		args = append(args, slog.String("code", fmt.Sprintf("%x", err)))
		c.logger.Error("Recv <-", args...)
	}

	return resp, err
}

func (c *TraceCard) BeginTransaction() error {
	c.logger.Info("BeginTransaction")

	return c.next.BeginTransaction()
}

func (c *TraceCard) EndTransaction() error {
	c.logger.Info("EndTransaction")

	return c.next.EndTransaction()
}
