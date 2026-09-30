# mwanachama-backend-insights — completed tasks

| Task | Title | Completed | Notes |
|------|-------|-----------|-------|
| I1 | Bootstrap repo: `go.mod`, `.gitignore`, `Makefile`, four-phase `documentation/` skeleton, `README.md`, `CLAUDE.md`, `todo.md`/`todo_done.md` | 2026-09-14 | New repo, module `github.com/aosanya/mwanachama-backend-insights`. Mirrors [mwanachama-backend-assetmanager](../../../mwanachama-backend-assetmanager)'s shape (GORM-direct, no HTTP layer of its own), scaled to one entity. |
| I2 | Write `1. requirements/requirements.md` + `2. design/architecture.md` — the 11-decision session record (no separate monitoring backend, no Hugging Face inference step, MCP write tool instead of a REST call, repo name/shape, why wakala-api not api-gateway hosts it, Insight field shape, append-only CRUD, no `ManagerResolver`, no `routes/`) | 2026-09-14 | See requirements.md's decision table for the full record and rationale — this session's back-and-forth (backend monitoring idea → Hugging Face → back to "just post via MCP") is preserved there, not just the final shape. |
| I3 | Write `models/insight.go`, `gormstore/insight.go`, `gormstore/tables.go`, root `doc.go`/`tables.go`/`errors.go`/`insight.go` (`InsightManager` interface, `insightManager` struct, `NewInsightManager`) | 2026-09-14 | `CreateInsight`/`GetInsight`/`ListInsights` only — no Update/Delete, matching assetmanager's Movement ledger. `Repo`/`Summary`/`Challenges`/`Suggestions` are plain text columns, not JSON-encoded arrays — an LLM caller writes prose into an MCP argument, not a hand-encoded array string. |
| I4 | Write `mcp/mcp.go` (`RegisterTools`, `ListResult[T]`, `summary()`) and `mcp/insight.go` (`insight_create`, `insight_list`) | 2026-09-14 | Follows `mwanachama-backend-agency/mcp`'s shape (AG13, the org's reference implementation) but skips its `ManagerResolver` indirection — that exists only because agency is mounted per-instance by a registry; insights has exactly one fixed instance, so tools close over a concrete `InsightManager` directly. |
| I5 | Unit tests (`testdb_test.go`+`insight_test.go`, `mcp/testdb_test.go`+`mcp/insight_test.go`, sqlite-backed via `glebarez/sqlite`) | 2026-09-14 | `go build ./...`, `go vet ./...`, `go test ./...` all clean. No Postgres integration test — nothing here needs one yet; add one later the same way assetmanager did if a real-Postgres-specific behavior ever needs covering. |
| I6 | Wired into `mwanachama-wakala-api`'s existing `/mcp` endpoint (not `mwanachama-backend-api-gateway`'s) | 2026-09-14 | api-gateway's `/mcp` is deliberately read-only by its own doc-comment convention; `insight_create` is a write, so this repo rides wakala-api's endpoint instead, which already has write-tool precedent (`agency_*`). Wiring kept thin per explicit instruction ("most of the code should be in the insight repo itself"): `go.mod` `replace`, one `InsightManager` in `cmd/server/main.go`'s `buildDeps` (instance label `"app"`), one `Deps` field, one `RegisterTools` call in `mcp.go`, one migration mirror (`000004_insights_tables.*.sql`) with `migrations_test.go`'s `wantVersions` extended to match. Full record of that side lives in wakala-api's own `CLAUDE.md` decision #17, not duplicated here. `go build`/`go vet`/`go test ./...` clean on that repo too. No `.mcp.json` change needed — rides the existing `wakala` entry. |

| I9 | Convert this repo to the declared-domain standard — `insights.blueprint.json` + `insights.operations.json`, storage through `mwanachama-backend-shared/specstore`, routes and MCP tools through `dispatch` | 2026-09-30 | Deleted `gormstore/`, root `tables.go`, five decode-call-encode handlers, `routes/wire.go`, the local `Route` type and `mcp/insight.go` — about 200 lines of structure the declarations now carry. Renamed `AgencyID`/`DraftID` to `SubjectID`/`ContextID` through to the wire on the owner's call (requirements #18); `Repo` kept as the one documented neutrality exception (#19). The published `insight_create`/`insight_list` tool names are preserved by the operations file's `tool` field and pinned by a real-MCP-client test; `insight_note_create`/`insight_note_list` arrive free, which is half of I8. Six guard tests added and each one proven to go red when the thing it guards is broken. Live rows adopted from `app_insights`/`app_insight_notes` by `Provision`, proven on both sqlite and real Postgres. `go test ./...` and `go vet -tags=integration ./...` clean. Left open: I10 (the studio's wire names) and I11 (uncapped lists, found in passing). |

## Archived board context

### Why this repo exists

The user had an idea mid-session: record the challenges/repetitive tasks a
Claude Code session hits, so they inform future shortcuts/tooling instead
of just recurring. First framed as "maybe a separate backend to monitor
token usage" — talked through and dropped (Claude Code already writes
local transcripts; token usage is already visible via Anthropic Console).
Then "use Hugging Face" — narrowed through a few rounds (inference vs.
Spaces vs. dataset storage; hosted vs. self-hosted; free-tier CPU vs. a
k8s deployment) before the user cut straight to the actual ask: "I just
want a summary of the chat posted using MCP... to add chat insights to
the database." That reframing is what this repo and its wakala-api wiring
actually implement — no external model, no separate hosting, just an MCP
write tool a session calls directly.

**Not yet built, on purpose**: the actual trigger (a Stop hook or slash
command that has a session call `insight_create` automatically at the
end). Deferred so the tool can be exercised manually first.
