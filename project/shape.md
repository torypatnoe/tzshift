# Shape: tzshift

> Forma living Shape — describes the full product as currently understood, deliberately larger than any single cycle's work. What a specific cycle commits to is recorded in that cycle's record (`cycles/cycle-N.md`). Per-cycle facts (appetite, scope boundary, DRI confirmation) live only in the Bet.

## Changelog
- 2026-06-22: Added **`tzshift config create`** (scaffold a starter config) and **`tzshift config list`** (print it). Dropped the leading `→` from `show` rows; cross-date rows now highlight the **whole line**; and the cross-date anchor is now the **source** zone (`←source`) rather than the local row.
- 2026-06-22: **Source zone is now optional.** `tzshift 14:30` (no zone) interprets the time in your local zone instead of erroring; the local `←you` row is then also the `←source`.
- 2026-06-22: **`show` drops the echo/header line.** The source zone is identified inline by a **`←source`** marker (a row that is both source and local shows `←you  ←source`); the assumed date is conveyed by the per-row dates rather than a header note.
- 2026-06-22: **`show` auto-includes your local time and the source zone.** The `+1`/`-1` day markers are anchored to your host's local date; that local row (`←you`) is now always shown — added automatically if your roster doesn't list your zone — so the markers always have a visible reference. The source zone you typed is likewise always shown. (Fixes the confusing case where every row showed `+1` against an off-screen anchor.)
- 2026-06-22: **`show` output reordered and abbreviation dropped.** Rows now read **label, date, time, UTC offset (`±HH:MM`)** — the three-letter zone abbreviation (PDT/JST) is no longer shown in `show` output; the numeric offset replaces it (clearer, unambiguous). `list` still shows letter-code abbreviations alongside the offset. Supersedes "Output abbreviations (PDT/JST) stay" in the entry below as it applies to `show` rows.
- 2026-06-19 (later, supersedes the abbreviation parts of the entry below): **No embedded abbreviation→zone mapping.** Dropped the shipped opinionated abbreviation table as a *source-zone resolver* — abbreviations are too ambiguous (IST = India/Israel/Ireland) and a shipped guess is exactly the silent-guess failure the tool exists to avoid. The source zone is now an **IANA name or a `[zones]` roster alias**; an optional user-defined `[abbreviations]` table is the escape hatch (the user's own data, still nothing embedded). Output also changes: every row now (a) shows the **date**, (b) **highlights cross-date rows** vs. the local row — color when stdout is a TTY, plus a `+1`/`-1` marker that survives a pipe, (c) is sorted **east-most (largest UTC offset) first** instead of config-file order, and (d) **marks the local `you` row**. Output abbreviations (PDT/JST) stay — computed per-instant *display*, never input resolvers. Affects AC1, AC4, AC6, AC8–AC10, AC14; adds AC15 (ordering) + AC16 (date display + cross-date highlight).
- 2026-06-19: Tool renamed `tz` → **`tzshift`** (`tz` is taken — see research). Subcommand model adopted (`show` default, `list`). **Language changed Swift → Go** after the Customer Narrative committed to Linux/macOS/unix support in phase 1: Swift's justification rested entirely on hypothesis milestones (M2/M3), and the known phase-1 requirement is a cross-unix CLI — Go's strength, Swift's weakness (retro lesson #3). DST decision decoupled from a specific library (states the requirement, not the runtime — retro lesson #4). Config dir `~/.config/tztool` → `~/.config/tzshift`. No-config fallback set to local + UTC + America/Denver.
- 2026-06-10: Initial living Shape, converted from the Cycle 1 per-cycle shape draft under the updated Forma workflow. Covers the CLI translator in detail; later surfaces (widget, mobile) noted as rough.

**Project:** [project.md](project.md)
**DRI:** Tory Patnoe
**Last updated:** 2026-06-19

---

## Rough solution narrative

**tzshift is subcommand-oriented.** The commands at phase 1: `tzshift show` (the workhorse), `tzshift list` (reference), and `tzshift config create|list` (scaffold / print the config file). `show` is the default — if no subcommand matches the first argument, the input is handed to `show` along with any parameters.

### tzshift show — the workhorse

An SRE is reading an incident channel. A teammate in Bangalore writes "rollback finished at 14:30." The SRE switches to their terminal (already open) and types:

```
$ tzshift 14:30 Asia/Kolkata
tokyo-team      2026-06-19 18:00 +09:00
india           2026-06-19 14:30 +05:30  ←source
legacy-billing  2026-06-19 09:00 +00:00
ny-dc           2026-06-19 05:00 -04:00
you             2026-06-19 02:00 -07:00  ←you
```

One command, instant answer, no mental math. The source zone is given as an **IANA name** (`Asia/Kolkata`) or a **roster alias** (`india`) — tzshift ships no three-letter input mapping, because an abbreviation like `IST` (India? Israel? Ireland?) is exactly the silent guess the tool refuses to make. The output rows come from a roster the SRE set up once:

```toml
# ~/.config/tzshift/config.toml
[zones]
tokyo-team     = "Asia/Tokyo"
india          = "Asia/Kolkata"
legacy-billing = "UTC"
ny-dc          = "America/New_York"
you            = "America/Los_Angeles"

# Optional escape hatch: define your OWN abbreviations as source zones.
# This is your data, not a shipped default — so it is never an ambiguous guess.
[abbreviations]
IST = "Asia/Kolkata"
```

The `[zones]` roster names the zones the SRE actually cares about — people, datacenters, legacy systems — in their own words (`alias = "IANA/Zone"`). With the optional `[abbreviations]` table, `tzshift 14:30 IST` works too — because *this user* declared what IST means. Output is plain text: pipeable, greppable, scriptable.

**Reading the output:**

- Rows are sorted **east-most first** (largest UTC offset at the top), so the output reads like a timeline regardless of roster order.
- Every row carries its **date**, so a conversion that crosses midnight is unmistakable.
- When a row's date differs from the **source** (`←source`) row it is **highlighted** — the whole line is colored when stdout is a TTY, and a `+1`/`-1` day marker survives a pipe:

```
$ tzshift 23:30 Asia/Kolkata
tokyo-team      2026-06-20 03:00 +09:00  +1
india           2026-06-19 23:30 +05:30  ←source
legacy-billing  2026-06-19 18:00 +00:00
ny-dc           2026-06-19 14:00 -04:00
you             2026-06-19 11:00 -07:00  ←you
```

- The **`←source` row is the anchor for the `+1`/`-1` markers** (the zone you typed; for epoch / current-time inputs there is no source, so the anchor falls back to your local date). `show` always includes the source row and your local **`←you`** row, adding either automatically when your roster doesn't already list it — so the day markers always reference a row you can actually see. There is **no echo/header line**; the `←source` marker replaces it (a row that is both source and local shows `←you  ←source`).

Variations the narrative must cover:

- `tzshift show` (or bare `tzshift`) with no timestamp — outputs the current time across the configured roster zones (glanceable "what time is it for my team right now").
- `tzshift 14:30 Asia/Kolkata` (or `tzshift 14:30 india`) — wall-clock time + source zone, the core case (no subcommand → routed to `show`).
- `tzshift 14:30` — wall-clock time with **no** source zone: interpreted in your own local zone (the `←you` row is then also `←source`).
- `tzshift 2026-06-09 14:30 Asia/Kolkata` — explicit date (DST correctness depends on it; defaults to today).
- `tzshift 1749571200` — epoch seconds, the log-file case.
- `tzshift 14:30 Asia/Kolkata --to Europe/Berlin` — one-off target zone (IANA name or alias) not in the roster.
- No config file yet → falls back silently to **local time + UTC + America/Denver**; `tzshift config create` scaffolds a starter file and `tzshift config list` prints the current one.

### tzshift list — reference

Lists the IANA zones tzshift knows and, for each, the abbreviation it is *currently* displaying (e.g. `America/Denver` → `MDT` today, `MST` in winter). It also echoes any `[abbreviations]` the user has defined, so they can see exactly how their own source-zone shortcuts resolve.

- `tzshift list` — enumerates IANA zone identifiers and each zone's computed current abbreviation, plus the user's configured `[abbreviations]` (if any). There is **no shipped abbreviation → zone mapping** to list, by design.

**OS-agnostic caveat (a real constraint, not a detail):** zone *identifiers* (`America/Denver`, `Asia/Kolkata`) come from the IANA database, which tzshift carries embedded so it does not depend on the host OS's zoneinfo layout. Abbreviations shown on output (`PDT`, `JST`) are **computed per-instant from that data for display only** — they are never used to *resolve* a source zone, because there is no canonical, unambiguous "abbreviation → zone" direction. So `list` shows each zone's computed current abbreviation; resolving a source zone is done from an IANA name or the user's own `[abbreviations]`, never from a shipped guess.

### Later surfaces (rough — shaped after the CLI ships and is measured)

- **macOS widget (M2):** glanceable team-time awareness on the desktop. Built natively (Swift) when validated, reusing tzshift's proven translation logic — not the Go CLI binary itself.
- **iOS / Android (M3/M4):** hypotheses only; mobile demand unvalidated.

## Hard design decisions

- **The tool is `tzshift`, not `tz`.** `tz` is already taken (research: the oz/tz TUI world clock). `tzshift` is subcommand-oriented (`show`, `list`), with `show` as the default when the first argument isn't a known subcommand.
- **Go as the implementation language.** Phase 1 requires a CLI that runs across Linux, macOS, and other unix environments and is trivial to distribute — Go's single static binary and cross-compilation are the best fit. The language is chosen for phase 1's *validated* requirement, not for the hypothesis milestones (M2/M3); when the macOS widget is validated it will be built natively in Swift, reusing the translation logic. (Earlier Swift choice reversed — see Changelog and Alternatives.)
- **Zone data from the IANA database, carried embedded.** Zone identifiers are enumerated from the IANA tz database bundled in the binary, so behaviour does not vary with the host OS's zoneinfo. Abbreviations are computed per-instant from that data **for display only**.
- **No shipped abbreviation→zone resolver — the source zone is an IANA name or the user's own alias.** Abbreviations are too ambiguous to resolve safely (IST = India/Israel/Ireland; CST = China/US Central), and shipping a default would *be* the silent guess the tool exists to eliminate. The source zone must be an IANA identifier or a `[zones]` roster alias. An optional user-defined `[abbreviations]` table lets a user declare their own shortcuts (their data, not shipped); an unknown or unresolvable source token is a hard error, never a guess.
- **DST correctness over cleverness.** All conversions go through the IANA tz database; never hand-roll UTC offsets. The date matters, so date-less input assumes "today" and says so, and every output row carries its date. *(The Spec names the concrete library; the Shape states only the requirement.)*
- **Output is for humans first, pipes second.** Rows are column-aligned with arrow prefixes, sorted **east-most (largest UTC offset) first** so output reads like a timeline independent of roster order; the local **`you` row is marked** rather than pinned. Rows whose date differs from the local row are **highlighted** — color when stdout is a TTY, plus a `+1`/`-1` marker that survives a pipe (off-TTY the marker, not color, carries the date-shift signal).

## Identified unknowns

1. **Input format breadth** — which timestamp formats appear most in real channels/logs? Start with the variations above; let usage (Measure) drive the next set.
2. **Abbreviation ergonomics** — resolved for M1 by *not* shipping a mapping (source = IANA name or user alias; optional user `[abbreviations]`). Open for the Measure: do users miss pasting a raw `IST` from chat enough to want a curated common-set opt-in later?
3. **Distribution across unix variants** — Go gives static binaries for Linux/macOS easily; confirm "any unix" (e.g. BSD) targets and the release/cross-compile matrix actually needed.
4. **Roster shape** — is a flat `[zones]` table enough, or do users immediately want groups/ordering?
5. **Awareness vs. translation** — does glanceable team-time awareness (the `show`-with-no-timestamp / widget territory) or arbitrary translation drive more daily use? Determines how M2 gets shaped.

## Alternatives considered (rejected)

- **Swift** — was the original choice for its native path to the macOS widget (M2) and iOS (M3). Rejected for phase 1: those are unvalidated hypothesis milestones, and letting them drive the phase-1 language conflicts with the now-committed cross-unix requirement, where Swift is weakest (Linux works but "any unix" isn't first-class; distribution is heavier). Swift returns as the natural choice for M2 *when M2 is validated*.
- **Rust** — also excellent cross-platform static binaries, and a single core could FFI into iOS/macOS later. Rejected for phase 1 in favour of Go's faster CLI development and simpler toolchain; revisit only if a shared-core-across-all-surfaces need is actually validated.
- **Interactive TUI dashboard** — richer, but slower to build, not pipeable, and doesn't serve the "translate this one timestamp" moment the narrative centers on. Glanceable awareness is `show`'s / the widget's territory.
- **Natural-language input ("tomorrow 2pm Bangalore time")** — attractive, but parsing ambiguity is a tarpit. Structured input until usage data argues otherwise.
