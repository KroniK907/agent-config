---
name: research
description: Research a question for a feature map - outside APIs, libraries, other software, our own systems - and post the findings as a comment on the map. Use when the user asks for research on a map or pauses a grill to look something up.
agent-config-sync: true
---

# Research

Findings, not decisions. Research feeds the map. The user settles decisions in a grill.

## 1. Frame

State the question in one sentence, along with what would answer it. That might be a recommendation, a comparison, or a measured number. Read the map (`gh issue view N`) and the code the question touches.

**Done when:** the question and its answer shape are in your reply.

## 2. Investigate

Start with primary sources: official docs, specs, the code, live systems you can read, and measurements you can run. Label any secondary source as secondary. When you weigh dependencies, check how mature and active each one is, and who maintains it. Find at least one serious alternative to your lean. If the evidence runs out, make one wider pass, then stop and name the gap.

**Done when:** every part of the answer shape has evidence or a named gap.

## 3. Post

Post a new comment on the map:

```markdown
## Research: <question>

**Summary:** 2-4 sentences answering the question, with your confidence.

**Findings:** evidence with links, `file:line` refs or numbers. Primary sources first.

**Alternatives:** the strongest other approach and what it costs.

**Gaps:** what is still unknown.

**For the grill:** the decisions this informs, each with a recommended answer.
```

In a grill, quote the summary back and resume the paused question. Outside a grill, link the comment in the map's **Related** line or the matching open question.

**Done when:** the comment URL is in your reply.
