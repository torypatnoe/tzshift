package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/torypatnoe/tztools/internal/tz"
)

// rows mirror the Shape's date-shift narrative block. Times are built in UTC so
// the formatted clock/date is exactly the field values, independent of host TZ.
func sampleRows() []tz.Row {
	at := func(y, mo, d, h, mi int) time.Time {
		return time.Date(y, time.Month(mo), d, h, mi, 0, 0, time.UTC)
	}
	return []tz.Row{
		{Label: "tokyo-team", OffsetSeconds: 9 * 3600, Time: at(2026, 6, 21, 3, 0), DayOffset: 1},
		{Label: "india", OffsetSeconds: 5*3600 + 1800, Time: at(2026, 6, 20, 23, 30), IsSource: true},
		{Label: "legacy-billing", OffsetSeconds: 0, Time: at(2026, 6, 20, 18, 0)},
		{Label: "ny-dc", OffsetSeconds: -4 * 3600, Time: at(2026, 6, 20, 14, 0)},
		{Label: "you", OffsetSeconds: -7 * 3600, Time: at(2026, 6, 20, 11, 0), IsLocal: true},
	}
}

func TestShowGolden(t *testing.T) {
	var buf bytes.Buffer
	Show(&buf, sampleRows(), Options{Color: false})
	want := strings.Join([]string{
		"tokyo-team      2026-06-21 03:00 +09:00  +1",
		"india           2026-06-20 23:30 +05:30  ←source",
		"legacy-billing  2026-06-20 18:00 +00:00",
		"ny-dc           2026-06-20 14:00 -04:00",
		"you             2026-06-20 11:00 -07:00  ←you",
		"",
	}, "\n")
	if got := buf.String(); got != want {
		t.Errorf("Show output mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestShowColorOnlyOnCrossDate(t *testing.T) {
	var buf bytes.Buffer
	Show(&buf, sampleRows(), Options{Color: true})
	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		crossDate := strings.HasPrefix(line, ansiHighlight+"tokyo-team")
		wrapped := strings.HasPrefix(line, ansiHighlight) && strings.HasSuffix(line, ansiReset)
		if strings.Contains(line, "tokyo-team") {
			if !wrapped {
				t.Errorf("cross-date row should be fully highlighted: %q", line)
			}
			_ = crossDate
		} else if strings.Contains(line, ansiHighlight) {
			t.Errorf("same-date row should not be highlighted: %q", line)
		}
	}
}

func TestListZoneFormat(t *testing.T) {
	at, _ := time.Parse(time.RFC3339, "2026-01-15T00:00:00Z") // winter: Denver = MST
	var buf bytes.Buffer
	List(&buf, []string{"UTC", "Asia/Kolkata", "America/Denver", "Pacific/Kiritimati"}, nil, at)
	out := buf.String()

	for _, want := range []string{"+00:00  UTC", "+05:30  IST", "-07:00  MST"} {
		if !strings.Contains(out, want) {
			t.Errorf("list output missing %q\n%s", want, out)
		}
	}
	// A numeric pseudo-abbreviation (+14) is suppressed: the row ends at the offset.
	if !strings.Contains(out, "+14:00\n") {
		t.Errorf("expected Pacific/Kiritimati row to end at its offset\n%s", out)
	}
	if strings.Contains(out, "+14:00  +14") {
		t.Errorf("numeric pseudo-abbreviation should be suppressed\n%s", out)
	}
}

func TestFormatOffset(t *testing.T) {
	cases := map[int]string{0: "+00:00", 5*3600 + 1800: "+05:30", -7 * 3600: "-07:00", 14 * 3600: "+14:00"}
	for sec, want := range cases {
		if got := formatOffset(sec); got != want {
			t.Errorf("formatOffset(%d) = %q, want %q", sec, got, want)
		}
	}
}

func TestShowSourceIsAlsoLocal(t *testing.T) {
	// When the row is both source and local, both markers appear.
	at := time.Date(2026, 6, 20, 11, 0, 0, 0, time.UTC)
	rows := []tz.Row{{Label: "you", OffsetSeconds: -7 * 3600, Time: at, IsLocal: true, IsSource: true}}
	var buf bytes.Buffer
	Show(&buf, rows, Options{Color: false})
	if got := buf.String(); !strings.Contains(got, "←you  ←source") {
		t.Errorf("expected both markers on one row, got: %s", got)
	}
}

func TestDayMarker(t *testing.T) {
	cases := map[int]string{0: "", 1: "+1", 2: "+2", -1: "-1"}
	for n, want := range cases {
		if got := dayMarker(n); got != want {
			t.Errorf("dayMarker(%d) = %q, want %q", n, got, want)
		}
	}
}
