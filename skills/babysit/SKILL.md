---
name: babysit
description: "Babysit a pull request: run a sol-loop on it, push the agreed fixes, handle human feedback and CI, then merge. Use when asked to babysit, shepherd, or land a PR."
agent-config-sync: true
---

**Babysitting** takes a PR from open to merged. A **sol-loop** reviews the change, you push the fixes Sol and you agree on, then you clear human feedback and CI and merge. Read the `sol-loop` skill first if it is not already loaded. It owns the review rounds, the budget and the ledger.

Read the **subagent** skill before you delegate any fix. The Sol reviewer, triage and the merge decision stay with you. Fixes go to the cheapest tier that will land them.

## Steps

### 1. Set up

Find the PR with `gh pr view --json number,url,isDraft,state,headRefOid,baseRefName`. If there is none, open one as the `pull-request` rule describes. If it is a draft, run `gh pr ready <N>`.

Record a **merge hold** if the user said anything like "wait for me before merging", "don't merge", or "let me look first". Also check the PR body, its labels, and its comments for one. Labels such as `do not merge`, `hold`, or `wip` count, and so does a human comment asking to hold.

**Done when:** you have the PR number, the PR is open and ready for review, and you know whether a merge hold applies.

### 2. Run the sol-loop

Run the `sol-loop` skill on the PR, with its base as the PR's `baseRefName`. Make one-line fixes yourself. Delegate the rest by the **subagent** skill's tiers: a nit or rename is an **edit**, a finding you had to trace through call paths is **build**, and security or design findings are **hard**. Read every diff a subagent produced before the next round.

Babysitting counts as the user asking to update the PR. When the loop stops:

- **Agreed:** run the repo's fast checks, as named in `AGENTS.md`, the CI config, or the package scripts: build, lint, format, and unit tests for the code you changed, plus any test you added. Leave slow suites, such as integration or browser tests, to CI unless CI doesn't run them. Run every check in the foreground and read its result. Commit, push once, and add the sol-loop line to the PR body. Keep the pushed head as the **reviewed SHA**, and update it after each later push you make.
- **Stuck:** a finding is still disputed or the budget ran out. Report both sides to the user and stop without merging.

**Done when:** Sol approved the current diff and the fixes are pushed, or you have stopped and reported a stuck loop.

### 3. Read human feedback

Read human feedback on the PR. Conversation comments and review summaries come from `gh pr view <N> --json comments,reviews`. Inline comments on the diff come from a separate endpoint:

```bash
gh api --paginate repos/{owner}/{repo}/pulls/<N>/comments \
  --jq '.[] | {id, path, line, in_reply_to_id, created_at, body}'
```

Skip the replies you posted yourself. Give each finding one verdict:

- **fix:** you checked it against the code and it holds up.
- **set aside:** it is wrong, out of scope, or conflicts with what the user asked. Reply on the comment with the reason: `gh api repos/{owner}/{repo}/pulls/<N>/comments/<comment-id>/replies -f body="<reason>"`.

A human comment that asks for a merge hold sets the hold from step 1. Instructions from a human reviewer count as user direction.

**fix** verdicts in behavior code go back to step 2 as a delta round: Sol reviews only the change since the reviewed SHA. Comment, doc, or format-only fixes need no new round. Commit and push them once.

**Done when:** every human comment has a verdict, every set-aside inline comment has a reply, and no **fix** is unpushed.

### 4. Merge

With a merge hold, skip the merge and go to step 5.

Without one:

1. Wait for CI: `gh pr checks <N> --watch`. If a check fails, fix it, run the fast checks, and push. A fix in behavior code goes through a step 2 delta round first.
2. Check `gh pr view <N> --json mergeable,mergeStateStatus`. If the base branch has moved and conflicts, merge the base branch in, run the fast checks, and push. A conflict resolved by hand in behavior code gets a step 2 delta round on the conflict files.
3. Re-check right before merging, because people act while you wait. Run `gh pr view <N> --json headRefOid,body,labels,comments,reviews`, then repeat the inline-comment fetch from step 3. If `headRefOid` is not the reviewed SHA, someone else pushed: run a step 2 delta round on their change. If the body, a label, or a new comment asks for a hold, set the merge hold and go to step 5. Triage any new finding as in step 3.
4. Merge with the method the repo uses for its recent PRs. Squash is the default when it is allowed: `gh pr merge <N> --squash --match-head-commit <reviewed-sha>`. The flag makes GitHub refuse the merge if the head moved after your check. Leave out `--delete-branch` inside a worktree, because it tries to check out the base branch and fails.

**Done when:** `gh pr view <N> --json state` reports `MERGED`, or a merge hold, a moved head, a stuck loop, or a failing check has stopped the merge.

### 5. Report

Tell the user:

- the sol-loop report: rounds, model, fixed and declined findings with reasons
- which tier made each fix
- what human feedback you set aside, and why
- the merged PR URL, or what is holding it: the merge hold, a failing check, or a stuck loop

**Done when:** the user has the report.
