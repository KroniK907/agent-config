---
name: grill-me
description: Stress-test a plan or design through sequential Q&A—one question per turn—walking depth-first until resolved, seeded from five canonical coverage zones. Use when the user wants to stress-test a plan, get grilled on their design, or mentions "grill me".
---

# Grill me

Explore the plan or design until you and the user share a clear picture. Work **depth-first**: choose **one** branch of the decision tree, drill into it until every detail that matters for that branch is settled, then move to the **next** branch. **Never** hop between unrelated branches in parallel; finish the current branch (or explicitly agree it is deferred) before starting another.

## Coverage zones (canonical branches)

Use these **five stack-neutral zones** so sessions do not miss whole categories of risk. They apply to web apps, backends, CLIs, workers, libraries, and **repo/meta** work (skills, rules, docs, CI)—see **Non-product lens** below.

| Zone | Reference (quick triage + deep prompts) |
|------|----------------------------------------|
| Surfaces & experience | [surfaces-and-experience.md](references/surfaces-and-experience.md) |
| Behavior & correctness | [behavior-and-correctness.md](references/behavior-and-correctness.md) |
| Boundaries & integration | [boundaries-and-integration.md](references/boundaries-and-integration.md) |
| Persistence & data | [persistence-and-data.md](references/persistence-and-data.md) |
| Change, risk & evidence | [change-risk-and-evidence.md](references/change-risk-and-evidence.md) |

### Session startup (before the first question)

1. **Mandatory lightweight pass:** Read **each** zone file only through **Quick triage** (and headings), **for every session**—no skipping files based on an early guess.
2. **Seed `Branches`:** In the **first reply**, the **Branches** list **must** include **exactly these five names** as rows (same wording as the table), each with status `in-progress`, `not started`, or `complete`.
3. **Close N/A early:** For any zone that is **not applicable**, set it to `complete` **in that first reply** (or as soon as confirmed) and add **one short reason** in **Follow Up** (what is absent or unchanged). If uncertain, keep it `not started` and ask **one** question that resolves scope—still **depth-first** (at most one row `in-progress`).
4. **Deep read:** After triage, **fully read** the reference body **only** for zones that remain plausibly in scope; use those prompts to drive questions **within** the current zone.

Add **extra** top-level rows to **Branches** only when the user asks or when the plan clearly requires a branch that is **orthogonal** to the five (e.g. a named compliance program). Do not replace the five with ad-hoc labels.

### Non-product lens (skills, rules, docs, CI)

Same five zones—**translate** them for artifact work:

- **Surfaces & experience:** Who reads this (humans, agents); clarity, structure, discoverability, examples, failure UX of instructions.
- **Behavior & correctness:** What the artifact **requires** vs **forbids**; edge cases; consistency with other rules/skills.
- **Boundaries & integration:** What tools/repos/paths this touches; how it composes with other skills or automation.
- **Persistence & data:** What is source of truth (paths, config keys); duplication vs links; versioning of templates.
- **Change, risk & evidence:** Rollout (enablement, migration), blast radius, review expectations, checks/tests/proof for the change.

## Session complete (escape from output format)

When you **reasonably believe** the grill-me session is done—every tracked branch in **Branches** is `complete` (or explicitly deferred with agreement), and you have no substantive depth-first follow-up left—you may send **one** closing reply that **does not** use the normal **Follow Up / Branches / Question / Recommendation** sections. **This is the only case** where that format is optional.

In that closing reply:

1. Give a **concise summary** of the decisions made across the session in concrete terms (what was settled, deferred, or explicitly left open).
2. Ask **exactly one** question: whether the user wants to explore **more** branches (or add new top-level branches) before closing.

Use plain prose and whatever headings or lists help readability; do not force the standard four-part template. If the user says yes or names new branches, the **next** reply resumes the full output format and interaction rules from a sensible new baseline (update **Branches** accordingly).

## Output format (required unless session complete)

Every reply **except** a **Session complete** closing turn (see above) must use this structure, in order:

**Follow Up:**
Describe the decision that was made in the last Q/A in concrete terms (what was agreed, settled, or deferred and why). If the user answered with further follow-up questions instead of directly answering your prior question, use this section to answer those user questions—do not leave them hanging. On the **first** reply in a session (no prior Q/A), briefly restate the plan or design you are stress-testing in concrete terms, then record **N/A** triage outcomes for any of the **five canonical zones** that you are closing out immediately (one short reason each). If nothing is settled yet, note that.

**Branches:**
A very short list of **top-level** decision areas. It **must** include the **five canonical coverage zones** (same names as the table), each with status `in-progress`, `not started`, or `complete`. At most **one** row should be `in-progress` (the depth-first focus). Optional additional rows are allowed per **Coverage zones**. Do **not** list sub-topics, nested bullets, or drill-down items here—only coarse branches. **Exception:** if the user explicitly asks you to track something as its own row in this list, add that branch by name and status. Update this list every turn as branches settle or focus shifts.

**Question:**
Pose **exactly one** next logical question. It must follow this skill’s interaction rules (depth-first, same branch until settled, one branch at a time). No multi-part questions, no bullet lists of questions, no trailing “also, …?” in the same section.

**Recommendation:**
Your recommendation for how to answer the **Question** above only—briefly, including why you lean that way. Do **not** use this section to recap prior Q/A, rehash **Follow Up**, or answer a different question.

## Interaction rules (non-negotiable)

1. Ask **exactly one** question per reply—in the **Question:** section only (or, on a **Session complete** closing turn only, as plain prose). No multi-part "Question 1 / Question 2", no bullet lists of questions, no trailing "also, ...?" in the same message.
2. **Wait** for the user's answer before asking the next question.
3. Use the answer to pick the next question along the **same** branch until that branch is nailed down; only then switch to a sibling or parent branch.
4. Put your suggested answer **only** in **Recommendation:**, scoped to the current **Question:**. Do not let it become a second question in the same turn.
5. If the question can be answered by exploring the codebase, **explore first**, then ask **one** follow-up only if something is still ambiguous.
6. **Coverage:** Follow **Session startup** every time: triage all five zone references, seed **Branches**, then deep-read only what is still in scope.

On normal turns, the **Question:** section must end with one clear question the user can answer in their next reply. **Session complete** closing turns have no **Question:** section; the single closing question still applies.
