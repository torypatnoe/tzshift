// Package tz is the pure, I/O-free translation engine for tzshift:
// source-zone resolution, timestamp parsing, DST-correct conversion, and
// east-most-first ordering. Keeping it free of config/CLI concerns lets the
// logic be reused (e.g. a future macOS widget) and table-tested directly.
package tz

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	// Embed the IANA tz database in the binary so conversions are identical
	// across host OSes and need no system zoneinfo (Spec AC5/AC12/AC14).
	_ "time/tzdata"
)

// Entry is one roster row: a user-chosen label mapped to an IANA zone.
type Entry struct {
	Label string
	Zone  string // IANA identifier, or "Local"/"UTC"
}

// ResolveZone resolves a source/target token to a Location.
// Resolution order (Spec AC8): roster alias -> user abbreviation -> literal
// IANA name. A token that matches none is a hard error (AC9) -- tzshift ships
// no abbreviation guesses.
func ResolveZone(token string, aliases, abbrev map[string]string) (*time.Location, string, error) {
	if iana, ok := aliases[token]; ok {
		return loadOr(token, iana, "alias")
	}
	if iana, ok := abbrev[token]; ok {
		return loadOr(token, iana, "abbreviation")
	}
	loc, err := time.LoadLocation(token)
	if err != nil {
		return nil, "", fmt.Errorf("unknown zone %q: use an IANA name (e.g. Asia/Kolkata), a roster alias, or define it under [abbreviations]", token)
	}
	return loc, token, nil
}

func loadOr(token, iana, kind string) (*time.Location, string, error) {
	loc, err := time.LoadLocation(iana)
	if err != nil {
		return nil, "", fmt.Errorf("%s %q points at %q, which is not a valid IANA zone", kind, token, iana)
	}
	return loc, iana, nil
}

// Kind distinguishes the three input shapes.
type Kind int

const (
	KindNow   Kind = iota // no timestamp -> current time (AC13)
	KindEpoch             // bare integer -> Unix seconds (AC3)
	KindWall              // HH:MM [date] zone (AC1/AC2)
)

// Request is a parsed-but-not-yet-resolved command input.
type Request struct {
	Kind             Kind
	Epoch            int64
	Hour, Minute     int
	Year, Month, Day int
	DateGiven        bool
	ZoneToken        string
}

var (
	reTime = regexp.MustCompile(`^([0-9]{1,2}):([0-9]{2})$`)
	reDate = regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})$`)
)

// ParseArgs classifies positional arguments into a Request. It does not touch
// config or the clock; resolution happens in Instant.
func ParseArgs(args []string) (*Request, error) {
	if len(args) == 0 {
		return &Request{Kind: KindNow}, nil
	}
	if len(args) == 1 {
		if n, err := strconv.ParseInt(args[0], 10, 64); err == nil {
			return &Request{Kind: KindEpoch, Epoch: n}, nil
		}
	}

	r := &Request{Kind: KindWall}
	timeSet := false
	for _, a := range args {
		switch {
		case reTime.MatchString(a):
			if timeSet {
				return nil, fmt.Errorf("more than one time value (%q)", a)
			}
			m := reTime.FindStringSubmatch(a)
			hh, _ := strconv.Atoi(m[1])
			mm, _ := strconv.Atoi(m[2])
			if hh > 23 || mm > 59 {
				return nil, fmt.Errorf("invalid time %q (expected HH:MM, 00:00-23:59)", a)
			}
			r.Hour, r.Minute, timeSet = hh, mm, true
		case reDate.MatchString(a):
			if r.DateGiven {
				return nil, fmt.Errorf("more than one date (%q)", a)
			}
			m := reDate.FindStringSubmatch(a)
			y, _ := strconv.Atoi(m[1])
			mo, _ := strconv.Atoi(m[2])
			d, _ := strconv.Atoi(m[3])
			if mo < 1 || mo > 12 || d < 1 || d > 31 {
				return nil, fmt.Errorf("invalid date %q (expected YYYY-MM-DD)", a)
			}
			r.Year, r.Month, r.Day, r.DateGiven = y, mo, d, true
		default:
			if r.ZoneToken != "" {
				return nil, fmt.Errorf("unexpected argument %q (source zone already set to %q)", a, r.ZoneToken)
			}
			r.ZoneToken = a
		}
	}
	if !timeSet {
		return nil, fmt.Errorf("missing time (expected HH:MM)")
	}
	// A missing source zone is allowed: an empty ZoneToken means "interpret the
	// time in the host's local zone" (resolved in Instant).
	return r, nil
}

// Resolved is the outcome of resolving a Request against config + clock.
type Resolved struct {
	Instant     time.Time // the absolute instant to translate
	SourceZone  string    // IANA name of the source (wall mode), else ""
	DateAssumed bool      // true when "today" was assumed (AC2)
}

// Instant resolves a Request to an absolute instant.
func (r *Request) Instant(aliases, abbrev map[string]string, now time.Time) (Resolved, error) {
	switch r.Kind {
	case KindNow:
		return Resolved{Instant: now}, nil
	case KindEpoch:
		return Resolved{Instant: time.Unix(r.Epoch, 0).UTC()}, nil
	case KindWall:
		loc, name := time.Local, "Local" // no source zone given -> local
		if r.ZoneToken != "" {
			var err error
			loc, name, err = ResolveZone(r.ZoneToken, aliases, abbrev)
			if err != nil {
				return Resolved{}, err
			}
		}
		y, mo, d := r.Year, time.Month(r.Month), r.Day
		assumed := false
		if !r.DateGiven {
			ny, nm, nd := now.In(loc).Date()
			y, mo, d, assumed = ny, nm, nd, true
		}
		return Resolved{
			Instant:     time.Date(y, mo, d, r.Hour, r.Minute, 0, 0, loc),
			SourceZone:  name,
			DateAssumed: assumed,
		}, nil
	}
	return Resolved{}, fmt.Errorf("unknown request kind")
}

// Row is one rendered output line: an instant expressed in one zone, with its
// computed abbreviation, day-offset vs the local row, and local/source markers.
type Row struct {
	Label         string
	Zone          string
	Time          time.Time
	Abbr          string
	OffsetSeconds int
	DayOffset     int // calendar days vs the local (host) date
	IsLocal       bool
	IsSource      bool
}

// Rows expresses instant across entries, marks the local and source rows,
// computes each row's day-offset vs the source zone's date, and sorts east-most
// first (AC15/AC16). sourceZone is the IANA zone the user typed, "Local" (no
// zone given), or "" (epoch / current-time modes); for the latter two the
// day-offset anchor falls back to the host's local date.
func Rows(instant time.Time, entries []Entry, sourceZone string) []Row {
	anchorLoc := time.Local
	if sourceZone != "" && sourceZone != "Local" {
		if loc, err := time.LoadLocation(sourceZone); err == nil {
			anchorLoc = loc
		}
	}
	refDate := dateOnly(instant.In(anchorLoc))
	_, localOff := instant.In(time.Local).Zone()

	rows := make([]Row, 0, len(entries))
	for _, e := range entries {
		loc, err := time.LoadLocation(e.Zone)
		if err != nil {
			continue // entries are validated before reaching here
		}
		lt := instant.In(loc)
		abbr, off := lt.Zone()
		rows = append(rows, Row{
			Label:         e.Label,
			Zone:          e.Zone,
			Time:          lt,
			Abbr:          abbr,
			OffsetSeconds: off,
			DayOffset:     daysBetween(refDate, dateOnly(lt)),
			IsSource:      sourceZone != "" && e.Zone == sourceZone,
		})
	}
	markLocal(rows, localOff)
	if sourceZone == "Local" {
		// No source zone was given: the source is the local row(s).
		for i := range rows {
			if rows[i].IsLocal {
				rows[i].IsSource = true
			}
		}
	}

	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].OffsetSeconds != rows[j].OffsetSeconds {
			return rows[i].OffsetSeconds > rows[j].OffsetSeconds // largest offset (east) first
		}
		return rows[i].Label < rows[j].Label
	})
	return rows
}

// markLocal flags the user's own row: a label of "you"/"local" wins; failing
// that, any row sharing the host-local UTC offset.
func markLocal(rows []Row, localOff int) {
	found := false
	for i := range rows {
		if strings.EqualFold(rows[i].Label, "you") || strings.EqualFold(rows[i].Label, "local") {
			rows[i].IsLocal = true
			found = true
		}
	}
	if found {
		return
	}
	for i := range rows {
		if rows[i].OffsetSeconds == localOff {
			rows[i].IsLocal = true
		}
	}
}

// SortedByOffset returns zone names ordered east-most first (largest UTC offset
// at the given instant), tie-broken by name. Used by `tzshift list --east`.
func SortedByOffset(names []string, at time.Time) []string {
	type zo struct {
		name string
		off  int
	}
	zos := make([]zo, 0, len(names))
	for _, n := range names {
		off := 0
		if loc, err := time.LoadLocation(n); err == nil {
			_, off = at.In(loc).Zone()
		}
		zos = append(zos, zo{n, off})
	}
	sort.SliceStable(zos, func(i, j int) bool {
		if zos[i].off != zos[j].off {
			return zos[i].off > zos[j].off
		}
		return zos[i].name < zos[j].name
	})
	out := make([]string, len(zos))
	for i, z := range zos {
		out[i] = z.name
	}
	return out
}

func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func daysBetween(ref, d time.Time) int {
	return int(d.Sub(ref).Round(24*time.Hour).Hours() / 24)
}
