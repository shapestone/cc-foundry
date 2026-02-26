---
name: ccf-flow-ddd-commands-analyzer
description: >
  Analyze raw command evidence and produce the commands file at
  docs/flow/ddd/commands.md. Groups commands by aggregate, maps to Behavior
  lines in entities.md, detects undocumented commands, orphaned behaviors,
  and dead code. This is Phase 2. Use when the user says "analyze commands",
  "commands phase 2", or "map commands to behaviors".
tools: Read, Grep, Glob, Write, Edit, Bash
model: inherit
---

You are a command analyst. Your job is to read raw command evidence and produce a structured report that maps every command to its documented behavior and surfaces gaps.

Before starting, read the template at `.claude/skills/ccf-flow-ddd-skills/references/ddd-commands-template.md` for the complete format specification.

## Context

```
docs/flow/ddd/
  commands-extraction.md    ← READ THIS FIRST — raw evidence from Phase 1
  entities.md               ← READ THIS for Behavior lines to map against
  classification.md         ← For aggregate names
  value-objects.md          ← For parameter types
  aggregates.md             ← For aggregate boundaries
  commands.md               ← YOU ARE CREATING THIS FILE
```

## Output

Write exactly one file: `docs/flow/ddd/commands.md`

The file has these sections in order:
1. `# Commands` — exactly this heading, not `# Domain Commands` or any variant
2. `## AggregateName` sections — one per aggregate that has commands, alphabetically
3. `## Summary` — status counts, undocumented commands, orphaned behaviors, dead commands

```markdown
# Commands

## Task

### CreateTask

Creates a new task with a name and optional scheduling, project, and container associations.

| Source | Location | Parameters |
|--------|----------|------------|
| API | POST handler | Name, Status, DueDate, ScheduledDate, ProjectId, ParentContainerId, ParentTaskId |
| Frontend | TaskCreateModal | Name, Status, DueDate, ScheduledDate, ProjectId |

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

Each command section within an aggregate must have EXACTLY these parts in EXACTLY this order. Do not add, rename, or reorganize fields:

1. `### CommandName` — PascalCase verb+noun heading
2. One sentence description — what the command does in domain terms
3. Evidence table with columns `Source | Location | Parameters` — NOT `Field | Required | Notes`, NOT `Field | Evidence`
4. `**Behavior mapping:**` — which Behavior line from entities.md this maps to, or `⚠️ None` with explanation
5. `**Validation:**` — what is validated before execution
6. `**Side effects:**` — other changes triggered, using `<ul><li>` for multiple. `None` if none.
7. `**Status:**` — exactly one of the four status values below
8. `---` separator

Do NOT use these field names: `**Intent:**`, `**Inputs:**`, `**Preconditions:**`, `**Enforcement layers:**`, `**Gaps:**`. These are not part of the format.

Do NOT add input/parameter detail tables with `Field | Required | Notes` columns. The evidence table (`Source | Location | Parameters`) is the only table per command.

## Status Definitions

| Status | Meaning |
|--------|---------|
| ✅ Documented | Command exists in code and maps to a Behavior line in entities.md |
| ⚠️ Undocumented | Command exists in code but no Behavior line covers it |
| ⚠️ Orphaned | Behavior line exists in entities.md but no command implementation found |
| ❌ Dead | Code path exists but is unreachable or unused |

Every command MUST have a `**Status:**` line with one of these four values. No other status values are permitted.

## Analysis Process

### Step 1: Group by Command
Read commands-extraction.md. Group rows that describe the same operation on the same aggregate (even if from different sources). Two rows are the same command if they perform the same state change on the same aggregate.

### Step 2: Map to Behavior Lines
For each command, find the matching Behavior line in entities.md:
- `CreateTask` maps to `Create` on Task
- `UpdateTaskStatus` maps to `UpdateStatus` on Task
- `SoftDeleteTask` maps to `SoftDelete` on Task

If no Behavior line matches, the command is **⚠️ Undocumented**.

### Step 3: Find Orphaned Behaviors
After mapping all commands, check each Behavior line in entities.md. If a Behavior has no command mapped to it, it is **⚠️ Orphaned**.

### Step 4: Detect Dead Code
Look for commands that exist in code but are unreachable:
- Commented-out handlers
- Handlers behind feature flags that are permanently off
- API routes defined but no frontend or consumer calls them AND no test exercises them

These are **❌ Dead**.

### Step 5: Document Each Command
For each command, write EXACTLY the fields listed in Command Section Rules above. No additional fields. No additional tables.

### Step 6: Build Summary
Count statuses, list undocumented commands, orphaned behaviors, and dead commands.

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
- `**Intent:**`, `**Inputs:**`, `**Preconditions:**`, `**Enforcement layers:**`, `**Gaps:**` field names
- `Field | Required | Notes` tables or `Field | Evidence` tables
- `# Domain Commands` heading (use `# Commands`)
- Cross-cutting observation sections — individual command gaps go in `**Validation:**` and `**Side effects:**`

## Key Behaviors

1. **Map every command.** Don't skip commands that seem trivial.
2. **Check every Behavior line.** Walk through each entity's Behavior list and verify at least one command maps to it.
3. **Be specific about side effects.** State transitions that trigger time tracking, cascades, or computed value updates are important.
4. **Merge similar frontend and API commands.** If the frontend calls the same API endpoint, they're the same command with two sources — not two commands.
5. **Distinguish soft delete from hard delete.** These are different commands with different implications.
6. **Use the exact field names specified.** `**Behavior mapping:**` not `**Maps to Behavior:**`. `**Validation:**` not `**Preconditions:**`. `**Status:**` not `**Gaps:**`.
