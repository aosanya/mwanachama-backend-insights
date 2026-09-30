# CLAUDE.md

Guidance for Claude Code working in this repository.

## Project: mwanachama-backend-insights

Records insights — a short summary of a stretch of work plus any repetitive
tasks/friction ("challenges") and shortcut/tooling ideas ("suggestions") it
surfaced, so a later pass can mine them for real product/tooling work
instead of the friction just recurring silently every session. Module path
`github.com/aosanya/mwanachama-backend-insights`.

**A declared domain since 2026-09-30 (I9).** The objects come from
`insights.blueprint.json` and the route table from
`insights.operations.json`; storage is
`mwanachama-backend-shared/specstore` and the routes are built by
`mwanachama-backend-shared/dispatch`. There are no row structs and no
`AutoMigrate`. This is the same shape `mwanachama-backend-catalog`,
`-agency`, `-permissions`, `-actor`, `-forms`, `-taskmanager` and
`-accounting` took; catalog is the reference and
`developer/documentation/2. design/architecture-spec-driven-modules.md` is
the strategy.

**Standalone library, no server of its own.** `InsightManager` is
append-only — no Update, no Delete, for either object. The way an insight
grows is a note beside it, never an edit.

## The rule that matters most: this module names no domain

**A word that means something in one domain and nothing in another does not
belong in this module.** It records that an observation was made; what the
observation is *about* belongs to whoever mounted it.

Two fields failed that test on 2026-09-30 and were renamed, wire included:

| Was | Why it failed | Is |
| --- | ------------- | -- |
| `Insight.AgencyID` / `agency_id` | a wakala word | `SubjectID` / `subject_id` |
| `Insight.DraftID` / `draft_id` | same | `ContextID` / `context_id` |

**`Repo` is the one documented exception.** A clinic has no repos, so it
fails the rule — but `repo` is a published argument of the live
`insight_create` tool, and a production tool's argument names may not change
across a conversion. It is deliberately absent from
`domain_agnostic_test.go`'s word list and described neutrally in the
blueprint. See requirements decision #19; do not "fix" it without reading
that.

`domain_agnostic_test.go` enforces the rest, over identifiers in `models/`,
`routes/`, `mcp/` and the root package, plus a second test over every
declared field, index and **stored enum value** — a stored enum value
outlives a rename, which is the worst version of this failure.

**MCP and REST are one table.** `mcp/mcp.go` registers what
`dispatch.Tools` builds from `insights.operations.json`, so an address
cannot reach one surface and miss the other. That had already happened:
`insight_note_create`/`insight_note_list` existed as REST handlers from
2026-09-14 and were still missing from `mcp/` when the conversion started.

**The published tool names are frozen.** `insight_create` and
`insight_list` are in the workspace `.mcp.json`. `dispatch` would have
derived `insights_insight_create` from the action id, so both operations
declare an explicit `tool` field. `mcp/insight_test.go` calls them by name
through a real MCP client, so a rename fails the suite rather than the
deployment. See [documentation/2. design/routes.md](documentation/2.%20design/routes.md).

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
insights directly, not through an AI agent's MCP call. It has no
caller-identity gate of its own: the mounting process
(`mwanachama-wakala-api`'s `requireCaller`, plus its own
`scopeInsightsToMembership`) wraps every route.

**`InsightNote` (added 2026-09-14)** answers "how does an insight grow
without duplicating itself": a follow-up remark on an insight already on
record (auto- or user-captured) is appended as a Note on that same row —
`CreateInsightNote`/`ListInsightNotes`, on the same `InsightManager`
interface (Note is subordinate to Insight, so it needs no manager of its
own) —
instead of spawning a near-duplicate Insight that repeats its Summary.
`Insight` also gained `SubjectID`/`ContextID` (renamed from
`AgencyID`/`DraftID` in I9 — which subject and revision an insight concerns,
both empty for an observation about the work itself) and `Source`
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

## What is superseded

All of it on 2026-09-30, by I9:

- **`gormstore/`** — `InsightRow`, `InsightNoteRow`, the four
  `*ToRow`/`*FromRow` converters, the `BeforeCreate` id hooks and
  `Migrate`'s `AutoMigrate` calls. `spec.Migrate` and `specstore`'s codec
  replace them, and root `tables.go` with its
  `TableNames`/`DefaultTableNames`/`Migrate` re-exports went with it.
  `NewInsightManager(db, *spec.Spec)` is the constructor now.
- **`routes/insight.go`, `routes/insightnote.go`, `routes/wire.go`,
  `routes/doc.go`** — five decode-call-encode handlers, `insightStatusFor`,
  `writeInsightErr`, the private `createInsightNoteBody`, and a local
  `writeJSON`/`writeErr`/`readJSON` trio that its own comment admitted was
  a byte-for-byte copy of assetmanager's. `httpwire` had all three.
- **The local `Route` struct and its `Pattern` method** — `httpwire.Route`,
  aliased, so no mounting process had to change.
- **`mcp/insight.go`** — `registerInsightTools` and the two hand-written
  param structs. `dispatch.Tools` derives both from the same operations
  file the REST table comes from.
- **The per-field doc comments in `models/`** — moved into the blueprint's
  `description` fields, which is where that prose lives now.

## Conventions

- Four-phase `documentation/` layout — see
  [documentation/README.md](documentation/README.md).
- `go test ./...` (sqlite via `glebarez/sqlite`) verifies a change here.
  The driver-specific provisioning lives in `postgres_integration_test.go`
  (`//go:build integration`, gated on `POSTGRES_URL`, run by `make
  test-pg`), which is **not** part of `go test ./...`.
- **Adding a field means editing `insights.blueprint.json` and the Go type
  together.** `TestEveryExampleFitsTheTypes` builds a manager over every
  shipped spec, so a column with no field to hold it fails there rather
  than dropping a value on every write. Never add a field to a domain spec
  to get it into the module.
- **A domain names objects; it does not re-declare them.** A domain spec
  supplies `instance`, the name and table each role lands in, its own
  indexes, and a default on a declared field — nothing else.
- **Assert on `spec.RawNameFor`, never on a physical table name.** The
  physical name is hashed; see `mwanachama-backend-shared`'s
  [declared-domains.md](../mwanachama-backend-shared/documentation/2.%20design/declared-domains.md).
- There are no hand-written route builders left to name: an address is an
  entry in `insights.operations.json`.
- This module has **no auth model**. A route needs a capability gate
  wrapped around it by whatever mounts it.

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
