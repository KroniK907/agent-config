---
name: commit
description: Stages and commits only working-tree changes attributable to the current agent chat (unless the user explicitly widens scope), using git and GitHub CLI for repo context and follow-ups. Ignores unrelated local edits from other sessions. Splits into multiple logical commits when this chat produced distinct changes; single commit for small cohesive edits. Use when the user wants to commit work from this chat, snapshot chat changes, or asks to git-commit agent output.
---

# Commit changes from an agent chat

## Goal

Turn edits from **this conversation** into one or more **clean commits**: correct staging, messages that match project style if obvious, and **logical splitting** when multiple independent changes landed in **this chat**.

**Default scope:** Only paths you can tie to **the current chat** (or a transcript/session the user explicitly points you at *for this run*). **Do not** stage or commit other dirty files just because `git status` lists them - those may be from another chat, manual work, or another branch of work unless the user clearly says otherwise.

## Preconditions

- Shell can run `git` and `gh` from the relevant repository root.
- If commits will be pushed or opened as a PR, ensure `gh auth status` succeeds when that step is in scope.

## Workflow

### 1. Locate the repo and inventory changes

Use `git status` / `git diff` to see the full working tree, but treat that listing as **inventory only**. You will stage **only** the subset that belongs to this chat (next section) - not “everything changed.”

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

### 2. Map changes to the chat (mandatory filter)

1. **Attribution:** Decide which paths were **actually produced or intentionally modified in this chat** (or in a **user-supplied** transcript/session the user asked you to use for this commit). Use the conversation: files touched, tasks completed, explicit user requests, tool edits traceable to this thread.
2. **`git diff` confirms content**, not inclusion: a file appearing in `git diff` does **not** mean it belongs in this commit if this chat never discussed or changed it.
3. **Exclude by default:** Leave untouched in the working tree (unstaged) any changed or untracked paths **not** attributable to this chat’s scope - **even if** the user would like a “clean” status. They must **explicitly** ask to commit “everything,” “all local changes,” named paths, or work from another context before you widen what you `git add`.
4. If it is **unclear** whether a path belongs to this chat, **ask one short question** or commit only the obvious subset and mention what was left out.
5. If the user points at a transcript or session, use it only to recall *which* files and *what* themes belong to *that* scope; still verify with `git diff`.

### 3. Decide one commit vs several

Apply this **only to paths you already attributed to this chat** in step 2 - not to the whole repo diff.

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

1. Stage only the paths that belong to that logical change **and** to this chat’s scope: `git add -- path1 path2` (avoid `git add .` unless every dirty path is intentionally in scope for this commit).
2. Re-check: `git diff --staged`
3. Commit with a message that matches obvious repo conventions; if none, use a short imperative subject and optional body (Conventional Commits is a safe default).

```bash
git commit -m "type(scope): short summary" -m "Optional body with rationale."
```

Repeat until **chat-attributed** work is committed. **Expect** other local changes to remain modified or untracked - that is normal when scope is chat-only.

### 5. Use `gh` after commits (when relevant)

`gh` does not replace `git` for staging/commits. Use it when the user wants GitHub next steps, for example:

- `gh pr create` after pushing a branch
- `gh issue view <n>` to reference an issue in the PR or commit message

Do not push or open a PR unless the user asks.

## Guardrails

- Do not commit secrets, `.env`, credentials, or large generated artifacts; if those appear in `git status`, stop and flag them.
- **Unless the user explicitly widens scope:** never commit paths that are not attributable to **this chat** (or the user-named transcript/session for this run). Unrelated dirty files are **out of scope**, not “optional extras.”
- Never force-push or rewrite history unless the user explicitly requests it.
