# DDD Classification Template

## File Location

`docs/flow/ddd/classification.md`

This file is produced FIRST, before entities, value objects, or aggregates. It captures every cross-cutting decision so the per-file extractors stay consistent.

## Format

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

## Column Rules — Entities Table

| Column | Content | Empty Value |
|--------|---------|-------------|
| Name | PascalCase entity name. | Never empty. |
| Aggregate Root | `Yes` or `No`. | Never empty. |
| Contains | Entities/refs within this aggregate boundary, using `<ul><li>`. Indicate relationship type. | `—` |
| Uses VOs | Value objects used as attribute types, using `<ul><li>`. | `—` |

## Column Rules — Value Objects Table

| Column | Content | Empty Value |
|--------|---------|-------------|
| Name | PascalCase VO name. | Never empty. |
| Kind | `Enum`, `Composite`, or `Derived`. | Never empty. |
| Owned By | Entity that primarily uses this VO. | `—` if shared/computed. |

## Structural Rules

- Both tables sorted alphabetically by Name
- No content other than the heading and two tables
- No implementation details: no file paths, table names, SQL types, HTTP endpoints, or code syntax
- Every domain concept should appear in exactly one table — either Entities or Value Objects

## Forbidden Content

- File paths, table names, column names, SQL types, HTTP details, code syntax
- `## Aggregates` section (aggregate boundaries are derived from the Entities table)
- `## Notes` or prose sections
