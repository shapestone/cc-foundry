---
description: Extract a single DDD concept. Pass the concept name as an argument - ubiquitous-language, classification, entities, value-objects, aggregates, rules, commands, or events.
argument-hint: <concept> e.g. entities, value-objects, aggregates, classification, ubiquitous-language, rules, commands, events
---

Extract a single DDD documentation file based on the argument provided.

Determine which concept the user wants from the argument: $ARGUMENTS

Then use the corresponding sub-agent:

- **ubiquitous-language** or **ul** → Use the `ccf-flow-ddd-ul-extractor` sub-agent to extract into `docs/flow/ddd/ubiquitous-language.md`
- **classification** or **classify** → Use the `ccf-flow-ddd-classifier` sub-agent to extract into `docs/flow/ddd/classification.md`
- **entities** or **entity** → Use the `ccf-flow-ddd-entity-extractor` sub-agent to extract into `docs/flow/ddd/entities.md`. It must read `classification.md` first.
- **value-objects** or **vo** or **vos** → Use the `ccf-flow-ddd-vo-extractor` sub-agent to extract into `docs/flow/ddd/value-objects.md`. It must read `classification.md` first.
- **aggregates** or **agg** or **aggs** → Use the `ccf-flow-ddd-agg-extractor` sub-agent to extract into `docs/flow/ddd/aggregates.md`. It must read `classification.md` first.
- **rules** or **business-rules** → Run TWO sub-agents in sequence:
  1. Use the `ccf-flow-ddd-rules-extractor` sub-agent to produce `docs/flow/ddd/rules-extraction.md`. It must read classification.md, entities.md, value-objects.md, and aggregates.md first.
  2. Then use the `ccf-flow-ddd-rules-analyzer` sub-agent to produce `docs/flow/ddd/business-rules.md`. It must read `rules-extraction.md` first.
- **commands** or **cmd** or **cmds** → Run TWO sub-agents in sequence:
  1. Use the `ccf-flow-ddd-commands-extractor` sub-agent to produce `docs/flow/ddd/commands-extraction.md`. It must first read the template at `.claude/skills/ccf-flow-ddd-skills/references/ddd-commands-extraction-template.md`, then read classification.md, entities.md, value-objects.md, and aggregates.md. The output must be a single flat table — no sections, no grouping, no file paths.
  2. Then use the `ccf-flow-ddd-commands-analyzer` sub-agent to produce `docs/flow/ddd/commands.md`. It must first read the template at `.claude/skills/ccf-flow-ddd-skills/references/ddd-commands-template.md`, then read `commands-extraction.md` and `entities.md`. The output must use the exact field names from the template: `**Behavior mapping:**`, `**Validation:**`, `**Side effects:**`, `**Status:**`.
- **events** or **domain-events** → Run TWO sub-agents in sequence:
  1. Use the `ccf-flow-ddd-events-extractor` sub-agent to produce `docs/flow/ddd/events-extraction.md`. It must first read the template at `.claude/skills/ccf-flow-ddd-skills/references/ddd-events-extraction-template.md`, then read classification.md, entities.md, aggregates.md, commands-extraction.md, and commands.md. The output must be a single flat table — no sections, no grouping, no file paths.
  2. Then use the `ccf-flow-ddd-events-analyzer` sub-agent to produce `docs/flow/ddd/events.md`. It must first read the template at `.claude/skills/ccf-flow-ddd-skills/references/ddd-events-template.md`, then read `events-extraction.md` and `commands.md`. The output must use the exact field names from the template: `**Triggered by:**`, `**Event type:**`, `**Subscribers:**`, `**Cascade effects:**`, `**Status:**`.

If the argument is empty or doesn't match any concept, list the valid options and ask the user to try again.

Do not do the work yourself — delegate to the named sub-agent.
