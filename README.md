# tzshift

> A fast, error-proof CLI for translating timestamps across timezones — built for SREs on distributed teams.

## The problem

SREs on distributed teams constantly receive timestamps posted in *someone else's* timezone — a teammate's "deploy at 14:30," a legacy system logging in the datacenter's local zone, an incident timeline mixing three regions. Translating them means mental math, a web search, or a world-clock app: slow, outside the terminal, and error-prone — especially across DST boundaries.

`tzshift` takes any posted timestamp and instantly prints it in your local time **and** the zones your team and systems actually use — one command, in the terminal, DST-correct, no mental math.

```
$ tzshift 23:30 Asia/Kolkata
tokyo-team      2026-06-21 03:00 +09:00  +1
india           2026-06-20 23:30 +05:30  ←source
legacy-billing  2026-06-20 18:00 +00:00
ny-dc           2026-06-20 14:00 -04:00
you             2026-06-20 11:00 -07:00  ←you
```

Rows are sorted east-most first, every row carries its date, and a row on a
different day than the source is marked (`+1`/`-1`) and the whole line is
highlighted on a TTY. The source zone you typed (`←source`) and your own local
time (`←you`) are always shown — added automatically if your roster doesn't
already list them — so the markers always have a visible anchor.

## Install & use

```
go install github.com/torypatnoe/tztools/cmd/tzshift@latest   # or: make build
```

```
tzshift 14:30 Asia/Kolkata          # wall-clock time in a zone (date = today)
tzshift 14:30                       # no zone -> interpreted in your local zone
tzshift 14:30 india                 # using a [zones] alias
tzshift 2026-06-09 14:30 UTC        # explicit date
tzshift 1749571200                  # epoch seconds
tzshift                             # current time across your roster
tzshift list                        # known zones (east-most first) + your abbreviations
tzshift config create               # scaffold ~/.config/tzshift/config.toml
tzshift config list                 # print your current config
```

The source zone is an **IANA name** (`Asia/Kolkata`) or a **roster alias** —
tzshift ships no three-letter abbreviation guesses, since `IST` is ambiguous
(India? Israel? Ireland?). Define your own under `[abbreviations]` if you want them.

Configure your roster once in `~/.config/tzshift/config.toml`:

```toml
[zones]
tokyo-team     = "Asia/Tokyo"
india          = "Asia/Kolkata"
legacy-billing = "UTC"
ny-dc          = "America/New_York"
you            = "America/Los_Angeles"

# optional: your own source-zone shortcuts
[abbreviations]
IST = "Asia/Kolkata"
```

With no config file, tzshift falls back to local + UTC + America/Denver. Run `tzshift config create` to scaffold one.

## Status

**M1 CLI built.** This project is run end-to-end on the [Forma](#how-this-project-is-managed) workflow: Idea → Research → Customer Narrative → Project → Shape → Spec → Tickets → **Build** are done for **Cycle 1 (M1 — the CLI translator)**. Next: dogfood Measure.

- **Language:** Go — a single static binary across Linux/macOS/unix (IANA tz database embedded via `time/tzdata`; build with `make release`).
- **Shape:** subcommand-oriented (`show` default, `list`), TOML roster, no shipped abbreviation mapping (IANA/alias source zones), DST-correct, date-aware east-most-first output.
- See [project/cycles/cycle-1.md](project/cycles/cycle-1.md) for the bet and tickets.

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

## License

Copyright 2026 Tory Patnoe.

Licensed under the [Apache License, Version 2.0](LICENSE). This covers everything in the repository, including the project documents under [project/](project/) and [research/](research/).
