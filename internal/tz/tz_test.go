package tz

import (
	"testing"
	"time"
)

func mustLoad(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("LoadLocation(%q): %v", name, err)
	}
	return loc
}

func TestResolveZone(t *testing.T) {
	aliases := map[string]string{"india": "Asia/Kolkata"}
	abbrev := map[string]string{"IST": "Asia/Kolkata"}

	cases := []struct {
		token   string
		want    string
		wantErr bool
	}{
		{"india", "Asia/Kolkata", false},        // roster alias
		{"IST", "Asia/Kolkata", false},          // user abbreviation
		{"Asia/Kolkata", "Asia/Kolkata", false}, // literal IANA
		{"UTC", "UTC", false},
		{"IST_typo", "", true}, // unknown -> hard error (AC9)
		{"PST", "", true},      // a 3-letter abbr we never shipped
	}
	for _, c := range cases {
		_, name, err := ResolveZone(c.token, aliases, abbrev)
		if c.wantErr {
			if err == nil {
				t.Errorf("ResolveZone(%q): expected error, got %q", c.token, name)
			}
			continue
		}
		if err != nil {
			t.Errorf("ResolveZone(%q): unexpected error %v", c.token, err)
		}
		if name != c.want {
			t.Errorf("ResolveZone(%q) = %q, want %q", c.token, name, c.want)
		}
	}
}

func TestParseArgs(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		kind    Kind
		check   func(*Request) bool
		wantErr bool
	}{
		{"empty -> now", nil, KindNow, nil, false},
		{"epoch", []string{"1749571200"}, KindEpoch, func(r *Request) bool { return r.Epoch == 1749571200 }, false},
		{"wall no date", []string{"14:30", "Asia/Kolkata"}, KindWall,
			func(r *Request) bool {
				return r.Hour == 14 && r.Minute == 30 && !r.DateGiven && r.ZoneToken == "Asia/Kolkata"
			}, false},
		{"wall with date", []string{"2026-06-09", "14:30", "UTC"}, KindWall,
			func(r *Request) bool { return r.DateGiven && r.Year == 2026 && r.Month == 6 && r.Day == 9 }, false},
		{"order independent", []string{"UTC", "14:30"}, KindWall,
			func(r *Request) bool { return r.ZoneToken == "UTC" && r.Hour == 14 }, false},
		{"missing zone -> local", []string{"14:30"}, KindWall,
			func(r *Request) bool { return r.Hour == 14 && r.ZoneToken == "" }, false},
		{"missing time", []string{"Asia/Kolkata"}, KindWall, nil, true},
		{"bad time", []string{"99:99", "UTC"}, KindWall, nil, true},
		{"two times", []string{"14:30", "15:30", "UTC"}, KindWall, nil, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, err := ParseArgs(c.args)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", r)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Kind != c.kind {
				t.Fatalf("kind = %v, want %v", r.Kind, c.kind)
			}
			if c.check != nil && !c.check(r) {
				t.Fatalf("check failed for %+v", r)
			}
		})
	}
}

// TestDSTConversion verifies fall-back and spring-forward boundaries for an
// instant converted into a zone (Spec AC5). Times are UTC so the test is
// independent of the host's local zone.
func TestDSTConversion(t *testing.T) {
	ny := mustLoad(t, "America/New_York")
	cases := []struct {
		utc      string
		wantClk  string
		wantAbbr string
	}{
		// Fall-back 2026-11-01: 02:00 EDT -> 01:00 EST (06:00Z boundary).
		{"2026-11-01T05:30:00Z", "01:30", "EDT"},
		{"2026-11-01T06:30:00Z", "01:30", "EST"}, // same wall clock, other side
		// Spring-forward 2026-03-08: 02:00 EST -> 03:00 EDT (07:00Z boundary).
		{"2026-03-08T06:30:00Z", "01:30", "EST"},
		{"2026-03-08T07:30:00Z", "03:30", "EDT"}, // 02:xx is skipped
	}
	for _, c := range cases {
		instant, _ := time.Parse(time.RFC3339, c.utc)
		local := instant.In(ny)
		if got := local.Format("15:04"); got != c.wantClk {
			t.Errorf("%s -> %s, want %s", c.utc, got, c.wantClk)
		}
		if abbr, _ := local.Zone(); abbr != c.wantAbbr {
			t.Errorf("%s -> abbr %s, want %s", c.utc, abbr, c.wantAbbr)
		}
	}
}

// TestWallClockDST checks that a date-bearing wall-clock input lands on the
// correct UTC instant on either side of DST.
func TestWallClockDST(t *testing.T) {
	aliases := map[string]string{}
	abbrev := map[string]string{}
	cases := []struct {
		args    []string
		wantUTC string
	}{
		{[]string{"2026-07-01", "14:30", "America/New_York"}, "2026-07-01T18:30:00Z"}, // EDT -4
		{[]string{"2026-01-01", "14:30", "America/New_York"}, "2026-01-01T19:30:00Z"}, // EST -5
	}
	for _, c := range cases {
		r, err := ParseArgs(c.args)
		if err != nil {
			t.Fatalf("ParseArgs(%v): %v", c.args, err)
		}
		res, err := r.Instant(aliases, abbrev, time.Now())
		if err != nil {
			t.Fatalf("Instant: %v", err)
		}
		if got := res.Instant.UTC().Format(time.RFC3339); got != c.wantUTC {
			t.Errorf("%v -> %s UTC, want %s", c.args, got, c.wantUTC)
		}
	}
}

func TestSortedByOffset(t *testing.T) {
	at, _ := time.Parse(time.RFC3339, "2026-06-20T00:00:00Z")
	in := []string{"America/New_York", "Asia/Tokyo", "UTC", "America/Los_Angeles"}
	got := SortedByOffset(in, at)
	want := []string{"Asia/Tokyo", "UTC", "America/New_York", "America/Los_Angeles"}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("position %d = %q, want %q (full: %v)", i, got[i], w, got)
		}
	}
}

func TestInstantNoZoneIsLocal(t *testing.T) {
	r, err := ParseArgs([]string{"14:30"})
	if err != nil {
		t.Fatalf("ParseArgs: %v", err)
	}
	res, err := r.Instant(map[string]string{}, map[string]string{}, time.Now())
	if err != nil {
		t.Fatalf("Instant: %v", err)
	}
	if res.SourceZone != "Local" {
		t.Errorf("SourceZone = %q, want Local", res.SourceZone)
	}
	if !res.DateAssumed {
		t.Errorf("expected date assumed for a date-less input")
	}
}

func TestInstantNowIsLocalSource(t *testing.T) {
	r, _ := ParseArgs(nil) // bare `tzshift` -> current time
	if r.Kind != KindNow {
		t.Fatalf("Kind = %v, want KindNow", r.Kind)
	}
	res, err := r.Instant(map[string]string{}, map[string]string{}, time.Now())
	if err != nil {
		t.Fatalf("Instant: %v", err)
	}
	if res.SourceZone != "Local" {
		t.Errorf("SourceZone = %q, want Local (current-time source is local)", res.SourceZone)
	}
}

func TestRowsAnchorToSource(t *testing.T) {
	instant, _ := time.Parse(time.RFC3339, "2026-06-22T04:30:00Z") // 10:00 in Kolkata
	entries := []Entry{
		{Label: "india", Zone: "Asia/Kolkata"},
		{Label: "denver", Zone: "America/Denver"},
	}
	rows := Rows(instant, entries, "Asia/Kolkata")
	for _, r := range rows {
		switch r.Zone {
		case "Asia/Kolkata":
			if r.DayOffset != 0 || !r.IsSource {
				t.Errorf("source row: DayOffset=%d IsSource=%v, want 0/true", r.DayOffset, r.IsSource)
			}
		case "America/Denver":
			if r.DayOffset != -1 { // 22:30 on the prior day vs the source date
				t.Errorf("denver DayOffset = %d, want -1", r.DayOffset)
			}
		}
	}
}

func TestRowsOrderingEastFirst(t *testing.T) {
	instant, _ := time.Parse(time.RFC3339, "2026-06-20T09:00:00Z")
	entries := []Entry{
		{Label: "you", Zone: "America/Los_Angeles"},
		{Label: "tokyo", Zone: "Asia/Tokyo"},
		{Label: "utc", Zone: "UTC"},
		{Label: "ny", Zone: "America/New_York"},
	}
	rows := Rows(instant, entries, "Asia/Tokyo")
	if !rows[0].IsSource || rows[1].IsSource {
		t.Errorf("expected only the tokyo row to be marked source")
	}
	wantOrder := []string{"tokyo", "utc", "ny", "you"} // descending UTC offset
	if len(rows) != len(wantOrder) {
		t.Fatalf("got %d rows, want %d", len(rows), len(wantOrder))
	}
	for i, w := range wantOrder {
		if rows[i].Label != w {
			t.Errorf("row %d = %q, want %q", i, rows[i].Label, w)
		}
	}
	// offsets must be non-increasing
	for i := 1; i < len(rows); i++ {
		if rows[i-1].OffsetSeconds < rows[i].OffsetSeconds {
			t.Errorf("ordering broken at %d: %d < %d", i, rows[i-1].OffsetSeconds, rows[i].OffsetSeconds)
		}
	}
}
