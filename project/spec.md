# Spec: tzshift

> Forma living Spec — single document, always the full current state. A `## Changelog` at the top records what changed each cycle; the per-cycle delta is rendered from the changelog + git diff, not maintained separately. The milestone label is taken from the cycle record's Bet, not declared here independently.

## Changelog
- 2026-06-19 (later): **No shipped abbreviation→zone mapping** — source zone is an IANA name or a `[zones]` alias, with an optional user-defined `[abbreviations]` escape hatch (AC8). Ambiguity AC reframed: an unknown/unresolvable source token is the hard error (AC9). Output now shows the **date** on every row, **sorts east-most (largest UTC offset) first** (overriding file order in AC6), **highlights cross-date rows** (color on TTY + persistent `+1`/`-1` marker), and **marks the local `you` row** (AC1, AC10). `list` no longer lists a shipped mapping — IANA zones + computed current abbreviation + the user's own `[abbreviations]` (AC14). Added AC15 (ordering) and AC16 (date display + cross-date highlight). Example commands updated to IANA/alias source zones.
- 2026-06-19: Tool renamed `tz` → **`tzshift`**; subcommands `show` (default) / `list` added (AC1, AC13–AC14). **Language Swift → Go**; distribution AC rewritten for cross-unix static binaries (AC12). No-config fallback corrected to local + UTC + America/Denver (AC7). Abbreviation ACs clarified: mapping is tzshift's shipped data, not OS-derived (AC8–AC9). Config dir → `~/.config/tzshift`. Org-standards check updated to Go/IANA.
- 2026-06-10: Initial Spec for Cycle 1 (M1 — CLI translator). Written against [cycle-1.md](cycles/cycle-1.md)'s Bet and the [living Shape](shape.md).

**Project:** [project.md](project.md)
**Living Shape:** [shape.md](shape.md)
**Cycle record:** [cycles/cycle-1.md](cycles/cycle-1.md)
**Milestone:** M1 — CLI translator *(from the Cycle 1 Bet)*
**DRI:** Tory Patnoe
**Status:** Draft — awaiting DRI review
**Last updated:** 2026-06-19

---

## Scope of this cycle

The detailed definition of the M1 CLI translator, `tzshift`. Bounds are inherited from the Cycle 1 Bet — anything in that Bet's *NOT in this cycle* list is out of scope here and is not re-litigated. This Spec turns the Shape's solution narrative into observable acceptance criteria.

`tzshift` is subcommand-oriented: `show` (the workhorse, and the default when the first argument isn't a known subcommand) takes a posted timestamp and prints it in the user's local time and each roster zone; `list` is reference output. The tool runs across Linux, macOS, and other unix environments.

## Acceptance criteria

Written as observable outcomes. "Done" for Cycle 1 means all of the following hold.

### Core translation (`show`)
- **AC1 — Wall-clock + source zone.** `tzshift 14:30 Asia/Kolkata` (or an alias, `tzshift 14:30 india`; routed to `show`) prints the input echoed with its resolved IANA zone, then one output row per roster entry, each showing the **date**, the local wall-clock time, the zone abbreviation, and the roster label. Rows are ordered and date-marked per AC15–AC16. Example output matches the Shape's narrative block.
- **AC2 — Explicit date.** `tzshift 2026-06-09 14:30 Asia/Kolkata` uses the given date for the conversion. With no date, the conversion assumes **today** and the output states that the date was assumed.
- **AC3 — Epoch seconds.** `tzshift 1749571200` interprets a bare integer as Unix epoch seconds and translates it to every roster zone.
- **AC4 — One-off target.** `tzshift 14:30 Asia/Kolkata --to Europe/Berlin` adds a single ad-hoc output row for a zone not in the roster, without modifying the roster. `--to` accepts an IANA name or a `[zones]` alias (never a shipped abbreviation).
- **AC5 — DST correctness.** A timestamp that falls on either side of a DST transition (verified with at least one spring-forward and one fall-back date per affected zone) converts to the correct local time. All conversions go through the IANA tz database; no hand-rolled offset math.
- **AC13 — Current-time mode.** `tzshift show` (or bare `tzshift`) with no timestamp prints the current time across the roster zones.

### Reference (`list`)
- **AC14 — Zone & abbreviation listing.** `tzshift list` enumerates the IANA zone identifiers tzshift knows and each zone's computed current abbreviation, plus the user's configured `[abbreviations]` (if any). There is **no shipped abbreviation → zone mapping** to list. Zone data is read from the IANA database **embedded in the binary** (not the host OS's zoneinfo), so output is identical across Linux/macOS/unix.

### Roster & configuration
- **AC6 — Roster file.** Output rows are driven by `[zones]` in `~/.config/tzshift/config.toml` (alias → IANA zone). File order does **not** determine output order — rows are sorted per AC15. The aliases defined here are also valid source-zone and `--to` tokens.
- **AC7 — No-config fallback.** With no config file present, the tool prints a helpful message showing how to create one, then falls back to translating into **local time + UTC + America/Denver**. It does not error.
- **AC8 — Source-zone resolution (no shipped mapping).** A source zone resolves from, in order: a `[zones]` alias, the user's `[abbreviations]` table (`abbr → IANA zone`, the user's own data), or a literal IANA identifier. tzshift ships **no** abbreviation→zone table — a bare three-letter abbreviation that the user has not defined does not resolve.
- **AC9 — Unknown source zone is a hard error, never a silent guess.** When a source token resolves to no zone (e.g. a bare `IST` with no matching alias or `[abbreviations]` entry), the tool exits non-zero with a message telling the user to use an IANA name or define an alias/abbreviation. It never picks a zone silently.

### Output & ergonomics
- **AC10 — Human-first, pipe-safe output.** Default output is column-aligned with arrow prefixes and is plain text — greppable and pipeable. Color is emitted only when stdout is a TTY (and suppressed by `--no-color`); all information conveyed by color is also conveyed by text (see AC16), so piped output loses no meaning.
- **AC11 — Exit codes.** Successful translation exits 0; bad input, an unresolvable source zone, and malformed config exit non-zero with a message on stderr.
- **AC15 — Chronological ordering (east-most first).** Output rows are sorted by UTC offset descending — the **east-most zone (largest offset) first**, down to the west-most. This overrides `[zones]` file order. Ordering is deterministic for a given instant; the local `you` row is marked (AC16) rather than positioned.
- **AC16 — Date display & cross-date highlight.** Every output row shows its date alongside the wall-clock time. A row whose local date differs from the local (`you`) row's date is **highlighted in color when stdout is a TTY** and is annotated with a persistent `+1` / `-1` (or larger) day-offset marker that survives a pipe / `--no-color`. The local row is marked (e.g. `←you`) so it is identifiable despite not being pinned to the top.

### Distribution
- **AC12 — Builds and runs across unix.** Builds from source via `go build` and runs on Linux and macOS (and other unix targets in the release matrix). A statically linked release binary is produced per target via cross-compilation, with no runtime dependency on the host's zoneinfo (IANA DB embedded). (Homebrew and other package managers remain out — see Bet.)

## Success metric

**What user behavior should change:** the SRE stops leaving the terminal (or doing mental math) to translate a posted timestamp.

**How it's measured this cycle:** the customer (DRI, dogfooding as the first SRE) records, over a week of real use, how often a posted timestamp is translated with `tzshift` versus the old path (browser/world-clock/mental math). Target: `tzshift` becomes the default reach-for tool for posted-timestamp translation, and at least one DST-boundary conversion is handled correctly that would previously have been error-prone. Because Cycle 1's appetite is small and the user base is one, this is a qualitative dogfood signal, not a quantitative funnel — that is the honest measurement available at M1 and is recorded as such in the Measure.

## Feature flag plan

**No runtime feature flag.** tzshift is a locally installed single-binary CLI with no server and no staged rollout surface; deploy and release collapse into one event — tagging and publishing a release binary. The Forma deploy/release separation does not apply at M1. New input formats or surfaces in later cycles ship as new released versions, not flagged code paths. *(Recorded explicitly so the gate is satisfied by a reasoned decision, not an omission.)*

## Scale considerations

- **Data sizes:** trivial. Config is a small hand-edited TOML file (tens of entries at most); each invocation translates one timestamp into a bounded roster. No large inputs.
- **Persistence:** none beyond the user-owned config file. The tool holds no state between invocations and writes nothing except optionally a scaffolded config on first run (if implemented).
- **Rate limits / throughput:** not applicable — no network calls, no shared service. Startup time is the only performance concern; a Go static binary cold-starts in single-digit milliseconds, comfortably under a sub-100ms perceived-latency target for a single conversion.
- **Concurrency:** single-shot process per invocation; no shared mutable state.

## Security sign-off

- **Auth:** none required. The tool runs locally as the invoking user, reads only the user's own config under `~/.config/tzshift/`, and makes no network requests. There is no server, no credential, and no multi-user surface.
- **Data handling:** timestamps and zone labels are processed in-memory and printed; nothing is transmitted or logged off-device.
- **Trust boundary:** the config file is user-authored; the tool must fail safe on malformed TOML (AC11) rather than crashing or executing config content. TOML is parsed as data only — no code evaluation.
- **Decision:** no auth, no network, local-only. Signed off for M1 by the DRI.

## Organization standards check

No standing Forma org-standards artifacts (security requirements, architecture standards, compliance rules) exist for this project yet — tzshift is a solo project, so the Shape's hard design decisions serve as the de facto standards. Checked against them:

- **Language / platform standard:** Go, cross-unix single static binary, per Shape ("Go as the implementation language" — chosen for phase 1's validated cross-unix requirement). Spec complies.
- **Time-correctness standard:** all conversions via the IANA tz database (embedded in the binary), no hand-rolled offsets, per Shape's "DST correctness over cleverness." Spec complies (AC5, AC12, AC14).
- **Deviations flagged for justification:** none.

When tzshift grows beyond a solo project, these should be promoted to standing org-standard documents with their own DRIs.
