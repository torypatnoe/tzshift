# Project: tztool

> Forma project-level artifact. Ties all cycles together. Updated after each cycle's Measure stage.

**DRI:** Tory Patnoe
**Customer Narrative:** [customer-narrative.md](customer-narrative.md)
**Living Shape:** [shape.md](shape.md)
**Living Spec:** [spec.md](spec.md)
**Status:** Active — Cycle 1 Shape + Spec drafted (Go, cross-unix), ready for Tickets
**Last updated:** 2026-06-19

---

## Milestone map

Per Forma, only the current milestone is fully resolved. Future milestones stay deliberately rough until the preceding cycle ships and is measured.

### M1 — `tzshift` CLI translator (current; shaped + specced in Cycle 1)
A cross-unix (Linux/macOS) command-line tool, `tzshift`, that takes any timestamp in any timezone and instantly translates it to the user's time and their configured team/system zones. Subcommand-oriented (`show` default, `list`). Built in Go. The core translation logic and zone-roster concept are proven here.

### M2 — macOS widget (rough)
Glanceable team-time awareness on the desktop. Built natively (Swift) when validated, reusing M1's translation *logic* and roster model — not the Go binary itself. Shape after M1's Measure.

### M3 — iOS app (hypothesis)
Translation + team awareness away from the desk. Untested assumption that mobile is where SREs need this; M1/M2 learning will confirm or kill.

### M4 — Android app (hypothesis)
Same as M3 for Android. Architecture/code-sharing decisions deferred until mobile demand is validated.

## Project hill chart

**Left side (figuring it out).**

- Understood: the customer, the core problem (timestamps posted in others' zones), terminal-first entry point.
- Unknown: what input formats matter most in practice, how the team/zone roster should be configured, whether glanceable awareness (M2) or arbitrary translation (M1) drives more daily use, whether mobile is needed at all.

## Cycle log

| Cycle | Milestone | Cycle record | Outcome / learning note |
|---|---|---|---|
| 1 | M1 — CLI translator | [bet placed](cycles/cycle-1.md) | — |
| 2 | M2 — MacOS Desktop Widget | not started | — |
