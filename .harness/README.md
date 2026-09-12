# Store Manager — Dedicated Harness (Go)

This harness specializes the Agentic SDLC Base Harness for the **Store Manager** project, a Go REST API for dropshipping sales management.

## Deviations from Base Harness

| Aspect | Base | Dedicated |
|---|---|---|
| Language | Go (generic) | Go (Vanilla: net/http, database/sql) |
| Framework | Generic | No framework — std library only |
| Database | Generic | MySQL via go-sql-driver/mysql |
| Architecture | Generic layered | Handlers/Services/Repository |
| Testing | Go testing/stdlib | Go testing + testify/sqlmock |
| Project skills | Go-focused | Go std lib + MySQL |

## Key files

| File | Purpose |
|---|---|
| `manifest.yaml` | Harness metadata and selection registry |
| `project-profile.md` | Domain and technical description |
| `architecture.md` | Architecture decision and rationale |
| `context/` | Domain glossary, vocabulary, project knowledge |
| `squads/catalog.yaml` | Specialization squad definitions |
| `agents/` | Agent role configurations |
| `skills/` | Project-specific skill overlays |
| `policies/` | Project-specific constraints |
| `workflows/` | Lifecycle workflow definitions |
| `tools/` | Tool capability configurations |

## How to operate

1. Start with `WF-001` (Discovery) if project context is incomplete.
2. Proceed through `WF-002` (Architecture) → `WF-004` (Planning) → `WF-005` (Implementation).
3. Review at `WF-006` (Verification) for each Story.
4. Release via `WF-007` (Release).
5. All merges follow the HITL policy — human gates are mandatory.

## Inherited policies

- HITL Policy (all mandatory gates apply)
- Branching Policy (main/develop/feature/story)
- Base Constitution (simplest good solution)