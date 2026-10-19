---
name: code-review
description: Review a PR, branch, or diff since a ref on two axes - standards and spec - and report or post the findings. Use when asked to code review a PR, all open PRs, a branch, or work in progress.
agent-config-sync: true
---

# Code review

Review the diff between a **fixed point** and the head on two separate axes:

- **Standards:** the repo's documented rules, plus the smell baseline below where the repo says nothing.
- **Spec:** whether the diff does what its issue or map decisions ask, no more and no less.

Report the two axes separately. A diff can build the wrong thing well, or the right thing badly, and merging the axes hides one behind the other.

## 1. Pin the diff

- **For a PR:** the fixed point is its base. Use `gh pr diff <N>`, or check it out and run `git diff origin/<base>...HEAD`.
- **For a branch or ref:** use the one the user named, and ask if they named none.
- **For several PRs:** run one review per PR.

**Done when:** the diff is non-empty and you know its commits.

## 2. Find the spec and the standards

- **Spec:** issues the PR body or commits name ("Part of #N", "Closes #N"). For a map, read the slice and the GM decisions it cites. If none of these exist, review Standards only and say so.
- **Standards:** `AGENTS.md`, `CLAUDE.md`, the repo rules files (`.claude/rules/`), and contributing docs.

**Done when:** you have a spec source (or a stated reason there is none) and a list of standards files.

## 3. Review

For a small diff, review it yourself. Otherwise give each axis a subagent at **build** tier on Sonnet 5.5 (see the `subagent` skill). Run them in the foreground and wait for both. With several PRs, launch one pair per PR in a single batch.

Each brief carries:
- the repo path and the diff command;
- the standards files or the spec text;
- for Standards: report each violation with `file:line` and the rule it breaks, and each baseline smell with the hunk quoted. Separate hard violations from judgement calls, and skip anything a linter already enforces;
- for Spec: report missing or partial requirements, scope creep, and wrong behavior, quoting the spec line for each;
- verify every finding against the code, and keep the report under 400 words.

**Smell baseline** (Fowler): Mysterious Name, Duplicated Code, Feature Envy, Data Clumps, Primitive Obsession, Repeated Switches, Shotgun Surgery, Divergent Change, Speculative Generality, Message Chains, Middle Man, Refused Bequest.

**Done when:** both axes have a report.

## 4. Report

```markdown
## Code review - <fixed point>...<head>
**Spec:** <#N / none>
## Standards
- **[major|minor|nit]** `file:line` - finding
## Spec
- ...
**Summary:** Standards N (worst: ...) - Spec M (worst: ...)
```

When the user asks for findings on the PR, post this as a PR comment (`gh pr comment <N> --body-file`). Otherwise reply in chat. Fix nothing unless asked.

**Done when:** the report is posted or in your reply.
