# Idea to User Story Workflow

**OUTPUT: User story documentation file ONLY. Do NOT implement code.**

## Phase 1: Discovery & Context Gathering

**Before asking questions, understand the codebase:**

1. **Explore Codebase:**
   - Use Task tool with `subagent_type=Explore`
   - Identify similar features already implemented
   - Understand current architecture patterns
   - Map out related components and services

2. **Review Skills for Patterns:**
   - `altalune-backend` skill for domain patterns
   - `altalune-frontend` skill for UI patterns
   - `altalune-authorization` skill for auth requirements
   - `altalune-chatbot` skill (if chatbot-related)
   - `CLAUDE.md` for architecture overview

3. **Analyze Existing Stories:**
   - Review `docs/stories/` for format and depth
   - Note acceptance criteria structure

## Phase 2: Iterative Clarification

**Use AskUserQuestion repeatedly until ALL aspects are clear.**

### Question Categories

#### User & Business Context
- Who is the target user?
- What is the primary business goal?
- What problem does this solve?
- Expected user workflow?

#### Functional Requirements
- Core CRUD operations needed?
- Data to display/capture?
- Search/filter/sort requirements?
- Validation rules?
- Business rules and constraints?

#### Technical Architecture
- New domain or extend existing?
- Reusable components/patterns?
- New protocol buffers needed?
- Database migrations?
- Project-scoped (partitioned)?

#### Security & Access Control
- Authentication requirements?
- Authorization rules (who can do what)?
- Project-level or user-level permissions?
- Audit/logging requirements?

#### User Experience
- UI/UX pattern (table, form, wizard)?
- Existing UI components to use?
- Responsive requirements?
- Feedback/notifications needed?

#### Scope & Boundaries
- What is IN scope?
- What is OUT of scope?
- Dependencies on other features?

## Phase 3: User Story Creation

### Determine Story Number

```bash
ls docs/stories/US*.md | sort -V | tail -1  # Get last story number
```

### Story Template

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

### Security Requirements

#### {Security Aspect}

- {Security requirement}

### Data Validation

#### {Field/Entity Name}

- {Validation rule}

### User Experience

#### Responsive Design

- {UX requirement}

#### Feedback and Notifications

- {UX requirement}

## Technical Requirements

### Backend Architecture

- {Reference patterns from BACKEND_GUIDE.md}
- {Database considerations}
- {Error handling requirements}

### Frontend Architecture

- {Reference patterns from FRONTEND_GUIDE.md}
- {Repository/service patterns}
- {Component architecture}

### API Design

- {Endpoints and methods}
- {Request/response structure}

## Out of Scope

- {Feature NOT included}
- {Deferred to future stories}

## Dependencies

- {Required existing features}
- {Third-party services}

## Definition of Done

- [ ] {Specific completion criterion}
- [ ] All CRUD operations implemented and tested
- [ ] Security requirements met
- [ ] Data validation comprehensive
- [ ] UI responsive and accessible
- [ ] Error handling comprehensive
- [ ] Code follows established patterns
- [ ] Tests written and passing
```

## Quality Checklist

- [ ] Story follows "As a... I want... So that..." format
- [ ] Acceptance criteria use Given/When/Then format
- [ ] Technical requirements reference architectural patterns
- [ ] Security considerations addressed
- [ ] Scope clearly defined (in-scope and out-of-scope)
- [ ] Dependencies identified
- [ ] Definition of Done is specific and measurable

## DO's

- Always explore codebase first
- Ask clarifying questions iteratively
- Reference existing patterns
- Include Given/When/Then for acceptance criteria
- Explicitly define out of scope
- Create comprehensive Definition of Done

## DON'Ts

- Skip exploration phase
- Ask generic questions without context
- Proceed without complete clarity
- Create vague acceptance criteria
- Omit security considerations
- Skip out-of-scope section
