---
name: design-modules
description: Shape deep-module interfaces and seams for a map's slices before building, and post the design as a comment on the map. Use when the user asks for design-modules, module or interface design, or a design gate before a build.
agent-config-sync: true
---

# Design modules

Turn a map's decisions into **deep modules**, each with a small interface hiding a lot of behavior. Find the **seams** where they meet, then post the design on the map so the build's planner and builders follow it. This is a design gate. It runs before any code is written.

## Vocabulary

- **Module:** anything with an interface and an implementation. A function, a package, or a slice that spans tiers.
- **Interface:** everything a caller must know: signatures, invariants, ordering, errors, configuration.
- **Depth:** behavior per unit of interface. A deep module puts a lot behind a little. A shallow one has an interface about as complex as its code.
- **Seam:** a place where behavior can change without editing that place. Interfaces live at seams. With one adapter the seam is hypothetical; with two, it is real.
- **Dependency category:**
  - **in-process:** test it directly.
  - **local-substitutable:** test with a local stand-in such as a disposable database.
  - **remote but owned:** define a port with an in-memory adapter for tests.
  - **true external:** take an injected port and mock it in tests.

## 1. Load

Read the map's decisions and **Slices**, and the code the new modules sit next to. Pick out the callers, which things change together, and which change apart.

**Done when:** you can list the callers and the decisions each module must honor.

## 2. Design

Use the fewest modules that stay deep. Split only where a seam improves depth or gives the build a cleaner layer boundary. For each module, write:
- the problem it solves and who calls it;
- the recommended interface;
- its seam and dependency category;
- how tests cross that interface.

Apply the deletion test. If deleting a module makes the complexity vanish, it was pass-through. If the complexity reappears across callers, the module earns its place. Sketch a second shape only when the user asks, or when two shapes are close.

**Done when:** every slice's new code has a module that owns it, and every module has an interface and a test surface.

## 3. Post

Post one comment on the map, `## Module design`. Open with an overview of how the modules relate, then give one block per module. End with the build order the seams imply. Get the user's approval before the build starts. A design choice that changes a decision goes back through `grill`.

**Done when:** the comment is posted and the user has approved it or asked for changes.
