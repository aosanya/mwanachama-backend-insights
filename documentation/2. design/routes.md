# routes/ and mcp/

Neither package holds a handler any more. The route table is declared in
`insights.operations.json` and built by
[`mwanachama-backend-shared/dispatch`](../../../mwanachama-backend-shared/documentation/2.%20design/dispatcher.md).
`routes/routes.go` is the sentinel map and the builder ladder;
`routes/tools.go` describes the fields; `mcp/mcp.go` registers the same
tools onto an MCP server.

Deleted in the 2026-09-30 conversion (I9): `routes/insight.go`,
`routes/insightnote.go`, `routes/wire.go`, `routes/doc.go` and
`mcp/insight.go` — about 200 lines of decode-call-encode and
hand-registered tool params.

## One table, two surfaces

`insights.operations.json` is the only place an address is written down.
`dispatch.Dispatch` turns it into `[]httpwire.Route` and `dispatch.Tools`
turns the same five operations into MCP tools, so a new address cannot
reach REST and miss MCP, which is what happened while the two were written
by hand: `insight_note_create` and `insight_note_list` existed as REST
handlers from 2026-09-14 and were still missing from `mcp/` when the
conversion started. They arrive free now.

## Why `routes/wire.go` and the local `Route` went

`routes/wire.go`'s own comment said it mirrored
`mwanachama-backend-assetmanager/routes`' copy "byte-for-byte on purpose",
and `routes/routes.go` carried its own `Route` struct with its own
`Pattern` method. `mwanachama-backend-shared/httpwire` already had all
three. Deleting them was a pure substitution — `Route` is now a type alias,
so nothing that mounted this package had to change.

## The published tool names are frozen

`insight_create` and `insight_list` are registered in the workspace
`.mcp.json` and called by dev sessions, so renaming either would break
every caller that named it. `dispatch` derives a tool's name from its
action id, which would have produced `insights_insight_create`. The
`tool` field on an operation overrides that, and both operations declare
it. `mcp/insight_test.go` calls the tools by name through a real MCP
client, so a rename fails the suite rather than the deployment.

The two new note operations take `insight_note_create` and
`insight_note_list`, which is what I8 on the board asked for.

## The address outranks the body

`create_note` binds `{from: path, as: insightID, into: InsightID}`, applied
after the body is decoded, so a caller cannot `POST
/insights/a/notes` with a body claiming `insight_id: b`. The old handler
did this by hand with a private `createInsightNoteBody` struct that simply
had no `insight_id` field; the declared binding is the same guarantee
stated once.
`TestUpsertInsight_TheAddressOutranksTheBody` holds it.

## Status codes, held exactly

The deleted `insightStatusFor` mapped three sentinels and defaulted to 500
with the body `"internal error"` rather than the sentinel's own text.
`insights.operations.json`'s `errors` map reproduces it:

| Sentinel | Status |
| --- | --- |
| `ErrInsightNotFound` | 404 |
| `ErrInvalidInsight` | 400 |
| `ErrInvalidInsightNote` | 400 |

Anything unmapped is redacted to 500 by the dispatcher, which is what the
old `writeInsightErr` did.

## Nothing is anonymous

`AnonymousActions` is empty, and `TestNothingIsAnonymous` pins it. Every
door here is the mounting host's to gate — `mwanachama-wakala-api` wraps
them in `requireCaller` plus its own `scopeInsightsToMembership`. The
allowlist names what is public, never what is protected, so an operation
added later and not named there arrives gated.

## What stays in the host

The membership scoping in `mwanachama-wakala-api`'s `insight_scope.go` and
`insight_mcp_scope.go` is deliberately not here. It reads the caller's
session and the agency registry, neither of which this module knows about.
That is the same boundary `mwanachama-backend-catalog` keeps: this module
has no auth model.
