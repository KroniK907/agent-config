---
name: commit
description: Stages and commits working-tree changes tied to an agent chat using git, with GitHub CLI for repo context and follow-ups. Splits into multiple logical commits when the chat produced distinct changes; uses a single commit for small cohesive edits. Use when the user wants to commit work from the current or a referenced chat, snapshot chat changes, or asks to git-commit agent output.
---

# Commit changes from an agent chat

## Goal

Turn edits from a conversation into one or more **clean commits**: correct staging, messages that match project style if obvious, and **logical splitting** when multiple independent changes landed in the same chat.

## Preconditions

- Shell can run `git` and `gh` from the relevant repository root.
- If commits will be pushed or opened as a PR, ensure `gh auth status` succeeds when that step is in scope.

## Workflow

### 1. Locate the repo and inventory changes

From the workspace (or path the user gives), run:

```bash
git rev-parse --show-toplevel
git status -sb
git diff
git diff --staged
```

Optional GitHub context (remote, default branch) without leaving the terminal:

```bash
gh repo view --json nameWithOwner,defaultBranchRef,url
```

### 2. Map changes to the chat

- Treat **paths with unstaged or staged edits** as the source of truth.
- Cross-check with the conversation: files mentioned, edits described, or tasks completed in that chat.
- If the user points at a transcript or session, use it only to recall *which* files and *what* themes changed; still verify with `git diff`.

### 3. Decide one commit vs several

**Prefer a single commit** when:

- The diff is small, or
- All edits serve one intent (one bugfix, one feature slice, one doc update).

**Split into multiple commits** when:

- Unrelated concerns are mixed (e.g. refactor + bugfix + formatting sweep),
- Changes would be hard to revert or cherry-pick as one unit,
- Different areas would benefit from different message scopes (`feat` vs `fix` vs `docs`).

Order commits **dependency-first** (lower-level or shared helpers before call sites).

Briefly tell the user the plan (e.g. "2 commits: fix X, then docs for Y") before running commands.

### 4. Stage and commit with git

For each commit:

1. Stage only the paths that belong to that logical change: `git add -- path1 path2` (avoid `git add .` unless the working tree for that commit is fully intentional).
2. Re-check: `git diff --staged`
3. Commit with a message that matches obvious repo conventions; if none, use a short imperative subject and optional body (Conventional Commits is a safe default).

```bash
git commit -m "type(scope): short summary" -m "Optional body with rationale."
```

Repeat until the working tree for chat-related work is clean (or until only intentionally excluded files remain).

### 5. Use `gh` after commits (when relevant)

`gh` does not replace `git` for staging/commits. Use it when the user wants GitHub next steps, for example:

- `gh pr create` after pushing a branch
- `gh issue view <n>` to reference an issue in the PR or commit message

Do not push or open a PR unless the user asks.

## Guardrails

- Do not commit secrets, `.env`, credentials, or large generated artifacts; if those appear in `git status`, stop and flag them.
- If unrelated local changes exist, commit only chat-related paths or ask the user whether to include them.
- Never force-push or rewrite history unless the user explicitly requests it.
