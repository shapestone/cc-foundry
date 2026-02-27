# Domain Services Template

## File Location

`docs/flow/ddd/services.md`

This file is Phase 2 of domain service extraction. It reads `services-extraction.md`, groups services by owning bounded context, maps which aggregates each service coordinates, and assesses encapsulation quality.

## Format

The file has a `# Services` heading (NOT `# Domain Services`), then `## ContextName` sections (one per bounded context that has services), then `## Cross-Context Services` for services that span context boundaries, then `## Summary`.

```markdown
# Services

## Task Management

### RebuildTimeBlocks

Orchestrates the full TimeBlock display lifecycle: clears overdue scheduled dates from tasks, computes block capacity from assigned task counts, and distributes or surfaces overdue tasks based on user settings.

| Aggregates | Source | Location |
|-----------|--------|----------|
| Task, TimeBlock | Frontend | TimeBlocks component rebuild |

**Type:** Orchestration
**Owning context:** Task Management
**Trigger:** Page load, task data change, settings change
**Steps:** <ul><li>Identify tasks with past scheduledDate</li><li>If autoReschedule enabled: clear scheduledDate via updateTask for each overdue task</li><li>If autoReschedule disabled: surface overdue tasks in decision surface</li><li>Compute blockCapacity for each block from active task count vs MaxLoad ceiling</li><li>Log auto-promotion or decision-surface event</li></ul>
**Encapsulation:** ⚠️ Scattered — orchestration lives in a Vue component rather than a dedicated service; mixes UI concerns with domain logic

**Status:** ⚠️ Scattered

---

### SyncStakeholders

Coordinates transactional consistency between Task and Stakeholder by replacing all stakeholder associations on every Task write.

| Aggregates | Source | Location |
|-----------|--------|----------|
| Task, Stakeholder | API | POST/PUT handler transaction |

**Type:** Coordination
**Owning context:** Task Management
**Trigger:** CreateTask, UpdateTask
**Steps:** <ul><li>Delete all task_stakeholders rows for the task</li><li>Insert new task_stakeholders rows for each stakeholder ID in the request</li><li>Both operations within the same database transaction as the Task INSERT/UPDATE</li></ul>
**Encapsulation:** ✅ Encapsulated — runs within the Task handler's transaction boundary

**Status:** ✅ Encapsulated

---

## Cross-Context Services

### FilterSystemProjects

Applies Project context business rules within Task Management UI by filtering out system-flagged projects from assignment dropdowns.

| Aggregates | Source | Location |
|-----------|--------|----------|
| Project, Task | Frontend | TaskCreateModal project dropdown |

**Type:** Policy
**Owning context:** Crosses Task Management → Work Organization
**Trigger:** TaskCreateModal render
**Steps:** <ul><li>Read all projects from Project store</li><li>Filter where isSystem is false</li><li>Display filtered list in Task creation dropdown</li></ul>
**Encapsulation:** ❌ Boundary violation — Task Management UI embeds Work Organization's system-flag business rule directly

**Status:** ❌ Boundary violation

---

## Summary

### Service Inventory

| Service | Type | Context | Aggregates | Status |
|---------|------|---------|-----------|--------|
| RebuildTimeBlocks | Orchestration | Task Management | Task, TimeBlock | ⚠️ Scattered |
| SyncStakeholders | Coordination | Task Management | Task, Stakeholder | ✅ Encapsulated |
| FilterSystemProjects | Policy | Cross-context | Project, Task | ❌ Boundary violation |

### Encapsulation Issues

| Service | Issue |
|---------|-------|
| RebuildTimeBlocks | Orchestration logic lives in a Vue component instead of a dedicated service |
| FilterSystemProjects | Task Management UI embeds Work Organization's system-flag rule |

### Service Type Distribution

| Type | Count |
|------|-------|
| Orchestration | 2 |
| Coordination | 3 |
| Computation | 4 |
| Policy | 1 |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

## Service Section Rules

Each service section must have EXACTLY these parts in EXACTLY this order:

1. `### ServiceName` — PascalCase verb+noun heading
2. One or two sentence description of what the service does
3. Evidence table with columns `Aggregates | Source | Location`
4. `**Type:**` — one of: `Orchestration`, `Coordination`, `Computation`, `Policy`
5. `**Owning context:**` — bounded context name from contexts.md, or `Crosses ContextA → ContextB`
6. `**Trigger:**` — what causes this service to execute
7. `**Steps:**` — the algorithm or coordination pattern, using `<ul><li>` for multiple steps
8. `**Encapsulation:**` — one of the status values with brief explanation
9. `**Status:**` — exactly one of the status values
10. `---` separator

These are the ONLY permitted fields. Do not add or rename fields.

## Forbidden Field Names

Do NOT use any of these field names:

- `**Intent:**`, `**Inputs:**`, `**Preconditions:**`
- `**Behavior mapping:**`, `**Validation:**`, `**Side effects:**`
- `**Triggered by:**`, `**Event type:**`, `**Subscribers:**`, `**Cascade effects:**`
- `**Aggregates:**`, `**Rationale:**`, `**Internal boundary quality:**`, `**Boundary status:**`

## Forbidden Table Formats

Do NOT use any of these table column layouts:

- `Source | Location | Parameters` — that's commands format
- `Trigger | Source | Location | Detail` — that's events format
- `Signal | Type | Detail` — that's contexts format
- `Field | Required | Notes` — that's API spec format

The ONLY table in each service section is `Aggregates | Source | Location`.

## Status Definitions

| Status | Meaning |
|--------|---------|
| ✅ Encapsulated | Service logic is contained within its owning context and uses well-defined interfaces to interact with other aggregates |
| ⚠️ Scattered | Service logic is spread across multiple locations (e.g., split between a component and a store, or duplicated in multiple handlers) |
| ⚠️ Misplaced | Service logic lives in the wrong layer (e.g., domain logic in a UI component, orchestration in a database trigger) |
| ❌ Boundary violation | Service crosses bounded context boundaries without a proper anti-corruption layer |

## Summary Section

The file must end with `## Summary` containing:

1. `### Service Inventory` — table of all services with type, context, aggregates, and status
2. `### Encapsulation Issues` — table of services with ⚠️ or ❌ status and their issues
3. `### Service Type Distribution` — count of each service type

If no issues exist, state "No encapsulation issues found."

## Structural Rules

- Context sections sorted alphabetically
- Services within contexts sorted alphabetically
- Cross-Context Services section after all context sections, before Summary
- Summary always last
- Service names use PascalCase verb+noun
- Context names match contexts.md
- Aggregate names match classification.md
- No code snippets or file paths

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`)
- Code snippets or backtick-wrapped code
- `# Domain Services` heading (must be `# Services`)
- Refactoring recommendations
- Fields from other formats
