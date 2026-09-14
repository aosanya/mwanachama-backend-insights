# mwanachama-backend-insights

Records session insights — a short summary of a working session plus any
repetitive tasks/friction ("challenges") and shortcut/tooling ideas
("suggestions") it surfaced. No gRPC, no sub-service shape.

Storage is GORM directly: `models/` holds the domain type, `gormstore/`
holds the row struct and migration. `InsightManager` is append-only —
`CreateInsight`/`GetInsight`/`ListInsights` only, no update or delete,
mirroring `mwanachama-backend-assetmanager`'s Movement ledger.

Ships an `mcp/` package exposing `insight_create`/`insight_list` as Model
Context Protocol tools, following `mwanachama-backend-agency/mcp`'s shape
(the org's reference implementation) — see that package's doc comment for
why this one skips the `ManagerResolver` indirection agency needs.
`mwanachama-wakala-api` mounts these tools on its existing `/mcp` endpoint;
this repo has no server of its own.

See [documentation/](documentation/) for the requirements write-up and
[documentation/3. implementation/todo_done.md](documentation/3.%20implementation/todo_done.md)
for what shipped.
