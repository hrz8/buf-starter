# Task Breakdown Workflow

**OUTPUT: Task documentation files ONLY. Do NOT implement code.**

## Critical Principles

### 1. Think Harder - Find Optimal Breakdown

**DON'T** rush to create tasks. **DO** understand:
- Full user story scope
- Existing codebase architecture
- Natural domain/technical boundaries
- What changes group together
- Sequential vs. parallel work

### 2. Efficiency Over Granularity

- **ONE task can handle MULTIPLE things** in same domain/context
- **Group related changes** together
- **NO micro-tasks** (separate tasks for each file/function)
- **NO over-engineering** into 20+ tiny tasks

### 3. Strict Scope Adherence

- Tasks ONLY implement what's explicitly in user story
- NO additions, improvements, or "nice-to-haves"
- NO refactoring unless explicitly mentioned
- NO extra features or clever additions

## Workflow

### Phase 1: Enter Plan Mode

**Required:** Use `EnterPlanMode` tool to begin.

### Phase 2: Deep Context Understanding

**Read user story completely** from `docs/stories/{filename}.md`
- Identify all acceptance criteria
- Understand technical requirements
- Note dependencies
- Check out-of-scope items

**Review codebase context:**
- Dev guidelines
- Similar implemented features
- Existing architecture patterns

### Phase 3: Plan the Breakdown

**Optimal task count:** 3-8 tasks typically

**Common task groups:**

1. **Database/Schema** (if significant DB changes)
   - Migrations + schema + partition setup

2. **Backend Implementation**
   - Proto definitions + domain layer (model/interface/repo/service/handler)

3. **Frontend Implementation**
   - Repository + service composable + components

4. **Integration**
   - Navigation + i18n + final integration

### Phase 4: Create Task Files

**Naming:** `T{number}-{kebab-case-description}.md`

Check existing task numbers:
```bash
ls docs/tasks/T*.md | sort -V | tail -1
```

**Task Template:**

```markdown
# Task T{number}: {Title}

**Story Reference:** US{X}-{story-name}.md
**Type:** {Backend|Frontend|Database|Integration|Fullstack}
**Priority:** {High|Medium|Low}
**Estimated Effort:** {X-Y hours}
**Prerequisites:** {T{n}-task-name} (if dependencies)

## Objective

{1-2 sentence description}

## Acceptance Criteria

- [ ] {Testable criterion from user story}

## Technical Requirements

{Detailed requirements with subsections}

## Implementation Details

{Guidance following established patterns}

## Files to Create

- `path/to/file.ts` - {Description}

## Files to Modify

- `path/to/existing.ts` - {What to modify}

## Testing Requirements

- {Testing approach}

## Commands to Run

```bash
{commands}
```

## Definition of Done

- [ ] {Completion criterion}
- [ ] Follows established patterns
- [ ] Tests passing
```

### Phase 5: Quality Verification

- [ ] Complete coverage of user story requirements
- [ ] No duplication or overlap
- [ ] No scope additions
- [ ] Clear dependencies
- [ ] Efficient grouping
- [ ] Implementable independently

### Phase 6: Exit Plan Mode

Use `ExitPlanMode` tool when tasks are complete.

## Anti-Patterns

### "File Per Task" - WRONG

```
T7-create-model-file.md
T8-create-repo-file.md
T9-create-service-file.md
```

**RIGHT:**
```
T7-entity-backend-domain-layer.md
  - model, repo, service, handler together
```

### "Operation Per Task" - WRONG

```
T7-implement-create-operation.md
T8-implement-read-operation.md
T9-implement-update-operation.md
```

**RIGHT:**
```
T7-entity-crud-operations.md
  - Create, Read, Update, Delete together
```

### "Feature Creep" - WRONG

```
T7-user-login-with-extras.md
"Implement login... also add password strength, social login, 2FA"
```
(But story only mentions basic login!)

**RIGHT:**
```
T7-user-basic-login.md
"Implement login as specified in user story"
```

## Good Example

**User Story:** Project Settings (view, update, delete)

```
T7-project-settings-backend-api.md
  - Proto definitions (Get, Update, Delete)
  - Domain layer (model, interface, repo, service, handler)
  - Validation and error handling

T8-project-settings-frontend-ui.md
  - Repository methods
  - Service composable
  - Settings page and form
  - Delete dialog

T9-project-settings-integration.md
  - Navigation config
  - Breadcrumbs
  - i18n translations
  - Final testing
```

**Why good:** 3 tasks, clear boundaries, efficient grouping.
