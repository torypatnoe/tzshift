// Package config loads tzshift's TOML roster and provides the no-config
// fallback. The file is the only persistent state and is parsed as data only.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/torypatnoe/tztools/internal/tz"
)

// Config is the resolved view of config.toml.
type Config struct {
	Zones         map[string]string // alias -> IANA zone
	Abbreviations map[string]string // user abbreviation -> IANA zone
}

type fileShape struct {
	Zones         map[string]string `toml:"zones"`
	Abbreviations map[string]string `toml:"abbreviations"`
}

// Path returns ~/.config/tzshift/config.toml, honoring XDG_CONFIG_HOME.
func Path() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(".config", "tzshift", "config.toml")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "tzshift", "config.toml")
}

// Load reads the config at Path(). The returned bool is false when no file
// exists (callers then use Fallback). A malformed file or bad zone is an error.
func Load() (*Config, bool, error) {
	return LoadFile(Path())
}

// LoadFile is Load against an explicit path (used in tests).
func LoadFile(path string) (*Config, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("reading %s: %w", path, err)
	}

	var f fileShape
	if err := toml.Unmarshal(data, &f); err != nil {
		return nil, true, fmt.Errorf("parsing %s: %w", path, err)
	}

	cfg := &Config{Zones: f.Zones, Abbreviations: f.Abbreviations}
	if cfg.Zones == nil {
		cfg.Zones = map[string]string{}
	}
	if cfg.Abbreviations == nil {
		cfg.Abbreviations = map[string]string{}
	}
	if err := cfg.validate(path); err != nil {
		return nil, true, err
	}
	return cfg, true, nil
}

func (c *Config) validate(path string) error {
	for label, zone := range c.Zones {
		if _, err := time.LoadLocation(zone); err != nil {
			return fmt.Errorf("%s: [zones] %q = %q is not a valid IANA zone", path, label, zone)
		}
	}
	for abbr, zone := range c.Abbreviations {
		if _, err := time.LoadLocation(zone); err != nil {
			return fmt.Errorf("%s: [abbreviations] %q = %q is not a valid IANA zone", path, abbr, zone)
		}
	}
	return nil
}

// Entries returns the roster as engine entries (unordered; Rows sorts them).
func (c *Config) Entries() []tz.Entry {
	out := make([]tz.Entry, 0, len(c.Zones))
	for label, zone := range c.Zones {
		out = append(out, tz.Entry{Label: label, Zone: zone})
	}
	return out
}

// Fallback is the no-config roster: host-local time, UTC, and America/Denver
// (Spec AC7). "Local" resolves to the host zone via time.LoadLocation.
func Fallback() []tz.Entry {
	return []tz.Entry{
		{Label: "local", Zone: "Local"},
		{Label: "UTC", Zone: "UTC"},
		{Label: "America/Denver", Zone: "America/Denver"},
	}
}

// FallbackMessage is the helpful first-run hint printed alongside fallback output.
func FallbackMessage() string {
	return fmt.Sprintf(`No config found at %s -- showing local, UTC, and America/Denver.
Create that file with a [zones] roster to choose your own, e.g.:

  [zones]
  you   = "America/Los_Angeles"
  ny-dc = "America/New_York"

  # optional: your own source-zone shortcuts
  [abbreviations]
  IST = "Asia/Kolkata"
`, Path())
}
