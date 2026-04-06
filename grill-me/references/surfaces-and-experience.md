# Surfaces & experience

## Quick triage (read every session)

**Likely N/A when:** there is no human-visible or operator-visible touchpoint (pure library internals, batch-only with no UI/CLI/reporting changes, invisible refactor). Still confirm nothing downstream renders errors or admin views.

**Common skip signals**

- Change is intentionally **invisible** to users and operators, and **no** docs, notifications, or exports change.
- **Headless** integration (API/machine-only) with **no** dashboard, email, or PDF impact.
- You are only moving code **within** a layer that already has a stable external contract.

If **any** new or changed **screen, message, document, command, or help text** exists, this zone is **in scope**—do not mark N/A without stating what the surface is.

## Deep prompts (load when not N/A)

- **Who** sees what, in which contexts (roles, admin vs self-serve, internal vs external)?
- **States:** loading, empty, success, validation errors, permission denied, partial failure, offline/slow—what should each surface show?
- **Flows:** entry points, exit points, and how this alters **existing** UX paths (dead ends, extra steps, changed affordances).
- **Information architecture:** where does this live in navigation; does it need discoverability work?
- **Accessibility & usability:** keyboard, focus, semantics, contrast, motion; error message clarity.
- **Localization:** copy constraints, pluralization, variable order, RTL if relevant.

### Web app overlay (use when relevant)

- **Routes/pages** affected; **new vs reused** components; layout constraints (grids, modals, drawers).
- **Responsive:** mobile vs tablet vs desktop—what differs (information priority, actions, breakpoints)?
- **SSR/CSR** or **hydration** assumptions that affect UX or SEO snippets.

### Non-web examples

- **CLI:** flags, defaults, `--help`, progressive output, non-interactive mode.
- **Docs/help center:** what must be updated so users/operators succeed.
