package appkit

import "time"

// Config contains the main application configuration.
type Config struct {
	Runmode  Runmode        `json:"runmode"  yaml:"runmode"  koanf:"runmode"`
	Name     string         `json:"name"     yaml:"name"     koanf:"name"`
	Shutdown ShutdownConfig `json:"shutdown" yaml:"shutdown" koanf:"shutdown"`
}

// IsDevmod reports whether the application is running in development mode.
func (c Config) IsDevmod() bool { return c.Runmode == RunmodeDev }

// ShutdownConfig configures the application's shutdown timing.
//
// Timeout specifies the main shutdown timeout. DrainDelay specifies an
// additional period used to drain in-flight work.
type ShutdownConfig struct {
	Timeout    time.Duration `json:"timeout"     yaml:"timeout"     koanf:"timeout"`
	DrainDelay DrainDelay    `json:"drain_delay" yaml:"drain_delay" koanf:"drain_delay"`
}

// TotalTimeout returns the total shutdown timeout, including the drain delay.
func (c ShutdownConfig) TotalTimeout() time.Duration {
	return c.Timeout + c.DrainDelay.Duration()
}
