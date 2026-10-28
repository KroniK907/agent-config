---
name: grill
description: Grill a feature map until every open question is a recorded decision and the build order is written. Use when the user asks to grill, sous-vide, settle decisions, or get a map ready to implement.
agent-config-sync: true
---

# Grill

A grilling session settles a map in waves. Each wave starts with a **bag**: the decisions you can propose, sorted from the questions only the user can answer. The user accepts the proposals. You then ask the rest one at a time and repeat until nothing is open. The session ends by writing everything to the map.

The map format is in the wayfinder folder's `MAP.md` (resolve this skill's folder with `readlink -f`; `MAP.md` is in its parent). Coverage zones are in [ZONES.md](ZONES.md). Read both once per session.

## 1. Load

- **With `#N`:** read the map's body and comments (`gh issue view N --comments`), the maps it links, and the code its questions touch.
- **Without a map:** grill from the user's description and file the map at the end.

Answer whatever the codebase can answer before asking.

**Done when:** you can state in one sentence what is being shaped, and its source.

## 2. Show the bag

1. **Areas.** Cut the work into named areas, one per scenario or artifact, with each behavior in exactly one area. Reuse the map's Open questions groups when it has them. Keep the names fixed across waves.
2. **Triage.** Cross each area with each zone in ZONES.md. A cell is N/A only with that zone's **Skip when** reason.
3. **Questions.** Write a question for each in-scope cell whose dependencies are settled. Give each one 2-4 answers and a stable id. Handed-back questions keep their ids and are scored first.
4. **Score.** A question is **Proposed** only when it has one obvious answer and you can name the fact that rules out the others. Everything else is **Ask**, including any layout and anything that means designing something new.

Reply with: **Shape** (one paragraph of what is settled, with sources), **Areas** (names only), **Proposed** (`id | decision`, each row a statement that reads on its own), **Ask** (`id | question`, no options), **N/A** (`id | question | why`). Then wait.

Read the reply for intent. "Accept" settles every Proposed row. "Accept except P3" moves P3 to Ask. "P3 is wrong" clears P3 and every row that depends on it, then re-scores them. An answer to an id settles it.

**Done when:** every Proposed row is accepted, replaced, or moved to Ask.

## 3. Ask one at a time

Work depth-first: finish one zone's Ask rows before moving to the next. Every reply uses exactly this format:

- **Follow Up:** the decision from the last answer, in concrete terms. Answer any question the user asked. List new questions under **Handed back** as `H#`, each with a one-line question and the id it depends on.
- **Question (id):** exactly one question.
- **Recommendation:** your answer and why. Say what the alternatives cost.

New questions an answer raises are handed back rather than asked. Ask a clarifying follow-up only when the answer leaves the current id unsettled.

**Layouts.** Every new or changed page, modal, drawer, sheet or panel needs a rough layout before the **Surfaces & experience** zone is complete. Use an ASCII layout for desktop and for phone, one surface at a time. When the user wants to see it, or two layouts are close, have a design subagent build a static HTML mock in a temp directory and serve it with `hatch`. Iterate until the user picks one. When they do, record the layout in the decision, and commit the mock if the repo keeps design mocks.

**Done when:** every Ask id in this bag is answered. Then go back to step 2 with the handed-back and newly unblocked questions. When a wave has nothing left to ask, go to step 4.

## Moves the user can call at any point

- **Pause for research:** follow the `research` skill, post the findings on the map, and resume with the same question.
- **Park an idea:** file a plain issue for a feature idea that is left out of this version, and link it from **Not yet specified**.
- **New map:** when a topic deserves its own map, create it with the `map` skill and link both maps.
- **Revise:** a changed answer replaces its decision. Re-score everything that depends on it.
- **Pause:** write what is settled so far (step 5), and say where the next session resumes.

## 4. Build order

The last decision is always **build order and evidence**. It covers the slices in build order (each one agent and one PR), what can run in parallel, the tests that prove each slice, and human steps such as infra applies or cutovers. Propose it, then get it accepted.

**Done when:** the user accepts the build order.

## 5. Write it to the map

Write after a zone completes, when the user pauses, and at the end:

1. Post one comment per decision (`**{MAP-SLUG}-GM-NNN: ...**`, see MAP.md). Put layouts and rationale inside it.
2. Update the body: tick the settled Open questions with their links and fill the Decisions table. Without a map, create one now with the body template.
   - **At a checkpoint:** leave **Phase** at `shaping` and **Next** at `/grill #N`, and note where the next session resumes.
   - **At the end** (every question settled, build order accepted): write **Slices** from the build order, set **Phase** to `shaped` and **Next** to `/build #N`.
3. Comment on each related map whose decisions this session affects. Name both sides' GM ids.

Close with a short summary of what was settled, what was parked and what is still open.

**Done when:** every settled decision has a GM comment linked from the body. At the end, the map also says `Next: /build #N`; at a checkpoint, it says `Next: /grill #N`.
