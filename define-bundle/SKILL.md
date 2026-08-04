---
name: define-bundle
description: Coalesce map decision-log clusters into draft build bundle issues and promote them on bundle approved. Use when a GM cluster is ready to implement while planning To Do or fog remain open, when wayfinder Route suggests bundling, or when the user mentions define-bundle or build bundles.
---

# Define bundle

Coalesce **`open`** decision-log rows into a **draft bundle issue** (`wayfinder:bundle`), then on human **`bundle approved`** promote it and sync map **Decision coverage** + log suffixes. Maps implement incrementally — do **not** wait for empty To Do, a PRD, or cleared fog.

Runs on a wayfinder map + its decision log. After approval, hand off to [create-tasks](../create-tasks/SKILL.md) or implement directly from the bundle when splitting is unnecessary.

## Not this skill

| Skill | When instead |
|-------|----------------|
| [wayfinder](../wayfinder/SKILL.md) | Chart, Materialize, Reconcile, Route only |
| [grill-me](../grill-me/SKILL.md) | Resolve unknowns → new binding GM rows |
| [create-tasks](../create-tasks/SKILL.md) | Split an **approved** bundle into implementation tasks |
| [write-a-prd](../write-a-prd/SKILL.md) | Small map-free scope only |

## Prerequisites

- Map issue (`wayfinder:map`) with **Decision coverage** table and linked decision log
- `gh` authenticated on the target repo
- Label `wayfinder:bundle` available

## Workflow

### 1. Load context

```text
gh issue view <map-num> --json body,title,url
gh issue view <log-num> --json body,title,url
```

From the map: slug, decision log link, **Decision coverage**, **To Do**, **Implementing**, **Notes**.

From the log body: all `{MAP-SLUG}-GM-*` rows. Parse scope:

- **`[global]`** in row text, or coverage status **`global`** → Constraints only (never claimed)
- **`- bundled via [#N]`** suffix → already scoped to a bundle
- Coverage **`open`** + no global tag → bundle candidates

### 2. Propose cluster

Present to the user:

- **Bundle name** (short, build-oriented)
- **Covered GM IDs** (contiguous cluster, one build slice)
- **Rationale** — why these rows ship together
- **Excluded rows** — globals, already bundled, or still foggy

Wait for confirmation or edits before creating/updating the bundle issue.

**Cluster heuristics:** same subsystem or skill (e.g. GM-013–018 → research skill); one vertical slice; rows that share deliverables. Do not bundle rows already **`scoped`**, **`assigned`**, or **`implemented`**.

### 3. Draft bundle issue

Create early with `gh issue create` or update an existing draft in place.

| Field | Value |
|-------|--------|
| Title | `Bundle: {short name}` |
| Label | `wayfinder:bundle` |
| Body | Per [REFERENCE.md](REFERENCE.md#bundle-issue-template) |
| **Status** | `draft` |

**Map section:** link parent map, slug, decision log — body link only (no GraphQL sub-issues).

**Decisions:** copy covered GM rows **verbatim** from the log body.

**Constraints:** copy every **`[global]`** log row verbatim — auto-included, not claimed.

Fill **Scope summary**, **Boundaries**, **Open questions**, **User stories or Outcomes** (`N/A — meta/infra` + bullet outcomes when no user stories).

Post or narrate the draft; end with: *Review the bundle — reply **bundle approved** when scope is accepted, or request edits.*

**Default:** do not run **`bundle approved`** writes without the explicit phrase.

### 4. On `bundle approved`

When the user says **`bundle approved`** (optionally naming the bundle issue):

1. **Bundle issue** — set **Status:** `approved` in body (`gh issue edit`)
2. **Map Decision coverage** — covered rows → **`scoped`**, **Linked issue** → bundle URL
3. **Decision log body** — append ` — bundled via [#N](url)` to each covered row (immutable paragraph text before suffix)
4. **Map Notes** — one-line approved-bundle link if helpful
5. **Comment** on bundle issue summarizing executed updates

Use `gh issue edit --body-file` for full body replacements. Requires `gh` auth.

**Do not** move rows to **Implementing** — that is [create-tasks](../create-tasks/SKILL.md).

### 5. Hand off

Tell the user:

- **Next:** [create-tasks](../create-tasks/SKILL.md) with the approved bundle link, **or** implement directly from the bundle when a single session needs no task split
- Planning **To Do** may stay open (per WF-ECO-GM-019)

## Interaction rules

1. **One bundle per build slice** — bundle-scoped GM rows are one-bundle-per-row (suffix on approval)
2. **Globals are inherited** — list in Constraints; never suffix or mark **`scoped`**
3. **Log body is authoritative** — binding prose lives in the log issue body only; comments are not binding
4. **Draft early** — create the bundle issue while scope is still being refined; update in place
5. **No implementation** — this skill bundles decisions; it does not write product code unless the bundle itself is meta/infra (then follow bundle deliverables)

## Quick start

User: "Bundle GM-013–018 for the research skill on map #12."

Load map #12 + log #11 → propose cluster → create draft `wayfinder:bundle` issue → user says **`bundle approved`** → sync coverage + suffixes → suggest create-tasks or direct implementation.

See [REFERENCE.md](REFERENCE.md) for bundle template, approval phrases, and global-row rules.
