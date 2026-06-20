// Command tzshift translates a posted timestamp into the user's local time and
// their configured roster of team/system zones. `show` is the default
// subcommand; `list` is reference output.
package main

import (
	"fmt"
	"os"
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

	// Subcommand dispatch: `list`, explicit `show`, else default to `show`.
	cmd := "show"
	if len(args) > 0 {
		switch args[0] {
		case "list":
			cmd, args = "list", args[1:]
		case "show":
			cmd, args = "show", args[1:]
		}
	}

	color := flags.color() // resolved against TTY

	switch cmd {
	case "list":
		return runList(flags.sort)
	default:
		return runShow(args, flags.to, color)
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
		loc, name, err := tz.ResolveZone(to, aliases, abbrev)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tzshift: --to:", err)
			return 1
		}
		_ = loc
		entries = append(entries, tz.Entry{Label: to, Zone: name})
	}

	rows := tz.Rows(res.Instant, entries)
	render.Show(os.Stdout, render.Header(req, res), rows, render.Options{Color: color})
	return 0
}

func runList(sortMode string) int {
	_, abbrev, _, notice := roster()
	if notice != "" {
		fmt.Fprint(os.Stderr, notice)
	}
	now := time.Now()
	names := tz.SortedByOffset(tz.ZoneNames, now) // default: east-most first, mirroring show
	if sortMode == "name" {
		names = tz.ZoneNames // explicit alphabetical (already sorted)
	}
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
		return map[string]string{}, map[string]string{}, config.Fallback(), config.FallbackMessage()
	}
	return cfg.Zones, cfg.Abbreviations, cfg.Entries(), ""
}

type flagSet struct {
	to      string
	noColor bool
	help    bool
	sort    string // `list --sort=name|offset` (default name)
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
		case a == "--sort":
			if i+1 >= len(argv) {
				return f, nil, fmt.Errorf("--sort needs a value (name or offset)")
			}
			i++
			f.sort = argv[i]
		case len(a) > 7 && a[:7] == "--sort=":
			f.sort = a[7:]
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
	if f.sort != "" && f.sort != "name" && f.sort != "offset" {
		return f, nil, fmt.Errorf("--sort must be name or offset (got %q)", f.sort)
	}
	return f, rest, nil
}

const usage = `tzshift — translate a timestamp into your roster's zones

Usage:
  tzshift [TIME [DATE] ZONE] [--to ZONE] [--no-color]
  tzshift list [--sort=name|offset]

Examples:
  tzshift 14:30 Asia/Kolkata        wall-clock time in a zone (date = today)
  tzshift 14:30 india               using a [zones] alias
  tzshift 2026-06-09 14:30 UTC      explicit date
  tzshift 1749571200                epoch seconds
  tzshift                           current time across your roster
  tzshift 14:30 UTC --to Europe/Berlin   add a one-off zone
  tzshift list                      known zones (east-most first) + your abbreviations
  tzshift list --sort=name          ...sorted alphabetically instead

Config: ~/.config/tzshift/config.toml  ([zones] and optional [abbreviations])
`
