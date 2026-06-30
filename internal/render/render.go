// Package render turns engine rows into the human-first, pipe-safe output:
// aligned columns, date on every row, cross-date markers, and TTY-gated color.
package render

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/torypatnoe/tztools/internal/tz"
)

const (
	dateFmt = "2006-01-02"
	timeFmt = "15:04"

	ansiHighlight = "\x1b[33m" // yellow: a row on a different calendar day
	ansiReset     = "\x1b[0m"
)

// Options controls presentation. Color must be false unless stdout is a TTY.
type Options struct {
	Color bool
}

// Show writes one aligned row per result. Rows read label first, then date and
// time, then UTC offset; the zone abbreviation is not shown. There is no header
// line — the source zone is identified inline by a ←source marker.
func Show(w io.Writer, rows []tz.Row, opts Options) {
	labelW := 0
	for _, r := range rows {
		if len(r.Label) > labelW {
			labelW = len(r.Label)
		}
	}
	for _, r := range rows {
		fmt.Fprintln(w, formatRow(r, labelW, opts.Color))
	}
}

func formatRow(r tz.Row, labelW int, color bool) string {
	line := fmt.Sprintf("%-*s  %s %s %s", labelW, r.Label, r.Time.Format(dateFmt), r.Time.Format(timeFmt), formatOffset(r.OffsetSeconds))
	if m := dayMarker(r.DayOffset); m != "" {
		line += "  " + m
	}
	if r.IsLocal {
		line += "  ←you"
	}
	if r.IsSource {
		line += "  ←source"
	}
	if color && r.DayOffset != 0 { // whole line highlighted on a cross-date row
		line = ansiHighlight + line + ansiReset
	}
	return line
}

// dayMarker is a plain-text day delta that survives a pipe (Spec AC16).
func dayMarker(n int) string {
	switch {
	case n > 0:
		return fmt.Sprintf("+%d", n)
	case n < 0:
		return fmt.Sprintf("%d", n)
	default:
		return ""
	}
}

// formatOffset renders a UTC offset in seconds as ±HH:MM.
func formatOffset(sec int) string {
	sign := "+"
	if sec < 0 {
		sign, sec = "-", -sec
	}
	return fmt.Sprintf("%s%02d:%02d", sign, sec/3600, (sec%3600)/60)
}

// isLetterCode reports whether an abbreviation is a real letter code (e.g. MDT)
// rather than a numeric pseudo-abbreviation (e.g. +0530), which Go returns for
// zones with no named abbreviation. Numeric ones always begin with + or -.
func isLetterCode(abbr string) bool {
	return abbr != "" && abbr[0] != '+' && abbr[0] != '-'
}

// List writes the reference output: every known IANA zone with its current
// abbreviation, then the user's configured [abbreviations] (Spec AC14).
func List(w io.Writer, zoneNames []string, abbrev map[string]string, now time.Time) {
	fmt.Fprintln(w, "ZONES  (identifier, current UTC offset, abbreviation)")
	nameW := 0
	for _, z := range zoneNames {
		if len(z) > nameW {
			nameW = len(z)
		}
	}
	for _, z := range zoneNames {
		offset, abbr := "", ""
		if loc, err := time.LoadLocation(z); err == nil {
			ab, off := now.In(loc).Zone()
			offset = formatOffset(off)
			if isLetterCode(ab) { // skip numeric pseudo-abbreviations (e.g. +0530)
				abbr = ab
			}
		}
		line := fmt.Sprintf("  %-*s  %-6s  %s", nameW, z, offset, abbr)
		fmt.Fprintln(w, strings.TrimRight(line, " "))
	}

	fmt.Fprintln(w)
	if len(abbrev) == 0 {
		fmt.Fprintln(w, "ABBREVIATIONS  (none configured -- add a [abbreviations] table to config.toml)")
		return
	}
	fmt.Fprintln(w, "ABBREVIATIONS  (your config: shortcut -> zone)")
	keys := make([]string, 0, len(abbrev))
	for k := range abbrev {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	abW := 0
	for _, k := range keys {
		if len(k) > abW {
			abW = len(k)
		}
	}
	for _, k := range keys {
		fmt.Fprintf(w, "  %-*s  %s\n", abW, k, abbrev[k])
	}
}
