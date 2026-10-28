---
name: build
description: Orchestrate building a shaped map as a stack of PRs, with cheaper subagents per slice, a sol-loop babysitter per PR, cross-model reviews per group of layers and an HTML report. Use when the user says to implement or build a map, or makes you the orchestrator.
agent-config-sync: true
---

# Build

You are the **orchestrator**. Delegates write the code; you plan, dispatch, track and merge. Your thread is usually the most expensive part of a build, so keep it **lean**: delegates read their briefs from files, report back in 8 lines or fewer, and you only make one-line fixes yourself.

Usage: `/build #N [slices] [merge]`. `merge` (or "merged into dev" in the request) pre-approves merging the stack into `dev`. Without it, stop when the stack is reviewed and clean.

Brief templates are in [BRIEFS.md](BRIEFS.md). Pick models with the `subagent` skill. Delegates babysit with the `babysit` skill and write PR bodies with the `pr` skill.

## 1. Load

Read the map body and every GM comment (`gh issue view N --comments`), and save them to `/tmp/{slug}-plan/map.md`. Choose the slices to build: the ones the user named, or all unbuilt slices.

If the map or the user asks for designs first (module interfaces with `design-modules`, or UI mocks), run that now and wait for approval.

**Done when:** `map.md` exists, the slice list is fixed, and any design approval is in.

## 2. Plan

Spawn one **planner** (build tier, read-only on the repo). Its brief is the planner template in BRIEFS.md. It writes `contract.md`, `common.md`, `build.md`, `babysit.md` and `questions.md` into `/tmp/{slug}-plan/`.

Read `questions.md`. Ask the user only those questions, then record the answers in `contract.md`. Post a **Stack status** comment on the map (template in BRIEFS.md), and set the map's **Phase** to `building`. You'll edit that comment in place from now on. It is the durable record if `/tmp` or this thread is lost.

**Done when:** the contract lists every layer with its branch, base, model and what it owns; questions are answered; Stack status is posted.

## 3. Dispatch

Each layer gets a **builder**. It builds, runs the checks, and opens a ready (not draft) PR against its base, then stops. Once the PR is open, a **babysitter** takes it through a sol-loop without merging. A builder can start as soon as the PR below it is open; it merges the base branch again right before it pushes.

- **Prompts are one line:** `Follow /tmp/{slug}-plan/build.md exactly. Your layer: NN-slug, branch {area}/NN-slug, base <base>, worktree /tmp/{slug}-wt/NN-slug.`
- **Models:**
  - Use Sonnet 5.5 for layers the contract fully specifies, and for UI.
  - Use Opus 5.5 at medium for layers with design risk: concurrency, security, migrations of live data, build tooling.
  - Use Opus 5.5 at high only for a layer that already failed once.
  - Babysitters run on Sonnet 5.5.
- **How to delegate:** on T3 Code, use `delegate_task`, which can also reach other providers (see `orchestrator_capabilities`). Elsewhere, use native subagents. Each delegate starts cold, never as a fork of this thread.
- **Concurrency:** model work runs remotely, but builds and tests run on this machine. Run up to about three Opus delegates or six Sonnet delegates at once. Start fewer when `uptime` shows a load average above twice `nproc`, and start no new builder until it drops. An overloaded machine slows every check and turns timing tests into flakes. Wait for completion notifications in the background instead of polling.
- **One full run per push:** builders run fast checks and their own tests, as common.md says. The PR's CI runs the full suites. Wait on them with the `babysit` skill's watcher, never a polling loop, and send a red check back to the layer.
- **Stalled delegates:** a delegate that reports back without its push or PR usually left a check running in the background. Send it straight back to read that output, or to rerun the check in the foreground, then finish its brief.
- **Sol unreachable:** if a babysitter reports `SOL UNAVAILABLE`, run that layer's sol-loop from this thread instead. You start the reviewer, check its findings, and send agreed fixes to a fix delegate. Never swap in another model.
- **Tracking:** after each report, update Stack status with the PR number, head SHA, check state and any deviation. If a delegate dies on a rate or session limit, write `resume.md` (template in BRIEFS.md) and relaunch it.
- **Long builds:** when this thread grows heavy, start a fresh orchestrator thread per wave. Point it at the plan directory and the Stack status comment.

Step 4 runs alongside this step. Start each group's review as soon as its layers are ready.

**Done when:** every layer's PR is open with a Sol verdict of `approved` (or held with a reported reason), and checks are green.

## 4. Review each group of layers

Don't wait for the whole stack. A cross-PR bug found at the end is usually in a low layer. Its fix then merges up through every layer above, and each sync of a long stack takes hours. The contract splits the layers into **review groups**. Review each group as soon as its layers' PRs are open with a Sol verdict of `approved`, while later groups keep building.

**Sync first.** A lower layer can change after the layers above it stopped pushing. Before each review round, walk the stack bottom-up to the group's top layer. Merge each branch's base into it with a merge commit, rerun its fast checks, and push it if the merge brought anything in. A push that only merges the base needs no babysit round, unless a conflict was resolved by hand in behavior code. Then that layer gets a sol-loop on the conflict files only. Hand a sync of more than a few layers to one delegate rather than running it in this thread.

**Round 1.** Run a read-only review of the group by a model from another family (GPT-6.1 Sol on T3 Code). In parallel, run an independent Claude pass. Both use the group review brief in BRIEFS.md on the group's synced top head. Send verified findings to fix delegates, one per owning layer. A fix goes on the layer that owns the code, or into a small follow-up PR on top.

**Later rounds.** Sync again. Each reviewer then checks only the change since the head it last reviewed, plus its open findings, with the re-review brief. It never reviews the full group again. After 5 rounds, stop and bring what's left to the user.

**Final pass.** When every group is clean, sync the full stack. One reviewer checks only the seams that join groups, then delta rounds as above.

**Done when:** both reviewers report `NO FINDINGS` on every group and on the final pass, or the user has the remaining findings.

## 5. Merge, if pre-approved

Merge bottom-up with merge commits:
1. Merge the lowest PR into `dev`.
2. Retarget the next PR to `dev` (`gh api -X PATCH repos/{owner}/{repo}/pulls/<N> -f base=dev`).
3. Wait for its checks, then merge it.
4. Repeat for each PR up the stack.

Don't rebase or force-push branches that have PRs open.

**Human gates:** stop and ask before any infra apply, any write to production data, merging to `dev` without approval, and anything in `release`. List the exact commands the human must run.

**Done when:** every PR is merged into `dev`, or you are stopped at a named gate.

## 6. Report

Build a self-contained HTML report covering:
- what was built, PR by PR;
- what changed from the spec, and why;
- problems hit and how they were solved;
- human steps still owed;
- screenshots of every UI state in the layers' notes files, at desktop and phone widths;
- where the time went: when each phase started and ended, and the longest waits (checks queued, reviews late, stalled delegates). Take the times from the thread and Stack status.

Serve it with `hatch`, link it from Stack status, and set **Phase** to `built`.

**Done when:** the user has the hatch URL and Stack status shows the final state of every PR.
