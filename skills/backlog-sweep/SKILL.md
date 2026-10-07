---
name: backlog-sweep
description: "Sweep a repo's open GitHub issues: close what shipped or went stale, narrow what is half done, and label small quick work priority: high. Run by hand or on a schedule."
disable-model-invocation: true
agent-config-sync: true
---

A **sweep** reads every open issue in the current repo against what has actually merged, then makes four kinds of change: close **shipped** issues, close **stale** ones, **narrow** half-done ones to what remains, and label **quick** ones `priority: high`. Every change cites its **evidence** in a comment on the issue, so anyone reading the issue later can check the call without redoing the sweep.

The user may pass `apply` (make changes without asking) or `dry-run` (report only). With neither, propose and wait for approval before changing anything.

## Steps

### 1. Gather

Resolve the repo (`gh repo view --json nameWithOwner,defaultBranchRef`) and fetch the default branch. Run the bundled script with an absolute path (resolve this skill's folder with `readlink -f` first):

```bash
<skill-dir>/scripts/gather.sh <owner/repo> "$(mktemp -d)/sweep" 100
```

It writes every open issue with its comments to `issues/<N>.md`, plus `index.tsv`, `merged-prs.tsv`, and `labels.txt`. Read `index.tsv`, then read every issue file. Run `gh` from outside the repo only with `-R <owner/repo>`.

Also read what the repo and the user already know about recent state: the default branch's `git log` since the oldest open issue, release notes or a changelog, and any memory or notes on releases and infrastructure.

**Done when:** you have read every issue file in full, including its comments, and you know which PRs merged since the oldest open issue was filed.

### 2. Classify

Give each issue exactly one verdict from the table in [Verdicts](#verdicts). For each verdict except **keep**, collect the evidence before you decide:

- **Shipped** needs a merged PR or commit that does what the issue asks. Confirm it reached the default branch: a PR merged into a stacked or feature branch counts only if `git merge-base --is-ancestor <sha> origin/<default>` succeeds. Then confirm the code is there: grep for the function, migration, config value, or file the issue names. A PR title alone is not evidence.
- **Stale** needs the thing that replaced it: a later decision, a different shipped design, a removed system, or a ticket that was skipped because the feature was built without it. Name it.
- **Narrow** needs the checklist split into done (with evidence) and remaining.
- **Quick** needs the remaining work to fit in one PR or one short human session, with no open decision, no blocker, and no pending design.

Grep the code and config for every issue whose steps name a file, flag, version, or table, even when no PR mentions it. Half-done issues show up this way: a migration that ran while a version pin stayed behind, for example.

**Done when:** every open issue has one verdict, and every non-keep verdict has a cited PR, commit, file, or decision.

### 3. Propose

Show the user one table, grouped by verdict, with the issue number, title, verdict, and one-line evidence. List the planning-tracker edits from [Planning trackers](#planning-trackers) under it.

In `dry-run`, stop here. In `apply`, continue without waiting, but apply only **shipped** and **quick** rows: their evidence is mechanical. Leave **stale** and **narrow** rows as proposals in the report, because each is a judgement call someone should see first. Otherwise wait for the user to approve, edit, or drop rows.

**Done when:** the user has approved the list, or the mode says to skip approval.

### 4. Apply

Make each approved change with `gh`. Write each comment as the evidence for the change in one to three sentences. Link PRs by number (`#123`) so GitHub cross-references them.

| Verdict | Action |
|---|---|
| shipped | `gh issue close N --reason completed --comment "<what shipped, in which PR>"` |
| stale | `gh issue close N --reason "not planned" --comment "<what replaced it>"` |
| narrow | Retitle to the remaining work when the old title no longer fits, then comment the done/remaining split as a checklist. Add `priority: high` if the remainder is quick. |
| quick | `gh issue edit N --add-label "priority: high"` plus a comment saying why it is quick and what to do first. |

Create the label once if it is missing: `gh label create "priority: high" --color B60205 --description "Small, quick, and worth doing next"`.

Re-check issues that already carry `priority: high`. Remove the label when the work grew, picked up a blocker, or shipped.

**Done when:** every approved row is applied, and `gh issue list --state open` no longer lists any issue you closed.

### 5. Report

Give the user a short summary: counts per verdict, the issues now labelled `priority: high` in suggested order, and anything you could not decide (missing access, a judgement call, a question only the user can answer).

**Done when:** the user has the summary, and every open question in it names its issue number.

## Verdicts

| Verdict | Meaning |
|---|---|
| **shipped** | A merged change on the default branch does what the issue asks. Partial delivery is **narrow**, not shipped. |
| **stale** | No longer worth doing as written: superseded by a later decision or design, made moot by a removed system, or a planning step the team skipped by building the feature. |
| **narrow** | Part of it shipped. The rest is still wanted. |
| **quick** | Still wanted, small, and unblocked. |
| **keep** | Still wanted and not quick, blocked, or waiting on a decision. |

Parked ideas are **keep**: an issue that says "not scheduled" or "feature idea" records a deliberate deferral, and its age is not evidence that it went stale. Planning trackers (maps, decision logs, research checklists) that still list open questions are **keep** too.

Security fixes, cost savings, and version pins that drift from production deserve **quick** whenever they are small, even if the issue is labelled for a human (`wf:hitl`). The label says who does it, not how big it is.

## Planning trackers

Some repos track features through map issues that list child tickets (wayfinder `{FeatureName}:Map` issues with **To Do**, **Implementing**, and **Completed** sections, plus decision logs). When you close a ticket that a map lists:

- Move its row from **To Do** or **Implementing** to **Completed**, with a one-line gist and the PR that delivered it.
- Update its decision rows in the coverage table if the map tracks status there (for example `open` → `implemented`).
- Update the map's **Phase** line when the feature it tracks has shipped.

Edit the body surgically: download it (`gh issue view N --json body -q .body > map.md`), change only those lines, `diff` against the live body, then `gh issue edit N --body-file map.md`. Close a map only when the user approves it explicitly. When a research or checklist issue has items a recent PR finished, comment which items are now done rather than rewriting the checklist.

## Running on a schedule

Run `/backlog-sweep apply` from a recurring job (`/schedule` for a cloud routine, `/loop` in a long-lived session). Weekly is enough for most repos. The job needs `gh` authenticated with issue write access and a checkout it can `git fetch` in.

## Gotchas

- `gh issue close --reason` takes `"not planned"` with a space.
- `gh pr list --json` has no closing-issue field. The gather script reads it through GraphQL, and most teams leave it empty, so search PR titles, bodies, and commit messages for the issue number and its subject.
- `gh` run outside a git checkout fails with "not a git repository" unless you pass `-R <owner/repo>`.
