# Resolution comment templates

Post as a **new top-level comment** on the implementation task at end-of-run. Copy the matching template; fill every section.

**Success:** after push and build work complete. **Blocked:** after any startup gate failure - no repo edits on that run.

---

## Success template

```markdown
## Implementation resolution

**Task:** [#N](task-url) - **Bundle:** [#B](bundle-url) - **Method:** `{method-name}`

### Summary

<2-4 sentences - what shipped; branch state; anything left for human review.>

### Method

<{method-name}> - <one line on what the action skill did; link to key paths or artifacts.>

### Code review

<paste [code-review return artifact](../../../actions/code-review/REFERENCE.md#implement-task-return-artifact) - Auto-fixes applied, Standards/Spec remaining, Summary table. Use _(skipped - no diff)_ when Method produced no file changes.>

### Commits

- `<short-hash>` - <message>
- …

### Done when

| Done when | Status | Notes |
|-----------|--------|-------|
| <bullet from task issue> | done / partial / not done | <evidence - path, comment link, or gap> |

### Next

- Review the pull request **{pr-url}**; wayfinder **Reconcile** closes the task once you're happy with it
- <optional follow-up - blocked dependents, doc gaps, etc.>
```

After posting, set **Status:** `awaiting-reconcile` and add **`wf:needs-review`**; keep **`wf:approved`**.

---

## Blocked template

Use when startup gates fail (missing label, invalid Method, branch checkout failure, AFK serial lock, etc.). **No commits** on this run.

```markdown
## Implementation resolution - Blocked

**Task:** [#N](task-url) - **Bundle:** [#B](bundle-url)

### Summary

Run stopped at startup - <gate name>. No repository edits on this run.

### Blocked reason

<Exact gate that failed - e.g. missing `wf:approved`, invalid **## Method**, missing `integrationBranch`, worktree add failed, `wf:afk-running` held by #other.>

### Method

Not dispatched - startup did not pass.

### Commits

None this run.

### Done when

| Done when | Status | Notes |
|-----------|--------|-------|
| <each task bullet> | not done | Blocked at startup |

### Next

- <Human action - fix label, set Method, release lock, merge blocker, etc.>
- Re-run **implement-task** when gates pass
```

Blocked runs leave **Status** as it was.
