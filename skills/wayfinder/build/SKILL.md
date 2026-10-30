---
name: build
description: Orchestrate building a shaped map as a stack of PRs, with cheaper subagents per slice, a sol-loop babysitter per PR, cross-model reviews per group of layers and an HTML report. Use when the user says to implement or build a map, or makes you the orchestrator.
agent-config-sync: true
---

# Build

You are the **orchestrator**. Delegates write the code; you plan, dispatch, track and merge. Your thread is usually the most expensive part of a build, so keep it **lean**: delegates read their briefs from files, report back in 8 lines or fewer, and you only make one-line fixes yourself. When the thread still grows long, hand the build to a fresh thread (see **Thread handoff**).

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
- **Tracking:** after each report, update Stack status with the PR number, head SHA, check state and any deviation. Keep `/tmp/{slug}-plan/tasks.md` current too: each running delegate (task ID, layer, role, model, start time), each `schedule_task` bound to this thread with its prompt, and each PR this thread watches. If a delegate dies on a rate or session limit, write `resume.md` (template in BRIEFS.md) and relaunch it.
- **Long builds:** check the handoff triggers in **Thread handoff** after each report.

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

The synced top head already holds every layer. It passed the seam reviews and has green CI. Don't merge the PRs one at a time. Each merge moves `dev`, the next branch has to merge it back in, and CI runs again on a head nobody will review. On a long stack that adds hours and checks nothing new. Merge the stack once:

1. **Confirm the stack.** CI is green on the top PR's head. Every layer's head is in the top head: `git merge-base --is-ancestor <layer head> <top head>` for each one. If a layer's head is missing, sync first.
2. **Retarget every PR to `dev`** (`gh api -X PATCH repos/{owner}/{repo}/pulls/<N> -f base=dev`). Push nothing to the layer branches.
3. **Merge the top PR** with a merge commit: `gh pr merge <top> --merge --match-head-commit <top head>`. If it conflicts with `dev`, merge `dev` into the top branch, wait for CI on that one push, then merge.
4. **Check the rest.** GitHub marks a PR merged once its head commit is in its base branch, so every lower PR closes as merged with the top one. Confirm each with `gh pr view <N> --json state`. A PR still open means its head isn't in `dev`. Report it; don't merge it on its own.

Don't rebase or force-push branches that have PRs open. Squash would break step 4, so always merge with a merge commit.

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

## Thread handoff

Every turn re-reads the whole thread. A thread that has run for hours gets slow and expensive, and its early instructions fade. Hand off at the next break (a group merged, a wave dispatched, a sync done) once any of these is true:
- the harness has compacted this thread's context;
- this thread has taken about 20 delegate reports;
- this thread has run about 4 hours.

Never fork this thread with `t3_thread_transfer`. A fork copies the context you are trying to drop. The successor starts cold and learns the build from files.

1. **Settle.** Start no new delegates. Bring Stack status and `tasks.md` up to date.
2. **Collect scripts.** Move every script and helper file this thread wrote outside the plan dir into `/tmp/{slug}-plan/`: status pushers, check wrappers, lock files, brief drafts. Fix the paths inside them and in the briefs that call them. The successor finds only what the plan dir holds.
3. **Write `HANDOFF.md`** in the plan dir from the template in BRIEFS.md. Copy the original command and every standing instruction the user gave since, word for word. A paraphrase drops merge approvals and gates.
4. **Release this thread's bindings.** Delete each `schedule_task` bound to this thread (`list_scheduled_tasks`, then `delete_scheduled_task`). Call `unwatch_pull_request` on each PR it watches. HANDOFF.md lists both so the successor recreates them. Leave running delegates alone, because cancelling one throws away its work.
5. **Launch the successor** with `t3_thread_launch`: this project, no `workspaceStrategy`, your model and runtime mode, title `{FeatureName} build orchestrator (wave N+1)`. Start the message with the original command, so the skill loads again:
   ```text
   /build #N [slices] [merge] Take over as orchestrator of the {FeatureName} build, wave N+1. Read /tmp/{slug}-plan/HANDOFF.md first. Recreate the schedules and PR watches it lists, bound to your thread. Layers it marks "old thread" stay with their delegates until the old thread relays their reports. Carry on until the build is done.
   ```
6. **Record it.** Put the successor's thread link in Stack status and at the top of HANDOFF.md.
7. **Relay.** A delegate reports to the thread that started it, so its result still lands here. From now on this thread only relays. When a delegate reports, send it to the successor with `t3_thread_send` in 3 lines or fewer: layer, PR, head, verdict, checks, notes file. Start nothing, fix nothing, and leave Stack status alone. If the successor asks you to cancel a delegate, `task_cancel` it. After the last report, send a final line saying so, and stop.

**As the successor,** check before you act. Read HANDOFF.md, Stack status, the open PRs and the notes files, and confirm each claim. A layer marked "old thread" belongs to its delegate. Don't sync, fix or push to its branch until the relay reports it done. If its PR and notes show it finished and its worktree has sat idle, ask the old thread to cancel the delegate first. Two agents pushing to one branch collide. Start your own `tasks.md` section, and hand off again at the same triggers.

Outside T3 Code, do steps 1 to 3, then give the user the launch message to paste into a new session.
