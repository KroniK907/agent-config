---
name: wayfinder
description: wayfinder, FeatureName:Map, wf:map, wf:decision-log, Chart handoff, Materialize, map-discovery artifact, Reconcile, Route, frontier, To Do tickets, sync chat to map, starting large feature, subfeature map, wf:needs-review, navigate wayfinder across sessions
disable-model-invocation: true
agent-config-sync: true
---

# Wayfinder

**Tracker and router** for large features: map skeleton → [feature-discovery](ideation/feature-discovery/SKILL.md) → tickets → sibling skills → **Reconcile** → [define-bundle](actions/define-bundle/SKILL.md) → [create-tasks](actions/create-tasks/SKILL.md). Templates, materialize rules, routing table, and GitHub ops: [REFERENCE.md](REFERENCE.md).

Wayfinder plans; it does not implement unless the map **Notes** say so. It does not run discovery, ideation, or grilling - it updates GitHub state and suggests the next skill.

## Working with the user

Act as a coding partner. When the user signs off on a draft - in any wording - act on it. When the conversation already shows the work is done and accepted, go ahead. Ask one short question only when you genuinely can't tell whether the user considers the work ready, or when a step would close, delete, or supersede something they may still want.

## When to use

| Situation | Mode |
|-----------|------|
| Rough feature idea, too big for one chat | **Chart** - skeleton → hand off to feature-discovery |
| Discovery capture ready | **Materialize** - create To Do tickets from the map-discovery artifact |
| Sibling skill session finished | **Reconcile** - resolution comment, decision log, map, close ticket |
| Existing map; need next step | **Route** - frontier + skill suggestion |
| Child subsystem needs its own planning | **Chart** a subfeature map; link from parent **Subfeatures** |

Skip wayfinder when the path is clear - use `sous-vide` for a grilling session, or implement directly.

## Modes

### Chart

1. **Seed** - the user's feature description (and target repo if not obvious).
2. **Name** - `{FeatureName}:Map`; derive the [map slug](REFERENCE.md#map-slug-and-decision-log-prefix).
3. **Target outcome** - one or two lines, usually a buildable PRD.
4. **Create** - decision log issue (`wf:decision-log`), then map issue (`wf:map`) linking it. **To Do** empty; **Phase:** `charting`. Local fallback: [plans/](utilities/plans/README.md).
5. **Hand off** - point the user to [feature-discovery](ideation/feature-discovery/SKILL.md) with the map link, target outcome, and seed. Chart creates no tickets.

### Materialize

1. **Load** the map and the map-discovery artifact: a `## Map discovery` block in chat, the latest map comment with **Status:** `ready for materialize`, or a user paste.
2. **Create To Do tickets** - one child issue per **Ticket candidates** row, labelled `wf:todo` + type + mode, titled per [ticket title conventions](REFERENCE.md#ticket-title-conventions). Wire blocked-by in a second pass.
3. **Update the map** - **To Do** table; **Fog** → **Not yet specified**; confirm **Out of scope suggestions** with the user; **Phase:** `deciding`; **Completed** gist *Map discovery materialized - N tickets*. Reply on the discovery comment with **Status:** `materialized`.
4. **Route** to the first frontier ticket.

Full rules: [materialize from map-discovery](REFERENCE.md#materialize-from-map-discovery).

### Reconcile

Turns a finished sibling session (grilling, research, prototype, task, constrain-fog) into tracker state. Details: [references/reconcile.md](references/reconcile.md).

1. **Load** the map, the session ticket, and the session output.
2. **Infer** the tracker delta per [reconcile inference](references/reconcile.md#reconcile-inference): decision-log rows, **Decision coverage**, map updates, new ticket candidates, bundle cluster suggestions, ticket invalidations, route hint.
3. **Post** the [resolution comment](references/reconcile.md#resolution-comment) on the ticket.
4. **Apply** with [`wf`](references/reconcile.md#apply-steps): append the log, edit the map, create tickets, apply invalidations. Close the ticket when its work is done; leave it open when gaps remain and say why.

If the user already said the session's outcome is accepted, apply in the same pass. Otherwise add **`wf:needs-review`**, summarize the delta in chat, and apply once they're happy with it. Reconcile only *suggests* bundle clusters - [define-bundle](actions/define-bundle/SKILL.md) creates bundles.

### Route

1. **Load** the map body (low-res) and decision log link.
2. **Frontier** - first open, unblocked, unclaimed **To Do** item ([frontier queries](REFERENCE.md#frontier-queries)).
3. **Implementation path** - if **Decision coverage** has a build-ready cluster of `open` rows, suggest [define-bundle](actions/define-bundle/SKILL.md). Approved bundles → [create-tasks](actions/create-tasks/SKILL.md). **`wf:approved`** **Implementing** tasks → [implement-task](orchestrators/implement-task/SKILL.md).
4. **Suggest** one next step and skill from the [routing table](REFERENCE.md#routing-table); a second only if genuinely ambiguous.

**Done when:** you have named one skill and one ticket (or bundle).

## Decision log

Each map owns a scoped log (`{MAP-SLUG}-GM-NNN`). Rows are appended as comments on the log issue. Rules: [REFERENCE.md](REFERENCE.md#map-slug-and-decision-log-prefix).

## Subfeatures

Large work may spawn child maps (`SearchPanel:Map`) linked under the parent's **Subfeatures**. See [REFERENCE.md](REFERENCE.md#subfeature-maps).

## Ecosystem

| Skill | Role |
|-------|------|
| [feature-discovery](ideation/feature-discovery/SKILL.md) | After Chart - posts the map-discovery comment |
| [constrain-fog](ideation/constrain-fog/SKILL.md) | Groom **Not yet specified** on a `Constrain:` ticket |
| [strategic-ideation](ideation/strategic-ideation/SKILL.md) | Scope/strategy: expand → tension → prune |
| [sous-vide](orchestrators/sous-vide/SKILL.md) | Orchestrate a grilling session |
| [grill-me](ideation/grill-me/SKILL.md) | Asks the sous-vide **Ask** list one question at a time |
| [design-modules](actions/design-modules/SKILL.md) | Module interface shaping before create-tasks |
| [define-bundle](actions/define-bundle/SKILL.md) | GM cluster → `wf:bundle` issue |
| [create-tasks](actions/create-tasks/SKILL.md) | Approved bundle → **Implementing** tasks |
| [one-off](orchestrators/one-off/SKILL.md) | Implement a map **To Do** ticket without the bundle pipeline |
| [implement-task](orchestrators/implement-task/SKILL.md) | **`wf:approved`** task → Method → code-review → PR → `awaiting-reconcile` |
| [code-review](actions/code-review/SKILL.md) | Standards + Spec review |
| [prototype](actions/prototype/SKILL.md) | **`wf:prototype`** Method |
| [write-code](actions/write-code/SKILL.md) | Default **`wf:task`** Method |
| [research](actions/research/SKILL.md) | `wf:research` tickets → findings comment |

Map-free path: [write-a-prd](../../write-a-prd/SKILL.md) → [prd-to-issues](../../prd-to-issues/SKILL.md). Cloud AFK: [REFERENCE.md](REFERENCE.md#ecosystem-integration).

In narration and map sections, refer to tickets by title, not bare `#42`.
