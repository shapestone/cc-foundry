# Context Extraction Template

## File Location

`docs/flow/ddd/contexts-extraction.md`

This file is Phase 1 of bounded context extraction. It captures raw evidence of coupling between aggregates — data dependencies, command boundaries, event flows, shared concepts, and module boundaries. The context analyzer reads this to propose bounded contexts.

## Format

The file has a `# Context Extraction` heading and a single flat table. Nothing else. No prose, no analysis, no sections, no grouping, no summary tables, no per-context subsections.

```markdown
# Context Extraction

| Signal | Type | From | To | Source | Location | Detail |
|--------|------|------|----|--------|----------|--------|
| Task references Project via FK | Data coupling | Task | Project | Database | task.project_id FK | ON DELETE SET NULL; Task can exist without Project |
| CreateTask syncs stakeholders in same transaction | Command coupling | Task | Stakeholder | API | POST handler | task_stakeholders join written in same transaction as task INSERT |
| SoftDeleteProject refetches task list | Event coupling | Project | Task | Frontend | projects store deleteProject | Task list refetched after project soft-delete to reflect nullified references |
| Task and TimeBlock share scheduling concept | Shared concept | Task | TimeBlock | Entity model | entities.md and aggregates.md | Both use ScheduledDate, ScheduledTimeBlockId; TimeBlock owns time ranges, Task references them |
| Container API exists with no UI entry points | Dead module | Container | — | Frontend | All views scanned | Store actions and API routes exist but no view component triggers them |
| EventLog written by all command handlers | Hub dependency | Task | EventLog | API | All POST/PUT/DELETE handlers | Every mutation fires logEvent to EventLog; EventLog has no awareness of callers |
```

THIS IS THE ONLY ACCEPTABLE FORMAT.

The output is ONE table under ONE heading. Not multiple tables. Not sections per context. Not sections per aggregate pair.

## Column Rules

| Column | Content | Empty Value |
|--------|---------|-------------|
| Signal | Short description of the coupling or boundary evidence found. Plain language. | Never empty. |
| Type | Category of coupling: `Data coupling`, `Command coupling`, `Event coupling`, `Shared concept`, `Shared table`, `Dead module`, `Hub dependency`, `Circular dependency`. | Never empty. |
| From | PascalCase aggregate name that initiates or owns the coupling. | Never empty. |
| To | PascalCase aggregate name that is coupled to. | `—` if the signal is about a single aggregate (e.g., Dead module). |
| Source | Where the evidence was found: `Database`, `API`, `Frontend`, `Backend`, `Entity model`, `Configuration`, `Migration`. | Never empty. |
| Location | Descriptive location — FK name, handler name, store name, document name. NOT a file path. | Never empty. |
| Detail | How the coupling manifests. One or two sentences. | Never empty. |

## What to Look For

### 1. Data Coupling
- Foreign key relationships between aggregate tables
- Shared lookup tables or join tables
- Cross-aggregate queries (JOINs, subqueries)
- Shared database sequences or ID generators

### 2. Command Coupling
- Commands that write to multiple aggregate tables in one transaction
- Commands that read from one aggregate to validate another
- Commands that must be called in a specific order across aggregates
- Handlers that call other handlers or services for different aggregates

### 3. Event Coupling
- Events from one aggregate that trigger reactions in another
- Frontend store actions that refetch or update other aggregate stores after a command
- Backend jobs that process events from one aggregate to update another
- Cascading side effects that cross aggregate boundaries

### 4. Shared Concepts
- Value objects or types used by multiple aggregates
- Business rules that span aggregate boundaries
- Domain terms that mean different things in different contexts
- Shared enums, constants, or configuration

### 5. Module Boundaries
- Aggregates that share a handler file or service class
- Aggregates with separate handler files that never interact
- Frontend stores that import from each other vs. those that are independent
- API route groupings (shared prefix vs. separate prefixes)

### 6. Dead or Orphaned Modules
- Aggregates with backend code but no frontend entry points
- Aggregates referenced in documentation but not in code
- API routes with no callers

## What to Exclude

- Coupling that exists only in documentation (e.g., entities.md mentions both but code doesn't connect them)
- Infrastructure coupling (shared database connection, shared HTTP server) — this is universal in a monolith
- Read-only coupling where one aggregate queries another without changing it (note: this IS coupling if it creates a runtime dependency, but exclude simple display joins)

## Structural Rules

- One row per coupling signal — the same pair of aggregates can appear multiple times with different signal types
- Sort by From aggregate alphabetically, then by To aggregate, then by Type
- Use domain language for Signal descriptions
- No file paths — use descriptive locations
- Detail column is one or two sentences, not multi-paragraph
- Include both directions if coupling is bidirectional (From=A, To=B AND From=B, To=A)

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`, `handlers/`, `stores/`, `repository/`)
- Per-signal sections with `### SignalName` headings — the output is a flat table
- Multiple tables — there is exactly one table
- Summary tables, coupling matrices, or observation sections
- Analysis or context proposals — that's Phase 2
