# Wayfinder reference

## Map slug and decision-log prefix

**Map title:** `{FeatureName}:Map` — PascalCase feature name, literal `:Map` suffix (GitHub issue title, map heading).

**Local filenames:** Use a dot, not a colon — `{FeatureName}.Map.md` and `{FeatureName}.Decision-Log.md`. Colons are invalid in paths on Windows.

**Map slug:** Derive from the feature name for stable IDs:

1. Take the part before `:Map` (e.g. `CommandPalette`).
2. Split on non-alphanumeric boundaries; take significant tokens.
3. Join with hyphens, uppercase, max ~12 chars: `CMD-PAL`, `WF-ECO`, `SEARCH-PANEL`.

**Decision-log prefix:** `{MAP-SLUG}-GM-` + three-digit sequence: `CMD-PAL-GM-001`.

**Subfeature maps:** Child slug **extends** parent when nested: parent `CMD-PAL`, child search UI → `CMD-PAL-SEARCH-GM-001`. Sibling subfeatures under the same parent share the parent prefix segment but distinct suffix: `CMD-PAL-ICONS`, `CMD-PAL-SEARCH`.

**Rules:**

- One authoritative log per map (GitHub issue labelled `wayfinder:decision-log` or section in local map file).
- `grill-me`, **Reconcile**, and ticket resolutions **append** rows; do not renumber existing rows.
- Wayfinder maps implement via [define-bundle](define-bundle/SKILL.md) + [create-tasks](create-tasks/SKILL.md); map-free work may still use `write-a-prd` consolidation.

### Decision-log row format

```markdown
**CMD-PAL-GM-012** — One paragraph: binding decision, constraints, pointers to appendices or tickets.
```

Optional source link: `(from [Palette IA grilling](#123))`

---

## Map body template

Use for GitHub issue body or `wayfinder/plans/{FeatureName}.Map.md`.

```markdown
# {FeatureName}:Map

**Phase:** charting | deciding | consolidating | implementing | complete
**Map slug:** `{MAP-SLUG}`
**Decision log:** [#NNN or path](link) — prefix `{MAP-SLUG}-GM-`

## Target outcome

<What done looks like for this map — usually a buildable PRD. One or two lines.>

## Notes

<Domain; sibling skills to consult; AFK automation repo; preferences for this effort.>

## Subfeatures

<!-- Child maps — integration boundary, not implementation detail store -->

- [{ChildName}:Map](link) — one-line scope

## To Do

<!-- Frontier lives here + open GitHub sub-issues. Agent-maintained table; human verifies on close. -->

| Ticket | Type | Mode | Assignee | Blocked by |
|--------|------|------|----------|------------|
| [Title](issue-link) | research / prototype / grilling / task | HITL / AFK | @user or unclaimed | — |

## Implementing

<!-- Minted implementation tasks from approved bundles. Separate from planning To Do frontier. -->

| Ticket | Bundle | Mode | Status | Blocked by |
|--------|--------|------|--------|------------|

## Completed

<!-- Row moves here after approved reconcile closes the ticket. One-line gist each. -->

- [Title](link) — gist of outcome

## Not yet specified

<!-- In-scope fog — not sharp enough to ticket yet -->

## Out of scope

<!-- Consciously excluded from this map's target outcome -->

## Decision coverage

<!-- Last section on the map body. Operational GM lifecycle index; binding prose lives in decision log only. -->

| GM ID | Status | Linked issue |
|-------|--------|--------------|
| {MAP-SLUG}-GM-001 | open | — |
```

**Section order:** **Decision coverage** is always the **last section** on the map body (after **Out of scope**). **Implementing** sits below **To Do**; planning and implementation frontiers stay separate.

**Index rule:** The map **lists** and **gists**; detail lives in ticket threads and the decision log. Do not paste full `GM-xx` paragraphs into the map body or Decision coverage table.

---

## Ticket types

| Label | Mode typical | Resolved by | Produces |
|-------|--------------|-------------|----------|
| `wayfinder:research` | HITL (v1) | [research](research/SKILL.md) | Structured findings comment; non-binding Proposed tracker updates |
| `wayfinder:prototype` | HITL | Stub code, outline, or `design-an-interface` | Asset link → comment |
| `wayfinder:grilling` | HITL | `grill-me` or `strategic-ideation` when Question is scope/strategy | `{MAP-SLUG}-GM-xx` rows in decision log |
| `wayfinder:task` | HITL or AFK | Agent checklist or human errand | Done-work record → comment |

Every To Do ticket is a **child issue** of the map, labelled `wayfinder:todo`.

### Ticket body template (grilling, prototype, task)

```markdown
## Question

<Single decision or investigation this ticket resolves — one session of work.>

## Map

Parent: [{FeatureName}:Map](#parent-issue-number)
```

### Research ticket template

For `wayfinder:research` tickets — full template in [research REFERENCE](research/REFERENCE.md#research-ticket-template). Required: **Question**, **Done when**, **Map**; optional: **Source hints**, **Perspectives**. Materialize and feature-discovery use this shape.

Script template: [scripts/issue-bodies/research.md](scripts/issue-bodies/research.md).

### Implementation task template

For tasks minted by [create-tasks](create-tasks/SKILL.md) from approved bundles — full template in [create-tasks REFERENCE](create-tasks/REFERENCE.md#task-issue-template). Labels: `wayfinder:task` or `:prototype` + `:hitl` or `:afk`; add **`wayfinder:approved`** when **Status:** `ready`.

---

## Map-discovery artifact

**Default:** [feature-discovery](feature-discovery/SKILL.md) posts a comment on the **map issue** whose body starts with `## Map discovery`. Not a separate issue — an creation-time artifact tied to the map.

| Field | Value |
|-------|--------|
| Created by | feature-discovery (on completion; optional partial comments while in progress) |
| Read by | wayfinder **Materialize** |
| Template | [feature-discovery REFERENCE — map-discovery artifact](feature-discovery/REFERENCE.md#map-discovery-artifact) |

**Local fallback:** `wayfinder/plans/{FeatureName}.Map-Discovery.md` when GitHub is unavailable.

Set **Status:** `ready for materialize` in the comment when discovery is complete.

**Materialize lookup:** `gh issue view <map-num> --comments` — use the latest comment containing `## Map discovery` with **Status:** `ready for materialize`, unless the user points at chat output or a specific comment.

---

## Materialize from map-discovery

Load the artifact from [feature-discovery](feature-discovery/REFERENCE.md#map-discovery-artifact).

| Artifact section | Map / GitHub action |
|------------------|---------------------|
| **Ticket candidates** — sharp Question | Create child issue + **To Do** row |
| **Fog** | Append to map **Not yet specified** |
| **Out of scope suggestions** | Confirm with user; then **Out of scope** |
| **Zone matrix** | Stays on map-discovery comment only; do not paste into map body |
| **Notes** | Merge into map **Notes** if durable |

After materialize: reply on the map-discovery comment thread with **Status:** `materialized`; add **Completed** gist on map (*Map discovery materialized — N tickets*).

**Create order:** Tickets → wire blockers → link sub-issues → update map body.

**Label each ticket:** `wayfinder:todo` + `wayfinder:research` | `:prototype` | `:grilling` | `:task` + `wayfinder:hitl` | `:afk`.

---

## Reconcile resolution template

Post as a **comment on the session ticket** (grilling, research, prototype, task). Non-binding until human approval — same spirit as research **Proposed tracker updates**, but Reconcile may also propose binding GM rows pending approval.

**Section order (fixed):**

```markdown
## Resolution — {ticket title}

### Session summary

<What was settled, deferred, or left open — concrete terms, not zone labels alone.>

### Decision log (proposed append)

**{MAP-SLUG}-GM-NNN** — … `[global]` when infrastructure applies map-wide; omit tag when bundle-scoped.
(from [{ticket title}](#ticket-num))

### Decision coverage (proposed)

| GM ID | Status | Linked issue |
|-------|--------|--------------|
| {MAP-SLUG}-GM-NNN | open / global | — |

### Map updates (proposed)

- **Completed gist:** …
- **Not yet specified:** …
- **Out of scope:** …
- **Notes:** …

### New ticket candidates

| Title | Type | Mode | Question | Blocked by |
|-------|------|------|----------|------------|
| … | research / prototype / grilling / task | HITL / AFK | One sharp Question | — or issue ref |

Omit table when none. Derive **Done when** bullets for `research` rows at materialize time (see [feature-discovery — research ticket shape](feature-discovery/REFERENCE.md#research-ticket-shape-materialize)).

### Bundle cluster suggestions

<!-- Non-binding — define-bundle owns draft/approved bundle issues. Reconcile suggests only. -->

| Suggested name | Covered GM IDs | Rationale | Excluded |
|----------------|----------------|-----------|----------|
| … | GM-012–015 | One vertical slice / subsystem | globals, already bundled, still foggy |

Omit table when no cluster is ready. Rows stay **`open`** until **`bundle approved`**.

### Ticket invalidations

- **Close:** [#N Title](link) — superseded by GM-… / merged into this session
- **Retitle / retype:** …
- **Move to Out of scope:** …

Omit when none.

### Route hint

<One recommended next step + skill — e.g. frontier ticket, define-bundle on cluster above, create-tasks on approved bundle.>

---

Ready for review — reply **Approved — reconcile and close** (or **Approved — reconcile, keep open**) when accepted. Edit any section in this comment before approving.
```

**Draft-only:** Do not append decision log, edit map body, create issues, or close tickets until an approval phrase.

---

## Reconcile inference

Run after loading map, decision log, **Decision coverage**, **To Do**, and the sibling session output.

### All session types

| Signal in session | Propose in resolution |
|-------------------|----------------------|
| Binding decision with durable prose | Decision log row + **Decision coverage** (`open` or `global`) |
| Explicit deferral (“decide later”, “needs spike”) | **Not yet specified** + often a **New ticket candidate** |
| Supersedes an open To Do **Question** | **Ticket invalidations** |
| Mis-scoped or rejected direction | **Out of scope** or invalidation |
| Durable preference for map **Notes** | **Map updates → Notes** |
| Next obvious frontier unchanged | **Route hint** → existing unblocked To Do |

### Grilling (`grill-me`, `strategic-ideation`)

Primary source for holistic inference.

| Signal | Propose |
|--------|---------|
| Branch marked `complete` with binding answers | GM row(s); split distinct decisions into separate IDs |
| Branch deferred with named follow-up | Fog item + typed ticket candidate (`research` for facts; `grilling` for decisions; `prototype` for layout/interface exploration) |
| **Surfaces & experience** — layout agreed but build order unclear | `prototype` or `grilling` ticket; optional bundle cluster if GM cluster is build-ready |
| Multiple `open` GMs describing one deliverable | **Bundle cluster suggestion** (see [define-bundle cluster heuristics](define-bundle/REFERENCE.md#route-heuristics-for-wayfinder)) |
| Infrastructure / cross-cutting constraint | `[global]` on log row; coverage **`global`** |
| Uncertain scope tag | Default **`[global]` when unsure**; human edits in review |

Do **not** create `wayfinder:bundle` issues in Reconcile — narrate clusters for [define-bundle](define-bundle/SKILL.md).

### Research

| Signal | Propose |
|--------|---------|
| Findings imply binding choices | GM rows only after human treats as decided (often via follow-up grilling) |
| **Proposed tracker updates** in findings comment | Merge into **New ticket candidates** / fog / Notes |
| **Done when** satisfied | **Completed gist**; close on full approval |
| Gaps needing more investigation | New `research` ticket or **keep open** |

### Prototype / task

| Signal | Propose |
|--------|---------|
| Asset delivered; decisions captured | GM rows if binding; **Completed gist** |
| Outcome opens new questions | Typed ticket candidates |
| Checklist incomplete | **keep open** or task invalidation |

### Materialize-on-approval (ticket candidates)

When resolution **New ticket candidates** are approved:

1. Create child issues — labels: `wayfinder:todo` + type + mode (same as [Materialize](REFERENCE.md#materialize-from-map-discovery)).
2. Append **To Do** rows; wire **blocked-by** in a second pass.
3. For `research`, include **Done when** in issue body per [research ticket template](research/REFERENCE.md#research-ticket-template).

Skip rows the human struck from the resolution comment before approving.

---

## Completed workflow and approval

1. Sibling skill (or Reconcile draft) → **resolution comment** on ticket ([template above](#reconcile-resolution-template)).
2. Human reviews → edits comment if needed → sends an **approval phrase** (below).
3. Agent **Reconcile** executes approved GitHub/map updates (including ticket materialize when candidates were not removed).
4. If resolution invalidates other tickets → update or close those; move mis-scoped items to **Out of scope**.

### Approval phrases

| User says | Agent may |
|-----------|-----------|
| **Approved — reconcile and close** | Close ticket issue; move row To Do → **Completed** with gist; append decision log **body**; add/update **Decision coverage** row(s) (last map section); update fog/Notes/Out of scope; **create To Do issues** from **New ticket candidates**; apply **Ticket invalidations** |
| **Approved — reconcile, keep open** | Post comment + decision log body + Decision coverage + map notes + optional ticket creates/invalidations; **do not** close or move source ticket to Completed |
| (no approval) | Post draft resolution comment only; **do not** close, edit map, append log, or create issues |

**Not Reconcile:** **`bundle approved`** → [define-bundle](define-bundle/SKILL.md). **`scope approved`** / **`tasks approved`** → [create-tasks](create-tasks/SKILL.md).

Synonyms accepted if unambiguous: “approve and close #N”, “reconcile and close this ticket”.

**Requires:** `gh` authenticated on the target repo for issue close/edit.

### Implementation task Reconcile

When closing a **`wayfinder:approved`** implementation task (from [create-tasks](create-tasks/SKILL.md)):

| User says | Agent may |
|-----------|-----------|
| **Approved — reconcile and close** | Close task; remove **`wayfinder:approved`**; move **Implementing** row → **Completed** gist; set Decision coverage **`implemented`** for bundle-scoped GMs shipped by this task |

Synonyms accepted if unambiguous: "reconcile and close task #N".

See [create-tasks REFERENCE — implementation Reconcile](create-tasks/REFERENCE.md#implementation-reconcile).

---

## Frontier queries

**Frontier** = rows in **To Do** whose linked issues are: **open**, **unblocked** (all blockers closed), **unclaimed** (no assignee) or assigned to current worker per session rules.

Use GitHub’s blocked-by graph for ordering. Open tickets not listed in **To Do** should not exist — the table is the human-facing frontier index.

---

## Routing table

Suggest-only — user starts the recommended skill. Map ticket **Type** → default skill:

| Ticket type | Default skill | Notes |
|-------------|---------------|-------|
| `grilling` (implementation) | `grill-me` | Ticket **Question** is the grill seed |
| `grilling` (scope/strategy) | `strategic-ideation` | When **Question** is bundling, roadmap, or scope shape |
| `research` | `research` | HITL v1; AFK deferred per map Notes |
| `prototype` | `design-an-interface` or inline stub | Per ticket **Question** |
| `task` | Agent checklist or human | Per ticket **Question** |
| GM cluster ready to build | `define-bundle` | While planning To Do or fog may stay open; see [define-bundle REFERENCE](define-bundle/REFERENCE.md#route-heuristics-for-wayfinder) |
| Approved bundle | `create-tasks` | Splits into **Implementing** tasks |
| Small scope, no map | `write-a-prd` → `prd-to-issues` | **Not** a map Route handoff |
| New feature, no map | wayfinder **Chart** | Then `feature-discovery` |

After sibling session: remind user to invoke wayfinder **Reconcile** (explicit invoke — see map fog if auto-reconcile is ever desired).

---

## Subfeature maps

**When:** A zone or subsystem is large enough for its own discovery + ticket graph but must stay consistent with the parent.

**Steps:**

1. **Chart** child map `{SubFeatureName}:Map` with its own slug and decision log.
2. Add link under parent **Subfeatures** with one-line boundary ("Owns search UX; parent owns shell IA").
3. Add parent **To Do** ticket if needed: *Integration review — align `{Child}-GM-*` with `{Parent}-GM-*`* (grilling, blocked by child frontier empty or milestone).
4. Child **Notes** must link parent map and list parent `GM-xx` rows that constrain it.

Cross-map conflicts → parent grilling ticket, not silent edits to child logs.

---

## Tracker operations

**Default:** GitHub issues on the **target repo** (or `KroniK907/skills` for meta/skills work). This is the **canonical tracker** when issues are enabled.

| Artifact | Label |
|----------|--------|
| Map | `wayfinder:map` |
| Decision log | `wayfinder:decision-log` |
| Build bundle | `wayfinder:bundle` |
| To Do ticket | `wayfinder:todo` + type + mode |
| Implementation task (draft) | `wayfinder:task` or `:prototype` + `:hitl` or `:afk` |
| Approved implementation task | above + **`wayfinder:approved`**; body **Status:** `ready` |

Map-discovery artifact = **comment on map issue** (no label).

**Sub-issues:** Link map → decision log and tickets via GitHub sub-issues. **Blocked-by:** Use native issue dependencies for frontier ordering.

**Chart create order:** Decision log → map (with log link) → link sub-issues → hand off to feature-discovery.

**Local fallback:** `wayfinder/plans/{FeatureName}.Map.md` and `{FeatureName}.Map-Discovery.md` only when GitHub is unavailable, or for export. See [plans/README.md](plans/README.md). Do not commit local files that duplicate active GitHub issues.

---

## Skills repo layout

In `KroniK907/skills`, ecosystem skills live under **`wayfinder/<skill>/`**. Map-frontier siblings (`feature-discovery`, `grill-me`, `research`, `define-bundle`, `create-tasks`, etc.) are peers of this hub skill. **`wayfinder/actions/`** holds **`implement-task` Method playbooks** — see [actions/PATTERNS.md](actions/PATTERNS.md). One-off utilities (`tdd`, `commit`, `write-a-prd`, PRD tools, etc.) stay at repo root.

**Method path validation** (WF-ECO-GM-030): default pool is skills at **`wayfinder/**/<name>/SKILL.md`** in the pinned pack. Repo-root one-offs are valid only when **## Method** explicitly names them.

| Path | Role |
|------|------|
| `wayfinder/SKILL.md` | Hub — Chart, Materialize, Reconcile, Route |
| `wayfinder/<skill>/` | Map-frontier sibling skills |
| `wayfinder/actions/<name>/` | Implementation Method playbooks (via `implement-task`) |
| `<one-off>/` (repo root) | Map-free or standalone utilities |

Install example: `npx skills@latest add KroniK907/skills/wayfinder/research`

---

## Ecosystem integration

Skills that **read** wayfinder maps:

| Skill | Reads | Writes |
|-------|-------|--------|
| `write-a-prd` | Map Completed + decision log | PRD issue |
| `prd-to-issues` | PRD | `agent-queue` issues |

Skills that **write** wayfinder state:

| Skill | Writes |
|-------|--------|
| `feature-discovery` | Map-discovery comment on map issue |
| `strategic-ideation` | Scope handoff (chat); Reconcile records on map |
| `grill-me` | Decision log `{MAP-SLUG}-GM-xx`; resolution comment on grilling ticket; Reconcile proposes holistic tracker delta (tickets, bundle clusters, route) |
| `define-bundle` | Draft/approved bundle issue; Decision coverage `scoped`; log `- bundled via [#N]` suffixes |
| `create-tasks` | Implementation task issues; **Implementing** table; coverage `assigned` on scope approval; `implemented` on Reconcile close |
| `wayfinder` | Map To Do / Completed / fog / Subfeatures; ticket create/close on approval |
| [research](research/SKILL.md) | Findings comment on research ticket; non-binding Proposed tracker updates |
| Cloud AFK automation (TBD) | PR + comment on AFK ticket; human approves close via Reconcile |

**Handoff chain:** Chart → feature-discovery → Materialize → sibling skills → Reconcile → `define-bundle` → `create-tasks` → implement → Reconcile. Map-free: grill-me → `write-a-prd` → `prd-to-issues`.

---

## Chart handoff message (template)

Post or narrate after skeleton creation:

```markdown
**Wayfinder Chart complete** — [{FeatureName}:Map](link) · decision log [#N](link)

**Next:** Run [feature-discovery](feature-discovery/SKILL.md) with:
- Map: #N
- Seed: …
- Target outcome: …

When the map-discovery comment has **Status:** `ready for materialize`, invoke **wayfinder Materialize** with the map link.
```
