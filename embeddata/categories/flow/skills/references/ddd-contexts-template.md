# Bounded Contexts Template

## File Location

`docs/flow/ddd/contexts.md`

This file is Phase 2 of bounded context extraction. It reads `contexts-extraction.md`, proposes bounded contexts by grouping related aggregates, maps coupling between contexts, and assesses boundary quality.

## Format

The file has a `# Bounded Contexts` heading (NOT `# Contexts` or `# Domain Contexts`), then `## ContextName` sections (one per proposed context), then `## Context Map` with coupling between contexts, then `## Summary`.

```markdown
# Bounded Contexts

## Task Management

The core scheduling and execution context. Owns the task lifecycle from creation through completion, including time tracking, scheduling, and hierarchical task structure.

**Aggregates:** Task, TimeBlock

**Rationale:** Task and TimeBlock are tightly coupled through scheduling fields (ScheduledDate, ScheduledTimeBlockId) and share drag-and-drop interactions. TimeBlock capacity calculations depend on task assignments. Commands frequently update both in the same user workflow.

| Signal | Type | Detail |
|--------|------|--------|
| Task references TimeBlock via FK | Data coupling | task.scheduled_time_block_id FK; ON DELETE SET NULL |
| BatchReorderTasks moves tasks across blocks | Command coupling | Drag-and-drop handler updates both order and scheduledTimeBlockId |
| TimeBlock rebuild clears overdue task dates | Event coupling | ClearOverdueScheduledDate modifies tasks during TimeBlock component lifecycle |

**Internal boundary quality:** ✅ Clean — Task and TimeBlock interact through well-defined scheduling fields

---

## Project Organization

Groups tasks into named projects with soft-delete lifecycle and hierarchical structure.

**Aggregates:** Project

**Rationale:** Project is referenced by Task but owns no task behavior. Deleting a project nullifies task references rather than cascading. Project has its own lifecycle (soft-delete, restore, system project protection) independent of task state.

| Signal | Type | Detail |
|--------|------|--------|
| Task references Project via FK | Data coupling | task.project_id FK; ON DELETE SET NULL; Task can exist without Project |
| SoftDeleteProject refetches task list | Event coupling | Frontend refetches tasks after project deletion to reflect nullified references |

**Internal boundary quality:** ✅ Clean — single aggregate context

---

## Context Map

### Task Management ↔ Project Organization

| Direction | Signal | Type | Detail |
|-----------|--------|------|--------|
| Task Management → Project Organization | Task references Project via FK | Data coupling | task.project_id FK; ON DELETE SET NULL |
| Task Management → Project Organization | BatchReorderTasks sets projectId | Command coupling | Drag across project groups changes task.project_id |
| Project Organization → Task Management | SoftDeleteProject refetches tasks | Event coupling | Frontend compensating read after delete |

**Boundary status:** ⚠️ Leaky boundary — Task writes project_id directly; no anti-corruption layer

---

## Summary

### Proposed Contexts

| Context | Aggregates | Internal Quality |
|---------|-----------|-----------------|
| Task Management | Task, TimeBlock | ✅ Clean |
| Project Organization | Project | ✅ Clean |

### Context Coupling

| From | To | Coupling Count | Boundary Status |
|------|----|---------------|-----------------|
| Task Management | Project Organization | 3 | ⚠️ Leaky boundary |

### Boundary Issues

| Issue | Contexts | Detail |
|-------|----------|--------|
| Direct FK write across boundary | Task Management → Project Organization | Task writes project_id directly rather than through a Project command |
```

THIS IS THE ONLY ACCEPTABLE FORMAT. Do not use any other format.

## Context Section Rules

Each context section must have EXACTLY these parts in EXACTLY this order:

1. `## ContextName` — human-readable name describing the subdomain
2. One or two sentence description of what this context owns
3. `**Aggregates:**` — comma-separated list of aggregates in this context
4. `**Rationale:**` — why these aggregates belong together (coupling evidence)
5. Evidence table with columns `Signal | Type | Detail` — key coupling signals within this context
6. `**Internal boundary quality:**` — one of: `✅ Clean`, `⚠️ Leaky boundary`, `❌ Circular dependency`
7. `---` separator

These are the ONLY permitted fields. Do not add or rename fields.

## Context Map Section Rules

The `## Context Map` section contains one subsection per context pair that has coupling:

1. `### ContextA ↔ ContextB` — bidirectional heading
2. Evidence table with columns `Direction | Signal | Type | Detail`
3. `**Boundary status:**` — one of the status values with brief explanation
4. `---` separator

## Forbidden Field Names

Do NOT use any of these field names:

- `**Intent:**`, `**Inputs:**`, `**Preconditions:**`
- `**Behavior mapping:**`, `**Validation:**`, `**Side effects:**`
- `**Triggered by:**`, `**Event type:**`, `**Subscribers:**`, `**Cascade effects:**`

## Status Definitions

| Status | Meaning |
|--------|---------|
| ✅ Clean boundary | Contexts interact only through well-defined, direction-aware interfaces (FK references, published events) |
| ⚠️ Leaky boundary | Contexts share transactions, write each other's fields directly, or have implicit knowledge of each other's internals |
| ⚠️ Missing boundary | Aggregates that logically belong to different contexts are tightly coupled with no separation |
| ❌ Circular dependency | Bidirectional command or data coupling where each context writes to the other |

## Summary Section

The file must end with `## Summary` containing:

1. `### Proposed Contexts` — table of contexts with their aggregates and internal quality
2. `### Context Coupling` — table of context pairs with coupling count and boundary status
3. `### Boundary Issues` — table of specific issues with affected contexts and detail

If no issues exist, state "No boundary issues found."

Do NOT add a `## Recommendations` or `## Refactoring Plan` section. The file documents what IS, not what SHOULD BE.

## Structural Rules

- Context sections sorted alphabetically
- Context Map subsections sorted alphabetically by the first context name
- Summary always last
- Aggregate names must match classification.md
- No code snippets or file paths
- Context names are human-readable subdomain names, not code module names

## Forbidden Content

- File paths (`.go`, `.ts`, `.vue`, `src/`)
- Code snippets or backtick-wrapped code
- `# Contexts` or `# Domain Contexts` heading (must be `# Bounded Contexts`)
- Refactoring recommendations or migration plans
- Per-aggregate sections (group by context, not by aggregate)
- Raw extraction data without analysis
