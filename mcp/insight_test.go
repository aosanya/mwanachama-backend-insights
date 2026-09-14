package mcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gosdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	insightsmcp "github.com/aosanya/mwanachama-backend-insights/mcp"
)

func ctxT(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func newTestServer(im mwanachamainsights.InsightManager) *gosdkmcp.Server {
	server := gosdkmcp.NewServer(&gosdkmcp.Implementation{Name: "test", Version: "0.0.0"}, nil)
	insightsmcp.RegisterTools(server, im)
	return server
}

func mcpConnect(t *testing.T, base string) *gosdkmcp.ClientSession {
	t.Helper()
	client := gosdkmcp.NewClient(&gosdkmcp.Implementation{Name: "test-client", Version: "0.0.0"}, nil)
	transport := &gosdkmcp.StreamableClientTransport{
		Endpoint:             base + "/mcp",
		DisableStandaloneSSE: true,
	}
	cs, err := client.Connect(ctxT(t), transport, nil)
	if err != nil {
		t.Fatalf("mcp connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func newTestHTTPServer(t *testing.T, im mwanachamainsights.InsightManager) string {
	t.Helper()
	server := newTestServer(im)
	mux := http.NewServeMux()
	mux.Handle("/mcp", gosdkmcp.NewStreamableHTTPHandler(func(*http.Request) *gosdkmcp.Server { return server }, nil))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv.URL
}

func TestRegisterTools_ListsExpectedTools(t *testing.T) {
	url := newTestHTTPServer(t, newTestManager(t))
	cs := mcpConnect(t, url)

	tools, err := cs.ListTools(ctxT(t), nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	names := map[string]bool{}
	for _, tl := range tools.Tools {
		names[tl.Name] = true
	}
	for _, w := range []string{"insight_create", "insight_list"} {
		if !names[w] {
			t.Errorf("tools/list missing %q", w)
		}
	}
}

func TestInsightCreateAndList(t *testing.T) {
	url := newTestHTTPServer(t, newTestManager(t))
	cs := mcpConnect(t, url)

	result, err := cs.CallTool(ctxT(t), &gosdkmcp.CallToolParams{
		Name: "insight_create",
		Arguments: map[string]any{
			"repo":        "mwanachama-backend-insights",
			"summary":     "Scaffolded a new domain repo end to end.",
			"challenges":  "Had to read four other repos to find the right pattern to mirror.",
			"suggestions": "A repo-scaffolding generator would save this every time.",
		},
	})
	if err != nil {
		t.Fatalf("CallTool insight_create: %v", err)
	}
	if result.IsError {
		t.Fatalf("insight_create: IsError = true (content: %+v)", result.Content)
	}

	result, err = cs.CallTool(ctxT(t), &gosdkmcp.CallToolParams{
		Name:      "insight_list",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool insight_list: %v", err)
	}
	got, ok := result.StructuredContent.(map[string]any)
	items, _ := got["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("insight_list structured content = %+v", result.StructuredContent)
	}
}

func TestInsightCreate_RequiresSummary(t *testing.T) {
	url := newTestHTTPServer(t, newTestManager(t))
	cs := mcpConnect(t, url)

	result, err := cs.CallTool(ctxT(t), &gosdkmcp.CallToolParams{
		Name:      "insight_create",
		Arguments: map[string]any{"repo": "some-repo"},
	})
	if err != nil {
		t.Fatalf("CallTool insight_create: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected IsError = true for a missing summary")
	}
}

func TestInsightList_FiltersByRepo(t *testing.T) {
	im := newTestManager(t)
	url := newTestHTTPServer(t, im)
	cs := mcpConnect(t, url)

	for _, args := range []map[string]any{
		{"repo": "repo-a", "summary": "first"},
		{"repo": "repo-b", "summary": "second"},
	} {
		if _, err := cs.CallTool(ctxT(t), &gosdkmcp.CallToolParams{Name: "insight_create", Arguments: args}); err != nil {
			t.Fatalf("CallTool insight_create: %v", err)
		}
	}

	result, err := cs.CallTool(ctxT(t), &gosdkmcp.CallToolParams{
		Name:      "insight_list",
		Arguments: map[string]any{"repo": "repo-a"},
	})
	if err != nil {
		t.Fatalf("CallTool insight_list: %v", err)
	}
	got, _ := result.StructuredContent.(map[string]any)
	items, _ := got["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 insight for repo-a, got %+v", got)
	}
}
