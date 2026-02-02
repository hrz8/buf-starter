# Idea to User Story Command

**Variables:** `idea_summary` = $ARGUMENTS

## Quick Start

1. Explore codebase using Task tool with `subagent_type=Explore`
2. Review `altalune-workflow` skill for story patterns
3. Ask clarifying questions using `AskUserQuestion`
4. Create user story in `docs/stories/US{N}-{slug}.md`

## Story Numbering

```bash
ls docs/stories/US*.md | sort -V | tail -1  # Get last story number
```

## Question Categories

- User & business context
- Functional requirements (CRUD, validation, business rules)
- Technical architecture (new domain vs extend, database)
- Security & access control
- User experience
- Scope boundaries (in/out of scope)

**For detailed guidance:** Use `altalune-workflow` skill for story template and question checklist.
