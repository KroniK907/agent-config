# Subagent model data

Checked 2026-10-08. Prices are first-party API list prices per million tokens, with no subscription or batch discounts. Scores are the Artificial Analysis Intelligence Index (v4.3). "Cost/task" is Artificial Analysis's measured cost per index task, which folds in how many tokens each model and effort level writes. "Time/task" is its mean end-to-end response time.

## Prices

| Model | Input | Output | Long-prompt surcharge |
|-------|------:|-------:|-----------------------|
| Claude Haiku 5.5 | $0.10 | $0.50 | Over 100K prompt tokens: $0.50 / $2.50 |
| Claude Sonnet 5.5 | $2.00 | $10.00 | none |
| Claude Opus 5.5 | $4.00 | $20.00 | none |
| Claude Fable 5.1 | $10.00 | $50.00 | none |
| GPT-6 Luna | $0.10 | $0.50 | Over 272K input tokens: $0.20 / $0.75 |
| GPT-6.1 Sol | $2.00 | $10.00 | Over 272K input tokens: $4.00 / $15.00 |
| GPT-6 Astra | $10.00 | $50.00 | Over 272K input tokens: $20.00 / $75.00 |
| GPT-5.6 Terra | $2.00 | $12.00 | Over 272K input tokens: $4.00 / $18.00 |

## Intelligence, cost, and speed by effort

`max` rows are listed to show what the skill gives up by skipping it.

### Claude

| Model | Effort | Index | Cost/task | Time/task |
|-------|--------|------:|----------:|----------:|
| Haiku 5.5 | low | 29 | $0.02 | 13s |
| Haiku 5.5 | medium | 34 | $0.05 | 17s |
| Haiku 5.5 | high | 38 | $0.08 | 31s |
| Haiku 5.5 | xhigh | 41 | $0.12 | 90s |
| Haiku 5.5 | max | 43 | $0.21 | 434s |
| Sonnet 5.5 | low | 36 | $0.35 | 6s |
| Sonnet 5.5 | medium | 41 | $0.48 | 7s |
| Sonnet 5.5 | high | 47 | $0.88 | 18s |
| Sonnet 5.5 | xhigh | 52 | $2.01 | 41s |
| Sonnet 5.5 | max | 56 | $5.46 | 465s |
| Opus 5.5 | low | 42 | $0.55 | 15s |
| Opus 5.5 | medium | 51 | $1.34 | 29s |
| Opus 5.5 | high | 54 | $1.82 | 44s |
| Opus 5.5 | xhigh | 56 | $3.46 | 138s |
| Opus 5.5 | max | 58 | $5.98 | 737s |
| Fable 5.1 | low | 47 | $2.37 | 13s |
| Fable 5.1 | high | 51 | $3.91 | 31s |
| Fable 5.1 | xhigh | 53 | $5.98 | 89s |

### GPT

| Model | Effort | Index | Cost/task | Time/task |
|-------|--------|------:|----------:|----------:|
| GPT-6 Luna | low | 22 | $0.0045 | 7s |
| GPT-6 Luna | medium | 30 | $0.02 | n/a |
| GPT-6 Luna | high | 33 | $0.03 | 26s |
| GPT-6 Luna | xhigh | 35 | $0.04 | 41s |
| GPT-6 Luna | max | 38 | $0.07 | 115s |
| GPT-6.1 Sol | low | 42 | $0.13 | 14s |
| GPT-6.1 Sol | medium | 48 | $0.21 | 19s |
| GPT-6.1 Sol | high | 50 | $0.32 | 71s |
| GPT-6.1 Sol | xhigh | 51 | $0.39 | 158s |
| GPT-6.1 Sol | max | 52 | $0.72 | 339s |
| GPT-6 Astra | low | 46 | $0.82 | 14s |
| GPT-6 Astra | high | 51 | $1.73 | 105s |
| GPT-6 Astra | xhigh | 52 | $2.31 | 237s |
| GPT-5.6 Terra | max | 42 | $1.40 | 142s |

## How the tiers were drawn

- **scout:** Haiku 5.5 `low` and Luna `medium` both land near index 30 for about $0.02 a task. Lookups and summaries do not need more.
- **edit:** Haiku 5.5 `high` (38, $0.08) is the last cheap step before latency climbs: `xhigh` adds three points but triples the time. Luna `high` is the starting effort OpenAI's Codex docs suggest for Luna.
- **build:** Sonnet 5.5 `high` (47, $0.88, 18s). Opus 5.5 `medium` scores four points more for half again the cost and runs slower. Sonnet 5.5 also leads the Artificial Analysis Coding Agent Index in Claude Code. On the GPT side, Sol `medium` (48, $0.21) is the step where Sol's gains flatten.
- **hard:** Opus 5.5 `high` (54, $1.82) beats Sonnet 5.5 `xhigh` (52, $2.01) on both score and cost. Sol `high` (50, $0.32) is GPT's best value before `xhigh` doubles the time for one point.
- **Dominated models:** Fable 5.1 never beats Opus 5.5 on cost per point. Astra trails Sol at equal scores. Terra at `max` scores what Sol scores at `low`, for ten times the cost.

## Caveats

- Haiku 5.5 writes two to three times more tokens per task than Luna. Equal list prices do not mean equal bills. In Cognition's FrontierCode runs, Haiku 5.5 at `max` scored 46.4% to Luna's 42.4%, at $1.33 per rollout against $0.10.
- Anthropic's own benchmarks put Haiku 5.5 well ahead of Luna on agentic coding (Terminal-Bench 4.0: 39.2% against 16.4%). Those runs used different harnesses.
- Early user reports describe Haiku 5.5 as strong on quick coding, extraction, and summaries, and weaker on long reasoning chains and judgement calls. That is why **build** and **hard** skip it.
- A subagent's context grows with every file it reads. A Haiku subagent that crosses 100K prompt tokens pays five times its base rate from then on.
- Index versions change often. Artificial Analysis shipped three index versions in September 2026, and scores from different versions are not comparable.

## Sources

- [Artificial Analysis LLM leaderboard](https://artificialanalysis.ai/leaderboards/models): index, cost per task, speed, latency
- [Artificial Analysis: Claude Haiku 5.5](https://artificialanalysis.ai/models/claude-haiku-5-5)
- [Artificial Analysis Coding Agent Index post, October 2026](https://x.com/ArtificialAnlys/status/2105814318294114720)
- [OpenAI API pricing](https://developers.openai.com/api/docs/pricing)
- [Anthropic: Introducing Claude Haiku 5.5](https://www.anthropic.com/claude-haiku-5-5)
- [The Neuron: Haiku 5.5 pricing, benchmarks and agent costs](https://www.theneuron.ai/news/haiku-5-5-claude-most-practical-launch/)
- [Kingy AI: Haiku 5.5 vs GPT-6 Luna](https://kingy.ai/blog/claude-haiku-5-5-vs-gpt-6-luna-benchmark/)
- [MyClaw: Haiku 5.5 vs GPT-6 Luna cost, coding and agents](https://myclaw.ai/blog/claude-haiku-5-5-vs-gpt-6-luna)
- [Codex subagents docs](https://learn.chatgpt.com/docs/agent-configuration/subagents)
- [SourceForge: Claude Haiku 5.5 reviews](https://sourceforge.net/software/product/Claude-Haiku-5.5/)
