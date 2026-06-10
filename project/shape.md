# Shape: tztool

> Forma living Shape — describes the full product as currently understood, deliberately larger than any single cycle's work. What a specific cycle commits to is recorded in that cycle's record (`cycles/cycle-N.md`). Per-cycle facts (appetite, scope boundary, DRI confirmation) live only in the Bet.

## Changelog
- 2026-06-10: Initial living Shape, converted from the Cycle 1 per-cycle shape draft under the updated Forma workflow. Covers the CLI translator in detail; later surfaces (widget, mobile) noted as rough.

**Project:** [project.md](project.md)
**DRI:** Tory Patnoe
**Last updated:** 2026-06-10

---

## Rough solution narrative

An SRE is reading an incident channel. A teammate in Bangalore writes "rollback finished at 14:30." The SRE switches to their terminal (already open) and types:

```
$ tz 14:30 IST
14:30 IST (Asia/Kolkata)
→ 02:00 PDT  you
→ 05:00 EDT  ny-dc
→ 09:00 UTC  legacy-billing
→ 18:00 JST  tokyo-team
```

One command, instant answer, no mental math. The rows come from a roster the SRE set up once:

```toml
# ~/.config/tztool/config.toml
[zones]
you = "America/Los_Angeles"
ny-dc = "America/New_York"
legacy-billing = "UTC"
tokyo-team = "Asia/Tokyo"
```

The roster names the zones the SRE actually cares about — people, datacenters, legacy systems — in their own words. Output is plain text: pipeable, greppable, scriptable.

Variations the narrative must cover:

- `tz 14:30 IST` — wall-clock time + zone, the core case
- `tz 2026-06-09 14:30 IST` — explicit date (DST correctness depends on it; defaults to today)
- `tz 1749571200` — epoch seconds, the log-file case
- `tz 14:30 IST --to Europe/Berlin` — one-off target zone not in the roster
- No config file yet → helpful message showing how to create one; falls back to local time + UTC

### Later surfaces (rough — shaped after the CLI ships and is measured)

- **macOS widget (M2):** glanceable team-time awareness on the desktop, reusing the roster and translation core.
- **iOS / Android (M3/M4):** hypotheses only; mobile demand unvalidated.

## Hard design decisions

- **Zone abbreviations are ambiguous.** IST is India, Israel, and Ireland; CST is China and US Central. The tool ships an opinionated default mapping for common SRE abbreviations, overridable in config (`[abbreviations]` table). Never guess silently when truly ambiguous — error with the candidates.
- **DST correctness over cleverness.** All conversions go through the IANA tz database (via Swift's `Foundation.TimeZone`); the date matters, so date-less input assumes "today" and says so in output.
- **Output is for humans first, pipes second.** Aligned columns, arrow prefixes; a `--no-color`/dumb-pipe degradation is acceptable to defer if time runs short.
- **Swift as the implementation language.** Native path to the macOS widget (M2) and iOS (M3), solid CLI via swift-argument-parser. Android (M4) would need a rewrite, but M4 is an unvalidated hypothesis — see alternatives.

## Identified unknowns

1. **Input format breadth** — which timestamp formats appear most in real channels/logs? Start with the four above; let usage (Measure) drive the next set.
2. **Abbreviation mapping** — is an opinionated default mapping good enough, or does ambiguity bite immediately?
3. **Swift CLI ergonomics** — startup time and distribution of a Swift binary (works on macOS; Linux support untested).
4. **Roster shape** — is a flat `[zones]` table enough, or do users immediately want groups/ordering?
5. **Awareness vs. translation** — does glanceable team-time awareness (widget territory) or arbitrary translation drive more daily use? Determines how M2 gets shaped.

## Alternatives considered (rejected)

- **Interactive TUI dashboard** — richer, but slower to build, not pipeable, and doesn't serve the "translate this one timestamp" moment the narrative centers on. Glanceable awareness is the widget's territory.
- **Rust core with FFI for future platforms** — maximum reuse across all four surfaces, but mobile milestones are unvalidated hypotheses; paying FFI plumbing cost now is over-speccing future cycles. Swift gives a native path to the next likeliest bets (macOS widget, iOS).
- **Natural-language input ("tomorrow 2pm Bangalore time")** — attractive, but parsing ambiguity is a tarpit. Structured input until usage data argues otherwise.
