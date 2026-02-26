# Commands Template

## File Location

`docs/flow/ddd/commands.md`

This file is Phase 2 of command extraction. It reads `commands-extraction.md`, groups commands by aggregate, maps them to documented Behavior lines in entities.md, and detects undocumented commands, orphaned behaviors, and dead code.

## Format

The file has a `# Commands` heading (NOT `# Domain Commands`), then `## AggregateName` sections (one per aggregate), then `## Summary`.

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
**Side effects:** <ul><li>Pending → InProgress: sets TimeTrackingStart</li><li>InProgress → Completed: sets TimeTrackingEnd</li><li>Completed → InProgress: clears TimeTrackingEnd</li><li>Any → Pending: clears both tracking timestamps</li></ul>

**Status:** ✅ Documented

---

### BatchReorderTasks

Reorders multiple tasks within a block or list in a single operation.

| Source | Location | Parameters |
|--------|----------|------------|
| API | PUT reorder handler | List<TaskId, Order> |

**Behavior mapping:** ⚠️ None — no Behavior line in entities.md covers batch reorder
**Validation:** Each task must exist for the user
**Side effects:** Updates order and optionally project_id for each task atomically

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

## Command Section Rules

Each command section within an aggregate must have EXACTLY these parts in EXACTLY this order:

1. `### CommandName` — PascalCase verb+noun heading
2. One sentence description — what the command does in domain terms
3. Evidence table with columns `Source | Location | Parameters`
4. `**Behavior mapping:**` — which Behavior line from entities.md this maps to, or `⚠️ None` if undocumented
5. `**Validation:**` — what the command validates before executing
6. `**Side effects:**` — what else changes as a result, using `<ul><li>` for multiple. `None` if no side effects.
7. `**Status:**` line — exactly one of the four status values
8. `---` separator

These are the ONLY permitted fields. Do not add or rename fields.

## Forbidden Field Names

Do NOT use any of these field names — they are not part of the format:

- `**Intent:**` — use the one-sentence description instead
- `**Inputs:**` — use the `Source | Location | Parameters` evidence table instead
- `**Preconditions:**` — fold into `**Validation:**`
- `**Enforcement layers:**` — fold into `**Validation:**`
- `**Gaps:**` — note validation gaps in `**Validation:**` and side effect gaps in `**Side effects:**`

## Forbidden Table Formats

Do NOT use any of these table column layouts:

- `Field | Required | Notes` — this is an API spec format, not the commands format
- `Field | Evidence` — this is a key-value format, not the commands format
- `Layer | Rule` — this is a rules format, not the commands format

The ONLY table in each command section is `Source | Location | Parameters`.

## Status Definitions

| Status | Meaning |
|--------|---------|
| ✅ Documented | Command exists in code and maps to a Behavior line in entities.md |
| ⚠️ Undocumented | Command exists in code but no Behavior line covers it |
| ⚠️ Orphaned | Behavior line exists in entities.md but no command implementation found |
| ❌ Dead | Code path exists but is unreachable, commented out, or behind a permanently false flag |

## Summary Section

The file must end with `## Summary` containing:

1. Status counts table
2. `### Undocumented Commands` — commands with no matching Behavior
3. `### Orphaned Behaviors` — Behavior lines with no matching command
4. `### Dead Commands` — unreachable code paths

If none exist in a category, state "None found."

Do NOT add a `## Cross-Cutting Gaps` or `## Observations` section. Gaps are noted within each command's `**Validation:**` and `**Side effects:**` fields.

## Structural Rules

- Aggregate sections sorted alphabetically
- Commands within aggregates sorted alphabetically
- Summary always last
- Domain language in command names — PascalCase verb+noun
- Parameters use domain types from entities.md and value-objects.md
- No code snippets or file paths

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`, `stores/`)
- HTTP methods as command names
- Code snippets or backtick-wrapped code
- `# Domain Commands` heading (must be `# Commands`)
- `## Summary Table` before the aggregate sections
- `## Cross-Cutting Gaps` or `## Observations` sections
- Per-command `Field | Required | Notes` input tables
- Raw extraction data without analysis
