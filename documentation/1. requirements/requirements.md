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

## Open questions

None. The one deferred item — actually wiring a Stop hook or slash command
so a session calls `insight_create` automatically — is explicitly a
follow-up once the tool exists and can be tested manually first, not an
open design question about this repo's own shape.
