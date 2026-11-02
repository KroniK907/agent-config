---
name: handoff
description: Hand a long-running orchestrator thread to a fresh thread that picks up from files, while the old thread relays its delegates' reports. Use when the user asks for a handoff or a fresh thread, or when a skill or a thread's context says it has run too long.
agent-config-sync: true
---

# Handoff

Every turn re-reads the whole thread. A thread that has run for hours gets slow and expensive, and its early instructions fade. A handoff moves the work to a fresh thread that starts cold and learns the state from files.

A calling skill (`build`, `sol-build`) supplies four things. Without one, work them out with the user:
- **State dir:** where the work's files live, such as `/tmp/{slug}-plan/`.
- **Status record:** the durable record outside this thread, such as a tracker comment edited in place.
- **Launch command:** the command that started the work, so the successor loads the same skill.
- **Template additions:** extra sections for HANDOFF.md.

## When

Hand off at the next break (a group merged, a wave dispatched, a sync done) once any of these is true:
- the harness has compacted this thread's context;
- this thread has taken about 20 delegate reports;
- this thread has run about 4 hours;
- the user asks for it.

Never fork this thread with `t3_thread_transfer`. A fork copies the context you are trying to drop.

## Steps

1. **Settle.** Start no new delegates. Bring the status record and the state dir's task list up to date.
2. **Collect scripts.** Move every script and helper file this thread wrote outside the state dir into it: status pushers, check wrappers, lock files, brief drafts. Fix the paths inside them and in the briefs that call them. The successor finds only what the state dir holds.
3. **Write `HANDOFF.md`** in the state dir from [HANDOFF.md](HANDOFF.md), plus the caller's additions. Copy the original command and every standing instruction the user gave since, word for word. A paraphrase drops merge approvals and gates.
4. **Release this thread's bindings.** Delete each `schedule_task` bound to this thread (`list_scheduled_tasks`, then `delete_scheduled_task`), usage-limit wake-ups included. Call `unwatch_pull_request` on each PR it watches. HANDOFF.md lists both so the successor recreates them. Leave running delegates alone, because cancelling one throws away its work.
5. **Launch the successor** with `t3_thread_launch`: this project, no `workspaceStrategy`, your model and runtime mode, a title naming the work and the next wave. Start the message with the launch command, so the skill loads again:
   ```text
   {launch command} Take over as orchestrator of {the work}, wave N+1. Read {state dir}/HANDOFF.md first. Recreate the schedules and PR watches it lists, bound to your thread. Work it marks "old thread" stays with its delegates until the old thread relays their reports. Carry on until the work is done.
   ```
6. **Record it.** Put the successor's thread link in the status record and at the top of HANDOFF.md.
7. **Relay.** A delegate reports to the thread that started it, so its result still lands here. From now on this thread only relays. When a delegate reports, send it to the successor with `t3_thread_send` in 3 lines or fewer: what it worked on, PR, head, verdict, checks, notes file. Start nothing, fix nothing, and leave the status record alone. If the successor asks you to cancel a delegate, `task_cancel` it. After the last report, send a final line saying so, and stop.

Outside T3 Code, do steps 1 to 3, then give the user the launch message to paste into a new session.

**Done when:** the successor's link is in the status record, this thread holds no schedules or PR watches, and every in-flight delegate has reported and been relayed.

## As the successor

Check before you act. Read HANDOFF.md, the status record, the open PRs and the notes files, and confirm each claim against GitHub. Work marked "old thread" belongs to its delegate. Don't sync, fix or push to its branch until the relay reports it done. If its PR and notes show it finished and its worktree has sat idle, ask the old thread to cancel the delegate first. Two agents pushing to one branch collide. Start your own section in the task list, and hand off again at the same triggers.
