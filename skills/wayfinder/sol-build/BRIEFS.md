# Sol build briefs

These templates change or add to `../build/BRIEFS.md`. Use that file for the planner brief, the common.md outline and resume.md. Replace `{...}` placeholders. Wherever those templates say `dev`, use the build's root branch `{root}`.

## Planner additions (append to build's planner brief)

```text
This build combines {one map | maps #N, #M, ...} and is rooted on {root}. Inputs: /tmp/{slug}-plan/maps/*.md
and cross-map.md, in place of map.md. cross-map.md beats any single map where they conflict. Order layers by
every map's build order, then by cross-map dependencies. Layer 01's base is {root}.

In contract.md, add a Seams table after the layer table:
  id | kind | producer (layer or existing path) | consumers (layers or existing paths) | contract
Kinds: type or schema, function or API signature, HTTP route and payload, table, migration order, event or
queue message, permission key, config or env var, route or path constant, generated code, UI component props.
List every point where one layer's code calls, reads or renders another layer's code, and every point where
a layer touches code that already exists in the repo. Contract = the exact shape both sides agree to, plus
the errors and empty states it carries. Seam ids are S01, S02, ...

Then add Seam groups: 2 to 4 groups by domain (for example payments, scheduling, setup). Each group lists
its layers, its seams, and the seams that join it to other groups. Every seam belongs to exactly one group.
These replace build's Review groups. Order the layers so each group's layers come together and the first
group closes as early as the dependencies allow, because each group is reviewed as soon as it closes.

Write review.md (from the template), not babysit.md. In common.md:
- PR section: add the sol-loop line to the PR body.
- Hard limits: allow one Sol reviewer per sol-loop round.
- Checks: follow build's three-way split. Before a sol-loop round, run the fast checks and the layer's own
  tests. Leave full suites to CI.
- Git: merge only the branches the contract lists for your layer. If you need another layer's code, stop
  and report it.
```

## build.md

```markdown
You build ONE layer of the {slug} stack. Your layer, branch, base and worktree are in your task.
1. Read common.md and contract.md (build only your row and its section; keep the shared names, and keep every
   seam your layer produces or consumes exactly as the Seams table states). Read earlier layers' notes-*.md.
2. Create your worktree: `git -C {repo} fetch origin && git -C {repo} worktree add {worktree} -b {branch} origin/{base}`
   (reuse it if it exists). Parallel-group layers merge their sibling branch before opening the PR.
3. Build, run every check in common.md, and commit.
4. Merge the base branch, rerun the fast checks, push, and open the PR. CI runs the full suites.
5. Follow review.md.
6. Write notes-NN.md and report per common.md, plus: Sol rounds, verdict, and any finding for a lower layer.
Do not merge.
```

## review.md

```markdown
Run the sol-loop skill on your layer's PR, with these overrides:
- Starting the one Sol reviewer it calls for is part of your task. Start nothing else.
- Base: the merge-base with origin/{your base}, not the default branch. Scope: your layer's commits only.
- Intent: your contract.md row. Binding: the maps' GM decisions and cross-map.md. Add to the brief:
  "Behavior that matches a GM decision or the Seams table is not a bug. Out of scope: {other layers}."
- Add the seams your layer produces or consumes (ids and contract) to the brief, and ask Sol to check
  that your side matches them.
- A finding in a lower layer's code: do not fix it. Report it with its evidence and say if it is a blocker.
- No route reaches Sol: report `SOL UNAVAILABLE` with what you tried. Sol is over its usage limit: report
  `SOL LIMIT <reset time>`. Do not swap in another model in either case.
- Push agreed fixes to your PR after the loop ends, merging your base first and rerunning the fast checks.
  Then wait for them with `pr-wait --until green --no-reviewers` from the `babysit` skill until CI finishes, and fix a red check before you report.
- Log each round (id -> fixed, disagreed or declined, and why) in notes-NN.md.
```

## deferred.md

```markdown
# Deferred findings ({slug})
| id | found on | owning layer | seam group | severity | claim and evidence | outcome |
|---|---|---|---|---|---|---|
| D01 | 41 (#N) | 35 | payments | should | preview API ignores the overridden amount, file:line | open / confirmed / not a bug / fixed in <sha> / follow-up layer NN |

## Follow-ups after merge
- <amendment or out-of-layer work>: decisions {MAP-SLUG}-GM-NNN, planned as layer NN.
```

## Sync brief

```text
Sync the {slug} stack bottom-up, from layer 01 on {root} to {top layer of the group, or the top layer}. Plan: /tmp/{slug}-plan/.
Read common.md and contract.md. For each layer in order: merge its base (and any sibling the contract lists)
with a merge commit, resolve conflicts, fix builds broken by changes below (signatures, nullability, renames)
without changing behavior, and run build, format and the short tests of the packages the merge touched.
Run every check in the foreground. Push each branch that changed; CI runs the full suites on each push.
Don't run full suites locally. Once all pushes are done, wait on each pushed PR with `pr-wait --until green --no-reviewers` from the `babysit` skill, and
reproduce a red check locally only to fix it. Never rebase or force-push. At the end, confirm every layer's head is in the top head
(`git merge-base --is-ancestor`). Log each layer (merged, conflicts, fixes, checks) in sync-rN-log.md, and
mark conflicts resolved by hand in behavior code.
Report in 8 lines or fewer: top head SHA, layers pushed, hand-resolved conflicts, check failures.
```

## Seam review brief (round 1, one per seam group)

```text
Review only. Do not edit, commit, push, comment on PRs, or start subagents.
Repo: {abs path}, detached at {group's top head}. Base: {merge-base with origin/{root}}. Stack: {PRs bottom-up with branches}.
Don't run full test suites. CI already ran them on every head.
Your seam group: {group name}: layers {list}; seams {ids}; seams joining other groups {ids}.
See: git diff {base} {top head}; contract.md (layers, Seams table, Seam groups); cross-map.md; maps/*.md.
GM decisions and the Seams table are binding. Behavior that matches them is not a bug.
Each layer already passed its own sol-loop. Ledgers: {per layer in your group: ids and outcomes}. Do not
re-review the middle of a module. Review the seams: where code in one layer calls, reads, renders or migrates
against code in another layer or in the existing repo. For each seam, read both sides and check:
- both sides use the same names, shapes, nullability, units and enum values;
- errors, empty results and partial failures cross the seam and are handled on the far side;
- auth, permission and tenant checks are not lost at the handoff;
- migration order and the data each layer assumes exists;
- front end and back end agree on routes, payloads and status codes;
- config, env and feature flags are set where they are read;
- races where two layers touch the same rows, and retries that repeat a side effect;
- existing callers of code a layer changed still work;
- the seams the contract missed: cross-layer calls in your group's layers not in the Seams table.
Also rule on each deferred.md item in your group: confirmed or not a bug, with evidence: {items}.
blocker = an intended flow breaks across a seam. should = a real mismatch callers will hit. nit = rest.
Write your reply to /tmp/{slug}-plan/review-rN-{group}.md and return it:
VERDICT: APPROVE or CHANGES
X1 [blocker|should|nit] <seam id or "unlisted"> <file>:<line> <claim>. Fix: <suggestion>. Owner: <layer>
D01 confirmed|not a bug <evidence>
```

## Seam re-review brief (rounds 2 and later)

```text
Review only, same rules and output as round 1. Group: {group name}. You last reviewed {previous top head};
the stack is now synced at {top head}. Review only: git diff {previous top head} {top head}, limited to
the fixes for your open ids and the seams those fixes touch. Ledger: {fixed: / disagree: / declined: with
evidence}. Mark each open id resolved or still open, with the reason. Raise a new finding only if a fix
caused it. Do not review the full stack again.
```

## HANDOFF.md additions

Add these sections to build's HANDOFF.md template:

```markdown
## Maps and root
Maps: {#N, #M}, each with its Stack status comment ID. Root: `{root}`. Cross-map answers: cross-map.md.

## Seam review
| Group | Round | Last reviewed head | Verdict | Open ids | Reviewer task ID, if running | Reply file |
|---|---|---|---|---|---|---|

Deferred: N open in deferred.md. Design choices made in fix rounds: {list}.

## Sol limits
Usage-limit state, and the reset time of any pending wake-up.
```

In the in-flight table, say for each builder whether its own sol-loop is running, and which round.

## Stack status comment

```markdown
**Stack status** (edited in place by the orchestrator)
Plan: `/tmp/{slug}-plan/` on {machine}. Maps: {#N, #M}. Root: `{root}`. Merge: {pre-approved | hold}.

| # | Branch | PR | Map | Model | State | Sol | Head | Notes |
|---|---|---|---|---|---|---|---|---|
| 01 | {area}/01-slug | #N | #N | sonnet | building / sol review / open / merged | 2 rounds, approved | abc1234 | deviation |

Seam review: round N; {group}: {approved | N open}, ... Deferred: N open. Design choices made: {list}.
Wake-up: {time, or none}. Follow-ups: {layers}. Human steps owed: {list}. Report: {hatch URL}.
```
