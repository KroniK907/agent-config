---
name: sol-build
description: Orchestrate building one or more grilled maps as a single stack of PRs, with a sol-loop review on every PR instead of reviewbot and a sol-loop seam review across the stack. Use when the user says sol build, wants several maps or grill sessions combined into one build, or wants Sol to review the PRs.
agent-config-sync: true
---

# Sol build

You are the **orchestrator**. This is the `build` skill with three changes:

1. **Several maps in one stack.** The input is one or more grilled maps. Their slices are merged into a single build order.
2. **Sol reviews each PR.** Each builder runs the `sol-loop` skill on its own layer before opening the PR. Reviewbot does not gate anything.
3. **Seam review.** The stack review is a sol-loop over the **seams**: code that links one layer to another, or a layer to the existing repo. The middle of a module is out of scope unless it breaks a seam.

Resolve this skill's folder with `readlink -f`. Then read `../build/SKILL.md` and `../build/BRIEFS.md`. Everything there applies unless this file says otherwise, including keeping the orchestrator lean, the model choices, concurrency, merge rules, human gates and the report. This skill's own templates are in [BRIEFS.md](BRIEFS.md). Read the `sol-loop` skill once, because builders and the seam review both follow it.

Usage: `/sol-build #N [#M ...] [slices] [merge]`.

## 1. Load

For each map, read the body and every GM comment (`gh issue view N --comments`). Save each one to `/tmp/{slug}-plan/maps/N.md`. `{slug}` is the map slug when there is one map. With several maps, use a short name for the combined build and confirm it with the user. Every map must be `shaped` (its grill finished, with **Slices** written). If one is not, stop and say which one, and point at `/grill #N`.

Choose the slices: the ones the user named, or every unbuilt slice of every map.

Then read the maps against each other. List each place where:
- two decisions conflict;
- two slices build the same thing;
- one map's slice depends on another map's slice.

Write the list to `/tmp/{slug}-plan/cross-map.md`. Ask the user about each conflict. Record each answer in `cross-map.md`, and post it as a GM comment on every map it changes.

If any map asks for designs first, run them now, as in `build`.

**Done when:** every map file exists, the slice list is fixed, every cross-map conflict has a recorded answer, and any design approval is in.

## 2. Plan

Spawn the planner with `build`'s planner brief, plus the planner additions in this skill's BRIEFS.md. It writes:
- `review.md` in place of `babysit.md`;
- a **Seams** table in `contract.md`.

Ask the user the questions in `questions.md`, as in `build`. Post a **Stack status** comment on every map, using this skill's template, and set each map's **Phase** to `building`.

**Done when:** the contract lists every layer and every seam, questions are answered, and Stack status is posted on every map.

## 3. Dispatch

Builders are dispatched as in `build`, with a one-line prompt that points at `build.md` (this skill's version). Each builder:
1. builds its layer;
2. runs the checks;
3. runs a sol-loop on its own diff against its base;
4. pushes;
5. opens a ready PR, with the sol-loop line in its body.

There is no babysitter. A builder may start the one Sol reviewer that sol-loop calls for, and nothing else.

- **Sol unreachable:** if a builder reports `SOL UNAVAILABLE`, run that layer's sol-loop from this thread instead. You start the Sol reviewer, check its findings, and send agreed fixes to a fix delegate on the layer's branch. Do not swap in another model.
- **CI:** after each PR opens, check `gh pr checks <N>`. Send a red check back to that layer as a fix delegate. Reviewbot comments, if any arrive, are not a gate. Skim them for a blocker the Sol loop missed and send that blocker to a fix delegate.
- **Lower-layer findings:** when a builder reports a finding that belongs to a lower layer, send a fix delegate to the lower layer. Then sync the stack (step 4).
- **Tracking:** record each layer's Sol rounds and verdict in Stack status, next to its PR, head SHA and checks.

**Done when:** every layer's PR is open with a Sol verdict of `approved` (or is held with a reported reason), and its checks are green.

## 4. Seam review

**Sync first.** Do this exactly as `build` does: merge each base into the branch above it, bottom-up, rerun checks and push. Do not run a new per-layer sol-loop after a sync, unless a merge conflict was resolved by hand in behavior code. In that case, the layer's builder (or a fix delegate) runs one sol-loop round on the conflict files.

Then run a sol-loop on the whole stack from this thread. Use the seam brief in BRIEFS.md:
- **Diff:** from the stack's root base (`dev`) to the synced top head.
- **Seams:** the Seams table from `contract.md`.
- **Already reviewed:** each layer's sol-loop ledger, so Sol does not re-review the middle of a module.

Check each finding yourself, as sol-loop requires. Send each agreed fix to a fix delegate on the layer that owns the producing side of the seam, or into a small follow-up PR on top. Sync again after each fix round, then start the next Sol round with the ledger.

Budget: sol-loop's 3 rounds, extended to 5 when a seam involves auth, permissions, data migrations or concurrency. When the budget runs out, or a finding is still disputed after two replies, stop and bring both sides to the user.

**Done when:** Sol returns `APPROVE` on the synced top head and no blocker or should is open, or the user has the open findings.

## 5. Merge and report

Merge and report exactly as `build` does. The report adds:
- a **Seams** section with each seam, the layers on either side, and what the seam review found and changed there;
- the Sol rounds for each PR.

Link the report from Stack status on every map, and set each map's **Phase** to `built`.

**Done when:** the user has the hatch URL, and every map's Stack status shows the final state of every PR.
