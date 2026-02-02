# User Story Breakdown Command

**Variables:** `user_story_path` = $ARGUMENTS

## Quick Start

1. Use `EnterPlanMode` tool to begin
2. Read user story from `docs/stories/$ARGUMENTS.md`
3. Review `altalune-workflow` skill for breakdown patterns
4. Create 3-8 tasks in `docs/tasks/T{N}-{slug}.md`
5. Use `ExitPlanMode` when done

## Key Principles

- **Efficiency over granularity** - Group related work together
- **Strict scope adherence** - Only implement what's in user story
- **No duplication** - Clear, unique scope per task
- **Clear dependencies** - Mark prerequisites

## Task Numbering

```bash
ls docs/tasks/T*.md | sort -V | tail -1  # Get last task number
```

**For detailed guidance:** Use `altalune-workflow` skill for task templates and anti-patterns.
