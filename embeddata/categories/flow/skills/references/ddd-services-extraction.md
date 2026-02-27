# Service Extraction Template

## File Location

`docs/flow/ddd/services-extraction.md`

This file is Phase 1 of domain service extraction. It captures raw evidence of stateless operations that orchestrate, coordinate, or compute across aggregate boundaries — logic that doesn't belong to a single entity or value object. The services analyzer reads this to produce the final services.md.

## Format

The file has a `# Service Extraction` heading and a single flat table. Nothing else. No prose, no analysis, no sections, no grouping, no summary tables, no per-service subsections.

```markdown
# Service Extraction

| Service | Type | Aggregates | Source | Location | Detail |
|---------|------|-----------|--------|----------|--------|
| RebuildTimeBlocks | Orchestration | Task, TimeBlock | Frontend | TimeBlocks component rebuild | Clears overdue scheduled dates, computes block capacity, distributes tasks across blocks on page load and data change |
| SyncStakeholders | Coordination | Task, Stakeholder | API | POST/PUT handler transaction | Delete-all then reinsert task_stakeholders join rows in same transaction as Task write |
| ComputeBlockCapacity | Computation | Task, TimeBlock | Frontend | TimeBlocks blockCapacity function | Counts active tasks assigned to each block against MaxLoad-scaled capacity ceiling |
| ComputeDerivedStatus | Computation | Task | Backend | task repository read path | Inspects child task statuses to compute parent DerivedStatus on every List and GetByID |
| FilterSystemProjects | Policy | Project, Task | Frontend | TaskCreateModal project dropdown | Filters out system projects from assignment dropdown, embedding Project policy in Task creation |
```

THIS IS THE ONLY ACCEPTABLE FORMAT.

The output is ONE table under ONE heading. Not multiple tables. Not sections per service. Not sections per aggregate.

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Service | PascalCase name describing the operation. Use verb+noun: `RebuildTimeBlocks`, `SyncStakeholders`, `ComputeBlockCapacity`. | Never empty. |
| Type | Category of service: `Orchestration`, `Coordination`, `Computation`, `Policy`. | Never empty. |
| Aggregates | Comma-separated PascalCase aggregate names involved. | Never empty — at least one aggregate. |
| Source | Where the service logic lives: `API`, `Frontend`, `Backend`, `Database`, `Migration`. | Never empty. |
| Location | Descriptive location — function name, component name, handler area. NOT a file path. | Never empty. |
| Detail | What the service does — the algorithm or coordination pattern. One or two sentences. | Never empty. |

## Type Definitions

- **Orchestration** — Multi-step process that coordinates work across aggregates or manages a workflow. Has sequencing, conditional logic, or state machine behavior. Examples: rebuild/reschedule algorithms, batch processing workflows, multi-aggregate lifecycle management.
- **Coordination** — Ensures transactional consistency across aggregates. Writes to multiple aggregate tables in a single transaction or manages compensating actions. Examples: sync join tables, atomic multi-aggregate updates, saga-like patterns.
- **Computation** — Derives values that span aggregate boundaries. Reads from multiple aggregates to produce a result that no single aggregate owns. Examples: capacity calculations, derived status from child aggregates, cross-aggregate statistics.
- **Policy** — Business decision logic that doesn't belong to a single entity. Applies rules, filters, or transformations based on cross-aggregate knowledge. Examples: system project filtering, scheduling conflict detection, access control decisions.

## What to Look For

### 1. Multi-Aggregate Operations
- Functions that read or write multiple aggregate types
- Transaction scopes that span multiple tables
- Batch operations that process items from different aggregates

### 2. Workflow Orchestration
- Multi-step processes with conditional branching
- Rebuild/refresh logic triggered by lifecycle events
- Scheduling or rescheduling algorithms
- State machine transitions that involve multiple aggregates

### 3. Cross-Aggregate Computations
- Derived values that require data from multiple aggregates
- Capacity or utilization calculations
- Aggregated statistics spanning multiple entity types
- Computed properties that don't belong to any single entity

### 4. Business Policy Logic
- Filtering or validation that crosses aggregate boundaries
- Access control decisions based on cross-aggregate state
- Default value computation using other aggregate data
- Business rule enforcement that spans contexts

### 5. Coordination Patterns
- Join table management (sync patterns)
- Compensating reads after cross-aggregate mutations
- Cache invalidation across aggregates
- Event-driven coordination between stores

## What to Exclude

- Simple CRUD operations on a single aggregate — those are commands
- Side effects and cascades — those are events
- Validation within a single aggregate — those are rules
- Pure read queries with no computation — those are queries

## Structural Rules

- One row per service — if the same logic appears in multiple locations, pick the primary location and note alternatives in Detail
- Sort by Service name alphabetically
- Use domain language for Service names — PascalCase verb+noun
- Aggregates column lists all aggregates touched, in PascalCase
- No file paths — use descriptive locations
- Detail column is one or two sentences

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`, `stores/`)
- Per-service sections with `### ServiceName` headings — the output is a flat table
- Multiple tables — there is exactly one table
- Summary, observation, or recommendation sections
- Analysis or status judgments — that's Phase 2
