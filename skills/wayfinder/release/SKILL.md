---
name: release
description: Cut a release or hotfix - ship dev to main with a version bump and user-facing changelog, babysit the release PR, then back-merge. Use when the user asks to cut a release, a patch or hotfix, or to merge dev into main.
agent-config-sync: true
---

# Release

A release ships `dev` to `main` on a `release/x.y.z` branch, then back-merges `main` into `dev`. Copy the repo's habits from its last release PR (`gh pr list --base main --state merged --limit 3`): where the version and changelog live, the body layout, and the merge method.

A **hotfix** is a patch release. It ships only the PRs the user names: branch from `main` and merge or cherry-pick just those.

## 1. Scope

List every PR on `dev` that is not on `main` (`git log origin/main..origin/dev --merges`). Choose the bump: the user's choice, otherwise a patch for fixes and a minor for new features.

The human steps come from the open release issue (`gh issue list --label wf:release --state open`). Run the `release-steps` skill's **Re-check** on it first. Later work may have done a step, dropped it or changed its command, and a PR may have landed without recording its steps. For a hotfix, take only the steps of the PRs it ships.

**Done when:** you have the version, the PR list, and a release issue whose open steps match today's code.

## 2. Release branch

On `release/x.y.z`, bump the version and write the changelog entry. The changelog is for **users**: one entry per change they will notice, in plain words, and none for deploys, refactors, CI or infra. Open the PR to `main` with the `pr` skill. Its body has:
- a table of each PR: what it does and its changelog entry, or `-`;
- what the deploy does;
- `Human steps: #R`, then the release issue's open steps copied under its **Before merge** and **After deploy** headings, with exact commands.

Leave closing keywords for maps out of the body. A map closes in step 4, once its prod steps are done.

**Done when:** the PR is open and every PR on `dev` is in its table.

## 3. Land it

Do every **Before merge** step first, or hand it to the user and wait. Babysit the PR with the `babysit` skill. Merge it when it is clean, unless the user said to hold. Watch the production deploy to the end. Then run the **After deploy** steps that need no human, and tick each one on the release issue. Stop at the **human gates**: an infra apply (run the plan, show what it changes, wait for approval), production data changes, and DNS or cutovers.

**Done when:** the deploy succeeded, and every post-merge step is either done or handed to the user with its command.

## 4. Back-merge

Open `chore/back-merge-x.y.z` from `main` into `dev`, babysit it, and merge it. Then finish with the `release-steps` skill's **After the release**: tick the done steps, close the release issue as `Release:x.y.z`, open the next `Release:Next` with every step still open, and set each shipped map to `released` and close it once its steps are all done.

**Done when:** `dev` contains the release commit, one `Release:Next` issue is open, and the user has a summary listing the version, the PRs shipped, and the human steps left with their release issue link.
