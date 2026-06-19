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
		{Label: "tokyo-team", Abbr: "JST", Time: at(2026, 6, 21, 3, 0), DayOffset: 1},
		{Label: "india", Abbr: "IST", Time: at(2026, 6, 20, 23, 30)},
		{Label: "legacy-billing", Abbr: "UTC", Time: at(2026, 6, 20, 18, 0)},
		{Label: "ny-dc", Abbr: "EDT", Time: at(2026, 6, 20, 14, 0)},
		{Label: "you", Abbr: "PDT", Time: at(2026, 6, 20, 11, 0), IsLocal: true},
	}
}

func TestShowGolden(t *testing.T) {
	var buf bytes.Buffer
	Show(&buf, "", sampleRows(), Options{Color: false})
	want := strings.Join([]string{
		"→ 2026-06-21 03:00 JST tokyo-team  +1",
		"→ 2026-06-20 23:30 IST india",
		"→ 2026-06-20 18:00 UTC legacy-billing",
		"→ 2026-06-20 14:00 EDT ny-dc",
		"→ 2026-06-20 11:00 PDT you  ←you",
		"",
	}, "\n")
	if got := buf.String(); got != want {
		t.Errorf("Show output mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestShowColorOnlyOnCrossDate(t *testing.T) {
	var buf bytes.Buffer
	Show(&buf, "", sampleRows(), Options{Color: true})
	out := buf.String()
	if !strings.Contains(out, ansiHighlight+"2026-06-21") {
		t.Error("expected the cross-date row to be highlighted")
	}
	// The same-date rows must not be wrapped in color.
	if strings.Contains(out, ansiHighlight+"2026-06-20") {
		t.Error("same-date rows should not be highlighted")
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
