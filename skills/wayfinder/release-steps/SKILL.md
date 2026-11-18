---
name: release-steps
description: Record the human steps that merged work needs on production, such as infra applies, secrets, Scheduler jobs, data imports, DNS or settings, on the repo's open release issue, and re-check that list against the code. Use when a build, feature branch or PR lands on dev with steps left for a human, or when asked to list or check what the next prod release needs.
agent-config-sync: true
---

# Release steps

Work merges to `dev` weeks before it reaches production. The **human steps** it needs on prod must outlast the build thread that knew about them: a tofu apply, a secret to create, a Scheduler job, a data import, a DNS record, a setting to fill in after deploy. They all go on one issue, the **release issue**, and the `release` skill works from it.

## The release issue

There is one open release issue per repo, labelled `wf:release` and titled `Release:Next`. Find it with `gh issue list --label wf:release --state open`. If none is open, create it from [TEMPLATE.md](TEMPLATE.md) (`gh issue create --label wf:release`). If the label is missing, create it first with `gh label create wf:release --description "Human steps owed before the next prod release"`.

Each human step is one checkbox line, in the section for when it must happen: **Before merge**, **After deploy**, **Verify on prod**, or **Any time** for leftovers on work already on prod. A step the code can't support yet goes under **Needs code before it can be done**. A step that blocks the release goes first in its section and says so.

```markdown
- [ ] **infra** Turn on the jobs tick for prod. Add "prod" to `jobs_tick_envs` in tfvars, then `tofu plan` and `tofu apply` from the main checkout. From #PR ({FeatureName}:Map #N). Dev: applied 2026-11-14. Check: `gcloud scheduler jobs describe nus-jobs-tick-prod`.
```

- **Kind** comes first, in bold: `infra`, `secret`, `scheduler`, `maps-vm`, `data`, `dns`, `third-party`, `setting`, `verify`.
- **Action:** the exact command, or the exact page and field. Never a secret value.
- **From:** the PRs and maps that need it. A step that several PRs need is one line naming them all.
- **Dev:** done with a date, not needed, or owed. A step done on dev only stays unticked.
- **Check:** how a human confirms it is done on prod.

Tick a box only when the step is done **on prod**. A step that later work made unnecessary moves to **Dropped** with the commit or file that made it so.

## 1. Collect

Gather the human steps for the work that just landed. Read, in this order:

1. The **Human steps** section of each merged PR body, and the Stack status comment's "Human steps owed".
2. The map's build-order decision, which names the human steps the grill expected.
3. The diff itself, for steps nobody wrote down. Compare the merged work against `main`: new or changed files under `infra/`, new secrets or env reads, Scheduler or cron targets, maps VM or GeoServer files, migrations that need data prepared first, third-party setup (webhooks, OAuth clients, Workspace delegation, mailboxes), and settings that start empty.

Check each step against the code as it is now, not as the PR described it. A later PR may have renamed the variable, changed the command, or done the step already. Read live state where you can do so read-only (`gcloud ... list` or `describe`). Never apply anything.

**Done when:** every step has a kind, an action, its PRs and maps, its dev state and a check, and every file the diff adds under `infra/` is either covered by a step or confirmed to need none.

## 2. Write

Other threads edit this issue too, so fetch the body right before you change it: `gh issue view R --json body -q .body > /tmp/release-R.md`. Add each step to its section. Where an existing step covers the same thing, add your PRs to it rather than writing a second line. Update the **Maps waiting** line. Upload with `gh issue edit R --body-file /tmp/release-R.md`, then read it back to confirm your lines are there.

Then point the work at it:

- Each map whose work landed: set **Phase** to `built` with `Next: ships in the next release; prod steps on Release:Next #R` (see the map format in `../MAP.md`).
- The Stack status comment: `Human steps owed: on #R`.
- Each PR that owes a step: one comment linking the release issue.

**Done when:** every step from step 1 is on the release issue, and every map and PR that owes a step links to it.

## 3. Re-check

Run this before a release, and whenever someone asks what prod still needs. For every unticked step, check the current `dev` code and live prod state again:

- Done on prod: tick it and add the evidence.
- No longer needed: move it to **Dropped** with the commit or file that made it so.
- Changed: rewrite the action to match the code now.

Then compare `main` with `dev` for steps nobody recorded, as step 1 does, and add them.

**Done when:** every unticked step was checked against today's code, and the diff from `main` to `dev` has no human step missing from the issue.

## After the release

The `release` skill copies the open steps into the release PR body and runs the ones that need no human. When the release has shipped:

1. Tick what was done. Rename the issue to `Release:x.y.z` and close it.
2. Open a new `Release:Next` that carries over every step still unticked. A hotfix ships only some PRs, so it carries over every step those PRs do not cover.
3. Set each shipped map to `released`. Close it if every step it owes is ticked; otherwise its Next names the open steps.

**Done when:** exactly one release issue is open, and every map the release shipped is `released`, closed unless a step it owes is still open on the new issue.
