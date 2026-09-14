// Package mcp exposes mwanachama-backend-insights' InsightManager as Model
// Context Protocol tools, mirroring mwanachama-backend-agency/mcp's shape
// (the org's reference implementation for this — see that package's doc
// comment and mwanachama-backend-actor/CLAUDE.md's "MCP tools" section).
// Unlike agency's mcp package, tools here close over a concrete
// InsightManager directly rather than a per-call ManagerResolver: agency
// needs that indirection because a host like mwanachama-wakala-api mounts
// many Agencies behind one server, while an Insight deployment has exactly
// one fixed instance.
package mcp

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

// RegisterTools registers every tool this package exposes onto server,
// backed by im.
func RegisterTools(server *mcp.Server, im mwanachamainsights.InsightManager) {
	registerInsightTools(server, im)
}

// ListResult wraps a List tool's slice: MCP's spec defines
// CallToolResult.structuredContent as a JSON object, and the go-sdk
// otherwise serializes a bare slice return type directly into it — a shape
// a schema-strict client rejects even though the go-sdk's own test client
// tolerates it. See mwanachama-backend-agency/mcp's identical type.
type ListResult[T any] struct {
	Items []T `json:"items" jsonschema:"the matching items"`
}

// summary builds a *mcp.CallToolResult carrying a one-line human-readable
// sentence as its Content, in place of the go-sdk's default behavior of
// JSON-dumping the whole return value there. Only ever worth calling on a
// tool's success path — see mwanachama-backend-agency/mcp's identical
// helper for why.
func summary(format string, args ...any) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(format, args...)}}}
}
