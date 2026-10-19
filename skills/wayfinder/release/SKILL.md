---
name: release
description: Cut a release or hotfix - ship dev to main with a version bump and user-facing changelog, babysit the release PR, then back-merge. Use when the user asks to cut a release, a patch or hotfix, or to merge dev into main.
agent-config-sync: true
---

# Release

A release ships `dev` to `main` on a `release/x.y.z` branch, then back-merges `main` into `dev`. Copy the repo's habits from its last release PR (`gh pr list --base main --state merged --limit 3`): where the version and changelog live, the body layout, and the merge method.

A **hotfix** is a patch release. It ships only the PRs the user names: branch from `main` and merge or cherry-pick just those.

## 1. Scope

List every PR on `dev` that is not on `main` (`git log origin/main..origin/dev --merges`). Choose the bump: the user's choice, otherwise a patch for fixes and a minor for new features. Note anything that needs a human: infra applies, migrations, secrets, DNS, data imports.

**Done when:** you have the version, the PR list, and the human steps.

## 2. Release branch

On `release/x.y.z`, bump the version and write the changelog entry. The changelog is for **users**: one entry per change they will notice, in plain words, and none for deploys, refactors, CI or infra. Open the PR to `main` with the `pr` skill. Its body has:
- a table of each PR: what it does and its changelog entry, or `-`;
- what the deploy does;
- the post-merge steps, with exact commands;
- `Closes #N` for each map whose slices have all merged.

**Done when:** the PR is open and every PR on `dev` is in its table.

## 3. Land it

Babysit the PR with the `babysit` skill. Merge it when it is clean, unless the user said to hold. Watch the production deploy to the end. Then run the post-merge steps that need no human. Stop at the **human gates**: an infra apply (run the plan, show what it changes, wait for approval), production data changes, and DNS or cutovers.

**Done when:** the deploy succeeded, and every post-merge step is either done or handed to the user with its command.

## 4. Back-merge

Open `chore/back-merge-x.y.z` from `main` into `dev`, babysit it, and merge it. Set each closed map's **Phase** to `released`.

**Done when:** `dev` contains the release commit, and the user has a summary listing the version, the PRs shipped, and any human steps left.
