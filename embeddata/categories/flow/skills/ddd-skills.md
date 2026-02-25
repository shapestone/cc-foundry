---
name: ccf-flow-ddd-skills
description: >
  Use this skill for creating or verifying DDD documentation. Triggers include:
  extracting domain terms from code, building a glossary, verifying glossary format,
  extracting DDD documentation, or when the user mentions "ubiquitous language",
  "domain terms", "glossary", "DDD", "domain model", or references docs/flow/ddd/.
---

# DDD Documentation Skills

## Ubiquitous Language

The ubiquitous language glossary lives at `docs/flow/ddd/ubiquitous-language.md`.

### Format

The file must contain a heading and a markdown table. Nothing else. No preamble paragraphs, no `## Section` headers per term, no `**Definition:**` labels, no prose blocks, no footnotes.

```markdown
# Ubiquitous Language

| Term | Aliases | Meaning | Notes |
|------|---------|---------|-------|
| Backlog | — | The set of tasks with no scheduled time block or date, available for future scheduling. | — |
| BlockCapacity | — | The maximum number of incomplete tasks a time block can hold, derived from duration and max load. | For the active block, only remaining time counts toward capacity. |
| Order | Purchase Order, PO | A confirmed request by a Customer to buy one or more Products at agreed prices. | "Order" in Shipping context means shipment request. |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

### Column Rules

- **Term**: PascalCase canonical name. Never empty.
- **Aliases**: Comma-separated alternative names, or `—` if none.
- **Meaning**: One sentence. Domain concept only. Never empty.
- **Notes**: Brief disambiguation, or `—` if none.

### Structural Rules

- Rows sorted alphabetically by Term
- One row per term
- No content outside the heading and table
- No implementation details anywhere: no file paths, table names, SQL types, HTTP endpoints, function signatures, migration numbers, or code syntax

### What to Include

- Entities and their identity concepts (e.g., Task, Project, Stakeholder)
- Value objects and enums (e.g., TaskStatus, Money, BlockCapacity)
- Aggregates if they have a distinct name from their root entity
- Domain operations and processes (e.g., CascadeAlgorithm, DailyShutdown)
- UI/UX concepts that have domain meaning (e.g., Backlog, Archive, DecisionSurface)

### What to Exclude

- Infrastructure terms (database, router, middleware, handler, repository)
- Framework/library names
- Generic programming concepts (list, map, string, interface)
- Internal variable names that don't represent domain concepts

## Agents

| Agent | Command | Purpose |
|-------|---------|---------|
| `ccf-flow-ddd-ul-extractor` | `/ccf-flow-ddd-ul-extract` | Scans code and docs, produces the glossary |
| `ccf-flow-ddd-ul-verifier` | `/ccf-flow-ddd-ul-verify` | Checks the glossary conforms to the format above |

## Verification

After creating or updating the glossary, always run the `ccf-flow-ddd-ul-verifier` agent to check the output.

## Full Template Reference

See `references/ddd-ul-template.md` for the complete specification including forbidden content patterns.
