# User Story Breakdown Command

**Variables:** `user_story_path` = $ARGUMENTS

## Purpose

Transform comprehensive user stories from `docs/stories/*.md` into focused, efficient task documents in `docs/tasks/*.md`. Tasks should balance scope and efficiency - not too granular, not too broad. Each task should be implementable independently or as part of a clear dependency chain, allowing executors (developers or LLMs) to focus on implementation without missing context from the user story.

## Critical Principles

### 1. Think Harder - Find the Optimal Breakdown

**DON'T** rush to create tasks. **DO** spend time understanding:

- The full user story scope and requirements
- The existing codebase architecture and patterns
- The natural boundaries in the domain/technical layers
- What changes naturally group together
- Which changes must be sequential vs. parallel

### 2. Efficiency Over Granularity

- ✅ **ONE task can handle MULTIPLE things** if they're in the same domain/context
- ✅ **Group related changes** together (e.g., "API endpoint + service layer + validation")
- ❌ **NO micro-tasks** (e.g., separate tasks for "create file", "add function", "add test")
- ❌ **NO over-engineering** the breakdown into 20+ tiny tasks

### 3. Strict Adherence to User Story

- ✅ Tasks **ONLY** implement what's **explicitly** in the user story
- ❌ **NO** additions, improvements, or "nice-to-haves"
- ❌ **NO** refactoring unless explicitly mentioned in user story
- ❌ **NO** extra features, optimizations, or clever additions

### 4. No Duplication or Overlap

- Each task has a **clear, unique scope**
- No tasks with overlapping responsibilities
- No redundant work across tasks

## Workflow

### Phase 1: Enter Plan Mode (REQUIRED)

**CRITICAL**: You **MUST** use the `EnterPlanMode` tool to begin. This ensures you:

- Have time and space to deeply analyze the user story
- Can explore the codebase thoroughly
- Can reference dev guidelines and patterns
- Can plan the optimal breakdown before creating tasks

```
Use EnterPlanMode tool → Triggers plan mode for breakdown analysis
```

### Phase 2: Deep Context Understanding (REQUIRED)

**Before creating ANY tasks**, you MUST understand:

#### A. User Story Analysis

1. **Read the user story completely** from `docs/stories/$ARGUMENTS.md`
2. **Identify all acceptance criteria** - each one may map to tasks
3. **Understand technical requirements** - backend, frontend, database, API
4. **Note dependencies** - what must be built before other parts
5. **Check out-of-scope items** - understand boundaries clearly

#### B. Codebase Context

1. **Review existing domain patterns**:
   - Read `docs/dev_guidelines/DOMAIN_ARCHITECTURE_GUIDE.md`
   - Read `docs/dev_guidelines/BACKEND_GUIDE.md`
   - Read `docs/dev_guidelines/FRONTEND_GUIDE.md`
   - Read `docs/dev_guidelines/EFFICIENCY_GUIDE.md`

2. **Explore related code**:
   - Find similar features already implemented
   - Identify reusable patterns and components
   - Understand the existing architecture
   - Map out which files will be modified vs. created

3. **Understand the project structure**:
   - Backend: Domain layers (model, interface, repo, service, handler)
   - Frontend: Repository pattern, service composables, components
   - Protocol buffers: Proto definitions and code generation
   - Database: Migrations, partitioning, cascade deletes

#### C. Find Natural Boundaries

Think harder about where to split tasks:

- **Domain boundaries**: Employee domain vs. Project domain
- **Layer boundaries**: Database → Backend → Frontend → Integration
- **Functional boundaries**: CRUD operations, validation, UI, etc.
- **Dependency boundaries**: What must be done before what?

### Phase 3: Clarification (If Needed)

**If anything is unclear or ambiguous:**

1. **Use `AskUserQuestion` tool** to clarify requirements
2. **Document all clarifications** received
3. **If clarification reveals NEW requirements not in user story:**
   - **STOP** creating tasks
   - **UPDATE the user story file FIRST** at `docs/stories/$ARGUMENTS.md`
   - Add new requirements to appropriate sections
   - **THEN** proceed with task breakdown including new requirements

**Remember**: Tasks must reflect the user story. If user story is incomplete, fix it first.

### Phase 4: Plan the Breakdown

**Now, think harder:**

#### Determine Optimal Task Count

- **Too few tasks** (1-2): Likely missing clear separation of concerns
- **Too many tasks** (15+): Likely over-granular, inefficient
- **Optimal range** (3-8 tasks): Usually the sweet spot for most user stories
- **Context matters**: A simple CRUD feature might be 3-4 tasks; a complex integration might be 6-8

#### Identify Task Groups

Common patterns (adapt to your user story):

1. **Database/Schema Changes**
   - One task for: migrations + schema updates + partition setup
   - Only if significant database changes needed

2. **Backend Implementation**
   - One task for: proto definitions + domain layer (model/interface/repo/service/handler)
   - OR split if backend is complex: "Core backend" + "Additional operations"

3. **Frontend Implementation**
   - One task for: repository + service composable + main components
   - OR split by major UI sections: "Form UI" + "Table/List UI"

4. **Integration & Polish**
   - One task for: navigation + i18n + final integration
   - OR split if complex: "Navigation & routes" + "Translations & polish"

**Think about dependencies:**
- Database → Backend → Frontend → Integration (typical flow)
- Some tasks can run in parallel (backend + frontend UI mockups)

### Phase 5: Create Task Files

For each identified task, create a file in `docs/tasks/` using:

**Naming Convention**: `T{number}-{kebab-case-description}.md`

The task number sequence should continue from the highest existing task number in `docs/tasks/`.

Example for breaking down `docs/stories/US2-project-settings.md`:
```
docs/tasks/T7-project-settings-backend-api.md
docs/tasks/T8-project-settings-frontend-ui.md
docs/tasks/T9-project-settings-integration.md
```

**IMPORTANT**: Always check existing task numbers first:
```bash
ls docs/tasks/T*.md | sort -V | tail -1  # Get the last task number
```

**Task File Template**:

```markdown
# Task T{number}: {Descriptive Title}

**Story Reference:** US{X}-{story-name}.md
**Type:** {Backend Foundation|Backend Implementation|Backend Integration|Frontend Foundation|Frontend UI|Frontend Integration|Database|Integration|Fullstack}
**Priority:** {High|Medium|Low}
**Estimated Effort:** {X-Y hours}
**Prerequisites:** {T{n}-task-name, T{m}-task-name} (only if there are dependencies)

## Objective

{Clear 1-2 sentence description of what this task accomplishes}

## Acceptance Criteria

- [ ] {Specific, testable criterion from user story}
- [ ] {Specific, testable criterion from user story}
- [ ] {Specific, testable criterion from user story}

## Technical Requirements

{Detailed technical requirements with subsections as needed}

### {Subsection Title 1}

{Details about specific technical requirements}

### {Subsection Title 2}

{More technical details, can include code examples}

## Implementation Details

{Detailed implementation guidance following established patterns}

### {Implementation Aspect 1}

{Specific guidance, patterns to follow, code examples}

### {Implementation Aspect 2}

{More implementation details}

## Files to Create

- `path/to/file1.ts` - {Brief description}
- `path/to/file2.go` - {Brief description}

## Files to Modify

- `path/to/existing/file.ts` - {What to modify}
- None (if no files need modification)

## Testing Requirements

{Testing requirements and strategies}

- {Testing approach 1}
- {Testing approach 2}
- {Manual testing commands if applicable}

## Commands to Run

```bash
# {Description of what these commands do}
{command 1}
{command 2}
```

## Validation Checklist

- [ ] {Validation item 1}
- [ ] {Validation item 2}
- [ ] {Validation item 3}

## Definition of Done

- [ ] {Specific completion criterion}
- [ ] {Specific completion criterion}
- [ ] {Follows established patterns and guidelines}
- [ ] {Code quality checks pass}
- [ ] {Tests are written and passing}

## Dependencies

- {T{n}: Task name that must be completed first}
- {Existing infrastructure or systems}
- {Third-party dependencies}

## Risk Factors

- **{Low|Medium|High} Risk**: {Description of risk and why}
- **{Low|Medium|High} Risk**: {Another risk factor}

## Notes

{Any additional context, considerations, gotchas, or important information}
- {Note 1}
- {Note 2}
```

### Phase 6: Quality Verification

Before finalizing tasks, verify:

- [ ] **Complete coverage**: All user story requirements are covered by tasks
- [ ] **No duplication**: No overlapping responsibilities between tasks
- [ ] **No additions**: Tasks only implement what's in user story
- [ ] **Clear dependencies**: Task order and dependencies are logical
- [ ] **Efficient grouping**: Related work is grouped together
- [ ] **Clear boundaries**: Each task has a well-defined scope
- [ ] **Implementable**: Each task can be completed independently (after dependencies)
- [ ] **Context preserved**: Each task references user story for full context

### Phase 7: Exit Plan Mode

Once all tasks are created and verified, use `ExitPlanMode` tool to complete the breakdown.

## Examples

### Example 1: Good vs. Bad Breakdown

**User Story**: US2 - Project Settings Management (view, update, delete project settings)

#### ❌ BAD: Over-Granular (Too Many Micro-Tasks)

```
T7-create-proto-definitions.md
T8-generate-go-code.md
T9-create-model-structs.md
T10-add-interface-methods.md
T11-implement-repo-getbyid.md
T12-implement-repo-update.md
T13-implement-repo-delete.md
T14-create-service-getproject.md
T15-create-service-updateproject.md
T16-create-service-deleteproject.md
T17-create-handler-methods.md
T18-create-repository-ts.md
T19-create-useprojectservice.md
T20-create-settings-page.md
T21-create-form-component.md
T22-create-delete-dialog.md
T23-add-navigation-config.md
T24-add-i18n-translations.md
```

**Why Bad**: 18 micro-tasks! Each step is too granular. Executor loses context switching between tiny tasks. Inefficient and over-complicated.

#### ✅ GOOD: Efficient Breakdown (Logical Grouping)

```
T7-project-settings-backend-api.md
  - Proto definitions (GetProject, UpdateProject, DeleteProject)
  - Domain layer (model, interface, repo, service, handler)
  - Validation rules and error handling
  - All three operations in one cohesive backend task

T8-project-settings-frontend-ui.md
  - Repository methods (project.ts)
  - Service composable (useProjectService.ts)
  - Settings page (pages/settings/project/index.vue)
  - Form component (ProjectSettingsForm.vue)
  - Delete dialog component (ProjectSettingsDeleteDialog.vue)
  - All frontend implementation together

T9-project-settings-integration.md
  - Navigation config updates
  - Breadcrumb configuration
  - i18n translations (en-US, id-ID)
  - Final integration testing
  - Polish and refinements
```

**Why Good**: 3 tasks with clear boundaries. Each task groups related work by domain/layer. Efficient, focused, and maintains context.

### Example 2: Strict Scope Adherence

**User Story**: US3 - Employee Directory (basic list and search)

#### ❌ BAD: Adding Scope

```
T10-employee-api-with-caching-rate-limiting.md
T11-employee-ui-advanced-filtering-export.md
T12-realtime-updates-websockets.md
```

**Why Bad**: Caching, rate limiting, advanced filtering, export, and WebSockets are NOT in the user story! These are additions.

#### ✅ GOOD: Exact Scope

```
T10-employee-list-search-backend.md
  - Proto definitions for ListEmployees
  - Repository method for basic search
  - Service layer implementation
  - Handler implementation

T11-employee-directory-frontend-ui.md
  - Repository method for employee list
  - Table component with basic search
  - Integration with backend API

T12-employee-navigation-translations.md
  - Add employee directory to navigation
  - i18n keys for employee feature
```

**Why Good**: Only implements what's explicitly in the user story. No extra features.

### Example 3: Dependency Management

**User Story**: US4 - Report Generation (backend processing + frontend display)

```
T13-report-database-schema.md
  - Migration for reports table
  - Partition setup if needed
  **Dependencies**: None (can start immediately)

T14-report-generation-backend-api.md
  - Proto definitions
  - Report generation service
  - Domain layer implementation
  **Dependencies**: T13-report-database-schema (needs database schema)

T15-report-display-frontend-ui.md
  - Can work on UI mockups in parallel with backend
  - Repository and service setup
  - Report display components
  **Dependencies**: None initially, but needs T14 for API integration

T16-report-integration-testing.md
  - Connect frontend to backend
  - End-to-end testing
  - Error handling and edge cases
  **Dependencies**: T14-report-generation-backend-api, T15-report-display-frontend-ui
```

**Clear dependency chain**: Database → Backend API → Frontend display → Integration

## Anti-Patterns to Avoid

### 1. The "File Per Task" Anti-Pattern

❌ **Wrong**:
```
T7-create-model-file.md
T8-create-repo-file.md
T9-create-service-file.md
T10-create-handler-file.md
```

✅ **Right**:
```
T7-entity-backend-domain-layer.md
  - Implement model, repo, service, handler together
```

### 2. The "Operation Per Task" Anti-Pattern

❌ **Wrong**:
```
T7-implement-create-operation.md
T8-implement-read-operation.md
T9-implement-update-operation.md
T10-implement-delete-operation.md
```

✅ **Right**:
```
T7-entity-crud-operations.md
  - Implement Create, Read, Update, Delete operations together
```

**Exception**: If Create/Read are simple but Update/Delete are complex, split into:
```
T7-entity-basic-crud.md (Create/Read)
T8-entity-complex-operations.md (Update/Delete)
```

### 3. The "Layer Per Task" Anti-Pattern

❌ **Wrong**:
```
T7-all-backend-changes.md
T8-all-frontend-changes.md
```

✅ **Right**: Split by functional areas, not just layers:
```
T7-entity-backend-implementation.md (backend domain + API)
T8-entity-frontend-ui.md (repository + service + components)
T9-entity-integration.md (navigation + i18n + polish)
```

### 4. The "Feature Creep" Anti-Pattern

❌ **Wrong**:
```
T7-user-login-with-extras.md
"Implement user login... also add password strength meter, social login, and 2FA"
```
(But user story only mentions basic login!)

✅ **Right**:
```
T7-user-basic-login.md
"Implement user login with email and password as specified in user story"
```

## Command Usage

```bash
# Break down a specific user story by filename
/breakdown US2-project-settings.md

# Break down by full path
/breakdown docs/stories/US3-employee-directory.md
```

## Reminders

### Critical DO's

1. ✅ **ALWAYS** use `EnterPlanMode` tool first
2. ✅ **ALWAYS** check existing task numbers and continue sequence
3. ✅ **ALWAYS** read the user story completely before planning
4. ✅ **ALWAYS** review dev guidelines and explore codebase
5. ✅ **ALWAYS** think harder about optimal breakdown
6. ✅ **ALWAYS** group related work together for efficiency
7. ✅ **ALWAYS** clarify ambiguities with `AskUserQuestion`
8. ✅ **ALWAYS** update user story FIRST if new requirements emerge
9. ✅ **ALWAYS** maintain strict scope adherence to user story
10. ✅ **ALWAYS** mark task dependencies clearly
11. ✅ **ALWAYS** verify complete coverage and no duplication

### Critical DON'Ts

1. ❌ **NEVER** skip plan mode
2. ❌ **NEVER** create micro-tasks for every file or function
3. ❌ **NEVER** add features not in the user story
4. ❌ **NEVER** create duplicate or overlapping tasks
5. ❌ **NEVER** proceed without understanding full context
6. ❌ **NEVER** create tasks before reading dev guidelines
7. ❌ **NEVER** ignore ambiguities - always clarify
8. ❌ **NEVER** create tasks without clear boundaries
9. ❌ **NEVER** forget to specify dependencies
10. ❌ **NEVER** overcomplicate the breakdown

## Success Criteria

A successful task breakdown:

1. **Covers** all user story requirements completely
2. **Groups** related work efficiently (not too granular)
3. **Maintains** clear boundaries between tasks
4. **Has** no duplication or overlap
5. **Adheres** strictly to user story scope (no additions)
6. **Specifies** dependencies clearly
7. **Preserves** context by referencing user story
8. **Follows** established patterns from dev guidelines
9. **Can be** implemented independently (after dependencies)
10. **Is** focused and efficient (typically 3-8 tasks)

## Final Note

**Think harder.** The goal is NOT to create the most tasks. The goal is to create the OPTIMAL number of tasks that balance:

- **Efficiency**: Minimize context switching, group related work
- **Clarity**: Each task has clear scope and boundaries
- **Implementability**: Each task can be completed independently
- **Context**: Each task preserves enough context from user story

Quality > Quantity. Efficiency > Granularity. Focus > Fragmentation.
