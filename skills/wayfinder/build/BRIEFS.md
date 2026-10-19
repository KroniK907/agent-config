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
   Size each layer for one agent and one PR (about 1500 changed lines or fewer). Infra is code only;
   say what a human must apply, and keep local dev and tests working without it. 250 lines or fewer.
2. common.md - the rules every layer follows (outline below). Take every rule and check from the repo's
   own files, and run each check command once to confirm it exists.
3. build.md and babysit.md - from the templates below.
4. questions.md - only blockers where the decisions are silent and a wrong guess would be costly.
   Leave it empty if there are none.

Report in 8 lines or fewer: layer count, `NN-slug(model)` in order with parallel groups marked,
and whether questions.md is empty.
```

## common.md outline

- **Context:** one paragraph on the feature and the stack, plus a pointer to `contract.md` and `map.md`. The decisions are binding; the contract beats reviewbot.
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
- **Checks:** the exact commands for build, vet, format, unit tests, integration tests and asset builds.
  - For UI layers: run the app against a disposable database and look at every new state at desktop and phone widths.
  - Prove each new integration test runs by making it fail once.
- **Hard limits:**
  - Leave production data, infra applies and DNS to a human. Write what they must run into the PR body.
  - Never merge.
- **PR:**
  - Open it ready, not draft, against the contract's base.
  - Write the body with the `pr` skill. It includes the line `Part of #N; base: <base>`, with no closing keywords.
- **Report:**
  - Report back in 10 lines or fewer: PR URL, head SHA, check results on one line, and one line per deviation from the contract.
  - Append `notes-NN.md` with the details, the names later layers need, and **Screenshot states**. List those one per line as `URL - how to reach it (role, seed data, clicks) - desktop/phone`, or write `none`.
- **Known quirks:** environment workarounds the repo or earlier builds hit.

## build.md

```markdown
You build ONE layer of the {slug} stack. Your layer, branch, base and worktree are in your task.
1. Read common.md and contract.md (build only your row and its section; keep the shared names).
   Read earlier layers' notes-*.md for names and deviations.
2. Create your worktree: `git -C {repo} fetch origin && git -C {repo} worktree add {worktree} -b {branch} origin/{base}`
   (reuse it if it exists). Parallel-group layers merge their sibling branch before opening the PR.
3. Build, run every check in common.md, commit, push, and open the PR. Do not babysit or merge.
4. Write notes-NN.md and report per common.md.
```

## babysit.md

```markdown
You babysit ONE PR of the {slug} stack. PR number, branch and worktree are in your task.
Read common.md, then follow the babysit skill with these overrides:
- MERGE HOLD: never merge. Stop at "review clean + checks green".
- Fix findings yourself. Spawn no subagents.
- Read every reviewbot review on the PR, not just the latest. An older review can hold an open major finding.
- No review after about 27 minutes: push one empty commit (`git commit --allow-empty -m "Request review"`) and
  wait again. If there is still nothing, report and stop.
- contract.md and map.md beat reviewbot. Set aside findings that contradict them, and reply on the comment.
- A finding that needs a change in a LOWER layer: report it rather than fixing it here.
- Merge the base branch before every push, so one push carries both. Rerun the checks in common.md first.
- At most 5 rounds. Log each round (finding -> fixed or set aside, and why) in notes-NN.md.
Report in 8 lines or fewer: rounds, final head SHA, review state, check state, anything for a lower layer.
```

## Stack review brief (one file per round and reviewer)

```text
Read-only review of the {FeatureName} stack: {list of PRs/branches, top branch}. Check out the top
branch detached in your worktree. No commits, pushes, PR comments, infra or database writes.
Local builds, tests and read-only requests to {dev URLs} are fine.
Read AGENTS.md, the repo rules and map.md. Its GM decisions are binding, so behavior that matches
a decision is not a bug. Out of scope: {list}.
Look for real defects, security first: authz gaps, data leaks across tenants, injection or XSS,
secrets in client code, CSRF. Then correctness, races, error handling, broken end-to-end flows,
migrations that fail on real data, and front end versus back end mismatches.
Prior findings and their resolutions: {list, or "none, round 1"}. Re-raise one only if its fix is wrong.
Each finding: severity, file:line, a concrete failure scenario, a suggested fix. Verify each against
the code. Write the findings to /tmp/{slug}-plan/review-rN-{reviewer}.md and return them.
If nothing holds up, say NO FINDINGS.
```

## Stack status comment

```markdown
**Stack status** (edited in place by the orchestrator)
Plan: `/tmp/{slug}-plan/` on {machine}. Merge: {pre-approved | hold}.

| # | Branch | PR | Model | State | Head | Notes |
|---|---|---|---|---|---|---|
| 01 | {area}/01-slug | #N | sonnet | review clean / babysitting / building / merged | abc1234 | deviation |

Stack review: round N, {open findings | NO FINDINGS}. Human steps owed: {list}. Report: {hatch URL}.
```

## resume.md

```markdown
A previous agent on this layer stopped mid-work ({reason}, {time}). First establish the state:
`git fetch`, your branch and its open PR, the latest reviewbot review and checks, and notes-NN.md.
Don't redo finished work. Check out the existing branch, then continue your brief.
```
