package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-insights/models"
)

// registerInsightTools registers insight_create and insight_list.
func registerInsightTools(server *mcp.Server, im mwanachamainsights.InsightManager) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insight_create",
		Title:       "Record Insight",
		Description: "Record a summary of a working session's challenges (repetitive tasks, friction points) and suggestions (shortcuts, tooling, or enhancements worth building).",
	}, mcpCreateInsight(im))
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insight_list",
		Title:       "List Insights",
		Description: "List recorded session insights, optionally filtered by repo.",
	}, mcpListInsights(im))
}

type mcpCreateInsightParams struct {
	Repo        string `json:"repo,omitempty" jsonschema:"which repo or project this session concerned, if any"`
	Summary     string `json:"summary" jsonschema:"a short summary of what happened in this session"`
	Challenges  string `json:"challenges,omitempty" jsonschema:"repetitive tasks or friction points observed in this session"`
	Suggestions string `json:"suggestions,omitempty" jsonschema:"shortcuts, tooling, or enhancements worth building based on this session"`
}

func (p mcpCreateInsightParams) toModel() models.Insight {
	return models.Insight{Repo: p.Repo, Summary: p.Summary, Challenges: p.Challenges, Suggestions: p.Suggestions}
}

func mcpCreateInsight(im mwanachamainsights.InsightManager) func(context.Context, *mcp.CallToolRequest, mcpCreateInsightParams) (*mcp.CallToolResult, models.Insight, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, p mcpCreateInsightParams) (*mcp.CallToolResult, models.Insight, error) {
		out, err := im.CreateInsight(ctx, p.toModel())
		return summary("Recorded insight (%s)", out.ID), out, err
	}
}

type mcpListInsightsParams struct {
	Repo string `json:"repo,omitempty" jsonschema:"restrict to insights recorded for this repo; empty lists every insight"`
}

func mcpListInsights(im mwanachamainsights.InsightManager) func(context.Context, *mcp.CallToolRequest, mcpListInsightsParams) (*mcp.CallToolResult, ListResult[models.Insight], error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, p mcpListInsightsParams) (*mcp.CallToolResult, ListResult[models.Insight], error) {
		out, err := im.ListInsights(ctx, mwanachamainsights.InsightFilter{Repo: p.Repo})
		return summary("Found %d insight(s)", len(out)), ListResult[models.Insight]{Items: out}, err
	}
}
