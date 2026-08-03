---
name: wayfinder
description: Plan and orchestrate large features through a structured ideation interview, a shared decision map with To Do and Completed tickets, scoped decision logs, and subfeature maps. Use when a feature is too big for one session, when charting or working a FeatureName:Map, or when coordinating research, grill-me, and implementation across multiple sessions.
disable-model-invocation: true
---

# Wayfinder

Opinionated planning for **large features**: ideation → map → To Do tickets → human-closed Completed → consolidated PRD. Part of a **skill ecosystem** — see [REFERENCE.md](REFERENCE.md) for templates, decision-log prefixes, tracker ops, and sibling-skill hooks.

**Plan, don't implement** unless the map **Notes** say otherwise.

## When to use

| Situation | Action |
|-----------|--------|
| Rough feature idea, too big for one chat | **Chart** — ideation interview → create map |
| Existing `FeatureName:Map` (URL, issue, or path) | **Work** — resolve next To Do ticket |
| All To Do done, decision log stable | Hand off to `write-a-prd` → `prd-to-issues` |
| Child subsystem needs its own planning | Create **subfeature map**; link from parent **Subfeatures** |

Skip wayfinder when the path is already clear — use `grill-me` or implement directly.

## Modes

### Chart (create a map)

1. **Ideation interview** — Walk the **five coverage zones** from [grill-me](../grill-me/SKILL.md) (Surfaces & experience through Change, risk & evidence). Read each zone’s **Quick triage** only; do not depth-first grill.
   - For each zone: capture what is **known**, **unknown**, or **needs research**.
   - **Valid answers:** concrete decisions, `unknown`, or `needs research` (with the question that must be answered).
   - Do **not** resolve unknowns in this phase — collect them.
2. **Name the map** — `{FeatureName}:Map` (e.g. `CommandPalette:Map`). Derive **map slug** per [REFERENCE.md](REFERENCE.md#map-slug-and-decision-log-prefix).
3. **Set target outcome** — What this map is working toward (usually a buildable PRD). One or two lines.
4. **Create the map artifact** — GitHub issue (preferred) labelled `wayfinder:map`, or local file under `wayfinder/plans/` only when GitHub is unavailable. Use the [map template](REFERENCE.md#map-body-template).
5. **Create To Do tickets** — One child issue per resolvable unknown/research gap from ideation. Labels: `wayfinder:todo` + type (`wayfinder:research` | `:prototype` | `:grilling` | `:task`) + mode (`wayfinder:hitl` | `:afk`). Wire **native GitHub blocking** in a second pass.
6. **Seed fog** — Anything in scope but not yet sharp enough → map **Not yet specified**.
7. **Stop** — Charting does not resolve tickets. Fire **AFK research** subagents only if the user asks in this session.

### Work (resolve a ticket)

1. **Load the map** — Low-res body only (not every ticket).
2. **Pick a ticket** — User-named, or first **frontier** item: open, unblocked, in **To Do**, unclaimed. **Claim** via assignee before work.
3. **Resolve by type** — See [ticket types](REFERENCE.md#ticket-types). Agent produces findings/resolution; **never closes the issue**.
4. **Post resolution** — Comment on the ticket: summary, assets/links, new `GM-xx` rows if applicable. End with: *Ready for human review — close when accepted.*
5. **Human closes** — Only the human moves a ticket to **Completed** on the map (see [REFERENCE.md](REFERENCE.md#completed-workflow)).
6. **Update map** — On human close: move row To Do → **Completed** (one-line gist). Graduate fog → new To Do tickets if needed. Add **Subfeatures** links when spawning child maps.

**One To Do ticket per agent session** (except parallel AFK research when charting).

## Coverage zones (ideation only)

Same five as grill-me — triage breadth, not depth:

| Zone | grill-me reference |
|------|-------------------|
| Surfaces & experience | [surfaces-and-experience.md](../grill-me/references/surfaces-and-experience.md) |
| Behavior & correctness | [behavior-and-correctness.md](../grill-me/references/behavior-and-correctness.md) |
| Boundaries & integration | [boundaries-and-integration.md](../grill-me/references/boundaries-and-integration.md) |
| Persistence & data | [persistence-and-data.md](../grill-me/references/persistence-and-data.md) |
| Change, risk & evidence | [change-risk-and-evidence.md](../grill-me/references/change-risk-and-evidence.md) |

## Decision log

Each map owns a **scoped decision log** (`{MAP-SLUG}-GM-NNN`). Tickets and `grill-me` append rows; the map links to the log issue/file. Full rules: [REFERENCE.md](REFERENCE.md#decision-log).

## Subfeatures

Large greenfield work may spawn child maps (`SearchPanel:Map`) linked under parent **Subfeatures**. Parent map stays the integration index; cross-map consistency is explicit To Do tickets on the parent. See [REFERENCE.md](REFERENCE.md#subfeature-maps).

## Ecosystem (sibling skills)

| Skill | Role in wayfinder |
|-------|-------------------|
| `grill-me` | Resolves `wayfinder:grilling` tickets; appends `{MAP-SLUG}-GM-xx` to map decision log |
| `feature-ideation` | Optional **before** wayfinder for very raw seeds |
| `design-an-interface` | Optional for `wayfinder:prototype` tickets |
| `write-a-prd` | After map To Do empty — merge decision log into PRD `GM-xx` |
| `prd-to-issues` | Implementation slices from PRD |

Other skills (`research`, cloud AFK pickup) — see [REFERENCE.md](REFERENCE.md#ecosystem-integration).

## Refer by name

In narration and map sections, use **ticket titles**, not bare `#42`. IDs live inside linked names.
