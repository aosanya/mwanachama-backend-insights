# Architecture — mwanachama-backend-insights

## Schema

One table, one type:

```
Insight
  id           string, primary key
  repo         string, indexed, optional
  summary      string, required
  challenges   string, optional (free text)
  suggestions  string, optional (free text)
  created_at   string (RFC3339)
```

No `updated_at`, no soft-delete column — an Insight is append-only, the
same reasoning as `mwanachama-backend-assetmanager`'s `MovementRow`.

## Package layout

```
models/insight.go       — the Insight domain type
gormstore/insight.go    — InsightRow, ToRow/FromRow, BeforeCreate UUID mint
gormstore/tables.go     — TableNames/DefaultTableNames/Migrate
doc.go, tables.go       — root-package wrappers over gormstore
errors.go               — ErrInsightNotFound, ErrInvalidInsight
insight.go              — InsightManager interface + insightManager impl
mcp/mcp.go              — RegisterTools, ListResult[T], summary()
mcp/insight.go          — insight_create, insight_list tool registrations
```

Same split as `mwanachama-backend-assetmanager` and
`mwanachama-backend-agency`: domain types and GORM plumbing each live in
their own subpackage, the root package exposes one manager interface, and
`mcp/` is the only network-facing surface — no `routes/` package exists
yet because nothing needs a plain REST path today.

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
