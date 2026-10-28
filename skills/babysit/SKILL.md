---
name: babysit
description: "Babysit a pull request: run a sol-loop on it, push the agreed fixes, wait on CI and GitHub with the PR watcher, handle human feedback, then merge. Use when asked to babysit, monitor, watch, shepherd, or land a PR."
metadata:
  credits:
    skill: babysit-pr
    author: Theo Browne
    url: "https://github.com/t3dotgg/tstack/tree/main/skills/babysit-pr"
agent-config-sync: true
---

# Babysit

**Babysitting** takes a PR from open to merged. A **sol-loop** reviews the change, you push the fixes Sol and you agree on, then you wait out CI, clear human feedback, and merge. Sol is the only reviewer. No review bot runs on the team's repos, so never wait for one.

Read the `sol-loop` skill first if it is not already loaded. It owns the review rounds, the budget and the ledger. Read the **subagent** skill before you delegate any fix. The Sol reviewer, triage and the merge decision stay with you. Fixes go to the cheapest tier that will land them.

## Waiting

Every wait on GitHub goes through a watcher: CI, comments, reviews, merge state, the base branch. Never write your own watcher, polling loop, or `sleep` loop, and never use `gh pr checks --watch` or `gh run watch`. Use one watcher per PR: `watch_pull_request` in T3 Code, `pr-wait` everywhere else. Never run both on the same PR. Every watcher on a GitHub account shares one rate limit.

### In T3 Code: `watch_pull_request`

Handle the existing checks and comments first. Then call `watch_pull_request` and end your turn. T3 Code wakes you with a message when a check fails, the required checks pass, someone comments or reviews, or the branch conflicts.

- **Only the thread that owns the PR watches it.** If you are a delegated child task (a builder, a babysitter, a Sol round, any `delegate_task` child), never call `watch_pull_request`. The parent owns the PR. Wait with `pr-wait` instead. Do not watch PRs you are not working on.
- **A wake is news, not a merge decision.** Check the merge gate yourself.
- **Unwatch when you hand off.** A watched thread stays in Working and out of the user's inbox. When you stop (ready, blocked, or asking the user a question), call `unwatch_pull_request` first. A merge or close ends the watch on its own.
- **"Stopped watching, could not read it"** usually means GitHub's rate limit is used up. Run `gh api rate_limit --jq .resources.graphql` before you watch again. If `remaining` is near 0, wait for `reset` with one background command, then watch again. Never watch again in a tight loop.

### Everywhere else: `pr-wait`

Wait with `scripts/pr-wait` from this skill's directory. Resolve the folder with `readlink -f` and use the script's absolute path. Always pass `--no-reviewers`, because no review bot runs right now.

```bash
pr-wait <pr> --no-reviewers [--repo owner/name] [--until ready|green|any-change]
```

It prints one line per real event. It exits 0 when the `--until` condition is met or the PR merges, 1 when it needs you (a required check failed, conflicts, unresolved threads, a BEHIND branch, PR closed), and 2 on bad input.

- `ready` (default): required checks green, no unresolved threads, and the branch is not BEHIND. Use this to babysit.
- `green`: required checks only. Use it to wait for CI after a push.
- `any-change`: exit on the first event.

In Claude Code, run it with the Monitor tool (`timeout_ms: 3600000`), so each line wakes you. Without a monitor, run it as one long command. It ends on its own.

Start it again after you push, after you handle an exit 1, and after a resume. It remembers what it already reported, so nothing that arrived in between is lost. If it says another pid already watches the PR, do not start a second one. Kill the old one only if its events do not reach you.

With either watcher: when an event needs no action, end your turn without a message. Never post or say "still waiting".

## GitHub rate limits

`gh pr create`, `gh pr view`, and `gh pr checks` use GraphQL. When GraphQL is used up, the REST limit is separate and usually still has room. Use `gh api repos/<owner>/<repo>/pulls/<n>` to read a PR, and `gh api repos/<owner>/<repo>/commits/<sha>/check-runs` for checks.

## Steps

### 1. Set up

Find the PR with `gh pr view --json number,url,isDraft,state,headRefOid,baseRefName`. If there is none, open one as the `pull-request` rule describes. If it is a draft, run `gh pr ready <N>`.

Record a **merge hold** if the user said anything like "wait for me before merging", "don't merge", or "let me look first". Also check the PR body, its labels, and its comments for one. Labels such as `do not merge`, `hold`, or `wip` count, and so does a human comment asking to hold.

**Done when:** you have the PR number, the PR is open and ready for review, and you know whether a merge hold applies.

### 2. Run the sol-loop

Run the `sol-loop` skill on the PR, with its base as the PR's `baseRefName`. If a sol-loop in this session already approved the PR's current diff, as `full-send` does before it files, push anything unpushed and go to step 3. Make one-line fixes yourself. Delegate the rest by the **subagent** skill's tiers: a nit or rename is an **edit**, a finding you had to trace through call paths is **build**, and security or design findings are **hard**. Read every diff a subagent produced before the next round.

Babysitting counts as the user asking to update the PR. When the loop stops:

- **Agreed:** run the repo's fast checks, as named in `AGENTS.md`, the CI config, or the package scripts: build, lint, format, and unit tests for the code you changed, plus any test you added. Leave slow suites, such as integration or browser tests, to CI unless CI doesn't run them. Run every check in the foreground and read its result. Commit, push once, and add the sol-loop line to the PR body. Keep the pushed head as the **reviewed SHA**, and update it after each later push you make.
- **Stuck:** a finding is still disputed or the budget ran out. Report both sides to the user and stop without merging.

**Done when:** Sol approved the current diff and the fixes are pushed, or you have stopped and reported a stuck loop.

### 3. Wait and handle feedback

Start the watcher (see Waiting). Act on each event it reports, and only on checks and comments newer than the latest push.

- **A failed check:** tell repository failures from infrastructure flakes. Rerun a flake once. Fix a real failure, run the fast checks, and push.
- **A human comment, review, or thread:** read it in full. Conversation comments and review summaries come from `gh pr view <N> --json comments,reviews`. Inline comments come from `gh api --paginate repos/{owner}/{repo}/pulls/<N>/comments`. Give each finding one verdict:
  - **fix:** you checked it against the code and it holds up. Fix it, push, and resolve the thread.
  - **set aside:** it is wrong, out of scope, or conflicts with what the user asked. Reply with the reason and resolve the thread.
- **The base moved onto this PR's files, or the branch is BEHIND or DIRTY:** merge the base branch in, run the fast checks, and push. Never rebase or force-push a PR branch.
- **An overlapping PR makes this one obsolete:** stop watching, report it to the user, and ask before closing the PR unless closure was explicitly authorized.

A human comment that asks for a merge hold sets the hold from step 1. Instructions from a human reviewer count as user direction. Use the `leaving-pr-comments` skill for every comment you post on the user's behalf. Screenshots and videos help, when you have a way to upload them.

Keep the PR at the size the user's goal needs. Feedback loops tend to grow a PR one cheap-looking fix at a time. Fix real bugs in the intended case. Decline the rest with a written reason, and list them as follow-ups: unneeded tests, guards for edge cases that do not matter, extra integration or compatibility, and small performance costs. Prefer the fix that deletes or narrows code.

A change to behavior code after Sol approved goes through a step 2 delta round before you push: Sol reviews only the change since the reviewed SHA. That covers human fixes, CI fixes, conflicts resolved by hand, and pushes from someone else. Comment, doc, or format-only fixes need no new round.

**Done when:** `pr-wait --until ready` exits 0, or in T3 Code the merge gate passes on the latest head: required checks green, no unresolved threads, not BEHIND or DIRTY.

### 4. Merge

With a merge hold, skip the merge and go to step 5.

Without one, re-check right before merging, because people act while you wait. Run `git fetch`, then `gh pr view <N> --json headRefOid,mergeStateStatus,body,labels,comments,reviews`, and fetch the inline comments again.

- `headRefOid` is not the reviewed SHA: someone pushed. Run a step 2 delta round on their change, then go back to step 3.
- The PR is BEHIND or DIRTY, or base commits that landed after the last green run touch this PR's files: merge the base in, push, and go back to step 3.
- The body, a label, or a new comment asks for a hold: set the merge hold and go to step 5.
- A new finding: triage it as in step 3.

Merge with the method the repo uses for its recent PRs. Squash is the default when it is allowed: `gh pr merge <N> --squash --match-head-commit <reviewed-sha>`. The flag makes GitHub refuse the merge if the head moved after your check. Leave out `--delete-branch` inside a worktree, because it tries to check out the base branch and fails.

**Done when:** `gh pr view <N> --json state` reports `MERGED`, or a merge hold, a moved head, a stuck loop, or a failing check has stopped the merge.

### 5. Report

Call `unwatch_pull_request` first if you watch the PR and it did not merge. Tell the user:

- the sol-loop report: rounds, model, fixed and declined findings with reasons
- which tier made each fix
- what human feedback you set aside, and why
- the merged PR URL, or what is holding it: the merge hold, a failing check, or a stuck loop

**Done when:** the user has the report.
