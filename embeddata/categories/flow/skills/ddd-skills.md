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
  classification.md         ← what is an entity, VO, or aggregate (phase 1)
  entities.md               ← entities with identity and lifecycle (phase 2)
  value-objects.md          ← immutable types with no identity (phase 2)
  aggregates.md             ← aggregate roots and their boundaries (phase 2)
  rules-extraction.md       ← raw rule evidence from codebase (phase 1)
  business-rules.md         ← analyzed rules with contradiction detection (phase 2)
  commands-extraction.md    ← raw command evidence from codebase (phase 1)
  commands.md               ← analyzed commands with coverage detection (phase 2)
  events-extraction.md      ← raw event/side-effect evidence from codebase (phase 1)
  events.md                 ← analyzed events with cascade detection (phase 2)
  contexts-extraction.md    ← raw coupling evidence between aggregates (phase 1)
  contexts.md               ← proposed bounded contexts with coupling map (phase 2)
  services-extraction.md    ← raw cross-aggregate service evidence (phase 1)
  services.md               ← analyzed services with encapsulation assessment (phase 2)
```

When extracting or creating a specific file, produce ONLY that file. Do not add sections for other concept types. For example, when creating `entities.md`, do not add `## Value Objects` or `## Aggregates` sections — those belong in their own files.

## Extraction Workflow

Extracting entities, value objects, and aggregates requires a two-phase approach because they depend on each other:

**Phase 1 — Classification:** Run the `ccf-flow-ddd-classify` command first. This scans the codebase and produces `docs/flow/ddd/classification.md`, which records what is an entity, what is a value object, what is an aggregate root, and how they relate. All cross-cutting decisions are made here.

**Phase 2 — Per-file extraction:** Run the individual extractors (`ccf-flow-ddd-entity-extract`, `ccf-flow-ddd-vo-extract`, `ccf-flow-ddd-agg-extract`). Each reads `classification.md` to ensure consistency — using the same type names, the same entity/VO/aggregate assignments, and the same cross-references.

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

## Classification (Phase 1)

File: `docs/flow/ddd/classification.md`

This file is produced FIRST, before entities, value objects, or aggregates. It captures every cross-cutting decision: what is an entity, what is a value object, what is an aggregate root, and how they reference each other. The per-file extractors read this to stay consistent.

### Format

The file has a `# DDD Classification` heading, then two tables: Entities and Value Objects.

```markdown
# DDD Classification

## Entities

| Name | Aggregate Root | Contains | Uses VOs |
|------|---------------|----------|----------|
| Task | Yes | <ul><li>Stakeholder (ref, M:N)</li><li>Tag (ref, M:N)</li><li>Child Tasks (ref)</li></ul> | <ul><li>TaskStatus</li><li>DerivedStatus</li></ul> |
| Project | Yes | — | — |
| Stakeholder | Yes | — | <ul><li>StakeholderType</li></ul> |
| TimeBlock | Yes | — | <ul><li>BlockCapacity (derived)</li></ul> |

## Value Objects

| Name | Kind | Owned By |
|------|------|----------|
| TaskStatus | Enum | Task |
| StakeholderType | Enum | Stakeholder |
| ContainerType | Enum | Container |
| DerivedStatus | Derived | Task |
| BlockCapacity | Derived | TimeBlock |
| DayStats | Derived | — (computed at request time) |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

### Column Rules — Entities Table

- **Name**: PascalCase entity name.
- **Aggregate Root**: `Yes` or `No`.
- **Contains**: Other entities or references within this aggregate's boundary, using `<ul><li>`. Use `—` if none. Indicate relationship type (ref, M:N, child).
- **Uses VOs**: Value objects used as attribute types, using `<ul><li>`. Use `—` if none.

### Column Rules — Value Objects Table

- **Name**: PascalCase VO name.
- **Kind**: `Enum`, `Composite`, or `Derived`.
- **Owned By**: The entity that primarily uses this VO, or `—` if shared/computed.

### Structural Rules

- Tables sorted alphabetically by Name
- No implementation details: no file paths, table names, SQL types, HTTP endpoints, or code syntax
- Every entity and value object in the codebase should appear in exactly one table
- If unsure whether something is an entity or VO, classify it and add a note

---

## Aggregates

File: `docs/flow/ddd/aggregates.md`

This file contains aggregates ONLY. An aggregate defines a consistency boundary — what is loaded and saved as a unit. Do not include entity details or value object details — those have their own files. Read `classification.md` first for entity/VO assignments.

### Format

The file has a `# Aggregates` heading, then one `## AggregateName` section per aggregate. Each section has: a one-sentence description, a boundary table, and an invariants line. Nothing else.

```markdown
# Aggregates

## Task

The central unit of work and scheduling. Loaded and saved as a single unit with its stakeholder and tag associations.

| Element | Type | Relationship |
|---------|------|-------------|
| Task | Root Entity | — |
| Stakeholder | Entity (ref) | M:N via join, replaced on every write |
| Tag | Entity (ref) | M:N via join, replaced on every write |
| Child Tasks | Entity (ref) | One-to-many, read-only on parent |

**Invariants:** <ul><li>Must have a non-empty name</li><li>Status must be a valid TaskStatus value</li><li>Cross-aggregate references (ProjectId, TimeBlockId) are by ID only</li></ul>

---

## Project

A named grouping of tasks with optional nesting.

| Element | Type | Relationship |
|---------|------|-------------|
| Project | Root Entity | — |

**Invariants:** <ul><li>Name must not be empty</li><li>Root-level names must be unique per user</li><li>System projects cannot be modified or deleted</li></ul>
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

### Aggregate Section Rules

Each aggregate section must have exactly these parts in this order:

1. `## AggregateName` — heading
2. One sentence description — what the aggregate represents and its consistency scope
3. Boundary table with columns `Element | Type | Relationship`
4. `**Invariants:**` followed by `<ul><li>` list of rules enforced at the aggregate boundary

### Column Rules for Boundary Table

- **Element**: PascalCase name of each entity, value object, or collection that is loaded/saved as part of this aggregate.
- **Type**: `Root Entity`, `Entity`, `Entity (ref)`, or `Value Object`. Use `(ref)` for references that cross aggregate boundaries.
- **Relationship**: How this element relates to the root: `M:N via join`, `One-to-many`, `read-only`, or `—` for the root itself.

### What Belongs in the Boundary Table

Include elements that are **loaded, saved, or computed as part of this aggregate**:
- The root entity itself
- Entities or collections fetched and persisted with the root (e.g., Stakeholder M:N associations on Task)
- Value objects composed inline (e.g., StatSnapshot within StatHistory)
- Derived/computed value objects (e.g., DerivedStatus on Task)

**Do NOT include bare ID-only cross-aggregate references.** If an entity just stores another aggregate's ID as a foreign key (e.g., ParentContainerId, ProjectId, ScheduledTimeBlockId), that is an attribute on the entity — it belongs in `entities.md`, not in the boundary table. The boundary table shows what is loaded as a unit, not what is pointed to.

### Structural Rules

- Aggregates sorted alphabetically by name
- Sections separated by `---`
- No entity attribute details — those live in `entities.md`
- No value object details — those live in `value-objects.md`
- No implementation details: no file paths, table names, SQL types, HTTP endpoints, or code syntax
- The file contains ONLY aggregates — no entity sections, no VO sections, no notes sections

---

## Business Rules

Business rules extraction uses the same two-phase approach as structural extraction.

### Phase 1 — Rule Extraction

File: `docs/flow/ddd/rules-extraction.md`

Scan the entire codebase for anything that enforces, validates, or constrains domain behavior. Collect raw evidence without analyzing or organizing it. Every rule found becomes a row in a flat table.

#### Format

```markdown
# Rule Extraction

| Rule | Aggregate | Source | Location | Detail |
|------|-----------|--------|----------|--------|
| Name must not be empty | Task | Validation | CreateTask handler | Returns error if name is empty string |
| Name must not be empty | Task | Database | Schema/migration | NOT NULL constraint on name column |
| MaxLoad >= 0 | TimeBlock | Validation | UpdateTimeBlock | Guard clause checks maxLoad >= 0 |
| MaxLoad 0–1000 | TimeBlock | VO logic | BlockCapacity calculation | Caps effective capacity using maxLoad/100, maxLoad checked <= 1000 |
| Email unique | User | Database | Schema/migration | UNIQUE constraint on email column |
| Email unique | User | Validation | CreateUser | Checks for existing email before insert |
| Deleted tasks excluded | Task | Query | Task repository | WHERE deleted_at IS NULL on all reads |
| Cascade delete children | Task | Database | Schema/migration | ON DELETE CASCADE on parent_task_id FK |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Column Rules

- **Rule**: Short, declarative statement of what must be true. Use domain language, not code.
- **Aggregate**: Which aggregate this rule applies to. Use PascalCase entity name from classification.md.
- **Source**: Where this evidence was found. Values: `Validation`, `Database`, `Query`, `VO logic`, `API`, `Frontend`, `Test`, `Documentation`, `Configuration`.
- **Location**: Specific location in the codebase — function name, file area, or constraint name. Not a file path, but a descriptive location (e.g., "CreateTask handler", "task table schema", "task_test assertions").
- **Detail**: Brief explanation of how the rule is enforced at this location.

#### What to Look For

- Guard clauses and validation checks in handlers/services
- Database constraints: NOT NULL, UNIQUE, CHECK, FK, CASCADE
- Query filters that are always applied (soft delete filters, user scoping)
- Conditional logic that enforces state transitions
- Error messages that imply business rules
- Test assertions that verify expected behavior
- Frontend validation (input masks, field limits, required fields)
- Configuration or constants that define limits
- Comments or documentation that state rules

#### Structural Rules

- One row per rule per source — the same rule enforced in two places gets two rows
- Sort by Aggregate alphabetically, then by Rule alphabetically within each aggregate
- No code snippets or file paths — use descriptive locations
- Include rules even if they seem redundant with entity invariants — this is raw evidence collection

### Phase 2 — Rule Analysis

File: `docs/flow/ddd/business-rules.md`

Read `rules-extraction.md` and group rules by aggregate. Cross-reference enforcement locations to detect contradictions, gaps, and partially enforced rules.

#### Format

```markdown
# Business Rules

## Task

### Name must not be empty

| Source | Location | Stated Rule |
|--------|----------|-------------|
| Validation | CreateTask handler | Name must not be empty |
| Validation | UpdateTask handler | Name must not be empty |
| Database | Task table schema | NOT NULL on name |

**Status:** ✅ Consistent

---

### MaxLoad range

| Source | Location | Stated Rule |
|--------|----------|-------------|
| Validation | UpdateTimeBlock handler | MaxLoad >= 0, no upper bound |
| VO logic | BlockCapacity calculation | MaxLoad 0–1000 |

**Status:** ⚠️ Contradiction — validation allows unbounded values, capacity calculation caps at 1000

---

## Cross-Aggregate Rules

### Cascade delete on parent task

| Source | Location | Stated Rule | Aggregates |
|--------|----------|-------------|-----------|
| Database | Task table FK | ON DELETE CASCADE on parent_task_id | Task |
| Query | Task repository | Filters by parent when loading children | Task |

**Status:** ✅ Consistent

---

## Summary

| Status | Count |
|--------|-------|
| ✅ Consistent | 12 |
| ⚠️ Contradiction | 2 |
| ⚠️ Partial | 3 |
| ❌ Unenforced | 1 |

### Contradictions

| Rule | Aggregates | Issue |
|------|-----------|-------|
| MaxLoad range | TimeBlock | Validation allows unbounded, VO caps at 1000 |

### Gaps

| Rule | Aggregates | Issue |
|------|-----------|-------|
| API name validation | Task | Database enforces NOT NULL but POST handler accepts empty |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Rule Section Rules

Each rule section within an aggregate must have exactly these parts:

1. `### Rule name` — short declarative statement
2. Evidence table with columns `Source | Location | Stated Rule`
3. `**Status:**` line with one of: `✅ Consistent`, `⚠️ Contradiction`, `⚠️ Partial`, `❌ Unenforced`
4. For non-consistent statuses, a brief explanation after the status
5. `---` separator before next rule

#### Cross-Aggregate Rules Section

Rules that span multiple aggregates go in `## Cross-Aggregate Rules`. Same format but the evidence table has an additional `Aggregates` column.

#### Summary Section

The file must end with a `## Summary` section containing:
1. Status counts table
2. `### Contradictions` — table listing each contradiction with aggregates and issue
3. `### Gaps` — table listing each gap with aggregates and issue

If no contradictions or gaps exist, state "None found."

#### Status Definitions

- **✅ Consistent** — all sources agree on the rule and it is enforced everywhere it should be
- **⚠️ Contradiction** — different sources state different versions of the rule (e.g., different ranges, different allowed values)
- **⚠️ Partial** — rule is enforced in some layers but missing in others where it should be (e.g., validated in backend but not in database, or vice versa)
- **❌ Unenforced** — rule is documented, implied by naming, or stated in tests but not actually enforced in code

---

## Domain Commands

Domain command extraction uses the same two-phase approach as rules extraction.

### Phase 1 — Command Extraction

File: `docs/flow/ddd/commands-extraction.md`

Scan the codebase for every operation that changes aggregate state. Collect raw evidence: what the command does, which aggregate it targets, what parameters it takes, and where it's invoked.

#### Format

```markdown
# Command Extraction

| Command | Aggregate | Source | Location | Parameters | Outcome |
|---------|-----------|--------|----------|------------|---------|
| CreateTask | Task | API | POST handler | Name, Status, DueDate, ScheduledDate, ProjectId | Creates a new task with defaults applied |
| UpdateTaskStatus | Task | API | PATCH status handler | Status | Changes status, triggers time tracking side effects |
| UpdateTaskStatus | Task | Frontend | TasksView status picker | Status | Also sets timeTrackingStart/End based on transition direction |
| SoftDeleteTask | Task | API | DELETE handler | TaskId | Sets deleted_at timestamp, excludes from normal queries |
| RestoreTask | Task | API | POST restore handler | TaskId | Clears deleted_at, returns task to active list |
| ReorderTasks | Task | API | PUT reorder handler | TaskId, Order | Updates sort position within block or list |
| AssignStakeholders | Task | API | PUT stakeholders handler | TaskId, StakeholderIds | Replaces all stakeholder associations |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Column Rules

- **Command**: PascalCase verb-noun name describing the operation. Derive from the code — use the handler/function name as a guide, normalized to domain language.
- **Aggregate**: PascalCase entity name from classification.md.
- **Source**: Where this command is implemented: `API`, `Frontend`, `Backend service`, `Scheduler`, `Migration`, `Seed`.
- **Location**: Descriptive location — handler name, function name, component name. Not a file path.
- **Parameters**: Comma-separated list of input parameters in domain terms. Use PascalCase attribute names from entities.md.
- **Outcome**: Brief description of what the command does — what state changes.

#### What to Scan

- API route handlers (POST, PUT, PATCH, DELETE)
- Frontend action functions that call the API
- Background jobs or schedulers that modify state
- Migration scripts that transform data
- Seed scripts that create initial data
- Service methods that orchestrate multiple changes

#### Structural Rules

- One row per command per source — same command in API and frontend gets two rows
- Sort by Aggregate alphabetically, then by Command within aggregate
- Use domain language, not implementation names
- Include commands even if they seem like simple CRUD — the analysis phase decides what's meaningful

### Phase 2 — Command Analysis

File: `docs/flow/ddd/commands.md`

Read `commands-extraction.md` and group by aggregate. Cross-reference with entities.md Behavior lines to detect coverage gaps.

#### Format

```markdown
# Commands

## Task

### CreateTask

Creates a new task with a name and optional scheduling, project, and container associations.

| Source | Location | Parameters |
|--------|----------|------------|
| API | POST handler | Name, Status, DueDate, ScheduledDate, ProjectId, ParentContainerId, ParentTaskId |
| Frontend | TaskCreateModal | Name, Status, DueDate, ScheduledDate, ProjectId |
| Frontend | TasksView inline add | Name, ScheduledDate, ScheduledTimeBlockId |

**Behavior mapping:** Create
**Validation:** Name must not be empty; Status defaults to pending
**Side effects:** None

**Status:** ✅ Documented

---

### UpdateTaskStatus

Transitions a task to a new lifecycle status, triggering time tracking side effects.

| Source | Location | Parameters |
|--------|----------|------------|
| API | PATCH status handler | Status |
| Frontend | TasksView status picker | Status |

**Behavior mapping:** UpdateStatus
**Validation:** Status must be a valid TaskStatus value
**Side effects:** <ul><li>Pending → InProgress: sets TimeTrackingStart</li><li>InProgress → Completed: sets TimeTrackingEnd</li></ul>

**Status:** ✅ Documented

---

### BatchReorderTasks

Reorders multiple tasks within a block or list in a single operation.

| Source | Location | Parameters |
|--------|----------|------------|
| API | PUT reorder handler | List<TaskId, Order> |

**Behavior mapping:** ⚠️ None — no Behavior line in entities.md covers batch reorder

**Status:** ⚠️ Undocumented

---

## Summary

| Status | Count |
|--------|-------|
| ✅ Documented | 18 |
| ⚠️ Undocumented | 3 |
| ⚠️ Orphaned | 1 |
| ❌ Dead | 0 |

### Undocumented Commands

| Command | Aggregate | Issue |
|---------|-----------|-------|
| BatchReorderTasks | Task | API endpoint exists but no Behavior line in entities.md |

### Orphaned Behaviors

| Behavior | Aggregate | Issue |
|----------|-----------|-------|
| Cancel | Task | Listed in entities.md Behavior but no command implementation found |

### Dead Commands

None found.
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Command Section Rules

Each command section within an aggregate must have exactly these parts:

1. `### CommandName` — PascalCase verb+noun heading
2. One sentence description of what the command does in domain terms
3. Evidence table with columns `Source | Location | Parameters`
4. `**Behavior mapping:**` — Behavior name from entities.md, or `⚠️ None` if undocumented
5. `**Validation:**` — what is validated before execution
6. `**Side effects:**` — other changes triggered, using `<ul><li>` for multiple. `None` if none.
7. `**Status:**` line
8. `---` separator

#### Summary Section

The file must end with `## Summary` containing:

1. Status counts table — `Status | Count`
2. `### Undocumented Commands` — commands with no matching Behavior
3. `### Orphaned Behaviors` — Behavior lines with no matching command
4. `### Dead Commands` — unreachable code paths

If none exist in a category, state "None found."

#### Status Definitions

- **✅ Documented** — command exists in code and maps to a Behavior line in entities.md
- **⚠️ Undocumented** — command exists in code but no Behavior line covers it
- **⚠️ Orphaned** — Behavior line exists in entities.md but no command implementation found
- **❌ Dead** — code path exists but is unreachable or unused

---

## Domain Events

Domain event extraction uses the same two-phase approach as rules and commands extraction.

### Phase 1 — Event Extraction

File: `docs/flow/ddd/events-extraction.md`

Scan the codebase for every side effect, cascade, reaction, and published event that happens as a consequence of a command. Collect raw evidence: what happened, what triggered it, which aggregate is affected, and where the effect originates.

#### Format

```markdown
# Event Extraction

| Event | Trigger | Aggregate | Source | Location | Detail |
|-------|---------|-----------|--------|----------|--------|
| TaskCreated | CreateTask | Task | Frontend | tasks store createTask | logEvent fired with taskId, name, scheduledTimeBlockId, projectId |
| StakeholdersReplaced | CreateTask | Task | API | POST handler | task_stakeholders join rows deleted and reinserted in same transaction |
| TimeTrackingStarted | UpdateTask | Task | Frontend | TasksView status picker | Sets timeTrackingStart when status transitions to InProgress |
| ChildTasksCascadeDeleted | HardDeleteTask | Task | Database | parent_task_id FK CASCADE | ON DELETE CASCADE removes child rows |
| ProjectReferencesNulled | SoftDeleteProject | Project | Database | task project_id FK SET NULL | Task project_id set to null when project deleted |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Column Rules

- **Event**: PascalCase past-tense name describing what happened. `TaskCreated` not `CreateTask`.
- **Trigger**: PascalCase command name from commands.md that causes this event.
- **Aggregate**: PascalCase entity name from classification.md that is affected.
- **Source**: Where the side effect originates: `API`, `Frontend`, `Backend`, `Database`, `Scheduler`, `Migration`.
- **Location**: Descriptive location — handler name, FK constraint name, component name. NOT a file path.
- **Detail**: How the side effect manifests. One or two sentences.

#### What to Scan

- Explicit event publications (logEvent calls, event bus emissions, message queue publishes)
- Database cascades (ON DELETE CASCADE, ON DELETE SET NULL, ON UPDATE CASCADE)
- Transactional side effects (join table replacements, timestamp mutations, computed field updates)
- Frontend reactions to command completion (store mutations, refetches, navigation)
- Backend reactions (scheduled cleanup, background processing)
- Implicit state transitions (timestamp mutations during status changes)

#### Structural Rules

- One row per event per source — same event from Database and Frontend gets two rows
- Sort by Aggregate alphabetically, then by Event within aggregate
- Use past-tense domain language for Event names
- Trigger must match a Command name from commands-extraction.md
- No file paths — use descriptive locations

### Phase 2 — Event Analysis

File: `docs/flow/ddd/events.md`

Read `events-extraction.md` and group by aggregate. Cross-reference with commands.md to map triggers. Classify events as explicit, implicit, silent, or undocumented.

#### Format

```markdown
# Events

## Task

### TaskCreated

A new task has been inserted into the database with all associated stakeholder and tag joins.

| Trigger | Source | Location | Detail |
|---------|--------|----------|--------|
| CreateTask | Frontend | tasks store createTask | logEvent fired with taskId, name, scheduledTimeBlockId, projectId |
| CreateTask | API | POST handler | Stakeholder and tag join rows inserted in same transaction |

**Triggered by:** CreateTask
**Event type:** ✅ Explicit — published via logEvent
**Subscribers:** None — logEvent is fire-and-forget to EventLog; no other system reacts
**Cascade effects:** None

**Status:** ✅ Explicit

---

## Summary

| Status | Count |
|--------|-------|
| ✅ Explicit | 8 |
| ⚠️ Implicit | 12 |
| ⚠️ Silent | 5 |
| ❌ Undocumented | 2 |

### Implicit Events

| Event | Aggregate | Trigger | Issue |
|-------|-----------|---------|-------|
| TimeTrackingStarted | Task | UpdateTask | Timestamp set as side effect with no named event |

### Silent Cascades

| Event | Aggregate | Trigger | Issue |
|-------|-----------|---------|-------|
| ChildTasksCascadeDeleted | Task | HardDeleteTask | DB cascade with no application awareness |

### Undocumented Events

| Event | Aggregate | Trigger | Issue |
|-------|-----------|---------|-------|
| StakeholderOrphaned | Stakeholder | DeleteTask | Stakeholder may lose all task associations |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Event Section Rules

Each event section within an aggregate must have exactly these parts:

1. `### EventName` — PascalCase past-tense heading
2. One sentence description
3. Evidence table with columns `Trigger | Source | Location | Detail`
4. `**Triggered by:**` — which command(s) cause this event
5. `**Event type:**` — one of the four status values with brief explanation
6. `**Subscribers:**` — what reacts to this event, or `None`
7. `**Cascade effects:**` — further changes, using `<ul><li>` for multiple. `None` if none.
8. `**Status:**` line
9. `---` separator

#### Summary Section

The file must end with `## Summary` containing:

1. Status counts table — `Status | Count`
2. `### Implicit Events` — events that happen but aren't named
3. `### Silent Cascades` — DB-level cascades the application doesn't know about
4. `### Undocumented Events` — side effects not mentioned in documentation

If none exist in a category, state "None found."

#### Status Definitions

- **✅ Explicit** — event is named and published via an event mechanism (logEvent, event bus, message queue)
- **⚠️ Implicit** — side effect happens as part of a command but is not modeled as a named event
- **⚠️ Silent** — cascade happens at the database level with no application-layer awareness
- **❌ Undocumented** — side effect found in code but not mentioned in any documentation or entity model

---

## Bounded Contexts

Bounded context extraction uses the same two-phase approach. Unlike other concepts, contexts are interpretive — they propose architectural boundaries based on coupling evidence rather than documenting code directly.

### Phase 1 — Context Extraction

File: `docs/flow/ddd/contexts-extraction.md`

Scan the codebase and existing DDD files for coupling signals between aggregates — FK relationships, shared transactions, cross-aggregate events, shared concepts, and module boundaries.

#### Format

```markdown
# Context Extraction

| Signal | Type | From | To | Source | Location | Detail |
|--------|------|------|----|--------|----------|--------|
| Task references Project via FK | Data coupling | Task | Project | Database | task.project_id FK | ON DELETE SET NULL; Task can exist without Project |
| CreateTask syncs stakeholders in same transaction | Command coupling | Task | Stakeholder | API | POST handler | task_stakeholders written in same transaction as task INSERT |
| SoftDeleteProject refetches task list | Event coupling | Project | Task | Frontend | projects store deleteProject | Task list refetched after project soft-delete |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Column Rules

- **Signal**: Short description of the coupling found.
- **Type**: `Data coupling`, `Command coupling`, `Event coupling`, `Shared concept`, `Shared table`, `Dead module`, `Hub dependency`, `Circular dependency`.
- **From**: PascalCase aggregate that initiates the coupling.
- **To**: PascalCase aggregate that is coupled to. `—` if about a single aggregate.
- **Source**: `Database`, `API`, `Frontend`, `Backend`, `Entity model`, `Configuration`, `Migration`.
- **Location**: Descriptive location. NOT a file path.
- **Detail**: How the coupling manifests. One or two sentences.

### Phase 2 — Context Analysis

File: `docs/flow/ddd/contexts.md`

Read `contexts-extraction.md` and group aggregates into proposed bounded contexts. Map coupling between contexts and assess boundary quality.

#### Format

```markdown
# Bounded Contexts

## Task Management

The core scheduling and execution context. Owns task lifecycle, time tracking, and hierarchical structure.

**Aggregates:** Task, TimeBlock

**Rationale:** Task and TimeBlock are tightly coupled through scheduling fields and share drag-and-drop interactions.

| Signal | Type | Detail |
|--------|------|--------|
| Task references TimeBlock via FK | Data coupling | task.scheduled_time_block_id FK; ON DELETE SET NULL |

**Internal boundary quality:** ✅ Clean — aggregates interact through well-defined scheduling fields

---

## Context Map

### Project Organization ↔ Task Management

| Direction | Signal | Type | Detail |
|-----------|--------|------|--------|
| Task Management → Project Organization | Task references Project via FK | Data coupling | task.project_id FK |

**Boundary status:** ⚠️ Leaky boundary — Task writes project_id directly

---

## Summary

### Proposed Contexts

| Context | Aggregates | Internal Quality |
|---------|-----------|-----------------|
| Task Management | Task, TimeBlock | ✅ Clean |

### Context Coupling

| From | To | Coupling Count | Boundary Status |
|------|----|---------------|-----------------|
| Task Management | Project Organization | 3 | ⚠️ Leaky boundary |

### Boundary Issues

| Issue | Contexts | Detail |
|-------|----------|--------|
| Direct FK write | Task Management → Project Organization | Task writes project_id directly |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Context Section Rules

Each context section must have exactly these parts:

1. `## ContextName` — human-readable subdomain name
2. One or two sentence description
3. `**Aggregates:**` — comma-separated PascalCase aggregate names
4. `**Rationale:**` — why these aggregates belong together
5. Evidence table with columns `Signal | Type | Detail`
6. `**Internal boundary quality:**` — one of the status values
7. `---` separator

#### Status Definitions

- **✅ Clean boundary** — contexts interact only through well-defined interfaces
- **⚠️ Leaky boundary** — contexts share transactions, write each other's fields, or have implicit knowledge
- **⚠️ Missing boundary** — logically separate aggregates are tightly coupled with no separation
- **❌ Circular dependency** — bidirectional command or data coupling

---

## Domain Services

Domain service extraction uses the same two-phase approach. Services capture stateless operations that orchestrate, coordinate, or compute across aggregate boundaries — logic that doesn't naturally belong to a single entity or value object.

### Phase 1 — Service Extraction

File: `docs/flow/ddd/services-extraction.md`

Scan the codebase for multi-aggregate orchestration, transactional coordination, cross-aggregate computations, and business policy logic that spans aggregate boundaries.

#### Format

```markdown
# Service Extraction

| Service | Type | Aggregates | Source | Location | Detail |
|---------|------|-----------|--------|----------|--------|
| RebuildTimeBlocks | Orchestration | Task, TimeBlock | Frontend | TimeBlocks component rebuild | Clears overdue dates, computes capacity, distributes tasks across blocks on page load |
| SyncStakeholders | Coordination | Task, Stakeholder | API | POST/PUT handler transaction | Delete-all then reinsert task_stakeholders in same transaction as Task write |
| ComputeBlockCapacity | Computation | Task, TimeBlock | Frontend | TimeBlocks blockCapacity function | Counts active tasks per block against MaxLoad-scaled ceiling |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Column Rules

- **Service**: PascalCase verb+noun name.
- **Type**: `Orchestration`, `Coordination`, `Computation`, or `Policy`.
- **Aggregates**: Comma-separated PascalCase aggregate names.
- **Source**: `API`, `Frontend`, `Backend`, `Database`, or `Migration`.
- **Location**: Descriptive location. NOT a file path.
- **Detail**: What the service does. One or two sentences.

#### Type Definitions

- **Orchestration** — Multi-step coordination across aggregates with sequencing or conditional logic.
- **Coordination** — Transactional consistency across aggregates (shared transactions, compensating actions).
- **Computation** — Derived values spanning aggregate boundaries.
- **Policy** — Business decision logic that doesn't belong to a single entity.

### Phase 2 — Service Analysis

File: `docs/flow/ddd/services.md`

Read `services-extraction.md` and group services by owning bounded context from contexts.md. Assess encapsulation quality.

#### Format

```markdown
# Services

## TaskManagement

### RebuildTimeBlocks

Orchestrates the full TimeBlock display lifecycle.

| Aggregates | Source | Location |
|-----------|--------|----------|
| Task, TimeBlock | Frontend | TimeBlocks component rebuild |

**Type:** Orchestration
**Owning context:** Task Management
**Trigger:** Page load, task data change
**Steps:** <ul><li>Identify overdue tasks</li><li>Clear or surface overdue tasks</li><li>Compute block capacity</li></ul>
**Encapsulation:** ⚠️ Scattered — orchestration in Vue component

**Status:** ⚠️ Scattered

---

## Cross-Context Services

### FilterSystemProjects

Applies Project context rules within Task Management UI.

| Aggregates | Source | Location |
|-----------|--------|----------|
| Project, Task | Frontend | TaskCreateModal project dropdown |

**Type:** Policy
**Owning context:** Crosses Task Management → Work Organization
**Trigger:** TaskCreateModal render
**Steps:** <ul><li>Read projects</li><li>Filter where isSystem is false</li><li>Display in dropdown</li></ul>
**Encapsulation:** ❌ Boundary violation

**Status:** ❌ Boundary violation

---

## Summary

### Service Inventory

| Service | Type | Context | Aggregates | Status |
|---------|------|---------|-----------|--------|

### Encapsulation Issues

| Service | Issue |
|---------|-------|

### Service Type Distribution

| Type | Count |
|------|-------|
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

#### Service Section Rules

Each service section must have exactly these parts:

1. `### ServiceName` — PascalCase verb+noun
2. One sentence description
3. Evidence table: `| Aggregates | Source | Location |`
4. `**Type:**` — Orchestration, Coordination, Computation, or Policy
5. `**Owning context:**` — context from contexts.md, or `Crosses A → B`
6. `**Trigger:**` — what causes execution
7. `**Steps:**` — algorithm, using `<ul><li>`
8. `**Encapsulation:**` — status with explanation
9. `**Status:**` — one of: `✅ Encapsulated`, `⚠️ Scattered`, `⚠️ Misplaced`, `❌ Boundary violation`
10. `---` separator

---

| Command | Purpose |
|---------|---------|
| `/ccf-flow-ddd-extract-all` | Full pipeline: UL → classify → entities → VOs → aggregates → rules → commands → events → contexts → services → xref verify |
| `/ccf-flow-ddd-extract <concept>` | Re-extract a single concept: `ul`, `classification`, `entities`, `value-objects`, `aggregates`, `rules`, `commands`, `events`, `contexts`, `services` |
| `/ccf-flow-ddd-xref-verify` | Verify consistency across all DDD files after manual edits |

## Sub-Agents

These are invoked automatically by the commands above. Users do not need to call them directly.

| Agent | Role |
|-------|------|
| `ccf-flow-ddd-ul-extractor` | Produces ubiquitous-language.md |
| `ccf-flow-ddd-ul-verifier` | Checks ubiquitous-language.md format |
| `ccf-flow-ddd-classifier` | Phase 1: produces classification.md |
| `ccf-flow-ddd-entity-extractor` | Phase 2: produces entities.md |
| `ccf-flow-ddd-entity-verifier` | Checks entities.md format |
| `ccf-flow-ddd-vo-extractor` | Phase 2: produces value-objects.md |
| `ccf-flow-ddd-vo-verifier` | Checks value-objects.md format |
| `ccf-flow-ddd-agg-extractor` | Phase 2: produces aggregates.md |
| `ccf-flow-ddd-agg-verifier` | Checks aggregates.md format |
| `ccf-flow-ddd-rules-extractor` | Phase 1: produces rules-extraction.md |
| `ccf-flow-ddd-rules-analyzer` | Phase 2: produces business-rules.md |
| `ccf-flow-ddd-rules-verifier` | Checks business-rules.md format |
| `ccf-flow-ddd-commands-extractor` | Phase 1: produces commands-extraction.md |
| `ccf-flow-ddd-commands-analyzer` | Phase 2: produces commands.md |
| `ccf-flow-ddd-commands-verifier` | Checks commands.md format |
| `ccf-flow-ddd-events-extractor` | Phase 1: produces events-extraction.md |
| `ccf-flow-ddd-events-analyzer` | Phase 2: produces events.md |
| `ccf-flow-ddd-events-verifier` | Checks events.md format |
| `ccf-flow-ddd-contexts-extractor` | Phase 1: produces contexts-extraction.md |
| `ccf-flow-ddd-contexts-analyzer` | Phase 2: produces contexts.md |
| `ccf-flow-ddd-contexts-verifier` | Checks contexts.md format |
| `ccf-flow-ddd-services-extractor` | Phase 1: produces services-extraction.md |
| `ccf-flow-ddd-services-analyzer` | Phase 2: produces services.md |
| `ccf-flow-ddd-services-verifier` | Checks services.md format |
| `ccf-flow-ddd-xref-verifier` | Cross-references all files for consistency |

## Verification

After creating or updating any DDD file, always run `/ccf-flow-ddd-xref-verify` to check cross-file consistency.

## Full Template References

See `references/ddd-ul-template.md` for ubiquitous language specification.
See `references/ddd-entity-template.md` for entity specification.
See `references/ddd-vo-template.md` for value object specification.
See `references/ddd-classification-template.md` for classification specification.
See `references/ddd-agg-template.md` for aggregate specification.
See `references/ddd-rules-extraction-template.md` for rule extraction specification.
See `references/ddd-business-rules-template.md` for business rules specification.
See `references/ddd-commands-extraction-template.md` for command extraction specification.
See `references/ddd-commands-template.md` for domain commands specification.
See `references/ddd-events-extraction-template.md` for event extraction specification.
See `references/ddd-events-template.md` for domain events specification.
See `references/ddd-contexts-extraction-template.md` for context extraction specification.
See `references/ddd-contexts-template.md` for bounded contexts specification.
See `references/ddd-services-extraction-template.md` for service extraction specification.
See `references/ddd-services-template.md` for domain services specification.
