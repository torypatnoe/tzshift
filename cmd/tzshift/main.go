// Command tzshift translates a posted timestamp into the user's local time and
// their configured roster of team/system zones. `show` is the default
// subcommand; `list` is reference output.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/torypatnoe/tztools/internal/config"
	"github.com/torypatnoe/tztools/internal/render"
	"github.com/torypatnoe/tztools/internal/tz"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(argv []string) int {
	flags, args, err := parseFlags(argv)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tzshift:", err)
		return 2
	}
	if flags.help {
		fmt.Fprint(os.Stdout, usage)
		return 0
	}

	// Subcommand dispatch: `list`, `config`, explicit `show`, else default `show`.
	cmd := "show"
	if len(args) > 0 {
		switch args[0] {
		case "list", "config", "show":
			cmd, args = args[0], args[1:]
		}
	}

	color := flags.color() // resolved against TTY

	switch cmd {
	case "list":
		return runList()
	case "config":
		return runConfig(args)
	default:
		return runShow(args, flags.to, color)
	}
}

// runConfig handles `tzshift config create` and `tzshift config list`.
func runConfig(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "tzshift config: expected a subcommand (create or list)")
		return 2
	}
	switch args[0] {
	case "create":
		path, err := config.Create()
		if err != nil {
			fmt.Fprintln(os.Stderr, "tzshift config create:", err)
			return 1
		}
		fmt.Println("created", path)
		return 0
	case "list":
		path := config.Path()
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tzshift config list: no config at %s (run 'tzshift config create')\n", path)
			return 1
		}
		os.Stdout.Write(data)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "tzshift config: unknown subcommand %q (expected create or list)\n", args[0])
		return 2
	}
}

func runShow(args []string, to string, color bool) int {
	req, err := tz.ParseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tzshift:", err)
		return 1
	}

	aliases, abbrev, entries, notice := roster()
	if notice != "" {
		fmt.Fprint(os.Stderr, notice)
	}

	res, err := req.Instant(aliases, abbrev, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "tzshift:", err)
		return 1
	}

	if to != "" {
		_, name, err := tz.ResolveZone(to, aliases, abbrev)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tzshift: --to:", err)
			return 1
		}
		entries = append(entries, tz.Entry{Label: to, Zone: name})
	}

	// Always show the zone the user typed and their own local time (the date
	// anchor for +N/-N), adding either only if the roster doesn't cover it. When
	// no source zone was given (SourceZone == "Local"), the local row below is
	// the source, so it isn't added separately.
	if res.SourceZone != "" && res.SourceZone != "Local" {
		entries = ensureZone(entries, tz.Entry{Label: res.SourceZone, Zone: res.SourceZone})
	}
	entries = ensureLocal(entries, res.Instant)

	rows := tz.Rows(res.Instant, entries, res.SourceZone)
	render.Show(os.Stdout, rows, render.Options{Color: color})
	return 0
}

// ensureZone appends e unless an entry already maps to the same IANA zone.
func ensureZone(entries []tz.Entry, e tz.Entry) []tz.Entry {
	for _, x := range entries {
		if x.Zone == e.Zone {
			return entries
		}
	}
	return append(entries, e)
}

// ensureLocal appends a `you` row for the host's local zone unless the roster
// already shows it — detected by a you/local label or a row at the same offset.
func ensureLocal(entries []tz.Entry, instant time.Time) []tz.Entry {
	_, localOff := instant.In(time.Local).Zone()
	for _, x := range entries {
		if strings.EqualFold(x.Label, "you") || strings.EqualFold(x.Label, "local") {
			return entries
		}
		if loc, err := time.LoadLocation(x.Zone); err == nil {
			if _, off := instant.In(loc).Zone(); off == localOff {
				return entries
			}
		}
	}
	return append(entries, tz.Entry{Label: "you", Zone: "Local"})
}

func runList() int {
	_, abbrev, _, notice := roster()
	if notice != "" {
		fmt.Fprint(os.Stderr, notice)
	}
	now := time.Now()
	names := tz.SortedByOffset(tz.ZoneNames, now) // east-most first, mirroring show
	render.List(os.Stdout, names, abbrev, now)
	return 0
}

// roster loads config or falls back, returning resolution maps, output entries,
// and any first-run notice to print on stderr.
func roster() (aliases, abbrev map[string]string, entries []tz.Entry, notice string) {
	cfg, found, err := config.Load()
	if err != nil {
		// Malformed config: surface it but degrade to the fallback roster so the
		// translation still succeeds. (Exit-code strictness is on parse/zone errors.)
		return map[string]string{}, map[string]string{}, config.Fallback(),
			fmt.Sprintf("tzshift: %v\n", err)
	}
	if !found {
		// No config is a normal state, not an error: fall back silently.
		// `tzshift --help` documents the file and the fallback.
		return map[string]string{}, map[string]string{}, config.Fallback(), ""
	}
	return cfg.Zones, cfg.Abbreviations, cfg.Entries(), ""
}

type flagSet struct {
	to      string
	noColor bool
	help    bool
}

// color reports whether ANSI color should be emitted: only on a TTY and only
// when --no-color was not given (Spec AC10/AC16).
func (f flagSet) color() bool {
	if f.noColor {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// parseFlags pulls tzshift's flags out of argv, leaving positional args (flags
// may be interspersed with the timestamp).
func parseFlags(argv []string) (flagSet, []string, error) {
	var f flagSet
	var rest []string
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		switch {
		case a == "--no-color":
			f.noColor = true
		case a == "-h" || a == "--help":
			f.help = true
		case a == "--to":
			if i+1 >= len(argv) {
				return f, nil, fmt.Errorf("--to needs a zone (IANA name or alias)")
			}
			i++
			f.to = argv[i]
		case len(a) > 5 && a[:5] == "--to=":
			f.to = a[5:]
		default:
			rest = append(rest, a)
		}
	}
	return f, rest, nil
}

const usage = `tzshift — translate a timestamp into your roster's zones

Usage:
  tzshift [TIME [DATE] ZONE] [--to ZONE] [--no-color]
  tzshift list
  tzshift config create | list

Examples:
  tzshift 14:30 Asia/Kolkata        wall-clock time in a zone (date = today)
  tzshift 14:30                      ...with no zone, time is in YOUR local zone
  tzshift 14:30 india               using a [zones] alias
  tzshift 2026-06-09 14:30 UTC      explicit date
  tzshift 1749571200                epoch seconds
  tzshift                           current time across your roster
  tzshift 14:30 UTC --to Europe/Berlin   add a one-off zone
  tzshift list                      known zones (east-most first) + your abbreviations
  tzshift config create             write a starter ~/.config/tzshift/config.toml
  tzshift config list               print your current config file

Config:
  ~/.config/tzshift/config.toml defines your roster and optional shortcuts:

    [zones]                          # label = IANA zone (output rows; also source aliases)
    you   = "America/Los_Angeles"
    ny-dc = "America/New_York"

    [abbreviations]                  # optional: your own source-zone shortcuts
    IST = "Asia/Kolkata"

  With no config file, tzshift falls back to local time, UTC, and America/Denver.
`
