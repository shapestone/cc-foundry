# Value Object Template

## File Location

`docs/flow/ddd/value-objects.md`

This file contains value objects ONLY. Entities go in `entities.md`. Aggregates go in `aggregates.md`.

## Format

The file has a `# Value Objects` heading, then one `## VOName` section per value object. No other sections.

Each section has exactly these parts in this order:

1. `## VOName` — heading
2. One sentence description
3. Attributes table with columns `Attribute | Type | Description`
4. `**Constraints:**` followed by `<ul><li>` list
5. `---` separator before next value object

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

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Attribute | PascalCase name. No identity field — value objects have no ID. | Never empty. |
| Type | Domain type: String, Decimal, Date, Boolean, Enum, or reference to another domain type. Use `Enum` for enumerated types. No SQL types, no language types. | Never empty. |
| Description | Short description. For Enum types, list valid values using `<ul><li>`. | Never empty. |

## Structural Rules

- Value objects sorted alphabetically by name
- Sections separated by `---`
- No content before the first VO other than `# Value Objects`
- No sections other than `## VOName` — no `## Entities`, `## Aggregates`, `## Notes`
- No implementation details

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`)
- Database details (`table`, `column`, `deleted_at`, `FK`, `migration`, `index`)
- SQL types (`TEXT`, `INTEGER`, `TIMESTAMP`, `REAL`, `DATETIME`)
- HTTP details (`GET /api`, `POST`, `HTTP 400`, `endpoint`)
- Code syntax (backtick-wrapped code, `function()`, `logEvent(`)
- Sections for entities or aggregates
- `## Notes` or `## Derived` sections
- Identity fields (value objects have no ID)

## What Is a Value Object

A value object is immutable, has no identity, and is compared by its attributes. It includes:

- Enums and type classifications (TaskStatus, StakeholderType, ContainerType)
- Small composite values (Money, DateRange)
- Computed/derived values that are not persisted (DerivedStatus, BlockCapacity, StatBar)

If a concept has a unique ID and mutable state, it is an entity and belongs in `entities.md`.
