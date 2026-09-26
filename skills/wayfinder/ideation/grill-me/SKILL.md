---
name: grill-me
description: grill me, stress-test plan, stress-test design, get grilled, wf:grilling, Grill: ticket, wayfinder Reconcile, depth-first Q&A, coverage zones, Surfaces & experience, layout-bearing surface, rough layout, five zones, sous-vide remainder, grilling session starts sous-vide, handed-back questions
agent-config-sync: true
---

# Grill me

[Sous-vide](../../orchestrators/sous-vide/SKILL.md) orchestrates a grilling session. If someone asks to grill a topic, plan, or design and sous-vide has not handed over an **Ask** table, follow sous-vide from the start. Do not ask the first grill question here.

Once sous-vide hands over the **Ask** table, work **depth-first**: settle every Ask row in one branch (or agree to defer it) before starting the next.

## Coverage zones

The five zones and their skip signals and prompts live in [ZONES.md](ZONES.md). Sous-vide already loaded it; do not re-read it. Load it only if it is not in context. The **Branches** labels are:

1. Surfaces & experience (UI/UX, flows, copy, a11y)
2. Behavior & correctness
3. Boundaries & integration
4. Persistence & data
5. Change, risk & evidence

## After sous-vide

The shape is settled. Ask only the rows in the bag's **Ask** table, one per turn. Do not re-ask Proposed or N/A rows, and do not re-triage zones.

A **handed-back question** is any new question an answer raises that is not an Ask id in this bag. Do not ask it. Give it a new id, set `depends-on` to the Ask id that raised it, list it under **Handed back** in **Follow Up**, and move to the next Ask row. Sous-vide scores it in the next wave, where the agent may settle it as Proposed without asking. A clarifying follow-up is allowed only when the user's answer does not settle the current Ask id itself.

A zone is `complete` when every cell in it is accepted, replaced, N/A, or answered. A zone with an open Ask row is `not started` until reached, then `in-progress`. If the user marks a shape row wrong, return to sous-vide for the cascade.

When every Ask id is answered, stop the grill format and return the answers and every handed-back question to sous-vide for the next wave. Skip the session-complete question. Reconcile waits until sous-vide finishes.

## UI layout articulation

When the plan adds or materially changes a layout-bearing surface (page/route, modal, dialog, drawer, sheet, side panel, substantial popover, each wizard step), do not mark **Surfaces & experience** `complete` until each such surface has a **rough layout** or a named deferral with an assumed default.

A rough layout gives the spatial commitments: major regions, where primary vs secondary actions sit, what scrolls vs stays fixed, and any structural choice (split view vs single column). ASCII, labeled zone lists, or short prose all work. Go one surface at a time. If the user describes several at once, reflect them in **Follow Up** and continue with one question.

Add extra **Branches** rows only when the user asks or the plan needs a branch orthogonal to the five (e.g. a named compliance program).

## Output format (every reply except session complete)

**Follow Up:** The decision from the last Q/A in concrete terms (agreed, settled, or deferred, and why). If the user asked questions instead of answering, answer them here. If the answer raised new questions, add a **Handed back** line with their ids and one-line questions. First reply: restate the plan and give one-line reasons for any zones closed as N/A.

**Branches:** The five zone labels verbatim, plus any extra top-level rows, each `in-progress`, `not started`, or `complete`. At most one `in-progress`. No sub-topics.

**Question:** Exactly one question on the current branch. No multi-part questions or lists.

**Recommendation:** Your suggested answer to that question and why. Nothing else.

## Rules

1. One question per reply. Wait for the answer.
2. Stay on the same branch until its Ask rows are answered, then move to a sibling or parent. New questions are handed back, not asked.
3. If the codebase can answer it, explore first; ask only what stays ambiguous.

## Session complete

When every branch is `complete` or deferred by agreement and no follow-up remains, send one closing reply without the four-part format: a concise summary of what was settled, deferred, or left open, then one question asking whether to explore more branches. If yes, resume the format with updated **Branches**.

On a `wf:grilling` ticket, tell the user to invoke [wayfinder](../../SKILL.md) **Reconcile** after they accept the summary. Reconcile proposes decision-log rows, tickets, and map updates; grill-me does not edit the map or log.
