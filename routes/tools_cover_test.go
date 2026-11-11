package routes_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/aosanya/mwanachama-backend-shared/dispatch"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-insights/routes"
)

func toolNamed(t *testing.T, tools []routes.Tool, name string) routes.Tool {
	t.Helper()
	for _, tl := range tools {
		if tl.Name == name {
			return tl
		}
	}
	t.Fatalf("no tool named %q", name)
	return routes.Tool{}
}

func TestFields_MarksStoreOwnedColumnsReadOnlyAndSummaryRequired(t *testing.T) {
	fields, err := routes.Fields()
	if err != nil {
		t.Fatalf("Fields: %v", err)
	}
	for _, key := range []string{"Insight.created_at", "InsightNote.created_at", "Insight.id"} {
		doc, ok := fields[key]
		if !ok {
			t.Fatalf("missing field doc %s", key)
		}
		if !doc.ReadOnly {
			t.Fatalf("%s must be read only", key)
		}
	}
	if !fields["Insight.summary"].Required {
		t.Fatalf("Insight.summary must be required")
	}
	if fields["Insight.summary"].Description == "" {
		t.Fatalf("Insight.summary must be described")
	}
}

func TestFields_DescribesTheListFilters(t *testing.T) {
	fields, err := routes.Fields()
	if err != nil {
		t.Fatalf("Fields: %v", err)
	}
	for _, key := range []string{"InsightFilter.repo", "InsightFilter.subject_id"} {
		if fields[key].Description == "" {
			t.Fatalf("filter %s has no description", key)
		}
	}
}

func TestBuildTools_PublishesTheFrozenNamesAndDerivesTheRest(t *testing.T) {
	tools, err := routes.BuildTools(newTestManager(t))
	if err != nil {
		t.Fatalf("BuildTools: %v", err)
	}
	var names []string
	for _, tl := range tools {
		names = append(names, tl.Name)
	}
	sort.Strings(names)
	got := strings.Join(names, ",")
	want := "insight_create,insight_list,insight_note_create,insight_note_list,insights_insight_read"
	if got != want {
		t.Fatalf("tool names = %s, want %s", got, want)
	}
}

func TestToolInsightCreate_ThenListAndNoteRoundTrip(t *testing.T) {
	tools := routes.Tools(newTestManager(t))
	ctx := context.Background()

	created, err := toolNamed(t, tools, "insight_create").Invoke(ctx, json.RawMessage(`{"summary":"Friction in deploys","subject_id":"s1"}`))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	in, ok := created.(mwanachamainsights.Insight)
	if !ok || in.ID == "" || in.Summary != "Friction in deploys" {
		t.Fatalf("created = %#v", created)
	}

	if _, err := toolNamed(t, tools, "insight_note_create").Invoke(ctx, json.RawMessage(`{"insightID":"`+in.ID+`","text":"still happening"}`)); err != nil {
		t.Fatalf("note create: %v", err)
	}
	listed, err := toolNamed(t, tools, "insight_list").Invoke(ctx, json.RawMessage(`{"subject_id":"s1"}`))
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	raw, _ := json.Marshal(listed)
	if !strings.Contains(string(raw), "Friction in deploys") {
		t.Fatalf("list lost the insight: %s", raw)
	}
}

func TestToolInsightCreate_RefusesAMissingSummary(t *testing.T) {
	tools := routes.Tools(newTestManager(t))

	_, err := toolNamed(t, tools, "insight_create").Invoke(context.Background(), json.RawMessage(`{"subject_id":"s1"}`))

	if err == nil {
		t.Fatalf("expected a refusal for an insight with no summary")
	}
}

func TestToolsFor_AuthorizerSeesTheDeclaredActionAndCanRefuse(t *testing.T) {
	var seen []string
	deny := errors.New("no")
	mount := routes.Mount{Authorize: func(_ context.Context, action string) error {
		seen = append(seen, action)
		return deny
	}}
	tools := routes.ToolsFor(newTestManager(t), mount)

	_, err := toolNamed(t, tools, "insight_create").Invoke(context.Background(), json.RawMessage(`{"summary":"x"}`))

	if err == nil {
		t.Fatalf("a refusing authorizer must stop the call")
	}
	if len(seen) != 1 || seen[0] != "insights.insight.create" {
		t.Fatalf("authorizer saw %v", seen)
	}
}

func TestRoutesFor_RefusedCallerGetsForbiddenAndWritesNothing(t *testing.T) {
	im := newTestManager(t)
	mount := routes.Mount{Authorize: func(context.Context, string) error { return dispatch.ErrForbidden }}
	mux := http.NewServeMux()
	for _, rt := range routes.RoutesFor(im, mount) {
		mux.HandleFunc(rt.Pattern(""), rt.Handler)
	}
	srv := newServerFor(t, mux)

	res, _ := do(t, http.MethodPost, srv+"/insights", `{"summary":"blocked"}`)

	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", res.StatusCode)
	}
	list, err := im.ListInsights(context.Background(), mwanachamainsights.InsightFilter{})
	if err != nil || len(list) != 0 {
		t.Fatalf("a forbidden call wrote %d insights (err %v)", len(list), err)
	}
}

func TestShape_ListsEveryRouteWithoutHandlersBound(t *testing.T) {
	shape := routes.Shape()
	if len(shape) != 5 {
		t.Fatalf("shape has %d routes, want 5", len(shape))
	}
	actions := map[string]bool{}
	for _, rt := range shape {
		actions[rt.Action] = true
	}
	for _, a := range []string{"insights.insight.create", "insights.insight.read", "insights.insight.list", "insights.note.create", "insights.note.list"} {
		if !actions[a] {
			t.Fatalf("shape lacks action %s", a)
		}
	}
}

func TestAnonymousActions_IsEmpty(t *testing.T) {
	if len(routes.AnonymousActions) != 0 {
		t.Fatalf("nothing may be anonymous: %v", routes.AnonymousActions)
	}
}

func newServerFor(t *testing.T, h http.Handler) string {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv.URL
}
