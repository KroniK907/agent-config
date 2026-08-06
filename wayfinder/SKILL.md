---
name: wayfinder
description: Bootstrap and maintain FeatureName:Map GitHub trackers - create map skeleton, materialize tickets from map-discovery artifact, reconcile session output after human approval, and suggest next frontier skill. Use when starting a large feature, syncing chat to map/issues, or navigating wayfinder To Do tickets across sessions.
disable-model-invocation: true
---

# Wayfinder

**Tracker and router** for large features: map skeleton → [feature-discovery](feature-discovery/SKILL.md) → tickets → sibling skills → **Reconcile** → [define-bundle](define-bundle/SKILL.md) → [create-tasks](create-tasks/SKILL.md). Part of a **skill ecosystem** - see [REFERENCE.md](REFERENCE.md) for templates, materialize rules, approval protocol, routing table, and GitHub ops.

**Plan, don't implement** unless the map **Notes** say otherwise. Wayfinder does **not** run discovery interviews, strategic ideation, or grilling - it creates/updates GitHub state and suggests what skill to use next.

## When to use

| Situation | Mode |
|-----------|------|
| Rough feature idea, too big for one chat | **Chart** - skeleton → hand off to feature-discovery |
| Discovery capture ready | **Materialize** - create To Do tickets from map-discovery artifact |
| After sibling skill session; user approved outcome | **Reconcile** - comment, decision log, close ticket, update map |
| Existing map; need next step | **Route** - frontier + skill suggestion |
| Child subsystem needs its own planning | **Chart** subfeature map; link from parent **Subfeatures** |

Skip wayfinder when the path is clear - use `grill-me` or implement directly.

## Modes

### Chart (create skeleton)

1. **Accept seed** - User’s one-paragraph feature description (and target repo if not obvious).
2. **Name the map** - `{FeatureName}:Map` (e.g. `CommandPalette:Map`). Derive **map slug** per [REFERENCE.md](REFERENCE.md#map-slug-and-decision-log-prefix).
3. **Set target outcome** - What this map works toward (usually a buildable PRD). One or two lines.
4. **Create artifacts** - GitHub issues (preferred): decision log (`wf:decision-log`) → map (`wf:map`) linking the log. Local fallback: [plans/](plans/README.md). **To Do** empty; **Phase:** `charting`.
5. **Hand off to feature-discovery** - Tell user to continue with [feature-discovery](feature-discovery/SKILL.md) in this or a new chat. Pass: **map issue** link, target outcome, seed. **Do not** interview zones in wayfinder.
6. **Stop** - Chart does not create tickets. Next step after discovery: **Materialize**.

### Materialize (capture → tickets)

1. **Load map** - Confirm map issue and slug.
2. **Ingest map-discovery artifact** - From (in order): current chat block with `## Map discovery`; latest map-issue comment with that heading and **Status:** `ready for materialize`; user paste. See [materialize rules](REFERENCE.md#materialize-from-map-discovery).
3. **Create To Do tickets** - One child issue per **Ticket candidates** row. Labels: `wf:todo` + type + mode. Titles per [ticket title conventions](REFERENCE.md#ticket-title-conventions). Wire **blocked-by** in a second pass per materialize rules.
4. **Update map** - Populate **To Do** table; copy **Fog** → **Not yet specified**; confirm **Out of scope suggestions** with user if present; set **Phase:** `deciding`. Append **Completed** gist: *Map discovery materialized - N tickets*. Reply on map-discovery comment: **Status:** `materialized`.
5. **Route** - Suggest first frontier ticket and skill (see **Route** below).

### Reconcile (sync session → GitHub)

Run when user explicitly invokes wayfinder after a sibling skill session.

1. **Load map + ticket + session output** - From user message, frontier context, or sibling skill thread (grilling Q&A, research findings, prototype outcome, etc.).
2. **Infer holistic tracker delta** - Per [Reconcile inference](REFERENCE.md#reconcile-inference), derive:
 - Decision-log rows (`{MAP-SLUG}-GM-NNN`) with `[global]` vs bundle-scoped tags
 - **Decision coverage** row additions/updates
 - Map diff: **Completed** gist, **Not yet specified**, **Out of scope**, **Notes**
 - **New ticket candidates** (research / prototype / grilling / task) for unresolved or follow-on work
 - **Bundle cluster suggestions** for define-bundle (name, GM IDs, rationale - draft bundle issues are **not** created here)
 - **Ticket invalidations** - close, retitle, or move superseded To Do items
 - **Route hint** - recommended next skill(s)
3. **Post resolution** - Comment on ticket using the [resolution template](REFERENCE.md#reconcile-resolution-template). Add label **`wf:needs-review`** to the ticket. End with: *Ready for review - reply **Approved - reconcile and close** (or **Approved - reconcile, keep open**) when accepted. Edit any section in this comment before approving.*
4. **On approval** - When user says an [approval phrase](REFERENCE.md#approval-phrases), agent executes approved sections: close issue (if full approval), move row To Do → **Completed**, append decision log **body**, update **Decision coverage**, update fog/Notes/Out of scope, **materialize approved ticket candidates** (create child issues + **To Do** rows), apply ticket invalidations. Remove label **`wf:needs-review`**. Requires `gh` auth on target repo.
5. **Partial approval** - **Approved - reconcile, keep open** → apply comments/log/map notes and optional ticket creates without closing the source ticket.

**Default:** Do not reconcile mid-session without user invoke. Sibling skills may remind: *Invoke wayfinder Reconcile when ready.*

**Grilling sessions:** Reconcile is the bridge from depth-first Q&A to map evolution - not only binding GM rows, but inferred frontier tickets, bundle-ready clusters, and route hints. **`bundle approved`** remains [define-bundle](define-bundle/SKILL.md); Reconcile **suggests** clusters only.

### Route (suggest next step)

1. **Load map** - Low-res body (not every ticket thread) + decision log link.
2. **Compute planning frontier** - First open, unblocked, unclaimed **To Do** item per [REFERENCE.md](REFERENCE.md#frontier-queries).
3. **Check implementation path** - If **Decision coverage** has a cluster of **`open`** rows ready to build (see [define-bundle route heuristics](define-bundle/REFERENCE.md#route-heuristics-for-wayfinder)), suggest [define-bundle](define-bundle/SKILL.md) alongside or instead of planning frontier when user wants to ship incrementally.
4. **Suggest** - One recommended next step + skill from [routing table](REFERENCE.md#routing-table). Optional second choice if ambiguous. User picks skill and starts work - wayfinder does not resolve tickets.

Approved bundles → suggest [create-tasks](create-tasks/SKILL.md). **`wf:approved`** **Implementing** tasks → suggest [implement-task](implement-task/SKILL.md). **`write-a-prd`** / **`prd-to-issues`** only for small map-free scope - not a map Route handoff.

## Decision log

Each map owns a **scoped decision log** (`{MAP-SLUG}-GM-NNN`). Sibling skills and **Reconcile** append rows; map links to log issue. Full rules: [REFERENCE.md](REFERENCE.md#decision-log).

## Subfeatures

Large greenfield work may spawn child maps (`SearchPanel:Map`) linked under parent **Subfeatures**. Parent stays integration index. See [REFERENCE.md](REFERENCE.md#subfeature-maps).

## Ecosystem (sibling skills)

| Skill | Role |
|-------|------|
| [feature-discovery](feature-discovery/SKILL.md) | Chart handoff - posts map-discovery comment on map issue |
| [constrain-fog](constrain-fog/SKILL.md) | Groom **Not yet specified** fog - **`Constrain:`** ticket + fog-resolution artifact |
| [strategic-ideation](strategic-ideation/SKILL.md) | Scope/strategy expand → tension → prune (ticket or pre-PRD) |
| [grill-me](grill-me/SKILL.md) | `wf:grilling` tickets → decision log |
| [design-modules](design-modules/SKILL.md) | Modules shaping (one or more) - bundle step before create-tasks; planning `wf:prototype` interface exploration |
| [define-bundle](define-bundle/SKILL.md) | GM cluster → draft/approved `wf:bundle` issue |
| [create-tasks](create-tasks/SKILL.md) | Approved bundle → **Implementing** tasks |
| [one-off](one-off/SKILL.md) | HITL map **To Do** implementation without bundle pipeline |
| [implement-task](implement-task/SKILL.md) | **`wf:approved`** tasks → Method → **code-review** → push → **`awaiting-reconcile`** |
| [code-review](code-review/SKILL.md) | Standards + Spec review; auto-fix obvious; invoked by implement-task or ad-hoc |
| [actions/prototype](actions/prototype/SKILL.md) | Bundle **`wf:prototype`** Method (LOGIC / UI branches via implement-task) |
| [actions/write-code](actions/write-code/SKILL.md) | Default bundle **`wf:task`** Method (TDD build via implement-task) |
| [research](research/SKILL.md) | `wf:research` tickets → findings comment |

Map-free path only: [write-a-prd](../../write-a-prd/SKILL.md) → [prd-to-issues](../../prd-to-issues/SKILL.md).

Cloud AFK automation - see [REFERENCE.md](REFERENCE.md#ecosystem-integration).

## Refer by name

In narration and map sections, use **ticket titles**, not bare `#42`. IDs live inside linked names.
