# Event Extraction Template

## File Location

`docs/flow/ddd/events-extraction.md`

This file is Phase 1 of domain event extraction. It captures raw evidence of every side effect, cascade, reaction, and published event found in the codebase. The events analyzer reads this to produce the final events.md.

## Format

The file has a `# Event Extraction` heading and a single flat table. Nothing else. No prose, no analysis, no sections, no grouping, no summary tables, no per-event subsections, no observation sections.

```markdown
# Event Extraction

| Event | Trigger | Aggregate | Source | Location | Detail |
|-------|---------|-----------|--------|----------|--------|
| TaskCreated | CreateTask | Task | Frontend | tasks store createTask | logEvent fired with taskId, name, scheduledTimeBlockId, projectId after successful POST |
| TaskCreated | CreateTask | Task | API | POST handler | No explicit event — only the frontend publishes this event name |
| StakeholdersReplaced | CreateTask | Task | API | POST handler | task_stakeholders join rows deleted and reinserted in same transaction |
| TimeTrackingStarted | UpdateTask | Task | Frontend | TasksView status picker | Sets timeTrackingStart when status transitions to InProgress |
| ChildTasksCascadeDeleted | HardDeleteTask | Task | Database | parent_task_id FK CASCADE | ON DELETE CASCADE removes child task rows when parent row is deleted |
| ProjectReferencesNulled | SoftDeleteProject | Project | Database | task project_id FK SET NULL | ON DELETE SET NULL clears task.project_id when project is deleted |
| OldEventsPurged | PurgeOldEvents | EventLog | Backend | Scheduler tick handler | Deletes event_log records older than 10 days every hour |
```

THIS IS THE ONLY ACCEPTABLE FORMAT.

The output is ONE table under ONE heading. Not multiple tables. Not sections per event. Not sections per aggregate. Not sections per command.

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Event | PascalCase past-tense name describing what happened. `TaskCreated` not `CreateTask`. `TimeTrackingStarted` not `SetTimeTracking`. | Never empty. |
| Trigger | PascalCase command name from commands.md that causes this event. | Never empty. |
| Aggregate | PascalCase entity name from classification.md that is affected. | Never empty. |
| Source | Where the side effect originates: `API`, `Frontend`, `Backend`, `Database`, `Scheduler`, `Migration`. | Never empty. |
| Location | Descriptive location — handler name, component name, FK constraint name. NOT a file path. | Never empty. |
| Detail | How the side effect manifests — what changes, what is published, what cascades. One or two sentences. | Never empty. |

## What to Look For

- Explicit event publications (logEvent calls, event bus emissions, message queue publishes)
- Database cascades triggered by FK constraints (ON DELETE CASCADE, ON DELETE SET NULL, ON UPDATE CASCADE)
- Transactional side effects (join table replacements, computed field updates, timestamp mutations)
- Frontend reactions to state changes (UI updates, store mutations, navigation changes)
- Backend reactions (scheduled cleanup, background processing, cache invalidation)
- Implicit state transitions (setting timeTrackingStart when status changes to InProgress)

## What to Exclude

- The command itself — only its consequences
- Pure reads that don't change state
- Logging to stdout/stderr (unless it's a domain event log like EventLog)
- Infrastructure events (connection pooling, health checks)

## Structural Rules

- One row per event per source — same event from Database and Frontend gets two rows
- Sort by Aggregate alphabetically, then by Event within aggregate, then by Source
- Use past-tense domain language for Event names
- Trigger must match a Command name from commands-extraction.md or commands.md
- No file paths — use descriptive locations
- Detail column is one or two sentences, not multi-paragraph

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`, `stores/`, `repository/`)
- Per-event sections with `### EventName` headings — the output is a flat table
- Multiple tables — there is exactly one table
- Summary tables, cascade tables, or observation sections
- `Field | Evidence` key-value tables — use the specified column format
- Analysis or status judgments — that's Phase 2
