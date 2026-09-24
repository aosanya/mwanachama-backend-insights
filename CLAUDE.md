# CLAUDE.md

Guidance for Claude Code working in this repository.

## Project: mwanachama-backend-insights

Records session insights — a short summary of a working session plus any
repetitive tasks/friction ("challenges") and shortcut/tooling ideas
("suggestions") it surfaced, so a later pass can mine them for real
product/tooling work instead of the friction just recurring silently every
session. Module path `github.com/aosanya/mwanachama-backend-insights`.
Sibling of [mwanachama-backend-assetmanager](../mwanachama-backend-assetmanager)
in spirit (same GORM-direct, no-HTTP-of-its-own shape), scaled down to one
entity.

**Standalone library, no server of its own.** `InsightManager` is
append-only — `CreateInsight`/`GetInsight`/`ListInsights` only, no
Update/Delete, the same as `mwanachama-backend-assetmanager`'s Movement
ledger. `models/` holds the domain type, `gormstore/` holds the row struct
and migration — same split as every other GORM-backed repo in this family.

**MCP tools live here, not in the mounting host** — `mcp/` exposes
`insight_create`/`insight_list`, following `mwanachama-backend-agency/mcp`'s
shape (the org's reference implementation for "a library's own manager
methods as MCP tools", see `mwanachama-backend-actor/CLAUDE.md`'s "MCP
tools" section). One deliberate simplification vs. agency's `mcp` package:
no `ManagerResolver` indirection. Agency needs that because
`mwanachama-wakala-api` mounts *many* Agencies behind one server; this
repo has exactly one fixed Insights instance for the whole deployment, so
`mcp.RegisterTools` takes a concrete `InsightManager` directly.

**Hosted on `mwanachama-wakala-api`, not `mwanachama-backend-api-gateway`.**
The gateway's `/mcp` (`git_*`/`taskmanager_*`) is deliberately read-only by
its own doc-comment convention; `insight_create` is a write, so this repo
rides wakala-api's endpoint instead — it already has write-tool precedent
(`agency_*`) and is already running and registered in the workspace
`.mcp.json`. See wakala-api's own `CLAUDE.md` for the wiring: one `go.mod`
dependency, one `InsightManager` constructed alongside its existing DB
setup, one field on `Deps`, one `RegisterTools` call, one migration
mirror. That wiring is deliberately thin — "most of the code should be in
the insight repo itself" was an explicit instruction, mirroring AG13's
same rationale for moving agency's MCP tools out of wakala-api.

**`routes/` REST package added 2026-09-14** — a human-facing surface
(`mwanachama-wakala-studio`'s agency-chat screen) needed to create/read
insights directly, not through an AI agent's MCP call. Mirrors
`mwanachama-backend-assetmanager/routes`'s shape exactly: `Route`/
`Route.Pattern`, `InsightRoutes`/`InsightNoteRoutes`/`Routes`, decode-call-
encode handlers, no caller-identity gate of its own (the mounting process —
`mwanachama-wakala-api`'s `requireCaller` — wraps every route). `mcp/` is
unchanged and remains the AI-agent surface.

**`InsightNote` (added 2026-09-14)** answers "how does an insight grow
without duplicating itself": a follow-up remark on an insight already on
record (auto- or user-captured) is appended as a Note on that same row —
`CreateInsightNote`/`ListInsightNotes`, on the same `InsightManager`
interface (Note is subordinate to Insight, the same way `AssetManager`
carries Movement's methods directly rather than a separate manager) —
instead of spawning a near-duplicate Insight that repeats its Summary.
`Insight` also gained `AgencyID`/`DraftID` (which wakala Agency/Draft an
insight concerns, both empty for a dev-session insight) and `Source`
(`SourceAuto`/`SourceUser`; empty on create defaults to `SourceAuto`, so
the existing `insight_create` MCP tool needed zero changes) and `Tags`
(comma-separated free text, same plain-string convention as
`Asset.AttributesJSON`). See
[documentation/1. requirements/requirements.md](documentation/1.%20requirements/requirements.md)'s
2026-09-14 decision block for the full record.

**Verification default**: `go test ./...` (sqlite-backed via
`glebarez/sqlite`) is the expected way to verify a change here — no
Postgres container needed for this repo's own tests. See
[[feedback_use_memory_backend_for_tests]].

## Open questions

None open — see
[documentation/1. requirements/requirements.md](documentation/1.%20requirements/requirements.md)'s
decision table. One deliberate follow-up, not an open design question:
actually wiring something (a Stop hook, a slash command) so a session
calls `insight_create` automatically at the end — not built yet, on
purpose, so the tool can be exercised manually first.

## Conventions

- Four-phase `documentation/` layout — see
  [documentation/README.md](documentation/README.md).
- Before adding new persistence code, read
  `mwanachama-backend-assetmanager`'s `gormstore/` package and
  `asset_manager.go`'s `AssetManager` interface as the reference shape
  this repo's GORM layer was ported from.
- Before touching `mcp/`, read `mwanachama-backend-agency/mcp/mcp.go`'s
  doc comment and `mwanachama-backend-actor/CLAUDE.md`'s "MCP tools"
  section — this repo follows that standard minus the `ManagerResolver`
  indirection (see above for why it doesn't need one).

## Code comments

Write code with no comments. Not one-liners above a function, not section
banners, not doc comments on exported symbols, not "why" notes next to a
tricky line. A name, a type, or a smaller function carries it instead.

Anything that genuinely needs explaining goes in this repo's `documentation/`
folder, under the phase it belongs to (`1. requirements`, `2. design`,
`3. implementation`, `4. qa`) — never inline.

**Why:** inline prose drifts out of sync with the code, duplicates what
`documentation/` already owns, and buries the explanation where nobody
looking for it will search.

**How to apply:**

- New code ships without comments. If a line seems to need one, rename or
  split until it doesn't.
- Touching code that already has comments: strip the ones in the code you are
  changing. Do not sweep untouched files unless asked.
- If the reasoning matters, add or update the matching `documentation/` page
  in the same change and leave nothing behind in the source.
- Machine-read directives are not comments and stay: build tags, `//go:embed`,
  `//go:generate`, linter pragmas (`//nolint`, `// eslint-disable-next-line`,
  `// ignore:`), license headers, codegen "do not edit" banners, and generated
  files as a whole.
- Commit messages, PR descriptions, and test names carry the narration that
  used to go in comments.

This rule is repeated verbatim in every mwanachama repo's `CLAUDE.md` so that
it reaches sessions that do not load this machine's user-level config —
scheduled cloud routines, other machines, and other agent harnesses.
