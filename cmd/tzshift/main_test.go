package main

import (
	"testing"
	"time"

	"github.com/torypatnoe/tztools/internal/tz"
)

func TestEnsureZone(t *testing.T) {
	base := []tz.Entry{{Label: "ny", Zone: "America/New_York"}}

	// Same IANA zone already present -> no addition (even via a different label).
	got := ensureZone(base, tz.Entry{Label: "Asia/Kolkata", Zone: "America/New_York"})
	if len(got) != 1 {
		t.Fatalf("expected no addition, got %d entries", len(got))
	}

	// New zone -> appended.
	got = ensureZone(base, tz.Entry{Label: "Asia/Kolkata", Zone: "Asia/Kolkata"})
	if len(got) != 2 || got[1].Zone != "Asia/Kolkata" {
		t.Fatalf("expected Asia/Kolkata appended, got %+v", got)
	}
}

func TestEnsureLocal(t *testing.T) {
	instant := time.Now()

	// No you/local label and (assumed) no offset match -> a you row is added.
	got := ensureLocal(nil, instant)
	if len(got) != 1 || got[0].Label != "you" || got[0].Zone != "Local" {
		t.Fatalf("expected a synthetic you/Local row, got %+v", got)
	}

	// An existing you label suppresses the addition.
	withYou := []tz.Entry{{Label: "you", Zone: "America/Los_Angeles"}}
	if got := ensureLocal(withYou, instant); len(got) != 1 {
		t.Fatalf("you-labelled roster should not gain a row, got %d", len(got))
	}
}
