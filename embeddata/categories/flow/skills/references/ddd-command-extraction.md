# Command Extraction Template

## File Location

`docs/flow/ddd/commands-extraction.md`

This file is Phase 1 of command extraction. It captures raw evidence of every state-changing operation found in the codebase. The commands analyzer reads this to produce the final commands.md.

## Format

The file has a `# Command Extraction` heading and a single flat table. Nothing else. No prose, no analysis, no sections, no grouping, no summary tables, no per-command subsections, no cascade observations.

```markdown
# Command Extraction

| Command | Aggregate | Source | Location | Parameters | Detail |
|---------|-----------|--------|----------|------------|--------|
| CreateTask | Task | API | POST handler | Name, Status, DueDate, ScheduledDate, ProjectId | Validates name not empty, defaults status to pending |
| CreateTask | Task | Frontend | TaskCreateModal submit | Name, Status, DueDate, ScheduledDate, ProjectId | Disables submit while name is empty |
| UpdateTaskStatus | Task | API | PATCH status handler | Status | Validates against TaskStatus enum |
| UpdateTaskStatus | Task | Frontend | TasksView status picker | Status | Also sets timeTrackingStart/End based on transition |
| SoftDeleteTask | Task | API | DELETE handler | TaskId | Soft delete — sets deleted_at timestamp |
| ReorderTasks | Task | API | PUT reorder handler | List<TaskId, Order> | Batch update of order values |
```

THIS IS THE ONLY ACCEPTABLE FORMAT.

The output is ONE table under ONE heading. Not multiple tables. Not sections per command. Not sections per aggregate.

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Command | PascalCase verb+noun name describing the operation. Infer from the code — use domain language. `CreateTask` not `POST task`. | Never empty. |
| Aggregate | PascalCase entity name from classification.md. | Never empty. |
| Source | Where the command is initiated: `API`, `Frontend`, `Backend`, `Scheduler`, `Migration`, `Seed`. | Never empty. |
| Location | Descriptive location — handler name, component name, function. NOT a file path. `POST handler` not `task_handlers.go`. | Never empty. |
| Parameters | Comma-separated list of input parameters in PascalCase domain types. | `—` if no parameters. |
| Detail | Brief description of what the command does, including side effects, defaults, and validation. One or two sentences. | Never empty. |

## What to Look For

- HTTP handlers that create, update, or delete records (POST, PUT, PATCH, DELETE)
- Frontend functions that call mutation APIs
- Backend services that modify aggregate state
- Scheduled jobs or cron tasks that change data
- Migration scripts that seed or transform data
- Event handlers that react to changes by making further changes

## What to Exclude

- Pure read operations (GET, list, search, filter) — these are queries, not commands
- Logging or telemetry that doesn't change domain state
- Infrastructure concerns (connection pooling, caching)

## Structural Rules

- One row per command per source — same command invoked from API and Frontend gets two rows
- Sort by Aggregate alphabetically, then by Command within each aggregate
- Use domain language for Command names, not HTTP verbs or function names
- No file paths — use descriptive locations
- Parameters use domain types from entities.md and value-objects.md
- Detail column is one or two sentences, not multi-paragraph

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`, `stores/`, `repository/`)
- HTTP methods as command names (use `CreateTask` not `POST task`)
- Code snippets or backtick-wrapped code
- Per-command sections with `### CommandName` headings — the output is a flat table
- Multiple tables — there is exactly one table
- Summary tables, cascade tables, or observation sections
- `Field | Evidence` key-value tables — use the specified column format
- Analysis or status judgments — that's Phase 2
