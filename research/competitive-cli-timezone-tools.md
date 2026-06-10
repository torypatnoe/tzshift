# Research: Existing timezone tools and why they fall short

> Forma Research-stage artifact (stage 2). Referenced by the [Customer Narrative](../project/customer-narrative.md).

**DRI:** Tory Patnoe
**Recorded:** 2026-06-10
**Method:** Working-knowledge survey of the tools SREs commonly reach for, scoped to the effort appropriate for a small internal tool (per Forma, research is effort-scalable). Findings are from direct familiarity with these tools, not formal market analysis — verify specifics before citing externally.

---

## Landscape

| Tool | What it does | Why it falls short for this customer |
|---|---|---|
| `date` / GNU coreutils (`TZ=... date -d`) | Shell-native conversion | Arcane syntax, differs across macOS/Linux, no team concept, easy to get wrong under pressure |
| `zdump` | Dumps zone offsets | Reference tool, not a translator; no parsing of arbitrary timestamps |
| `tz` (oz/tz, TUI) | Terminal world clock with zone list | Awareness only — doesn't translate a given arbitrary timestamp; no team roster semantics |
| worldtimebuddy / Every Time Zone (web) | Visual zone comparison | Requires leaving the terminal; manual setup per conversion; not scriptable |
| Clocker / menubar world clocks (macOS) | Glanceable multi-zone clocks | Current-time awareness only; can't translate "14:30 IST from this log line" |
| Slack timestamp formatting | Renders times in reader's local zone | Only works if the sender formats it correctly — the problem is precisely that they don't |

## Finding: the gap

No existing tool combines all three of:

1. **Translation of arbitrary posted timestamps** — wall-clock times, epoch seconds, log-line formats — rather than just showing current time in other zones.
2. **A team/zone roster** — output expressed in the people, datacenters, and systems the SRE actually cares about, named in their own words.
3. **Presence inside the SRE's working surfaces** — terminal first, then glanceable widgets and mobile.

Each existing tool covers at most one of these. This gap is what the Customer Narrative's problem statement rests on.
