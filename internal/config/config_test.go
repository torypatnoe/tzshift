package config

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadValid(t *testing.T) {
	path := write(t, `
[zones]
you   = "America/Los_Angeles"
india = "Asia/Kolkata"

[abbreviations]
IST = "Asia/Kolkata"
`)
	cfg, found, err := LoadFile(path)
	if err != nil || !found {
		t.Fatalf("LoadFile: found=%v err=%v", found, err)
	}
	if cfg.Zones["india"] != "Asia/Kolkata" {
		t.Errorf("zones[india] = %q", cfg.Zones["india"])
	}
	if cfg.Abbreviations["IST"] != "Asia/Kolkata" {
		t.Errorf("abbreviations[IST] = %q", cfg.Abbreviations["IST"])
	}
	if len(cfg.Entries()) != 2 {
		t.Errorf("Entries() = %d, want 2", len(cfg.Entries()))
	}
}

func TestLoadMissing(t *testing.T) {
	_, found, err := LoadFile(filepath.Join(t.TempDir(), "nope.toml"))
	if err != nil {
		t.Fatalf("missing file should not error: %v", err)
	}
	if found {
		t.Fatal("found should be false for a missing file")
	}
}

func TestLoadMalformed(t *testing.T) {
	path := write(t, "this is = not = valid toml [[[")
	if _, _, err := LoadFile(path); err == nil {
		t.Fatal("expected error for malformed TOML")
	}
}

func TestLoadBadZone(t *testing.T) {
	path := write(t, `
[zones]
oops = "Not/AZone"
`)
	if _, _, err := LoadFile(path); err == nil {
		t.Fatal("expected error for invalid IANA zone")
	}
}

func TestCreate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	path, err := Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if path != Path() {
		t.Errorf("Create path = %q, want %q", path, Path())
	}

	// The template must parse and round-trip through Load.
	cfg, found, err := LoadFile(path)
	if err != nil || !found {
		t.Fatalf("generated template failed to load: found=%v err=%v", found, err)
	}
	if len(cfg.Zones) == 0 {
		t.Error("template should define some [zones]")
	}

	// Second create must refuse to overwrite.
	if _, err := Create(); err == nil {
		t.Error("Create should refuse to overwrite an existing config")
	}
}

func TestFallback(t *testing.T) {
	fb := Fallback()
	if len(fb) != 3 {
		t.Fatalf("Fallback() = %d entries, want 3", len(fb))
	}
	if fb[0].Zone != "Local" || fb[1].Zone != "UTC" || fb[2].Zone != "America/Denver" {
		t.Errorf("unexpected fallback zones: %+v", fb)
	}
}
