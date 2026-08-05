---
name: code-review
description: Two-axis code review (Standards + Spec) of changes since a pinned ref — commit, branch, tag, or merge-base. Runs parallel generalPurpose sub-agents and aggregates under separate headings. Use when the user wants to review a branch, a PR, work-in-progress changes, review since a ref, or asks "what changed since X".
---

# Code review

Two-axis review of the diff between `HEAD` and a **fixed point** the user supplies:

- **Standards** — does the change conform to this repo's documented coding standards (plus the Fowler smell baseline where the repo is silent)?
- **Spec** — does the change faithfully implement the originating issue / spec?

Both axes run as **parallel sub-agents** (Task tool, `generalPurpose`) so they do not pollute each other's context. This skill aggregates their findings under separate headings — do **not** merge or rerank across axes.

Detail: [REFERENCE.md](REFERENCE.md) — smell baseline, spec lookup order, sub-agent prompt templates, output format.

## Not this skill

| Skill | When instead |
|-------|----------------|
| `review-bugbot` | Fixed-prompt diff review for bugs — no spec lookup, no standards doc sweep |
| `review-security` | Fixed-prompt diff review for security — no spec lookup |
| [implement-task](../wayfinder/implement-task/SKILL.md) | Pick up approved implementation tasks — not ad-hoc review |

## Process

### 1. Pin the fixed point

Whatever the user said — commit SHA, branch name, tag, `main`, `HEAD~5`, etc. If unspecified, ask for it.

Capture once:

```powershell
git diff <fixed-point>...HEAD
git log <fixed-point>..HEAD --oneline
```

Use three-dot diff (merge-base comparison). Gate before sub-agents: `git rev-parse <fixed-point>` must succeed and the diff must be non-empty.

### 2. Identify the spec source

Follow [spec lookup order](REFERENCE.md#spec-lookup-order). Primary path: issue refs in commit messages → `gh issue view`. No Matt Pocock tracker dependency.

If no spec is found after the lookup chain, ask the user. If they confirm there is none, skip the Spec sub-agent and note **no spec available** in the final report.

### 3. Identify the standards sources

Collect repo docs that define how code should be written — e.g. `CONTRIBUTING.md`, `CODING_STANDARDS.md`, `.cursor/rules/`, [wayfinder REFERENCE](../wayfinder/REFERENCE.md) conventions for skills work.

Always include the [Fowler smell baseline](REFERENCE.md#fowler-smell-baseline) from REFERENCE — repo documented standards override baseline smells.

### 4. Spawn both sub-agents in parallel

Send **one message** with two Task tool calls (`subagent_type: generalPurpose`). Use prompt templates from [REFERENCE.md](REFERENCE.md#sub-agent-prompts).

- **Standards** — diff command, commit list, standards file list, smell baseline pasted in full
- **Spec** — diff command, commit list, fetched spec (issue body, file path, or user paste)

If spec is missing (step 2), skip Spec sub-agent only.

### 5. Aggregate

Present reports under `## Standards` and `## Spec` headings — verbatim or lightly cleaned. Do **not** merge or rerank findings across axes.

End with a one-line summary: total findings per axis and the worst issue **within each axis** (if any). Do not pick a single winner across axes.

See [output format](REFERENCE.md#output-format) and [_Why two axes_](REFERENCE.md#why-two-axes).
