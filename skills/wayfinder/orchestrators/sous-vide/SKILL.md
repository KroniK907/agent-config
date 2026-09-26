---
name: sous-vide
description: sous-vide, wayfinder orchestrator, pre-grill a plan, sign off proposed decisions, grill-me the ask list, next wave after grill-me
agent-config-sync: true
---

# Sous-vide

Orchestrate a plan into a shape. Score one wave, the user settles **Proposed**, [grill-me](../../ideation/grill-me/SKILL.md) asks **Ask**, then the next wave starts. Repeat until a wave has no **Ask** rows and nothing left to score.

Use this when a plan is concrete enough to cut into scenarios, including when the user asks for a grilling session. Use [strategic-ideation](../../ideation/strategic-ideation/SKILL.md) when the scope is still unstable. If the plan is too thin to slice, stop and point them there.

This skill does not edit a map, a decision log, or code. The only reads for triage are the five grill-me zone files. Scoring stays in this session. Wayfinder **Reconcile** runs after the loop finishes, and only on a `wf:grilling` ticket.

## Loop

### 1. Take the plan

On the first wave, restate the thing being shaped in one sentence and name the source. A correction restarts the loop from slicing.

**Done when:** That sentence is in the reply, or a later wave is continuing from settled rows.

### 2. Cut slices

A **slice** is one scenario or one artifact. Re-score that slice when an answer changes. Put each in-scope behavior in one slice. Keep the same slice names on later waves.

**Done when:** Every slice has a name, and each behavior has one owning slice.

### 3. Triage cells

A **cell** is one slice crossed with one grill-me zone. On the first wave, read only **Quick triage** in each zone file, for every slice:

- [Surfaces & experience](../../ideation/grill-me/references/surfaces-and-experience.md)
- [Behavior & correctness](../../ideation/grill-me/references/behavior-and-correctness.md)
- [Boundaries & integration](../../ideation/grill-me/references/boundaries-and-integration.md)
- [Persistence & data](../../ideation/grill-me/references/persistence-and-data.md)
- [Change, risk & evidence](../../ideation/grill-me/references/change-risk-and-evidence.md)

A cell is N/A only with one reason from that file's skip signals. If the plan does not state the absence, the cell stays in scope. Later waves reuse this triage.

**Done when:** Every cell is `in scope` or `N/A` plus that reason.

### 4. Write questions for this wave

A **wave** is every still-open question whose `depends-on` ids are settled. Settled means accepted, replaced, answered in grill-me, or N/A. Leave `depends-on` empty when no earlier answer can change this one.

Write the questions before you score them. Each has two to four answers and a stable id. A layout-bearing surface follows grill-me **UI layout articulation** and is Ask. Read a zone's deep prompts only to invent questions for in-scope cells.

**Done when:** Every unblocked in-scope cell has a question, or a one-line note that it has no decision of its own. Questions that still depend on an open row wait for a later wave.

### 5. Score the wave

Pick the option that fits the plan and the settled rows, or `other` if none fit. Then pick a closedness level, low to high: `several real designs`, `one lean and a real alternative`, `one obvious answer`.

**Proposed** only when closedness is `one obvious answer` and the decision names the fact that rules the others out. Anything else is **Ask**, including `other`, a layout, and a question that still means writing something new.

**Done when:** Every question in this wave is Proposed or Ask.

### 6. Show the bag

Reply in this order, then wait.

**Shape.** One short paragraph of what is already settled.

**Slices.** Names only.

**Proposed.** A table with `id` and `decision`. The decision is one statement that states the question and the answer together. A reader needs no other row to understand it.

**Ask.** A table with `id` and `question`. No why, no options. The user does not answer these here.

**N/A.** A table with `id`, `question`, and `why`.

| Reply | Effect |
| --- | --- |
| Accept | Proposed rows join the shape |
| Accept except `<id>` | Those ids move to Ask. The rest join the shape |
| Not the shape: `<id>` | Clear that id and its cascade, then re-score |
| Answer `<id>`: `<text>` | That text is settled. Re-score dependents |

A **cascade** is the changed row plus every row in that slice whose `depends-on` chain includes it. List the cleared ids and their old answers, then the new Proposed and Ask entries. Other slices stay as they are.

**Done when:** Every Proposed row is accepted, replaced, or moved to Ask.

### 7. Grill the ask list

When the bag is settled and **Ask** has rows, follow [grill-me](../../ideation/grill-me/SKILL.md) in this same session. Pass the shape as settled and the Ask table as the only open work. Grill-me asks those questions one at a time. When every Ask id from this bag is answered, those answers are settled.

**Done when:** Every Ask id from this bag has an answer, or **Ask** was empty.

### 8. Next wave or finish

If any open question is now unblocked, go back to step 4 and show a new bag. If none are, the loop is finished. On a `wf:grilling` ticket, tell the user to invoke wayfinder **Reconcile**.

**Done when:** The reply is the next bag, or the loop is finished and Reconcile is named when a grilling ticket is in play.
