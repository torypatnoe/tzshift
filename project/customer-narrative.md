# Customer Narrative: tztool

> Forma project-level artifact. This document belongs to the Project and does not change cycle to cycle. If the customer or the problem changes, that is a new Project.

**DRI:** Tory Patnoe
**Status:** Draft — awaiting DRI review
**Last updated:** 2026-06-09

---

## Who is the customer?

SRE team members — both individual contributors and team leads — working on distributed teams that span multiple timezones. They live in the terminal, coordinate across regions daily, and interact with systems and people that each report time in their own local zone.

## What is their problem?

Timestamps arrive in *someone else's* timezone, constantly:

- A teammate in another region posts "deploy at 14:30" — their 14:30.
- A legacy system writes logs in the datacenter's local zone, not UTC.
- An incident timeline mixes timestamps from systems in three different zones.

Translating these is done with mental math, a web search, or a world-clock app — all slow, all outside the SRE's working context (the terminal), and the mental-math path produces real errors: misread incident timelines, missed handoffs, wrong maintenance windows. DST transitions make even practiced conversions unreliable.

## What can the customer do after this exists that they cannot do today?

**Instant team-time awareness.** In one action, from wherever they are working, an SRE can:

- Take any timestamp in any timezone and see it in their own time and in their team's times — without leaving their workflow and without doing math in their head.
- Trust the answer completely, including across DST boundaries.

Today this requires context-switching to a website or app, manually setting up the conversion, and hoping they picked the right zone variant. After tztool exists, the translation is a single fast action with zero mental arithmetic.

## What we are explicitly NOT solving (non-goals)

- **On-call schedule management** — not replacing PagerDuty/Opsgenie; we will not manage rotations or schedules.
- **Calendar/meeting booking** — we will not book meetings or integrate with calendars as a scheduler.
- **General world-clock app** — not for travelers or the general public; SRE workflows only.
- **Team chat/notifications** — we will not send alerts or messages to teammates.

## Competitive research

| Tool | What it does | Why it falls short for this customer |
|---|---|---|
| `date` / GNU coreutils (`TZ=... date -d`) | Shell-native conversion | Arcane syntax, differs across macOS/Linux, no team concept, easy to get wrong under pressure |
| `zdump` | Dumps zone offsets | Reference tool, not a translator; no parsing of arbitrary timestamps |
| `tz` (oz/tz, TUI) | Terminal world clock with zone list | Awareness only — doesn't translate a given arbitrary timestamp; no team roster semantics |
| worldtimebuddy / Every Time Zone (web) | Visual zone comparison | Requires leaving the terminal; manual setup per conversion; not scriptable |
| Clocker / menubar world clocks (macOS) | Glanceable multi-zone clocks | Current-time awareness only; can't translate "14:30 IST from this log line" |
| Slack timestamp formatting | Renders times in reader's local zone | Only works if the sender formats it correctly — the problem is precisely that they don't |

**Gap:** no tool combines (a) translation of arbitrary posted timestamps, (b) a team/zone roster so output is in *the people and systems you care about*, and (c) presence inside the SRE's working surfaces — terminal first, then glanceable widgets and mobile.

## Hypothesis

If translation of any timestamp into team-relevant times is a single instant action inside the SRE's existing workflow, timezone-related coordination errors and context switches drop to near zero.
