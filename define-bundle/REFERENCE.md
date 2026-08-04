# Define bundle reference

## Bundle issue template

Use for GitHub issue body. Title: `Bundle: {short name}`.

```markdown
# Bundle: {short name}

**Status:** draft | approved

## Map

Parent: [{FeatureName}:Map](map-issue-url) · slug `{MAP-SLUG}` · decision log [#N](log-url)

## Name

{One line — build-oriented name}

## Scope summary

{What this bundle delivers; 1–3 paragraphs. List concrete deliverables when meta/infra.}

## Boundaries

**In scope**
- …

**Out of scope**
- …

**Dependencies**
- …

## Decisions

*(Covered GM rows — binding prose verbatim from decision log)*

**{MAP-SLUG}-GM-NNN** — …

## Constraints

*(Auto-included `[global]` rows — not claimed by this bundle)*

**{MAP-SLUG}-GM-NNN** — …

## Open questions

1. …

## User stories or Outcomes

N/A — meta/infra

- …
```

For product bundles with user stories, replace the `N/A` line with story bullets and keep **Outcomes** as acceptance criteria.

---

## Approval phrases

| User says | Agent may |
|-----------|-----------|
| **bundle approved** | Set bundle Status `approved`; update map Decision coverage (`scoped` + link); append log suffixes; optional Notes line |
| (edits requested) | Update draft bundle body in place; keep Status `draft` |
| (no approval) | Narrate or post draft only; **do not** write coverage or suffixes |

Synonyms accepted if unambiguous: "approve the bundle", "approve bundle #N".

**Separate from Reconcile:** `bundle approved` is owned by **define-bundle**, not wayfinder Reconcile. Reconcile still owns grilling ticket close + new GM row append.

---

## Global vs bundle-scoped rows

Per WF-ECO-GM-021:

| Signal | Treatment in bundle |
|--------|---------------------|
| `[global]` in log row text | **Constraints** only |
| Decision coverage status **`global`** | **Constraints** only |
| Coverage **`open`**, no suffix | May be **claimed** in **Decisions** |
| `- bundled via [#N]` suffix | Already claimed — exclude |
| Coverage **`scoped`** / **`assigned`** / **`implemented`** | Exclude from new bundles |

On **`bundle approved`**, only rows listed in **Decisions** receive the suffix and **`scoped`** coverage update.

Reconcile proposes **`[global]`** vs bundle-scoped when appending new rows; human confirms. Default **`[global]` when unsure**.

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
| `open` | Decided, not yet bundled | — |
| `scoped` | In approved bundle | Bundle issue |
| `assigned` | Implementation task exists | Task issue |
| `implemented` | Shipped | Task issue (closed) |
| `global` | Infrastructure; never bundled | — |

---

## Log suffix format

Append to the **end** of each covered row paragraph (after any `(from …)` source link):

```markdown
 — bundled via [#16](https://github.com/org/repo/issues/16)
```

Do not alter binding prose before the suffix.

---

## GitHub operations

Use **`gh` only** — no GraphQL for issue create/edit/link.

```powershell
gh issue create --repo owner/repo --title "Bundle: …" --label "wayfinder:bundle" --body-file path\to\body.md
gh issue edit <num> --body-file path\to\body.md
gh issue comment <num> --body "…"
```

Link parent map in the bundle **Map** section (markdown link). Native parent/child sub-issues are optional.

---

## Route heuristics (for wayfinder)

Suggest **define-bundle** when:

- User explicitly asks to bundle or implement from the decision log
- **Decision coverage** has a cluster of **`open`** rows clearly describing one deliverable (e.g. GM-013–018)
- User wants to ship while planning **To Do** or **Not yet specified** remain non-empty

Prefer planning frontier skills (grill-me, research ticket, etc.) when rows are still **`open`** because work is incomplete — not because fog exists elsewhere on the map.

After approval, suggest **create-tasks** (or direct implementation if create-tasks is unavailable).
