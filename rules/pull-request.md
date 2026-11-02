---
description: As soon as this session starts changing files in a git repo, open a pull request to track the work.
agent-config-sync: true
---

# Pull request as soon as you edit

When this session starts changing files in a git repo, open a pull request right after the first edit. Do not wait until the work is finished. Push later changes to the same branch so the pull request tracks the work as it goes. Questions that do not edit files stop with no pull request.

**Done when:** `gh pr view --json url,baseRefName` returns a URL whose `baseRefName` is the target branch, or the session changed no files. On a merge-train repo, the branch is pushed and the PR exists only once Sol approved the work.

1. Pick the target branch. Use `dev` if it exists on the remote, then `staging`, then `main`. Do not ask. The user will correct the target if needed.
2. Stay in the current worktree when its branch is not the target branch. A host worktree counts. Create a new worktree only when the checkout is the target branch itself. Do not commit on the target branch.
3. After the first edit, commit the files this session changed and push. Open the pull request with `gh pr create --draft --base <target>` when `gh pr view` does not already show one. Commit and push again after each further round of edits.
4. If `git diff origin/<target>...HEAD` is empty, stop. There is nothing to review.
5. When the work is done, mark the pull request ready with `gh pr ready`.

## Merge-train repos

Some repos run a merge train: a bot that merges every open, non-draft PR once its checks pass. The repo's AGENTS.md or CI docs say so. There, a PR opened early can merge before anyone reviews it. Change the steps above:

- After each round of edits, commit and push the branch. Open no PR yet, not even a draft.
- When the work is done, run the `sol-loop` skill on the branch, or the `babysit` skill, which runs it.
- Once Sol approves and the repo's fast checks pass on the pushed head, open the PR ready with `gh pr create --base <target>`. The train takes it from there.
- If the loop ends stuck or the user stops the work, leave the branch pushed with no PR and say so.

Skip the pull request only when the user says not to open one.
