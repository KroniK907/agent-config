# Define bundle reference

## Bundle issue template

Use for GitHub issue body. Title: `Bundle: {short name}`.

```markdown
# Bundle: {short name}

**Status:** draft | approved

## Map

Parent: [{FeatureName}:Map](map-issue-url) - slug `{MAP-SLUG}` - decision log [#N](log-url)

## Name

{One line - build-oriented name}

## Scope summary

{What this bundle delivers; 1-3 paragraphs. List concrete deliverables when meta/infra.}

## Boundaries

**In scope**
- …

**Out of scope**
- …

**Dependencies**
- …

## Decisions

Decision log: [#N](log-url). IDs only. Load prose with `wf log <log> --ids`.

- **{MAP-SLUG}-GM-NNN** - one-line summary

## Constraints

All `[global]` rows in the decision log apply. Load them with `wf log <log> --global`.

- **{MAP-SLUG}-GM-NNN** - one-line summary

## Open questions

1. …

## User stories or Outcomes

N/A - meta/infra

- …
```

For product bundles with user stories, replace the `N/A` line with story bullets and keep **Outcomes** as acceptance criteria.

---

## Approval phrases

| User says | Agent may |
|-----------|-----------|
| **bundle approved** | Set bundle Status `approved`; remove **`wf:needs-review`**; update map Decision coverage (`scoped` + link); optional Notes line. Do not create a git branch |
| (edits requested) | Update draft bundle body in place; keep Status `draft`; keep **`wf:needs-review`** |
| (no approval) | Narrate or post draft only; add **`wf:needs-review`**; **do not** write coverage |

Synonyms accepted if unambiguous: "approve the bundle", "approve bundle #N".

**Separate from Reconcile:** `bundle approved` is owned by **define-bundle**, not wayfinder Reconcile. Reconcile still owns grilling ticket close + new GM row append.

---

## Global vs bundle-scoped rows

| Signal | Treatment in bundle |
|--------|---------------------|
| `[global]` in log row text | **Constraints** only |
| Decision coverage status **`global`** | **Constraints** only |
| Coverage **`open`** | May be **claimed** in **Decisions** |
| Coverage **`scoped`** / **`assigned`** / **`implemented`** | Already claimed - exclude |

On **`bundle approved`**, rows listed in **Decisions** get coverage **`scoped`** and a link to the bundle. Do not edit the decision log.

Reconcile proposes **`[global]`** vs bundle-scoped when appending new rows; human confirms. Default **`[global]` when unsure**. Reconcile may also propose **bundle cluster suggestions** in the resolution comment - human runs define-bundle to draft/approve bundle issues.

---

## Decision coverage updates

On **`bundle approved`**, for each covered GM ID in **Decisions**:

```markdown
| {MAP-SLUG}-GM-NNN | scoped | [#bundle](bundle-url) |
```

Do not change status for **Constraints** rows.

Later lifecycle (create-tasks / implementation Reconcile):

| Status | Meaning | Linked issue |
|--------|---------|--------------|
| `open` | Decided, not yet bundled | - |
| `scoped` | In approved bundle | Bundle issue |
| `assigned` | Implementation task exists | Task issue |
| `implemented` | Shipped | Task issue (closed) |
| `global` | Infrastructure; never bundled | - |

---

## Where tasks land

define-bundle does not create a branch. Each task that [implement-task](../../orchestrators/implement-task/SKILL.md) picks up gets a worktree and a pull request. The base branch is `integrationBranch` in the app repo `.cursor/agent-manifest.json`. [create-tasks](../create-tasks/SKILL.md) does not create branches either.

---

## Route heuristics (for wayfinder)

Suggest **define-bundle** when:

- User explicitly asks to bundle or implement from the decision log
- **Decision coverage** has a cluster of **`open`** rows clearly describing one deliverable
- User wants to ship while planning **To Do** or **Not yet specified** remain non-empty

Prefer planning frontier skills (grill-me, research ticket, etc.) when rows are still **`open`** because work is incomplete - not because fog exists elsewhere on the map.

After approval, suggest **[design-modules](design-modules/SKILL.md)** when module interface shape is still open (one or more modules), then **create-tasks** (or direct implementation if create-tasks is unavailable). Narrate that each task will open its own pull request through implement-task.
