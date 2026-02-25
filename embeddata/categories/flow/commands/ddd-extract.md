---
description: Extract a single DDD concept. Pass the concept name as an argument - ubiquitous-language, classification, entities, value-objects, or aggregates.
argument-hint: <concept> e.g. entities, value-objects, aggregates, classification, ubiquitous-language
---

Extract a single DDD documentation file based on the argument provided.

Determine which concept the user wants from the argument: $ARGUMENTS

Then use the corresponding sub-agent:

- **ubiquitous-language** or **ul** → Use the `ccf-flow-ddd-ul-extractor` sub-agent to extract into `docs/flow/ddd/ubiquitous-language.md`
- **classification** or **classify** → Use the `ccf-flow-ddd-classifier` sub-agent to extract into `docs/flow/ddd/classification.md`
- **entities** or **entity** → Use the `ccf-flow-ddd-entity-extractor` sub-agent to extract into `docs/flow/ddd/entities.md`. It must read `classification.md` first.
- **value-objects** or **vo** or **vos** → Use the `ccf-flow-ddd-vo-extractor` sub-agent to extract into `docs/flow/ddd/value-objects.md`. It must read `classification.md` first.
- **aggregates** or **agg** or **aggs** → Use the `ccf-flow-ddd-agg-extractor` sub-agent to extract into `docs/flow/ddd/aggregates.md`. It must read `classification.md` first.

If the argument is empty or doesn't match any concept, list the valid options and ask the user to try again.

Do not do the work yourself — delegate to the named sub-agent.
