# Wayfinder reference

## Map slug and decision-log prefix

**Map title:** `{FeatureName}:Map` - PascalCase feature name, literal `:Map` suffix (GitHub issue title, map heading).

**Local filenames:** Use a dot, not a colon - `{FeatureName}.Map.md` and `{FeatureName}.Decision-Log.md`. Colons are invalid in paths on Windows.

**Map slug:** Derive from the feature name for stable IDs:

1. Take the part before `:Map` (e.g. `CommandPalette`).
2. Split on non-alphanumeric boundaries; take significant tokens.
3. Join with hyphens, uppercase, max ~12 chars: `CMD-PAL`, `WF-ECO`, `SEARCH-PANEL`.

**Decision-log prefix:** `{MAP-SLUG}-GM-` + three-digit sequence: `CMD-PAL-GM-001`.

**Subfeature maps:** Child slug **extends** parent when nested: parent `CMD-PAL`, child search UI → `CMD-PAL-SEARCH-GM-001`. Sibling subfeatures under the same parent share the parent prefix segment but distinct suffix: `CMD-PAL-ICONS`, `CMD-PAL-SEARCH`.

**Rules:**

- One authoritative log per map (GitHub issue labelled `wf:decision-log` or section in local map file).
- `grill-me`, **Reconcile**, and ticket resolutions append rows as **comments** on the log issue (`wf log-append`). Do not rewrite the log body to add a row. A later entry with the same ID replaces the earlier one. Do not renumber existing rows. Binding text is the issue body plus those comments.
- Wayfinder maps implement via [define-bundle](actions/define-bundle/SKILL.md) + [create-tasks](actions/create-tasks/SKILL.md); map-free work may still use `write-a-prd` consolidation.

### Decision-log row format

```markdown
**CMD-PAL-GM-012** - One paragraph: binding decision, constraints, pointers to appendices or tickets.
```

Optional source link: `(from [Palette IA grilling](#123))`

---

## Map body template

Use for GitHub issue body or `wayfinder/utilities/plans/{FeatureName}.Map.md`.

```markdown
# {FeatureName}:Map

**Phase:** charting | deciding | consolidating | implementing | complete
**Map slug:** `{MAP-SLUG}`
**Decision log:** [#NNN or path](link) - prefix `{MAP-SLUG}-GM-`

## Target outcome

<What done looks like for this map - usually a buildable PRD. One or two lines.>

## Notes

<Domain; sibling skills to consult; AFK automation repo; preferences for this effort.>

## Subfeatures

<!-- Child maps - integration boundary, not implementation detail store -->

- [{ChildName}:Map](link) - one-line scope

## To Do

<!-- Frontier lives here + open GitHub sub-issues. Agent-maintained table; human verifies on close. -->

| Ticket | Type | Mode | Assignee | Blocked by |
|--------|------|------|----------|------------|
| [Title](issue-link) | research / prototype / grilling / task | HITL / AFK | @user or unclaimed | - |

## Implementing

<!-- Minted implementation tasks from approved bundles. Separate from planning To Do frontier. -->

| Ticket | Bundle | Mode | Status | Blocked by |
|--------|--------|------|--------|------------|

## Completed

<!-- Row moves here after approved reconcile closes the ticket. One-line gist each. -->

- [Title](link) - gist of outcome

## Not yet specified

<!-- In-scope fog - not sharp enough to ticket yet -->

## Out of scope

<!-- Consciously excluded from this map's target outcome -->

## Decision coverage

<!-- Last section on the map body. Operational GM lifecycle index; binding prose lives in decision log only. -->

| GM ID | Status | Linked issue |
|-------|--------|--------------|
| {MAP-SLUG}-GM-001 | open | - |
```

**Section order:** **Decision coverage** is always the **last section** on the map body (after **Out of scope**). **Implementing** sits below **To Do**; planning and implementation frontiers stay separate.

**Index rule:** The map **lists** and **gists**; detail lives in ticket threads and the decision log. Do not paste full `GM-xx` paragraphs into the map body or Decision coverage table.

---

## Ticket types

| Label | Mode typical | Resolved by | Produces |
|-------|--------------|-------------|----------|
| `wf:research` | HITL (v1) | [research](actions/research/SKILL.md) | Structured findings comment; non-binding Proposed tracker updates |
| `wf:prototype` | HITL | Stub code, outline, or `design-modules` | Asset link → comment |
| `wf:grilling` | HITL | `grill-me` or `strategic-ideation` when Question is scope/strategy | `{MAP-SLUG}-GM-xx` rows in decision log |
| `wf:task` | HITL or AFK | Agent checklist or human errand | Done-work record → comment |

Every To Do ticket is a **child issue** of the map, labelled `wf:todo`.

### Ticket title conventions

Issue **titles** are the first signal agents and humans see in map **To Do** rows, Route suggestions, and chat citations. Use a **type prefix** so the title names the **process** (for planning tickets) or the **deliverable** (for implementation tickets) - not the other way around.

| Prefix | Type label (typical) | Route to skill(s) | Title names |
|--------|----------------------|-------------------|-------------|
| **Grill:** | `wf:grilling` | [grill-me](ideation/grill-me/SKILL.md) | The decision or contract to stress-test - depth-first Q&A |
| **Ideate:** | `wf:grilling` | [strategic-ideation](ideation/strategic-ideation/SKILL.md) - [feature-ideation](feature-ideation/SKILL.md) (stub → strategic-ideation) | Scope/strategy or feature shape to expand → tension → prune |
| **Constrain:** | `wf:grilling` | [constrain-fog](ideation/constrain-fog/SKILL.md) | A **Not yet specified** fog line or cluster to sharpen |
| **Research:** | `wf:research` | [research](actions/research/SKILL.md) | The investigation - facts, prior art, survey |
| **Prototype:** | `wf:prototype` | [prototype](actions/prototype/SKILL.md) - [design-modules](actions/design-modules/SKILL.md) on planning **To Do** | What to explore - throwaway demo, layout, interface variants |
| **Task:** | `wf:task` | [one-off](orchestrators/one-off/SKILL.md) - [implement-task](orchestrators/implement-task/SKILL.md) - [define-bundle](actions/define-bundle/SKILL.md) | What to ship or bundle - repo deliverable, implementation run, or GM cluster |
| **Organize:** | `wf:task` | [wayfinder](SKILL.md) - [one-off](orchestrators/one-off/SKILL.md) | Tracker/map housekeeping - fog sort, map sync, labels, tables, trivial edits |

**Promoted bundle issues** (output of define-bundle, not frontier entry titles) stay **`Bundle: {short name}`** per [define-bundle REFERENCE](actions/define-bundle/REFERENCE.md). **Implementing** tasks from create-tasks stay **`Task: {short name}`**.

### Route by prefix (when several skills share a prefix)

| Prefix | Pick skill when… |
|--------|------------------|
| **Constrain:** | Map-scoped fog grooming → **constrain-fog** (auto-creates session ticket; artifact on ticket body). |
| **Ideate:** | Scope/strategy expand → tension → prune → **strategic-ideation** (default). Legacy **feature-ideation** invoke resolves to the same skill. |
| **Prototype:** | Planning **To Do** exploration → **design-modules** or inline stub. **`wf:approved`** on **Implementing** → **implement-task** → Method **prototype**. |
| **Task:** | Map **To Do** repo deliverable → **one-off**. **`wf:approved`** on **Implementing** → **implement-task**. User wants to group GM cluster → **define-bundle**. |
| **Organize:** | Chart / Materialize / Reconcile / map table or **Not yet specified** edits → **wayfinder**. Trivial checklist only → **one-off**. |

**Rules:**

1. **Planning titles = process, not artifact** - `Grill:` / `Ideate:` / `Constrain:` / `Research:` / `Prototype:` titles describe the **session**; put binding deliverables in **## Question** or a follow-on **`Task:`** ticket after Reconcile.
2. **Implementation titles = deliverable** - `Task:` titles name what lands in the repo, bundle, or tracker when done.
3. **Prefix matches Type column** - map **To Do** Type must agree with the title prefix; fix mismatches via retitle or retype.
4. **Materialize and Reconcile apply prefixes** - when creating issues from **Ticket candidates**, set the title from the Type column using this table (do not copy the Question verbatim as the title).
5. **Noun after prefix is short** - topic or subsystem name; details live in the issue body.

**Examples:**

| Weak title | Strong title | Why |
|------------|--------------|-----|
| Design constrain-fog skill for map fog resolution | **Grill:** constrain-fog skill design | "Design … skill" reads like implementation; `Grill:` signals Q&A first |
| Specify research ticket workflow | **Grill:** research ticket workflow | "Specify" is ambiguous; grilling resolves the contract |
| Cloud automations for AFK pickup | **Research:** cloud automations for AFK pickup | Names the investigation |
| Subfeature map worked example in REFERENCE | **Prototype:** subfeature map worked example | Names exploration, not a shipped doc yet |
| implement create-tasks skill | **Task:** implement create-tasks skill | Deliverable prefix |
| Group GM-012-015 into first bundle | **Task:** define-bundle for palette shell | Bundling work; Route → define-bundle |
| Decision coverage backfill | **Organize:** decision coverage backfill | Tracker errand; Route → wayfinder or one-off |
| Clear routing-table fog lines | **Organize:** routing table fog | Tracker sort; Route → wayfinder |

**Agent cue:** When the user cites a map ticket by `#N` or title, read the **prefix** first - it narrows the skill set. When the prefix maps to **one** skill, start there. When it maps to **several**, use **## Question** and map context (To Do vs Implementing, fog vs deliverable vs map sync) per the table above - do not treat body prose as permission to skip the prefix family (e.g. `Grill:` → implement).

### Ticket body template (grilling, prototype, task)

```markdown
## Question

<Single decision or investigation this ticket resolves - one session of work.>

## Map

Parent: [{FeatureName}:Map](#parent-issue-number)
```

### Research ticket template

For `wf:research` tickets - full template in [research REFERENCE](actions/research/REFERENCE.md#research-ticket-template). Required: **Question**, **Done when**, **Map**; optional: **Source hints**, **Perspectives**. Materialize and feature-discovery use this shape.

Script template: [scripts/issue-bodies/research.md](scripts/issue-bodies/research.md).

### Implementation task template

For tasks minted by [create-tasks](actions/create-tasks/SKILL.md) from approved bundles - full template in [create-tasks REFERENCE](actions/create-tasks/REFERENCE.md#task-issue-template). Labels: `wf:task` or `:prototype` + `:hitl` or `:afk`; add **`wf:approved`** when **Status:** `ready`.

---

## Map-discovery artifact

**Default:** [feature-discovery](ideation/feature-discovery/SKILL.md) posts a comment on the **map issue** whose body starts with `## Map discovery`. Not a separate issue - an creation-time artifact tied to the map.

| Field | Value |
|-------|--------|
| Created by | feature-discovery (on completion; optional partial comments while in progress) |
| Read by | wayfinder **Materialize** |
| Template | [feature-discovery REFERENCE - map-discovery artifact](ideation/feature-discovery/REFERENCE.md#map-discovery-artifact) |

**Local fallback:** `wayfinder/utilities/plans/{FeatureName}.Map-Discovery.md` when GitHub is unavailable.

Set **Status:** `ready for materialize` in the comment when discovery is complete.

**Materialize lookup:** `gh issue view <map-num> --comments` - use the latest comment containing `## Map discovery` with **Status:** `ready for materialize`, unless the user points at chat output or a specific comment.

---

## Materialize from map-discovery

Load the artifact from [feature-discovery](ideation/feature-discovery/REFERENCE.md#map-discovery-artifact).

| Artifact section | Map / GitHub action |
|------------------|---------------------|
| **Ticket candidates** - sharp Question | Create child issue + **To Do** row; title per [ticket title conventions](REFERENCE.md#ticket-title-conventions) |
| **Fog** | Append to map **Not yet specified** |
| **Out of scope suggestions** | Confirm with user; then **Out of scope** |
| **Zone matrix** | Stays on map-discovery comment only; do not paste into map body |
| **Notes** | Merge into map **Notes** if durable |

After materialize: reply on the map-discovery comment thread with **Status:** `materialized`; add **Completed** gist on map (*Map discovery materialized - N tickets*).

**Create order:** Tickets → wire blockers → link sub-issues → update map body.

**Label each ticket:** `wf:todo` + `wf:research` | `:prototype` | `:grilling` | `:task` + `wf:hitl` | `:afk`.

---

## Reconcile

Resolution template, inference, approval phrases, and the wf apply steps live in [references/reconcile.md](references/reconcile.md).

- [Resolution template](references/reconcile.md#reconcile-resolution-template)
- [Inference](references/reconcile.md#reconcile-inference)
- [Approval phrases](references/reconcile.md#approval-phrases)
- [Map and issue body edits](references/reconcile.md#map-and-issue-body-edits-reconcile)

---

## Frontier queries

**Frontier** = rows in **To Do** whose linked issues are: **open**, **unblocked** (all blockers closed), **unclaimed** (no assignee) or assigned to current worker per session rules.

Use GitHub’s blocked-by graph for ordering. Open tickets not listed in **To Do** should not exist - the table is the human-facing frontier index.

---

## Routing table

Suggest-only - user starts the recommended skill. Map ticket **Type** → default skill:

| Ticket type | Default skill | Notes |
|-------------|---------------|-------|
| `grilling` (implementation) | `grill-me` | Ticket **Question** is the grill seed |
| `grilling` (scope/strategy) | `strategic-ideation` | When **Question** is bundling, roadmap, or scope shape |
| `research` | `research` | HITL v1; AFK deferred per map Notes |
| `prototype` (To Do) | `design-modules` or inline stub | Planning frontier; per ticket **Question** |
| `prototype` (Implementing) | [implement-task](orchestrators/implement-task/SKILL.md) → Method **`prototype`** | [actions/prototype](actions/prototype/SKILL.md); bundle tasks only |
| `task` (Implementing) | [implement-task](orchestrators/implement-task/SKILL.md) | Default Method **`write-code`** for normal build; **`prototype`** for throwaway demos |
| Approved bundle (post-approval) | [design-modules](actions/design-modules/SKILL.md) | Optional HITL modules shaping (one or more) before [create-tasks](actions/create-tasks/SKILL.md) |
| `task` (To Do) | [one-off](orchestrators/one-off/SKILL.md) | Map-scoped repo deliverables; trivial checklist-only errands stay *Agent checklist or human* |
| GM cluster ready to build | `define-bundle` | While planning To Do or fog may stay open; see [define-bundle REFERENCE](actions/define-bundle/REFERENCE.md#route-heuristics-for-wayfinder) |
| Approved bundle | `create-tasks` | Splits into **Implementing** tasks. Each task later gets a worktree and a pull request |
| Small scope, no map | `write-a-prd` → `prd-to-issues` | **Not** a map Route handoff |
| New feature, no map | wayfinder **Chart** | Then `feature-discovery` |

After sibling session: remind user to invoke wayfinder **Reconcile** (explicit invoke - see map fog if auto-reconcile is ever desired).

### Route heuristics - constrain-fog

Suggest [constrain-fog](ideation/constrain-fog/SKILL.md) when **all** of:

1. Map **To Do** table is **empty** (no open frontier rows)
2. **Not yet specified** is **non-empty**

**Never** auto-suggest constrain-fog when open **To Do** items exist - user may **explicitly invoke** constrain-fog anytime.

When **To Do** has items, Route the planning frontier per [frontier queries](#frontier-queries) instead.

**Contrast:**

| Skill | When |
|-------|------|
| **feature-discovery** | Post-Chart whole-map breadth-first triage → **`## Map discovery`** → Materialize |
| **constrain-fog** | Existing map; groom **Not yet specified** → **`## Fog resolution`** on **`Constrain:`** ticket → Reconcile |
| **grill-me / strategic-ideation / research** | Existing typed tickets on **To Do** |

---

## Subfeature maps

**When:** A zone or subsystem is large enough for its own discovery + ticket graph but must stay consistent with the parent.

**Steps:**

1. **Chart** child map `{SubFeatureName}:Map` with its own slug and decision log.
2. Add link under parent **Subfeatures** with one-line boundary ("Owns search UX; parent owns shell IA").
3. Add parent **To Do** ticket if needed: *Integration review - align `{Child}-GM-*` with `{Parent}-GM-*`* (grilling, blocked by child frontier empty or milestone).
4. Child **Notes** must link parent map and list parent `GM-xx` rows that constrain it.

Cross-map conflicts → parent grilling ticket, not silent edits to child logs.

---

## Tracker operations

**Default:** GitHub issues on the **target repo** (or `KroniK907/agent-config` for meta/skills work). This is the **canonical tracker** when issues are enabled.

| Artifact | Label |
|----------|--------|
| Map | `wf:map` |
| Decision log | `wf:decision-log` |
| Build bundle | `wf:bundle` - no feature branch on the bundle |
| To Do ticket | `wf:todo` + type + mode |
| Implementation task (draft) | `wf:task` or `:prototype` + `:hitl` or `:afk` |
| Approved implementation task | above + **`wf:approved`**; body **Status:** `ready` \| `awaiting-reconcile` |
| AFK run lock | **`wf:afk-running`** on current AFK task |
| Awaiting approval | **`wf:needs-review`** - add when agent posts draft awaiting human gate phrase; remove when phrase received |

Map-discovery artifact = **comment on map issue** (no label).

### `wf:needs-review`

Bright-red queue signal: an agent finished a draft step and a **human approval phrase** is required before tracker writes or close.

| Add label | When | Remove label |
|-----------|------|--------------|
| [define-bundle](actions/define-bundle/SKILL.md) | Draft bundle issue posted | **`bundle approved`** |
| [create-tasks](actions/create-tasks/SKILL.md) | Draft task issue(s) posted | **`tasks approved`** (or **`scope approved`** if no further task promotion pending) |
| [wayfinder](SKILL.md) **Reconcile** | Resolution draft comment posted on session ticket | **`Approved - reconcile and close`** or **`Approved - reconcile, keep open`** |
| [implement-task](orchestrators/implement-task/SKILL.md) | Success end-of-run (**Status:** `awaiting-reconcile`) | wayfinder **Reconcile** on **`Approved - reconcile and close`** |

```powershell
gh issue edit <num> --add-label "wf:needs-review"
gh issue edit <num> --remove-label "wf:needs-review"
```

**Sub-issues:** Link map → decision log and tickets via GitHub sub-issues. **Blocked-by:** Use native issue dependencies for frontier ordering.

**Chart create order:** Decision log → map (with log link) → link sub-issues → hand off to feature-discovery.

**Local fallback:** `wayfinder/utilities/plans/{FeatureName}.Map.md` and `{FeatureName}.Map-Discovery.md` only when GitHub is unavailable, or for export. See [plans/README.md](utilities/plans/README.md). Do not commit local files that duplicate active GitHub issues.

---

## Skills repo layout

**Method path validation:** default pool is skills at **`skills/wayfinder/**/<name>/SKILL.md`** in the pinned pack. Skills under `skills/` are valid only when **## Method** explicitly names them. AFK app repos: [AFK-BOOTSTRAP.md](utilities/AFK-BOOTSTRAP.md).

---

## Ecosystem integration

Per-skill roles are in the [SKILL.md Ecosystem table](SKILL.md#ecosystem-related-skills). Extra detail:

- **Map readers (map-free path):** `write-a-prd` reads Completed + decision log and writes a PRD issue; `prd-to-issues` turns it into `agent-queue` issues.
- **Coverage transitions:** `define-bundle` sets `scoped` on **`bundle approved`** (no git branch). `create-tasks` sets `assigned` on scope approval, adds **`wf:approved`** on **`tasks approved`**; Reconcile sets `implemented`.
- **Cloud AFK automation:** comment **`Approved - AFK implement`** plus label **`wf:approved`** runs [implement-task](orchestrators/implement-task/SKILL.md) in a task worktree (push, pull request, resolution comment); a human Reconcile closes the task. Setup: [AFK-BOOTSTRAP.md](utilities/AFK-BOOTSTRAP.md).
- **Route hint:** a request to review a branch, PR, WIP, or diff outside implement-task goes to [code-review](actions/code-review/SKILL.md) ad-hoc. During implement-task, code-review runs automatically after Method.

---

## Chart handoff message (template)

Post or narrate after skeleton creation:

```markdown
**Wayfinder Chart complete** - [{FeatureName}:Map](link) - decision log [#N](link)

**Next:** Run [feature-discovery](ideation/feature-discovery/SKILL.md) with:
- Map: #N
- Seed: …
- Target outcome: …

When the map-discovery comment has **Status:** `ready for materialize`, invoke **wayfinder Materialize** with the map link.
```
