# Reconcile

Execution detail for wayfinder **Reconcile**. The hub skill states when to run it. This file is the resolution template, inference rules, approval phrases, and the `wf` commands that apply an approval.

`wf` is `go run <wayfinder>/utilities/wf/wf.go`. `<wayfinder>` is the installed skill root (`skills/wayfinder` in this repo, or `.cursor/skills/wayfinder` / `~/.cursor/skills/wayfinder` after apply). Optional `-R owner/repo` targets another repo. Needs `gh` auth and Go.

## Reconcile resolution template

Post as a **comment on the session ticket** (grilling, research, prototype, task). Non-binding until human approval.

**Section order (fixed):**

```markdown
## Resolution - {ticket title}

### Session summary

<What was settled, deferred, or left open - concrete terms, not zone labels alone.>

### Decision log (proposed append)

**{MAP-SLUG}-GM-NNN** - ... `[global]` when infrastructure applies map-wide; omit tag when bundle-scoped.
(from [{ticket title}](#ticket-num))

### Decision coverage (proposed)

| GM ID | Status | Linked issue |
|-------|--------|--------------|
| {MAP-SLUG}-GM-NNN | open / global | - |

### Map updates (proposed)

- **Completed gist:** ...
- **Not yet specified:** ...
- **Out of scope:** ...
- **Notes:** ...

### New ticket candidates

| Title | Type | Mode | Question | Blocked by |
|-------|------|------|----------|------------|
| **Grill:** ... / **Research:** ... / etc. | research / prototype / grilling / task | HITL / AFK | One sharp Question | - or issue ref |

Title prefix must match Type - see [ticket title conventions](../REFERENCE.md#ticket-title-conventions).

Omit table when none. Derive **Done when** bullets for `research` rows at materialize time (see [feature-discovery - research ticket shape](../ideation/feature-discovery/REFERENCE.md#research-ticket-shape-materialize)).

### Bundle cluster suggestions

<!-- Non-binding - define-bundle owns draft/approved bundle issues. Reconcile suggests only. -->

| Suggested name | Covered GM IDs | Rationale | Excluded |
|----------------|----------------|-----------|----------|
| ... | {MAP-SLUG}-GM-012-015 | One vertical slice / subsystem | globals, already bundled, still foggy |

Omit table when no cluster is ready. Rows stay **`open`** until **`bundle approved`**.

### Ticket invalidations

- **Close:** [#N Title](link) - superseded by {MAP-SLUG}-GM-NNN / merged into this session
- **Retitle / retype:** ...
- **Move to Out of scope:** ...

Omit when none.

### Route hint

<One recommended next step + skill - e.g. frontier ticket, define-bundle on cluster above, create-tasks on approved bundle.>

---

Ready for review - reply **Approved - reconcile and close** (or **Approved - reconcile, keep open**) when accepted. Edit any section in this comment before approving.
```

**Draft-only:** Do not append the decision log, edit the map, create issues, or close tickets until an approval phrase.

## Reconcile inference

Run after loading the map sections you need (`wf section`), the decision log rows you need (`wf log --list` or `--ids`), **Decision coverage**, **To Do**, and the sibling session output.

### All session types

| Signal in session | Propose in resolution |
|-------------------|----------------------|
| Binding decision with durable prose | Decision log row + **Decision coverage** (`open` or `global`) |
| Explicit deferral ("decide later", "needs spike") | **Not yet specified** + often a **New ticket candidate** |
| Supersedes an open To Do **Question** | **Ticket invalidations** |
| Mis-scoped or rejected direction | **Out of scope** or invalidation |
| Durable preference for map **Notes** | **Map updates -> Notes** |
| Next obvious frontier unchanged | **Route hint** -> existing unblocked To Do |

### Grilling (`grill-me`, `strategic-ideation`)

Primary source for full-session inference.

| Signal | Propose |
|--------|---------|
| Branch marked `complete` with binding answers | GM row(s); split distinct decisions into separate IDs |
| Branch deferred with named follow-up | Fog item + typed ticket candidate (`research` for facts; `grilling` for decisions; `prototype` for layout/interface exploration) |
| **Surfaces & experience** - layout agreed but build order unclear | `prototype` or `grilling` ticket; optional bundle cluster if GM cluster is build-ready |
| Multiple `open` GMs describing one deliverable | **Bundle cluster suggestion** (see [define-bundle cluster heuristics](../actions/define-bundle/REFERENCE.md#route-heuristics-for-wayfinder)) |
| Infrastructure / cross-cutting constraint | `[global]` on log row; coverage **`global`** |
| Uncertain scope tag | Default **`[global]` when unsure**; human edits in review |

Do **not** create `wf:bundle` issues in Reconcile. Narrate clusters for [define-bundle](../actions/define-bundle/SKILL.md).

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

1. Create child issues with `gh issue create`. Labels: `wf:todo` + type + mode (same as [Materialize](../REFERENCE.md#materialize-from-map-discovery)).
2. Append **To Do** rows with `wf map-edit --append`; wire **blocked-by** in a second pass.
3. For `research`, include **Done when** in the issue body per [research ticket template](../actions/research/REFERENCE.md#research-ticket-template).

Skip rows the human struck from the resolution comment before approving.

### Constrain-fog

Run when loading a **`Constrain:`** ticket whose body contains **`## Fog resolution`** with **Status:** `ready for reconcile`.

| Signal in artifact | Propose in resolution |
|--------------------|----------------------|
| **New ticket candidates** rows | **New ticket candidates** table (materialize on approval - same labels as Materialize) |
| **Cleanup** - `out of scope` | **Map updates -> Out of scope** |
| **Cleanup** - `deleted` or **Per-item** - remove | Drop line from **Not yet specified** |
| **Remaining fog** / rewrite / split | **Map updates -> Not yet specified** |
| **Ticket invalidations** | **Ticket invalidations** section |
| **Route hint** in artifact | **Route hint** (merge with standard inference) |
| **Session summary** | **Session summary** |
| Binding decision prose | **Grill / Ideate** ticket candidates only - **not** GM rows unless the user explicitly requested during constrain-fog |

Do **not** use **Materialize** for constrain-fog output. Do **not** append **`## Map discovery`** to the map issue.

## Completed workflow and approval

1. Sibling skill (or Reconcile draft) posts a **resolution comment** on the ticket ([template above](#reconcile-resolution-template)).
2. Human reviews, edits the comment if needed, and sends an **approval phrase** (below).
3. Agent **Reconcile** runs the approved GitHub and map updates with `wf` (including ticket materialize when candidates were not removed).
4. If the resolution invalidates other tickets, update or close those and move mis-scoped items to **Out of scope**.

### Approval phrases

| User says | Agent may |
|-----------|-----------|
| **Approved - reconcile and close** | `wf log-append` the approved rows; `wf map-edit` to complete the row, update **Decision coverage**, and update fog/Notes/Out of scope; `gh issue close` the ticket and remove **`wf:needs-review`**; `gh issue create` for approved **New ticket candidates**; apply **Ticket invalidations** |
| **Approved - reconcile, keep open** | `wf log-append`; `wf map-edit` for coverage, fog, and Notes; optional `gh issue create` / invalidations; remove **`wf:needs-review`**. Do not close the source ticket or move it to **Completed** |
| (no approval) | Post the draft resolution comment only. Do not close, edit the map, append the log, or create issues |

**Not Reconcile:** **`bundle approved`** is [define-bundle](../actions/define-bundle/SKILL.md). **`scope approved`** / **`tasks approved`** are [create-tasks](../actions/create-tasks/SKILL.md).

Synonyms are fine when unambiguous: "approve and close #N", "reconcile and close this ticket".

**Requires:** `gh` on the target repo, and Go so `wf` can run.

### Approved reconcile steps

Run these in order. Do not re-download the decision log or the map to splice a few lines by hand.

1. **Decision log** - write the approved rows to a small markdown file (one `**{MAP-SLUG}-GM-NNN** - ...` paragraph per row). Post them as a comment:

   ```text
   go run <wayfinder>/utilities/wf/wf.go log-append <log-issue> rows.md
   ```

   `log-append` checks the ID prefix, that IDs are consecutive, and that no ID already exists. Binding text is the log body plus later comments. A later entry with the same ID replaces the earlier one. Use `--amend` only when replacing rows you just posted in this session.

2. **Map** - one `map-edit` call. It validates, uploads, and re-fetches.

   ```text
   go run <wayfinder>/utilities/wf/wf.go map-edit <map-issue> --complete N --gist "..." --coverage {MAP-SLUG}-GM-NNN=open --append "Not yet specified=..." 
   ```

   Flags: `--complete N --gist TEXT` moves a **To Do** or **Implementing** row to **Completed**. `--coverage ID=status[=link]` updates **Decision coverage**. `--append "Section=line"` adds a fog, Notes, or Out of scope line. `--remove "Section=substring"` drops a matching line. `--dry-run` prints the edit and does not upload.

3. **Close** - on full approval, `gh issue close <ticket>` and `gh issue edit <ticket> --remove-label wf:needs-review`. On keep-open, remove the label only.

4. **New tickets** - `gh issue create` for each approved candidate, then `wf map-edit --append "To Do=..."` for the new rows.

Read a section before editing it with `wf section <map-issue> "To Do" "Decision coverage"`. Read decision text with `wf log <log-issue> --ids A,B` or `--global`. `--list` prints IDs and one-line summaries. `--next` prints the next free ID.

### Map and issue body edits (Reconcile)

Prefer `wf map-edit` and `wf log-append`. Use a full body replacement only when the edit is not a complete, coverage, append, or remove.

1. **Fetch** - `wf body get <issue> body.md`
2. **Edit the file** - keep section breaks. ASCII hyphen only.
3. **Validate** - `wf validate body.md`. Stop if it fails.
4. **Upload** - `wf body put <issue> body.md`

`wf validate` checks line count (at least 40), that `## To Do`, `## Completed`, and `## Decision coverage` each sit on their own line, and that the file has no mojibake or em dash, en dash, or middle dot.

Do not round-trip a body through a PowerShell variable. That collapses newlines onto one line.

### Implementation task Reconcile

When closing a **`wf:approved`** implementation task (from [create-tasks](../actions/create-tasks/SKILL.md)):

| User says | Agent may |
|-----------|-----------|
| **Approved - reconcile and close** | Close the task; remove **`wf:approved`** and **`wf:needs-review`**; `wf map-edit --complete` the **Implementing** row; set Decision coverage **`implemented`** for bundle-scoped GMs this task shipped |

Synonyms are fine when unambiguous: "reconcile and close task #N".

See [create-tasks REFERENCE - implementation Reconcile](../actions/create-tasks/REFERENCE.md#implementation-reconcile).
