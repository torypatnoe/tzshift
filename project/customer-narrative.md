# Customer Narrative: tztool

> Forma project-level artifact. Cycles may sharpen this narrative (recorded in the changelog); they cannot change the customer or the problem — that would be a new Project.

## Changelog
- 2026-06-10: Competitive research moved to its own Research-stage artifact ([research/competitive-cli-timezone-tools.md](../research/competitive-cli-timezone-tools.md)); narrative now references it
- 2026-06-09: Initial draft

**DRI:** Tory Patnoe
**Status:** Draft — awaiting DRI review
**Last updated:** 2026-06-10

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

## Research findings referenced

Competitive landscape recorded in [research/competitive-cli-timezone-tools.md](../research/competitive-cli-timezone-tools.md). Its finding: no existing tool combines (a) translation of arbitrary posted timestamps, (b) a team/zone roster so output is in *the people and systems you care about*, and (c) presence inside the SRE's working surfaces — terminal first, then glanceable widgets and mobile. Every claim this narrative makes about existing solutions traces to that document.

## Hypothesis

If translation of any timestamp into team-relevant times is a single instant action inside the SRE's existing workflow, timezone-related coordination errors and context switches drop to near zero.
