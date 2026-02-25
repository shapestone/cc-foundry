---
description: Run full DDD extraction - ubiquitous language, classification, entities, value objects, aggregates. Chains sub-agents sequentially.
---

Run the following DDD extraction steps in order. For each step, use the named sub-agent. Wait for each sub-agent to complete before starting the next. Do not do the work yourself — delegate to the sub-agent.

## Step 1: Ubiquitous Language
Use the `ccf-flow-ddd-ul-extractor` sub-agent to extract the ubiquitous language glossary into `docs/flow/ddd/ubiquitous-language.md`.

## Step 2: Classification (Phase 1)
Use the `ccf-flow-ddd-classifier` sub-agent to classify domain concepts into entities and value objects. Write `docs/flow/ddd/classification.md`.

## Step 3: Entities (Phase 2)
Use the `ccf-flow-ddd-entity-extractor` sub-agent to extract entities into `docs/flow/ddd/entities.md`. It must read `classification.md` first.

## Step 4: Value Objects (Phase 2)
Use the `ccf-flow-ddd-vo-extractor` sub-agent to extract value objects into `docs/flow/ddd/value-objects.md`. It must read `classification.md` first.

## Step 5: Aggregates (Phase 2)
Use the `ccf-flow-ddd-agg-extractor` sub-agent to extract aggregates into `docs/flow/ddd/aggregates.md`. It must read `classification.md` first.

## Step 6: Cross-Reference Verification
Use the `ccf-flow-ddd-xref-verifier` sub-agent to verify that all four files (classification, entities, value-objects, aggregates) are consistent with each other.

## After all steps

Print a summary:

```
DDD Extraction Complete

1. Ubiquitous Language: ✅/❌
2. Classification:      ✅/❌
3. Entities:            ✅/❌
4. Value Objects:       ✅/❌
5. Aggregates:          ✅/❌
6. Cross-Reference:     ✅/❌
```

If any step fails, continue with the remaining steps. Note failures in the summary.

$ARGUMENTS
