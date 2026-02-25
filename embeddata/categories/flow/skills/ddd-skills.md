---
name: ccf-flow-ddd-skills
description: >
  Use this skill for creating or verifying DDD documentation. Triggers include:
  extracting domain terms from code, building a glossary, verifying glossary format,
  extracting DDD documentation, or when the user mentions "ubiquitous language",
  "domain terms", "glossary", "DDD", "domain model", "entities", "aggregates",
  "value objects", or references docs/flow/ddd/.
---

# DDD Documentation Skills

## Document Structure

DDD documentation is split into separate files by concept type. Each file covers one concept only. Do not combine concepts into a single file.

```
docs/flow/ddd/
  ubiquitous-language.md    ← domain glossary (table)
  entities.md               ← entities with identity and lifecycle
  value-objects.md          ← immutable types with no identity
  aggregates.md             ← aggregate roots and their boundaries (planned)
```

When extracting or creating a specific file, produce ONLY that file. Do not add sections for other concept types. For example, when creating `entities.md`, do not add `## Value Objects` or `## Aggregates` sections — those belong in their own files.

---

## Ubiquitous Language

File: `docs/flow/ddd/ubiquitous-language.md`

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

---

## Entities

File: `docs/flow/ddd/entities.md`

This file contains entities ONLY. An entity has identity (a unique ID), mutable state, and a lifecycle. Do not include aggregates or value objects — those have their own files.

### Format

The file has a `# Entities` heading, then one `## EntityName` section per entity. Each section has: a one-sentence description, an attributes table, a behavior line, and an invariants line. Nothing else. No `## Value Objects` sections, no `## Aggregates` sections, no `## Notes` sections.

```markdown
# Entities

## Task

The fundamental unit of work. Identified by a UUID, belongs to a user, tracks status through its lifecycle.

| Attribute | Type | Description |
|-----------|------|-------------|
| TaskId | UUID | Unique identity |
| Name | String | Required. The name of the task |
| Status | TaskStatus | Current lifecycle state. Defaults to Pending |
| DueDate | Date | Optional. Deadline for completion |
| ScheduledDate | Date | Optional. When the user plans to work on it |
| ParentTaskId | TaskId | Optional. Reference to parent task |
| ProjectId | ProjectId | Optional. Reference by ID only |
| Order | Decimal | Sort position within a block or list |

**Behavior:** <ul><li>Place</li><li>UpdateStatus</li><li>Reschedule</li><li>Backlog</li><li>Cancel</li></ul>

**Invariants:** <ul><li>Name must not be empty</li><li>Must have a valid status</li></ul>

---

## Stakeholder

A person, team, or organization relevant to one or more tasks.

| Attribute | Type | Description |
|-----------|------|-------------|
| StakeholderId | UUID | Unique identity |
| Name | String | Required. Unique per user |
| Type | StakeholderType | Person, Team, or Organization |

**Behavior:** <ul><li>Create</li><li>Update</li><li>Delete</li></ul>

**Invariants:** <ul><li>Name must not be empty</li><li>Name must be unique per user</li></ul>
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

### Entity Section Rules

Each entity section must have exactly these parts in this order:

1. `## EntityName` — heading
2. One sentence description — domain meaning, not implementation
3. Attributes table with columns `Attribute | Type | Description`
4. `**Behavior:**` followed by `<ul><li>` list of domain operations
5. `**Invariants:**` followed by `<ul><li>` list of rules that must always be true

### Column Rules for Attributes Table

- **Attribute**: PascalCase name. First row is always the identity field.
- **Type**: Domain type (UUID, String, Date, Decimal, Boolean, or a reference to another domain type like TaskStatus or ProjectId). No SQL types, no language-specific types.
- **Description**: Short description. Use "Required" or "Optional" prefix. Use "Reference by ID only" for cross-aggregate references.

### Structural Rules

- Entities sorted alphabetically by name
- Sections separated by `---`
- No implementation details: no file paths, table names, SQL types, column names, HTTP endpoints, migration numbers, or code syntax
- Cross-aggregate references use ID types (e.g., `ProjectId`) not object references
- The file contains ONLY entities — no value object sections, no aggregate sections, no notes sections

### What Is an Entity vs Other Concepts

- **Entity** (goes in `entities.md`): Has a unique ID, mutable state, and a lifecycle.
- **Value Object** (goes in `value-objects.md`): Immutable, no identity, compared by attributes. Examples: Money, TaskStatus, StakeholderType.
- **Aggregate** (goes in `aggregates.md`): A consistency boundary around one or more entities. Defines what is loaded and saved as a unit.

---

## Value Objects

File: `docs/flow/ddd/value-objects.md`

This file contains value objects ONLY. A value object is immutable, has no identity, and is compared by its attributes. Do not include entities or aggregates — those have their own files.

### Format

The file has a `# Value Objects` heading, then one `## VOName` section per value object. Each section has: a one-sentence description, an attributes table, and a constraints line. Nothing else. No `## Entities` sections, no `## Aggregates` sections, no `## Notes` sections.

```markdown
# Value Objects

## Money

An amount of currency. Immutable. Compared by amount and currency together.

| Attribute | Type | Description |
|-----------|------|-------------|
| Amount | Decimal | The monetary value |
| Currency | String | ISO 4217 currency code |

**Constraints:** <ul><li>Amount must be ≥ 0</li><li>Currency is required</li><li>Arithmetic only valid when currencies match</li></ul>

---

## TaskStatus

The lifecycle state of a task.

| Attribute | Type | Description |
|-----------|------|-------------|
| Value | Enum | <ul><li>Pending</li><li>InProgress</li><li>Completed</li><li>Deferred</li><li>Cancelled</li><li>Blocked</li></ul> |

**Constraints:** <ul><li>Must be a valid status value</li></ul>
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

### Value Object Section Rules

Each value object section must have exactly these parts in this order:

1. `## VOName` — heading
2. One sentence description — domain meaning, not implementation
3. Attributes table with columns `Attribute | Type | Description`
4. `**Constraints:**` followed by `<ul><li>` list of value constraints

### Column Rules for Attributes Table

- **Attribute**: PascalCase name. No identity field (value objects have no ID).
- **Type**: Domain type (String, Decimal, Date, Boolean, Enum, or a reference to another domain type). Use `Enum` for enumerated types, with values listed as `<ul><li>` in the Description column. No SQL types, no language-specific types.
- **Description**: Short description. For Enum types, list valid values using `<ul><li>`.

### Structural Rules

- Value objects sorted alphabetically by name
- Sections separated by `---`
- No implementation details: no file paths, table names, SQL types, column names, HTTP endpoints, migration numbers, or code syntax
- The file contains ONLY value objects — no entity sections, no aggregate sections, no notes sections

### What Is a Value Object vs Other Concepts

- **Value Object** (goes in `value-objects.md`): Immutable, no identity, compared by attributes. Includes enums, type classifications, and small composite values.
- **Entity** (goes in `entities.md`): Has a unique ID, mutable state, and a lifecycle.
- **Aggregate** (goes in `aggregates.md`): A consistency boundary around one or more entities.

---

## Agents

| Agent | Command | Purpose |
|-------|---------|---------|
| `ccf-flow-ddd-ul-extractor` | `/ccf-flow-ddd-ul-extract` | Scans code and docs, produces the glossary |
| `ccf-flow-ddd-ul-verifier` | `/ccf-flow-ddd-ul-verify` | Checks the glossary conforms to the format above |
| `ccf-flow-ddd-entity-extractor` | `/ccf-flow-ddd-entity-extract` | Scans code and docs, produces the entities file |
| `ccf-flow-ddd-entity-verifier` | `/ccf-flow-ddd-entity-verify` | Checks the entities file conforms to the format above |
| `ccf-flow-ddd-vo-extractor` | `/ccf-flow-ddd-vo-extract` | Scans code and docs, produces the value objects file |
| `ccf-flow-ddd-vo-verifier` | `/ccf-flow-ddd-vo-verify` | Checks the value objects file conforms to the format above |

## Verification

After creating or updating any DDD file, always run the corresponding verifier agent to check the output.

## Full Template References

See `references/ddd-ul-template.md` for ubiquitous language specification.
See `references/ddd-entity-template.md` for entity specification.
See `references/ddd-vo-template.md` for value object specification.
