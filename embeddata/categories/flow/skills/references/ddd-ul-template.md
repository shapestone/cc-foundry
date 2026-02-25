# Ubiquitous Language Template

## File Location

`docs/flow/ddd/ubiquitous-language.md`

## Format

The file must be a single markdown table with exactly 4 columns:

```markdown
# Ubiquitous Language

| Term | Aliases | Meaning | Notes |
|------|---------|---------|-------|
| Order | Purchase Order, PO | A confirmed request by a Customer to buy one or more Products at agreed prices. | "Order" in Shipping context means shipment request. |
| Money | — | A monetary amount with a currency code. | — |
```

## Column Rules

| Column | Required | Content | Empty Value |
|--------|----------|---------|-------------|
| Term | Yes | PascalCase canonical name as it appears in domain models. | Never empty. |
| Aliases | No | Comma-separated alternative names found in code, UI, database, or docs. | `—` |
| Meaning | Yes | One sentence describing the domain concept. No implementation details. | Never empty. |
| Notes | No | Brief disambiguation or edge case only. | `—` |

## Structural Rules

- File starts with `# Ubiquitous Language`
- Exactly one table in the file
- Rows sorted alphabetically by Term
- One row per term, no sub-rows or multi-line cells
- No content outside the heading and table (no preamble paragraphs, no footnotes)

## Forbidden Content

The following must never appear anywhere in the file:

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`)
- Table or column names (`tasks table`, `deleted_at`, `FK`, `migration`)
- SQL types (`TEXT`, `INTEGER`, `TIMESTAMP`, `NULL`, `CASCADE`)
- HTTP details (`GET /api`, `POST`, `HTTP 400`, `endpoint`)
- Code syntax (backtick-wrapped code, `function()`, `logEvent(`)
- Line numbers or file locations
