// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	iso "cunicu.li/go-iso7816"
)

var _ iso.PCSCCard = (*MockCard)(nil)

type call struct {
	Method   string
	Command  []byte
	Response []byte
}

// MockCard is a wrapper around iso.PCSCCard which
// reads/writes transcripts of the commands (APDUs)
// from a file to emulate a real smart card with a mock object.
type MockCard struct {
	mock.Mock
	test *testing.T

	calls []call
	next  iso.PCSCCard
}

// NewMockCard creates a new smart card mock.
//
// If realCard is nil, we attempt to load a transcript of the
// expected commands from a transcript file at:
// `mockdata/t.Name()`.
//
// If realCard is not nil, all exchanged commands will
// be recorded and written the the transcript file above during Close().
func NewMockCard(t *testing.T, realCard iso.PCSCCard) (c *MockCard, err error) {
	c = &MockCard{
		test: t,
		next: realCard,
	}

	c.Mock.Test(c.test)

	if c.next == nil {
		if err := c.LoadTranscript(); err != nil {
			return nil, fmt.Errorf("failed to load transcript: %w", err)
		}
	}

	return c, nil
}

// Close invokes WriteTranscript() in case a real-card
// was passed to NewMockCard().
func (c *MockCard) Close() error {
	if c.next != nil {
		if err := c.WriteTranscript(); err != nil {
			return fmt.Errorf("failed to write transcript: %w", err)
		}
	}

	c.Mock.AssertExpectations(c.test)

	return nil
}

func (c *MockCard) Transmit(cmd []byte) (resp []byte, err error) {
	if c.next != nil {
		resp, err = c.next.Transmit(cmd)

		c.calls = append(c.calls, call{
			Method:   "Transmit",
			Command:  cmd,
			Response: resp,
		})
	} else {
		args := c.Mock.MethodCalled("Transmit", cmd)

		resp = args.Get(0).([]byte) //nolint:forcetypeassert
		err = args.Error(1)
	}

	return resp, err
}

func (c *MockCard) BeginTransaction() error {
	if c.next != nil {
		c.calls = append(c.calls, call{
			Method: "BeginTransaction",
		})

		return c.next.BeginTransaction()
	}

	args := c.Mock.MethodCalled("BeginTransaction")
	return args.Error(0)
}

func (c *MockCard) EndTransaction() error {
	if c.next != nil {
		c.calls = append(c.calls, call{
			Method: "EndTransaction",
		})

		return c.next.EndTransaction()
	}

	args := c.Mock.MethodCalled("EndTransaction")
	return args.Error(0)
}

// LoadTranscript loads the a command transcript from
// `mockdata/t.Name()` and configures the mock object
// with the expected calls to Transmit(), BeginTransaction()
// and EndTransaction().
func (c *MockCard) LoadTranscript() error {
	fn := filepath.Join("mockdata", c.test.Name())
	f, err := os.OpenFile(fn, os.O_RDONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open transcript: %w", err)
	}

	c.test.Logf("Mock transcript loaded from: %s", fn)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) < 1 || strings.HasPrefix(parts[0], "#") {
			continue
		}

		action := parts[0]
		if action != "on" {
			continue
		}

		method := parts[1]
		args := parts[2:]

		switch method {
		case "Transmit":
			cmd, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode command: %w", err)
			}

			resp, err := hex.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}

			c.Mock.On(method, cmd).Return(resp, nil).Once()

		case "BeginTransaction", "EndTransaction":
			c.Mock.On(method).Return(nil)
		}
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close transcript: %w", err)
	}

	return nil
}

// WriteTranscript writes recorded commands exchanged with
// a real card to the  transcript file at `mockdata/t.Name()`.
func (c *MockCard) WriteTranscript() error {
	if err := os.MkdirAll("mockdata", 0o755); err != nil {
		return fmt.Errorf("failed to create mockdata directory: %w", err)
	}

	fn := filepath.Join("mockdata", c.test.Name())
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open transcript: %w", err)
	}

	fmt.Fprintf(f, "# Mockfile/v1 created=%s\n", time.Now().Format(time.RFC3339))

	for _, call := range c.calls {
		switch call.Method {
		case "Transmit":
			fmt.Fprintln(f, "on", "Transmit", hex.EncodeToString(call.Command), hex.EncodeToString(call.Response))

		default:
			fmt.Fprintln(f, "on", call.Method)
		}
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close transcript: %w", err)
	}

	return nil
}
