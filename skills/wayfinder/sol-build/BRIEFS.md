# Sol build briefs

These templates change or add to `../build/BRIEFS.md`. Use that file for the planner brief, the common.md outline and resume.md. Replace `{...}` placeholders.

## Planner additions (append to build's planner brief)

```text
This build combines {one map | maps #N, #M, ...}. Inputs: /tmp/{slug}-plan/maps/*.md and cross-map.md,
in place of map.md. cross-map.md beats any single map where they conflict. Order layers by every map's
build order, then by cross-map dependencies.

In contract.md, add a Seams table after the layer table:
  id | kind | producer (layer or existing path) | consumers (layers or existing paths) | contract
Kinds: type or schema, function or API signature, HTTP route and payload, table, migration order, event or
queue message, permission key, config or env var, route or path constant, generated code, UI component props.
List every point where one layer's code calls, reads or renders another layer's code, and every point where
a layer touches code that already exists in the repo. Contract = the exact shape both sides agree to, plus
the errors and empty states it carries. Seam ids are S01, S02, ...

Write review.md (from the template), not babysit.md. In common.md, the PR section adds the sol-loop line
to the PR body, and the Hard limits section allows one Sol reviewer per sol-loop round.
```

## build.md

```markdown
You build ONE layer of the {slug} stack. Your layer, branch, base and worktree are in your task.
1. Read common.md and contract.md (build only your row and its section; keep the shared names, and keep every
   seam your layer produces or consumes exactly as the Seams table states). Read earlier layers' notes-*.md.
2. Create your worktree: `git -C {repo} fetch origin && git -C {repo} worktree add {worktree} -b {branch} origin/{base}`
   (reuse it if it exists). Parallel-group layers merge their sibling branch before review.
3. Build, run every check in common.md, and commit. Do not push yet.
4. Follow review.md.
5. Merge the base branch, rerun the checks, push, and open the PR. Do not babysit or merge.
6. Write notes-NN.md and report per common.md, plus: Sol rounds, verdict, and any finding for a lower layer.
```

## review.md

```markdown
Run the sol-loop skill on your layer, with these overrides:
- Starting the one Sol reviewer it calls for is part of your task. Start nothing else.
- Base: the merge-base with origin/{your base}, not the default branch. Scope: your layer's commits only.
- Intent: your contract.md row. Binding: the maps' GM decisions and cross-map.md. Add to the brief:
  "Behavior that matches a GM decision or the Seams table is not a bug. Out of scope: {other layers}."
- Add the seams your layer produces or consumes (ids and contract) to the brief, and ask Sol to check
  that your side matches them.
- A finding in a lower layer's code: do not fix it. Report it.
- When no route reaches Sol, do not swap in another model. Report `SOL UNAVAILABLE` with what you tried, and
  open the PR without the sol-loop line.
- Log each round (id -> fixed, disagreed or declined, and why) in notes-NN.md.
- Push only after the loop ends, as part of build.md step 5.
```

## Seam review brief (round 1; later rounds follow sol-loop step 4)

```text
Review only. Do not edit, commit, push, comment on PRs, or start subagents.
Repo: {abs path}, detached at {top head}. Base: {merge-base with origin/dev}. Stack: {PRs bottom-up with branches}.
See: git diff {base} {top head}; contract.md (layers and Seams table); cross-map.md; maps/*.md.
GM decisions and the Seams table are binding. Behavior that matches them is not a bug.
Each layer already passed its own sol-loop. Ledgers: {per layer: ids and outcomes}. Do not re-review the
middle of a module. Review the seams: where code in one layer calls, reads, renders or migrates against code
in another layer or in the existing repo. For each seam, read both sides and check:
- both sides use the same names, shapes, nullability, units and enum values;
- errors, empty results and partial failures cross the seam and are handled on the far side;
- auth, permission and tenant checks are not lost at the handoff;
- migration order and the data each layer assumes exists;
- front end and back end agree on routes, payloads and status codes;
- config, env and feature flags are set where they are read;
- existing callers of code a layer changed still work;
- the seams the contract missed: cross-layer calls not in the Seams table.
blocker = an intended flow breaks across a seam. should = a real mismatch callers will hit. nit = rest.
Reply:
VERDICT: APPROVE or CHANGES
X1 [blocker|should|nit] <seam id or "unlisted"> <file>:<line> <claim>. Fix: <suggestion>. Owner: <layer>
```

## Stack status comment

```markdown
**Stack status** (edited in place by the orchestrator)
Plan: `/tmp/{slug}-plan/` on {machine}. Maps: {#N, #M}. Merge: {pre-approved | hold}.

| # | Branch | PR | Map | Model | State | Sol | Head | Notes |
|---|---|---|---|---|---|---|---|---|
| 01 | {area}/01-slug | #N | #N | sonnet | building / sol review / open / merged | 2 rounds, approved | abc1234 | deviation |

Seam review: round N, {open findings | approved}. Human steps owed: {list}. Report: {hatch URL}.
```
