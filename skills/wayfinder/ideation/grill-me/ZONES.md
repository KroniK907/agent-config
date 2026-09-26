# Coverage zones

Five stack-neutral zones. Read this file **once per session**. If it is already in context (sous-vide, feature-discovery, or an earlier turn loaded it), do not re-read it.

Each zone has **Skip when** (a cell is N/A only with one of these reasons, stated) and **In scope if** (any hit keeps it open), then **Prompts** for inventing questions on in-scope cells. Web and non-web hints are examples, not checklists.

## Surfaces & experience (UI/UX, flows, copy, a11y)

- **Skip when:** no human- or operator-visible touchpoint changes - invisible refactor, headless API with no dashboard/email/PDF/export impact, code moved within a layer behind a stable contract. Confirm nothing downstream renders errors or admin views.
- **In scope if:** any new or changed screen, message, document, command, or help text. Name the surface before calling it N/A.
- **Layout pass:** every layout-bearing surface (page/route, modal, drawer/sheet, side panel, wizard step, full-screen view) needs a rough layout or a named deferral with an assumed default - see grill-me **UI layout articulation**.
- **Prompts:** who sees what (roles, internal/external); states (loading, empty, error, denied, partial failure, offline); entry/exit points and changes to existing paths; navigation and discoverability; accessibility (keyboard, focus, contrast, motion, error clarity); localization.
- **Web:** routes affected, new vs reused components, responsive differences, SSR/hydration effects. **Non-web:** CLI flags, defaults, `--help`, non-interactive mode; docs operators need.

## Behavior & correctness

- **Skip when:** behavior is literally unchanged (pure rename/move) or the work only documents existing behavior.
- **In scope if:** any branch, validation, state machine, retry, ordering, or consistency rule is added or altered.
- **Prompts:** happy paths trigger to outcome; edge cases (duplicates, conflicts, partial input, double-submit, stale data, idempotency); invariants; interaction with existing features and flags; concurrency (locking vs last-write-wins); failure behavior (compensation, partial commits, visible vs silent).
- **Web:** client vs server truth, optimistic UI, multi-tab sessions. **Non-web:** at-least-once delivery, poison messages; library API contracts and error types.

## Boundaries & integration

- **Skip when:** the change stays inside one process with no new or altered external contract - identical wire formats, schemas, call graphs.
- **In scope if:** any new endpoint, event, topic, webhook, third-party call, file format, or auth boundary moves.
- **Prompts:** trust (authn, authz, rate limits, credentials); contracts (shapes, errors, versioning, deprecation); transport (timeouts, retries, idempotency keys, duplicates); routing and discovery; side integrations (email, payment, identity) and their fallbacks.
- **Web:** routes, middleware order, CORS, cookies, CSRF, BFF vs direct. **Non-web:** queue ordering keys, replay, dead letters; file drops, schemas, checksums, partial uploads.

## Persistence & data

- **Skip when:** no durable state changes - no DB, object store, written files, cache semantics, retention, or PII impact (check derived fields and audit needs).
- **In scope if:** any schema, index, migration, new entity, PII, encryption, or at-rest permission change.
- **Prompts:** entities and fields (nullability, uniqueness, lifecycle); migrations (online/offline, backfill, rollback, dual-write); access control at rest; query patterns (pagination, hot paths, N+1); retention and compliance (GDPR, audit logs, encryption, key rotation).
- **Web:** API payload vs stored shape, caching, eventual consistency seen by UI. **Non-web:** warehouse vs OLTP, CDC, event-sourcing snapshots.

## Change, risk & evidence

- **Skip when:** personal or local change with no production path, collaborators, or rollback concern. Anything merged for others to run needs at least a light pass.
- **In scope if:** any rollout, customer impact, security sensitivity, SLO risk, or need to prove correctness.
- **Prompts:** rollout (flags, canary, kill switch); backward compatibility (old clients, jobs, data); operational readiness (runbooks, on-call); observability (logs, metrics, alerts); security and abuse deltas; evidence (tests, staging checks, acceptance criteria, sign-off).
- **Web:** cache/CDN invalidation, SEO, bundle size. **Non-web:** CI gates, migration windows, consumer contract tests.

## Non-product lens (skills, rules, docs, CI)

- **Surfaces:** who reads it (humans, agents); structure, reading order, discoverability, examples. Layout pass applies only if the artifact implies screens.
- **Behavior:** what it requires vs forbids; edge cases; consistency with other rules and skills.
- **Boundaries:** tools, repos, paths touched; composition with other skills or automation.
- **Persistence:** source of truth (paths, config keys); duplication vs links; template versioning.
- **Change:** rollout, blast radius, review expectations, proof the change works.
