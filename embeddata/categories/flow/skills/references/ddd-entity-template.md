# Entity Template

## File Location

`docs/flow/ddd/entities.md`

This file contains entities ONLY. Value objects go in `value-objects.md`. Aggregates go in `aggregates.md`.

## Format

The file has a `# Entities` heading, then one `## EntityName` section per entity. No other sections.

Each section has exactly these parts in this order:

1. `## EntityName` — heading
2. One sentence description
3. Attributes table with columns `Attribute | Type | Description`
4. `**Behavior:**` followed by `<ul><li>` list
5. `**Invariants:**` followed by `<ul><li>` list
6. `---` separator before next entity

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

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Attribute | PascalCase name. First row is always the identity field. | Never empty. |
| Type | Domain type: UUID, String, Date, Decimal, Boolean, or reference to domain type (TaskStatus, ProjectId). No SQL types, no language types. | Never empty. |
| Description | Short description. Prefix with "Required" or "Optional". Use "Reference by ID only" for cross-aggregate refs. | Never empty. |

## Structural Rules

- Entities sorted alphabetically by name
- Sections separated by `---`
- No content before the first entity other than `# Entities`
- No sections other than `## EntityName` — no `## Value Objects`, `## Aggregates`, `## Notes`
- No implementation details: no file paths, table names, SQL column names, HTTP endpoints, migration numbers, or code syntax

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`)
- Database details (`table`, `column`, `deleted_at`, `FK`, `migration`, `index`)
- SQL types (`TEXT`, `INTEGER`, `TIMESTAMP`, `REAL`, `DATETIME`)
- HTTP details (`GET /api`, `POST`, `HTTP 400`, `endpoint`)
- Code syntax (backtick-wrapped code, `function()`, `logEvent(`)
- Sections for value objects or aggregates
- `## Notes` or `## Derived and Computed Attributes` sections

## What Is an Entity

An entity has identity (a unique ID), mutable state, and a lifecycle. If a concept is immutable or has no identity, it is a value object and belongs in `value-objects.md`, not here.
