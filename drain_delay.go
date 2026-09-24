package appkit

import (
	"context"
	"fmt"
	"time"
)

// DrainDelay specifies how long an application should wait while draining
// in-flight work during shutdown.
//
// DrainDelay is represented as a time.Duration and supports the standard text
// marshaling interfaces for configuration serialization.
type DrainDelay time.Duration

// Duration returns the drain delay as a time.Duration.
func (d DrainDelay) Duration() time.Duration { return time.Duration(d) }

// String returns the drain delay in the standard time.Duration format.
func (d DrainDelay) String() string { return d.Duration().String() }

// Wait blocks for the configured drain delay or until ctx is cancelled.
//
// If the delay completes first, Wait returns nil. If the context is cancelled
// first, Wait returns the context's error.
func (d DrainDelay) Wait(ctx context.Context) error {
	t := time.NewTimer(d.Duration())
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Set parses s and updates the drain delay.
//
// The value must use the duration format accepted by [time.ParseDuration].
// An empty string sets the delay to zero. Negative durations are rejected.
func (d *DrainDelay) Set(s string) error {
	return d.UnmarshalText([]byte(s))
}

// Get returns the drain delay as an any value.
//
// Get is provided for compatibility with configuration interfaces that expect
// a getter returning any.
func (d DrainDelay) Get() any {
	return d
}

// MarshalText implements [encoding.TextMarshaler].
//
// The delay is serialized using the standard time.Duration string format.
// Negative delays are rejected.
func (d DrainDelay) MarshalText() ([]byte, error) {
	if d < 0 {
		return nil, fmt.Errorf("invalid drain delay %s", d)
	}
	return []byte(d.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
//
// The value must use the duration format accepted by [time.ParseDuration].
// An empty value sets the delay to zero. Negative durations are rejected.
//
// A nil receiver returns an error.
func (d *DrainDelay) UnmarshalText(text []byte) error {
	if d == nil {
		return fmt.Errorf("nil DrainDelay")
	}
	if len(text) == 0 {
		*d = 0
		return nil
	}

	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return fmt.Errorf("invalid drain delay %q: %w", text, err)
	}

	if parsed < 0 {
		return fmt.Errorf("drain delay must not be negative: %s", parsed)
	}

	*d = DrainDelay(parsed)
	return nil
}
