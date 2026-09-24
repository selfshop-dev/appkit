package appkit

import (
	"bytes"
	"fmt"
)

// Runmode identifies the application execution mode.
type Runmode int8

const (
	// RunmodeDev is the development execution mode.
	RunmodeDev Runmode = iota - 1

	// RunmodeProd is the production execution mode.
	RunmodeProd

	_minRunmode = RunmodeDev
	_maxRunmode = RunmodeProd

	// InvalidRunmode is the sentinel value used for an invalid run mode.
	InvalidRunmode = _maxRunmode + 1
)

// ParseRunmode parses a textual run mode.
//
// Parsing is case-insensitive. An empty string is treated as RunmodeProd.
// ParseRunmode returns InvalidRunmode together with an error when the value
// is not recognized.
func ParseRunmode(text string) (Runmode, error) {
	var r Runmode
	err := r.UnmarshalText([]byte(text))
	return r, err
}

// String returns the textual representation of the run mode.
//
// RunmodeDev is represented as "dev" and RunmodeProd as "prod". Unknown
// values are represented as "unknown(<value>)".
func (r Runmode) String() string {
	switch r {
	case InvalidRunmode:
		return "invalid"
	case RunmodeDev:
		return "dev"
	case RunmodeProd:
		return "prod"
	default:
		return fmt.Sprintf("unknown(%d)", r)
	}
}

// Set parses s and updates the run mode.
//
// Parsing is case-insensitive. An empty string is treated as RunmodeProd.
// Invalid values set the receiver to InvalidRunmode and return an error.
func (r *Runmode) Set(s string) error {
	return r.UnmarshalText([]byte(s))
}

// Get returns the run mode as an any value.
//
// Get is provided for compatibility with configuration interfaces that expect
// a getter returning any.
func (r Runmode) Get() any {
	return r
}

// MarshalText implements [encoding.TextMarshaler].
//
// RunmodeDev and RunmodeProd are serialized as "dev" and "prod",
// respectively. Invalid run modes return an error.
func (r Runmode) MarshalText() ([]byte, error) {
	if r < _minRunmode ||
		r > _maxRunmode {
		return nil, fmt.Errorf("invalid runmode %d", int(r))
	}
	return []byte(r.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
//
// Parsing is case-insensitive. An empty string is treated as RunmodeProd.
// If the value is not recognized, the receiver is set to InvalidRunmode and
// an error is returned.
//
// A nil receiver returns an error.
func (r *Runmode) UnmarshalText(text []byte) error {
	if r == nil {
		return fmt.Errorf("nil Runmode")
	}
	if r.unmarshalText(text) ||
		r.unmarshalText(bytes.ToLower(text)) {
		return nil
	}
	return fmt.Errorf("unrecognized runmode: %q", text)
}

func (r *Runmode) unmarshalText(text []byte) bool {
	switch string(text) {
	case "dev":
		*r = RunmodeDev
	case "prod", "": // allow empty string
		*r = RunmodeProd
	default:
		*r = InvalidRunmode
		return false
	}
	return true
}
