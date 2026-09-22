# Define bundle reference

## Bundle issue template

Use for GitHub issue body. Title: `Bundle: {short name}`.

```markdown
# Bundle: {short name}

**Status:** draft | approved

**Branch:** afk/bundle-{issue-num}-{slug}   *(draft: proposed; approved: confirmed - see [Bundle branch](#bundle-branch-wf-eco-gm-027))*

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
| **bundle approved** | Set bundle Status `approved`; remove **`wf:needs-review`**; create and push bundle git branch; persist **Branch:** on bundle body; update map Decision coverage (`scoped` + link); optional Notes line |
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

## Bundle branch (WF-ECO-GM-027)

All implementation tasks for a bundle commit and **push** on one shared branch. **Agents never open PRs** - humans open a PR when the bundle is complete.

### Naming pattern

```text
afk/bundle-{issue-num}-{slug}
```

| Part | Rule |
|------|------|
| `{issue-num}` | Bundle issue number (e.g. `23` for `#23`) |
| `{slug}` | Short kebab-case from bundle **Name** (e.g. `afk-v1` from "AFK v1 - skills…") |

Example: bundle [#23](https://github.com/KroniK907/skills/issues/23) → **`afk/bundle-23-afk-v1`**

### Draft bundle

When creating or updating a draft bundle issue, **propose the full `Branch:` line** in the body (after **Status**). Use the naming pattern above; human may edit before approval.

```markdown
**Status:** draft

**Branch:** afk/bundle-23-afk-v1
```

Do **not** create the git branch while Status is `draft`.

### On `bundle approved`

After setting Status `approved` and before handoff to create-tasks:

1. **Confirm branch name** - use the draft **Branch:** line; human may edit the slug in chat before approval (e.g. "bundle approved with branch `afk/bundle-23-my-name`")
2. **Create branch** - from default branch (usually `main`):

```powershell
git fetch origin
git checkout main
git pull origin main
git checkout -b <branch-name>
git push -u origin <branch-name>
```

3. **Persist** - set **Branch:** on the bundle body to the confirmed name (`gh issue edit`)

[implement-task](../../orchestrators/implement-task/SKILL.md) reads **Branch:** from the approved bundle at pickup; [create-tasks](../create-tasks/SKILL.md) does not create branches.

### Human edit gate

The **Branch:** line is the human edit gate at **`bundle approved`**. Agent proposes in draft; human confirms or renames when approving. Once approved, branch name is stable for the bundle lifecycle unless human explicitly edits the bundle body.

---

## Route heuristics (for wayfinder)

Suggest **define-bundle** when:

- User explicitly asks to bundle or implement from the decision log
- **Decision coverage** has a cluster of **`open`** rows clearly describing one deliverable
- User wants to ship while planning **To Do** or **Not yet specified** remain non-empty

Prefer planning frontier skills (grill-me, research ticket, etc.) when rows are still **`open`** because work is incomplete - not because fog exists elsewhere on the map.

After approval, suggest **[design-modules](design-modules/SKILL.md)** when module interface shape is still open (one or more modules), then **create-tasks** (or direct implementation if create-tasks is unavailable). Narrate that the bundle **Branch:** is set and pushed - create-tasks and implement-task use it for all bundle task commits.
