# Architecture — mwanachama-backend-insights

## Schema

Two tables, two types — `InsightNote` (added 2026-09-14) is how an insight
grows over time without duplicating itself:

```
Insight
  id           string, primary key
  repo         string, indexed, optional
  agency_id    string, indexed, optional
  draft_id     string, optional
  source       string, "auto" or "user"
  summary      string, required
  tags         string, optional (comma-separated free text)
  challenges   string, optional (free text)
  suggestions  string, optional (free text)
  created_at   string (RFC3339)

InsightNote
  id           string, primary key
  insight_id   string, indexed, required (references Insight.id)
  source       string, "auto" or "user"
  text         string, required
  created_at   string (RFC3339)
```

No `updated_at`, no soft-delete column on either — both are append-only,
the same reasoning as `mwanachama-backend-assetmanager`'s `MovementRow`.
`repo`/`challenges`/`suggestions` are the original Claude-Code-session
shape; `agency_id`/`draft_id` scope a wakala-chat insight instead — a row
uses one set or the other, never both.

## Package layout

```
models/insight.go       — the Insight domain type, Source constants
models/insightnote.go   — the InsightNote domain type
gormstore/insight.go    — InsightRow, ToRow/FromRow, BeforeCreate UUID mint
gormstore/insightnote.go — InsightNoteRow, same shape
gormstore/tables.go     — TableNames/DefaultTableNames/Migrate (both tables)
doc.go, tables.go       — root-package wrappers over gormstore
errors.go               — ErrInsightNotFound, ErrInvalidInsight, ErrInvalidInsightNote
insight.go              — InsightManager interface + insightManager impl
insightnote.go          — CreateInsightNote/ListInsightNotes methods
mcp/mcp.go              — RegisterTools, ListResult[T], summary()
mcp/insight.go          — insight_create, insight_list tool registrations
routes/                 — plain REST surface for a human caller (added
                           2026-09-14): doc.go, wire.go, routes.go
                           (Route/Route.Pattern/InsightRoutes/
                           InsightNoteRoutes/Routes), insight.go,
                           insightnote.go
```

Same split as `mwanachama-backend-assetmanager` and
`mwanachama-backend-agency`: domain types and GORM plumbing each live in
their own subpackage, the root package exposes one manager interface.
Two network-facing surfaces now exist side by side: `mcp/` for an AI
agent, `routes/` for a human caller (`mwanachama-wakala-studio`) — neither
depends on the other, and `InsightNote` has no MCP tools yet (deliberately
deferred, see requirements.md's 2026-09-14 decisions).

## How this reaches a live MCP client

This repo has no server of its own. `mwanachama-wakala-api`:

1. Depends on this module (local `replace` directive, same as its
   `mwanachama-backend-agency` dependency).
2. At boot, runs `Migrate` and constructs one `InsightManager` against its
   existing `*gorm.DB`, using instance label `"app"` (table `app_insights`)
   — a fixed, compile-time instance, the same style
   `mwanachama-backend-git`/`mwanachama-backend-taskmanager` use when
   mounted on `mwanachama-backend-api-gateway`, not agency's per-registry
   random slug (there is exactly one Insights instance, ever, for this
   deployment).
3. Calls `insightsmcp.RegisterTools(server, insightsManager)` inside its
   existing `newMCPServer()`, alongside `agencymcp.RegisterTools`.
4. Adds a hand-written SQL migration mirror
   (`internal/store/postgres/migrations/000004_insights_tables.*.sql`) so
   `cmd/migrate up` alone can provision a fresh DB, matching every other
   GORM domain's requirement in this codebase family.

No change to the container's `.mcp.json` — `insight_create`/`insight_list`
ride the existing `wakala` MCP server entry, same URL, same
`WAKALA_API_KEY` bearer auth as every `agency_*` tool.

## Why not api-gateway

`mwanachama-backend-api-gateway`'s `/mcp` (`git_*`/`taskmanager_*` tools)
is deliberately read-only — every tool there wraps a read method on its
underlying manager, by explicit doc-comment convention. `insight_create`
is a write. Rather than making it the first write-capable exception to
that convention, this repo rides `mwanachama-wakala-api`'s endpoint
instead, which already has write-tool precedent (`agency_create_*`) and
is already running and registered.
