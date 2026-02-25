# Rule Extraction Template

## File Location

`docs/flow/ddd/rules-extraction.md`

This file is Phase 1 of business rules extraction. It captures raw evidence of every rule found in the codebase without analysis or organization. The rules analyzer reads this to produce the final business-rules.md.

## Format

The file has a `# Rule Extraction` heading and a single flat table. Nothing else.

```markdown
# Rule Extraction

| Rule | Aggregate | Source | Location | Detail |
|------|-----------|--------|----------|--------|
| Name must not be empty | Task | Validation | CreateTask handler | Returns error if name is empty string |
| Name must not be empty | Task | Database | Task table schema | NOT NULL constraint on name column |
| MaxLoad >= 0 | TimeBlock | Validation | UpdateTimeBlock handler | Guard clause checks maxLoad >= 0 |
| MaxLoad 0–1000 | TimeBlock | VO logic | BlockCapacity calculation | Caps effective capacity, maxLoad checked <= 1000 |
| Email unique | User | Database | User table schema | UNIQUE constraint on email column |
| Email unique | User | Validation | CreateUser handler | Checks for existing email before insert |
| Deleted tasks excluded | Task | Query | Task repository | WHERE deleted_at IS NULL on all reads |
| Cascade delete children | Task | Database | Task table FK | ON DELETE CASCADE on parent_task_id |
```

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Rule | Short declarative statement in domain language. Not code. | Never empty. |
| Aggregate | PascalCase entity name from classification.md. | Never empty. |
| Source | Where evidence was found: `Validation`, `Database`, `Query`, `VO logic`, `API`, `Frontend`, `Test`, `Documentation`, `Configuration`. | Never empty. |
| Location | Descriptive location — function name, constraint name, file area. Not a file path. | Never empty. |
| Detail | Brief explanation of how the rule is enforced at this location. | Never empty. |

## What to Look For

- Guard clauses and validation checks in handlers/services
- Database constraints: NOT NULL, UNIQUE, CHECK, FK, CASCADE
- Query filters always applied (soft delete, user scoping)
- Conditional logic enforcing state transitions
- Error messages implying business rules
- Test assertions verifying expected behavior
- Frontend validation (input masks, field limits, required fields)
- Configuration or constants defining limits

## Structural Rules

- One row per rule per source — same rule in two places gets two rows
- Sort by Aggregate alphabetically, then by Rule within each aggregate
- No code snippets or file paths
- Include rules even if they seem redundant with entity invariants

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`)
- Code snippets or backtick-wrapped code
- `## Section` headings — the file is one heading and one table
- Analysis, grouping, or status judgments — that's Phase 2
