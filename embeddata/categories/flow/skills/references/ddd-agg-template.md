# Aggregate Template

## File Location

`docs/flow/ddd/aggregates.md`

This file contains aggregates ONLY. Entity details go in `entities.md`. Value object details go in `value-objects.md`. Read `classification.md` first for entity/VO assignments.

## Format

The file has a `# Aggregates` heading, then one `## AggregateName` section per aggregate. No other sections.

Each section has exactly these parts in this order:

1. `## AggregateName` — heading
2. One sentence description — what the aggregate represents and its consistency scope
3. Boundary table with columns `Element | Type | Relationship`
4. `**Invariants:**` followed by `<ul><li>` list
5. `---` separator before next aggregate

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

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Element | PascalCase name of entity or reference within the aggregate. | Never empty. |
| Type | `Root Entity`, `Entity`, `Entity (ref)`, or `Value Object`. Use `(ref)` for cross-boundary references. | Never empty. |
| Relationship | How this element relates to the root: `M:N via join`, `One-to-many`, `read-only`, etc. | `—` for the root itself. |

## Structural Rules

- Aggregates sorted alphabetically by name
- Sections separated by `---`
- No entity attribute details — those live in `entities.md`
- No value object details — those live in `value-objects.md`
- No sections other than `## AggregateName`
- No implementation details

## Forbidden Content

- File paths, table names, column names, SQL types, HTTP details, code syntax
- Entity attribute lists (belong in `entities.md`)
- Value object definitions (belong in `value-objects.md`)
- `## Entities`, `## Value Objects`, `## Notes` sections
