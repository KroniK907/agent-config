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
- `grill-me` and ticket resolutions **append** rows; do not renumber until `write-a-prd` consolidates.
- PRD consolidation may merge/dedupe into plain `GM-001` or keep prefixed IDs — document choice in map **Notes**.

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

## Completed

<!-- Human moves rows here after closing the ticket issue. One-line gist each. -->

- [Title](link) — gist of outcome

## Not yet specified

<!-- In-scope fog — not sharp enough to ticket yet -->

## Out of scope

<!-- Consciously excluded from this map's target outcome -->
```

**Index rule:** The map **lists** and **gists**; detail lives in ticket threads and the decision log. Do not paste full `GM-xx` paragraphs into the map body.

---

## Ticket types

| Label | Mode typical | Resolved by | Produces |
|-------|--------------|-------------|----------|
| `wayfinder:research` | AFK | Task subagent / future `research` skill | Facts, links, constraints → ticket comment |
| `wayfinder:prototype` | HITL | Stub code, outline, or `design-an-interface` | Asset link → comment |
| `wayfinder:grilling` | HITL | `grill-me` scoped to ticket **Question** | `{MAP-SLUG}-GM-xx` rows in decision log |
| `wayfinder:task` | HITL or AFK | Agent checklist or human errand | Done-work record → comment |

Every To Do ticket is a **child issue** of the map, labelled `wayfinder:todo`.

### Ticket body template

```markdown
## Question

<Single decision or investigation this ticket resolves — one session of work.>

## Map

Parent: [{FeatureName}:Map](#parent-issue-number)
```

---

## Completed workflow

1. Agent finishes work → **resolution comment** on ticket (never closes).
2. Human reviews → **closes** the GitHub issue (or marks done in local tracker).
3. Human or agent-with-explicit-instruction updates map: remove from **To Do** table, append to **Completed** with gist.
4. If resolution invalidates other tickets → close or update those tickets; move mis-scoped items to **Out of scope**.

**Agents must not** close `wayfinder:todo` issues. Cloud automations deliver PRs/comments; human still closes.

---

## Frontier queries

**Frontier** = rows in **To Do** whose linked issues are: **open**, **unblocked** (all blockers closed), **unclaimed** (no assignee) or assigned to current worker per session rules.

Use GitHub’s blocked-by graph for ordering. Open tickets not listed in **To Do** should not exist — the table is the human-facing frontier index.

---

## Subfeature maps

**When:** A zone or subsystem is large enough for its own ideation + ticket graph but must stay consistent with the parent.

**Steps:**

1. Create child map `{SubFeatureName}:Map` with its own slug and decision log.
2. Add link under parent **Subfeatures** with one-line boundary ("Owns search UX; parent owns shell IA").
3. Add parent **To Do** ticket if needed: *Integration review — align `{Child}-GM-*` with `{Parent}-GM-*`* (grilling, blocked by child frontier empty or milestone).
4. Child **Notes** must link parent map and list parent `GM-xx` rows that constrain it.

Cross-map conflicts → parent grilling ticket, not silent edits to child logs.

---

## Tracker operations

**Default:** GitHub issues on the **target repo** (or `KroniK907/skills` for meta/skills work).

| Artifact | Label |
|----------|--------|
| Map | `wayfinder:map` |
| Decision log | `wayfinder:decision-log` |
| To Do ticket | `wayfinder:todo` + type + mode |

**Local fallback:** `wayfinder/plans/{FeatureName}.Map.md` and `{FeatureName}.Decision-Log.md` when GitHub is unavailable. Same body template; To Do/Completed as markdown tables.

**Create order:** Map → decision log issue → To Do tickets → wire blockers.

---

## Ecosystem integration

Skills that **read** wayfinder maps:

| Skill | Reads | Writes |
|-------|-------|--------|
| `write-a-prd` | Map Completed + decision log | PRD issue |
| `prd-to-issues` | PRD | `agent-queue` issues |

Skills that **write** wayfinder state (planned updates — see bootstrap map):

| Skill | Writes |
|-------|--------|
| `grill-me` | Decision log `{MAP-SLUG}-GM-xx`; resolution comment on grilling ticket |
| `wayfinder` | Map To Do / Completed / fog / Subfeatures |
| `research` (TBD) | Research ticket comment; optional map note |
| Cloud AFK automation (TBD) | PR + comment on AFK ticket; human closes |

**Handoff chain:** wayfinder (To Do empty) → optional integration `grill-me` → `write-a-prd` → `prd-to-issues` → implement → reconcile map (future).

---

## Ideation interview capture template

Use during **Chart** before tickets exist:

```markdown
## Ideation capture — {FeatureName}

| Zone | Known | Unknown / needs research |
|------|-------|---------------------------|
| Surfaces & experience | … | … |
| Behavior & correctness | … | … |
| Boundaries & integration | … | … |
| Persistence & data | … | … |
| Change, risk & evidence | … | … |
```

Each **unknown / needs research** cell item → candidate To Do ticket (when the question is sharp) or **Not yet specified** (when not).
