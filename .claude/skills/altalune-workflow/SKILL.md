---
name: altalune-workflow
description: |
  Project workflow management for Altalune development including user story creation, task breakdown, and task implementation. Use when: (1) Converting ideas/feature requests into comprehensive user stories (/idea command), (2) Breaking down user stories into implementable tasks (/breakdown command), (3) Implementing specific tasks following architectural patterns (/work command). This skill consolidates the idea-to-implementation workflow for efficient development.
---

# Altalune Development Workflow

## Overview

This skill provides three interconnected workflows for turning ideas into implemented features:

1. **Idea → User Story** - Transform feature ideas into comprehensive user stories
2. **User Story → Tasks** - Break down user stories into implementable tasks
3. **Task → Implementation** - Implement tasks following architectural patterns

## Important: Output Distinction

| Command | Output | NOT the Output |
|---------|--------|----------------|
| `/idea` | User story document (`docs/stories/US*.md`) | Actual code implementation |
| `/breakdown` | Task files (`docs/tasks/T*.md`) | Actual code implementation |
| `/work` | Actual code implementation | Documentation only |

**The `/idea` and `/breakdown` commands produce DOCUMENTATION ONLY.** They do not write actual implementation code. After these commands complete, the user must explicitly invoke `/work` or manually request implementation.

## Workflow Commands

### /idea - Create User Story from Idea

**Usage:** `/idea <brief description>`

**Example:** `/idea add employee directory feature`

**Process:**
1. Explore codebase to understand relevant architecture
2. Review dev guidelines for patterns
3. Ask clarifying questions iteratively (AskUserQuestion)
4. Create comprehensive user story in `docs/stories/US{N}-{slug}.md`

**Output:** User story markdown file with acceptance criteria, technical requirements, and definition of done.

**IMPORTANT:** This command ends after creating the user story file. Do NOT proceed to implementation unless explicitly asked.

### /breakdown - Task Breakdown from User Story

**Usage:** `/breakdown <story-filename>`

**Example:** `/breakdown US2-project-settings.md`

**Process:**
1. Enter plan mode (EnterPlanMode)
2. Read user story completely
3. Review dev guidelines and codebase patterns
4. Plan optimal task breakdown (3-8 tasks typically)
5. Create task files in `docs/tasks/T{N}-{slug}.md`
6. Exit plan mode (ExitPlanMode)

**Output:** Task markdown files ready for implementation.

**Key Principles:**
- Group related work together (efficiency over granularity)
- Strict adherence to user story scope (no additions)
- No overlapping or duplicate tasks
- Clear dependencies between tasks

**IMPORTANT:** This command ends after creating task files. Do NOT proceed to implementation unless explicitly asked.

### /work - Implement Task (WRITES REAL CODE)

**Usage:** `/work <task-id>`

**Example:** `/work T7`

**Process:**
1. Read task from `docs/tasks/T{N}*.md`
2. Reference linked user story for context
3. Use `context7` MCP to lookup library docs if needed
4. Follow appropriate dev skills:
   - Backend: `altalune-backend` skill
   - Frontend: `altalune-frontend` skill
   - Chatbot modules: `altalune-chatbot` skill
   - Authorization: `altalune-authorization` skill
5. **Write actual code** - Create/modify files as specified in task
6. Run quality checks (`buf generate`, `pnpm lint:fix`, etc.)
7. **Test with MCP tools:**
   - Frontend: Use `playwright` MCP to verify UI
   - Backend: Use `postgres` MCP to verify database state

**Output:** Real code implementation - proto files, Go code, Vue components, etc.

**This is the ONLY workflow command that produces actual code changes.**

## Workflow Sequence

```
/idea "feature description"
    ↓ Creates: docs/stories/US{N}-{slug}.md
    ↓ STOPS HERE (documentation only)

/breakdown US{N}-{slug}.md
    ↓ Creates: docs/tasks/T{N}-{slug}.md (multiple files)
    ↓ STOPS HERE (documentation only)

/work T{N}
    ↓ IMPLEMENTS: Actual code changes
    ↓ Uses: altalune-backend, altalune-frontend, or altalune-chatbot skills
```

## Reference Files

- **[idea-workflow.md](references/idea-workflow.md)** - Complete user story creation workflow
- **[breakdown-workflow.md](references/breakdown-workflow.md)** - Task breakdown patterns and anti-patterns
- **[work-workflow.md](references/work-workflow.md)** - Task implementation checklist
