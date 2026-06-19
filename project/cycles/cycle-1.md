# Cycle 1: CLI Translator (M1)

> Forma cycle record — the Bet (written at cycle start), the hill chart position (updated during Build), and the Measure (written at cycle end) live together so the measure is always judged against the bet that was made.

**Project:** [project.md](../project.md)
**Living Shape:** [shape.md](../shape.md) (as of 2026-06-10)

---

## The Bet

- **Milestone:** M1 — CLI translator
- **Goal:** An SRE can translate any posted timestamp (wall-clock + zone, with optional date, or epoch seconds) into their own time and their configured roster's times with a single `tzshift` command, trusting the result across DST boundaries.
- **Appetite:** Small — 1–2 weeks of one engineer
- **DRI:** Tory Patnoe (confirmed)

### In scope

- `tzshift show` (default) command: wall-clock + zone, explicit date, epoch seconds, `--to` one-off target; bare invocation shows current time across roster
- `tzshift list` reference command: zones + abbreviations
- TOML roster at `~/.config/tzshift/config.toml` (`[zones]` table) with `[abbreviations]` override
- Opinionated default abbreviation mapping (shipped, not OS-derived); hard error with candidates when ambiguous
- Helpful no-config fallback (local + UTC + America/Denver)
- Go, cross-unix (Linux + macOS), build-from-source + statically linked release binaries

### NOT in this cycle

- macOS widget, iOS, Android (M2–M4)
- Natural-language time parsing
- On-call schedules, calendars, notifications (project non-goals)
- Homebrew/package-manager distribution
- Native macOS widget / iOS / Android surfaces (M2–M4)
- Watching/auto-translating clipboard or log streams
- Roster groups, per-zone working hours, or any team metadata beyond name → zone

## Hill chart

**Left side (figuring it out)** — Spec drafted 2026-06-18 ([spec.md](../spec.md)). Unknowns being resolved: Swift CLI cold-start, abbreviation-mapping coverage, roster shape. Not yet on the right side — no code written.

## The Measure

*(written at cycle end)*

- **Result against the goal:** —
- **Learning note:** —
- **Next cycle input:** —
