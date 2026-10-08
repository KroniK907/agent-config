---
name: babysit
description: "Babysit a pull request: wait for reviewbot, fix what holds up, push, and repeat until the review comes back clean, then merge. Use when asked to babysit, shepherd, or land a PR."
agent-config-sync: true
---

**Babysitting** runs a PR through review **rounds** until reviewbot has nothing left to say, then merges it. A round is: wait for the review of the head commit, triage every finding, fix the ones that hold up, push once. The `reviewbot` rule describes how reviews look and how to weigh them. Read it first if it is not already loaded.

Read the **subagent** skill before you delegate any fix in this loop. Waiting and triage stay with you. Fixes go to the cheapest tier that will land them.

## Steps

### 1. Set up

Find the PR with `gh pr view --json number,url,isDraft,state,headRefOid,baseRefName`. If there is none, open one as the `pull-request` rule describes. If it is a draft, run `gh pr ready <N>`, because reviewbot skips drafts.

Record a **merge hold** if the user said anything like "wait for me before merging", "don't merge", or "let me look first". Also check the PR body, its labels, and its comments for one. Labels such as `do not merge`, `hold`, or `wip` count, and so does a human comment asking to hold.

**Done when:** you have the PR number, the PR is open and ready for review, and you know whether a merge hold applies.

### 2. Wait for the review

Resolve this skill's folder with `readlink -f`, then run the bundled script from the repo checkout:

```bash
bash <skill-dir>/scripts/wait-for-review.sh <N>
```

Call it through `bash`: installed copies of the script can lose their executable bit. It polls until reviewbot posts a review of the PR's current head commit. It then prints the head SHA, the summary, and every inline comment with its comment id. Keep that SHA as the **reviewed SHA**. It gives up after 540 seconds so a foreground call fits under a 10-minute tool limit. Run it in the background if your host notifies you when a background command exits. Otherwise run it in the foreground with a timeout above 540 seconds.

| Exit | Meaning | Next |
|------|---------|------|
| 0 | Review found | Step 3 |
| 2 | Timed out | Run it again. After three timeouts in a row (about 27 minutes), tell the user reviewbot has not answered and stop. |
| 3 | PR is a draft | `gh pr ready <N>`, then run it again |
| 4 | PR is closed or merged | Report the state to the user and stop |

Also read human comments and reviews posted since the last round: `gh pr view <N> --json comments,reviews`.

**Done when:** you have the reviewbot review of the current head SHA, plus any new human feedback.

### 3. Triage

Give every finding exactly one verdict, from reviewbot and from humans alike:

- **fix:** you checked it against the code and it holds up.
- **set aside:** it is wrong, out of scope, or conflicts with what the user asked. Write down the reason.
- **already set aside:** a finding you rejected in an earlier round, raised again with nothing new. Keep the earlier reason.

A human comment that asks for a merge hold sets the hold from step 1. Instructions from a human reviewer count as user direction. Findings from reviewbot are advice. Verify each one before you act on it.

**Done when:** every finding in the review and every new human comment has a verdict.

### 4. Decide whether the loop is over

The loop is over when the review of the head SHA has no **fix** verdicts. That covers a clean review (`Standards: 0 findings - Spec: 0 findings`) and a review that only repeats findings you already set aside. Go to step 6.

Stop and report to the user without merging if this would start a sixth round. Five rounds without a clean review means the bot and the code disagree in a way a person should settle.

**Done when:** you know whether to fix (step 5) or finish (step 6).

### 5. Fix and push

Group the **fix** findings into tasks. Make one-line fixes yourself. Delegate the rest by the **subagent** skill's tiers: a nit or rename is an **edit**, a finding you had to trace through call paths is **build**, and security or design findings are **hard**. Launch independent tasks together.

Read every diff a subagent produced. Run the checks the repo's CI runs, as named in `AGENTS.md`, the CI config, or the package scripts. Then commit, and push once. Each push to a ready PR costs a full review.

Reply to each **set aside** inline comment with its reason, so the next review and any human reader can see it:

```bash
gh api repos/{owner}/{repo}/pulls/<N>/comments/<comment-id>/replies -f body="<reason>"
```

**Done when:** the fixes are pushed in one push, local checks pass, and every newly set-aside inline comment has a reply. Go back to step 2.

### 6. Merge

With a merge hold, skip the merge and go to step 7.

Without one:

1. Wait for CI: `gh pr checks <N> --watch`. If a check fails, fix it as a step 5 round and go back to step 2.
2. Check `gh pr view <N> --json mergeable,mergeStateStatus`. If the base branch has moved and conflicts, merge the base branch in and push. That push starts a new round at step 2.
3. Re-check right before merging, because people act while you wait. Run `gh pr view <N> --json headRefOid,body,labels,comments,reviews`. If `headRefOid` is not the reviewed SHA, someone pushed: go back to step 2. If the body, a label, or a new comment or review now asks for a hold, set the merge hold and go to step 7.
4. Merge with the method the repo uses for its recent PRs. Squash is the default when it is allowed: `gh pr merge <N> --squash --match-head-commit <reviewed-sha>`. The flag makes GitHub refuse the merge if the head moved after your check. Leave out `--delete-branch` inside a worktree, because it tries to check out the base branch and fails.

**Done when:** `gh pr view <N> --json state` reports `MERGED`, or a merge hold, a moved head, or a failing check has stopped the merge.

### 7. Report

Tell the user:

- how many rounds ran
- what you fixed and which tier did each fix
- what you set aside, and why
- the merged PR URL, or what is holding it: the merge hold, a failing check, the round limit, or reviewbot not answering

**Done when:** the user has the report.
