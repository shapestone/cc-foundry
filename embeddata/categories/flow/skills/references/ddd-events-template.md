# Domain Events Template

## File Location

`docs/flow/ddd/events.md`

This file is Phase 2 of domain event extraction. It reads `events-extraction.md`, groups events by triggering command, maps them to any explicit event infrastructure, and detects implicit or silent side effects that aren't modeled as events.

## Format

The file has a `# Events` heading (NOT `# Domain Events`), then `## AggregateName` sections (one per aggregate), then `## Summary`.

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
**Subscribers:** None — logEvent is fire-and-forget to EventLog; no other system reacts to this event
**Cascade effects:** None

**Status:** ✅ Explicit

---

### TimeTrackingStarted

The timeTrackingStart timestamp has been set on a task transitioning to InProgress status.

| Trigger | Source | Location | Detail |
|---------|--------|----------|--------|
| UpdateTask | Frontend | TasksView status picker | Sets timeTrackingStart to now when status transitions to InProgress |
| UpdateTask | Frontend | TaskItem status picker | Same transition logic as TasksView |

**Triggered by:** UpdateTask (status transition to InProgress)
**Event type:** ⚠️ Implicit — timestamp is set as a side effect with no named event
**Subscribers:** None — no system reacts to this timestamp change
**Cascade effects:** None

**Status:** ⚠️ Implicit

---

### ChildTasksCascadeDeleted

Child task rows are removed by the database when a parent task row is hard-deleted.

| Trigger | Source | Location | Detail |
|---------|--------|----------|--------|
| HardDeleteTask | Database | parent_task_id FK CASCADE | ON DELETE CASCADE removes child rows |

**Triggered by:** HardDeleteTask (or any direct DELETE on task table)
**Event type:** ⚠️ Silent — database cascade with no application-level awareness
**Subscribers:** None — application has no hook into this cascade
**Cascade effects:** <ul><li>Child task_stakeholders and task_tags join rows also cascade-deleted</li><li>Grandchild tasks also cascade-deleted recursively</li></ul>

**Status:** ⚠️ Silent

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
| StakeholderOrphaned | Stakeholder | DeleteTask | Stakeholder may have no remaining task associations after task deletion |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

## Event Section Rules

Each event section within an aggregate must have EXACTLY these parts in EXACTLY this order. Do not add, rename, or reorganize fields:

1. `### EventName` — PascalCase past-tense heading
2. One sentence description — what happened, in domain terms
3. Evidence table with columns `Trigger | Source | Location | Detail`
4. `**Triggered by:**` — which command(s) from commands.md cause this event
5. `**Event type:**` — one of the four status values with brief explanation
6. `**Subscribers:**` — what reacts to this event, or `None` if nothing
7. `**Cascade effects:**` — what further changes result, using `<ul><li>` for multiple. `None` if none.
8. `**Status:**` — exactly one of the four status values
9. `---` separator

These are the ONLY permitted fields. Do not add or rename fields.

## Forbidden Field Names

Do NOT use any of these field names — they are not part of the format:

- `**Intent:**`
- `**Inputs:**`
- `**Preconditions:**`
- `**Enforcement layers:**`
- `**Gaps:**`
- `**Behavior mapping:**`
- `**Validation:**`

## Forbidden Table Formats

Do NOT use any of these table column layouts:

- `Source | Location | Parameters` — that's the commands format
- `Field | Required | Notes` — that's an API spec format
- `Field | Evidence` — that's a key-value format

The ONLY table in each event section is `Trigger | Source | Location | Detail`.

## Status Definitions

| Status | Meaning |
|--------|---------|
| ✅ Explicit | Event is named and published via an event mechanism (logEvent, event bus, message queue) |
| ⚠️ Implicit | Side effect happens as part of a command but is not modeled as a named event |
| ⚠️ Silent | Cascade happens at the database level with no application-layer awareness |
| ❌ Undocumented | Side effect found in code but not mentioned in any documentation or entity model |

Every event MUST have a `**Status:**` line with one of these four values. No other status values are permitted.

## Summary Section

The file must end with `## Summary` containing:

1. Status counts table — `Status | Count`
2. `### Implicit Events` — events that happen but aren't named
3. `### Silent Cascades` — DB-level cascades the application doesn't know about
4. `### Undocumented Events` — side effects not mentioned in documentation

If none exist in a category, state "None found."

Do NOT add a `## Cross-Cutting Observations` or `## Gaps` section.

## Structural Rules

- Aggregate sections sorted alphabetically
- Events within aggregates sorted alphabetically
- Summary always last
- Event names use past tense — `TaskCreated` not `CreateTask`, `TimeTrackingStarted` not `StartTimeTracking`
- Trigger names match command names from commands.md
- No code snippets or file paths

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`, `stores/`)
- Code snippets or backtick-wrapped code
- `# Domain Events` heading (must be `# Events`)
- `## Summary Table` before the aggregate sections
- `## Cross-Cutting Observations` or `## Gaps` sections
- Per-event `Field | Required | Notes` input tables
- Raw extraction data without analysis
