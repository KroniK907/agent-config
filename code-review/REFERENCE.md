# Code review reference

Disclosed reference for [SKILL.md](SKILL.md) — smell baseline, spec lookup, sub-agent prompts, aggregation format.

---

## Spec lookup order

Resolve the originating spec before spawning the Spec sub-agent:

1. **Issue refs in commit messages** — `#123`, `Closes #45`, `Fixes org/repo#67`, etc. Fetch each via:

   ```powershell
   gh issue view <num> --json body,title,url
   ```

   For wayfinder work, also scan issue bodies for **What to build**, **Done when**, **Decisions**, and parent bundle / map links.

2. **User-provided path** — argument or explicit path in chat.

3. **Repo spec files** — `docs/`, `specs/`, `.scratch/` matching branch name or feature slug.

4. **Ask the user** — if nothing resolves. If they confirm no spec, skip Spec sub-agent.

**Not used:** Matt Pocock `docs/agents/issue-tracker.md` or external tracker conventions unless the repo explicitly documents them.

---

## Standards sources

Scan the repo for documented coding standards, in rough priority:

| Source | Examples |
|--------|----------|
| Project docs | `CONTRIBUTING.md`, `CODING_STANDARDS.md`, `AGENTS.md`, `CLAUDE.md` |
| Cursor rules | `.cursor/rules/*.mdc` |
| Skill conventions | [writing-for-agents](../writing-for-agents/SKILL.md), [wayfinder REFERENCE](../wayfinder/REFERENCE.md) for skills-repo work |
| Language/framework defaults | Only when repo is silent on the topic |

Two binding rules for the smell baseline:

- **Repo overrides** — a documented repo standard always wins; suppress baseline smells the repo endorses.
- **Judgement calls** — baseline smells are labelled heuristics ("possible Feature Envy"), never hard violations. Skip anything tooling already enforces (formatters, linters, typecheckers).

---

## Fowler smell baseline

Fixed set from Fowler (*Refactoring*, ch. 3). Match against the diff; each entry is *what it is* → *how to fix*:

- **Mysterious Name** — a function, variable, or type whose name doesn't reveal what it does or holds. → rename it; if no honest name comes, the design's murky.
- **Duplicated Code** — the same logic shape appears in more than one hunk or file in the change. → extract the shared shape, call it from both.
- **Feature Envy** — a method that reaches into another object's data more than its own. → move the method onto the data it envies.
- **Data Clumps** — the same few fields or params keep travelling together (a type wanting to be born). → bundle them into one type, pass that.
- **Primitive Obsession** — a primitive or string standing in for a domain concept that deserves its own type. → give the concept its own small type.
- **Repeated Switches** — the same `switch`/`if`-cascade on the same type recurs across the change. → replace with polymorphism, or one map both sites share.
- **Shotgun Surgery** — one logical change forces scattered edits across many files in the diff. → gather what changes together into one module.
- **Divergent Change** — one file or module is edited for several unrelated reasons. → split so each module changes for one reason.
- **Speculative Generality** — abstraction, parameters, or hooks added for needs the spec doesn't have. → delete it; inline back until a real need shows.
- **Message Chains** — long `a.b().c().d()` navigation the caller shouldn't depend on. → hide the walk behind one method on the first object.
- **Middle Man** — a class or function that mostly just delegates onward. → cut it, call the real target direct.
- **Refused Bequest** — a subclass or implementer that ignores or overrides most of what it inherits. → drop the inheritance, use composition.

Paste this section in full into the Standards sub-agent prompt — the sub-agent has no other access to it.

---

## Sub-agent prompts

Send both in **one message** as parallel Task calls with `subagent_type: generalPurpose`.

### Standards sub-agent

```text
Review this diff for coding standards compliance.

Repository: <absolute path>
Diff command: git diff <fixed-point>...HEAD
Commits: <paste git log --oneline output>

Standards sources (read these files if they exist):
- <list each path found in step 3>

Fowler smell baseline (apply when repo is silent; repo docs override):
<paste full smell baseline section>

Brief: Report — per file/hunk where relevant — (a) every place the diff violates a documented standard: cite the standard (file + rule); and (b) any baseline smell you spot: name it and quote the hunk. Distinguish hard violations from judgement calls — documented-standard breaches can be hard, but baseline smells are always judgement calls, and a documented repo standard overrides the baseline. Skip anything tooling enforces. Under 400 words.
```

### Spec sub-agent

```text
Review this diff against the spec.

Repository: <absolute path>
Diff command: git diff <fixed-point>...HEAD
Commits: <paste git log --oneline output>

Spec source: <issue #N title + URL, or file path>
Spec contents:
<paste fetched issue body or spec file contents>

Brief: Report: (a) requirements the spec asked for that are missing or partial; (b) behaviour in the diff that wasn't asked for (scope creep); (c) requirements that look implemented but where the implementation looks wrong. Quote the spec line for each finding. Under 400 words.
```

---

## Output format

After both sub-agents return, present:

```markdown
## Code review — `<fixed-point>`...HEAD

**Commits:** N · **Spec:** <issue #N / path / none>

## Standards

<sub-agent report — verbatim or lightly cleaned>

## Spec

<sub-agent report — or "Skipped — no spec available">

---

**Summary:** Standards: N findings (worst: …) · Spec: M findings (worst: …)
```

Rules:

- Keep **Standards** and **Spec** as separate top-level sections — never merge lists or rerank across axes.
- One-line summary counts findings **per axis** and names the worst issue **within each axis** only.

---

## Why two axes

A change can pass one axis and fail the other:

- Code that follows every standard but implements the wrong thing → **Standards pass, Spec fail.**
- Code that does exactly what the issue asked but breaks project conventions → **Spec pass, Standards fail.**

Reporting them separately stops one axis from masking the other.
