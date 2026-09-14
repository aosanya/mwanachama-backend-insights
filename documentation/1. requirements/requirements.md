# Requirements — mwanachama-backend-insights

## Problem

Working across the mwanachama repo family in Claude Code, a session hits
repetitive tasks and friction that never get written down anywhere — they
just get repeated next session. The idea: at (or near) the end of a
session, have Claude post a short summary — what happened, what was
repetitive/annoying, what shortcut or tooling would have helped — into a
database, so a later pass can mine it for actual product/tooling work.

## Scope decisions (session of 2026-09-13/14)

| # | Question | Decision |
| - | -------- | -------- |
| 1 | Separate monitoring backend vs. reuse existing infra | No new backend for token-usage/monitoring — Claude Code already writes local session transcripts, and this repo only needs to record the *qualitative* summary, not raw usage metrics. |
| 2 | External model (Hugging Face) for the summarization step | Rejected. The summary is written by whichever Claude session is ending, in-conversation — no separate inference step, no separate hosting (HF Spaces/Inference/self-host) needed. |
| 3 | How the write actually happens | An MCP tool call (`insight_create`), not a REST call from an external routine — this repo's `mcp/` package. |
| 4 | Repo name | `insights` (bare-word style, matching `actor`/`assetmanager`/`orgsettings`) — general enough to cover challenges/shortcuts/suggestions, not narrowed to "chat". |
| 5 | Repo shape vs. bolt onto an existing domain | New standalone repo, mirroring `mwanachama-backend-assetmanager`'s shape (`models/`+`gormstore/`+root manager) — insights don't belong to any existing domain's purpose, and cramming them in would tie an unrelated concern to that repo's release cycle. |
| 6 | Where the MCP tool is hosted | `mwanachama-wakala-api`'s existing `/mcp` endpoint, not `mwanachama-backend-api-gateway`'s. The gateway's `/mcp` is deliberately read-only (every tool there wraps a read method); `insight_create` is a write, so it doesn't fit that endpoint's own convention. wakala-api already has precedent for write-capable MCP tools (`agency_*`) and is already running/registered in `.mcp.json`. |
| 7 | How much code lives in wakala-api vs. this repo | Most of it stays here, mirroring `mwanachama-backend-agency`'s AG13 decision (moving its MCP tools out of wakala-api and into the library itself, to stop drift between the model and the tool code). wakala-api only gets: the go.mod dependency, one `InsightManager` constructed alongside its existing DB setup, one field on `Deps`, one `RegisterTools` call, and a migration mirror. |
| 8 | Insight shape | `Repo` (optional), `Summary` (required), `Challenges` (free text), `Suggestions` (free text), `CreatedAt`. Plain text, not JSON-encoded arrays — an LLM caller writes prose into an MCP tool argument, not a hand-encoded array string, mirroring `mwanachama-backend-assetmanager`'s deliberate choice to keep `Asset.AttributesJSON` a plain string rather than promote it. |
| 9 | CRUD surface | Create + Get + List only, no Update/Delete — an Insight is an append-only historical record, same as `mwanachama-backend-assetmanager`'s Movement ledger. |
| 10 | Multi-instance / `ManagerResolver` indirection | Not needed. Agency's `mcp/` package takes a `ManagerResolver` because `mwanachama-wakala-api` mounts many Agencies behind one server; Insights has exactly one fixed instance for the whole wakala-api deployment, so its `mcp/` tools close over a concrete `InsightManager` directly. |
| 11 | REST `routes/` surface | Not built. Nothing asked for one — MCP is the only surface this repo needs today. Can be added later the same way `mwanachama-backend-assetmanager` added `routes/` after the fact, without a redesign. |

## Scope decisions (session of 2026-09-14, user-posted insights)

`mwanachama-wakala-studio`'s agency-chat screen needed a way for a human
operator to discuss with the agent and then explicitly capture an insight
about the agency being designed — distinct from any *automatic* insight
the agent itself might record. This required extending the repo's
original Claude-Code-session-shaped model to also serve that product use
case.

| # | Question | Decision |
| - | -------- | -------- |
| 12 | How to tell a user-captured insight from an agent-captured one | New `Source` field (`SourceAuto`/`SourceUser` in `models`), empty on create defaults to `SourceAuto` — the existing `insight_create` MCP tool needed zero code changes to keep landing as `"auto"`. |
| 13 | How to scope an insight to a wakala Agency/Draft | New optional `AgencyID` (indexed, mirrors `Repo`'s existing index) and `DraftID` fields — both empty for a dev-session insight, both populated for a wakala-chat insight. |
| 14 | How to capture insight content beyond Summary | A lightweight `Tags` field, comma-separated free text — same plain-string convention as `mwanachama-backend-assetmanager`'s `Asset.AttributesJSON` (decision #8's precedent, extended). No fixed category enum — nothing today needs filtering by a closed set of categories. |
| 15 | How to avoid a user's follow-up remark duplicating an insight already on record | New `InsightNote` type — an append-only follow-up attached to an existing Insight (`CreateInsightNote`/`ListInsightNotes`, same `InsightManager` interface, not a separate manager). Considered a self-referencing `RelatedInsightID` on Insight instead; rejected because every follow-up would still be a full new Insight row repeating a Summary, even for a one-line addition — a Note never touches the parent's Summary/Tags. |
| 16 | REST surface | Built — `routes/` package (superseding decision #11's "not built"), mirroring `mwanachama-backend-assetmanager/routes`'s shape: `InsightRoutes`/`InsightNoteRoutes`/`Routes`, mounted by `mwanachama-wakala-api` at `/insights` and `/insights/{insightID}/notes`, `requireCaller`-gated there. `mcp/` is unchanged — this surface is for a human caller, not an AI agent. |

## Open questions

None. Two deliberate follow-ups, not open design questions about this
repo's own shape:

- Wiring a Stop hook or slash command so a Claude Code session calls
  `insight_create` automatically — unchanged from the original scope.
- Wiring the wakala agent side of "automatic insights" (the agent itself
  calling `insight_create`/a future `insight_note_create` mid-conversation)
  — deferred until the human-facing path above is proven; no MCP tool
  changes were made in this pass.
