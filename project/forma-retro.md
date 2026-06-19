# Forma dogfood retro — lessons from running tztool through the workflow

> Destination: this is feedback for the **`../../forma-program-management`** repo, not a tztool artifact. It extends the dogfood validation already logged there ([cycles/cycle-1.md](../../forma-program-management/cycles/cycle-1.md) → "Review 1 — AI dogfood pass (tztool)"). Treat each entry below as a candidate validation-log entry / Shape changelog item for that repo. Captured 2026-06-18/19 while taking tztool from Idea → Spec.

---

## Context

tztool is the first real project run end-to-end on Forma. The methodology itself is being validated by being used. This file records what the workflow caught, what it let slip, and concrete suggestions for the Forma Shape — so Forma improves from real use rather than from theory.

---

## What worked well

- **Living-doc + per-cycle-record split was easy to follow.** Knowing that `shape.md`/`spec.md` are always-current and `cycles/cycle-N.md` is the only per-cycle file removed all "where does this go?" friction. The "per-cycle facts live only in the Bet" rule prevented duplication in practice, not just in principle.
- **The Research stage paid for itself immediately.** Forcing a recorded competitive survey before the narrative meant the `tz`/`worldtimebuddy`/`Clocker` gap was documented and the narrative's claims traced to it. (See the linkage failure below for the flip side.)
- **"Sharpen, not change" held under pressure.** When the customer was narrowed to "SREs as beachhead," the rule gave a clean test for whether that was a sharpen (yes) or a new Project (no). The distinction was genuinely useful, not bureaucratic.
- **Outcome-vs-mechanism in the Customer Narrative.** Separating the durable customer gain from how the tool delivers it made it obvious which parts a future cycle is allowed to refine.
- **AI-as-gate caught real Spec gaps.** The feature-flag question forced an explicit "no flag, deploy=release" decision instead of a silent omission; the org-standards check surfaced that none exist yet for a solo project.

## What didn't work well — and suggested Forma changes

### 1. No forcing function that the *previous* artifact was actually reviewed before the next gate opened
The Spec was written and reviewed before the Shape was reviewed. Issues that a Shape review would have caught (below) slipped straight through into the Spec. Forma says "each stage is a gate, don't skip forward," but nothing made the Shape review a precondition for opening the Spec.
- **Suggestion:** a gate is not passed until the previous artifact has a recorded review/sign-off (even a one-line AI-gate pass). Opening the Spec should require a reviewed Shape, not just an existing one.

### 2. Research findings did not propagate downstream — the gate didn't cross-check
The research artifact explicitly recorded that a terminal tool named **`tz` already exists**. The Shape then named the product `tz` anyway. The naming collision survived all the way to the Spec before a human caught it. Research informed the narrative but was never re-consulted at Shape/Spec time.
- **Suggestion:** downstream gates (Shape, Spec) should cross-check against recorded Research findings, not just the Customer Narrative. "This Shape names the tool `tz`; research finding X says `tz` is taken — flag." This is exactly the kind of check Forma's AI-gate vision should perform.

### 3. A hypothesis milestone drove a foundational phase-1 technical decision (over-speccing, in disguise)
The Shape chose Swift as the implementation language *because* of M2 (macOS widget) and M3 (iOS) — both explicitly **hypothesis** milestones. Meanwhile the known M1 requirement (a cross-platform Linux/macOS/unix CLI) is the thing Swift is *weakest* at. Forma loudly warns against over-speccing future cycles, but the principle was only framed around scope/planning, not technical/architecture choices. The anti-pattern wore a different hat and the workflow didn't recognize it.
- **Suggestion:** extend the "keep future milestones rough" principle explicitly to foundational technical decisions. A phase-1 tech choice should be justified by phase-1's *validated* requirements; optimizing it for unvalidated future milestones is the same mistake as over-speccing them. The Shape gate should ask: "is this technical decision driven by a hypothesis milestone?"

### 4. The Shape coupled itself to a specific implementation
The Shape's "hard design decisions" baked in `Swift's Foundation.TimeZone`. That tied the Shape (a solution-narrative artifact) to a library choice, which the Spec then inherited as if settled. When the language decision reopened, the Shape's design rationale had to be untangled from the implementation detail.
- **Suggestion:** Shape decisions should state the *requirement* ("all conversions via the IANA tz database; never hand-roll offsets") and leave the *library/runtime* to the Spec. Keep the Shape implementation-light so a language change doesn't rewrite the product thinking.

### 5. The org-standards check is awkward for solo / very small projects
The Spec gate's "organization standards check" assumes standing org-standard artifacts exist. For a solo project there are none, so the check degraded into "the Shape's decisions are the de facto standard." It still added value, but the gate's framing didn't anticipate the small-team case Forma explicitly targets.
- **Suggestion:** define a lightweight default for the small-team case — either an explicit "no org standards yet; using Shape decisions as standards" recorded outcome, or a starter org-standards template Forma seeds for new orgs.

## Meta-observations

- **Dogfooding is doing exactly what cycle-1 of Forma bet it would.** Every issue above is a workflow gap found by *using* the workflow, not inspecting it. That validates the dogfood approach itself — and argues for keeping a real downstream project in the loop as Forma evolves.
- **The cross-repo split (methodology vs. project) works, but the linkage is manual.** tztool's `CLAUDE.md` points at `../forma-program-management` by hand. If Forma ships, the relationship between the methodology and a project instance probably wants to be a first-class concept, not a relative path.

---

*Next: fold these into `../../forma-program-management/cycles/cycle-1.md`'s validation log and/or the Shape changelog when updating that repo.*
