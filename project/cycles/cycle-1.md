# Cycle 1: CLI Translator (M1)

> Forma cycle record — the Bet (written at cycle start), the hill chart position (updated during Build), and the Measure (written at cycle end) live together so the measure is always judged against the bet that was made.

**Project:** [project.md](../project.md)
**Living Shape:** [shape.md](../shape.md) (as of 2026-06-10)

---

## The Bet

- **Milestone:** M1 — CLI translator
- **Goal:** An SRE can translate any posted timestamp (wall-clock + zone, with optional date, or epoch seconds) into their own time and their configured roster's times with a single `tz` command, trusting the result across DST boundaries.
- **Appetite:** Small — 1–2 weeks of one engineer
- **DRI:** Tory Patnoe (confirmed)

### In scope

- One-shot `tz` command: wall-clock + zone, explicit date, epoch seconds, `--to` one-off target
- TOML roster at `~/.config/tztool/config.toml` (`[zones]` table) with `[abbreviations]` override
- Opinionated default abbreviation mapping; hard error with candidates when ambiguous
- Helpful no-config fallback (local time + UTC)
- Swift, macOS, build-from-source + release binary

### NOT in this cycle

- macOS widget, iOS, Android (M2–M4)
- Natural-language time parsing
- On-call schedules, calendars, notifications (project non-goals)
- Homebrew/package-manager distribution
- Linux support
- Watching/auto-translating clipboard or log streams
- Roster groups, per-zone working hours, or any team metadata beyond name → zone

## Hill chart

**Not started** — bet placed 2026-06-10; cycle begins at Spec.

## The Measure

*(written at cycle end)*

- **Result against the goal:** —
- **Learning note:** —
- **Next cycle input:** —
