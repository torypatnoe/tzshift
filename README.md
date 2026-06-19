# tzshift

> A fast, error-proof CLI for translating timestamps across timezones — built for SREs on distributed teams.

## The problem

SREs on distributed teams constantly receive timestamps posted in *someone else's* timezone — a teammate's "deploy at 14:30," a legacy system logging in the datacenter's local zone, an incident timeline mixing three regions. Translating them means mental math, a web search, or a world-clock app: slow, outside the terminal, and error-prone — especially across DST boundaries.

`tzshift` takes any posted timestamp and instantly prints it in your local time **and** the zones your team and systems actually use — one command, in the terminal, DST-correct, no mental math.

```
$ tzshift 14:30 IST
14:30 IST (Asia/Kolkata)
→ 02:00 PDT  you
→ 05:00 EDT  ny-dc
→ 09:00 UTC  legacy-billing
→ 18:00 JST  tokyo-team
```

*(Planned interface — see [project/shape.md](project/shape.md). The CLI is specced but not yet built.)*

## Status

**Planning complete through the Spec gate; build not yet started.** This project is run end-to-end on the [Forma](#how-this-project-is-managed) workflow: Idea → Research → Customer Narrative → Project → Shape → Spec are done for **Cycle 1 (M1 — the CLI translator)**. Next gate is Tickets → Build.

- **Language:** Go — a single static binary across Linux/macOS/unix.
- **Shape:** subcommand-oriented (`show` default, `list`), TOML roster, opinionated abbreviation mapping, IANA tz database for DST correctness.
- See [project/cycles/cycle-1.md](project/cycles/cycle-1.md) for the current bet.

## Repo structure

| Path | Purpose |
|------|---------|
| [project/](project/) | Forma project artifacts: idea, customer narrative, project map, living shape, living spec |
| [project/cycles/](project/cycles/) | One record per cycle — the Bet (at start) + the Measure (at close) |
| [research/](research/) | Competitive analysis and other inputs to the Customer Narrative |
| [CLAUDE.md](CLAUDE.md) | Conventions for working in this repo |

## How this project is managed

tzshift is managed with **Forma**, an opinionated program-management workflow. The methodology itself lives in a sibling repository, [`forma-program-management`](https://github.com/torypatnoe/forma-program-management), and this repo dogfoods it.

**Repo layout dependency:** some documents here link to the methodology repo by relative path (`../forma-program-management`). To follow those links, clone both repositories as siblings:

```
parent/
  ├── tzshift/                  (this repo)
  └── forma-program-management/ (the methodology — keep this exact name)
```

[project/forma-retro.md](project/forma-retro.md) records lessons from dogfooding Forma on this project, destined for the methodology repo.
