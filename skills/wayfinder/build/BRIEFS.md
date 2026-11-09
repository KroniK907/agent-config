# Build briefs

Templates for the plan directory `/tmp/{slug}-plan/`. Replace `{...}` placeholders. Delegates read these files, so the prompt that launches each one stays a single line.

## Planner brief (the orchestrator sends this)

```text
You plan a stacked-PR build of {org/repo} #N ({FeatureName}:Map). The repo checkout at {path} is read-only
for you: no edits, commits or PRs. Write only under /tmp/{slug}-plan/.

Inputs: /tmp/{slug}-plan/map.md (map + binding decisions), AGENTS.md, the repo's rules files, its
architecture doc. Slices to build: {list}.

Write:
1. contract.md - the stack as ordered layers following the map's build order. A table with
   NN | branch {area}/NN-slug | base (layer below; 01 on dev) | model (sonnet = well specified or UI,
   opus = design risk) | migration file name (distinct, later than every existing one) | owns.
   Then: shared names other layers depend on (types, tables, routes, permission keys, path constants),
   parallel groups (layers that start from the same base and merge their sibling before opening a PR),
   out of scope per layer, and "Defaults chosen" for anything the decisions leave open.
   Default to one layer per map slice. Split a slice only when it passes about 3000 changed lines, not
   counting tests and generated files, or when it mixes design-risk code with UI. Every layer pays a
   fixed cost in checks, reviews and merges, so fewer layers build faster. Infra is code only;
   say what a human must apply, and keep local dev and tests working without it.
   Then add Review groups: 2 to 4 groups of layers, usually one per slice or domain. Order the layers
   so each group's layers come together and the first group closes as early as the dependencies allow.
   Each group lists its layers and the layers below it that it calls. 250 lines or fewer.
2. common.md - the rules every layer follows (outline below). Take every rule and check from the repo's
   own files, and run each check command once to confirm it exists. Find out which checks the PR's CI
   runs, and whether it runs on this machine. Find out whether the repo runs a merge train (a bot that
   merges every open, non-draft PR once its checks pass) and which label keeps a PR off it. Check whether worktrees and the CI runner share the
   build cache (for example `go env GOCACHE`). Write both into common.md.
3. build.md and babysit.md - from the templates below.
4. questions.md - only blockers where the decisions are silent and a wrong guess would be costly.
   Leave it empty if there are none.

Report in 8 lines or fewer: layer count, `NN-slug(model)` in order with parallel groups marked,
and whether questions.md is empty.
```

## common.md outline

- **Context:** one paragraph on the feature and the stack, plus a pointer to `contract.md` and `map.md`. The decisions are binding; the contract beats Sol findings.
- **Read first:** AGENTS.md, every repo rules file, the architecture doc.
- **Git:**
  - Work only in the given worktree, and leave `git stash` alone.
  - Use the commit identity and `Co-Authored-By` trailer for the model the delegate actually runs.
  - Use merge commits. Never rebase or force-push a branch that has a PR.
  - Merge the base branch right before every push. Regenerate generated files after merge conflicts, for example migration hash files.
- **Code rules:**
  - The repo's own rules, summarized, with each one pointing at its file.
  - No planning ids (`GM-`, issue numbers) in code, tests or user-visible text.
  - Only the libraries the decisions allow.
- **Checks:** the exact commands for build, vet, format, unit tests, integration tests and asset builds. Split them three ways so each push runs the full suite once, in CI:
  - **Fast, before every push:** build, vet, format, and unit tests for the packages the layer touches.
  - **Own tests:** only the layer's new or changed integration tests, narrowed with the runner's filter (for example `-run`). Prove each new one runs by making it fail once.
  - **Full suites:** left to the PR's CI checks. Wait for them with `pr-wait --until green --no-reviewers` from the `babysit` skill, never a polling loop of your own. Run one locally only when CI doesn't run it, or to reproduce a red check.
  - **Foreground only:** run every check in the foreground and read its result before you stop. If a lock or queue makes a check wait, wait in the foreground with a long timeout. Never end your turn with a check still running in the background.
  - **Flakes:** the machine is shared. Rerun a timeout, or a failing test the layer doesn't touch, once and alone before you call it a failure. Report it either way.
  - For UI layers: run the app against a disposable database and look at every new state at desktop and phone widths.
- **Hard limits:**
  - Leave production data, infra applies and DNS to a human. Write what they must run into the PR body.
  - Never merge.
- **PR:**
  - Open no PR, not even a draft, until Sol approved the layer and the fast checks and own tests pass on the pushed head. This overrides the `pull-request` rule's draft-on-first-edit step.
  - Then open it ready, not draft, against the contract's base.
  - On a merge-train repo, open it with the train's opt-out label (`gh pr create --label <label>`), and say so here with the label's name.
  - Write the body with the `pr` skill. It includes the line `Part of #N; base: <base>`, with no closing keywords.
- **Report:**
  - Post `PUSHED <branch> <sha>` as an interim message as soon as the branch is first pushed.
  - Report back in 10 lines or fewer: branch, PR URL once one exists, head SHA, check results on one line, and one line per deviation from the contract.
  - Append `notes-NN.md` with the details, the names later layers need, and **Screenshot states**. List those one per line as `URL - how to reach it (role, seed data, clicks) - desktop/phone`, or write `none`.
- **Known quirks:** environment workarounds the repo or earlier builds hit, including known flaky tests and whether the build cache is shared.

## build.md

```markdown
You build ONE layer of the {slug} stack. Your layer, branch, base and worktree are in your task.
**Never end your turn while a command runs.** Run every check in the foreground with a long timeout and read
its result. Don't run the full suite ({full suite command}) before the PR. CI runs it.
1. Read common.md and contract.md (build only your row and its section; keep the shared names).
   Read earlier layers' notes-*.md for names and deviations.
2. Create your worktree: `git -C {repo} fetch origin && git -C {repo} worktree add {worktree} -b {branch} origin/{base}`
   (reuse it if it exists).
3. Build, run every check in common.md, commit, and push the branch. Open no PR; the babysitter opens it
   after Sol approves. Do not babysit or merge.
4. Write notes-NN.md and report per common.md.
```

## babysit.md

```markdown
You babysit ONE layer of the {slug} stack. Branch, base and worktree are in your task. The branch is pushed
and has no PR yet.
**Never end your turn while a command runs, or before your final report.** Run checks in the foreground. Start
each Sol round with `delegate_task` `mode: "wait"`, as sol-loop says for delegated runs.
Read common.md, then follow the babysit skill with these overrides:
- MERGE HOLD: never merge. Stop at "Sol approved + PR open + checks green".
- Run the sol-loop on the pushed branch. Open the PR as common.md says only after Sol approves and the fast
  checks and own tests pass on the pushed head.
- Fix findings yourself. Start only the one Sol reviewer per round that sol-loop calls for, and nothing else.
- sol-loop base: the merge-base with origin/{your base}. Scope: this layer's commits only. Intent: your
  contract.md row. Add to the brief: "Behavior that matches a GM decision in map.md or contract.md is not a bug.
  Out of scope: {other layers}."
- contract.md and map.md beat Sol. Decline findings that contradict them, citing the decision.
- A finding that needs a change in a LOWER layer: report it rather than fixing it here.
- No route reaches Sol: report `SOL UNAVAILABLE` with what you tried. Do not swap in another model.
- Merge the base branch before every push, so one push carries both. Rerun the fast checks in common.md first.
  Leave full suites to CI and wait for them with `pr-wait --until green --no-reviewers` from the `babysit` skill.
- Parallel-group layer: merge your sibling branch before you open the PR. If the sibling isn't pushed with Sol
  approval yet, report `WAITING ON NN` with your branch pushed and no PR. Don't wait for it.
- Log each sol-loop round (id -> fixed, disagreed or declined, and why) in notes-NN.md.
Report in 8 lines or fewer: PR URL, Sol rounds and verdict, final head SHA, check state, anything for a lower layer.
```

## finish-NN.md (parallel-group finisher)

The orchestrator writes one when a waiting layer's sibling opens its PR. The prompt is one line:
`Follow /tmp/{slug}-plan/finish-NN.md exactly.`

```markdown
You finish layer NN ({NN-slug}) of the {slug} stack. Its earlier agent stopped on purpose, waiting for layer MM.
State: branch {branch} is pushed at {sha}, Sol approved it in {N} rounds (ledger in notes-NN.md), no PR yet.
Worktree: {worktree} (exists). Layer MM is now pushed with Sol approval: origin/{sibling branch} at {sha}, PR #{N}.
**Never end your turn while a command runs, or before your final report.** Sol rounds use `mode: "wait"`.
Read common.md, contract.md row NN and its parallel group, and notes-NN.md and notes-MM.md. Then:
1. `git fetch origin`, check out {branch} at origin's head, and merge origin/{sibling branch} with a merge commit.
   Keep both sides in conflicts and regenerate generated files. Fix builds the merge broke without changing
   the sibling's behavior.
2. Rerun the fast checks and your own tests in the foreground.
3. If you resolved a conflict by hand in behavior code, or reworked code, run one sol-loop round on only those
   files: diff since {sha}, excluding the sibling's commits. Otherwise run no new round.
4. Push, open the PR against {PR base} as common.md says, and wait for CI with `pr-wait --until green --no-reviewers`.
   Fix a red check.
5. Append the merge to notes-NN.md and report per common.md.
```

## Group review brief (round 1, one file per group and reviewer)

```text
Read-only review of the {group} group of the {FeatureName} stack: layers {list of PRs/branches},
top {branch} at {head}. Base: {merge-base with origin/dev}. Check out the top head detached in your
worktree. No commits, pushes, PR comments, infra or database writes. Local builds, tests and read-only
requests to {dev URLs} are fine. Don't run full test suites; CI already ran them on every head.
Read AGENTS.md, the repo rules and map.md. Its GM decisions are binding, so behavior that matches
a decision is not a bug. Out of scope: {list}.
Each PR already passed its own review. Look at how this group's layers work together and with the
layers below them: calls, payloads, permissions, migrations and races that cross a PR boundary.
Security first: authz gaps, data leaks across tenants, injection or XSS, secrets in client code, CSRF.
Then correctness, error handling, broken end-to-end flows, migrations that fail on real data, and
front end versus back end mismatches.
Each finding: severity, file:line, a concrete failure scenario, a suggested fix, the owning layer.
Verify each against the code. Write the findings to /tmp/{slug}-plan/review-rN-{group}-{reviewer}.md
and return them. If nothing holds up, say NO FINDINGS.
```

## Group re-review brief (rounds 2 and later)

```text
Same rules and output as round 1. Group: {group}. You last reviewed {previous head}; the group is now
synced at {head}. Review only git diff {previous head} {head}, limited to the fixes for your open
findings and the code those fixes touch. Prior findings and resolutions: {list}. Mark each one
resolved or still open, with the reason. Raise a new finding only if a fix caused it. Do not review
the full group again.
```

The final pass uses the round 1 brief with the whole stack as the group, scoped to the seams that join groups.

## Stack status comment

```markdown
**Stack status** (edited in place by the orchestrator)
Plan: `/tmp/{slug}-plan/` on {machine}. Merge: {pre-approved | hold}.

| # | Branch | PR | Model | State | Head | Notes |
|---|---|---|---|---|---|---|
| 01 | {area}/01-slug | #N | sonnet | Sol approved / babysitting / building / merged | abc1234 | deviation |

Stack review: round N, {open findings | NO FINDINGS}. Human steps owed: {list}. Report: {hatch URL}.
```

## resume.md

```markdown
A previous agent on this layer stopped mid-work ({reason}, {time}). First establish the state:
`git fetch`, your branch and its open PR, its checks, and notes-NN.md (including the sol-loop ledger).
Don't redo finished work. Check out the existing branch, then continue your brief.
```
