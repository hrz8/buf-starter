---
name: altalune-chatbot
description: |
  Chatbot module development for Altalune. Use when creating new chatbot configuration modules (non-entity JSONB-based configs). Covers the complete flow: (1) Proto schema definition with validation, (2) Code generation (Go, TypeScript, JSON Schema), (3) Backend defaults in repo.go, (4) Frontend module plugin system (defaults, metadata, auto-discovery). This is for configuration modules like llm, prompt, widget, mcpServer - NOT for standard CRUD entities.
---

# Chatbot Module Development

## Overview

Chatbot modules are **configuration plugins** stored as JSONB, not standard CRUD entities. Each project has ONE chatbot config row containing ALL module configs.

**Module vs Entity:**
| Aspect | Regular Entity | Chatbot Module |
|--------|---------------|----------------|
| Storage | Separate table rows | JSONB field in single config row |
| CRUD | Full Create/Read/Update/Delete | Get/Update only (lazy init) |
| Frontend | Repository + Service | Pinia Store + Plugin System |

## Quick Start: Add New Module

```bash
# 1. Create proto
# 2. Generate code
buf generate
# 3. Update backend defaults
# 4. Create frontend module files
# 5. Test
```

## Step-by-Step Guide

See **[module-development.md](references/module-development.md)** for complete walkthrough with examples.

## Existing Modules

| Module | Purpose | Proto |
|--------|---------|-------|
| `llm` | LLM provider/model config | `chatbot/modules/v1/llm.proto` |
| `prompt` | System prompt | `chatbot/modules/v1/prompt.proto` |
| `mcpServer` | MCP tool servers | `chatbot/modules/v1/mcp_server.proto` |
| `widget` | Embeddable widget CORS | `chatbot/modules/v1/widget.proto` |

## Key Files

**Backend:**
- `api/proto/chatbot/modules/v1/` - Proto schemas
- `internal/domain/chatbot/model.go` - `ValidModuleNames` map
- `internal/domain/chatbot/repo.go` - `defaultChatbotModulesConfig`

**Frontend:**
- `frontend/app/lib/chatbot-modules/` - Module plugin system
- `frontend/gen/jsonschema/` - Generated JSON schemas (auto-discovered)

## Reference Files

- **[module-development.md](references/module-development.md)** - Complete step-by-step guide
