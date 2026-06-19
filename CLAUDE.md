# tzshift — CLI Timezone Translator for SREs

## What this is

tzshift is a cross-unix (Linux/macOS) command-line tool that takes any posted timestamp (a teammate's wall-clock time, a legacy log's epoch seconds) and instantly translates it into the user's local time **and** their configured team/system zones — terminal-first, DST-correct, no mental math. Subcommand-oriented (`show` default, `list`). Target user: SREs on distributed teams. (The repo is `tztools`; the tool binary is `tzshift` — `tz` was already taken.)

This project is managed with the **Forma** program-management workflow. The methodology lives in the sibling repo `../forma-program-management` — read its `CLAUDE.md` and `project/shape.md` for the full model. The rules below are the parts that govern how files in *this* repo are created and edited.

## Repo Structure

This repo is structured to match Forma's workflow.

```
project/
  idea.md                — one-sentence problem + owner
  customer-narrative.md   — who/problem/after/non-goals (sharpened by cycles, never changed)
  project.md              — milestone map + project hill chart, updated after each Measure
  shape.md                — living Shape: full product as understood; changelog at top
  spec.md                 — living Spec: full current state; changelog at top
  cycles/
    cycle-N.md            — the ONLY per-cycle file: Bet + hill chart + Measure
research/
  competitive-cli-timezone-tools.md — Research-stage artifact, cited by the Customer Narrative
```

## File Conventions — Read Before Creating Any Files

**Artifacts are single living documents, not per-cycle copies.** Do NOT create per-cycle shape/spec files or `cycles/cycle-N/` subdirectories. That two-sources-of-truth pattern is explicitly rejected in Forma. `cycles/` holds exactly one flat `cycle-N.md` per cycle — nothing else.

- `project/shape.md` and `project/spec.md` — one file each, always the full current state, with a `## Changelog` at the top. Git history is the version record; the per-cycle delta is *rendered* from the changelog + diff, not maintained as a second file.
- `cycles/cycle-N.md` — opened at cycle start with **the Bet** (milestone, goal, appetite, in/out scope, DRI), carries the **hill chart** position during Build, closed at cycle end with **the Measure** (result vs. goal, learning note). Once closed, immutable.
- **Per-cycle facts (appetite, scope boundary, DRI confirmation) live ONLY in the Bet** — never duplicated into shape.md or spec.md. Two homes for one fact is how drift starts.
- New research → `research/`. Shape change → edit `shape.md` + changelog entry. Spec change → edit `spec.md` + changelog entry. Cycle start/close → `cycles/cycle-N.md`.

## The workflow (where tztool sits)

```
PROJECT LEVEL:  Idea → Research → Customer Narrative → Project
CYCLE LEVEL:    Shape → Spec → Tickets → Build → Ship → Measure   (repeats per cycle)
```

Each stage is a gate: don't skip forward, you may move back, every artifact has one DRI. Hard gate: no cycle without a completed Customer Narrative. Cycles *sharpen* the narrative (better wording, validated evidence) but cannot *change* the customer or problem — that would be a new Project.

## Current Status

- **Idea / Research / Customer Narrative / Project / Shape** ✅ complete (Customer Narrative is a draft awaiting validation).
- **Cycle 1 — M1 `tzshift` CLI translator:** Bet placed; Shape + **Spec drafted** (`project/spec.md`). Next: Tickets → Build.
- **Milestone map:** M1 CLI (current) → M2 macOS widget (rough) → M3/M4 mobile (hypothesis).

## Key product decisions (from the Shape)

- **DST correctness over cleverness** — all conversions via the IANA tz database (embedded in the binary); never hand-roll offsets. Date-less input assumes "today" and says so.
- **Ambiguous zone abbreviations never guess silently** — opinionated default mapping (tzshift's own shipped data, not OS-derived), `[abbreviations]` config override, hard error with candidates when truly ambiguous.
- **Go, cross-unix** — single static binary for Linux/macOS, chosen for phase 1's validated cross-unix requirement (not for hypothesis milestones). M2 macOS widget will be built natively in Swift when validated, reusing the translation logic.
- **`show` is the default subcommand** — bare `tzshift <timestamp>` routes to `show`; `tzshift list` is reference output.
- **Output is human-first, pipe-second** — aligned columns + arrows, plain text, greppable.

## Working with Claude here

- Think like a PM, not just an engineer — decisions trace back to the customer narrative and the core workflow.
- Capture design decisions with reasoning in `shape.md` (+ changelog); keep `research/` grounded in real data, not assertions.
- Flag when an assumption needs user validation rather than encoding it silently.
- Keep `../forma-program-management` as the source of truth for the *methodology*; this repo is the source of truth for *tztool*.
