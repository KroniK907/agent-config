---
name: feature-ideation
description: Expand a rough feature seed through alternating ideation and tension rounds at idea level, then prune to clear boundaries for grill-me or PRD prep. Use when the user wants to ideate on a new feature, balloon scope before cutting it down, alternate expansion with strategic tensions, surface conflicts with existing initiatives, or produce a scoped description before implementation grilling.
---

# Feature ideation (expand → tension → prune)

Take a **simple seed** (or chat history) and grow an **oversized** picture through **ideation**, stress it with **idea-level tensions** (like grill-me but not implementation), repeat bands until shape stabilizes, then **prune** to a single bounded description suitable for **[grill-me](../grill-me/SKILL.md)**.

## Roles (fixed)

| Mode | Job |
|------|-----|
| **Ideation** | Expansion: new contexts, use cases, surfaces, future-proofing, synergy if adjacent areas grow. Add branches and possibilities—bias toward **too big** before pruning. |
| **Tension** | Idea/strategy stress-test: conflicts with existing features or roadmap, separation vs bundling, **internal consistency** of goals. **One question per assistant reply.** Not API/load-testing—that belongs in grill-me. |
| **Prune** | Consolidate into a **handoff artifact** (template in [REFERENCE.md](REFERENCE.md)): positive **In scope** + **Boundaries**; **no** default “out of scope / deferred” dumps for grill-me (see REFERENCE handoff rules). |

## Rhythm: alternating bands

Default loop:

1. **Ideation band:** **3 rounds** (each round = **one full exchange**: user message → assistant reply).
2. **Tension band:** **2 rounds** (each tension reply asks **exactly one** question).

**One assistant reply = one round only.** Never combine rounds in a single message (for example, do **not** complete ideation round 3 and start tension round 1 in the same reply). When a round ends, **stop** and **wait for the user’s next message** before advancing the round counter or switching phase/band—even at ideation→tension or tension→ideation boundaries.

Repeat until entering **prune**. Adjust band lengths only when the user asks or when the seed is tiny—see [REFERENCE.md](REFERENCE.md).

Forked takes (optional): During ideation, occasionally offer **2–3 deliberately different, intentionally oversized** interpretations of the seed so expansion stays vivid—merge user picks into the running picture.

## Entering prune mode

Use prune when **either**:

1. **User asks** (e.g. “let’s prune,” “finalize scope,” “hand off to grill”).
2. **Assistant softly suggests**—when ideation loops repeat without new ground, expansions keep collapsing, or tensions already surfaced dominate—**never** force; invite (“Want to prune, or one more ideation band?”).

If the user prefers more cycles, resume alternating bands.

## Every reply: recap first (all phases)

**Recap is mandatory** on every assistant turn (including the first—see below). It is the **first** substantive block in the message.

1. **Recap:** Short summary of **the last exchange**: what was decided, proposed, rejected, or how you interpreted the user’s last message. Gives the user a chance to correct drift. On the **first** reply with no prior exchange: state the seed or history you are using and that the session is starting in **IDEATION** at **Round 1/3**.
2. **Session state:** One line per [REFERENCE.md](REFERENCE.md) (`Phase`, `Band`, `Round`, `Next`).
3. **Body:** Ideation exploration, tension **Question** (single), or prune synthesis—per current phase.

Rename or merge Recap + Session state only if the client UI makes a single combined header clearer—the **information** must always appear.

## Output shape by phase

### Ideation band

After Recap + Session state:

- Drive **conversational** expansion (not a single-question grill). You may propose several related directions in one reply; stay grounded in the user’s seed and accumulated list of in/out scope as it emerges.
- Track **working** scope notes as they emerge (labels provisional until prune). During **prune**, follow [REFERENCE.md](REFERENCE.md) handoff rules for grill-me: affirmative boundaries; avoid negative backlog lists in the grill pack.

### Tension band

After Recap + Session state:

Mirror grill-me discipline for the **question** only:

**Question:**  
Exactly **one** tension question (conflict, consistency, boundary with existing work). No multi-part or stacked questions.

Optional short **Recommendation:** line—your lean on how to answer *that* question only—same spirit as grill-me.

### Prune

After Recap + Session state:

- Synthesize threads into one coherent picture.
- Produce the handoff per [REFERENCE.md](REFERENCE.md): **In scope**, **Boundaries** (positive edges), optional tiny **Explicitly not building** only when misleading defaults exist; **Tensions**; **Ready for grill-me**. Put backlog/deferred material only under **Notes for PRD**, not in the grill-me seed.
- Offer edits until the user accepts the block.

## Mechanical tracking

Maintain internally:

- **Phase:** `IDEATION` | `TENSION` | `PRUNE`.
- **Band:** increments each time you finish an ideation segment + tension segment cycle (or document your counting if you reset per user request).
- **Round-within-band:** e.g. ideation **2/3**, then tension **1/2**.

Expose this in **Session state** every turn so “what’s next” is obvious.

## Interaction rules (non-negotiable)

1. **Round** = **one full exchange** (user message → assistant reply).
2. **No stacked rounds:** Each assistant message must address **exactly one** round in the current phase. You may describe in **Session state** what comes *next* after the user replies, but you must **not** perform the next round (or next phase) until the user sends another message.
3. **Recap** first every time; never skip because the reply is long or late in session.
4. **Tension:** exactly **one** question per tension reply.
5. **Prune suggestion:** soft invitation only; user controls exit from expansion loops.
6. **Grill-me handoff content:** Give grill-me mostly **what to build** and open decisions on that slice. **Do not** fill the grill seed with big out-of-scope or “later” lists—see [REFERENCE.md](REFERENCE.md)—unless the user asks for a separate human-only appendix.
7. **Not a replacement for grill-me:** If ideation and grilling happen in one chat, **label the mode** each turn so behavior stays correct.

## Quick start

User provides a one-sentence seed (or points to prior chat). Respond with Recap (start rules) + Session state + first ideation expansion. After **3** ideation rounds, switch to **2** tension rounds, then repeat until prune.

See [REFERENCE.md](REFERENCE.md) for the handoff template and session-state examples.
