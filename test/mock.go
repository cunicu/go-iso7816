// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bufio"
	"cmp"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	iso "cunicu.li/go-iso7816"
)

var ErrMalformedMockfile = fmt.Errorf("invalid mockfile")

var _ iso.PCSCCard = (*MockCard)(nil)

type call struct {
	Start, End time.Time
	Method     string
	Command    []byte
	Response   []byte
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
		start := time.Now()

		resp, err = c.next.Transmit(cmd)

		c.calls = append(c.calls, call{
			Start:    start,
			End:      time.Now(),
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
		start := time.Now()

		err := c.next.BeginTransaction()

		c.calls = append(c.calls, call{
			Start:  start,
			End:    time.Now(),
			Method: "BeginTransaction",
		})

		return err
	}

	args := c.Mock.MethodCalled("BeginTransaction")
	return args.Error(0)
}

func (c *MockCard) EndTransaction() error {
	if c.next != nil {
		start := time.Now()

		err := c.next.EndTransaction()

		c.calls = append(c.calls, call{
			Start:  start,
			End:    time.Now(),
			Method: "EndTransaction",
		})

		return err
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

	if !scanner.Scan() {
		return ErrMalformedMockfile
	} else if firstLine := scanner.Text(); firstLine != "mockfile" {
		return fmt.Errorf("%w: invalid first file: %s", ErrMalformedMockfile, firstLine)
	}

	for scanner.Scan() {
		line := scanner.Text()
		cols := splitColumns(line)

		if len(cols) < 1 || strings.HasPrefix(cols[0], "#") {
			continue
		}

		action := cols[0]
		if action != "on" {
			continue
		}

		start, err := strconv.ParseFloat(cols[1], 64)
		if err != nil {
			return fmt.Errorf("failed to parse start column: %w", err)
		}

		end, err := strconv.ParseFloat(cols[2], 64)
		if err != nil {
			return fmt.Errorf("failed to parse end column: %w", err)
		}

		method := cols[3]
		args := cols[4:]

		var call *mock.Call

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

			c.Mock.On(method, cmd).
				Return(resp, nil)

		case "BeginTransaction", "EndTransaction":
			c.Mock.On(method).
				Return(nil)
		}

		call.
			Once().
			After(time.Duration((end - start) * float64(time.Millisecond)))
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
	defer f.Close()

	fmt.Fprintln(f, "mockfile")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "file.version", "v2")
	fmt.Fprintln(f, "file.created", time.Now().Format(time.RFC3339))

	if hostname, err := os.Hostname(); err == nil {
		fmt.Fprintln(f, "file.creator", fmt.Sprintf("%s@%s", os.Getenv("USER"), hostname))
	}

	fmt.Fprintln(f)

	if mc, ok := c.next.(iso.MetadataCard); ok {
		forEachSorted(mc.Metadata(), func(key, value string) {
			fmt.Fprintln(f, "meta", key, value)
		})
	}

	if len(c.calls) == 0 {
		return nil
	}

	fmt.Fprintf(f, "\n#  %s %s method\n", fmt.Sprintf("%8s", "start"), fmt.Sprintf("%8s", "end"))

	start := c.calls[0].Start
	for _, call := range c.calls {
		callStart := fmt.Sprintf("%8.3f", inMilliseconds(call.Start.Sub(start)))
		callEnd := fmt.Sprintf("%8.3f", inMilliseconds(call.Start.Sub(start)))

		switch call.Method {
		case "Transmit":
			fmt.Fprintln(f, "on", callStart, callEnd, call.Method, hex.EncodeToString(call.Command), hex.EncodeToString(call.Response))

		default:
			fmt.Fprintln(f, "on", callStart, callEnd, call.Method)
		}
	}

	return nil
}

func inMilliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

func splitColumns(line string) (cols []string) {
	r := strings.NewReader(line)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	for s.Scan() {
		cols = append(cols, s.Text())
	}

	return cols
}

func forEachSorted[V any, K cmp.Ordered](s map[K]V, cb func(k K, v V)) {
	ks := []K{}

	for k := range s {
		ks = append(ks, k)
	}

	slices.Sort(ks)

	for _, k := range ks {
		cb(k, s[k])
	}
}
