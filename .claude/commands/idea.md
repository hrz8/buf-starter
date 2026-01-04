# Idea to User Story Command

**Variables:** `idea_summary` = $ARGUMENTS

## Purpose

Transform a simple idea or feature request into a comprehensive, well-structured user story following standard architecture and patterns.

## Workflow

### Phase 1: Discovery & Context Gathering (REQUIRED)

**Objective**: Understand the codebase deeply before asking questions

1. **Explore Codebase**:

   - Use Task tool with `subagent_type=Explore` to understand relevant parts of the codebase
   - Identify similar features or patterns already implemented
   - Understand current architecture patterns (backend domains, frontend repositories/services)
   - Map out related components, services, and database tables
   - Identify reusable code and established patterns

2. **Review Guidelines**:

   - Read `docs/dev_guidelines/BACKEND_GUIDE.md` for backend patterns
   - Read `docs/dev_guidelines/FRONTEND_GUIDE.md` for frontend patterns
   - Read `docs/dev_guidelines/DOMAIN_ARCHITECTURE_GUIDE.md` for domain structure
   - Read `docs/dev_guidelines/EFFICIENCY_GUIDE.md` for best practices
   - Read `CLAUDE.md` for architecture overview

3. **Analyze Existing Stories**:
   - Review existing user stories in `docs/stories/` for format and depth
   - Understand the level of detail expected
   - Note patterns in acceptance criteria structure

### Phase 2: Iterative Clarification (REQUIRED)

**Objective**: Ask targeted questions based on discovery findings

**CRITICAL**: Use `AskUserQuestion` tool repeatedly until ALL aspects are clear. Do NOT proceed to Phase 3 until you have complete clarity.

**Question Categories to Cover**:

#### User & Business Context

- Who is the target user for this feature?
- What is the primary business goal?
- What problem does this solve?
- What is the expected user workflow?
- Are there any compliance or regulatory requirements?

#### Functional Requirements

- What are the core CRUD operations needed?
- What data needs to be displayed/captured?
- What are the search/filter/sort requirements?
- What are the validation rules?
- What are the business rules and constraints?
- Are there any batch operations needed?

#### Technical Architecture

- Does this need a new backend domain or extend existing?
- Which existing components/patterns should be reused?
- Does this require new protocol buffers or extend existing?
- Does this need database migrations or use existing tables?
- Should this be partitioned by project_id? (if multi-tenant)
- What are the performance requirements?

#### Security & Access Control

- What are the authentication requirements?
- What are the authorization rules (who can do what)?
- Are there project-level or user-level permissions?
- What data should be encrypted or protected?
- Are there any audit/logging requirements?

#### User Experience

- What is the UI/UX pattern (table, form, wizard, dashboard)?
- Which existing UI components should be used?
- What are the responsive design requirements?
- What feedback/notifications are needed?
- What are the error handling expectations?

#### Scope & Boundaries

- What is explicitly IN scope for this story?
- What is explicitly OUT of scope (for future stories)?
- What are the dependencies on other features?
- Are there any breaking changes or migration needs?

#### Definition of Done

- What are the quality standards for completion?
- What testing is required?
- What documentation needs to be updated?
- Are there any deployment considerations?

**Clarification Strategy**:

1. Start with high-level questions (user, business goal)
2. Then drill into functional requirements
3. Follow with technical architecture decisions
4. Cover security and access control
5. Discuss UX and interface details
6. Finalize scope and boundaries
7. Confirm definition of done criteria

**Use Multi-Question Batches**: When asking questions, group related questions together using `AskUserQuestion` with multiple questions to be efficient.

### Phase 3: User Story Creation (Only After Complete Clarity)

**CRITICAL**: Only proceed when you have comprehensive answers to all questions from Phase 2.

1. **Determine Story Number**:

   - List all existing stories in `docs/stories/`
   - Find the highest US number (e.g., US1, US2, etc.)
   - Increment by 1 for the new story number

2. **Generate Filename**:

   - Format: `US{number}-{slug}.md`
   - Slug: lowercase, hyphen-separated, 2-4 words describing the feature
   - Example: `US2-employee-management.md`, `US3-reporting-dashboard.md`

3. **Create User Story Document** following this template:

```markdown
# User Story US{N}: {Title}

## Story Overview

**As a** {user role}
**I want** {capability or goal}
**So that** {business value or benefit}

## Acceptance Criteria

### Core Functionality

#### {Feature/Operation Name 1}

- **Given** {context or precondition}
- **When** {action or event}
- **Then** {expected outcome}
- **And** {additional outcomes}

#### {Feature/Operation Name 2}

- **Given** {context}
- **When** {action}
- **Then** {outcome}
- **And** {additional outcomes}

{Repeat for each major feature/operation}

### Security Requirements

#### {Security Aspect 1}

- {Security requirement}
- {Security requirement}

#### {Security Aspect 2}

- {Security requirement}
- {Security requirement}

### Data Validation

#### {Field/Entity Name 1}

- {Validation rule}
- {Validation rule}

#### {Field/Entity Name 2}

- {Validation rule}
- {Validation rule}

### User Experience

#### Responsive Design

- {UX requirement}
- {UX requirement}

#### Feedback and Notifications

- {UX requirement}
- {UX requirement}

#### Integration with Existing UI

- {UX requirement}
- {UX requirement}

## Technical Requirements

### Backend Architecture

- {Technical requirement - reference to patterns from BACKEND_GUIDE.md}
- {Database considerations - partitioning, migrations, etc.}
- {Error handling and logging requirements}
- {Dual ID system usage if applicable}
- {Protovalidate integration}

### Frontend Architecture

- {Technical requirement - reference to patterns from FRONTEND_GUIDE.md}
- {Repository layer details}
- {Service composable patterns}
- {Component architecture}
- {Form and validation approach}

### API Design

- {API endpoints and methods}
- {Request/response structure}
- {Status codes and error handling}
- {Pagination, filtering, sorting support}

## Out of Scope

- {Feature or capability NOT included in this story}
- {Feature deferred to future stories}
- {Functionality that won't be implemented}

## Dependencies

- {Required existing features or infrastructure}
- {Third-party services or libraries}
- {Database schema requirements}
- {Authentication/authorization dependencies}

## Definition of Done

- [ ] {Specific completion criterion}
- [ ] {Specific completion criterion}
- [ ] All CRUD operations (if applicable) are implemented and tested
- [ ] Security requirements are met
- [ ] Data validation is comprehensive
- [ ] User interface is responsive and accessible
- [ ] Error handling is comprehensive
- [ ] Code follows established patterns and guidelines
- [ ] Unit and integration tests are written
- [ ] Documentation is updated
- [ ] Code is reviewed and approved
- [ ] Feature is deployed and verified in staging environment

## Notes

- {Additional context or considerations}
- {Links to related stories or documentation}
- {Future enhancement ideas}
- {Technical debt or trade-offs}
```

4. **Quality Checklist**:

   - [ ] Story follows the "As a... I want... So that..." format
   - [ ] Acceptance criteria use Given/When/Then format
   - [ ] Technical requirements reference specific architectural patterns
   - [ ] Security considerations are explicitly addressed
   - [ ] Scope is clearly defined (in-scope and out-of-scope)
   - [ ] Dependencies are identified
   - [ ] Definition of Done is specific and measurable
   - [ ] Story is comprehensive but focused on a single feature/epic

5. **Save User Story**:
   - Create file at `docs/stories/US{number}-{slug}.md`
   - Ensure proper formatting and line breaks
   - Verify all sections are complete

## Important Reminders

### DO's:

- ✅ ALWAYS explore the codebase first before asking questions
- ✅ ALWAYS read relevant dev guidelines before clarifying
- ✅ ALWAYS ask clarifying questions iteratively using AskUserQuestion
- ✅ ALWAYS reference existing patterns and components discovered during exploration
- ✅ ALWAYS make questions specific based on codebase findings
- ✅ ALWAYS ensure complete clarity before creating the user story
- ✅ ALWAYS follow the US1 example format for structure and depth
- ✅ ALWAYS include Given/When/Then format for acceptance criteria
- ✅ ALWAYS explicitly define what's out of scope
- ✅ ALWAYS create comprehensive Definition of Done checklist

### DON'Ts:

- ❌ NEVER skip the exploration phase
- ❌ NEVER ask generic questions without codebase context
- ❌ NEVER proceed to story creation without complete clarity
- ❌ NEVER create vague or incomplete acceptance criteria
- ❌ NEVER forget to reference architectural patterns from guidelines
- ❌ NEVER omit security considerations
- ❌ NEVER skip the out-of-scope section
- ❌ NEVER create stories without clear business value

## Success Criteria

A successful user story:

1. Is grounded in deep understanding of the existing codebase
2. References specific patterns and components already in use
3. Has crystal-clear requirements from iterative clarification
4. Follows the established format and depth of US1
5. Provides comprehensive acceptance criteria with Given/When/Then
6. Addresses all technical, security, and UX considerations
7. Clearly defines scope boundaries
8. Has measurable definition of done
9. Can be directly implemented by following the story alone

## Example Usage

```bash
# Simple idea
claude /idea "add employee directory feature"

# More specific idea
claude /idea "implement real-time notifications for project updates"

# Brief concept
claude /idea "user preferences and settings management"
```

**Goal**: Transform raw ideas into implementation-ready user stories that fully leverage Altalune's architecture and can guide development without ambiguity.
