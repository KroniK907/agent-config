# Skill mechanics

The skill-specific branch of [`writing-for-agents`](SKILL.md). Covers frontmatter, invocation choice, storage, and router skills. Universal writing rules live in `SKILL.md`.

## Storage

| Type | Path | Scope |
|------|------|-------|
| Personal | `~/.cursor/skills/<name>/` | Available across all projects |
| Project | `.cursor/skills/<name>/` | Shared with anyone using the repository |

Never create skills in `~/.cursor/skills-cursor/`. That directory is reserved for Cursor built-in skills.

Layout:

```
skill-name/
├── SKILL.md              # Required - main instructions
├── REFERENCE.md          # Optional - disclosed reference
└── scripts/              # Optional - utility scripts
```

## Frontmatter

Every `SKILL.md` requires YAML frontmatter:

```yaml
---
name: skill-name
description: Brief description with trigger branches. Use when ...
---
```

| Field | Requirements |
|-------|--------------|
| `name` | Max 64 chars, lowercase letters/numbers/hyphens only |
| `description` | Max 1024 chars; the skill's top-level context pointer when model-invoked |

## Invocation

Two choices, trading the two loads.

**Model-invoked.** Keep a `description` so the agent can fire the skill on its own and other skills can reach it. You can still type its name. The description stays loaded every turn. That is permanent context load in exchange for discoverability. A model-invoked skill whose content is all reference can also be shared reference for other skills. Omit `disable-model-invocation`. Write a model-facing description with distinct trigger branches. Follow pointer rules in `SKILL.md`.

**User-invoked.** Set `disable-model-invocation: true`. Only the human typing the name can invoke it. Zero context load, but the human must remember it exists. The `description` becomes human-facing summary text. Trigger lists come out.

Pick model-invocation when the agent or another skill must reach the skill on its own. Pick user-invocation when only the human invokes it by hand.

Shared reference that two user-invoked skills both need cannot live in either skill. Push it to a plain file outside the skill system.

## Splitting by invocation

Split off a model-invoked skill when you have a distinct leading word that should trigger it on its own, or another skill must reach it. You pay context load for the new always-loaded description, so the independent reach has to be worth it. The sequence cut lives in `SKILL.md`.

## Router skills

When user-invoked skills multiply past what you can remember, add a **router skill**. One user-invoked skill names the others and when to reach for each. The human remembers one entry point instead of many. A router can only hint. User-invoked skills have no description, so nothing but the human can reach them.

## Skill creation workflow

When helping a user create a new skill:

1. **Gather requirements** - purpose, trigger scenarios, scripts vs instructions only, reference materials, personal vs project location.

   **Done when:** You can state purpose, invocation choice, and storage path.

2. **Draft** - apply `SKILL.md` levers (hierarchy, pointers, completion criteria, leading words). Keep SKILL.md under ~500 lines; disclose the rest.

   **Done when:** `SKILL.md` exists with frontmatter, steps, and pointers to disclosed reference.

3. **Review** - verify description triggers, terminology, no time-sensitive info, references one level deep.

   **Done when:** Description has distinct trigger branches (if model-invoked), no stale env caches, and the user has seen the draft.

If the user supplies exact wording for the skill, use it **verbatim** in `SKILL.md`.
