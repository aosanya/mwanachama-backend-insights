# mwanachama-backend-insights — documentation

## Layout

Four folders, in SDLC order, and everything lives under one of them.

| Folder | What's inside |
|--------|---------------|
| [1. requirements/](1.%20requirements/) | Problem, vision and scope decisions for recording session insights. |
| [2. design/](2.%20design/) | The Insight schema and how it's exposed via MCP. |
| [3. implementation/](3.%20implementation/) | The work: `todo.md` (open board), `todo_done.md` (completed rows + board context). |
| [4. qa/](4.%20qa/) | Test coverage and results. |

## What this repo is

A standalone Go library recording session insights — a short summary of a
working session plus repetitive tasks/friction ("challenges") and
shortcut/tooling ideas ("suggestions") it surfaced. Modelled on
[mwanachama-backend-assetmanager](../mwanachama-backend-assetmanager)'s
GORM-direct shape. No HTTP/gRPC layer of its own — exposed via an `mcp/`
package that [mwanachama-wakala-api](../mwanachama-wakala-api) mounts on
its existing `/mcp` endpoint, since api-gateway's own `/mcp` is
deliberately read-only and this domain needs a write tool
(`insight_create`).
