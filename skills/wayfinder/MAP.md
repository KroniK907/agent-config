# Map issue

One GitHub issue per feature, titled `{FeatureName}:Map` and labelled `wf:map`. It is both the map and the decision log. `map`, `grill` and `build` all read and write it.

## Decisions

- Each decision is its own comment whose first line is `**{MAP-SLUG}-GM-NNN: <one-line decision>.**`, followed by the detail: bullets, rationale, ASCII layouts.
- IDs count up from `001` per map. Find the next one from the existing comments: `gh issue view N --comments`.
- A later comment that reuses an ID replaces the earlier one. Revise by posting the replacement, then repoint the body links.
- Link a decision by its comment URL (`#issuecomment-...`). Other maps cite it as `{MAP-SLUG}-GM-NNN` with that link.
- Research findings and cross-map notes are also comments on the map, not separate issues.

## Body template

```markdown
# {FeatureName}:Map

**Phase:** <shaping | shaped | building | built | released>. <one line of status>. Next: <`/grill #N` | `/build #N` | ships in the next release; prod steps on Release:Next #R>.
**Map slug:** `{MAP-SLUG}`. Decision IDs are `{MAP-SLUG}-GM-NNN`.
**Parent:** <parent map, if any>. **Related:** <#N (why)>, ...

This one issue is both the map and the decision log. Each decision is a comment starting with its ID. The **Decisions** table at the bottom indexes them. If a later comment reuses an ID, it replaces the earlier one.

## Target outcome
1-3 numbered outcomes a user would notice.

## What's wrong today
Bullets: the current behavior and why it falls short.

## Shape
ASCII diagram of the parts and how they connect.

## Open questions
### <Area>
- [ ] **<Question>.** <one-line answer once settled> [{MAP-SLUG}-GM-NNN](comment-url)

## Slices
1. **<Name>.** What it delivers. ([GM-NNN](url), ...)

Which slices can run in parallel and which wait on which.

## Not yet specified
Known gaps parked for later.

## Out of scope
What this map leaves to other maps or issues, with links.

## Decisions
| ID | Decision |
|---|---|
| [{MAP-SLUG}-GM-001](comment-url) | <one line> |
```

**Slices** stays empty until the grill's last decision (build order and evidence) fills it. Each slice is sized for one agent and one PR.

## Phase

The **Phase** line is the first thing a reader trusts, so whoever moves the work also moves the Phase in the same session.

| Phase | True when | Next |
|---|---|---|
| `shaping` | The map is filed and not fully grilled. | `/grill #N` |
| `shaped` | Every question is settled and **Slices** is written. Nothing is built. | `/build #N` or `/sol-build #N` |
| `building` | A build is running, or some slices are merged and others are not. Slices merged only into a feature branch count here, because only `dev` leads to a release. | The open slices, or the feature branch merge |
| `built` | Every slice is merged into `dev`, and its human steps are on the release issue (`release-steps` skill). | `ships in the next release; prod steps on Release:Next #R` |
| `released` | The work is on `main`, and every human step it needs is done on prod. Close the map. | none |

A map that splits slices into child maps says so on the Phase line and lists each child with its phase. A hub map's Phase is the earliest phase of its unfinished children.

## Editing the body

Large bodies go through the `wf` helper so only what you need enters context and uploads are verified. Resolve the wayfinder folder with `readlink -f` first.

```bash
go run <wayfinder>/utilities/wf/wf.go section <N> "Slices" "Open questions"   # print sections
go run <wayfinder>/utilities/wf/wf.go body get <N> /tmp/map-N.md              # edit the file...
go run <wayfinder>/utilities/wf/wf.go body put <N> /tmp/map-N.md              # ...then upload and verify
```

Refer to other issues by title in prose, with the link. Keep planning IDs (`GM-`, issue numbers) out of code; PR bodies may say "Part of #N".
