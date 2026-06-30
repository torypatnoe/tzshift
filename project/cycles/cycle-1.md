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

## Tickets

Derived from the [Spec](../spec.md)'s acceptance criteria (AC1–AC16). Lean stdlib build: stdlib `flag` + manual dispatch, `BurntSushi/toml`, `golang.org/x/term`, embedded IANA DB via `time/tzdata`. Package layout: `cmd/tzshift` + pure `internal/tz` engine + `internal/config` + `internal/render`.

**Phase 0 — Scaffold**
- T1: `go mod init`, package layout, `time/tzdata` blank import, buildable stub. → AC12

**Phase 1 — Engine (pure, table-tested)**
- T2: source-zone resolution (alias → user `[abbreviations]` → IANA; unresolved = error). → AC8, AC9
- T3: input parsing (`HH:MM [YYYY-MM-DD] <zone>`, bare int = epoch, empty = now; "date assumed" flag). → AC1, AC2, AC3, AC13
- T4: DST-correct conversion; result carries in-zone instant + day-offset vs. local. → AC5
- T5: order east-most first (UTC offset descending). → AC15

**Phase 2 — Config**
- T6: load `~/.config/tzshift/config.toml` (`[zones]`, `[abbreviations]`); malformed = error, data-only parse. → AC6, AC11
- T7: no-config fallback (local + UTC + America/Denver) + helpful message; never errors. → AC7

**Phase 3 — Render**
- T8: aligned rows (date, wall-clock, computed abbr, label); `←you` marker; plain text. → AC1, AC10
- T9: cross-date highlight — `+N`/`-N` marker always; color only on TTY; `--no-color`. → AC16, AC10

**Phase 4 — CLI**
- T10: dispatch (`show` default) + `--to <iana|alias>` ad-hoc row. → AC1, AC4
- T11: `list` (IANA zones + computed current abbr + user `[abbreviations]`). → AC14
- T12: exit-code sweep (0 ok; non-zero + stderr for bad input / unresolved zone / bad config). → AC11

**Phase 5 — Distribution & verification**
- T13: cross-compile matrix (linux/darwin × amd64/arm64), static binaries. → AC12
- T14: DST golden tests (spring-forward + fall-back) + end-to-end golden output. → AC5 + regression

## Hill chart

**Right side (making it happen)** — Spec finalized 2026-06-19 (Go, no shipped abbreviation mapping, date-aware east-most-first output; [spec.md](../spec.md)). Tickets cut from AC1–AC16; build under way. Remaining unknowns being resolved in code: cross-date highlight ergonomics off-TTY, roster shape (flat `[zones]` sufficient for M1), DST golden coverage.

## The Measure

*(written at cycle end)*

- **Result against the goal:** —
- **Learning note:** —
- **Next cycle input:** —
