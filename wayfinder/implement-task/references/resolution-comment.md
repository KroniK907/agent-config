# Resolution comment templates

Post as a **new top-level comment** on the implementation task at end-of-run. Copy the matching template; fill every section.

**Success:** after push and build work complete. **Blocked:** after any startup gate failure — no repo edits on that run.

---

## Success template

```markdown
## Implementation resolution

**Task:** [#N](task-url) · **Bundle:** [#B](bundle-url) · **Method:** `{method-name}`

### Summary

<2–4 sentences — what shipped; branch state; anything left for human review.>

### Method

<{method-name}> — <one line on what the action skill did; link to key paths or artifacts.>

### Commits

- `<short-hash>` — <message>
- …

### Done when

| Done when | Status | Notes |
|-----------|--------|-------|
| <bullet from task issue> | done / partial / not done | <evidence — path, comment link, or gap> |

### Next

- Human: review this comment and diff on bundle branch **`{branch-name}`**
- Invoke wayfinder **Reconcile** with **`Approved — reconcile and close`** when accepted
- <optional follow-up — blocked dependents, doc gaps, etc.>

### Reconcile

<!-- Orchestrator does NOT post approval phrases. Human only. -->

Pending **`Approved — reconcile and close`** — agent must not close this task or remove **`wayfinder:approved`**.
```

After posting, set task body **Status:** `awaiting-reconcile`. Do **not** remove **`wayfinder:approved`**.

---

## Blocked template

Use when startup gates fail (missing label, invalid Method, branch checkout failure, AFK serial lock, etc.). **No commits** on this run.

```markdown
## Implementation resolution — Blocked

**Task:** [#N](task-url) · **Bundle:** [#B](bundle-url)

### Summary

Run stopped at startup — <gate name>. No repository edits on this run.

### Blocked reason

<Exact gate that failed — e.g. missing `wayfinder:approved`, invalid **## Method**, could not checkout `afk/bundle-23-…`, `wayfinder:afk-running` held by #other.>

### Method

Not dispatched — startup did not pass.

### Commits

None this run.

### Done when

| Done when | Status | Notes |
|-----------|--------|-------|
| <each task bullet> | not done | Blocked at startup |

### Next

- <Human action — fix label, set Method, release lock, merge blocker, etc.>
- Re-run **implement-task** when gates pass

### Reconcile

Not ready — resolve blocker first. Do not post **`Approved — reconcile and close`** for a blocked run.
```

Do **not** set **`awaiting-reconcile`** on blocked runs. Leave **Status:** `ready` (or prior state).
