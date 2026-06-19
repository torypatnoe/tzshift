# Shape: tzshift

> Forma living Shape — describes the full product as currently understood, deliberately larger than any single cycle's work. What a specific cycle commits to is recorded in that cycle's record (`cycles/cycle-N.md`). Per-cycle facts (appetite, scope boundary, DRI confirmation) live only in the Bet.

## Changelog
- 2026-06-19: Tool renamed `tz` → **`tzshift`** (`tz` is taken — see research). Subcommand model adopted (`show` default, `list`). **Language changed Swift → Go** after the Customer Narrative committed to Linux/macOS/unix support in phase 1: Swift's justification rested entirely on hypothesis milestones (M2/M3), and the known phase-1 requirement is a cross-unix CLI — Go's strength, Swift's weakness (retro lesson #3). DST decision decoupled from a specific library (states the requirement, not the runtime — retro lesson #4). Config dir `~/.config/tztool` → `~/.config/tzshift`. No-config fallback set to local + UTC + America/Denver.
- 2026-06-10: Initial living Shape, converted from the Cycle 1 per-cycle shape draft under the updated Forma workflow. Covers the CLI translator in detail; later surfaces (widget, mobile) noted as rough.

**Project:** [project.md](project.md)
**DRI:** Tory Patnoe
**Last updated:** 2026-06-19

---

## Rough solution narrative

**tzshift is subcommand-oriented.** Two commands matter at phase 1: `tzshift show` (the workhorse) and `tzshift list` (reference). `show` is the default — if no subcommand matches the first argument, the input is handed to `show` along with any parameters.

### tzshift show — the workhorse

An SRE is reading an incident channel. A teammate in Bangalore writes "rollback finished at 14:30." The SRE switches to their terminal (already open) and types:

```
$ tzshift 14:30 IST
14:30 IST (Asia/Kolkata)
→ 02:00 PDT  you
→ 05:00 EDT  ny-dc
→ 09:00 UTC  legacy-billing
→ 18:00 JST  tokyo-team
```

One command, instant answer, no mental math. The rows come from a roster the SRE set up once:

```toml
# ~/.config/tzshift/config.toml
[zones]
you = "America/Los_Angeles"
ny-dc = "America/New_York"
legacy-billing = "UTC"
tokyo-team = "Asia/Tokyo"
```

The roster names the zones the SRE actually cares about — people, datacenters, legacy systems — in their own words. Output is plain text: pipeable, greppable, scriptable.

Variations the narrative must cover:

- `tzshift show` (or bare `tzshift`) with no timestamp — outputs the current time across the configured roster zones (glanceable "what time is it for my team right now").
- `tzshift 14:30 IST` — wall-clock time + zone, the core case (no subcommand → routed to `show`).
- `tzshift 2026-06-09 14:30 IST` — explicit date (DST correctness depends on it; defaults to today).
- `tzshift 1749571200` — epoch seconds, the log-file case.
- `tzshift 14:30 IST --to Europe/Berlin` — one-off target zone not in the roster.
- No config file yet → helpful message showing how to create one; falls back to **local time + UTC + America/Denver**.

### tzshift list — reference

Lists the timezones tzshift knows about and the abbreviations it resolves. Because **zone abbreviations are ambiguous** (IST is India, Israel, and Ireland; CST is China and US Central), `list` is how an SRE discovers what tzshift will do with a given abbreviation before relying on it.

- `tzshift list` — enumerates IANA zone identifiers and each zone's current abbreviation, plus tzshift's shipped abbreviation → zone mapping.

**OS-agnostic caveat (a real constraint, not a detail):** zone *identifiers* (`America/Denver`, `Asia/Kolkata`) come from the IANA database, which tzshift carries embedded so it does not depend on the host OS's zoneinfo layout. But there is **no canonical OS-provided "zone → abbreviation" table** — abbreviations are time-dependent (`America/Denver` is MST or MDT depending on the date) and ambiguous. So `list` shows each zone's *computed current* abbreviation, and the ambiguous-abbreviation resolution (IST → candidates) is tzshift's **own shipped opinionated mapping**, overridable in config — not something read from the OS.

### Later surfaces (rough — shaped after the CLI ships and is measured)

- **macOS widget (M2):** glanceable team-time awareness on the desktop. Built natively (Swift) when validated, reusing tzshift's proven translation logic — not the Go CLI binary itself.
- **iOS / Android (M3/M4):** hypotheses only; mobile demand unvalidated.

## Hard design decisions

- **The tool is `tzshift`, not `tz`.** `tz` is already taken (research: the oz/tz TUI world clock). `tzshift` is subcommand-oriented (`show`, `list`), with `show` as the default when the first argument isn't a known subcommand.
- **Go as the implementation language.** Phase 1 requires a CLI that runs across Linux, macOS, and other unix environments and is trivial to distribute — Go's single static binary and cross-compilation are the best fit. The language is chosen for phase 1's *validated* requirement, not for the hypothesis milestones (M2/M3); when the macOS widget is validated it will be built natively in Swift, reusing the translation logic. (Earlier Swift choice reversed — see Changelog and Alternatives.)
- **Zone data from the IANA database, carried embedded.** Zone identifiers are enumerated from the IANA tz database bundled in the binary, so behaviour does not vary with the host OS's zoneinfo. Abbreviations are computed per-instant from that data; the ambiguous-abbreviation → zone mapping is tzshift's own shipped, config-overridable data (there is no OS abbreviation list to read).
- **Zone abbreviations are ambiguous — never guess silently.** Ship an opinionated default mapping for common SRE abbreviations, overridable in config (`[abbreviations]` table). When truly ambiguous and unresolved by config, error with the candidates.
- **DST correctness over cleverness.** All conversions go through the IANA tz database; never hand-roll UTC offsets. The date matters, so date-less input assumes "today" and says so in output. *(The Spec names the concrete library; the Shape states only the requirement.)*
- **Output is for humans first, pipes second.** Aligned columns, arrow prefixes; a `--no-color`/dumb-pipe degradation is acceptable to defer if time runs short.

## Identified unknowns

1. **Input format breadth** — which timestamp formats appear most in real channels/logs? Start with the variations above; let usage (Measure) drive the next set.
2. **Abbreviation mapping** — is an opinionated default mapping good enough, or does ambiguity bite immediately?
3. **Distribution across unix variants** — Go gives static binaries for Linux/macOS easily; confirm "any unix" (e.g. BSD) targets and the release/cross-compile matrix actually needed.
4. **Roster shape** — is a flat `[zones]` table enough, or do users immediately want groups/ordering?
5. **Awareness vs. translation** — does glanceable team-time awareness (the `show`-with-no-timestamp / widget territory) or arbitrary translation drive more daily use? Determines how M2 gets shaped.

## Alternatives considered (rejected)

- **Swift** — was the original choice for its native path to the macOS widget (M2) and iOS (M3). Rejected for phase 1: those are unvalidated hypothesis milestones, and letting them drive the phase-1 language conflicts with the now-committed cross-unix requirement, where Swift is weakest (Linux works but "any unix" isn't first-class; distribution is heavier). Swift returns as the natural choice for M2 *when M2 is validated*.
- **Rust** — also excellent cross-platform static binaries, and a single core could FFI into iOS/macOS later. Rejected for phase 1 in favour of Go's faster CLI development and simpler toolchain; revisit only if a shared-core-across-all-surfaces need is actually validated.
- **Interactive TUI dashboard** — richer, but slower to build, not pipeable, and doesn't serve the "translate this one timestamp" moment the narrative centers on. Glanceable awareness is `show`'s / the widget's territory.
- **Natural-language input ("tomorrow 2pm Bangalore time")** — attractive, but parsing ambiguity is a tarpit. Structured input until usage data argues otherwise.
