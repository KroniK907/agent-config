---
name: sol-build
description: Orchestrate building one or more grilled maps as a single stack of PRs, with a sol-loop review on every PR by its builder and parallel Sol seam reviews across the stack. Use when the user says sol build, wants several maps or grill sessions combined into one build, or wants Sol to review the PRs.
agent-config-sync: true
---

# Sol build

You are the **orchestrator**. This is the `build` skill with four changes:

1. **Several maps in one stack.** The input is one or more grilled maps. Their slices are merged into a single build order, on a root branch you choose.
2. **Sol reviews each PR.** Each builder opens its PR, then runs the `sol-loop` skill on its own layer. There are no separate babysitters.
3. **Deferred findings.** A finding that belongs to a lower layer goes into `deferred.md`, unless it is a blocker. Those findings are fixed together in the seam fix round.
4. **Seam review.** The stack review is one Sol loop per **seam group**, started as soon as that group's layers are done. Each loop looks at code that links layers to each other, or to the existing repo. Each later round checks only the fixes.

Resolve this skill's folder with `readlink -f`. Then read `../build/SKILL.md` and `../build/BRIEFS.md`. Everything there applies unless this file says otherwise, including keeping the orchestrator lean, the model choices, concurrency, human gates and the report. This skill's own templates are in [BRIEFS.md](BRIEFS.md). Read the `sol-loop` skill once, because builders and seam reviewers both follow it.

Usage: `/sol-build #N [#M ...] [slices] [into <branch>] [merge]`.
- `into <branch>`: the **root** the stack is built on and merged into, for example a feature branch. The default is `dev`. Wherever `build` says `dev`, read root.
- `merge`: pre-approves merging the reviewed stack into root.

## 1. Load

For each map, read the body and every GM comment (`gh issue view N --comments`). Save each one to `/tmp/{slug}-plan/maps/N.md`. `{slug}` is the map slug when there is one map. With several maps, use a short name for the combined build and confirm it with the user. Every map must be `shaped`, with its grill finished. If one is not, stop, say which one, and point at `/grill #N`.

Choose the slices: the ones the user named, or every unbuilt slice of every map.

Read the maps against each other. List each place where:
- two decisions conflict;
- two slices build the same thing;
- one map's slice depends on another map's slice.

Write the list to `/tmp/{slug}-plan/cross-map.md`. Ask the user about each conflict. Record each answer there, and post it as a GM comment on every map it changes. Create an empty `/tmp/{slug}-plan/deferred.md`.

If any map asks for designs first, run them now, as in `build`.

**Done when:**
- every map file exists;
- the slice list and root are fixed;
- every cross-map conflict has a recorded answer;
- any design approval is in.

## 2. Plan

Spawn the planner with `build`'s planner brief and the planner additions in this skill's BRIEFS.md. It writes `review.md` in place of `babysit.md`. It also adds two things to `contract.md`:
- a **Seams** table;
- **Seam groups**: 2 to 4 groups, each a set of layers and the seams that join them.

Ask the user the questions in `questions.md`, as in `build`. Post a **Stack status** comment on every map, using this skill's template, and set each map's **Phase** to `building`.

**Done when:** the contract lists every layer, seam and seam group; questions are answered; and Stack status is posted on every map.

## 3. Dispatch

Use one-line builder prompts that point at this skill's `build.md`. Each builder:
1. builds;
2. runs the checks in the foreground;
3. pushes and opens a ready PR;
4. runs a sol-loop on its layer against its base, and pushes the agreed fixes;
5. reports back.

A builder above can start once the PR below it is open. There are no babysitters. A builder may start the one Sol reviewer per round that sol-loop calls for, and nothing else.

- **Lower-layer findings:** a builder reports a finding in a lower layer's code; it does not fix it. If the finding is a blocker, send a fix delegate to the owning layer now, then sync the layers above it. Otherwise, add it to `deferred.md` with its id, the owning layer and the evidence.
- **Sol unreachable:** if a builder reports `SOL UNAVAILABLE`, run that layer's sol-loop from this thread instead. You start the reviewer, check its findings, and send agreed fixes to a fix delegate. Never swap in another model.
- **Usage limits:** if Sol, or this thread, is close to a usage limit, use `schedule_task` to wake this thread 5 minutes after the reset. Record the wake-up in Stack status. Layers that need no Sol keep building. A builder that hits the limit opens its PR without the sol-loop line, reports `SOL LIMIT <reset time>`, and the layer's loop runs after the wake-up.
- **CI:** after each push, wait on CI with the `babysit` skill's watcher, never a polling loop. Send a red check back to the layer as a fix delegate.
- **Amendments mid-build:** a decision that lands on a map during the build does not change an in-flight layer. Add it to `deferred.md`, and add a follow-up layer to `contract.md` that starts after the stack merges.
- **Tracking:** record each layer's Sol rounds and verdict in Stack status, with its PR, head SHA and checks.

**Done when:** every layer's PR is open with a Sol verdict of `approved` (or is held with a reported reason), and its checks are green.

## 4. Seam review

Review each seam group as soon as every layer in it has its PR open with a Sol verdict of `approved`. Don't wait for the full stack. Later groups keep building meanwhile. A fix found early lands in a low layer before many layers sit on top of it, so the sync after it stays short.

**Sync.** Hand the sync to one delegate on the sync brief in BRIEFS.md. It syncs up to the group's top layer. Do not run it in this thread. The delegate:
1. merges each branch into the layer above it, bottom-up;
2. fixes builds that the merges broke;
3. reruns the checks;
4. pushes;
5. confirms every layer is in the top head.

After a sync, run a new per-layer sol-loop only where a merge conflict was resolved by hand in behavior code. That loop covers only the conflict files.

**Round 1.** Start one Sol reviewer for the group on the seam brief in BRIEFS.md. Groups that close together start together. Each reviewer:
- reads the full diff from root to the group's synced top head;
- reviews the seams in its group;
- is given the per-layer ledgers, so it doesn't re-review the middle of a module;
- rules on the `deferred.md` items in its group: confirmed or not a bug.

Each group keeps its own sol-loop ledger. Running several loops at once overrides sol-loop's one-reviewer rule. Each loop still has one reviewer.

**Fix round.** Check every finding yourself, as sol-loop requires. Then send **one combined fix round**, with one fix delegate per owning layer or per group of nearby layers. Each delegate gets every confirmed finding and deferred item for its layers. A fix goes on the layer that produces the seam, or into a small follow-up PR on top. Use Opus for fixes in money, auth, permissions, migrations or concurrency. When a fix needs a design choice, pick one and say so in Stack status. The user can overrule it when they read it.

**Later rounds.** Sync again. Then each group's reviewer checks only what changed since the head it last approved, plus its open ids. It does not review the full stack again. A group that approved, and whose seams the fix round did not touch, is done.

**Cross-group pass.** When every group has approved, sync the full stack. One Sol reviewer checks only the seams that join groups, with the round 1 brief, then delta rounds as above.

Budget: sol-loop's 3 rounds per group, extended to 5 when a group's seams involve money, auth, permissions, migrations or concurrency. When the budget runs out, or a finding is still disputed after two replies, stop that group. Bring both sides to the user.

**Done when:** every group's Sol, and the cross-group pass, returns `APPROVE` on the synced top head with no blocker or should open, or the user has the open findings. Every `deferred.md` item is fixed, cleared, or moved to a follow-up layer.

## 5. Merge and report

Merge into root bottom-up, as `build` merges into `dev`. Then dispatch any follow-up layers from `deferred.md`, through steps 3 and 4.

The report adds:
- a **Seams** section with each group, its seams, and what its review found and changed;
- the Sol rounds for each PR;
- the design choices you made in fix rounds;
- the follow-up layers.

Link the report from Stack status on every map, and set each map's **Phase** to `built`.

**Done when:** the user has the hatch URL, and every map's Stack status shows the final state of every PR.
