# Research reference

## Research ticket template

Use when [feature-discovery](../../ideation/feature-discovery/SKILL.md) **Materialize** or wayfinder **Materialize** creates a `wf:research` ticket. Title: **`Research: {investigation topic}`** - see [ticket title conventions](../../REFERENCE.md#ticket-title-conventions).

```markdown
## Question

<Single investigation this ticket resolves - one session of fact-gathering.>

## Done when

- [ ] <Acceptance bullet 1 - verifiable from findings>
- [ ] <Acceptance bullet 2>

## Map

Parent: [{FeatureName}:Map](#parent-issue-number)

## Source hints

<Optional - URLs, repos, docs, code paths to start from.>

## Perspectives

<Optional - viewpoints or stakeholders to seek; alternate framings to cover.>
```

**Required:** `## Question`, `## Done when`, `## Map`.
**Optional:** `## Source hints`, `## Perspectives`.

Labels: `wf:todo` + `wf:research` + `wf:hitl` (v1 default) or `wf:afk` (deferred).

---

## Findings comment template

Post as a **new top-level comment** on the research ticket each session.

```markdown
## Research findings

**Ticket:** [#N](url) - **Map:** [{FeatureName}:Map](map-url)
**Session:** <one-line scope note>

### Summary

<2-4 sentences - direct answer to Question; confidence level.>

### Findings

<Primary evidence first - facts, links, code refs, quotes with sources. Label secondary sources explicitly.>

### Gaps & follow-ups

<What remains unknown; suggested follow-up questions or tickets.>

### Viewpoints/alternatives

<≥1 alternate framing, stakeholder view, or competing approach - even when consensus exists.>

### Coverage

| Done when | Status | Notes |
|-----------|--------|-------|
| <bullet from ticket> | satisfied / partial / not satisfied | <brief evidence> |

### Invalid premise

<!-- Include this section ONLY when the ticket Question is logically wrong, mis-categorized, or impossible - NOT when sources are missing. -->

<Why the ticket premise fails; suggested reticket or handoff skill.>

### Proposed tracker updates

<!-- Non-binding - Reconcile applies these once the user accepts them. -->

- **Map fog:** …
- **New ticket candidate:** …
- **Notes for map:** …
- **Decision log (for grilling, not direct append):** …
```

**Section order (fixed):** Summary → Findings → Gaps & follow-ups → Viewpoints/alternatives → Coverage → Proposed tracker updates.

**Invalid premise** sits after Coverage, before Proposed tracker updates, and **only when triggered**.

---

## Behavior rules

### Source hierarchy

| Tier | Examples | Weight |
|------|----------|--------|
| Primary | Official docs, source code, specs, standards, first-party APIs | Highest - lead Findings |
| Secondary | Blog posts, Stack Overflow, third-party tutorials | Lower - label explicitly; use to fill gaps only |

### Viewpoints/alternatives

Always seek **≥1 alternate viewpoint** - competing approach, dissenting opinion, stakeholder with different goals, or "do nothing" baseline. If none found, state that and why.

### Scope expansion

When **Coverage** shows `not satisfied` on a **valid premise** (question is sensible but evidence is thin):

1. Run **one** scope-expansion pass - broaden search, adjacent docs, related code paths
2. Re-assess Coverage
3. Stop - do not loop expansion; record remaining gaps in **Gaps & follow-ups**

### Coverage mapping

Each **Done when** bullet maps to exactly one row:

| Status | Meaning |
|--------|---------|
| `satisfied` | Evidence fully meets the bullet |
| `partial` | Some evidence; material gaps remain |
| `not satisfied` | No meaningful evidence after investigation (+ optional expansion pass) |

### Invalid premise

Use **only** for logical/category ticket errors:

- Question asks for binding decision (→ grill-me / strategic-ideation)
- Question is impossible or contradictory
- Ticket type wrong (should be prototype, grilling, or task)

**Not invalid premise:** couldn't find sources, partial answers, or disagreeing experts.

---

## Handoff

Findings are fact-gathering, not decisions. Decisions go through [grill-me](../../ideation/grill-me/SKILL.md) or [strategic-ideation](../../ideation/strategic-ideation/SKILL.md); tracker sync goes through wayfinder **Reconcile**, which merges **Proposed tracker updates** into its [resolution comment](../../references/reconcile.md#resolution-comment). Each follow-up session posts a new comment.

---

## Route trigger (for wayfinder)

| Trigger | Suggest `research` |
|---------|-------------------|
| Open unblocked `wf:research` ticket on map frontier | Yes - default |
| User explicitly invokes research / names a research ticket | Yes |
| User asks general question without a ticket | No - answer inline or suggest Materialize first |

Explicit user invoke always valid even when ticket is not frontier.

---

