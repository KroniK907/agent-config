# Reconcile

Execution detail for wayfinder **Reconcile**: resolution template, inference rules, and the `wf` commands that apply it.

`wf` is `go run <wayfinder>/utilities/wf/wf.go`, where `<wayfinder>` is the installed skill root (`skills/wayfinder` here, or the applied copy under `.cursor/skills/`, `~/.cursor/skills/`, `.claude/skills/`, or `~/.claude/skills/`). Optional `-R owner/repo` targets another repo. Needs `gh` auth and Go.

## Resolution comment

Post on the session ticket. Omit sections that have nothing in them.

```markdown
## Resolution - {ticket title}

### Session summary

<What was settled, deferred, or left open - concrete terms, not zone labels.>

### Decision log

**{MAP-SLUG}-GM-NNN** - ... `[global]` when it applies map-wide; omit when bundle-scoped.
(from [{ticket title}](#ticket-num))

### Decision coverage

| GM ID | Status | Linked issue |
|-------|--------|--------------|
| {MAP-SLUG}-GM-NNN | open / global | - |

### Map updates

- **Completed gist:** ...
- **Not yet specified:** ...
- **Out of scope:** ...
- **Notes:** ...

### New ticket candidates

| Title | Type | Mode | Question | Blocked by |
|-------|------|------|----------|------------|
| **Grill:** ... / **Research:** ... | research / prototype / grilling / task | HITL / AFK | One sharp question | - or issue ref |

### Bundle cluster suggestions

| Suggested name | Covered GM IDs | Rationale | Excluded |
|----------------|----------------|-----------|----------|
| ... | {MAP-SLUG}-GM-012-015 | One vertical slice | globals, already bundled, still foggy |

### Ticket invalidations

- **Close:** [#N Title](link) - superseded by {MAP-SLUG}-GM-NNN
- **Retitle / retype / move to Out of scope:** ...

### Route hint

<One recommended next step + skill.>
```

Ticket titles follow [ticket title conventions](../REFERENCE.md#ticket-title-conventions). Research candidates get **Done when** bullets per the [research ticket template](../actions/research/REFERENCE.md#research-ticket-template). Bundle suggestions stay suggestions - [define-bundle](../actions/define-bundle/SKILL.md) creates bundles.

## Reconcile inference

Load what you need first: `wf section <map> "To Do" "Decision coverage"`, `wf log <log> --list` or `--ids`, and the session output.

### All session types

| Signal in session | Propose |
|-------------------|---------|
| Binding decision | Decision log row + coverage (`open` or `global`) |
| Explicit deferral ("decide later", "needs spike") | **Not yet specified**, often plus a ticket candidate |
| Supersedes an open To Do **Question** | Ticket invalidation |
| Rejected or mis-scoped direction | **Out of scope** or invalidation |
| Durable preference | **Notes** |

### Grilling (`grill-me`, `sous-vide`, `strategic-ideation`)

- Split distinct binding answers into separate GM rows.
- Deferred branches become fog plus a typed ticket: `research` for facts, `grilling` for decisions, `prototype` for layout or interface exploration.
- Several `open` GMs describing one deliverable → bundle cluster suggestion ([heuristics](../actions/define-bundle/REFERENCE.md#route-heuristics-for-wayfinder)).
- Cross-cutting constraints are `[global]` with coverage `global`. Default to `[global]` when unsure.

### Research

- Findings become GM rows only once the user treats them as decided - often via a follow-up grilling ticket.
- Merge the findings comment's **Proposed tracker updates** into candidates, fog, and Notes.
- **Done when** satisfied → close. Gaps → new `research` ticket or keep open.

### Prototype / task

- Binding decisions → GM rows; delivered asset → **Completed** gist.
- New questions → typed ticket candidates. Checklist incomplete → keep open.

### Constrain-fog

For a `Constrain:` ticket whose body has `## Fog resolution` with **Status:** `ready for reconcile`. Map its **New ticket candidates**, cleanup (out of scope / delete), remaining fog, invalidations, route hint, and session summary onto the matching resolution sections. Binding decision prose becomes **Grill:** / **Ideate:** candidates, not GM rows, unless the user asked for GM rows during the session. Do not use Materialize or append `## Map discovery` for constrain-fog output.

## Apply steps

Apply once the outcome is accepted (see the hub skill). Honor any edits the user made to the resolution comment, and skip candidates they struck.

1. **Decision log** - write the rows to a markdown file (one `**{MAP-SLUG}-GM-NNN** - ...` paragraph each) and post:

   ```text
   wf log-append <log-issue> rows.md
   ```

   It checks the prefix, that IDs are consecutive, and that none exist. A later entry with the same ID replaces the earlier one. Use `--amend` only to replace rows you just posted.

2. **Map** - one `map-edit` call; it validates, uploads, and re-fetches.

   ```text
   wf map-edit <map-issue> --complete N --gist "..." --coverage {MAP-SLUG}-GM-NNN=open --append "Not yet specified=..."
   ```

   `--complete N --gist TEXT` moves a **To Do** or **Implementing** row to **Completed**. `--coverage ID=status[=link]` updates coverage. `--append "Section=line"` / `--remove "Section=substring"` edit fog, Notes, Out of scope, or To Do. `--dry-run` prints without uploading.

3. **New tickets** - `gh issue create` with `wf:todo` + type + mode, then `wf map-edit --append "To Do=..."`; wire blocked-by in a second pass.

4. **Invalidations** - close, retitle, or move superseded tickets.

5. **Close** - when the ticket's work is done: `gh issue close <ticket>` and remove **`wf:needs-review`**. Keeping it open: remove the label only and say what's left.

Inspect with `wf section <map> "To Do"`, `wf log <log> --ids A,B` / `--global` / `--list` / `--next`. Edit through `wf` rather than re-downloading bodies and splicing by hand.

### Map and issue body edits (Reconcile)

Use a full body replacement only when the change isn't a complete, coverage, append, or remove:

1. `wf body get <issue> body.md`
2. Edit the file; keep section breaks; ASCII hyphen only.
3. `wf validate body.md` - stop if it fails (checks line count, section headings on their own lines, no mojibake or em/en dash or middle dot).
4. `wf body put <issue> body.md`

Don't round-trip a body through a PowerShell variable - it collapses newlines.

## Implementation tasks

Closing a shipped **`wf:approved`** task (from [create-tasks](../actions/create-tasks/SKILL.md) or [one-off](../orchestrators/one-off/SKILL.md)), once the user is satisfied with the PR:

- Close the task; remove **`wf:approved`** and **`wf:needs-review`**.
- `wf map-edit --complete` the **Implementing** row (one-off tickets: the **To Do** row).
- Set coverage `implemented` for bundle-scoped GMs this task fully shipped. One-off tickets touch coverage only when the body names GM rows.
