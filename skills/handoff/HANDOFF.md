# HANDOFF.md template

Write it fresh at each handoff. Fill every section, or write `none`. Add the calling skill's sections after **Setup that works**.

```markdown
# {Work name}: orchestrator handoff, wave N+1 ({date, time, zone})

Previous orchestrator: {old thread link}. Successor: {new thread link, added after launch}.
You take over as orchestrator. Follow the `{skill}` skill. This file, the state dir and the
status record hold the state. Check each claim against GitHub before you act on it.

## Original command
{the command that started the work, word for word}

## Standing instructions
{everything the user has said since, word for word, each with its date: merge approvals, gates,
what to ask before, what never to do, recurring checks}

## Setup that works
- State dir: {path}, with what each file holds when the name doesn't say.
- Status record: {where it is and how to edit it}.
- Scripts: `path` - what it does - how to call it. One line each.
- Delegates: the one-line prompt, the model and effort per kind of task, the reviewer model and options.
- Merge steps as they actually ran, with workarounds.
- Traps: failures that came back more than once, and their fixes.

## Done
PRs and follow-ups merged. Reports served.

## In flight (old thread)
| Task | Role | Task ID | PR | Last known state | How to tell it finished |
|---|---|---|---|---|---|

## Recreate in your thread
- Schedule: every {interval}. Prompt: {word for word}
- Watch: PR #N

## Next
What to start, in order, with the condition each waits on. Open questions for the user. Human steps owed.
```
