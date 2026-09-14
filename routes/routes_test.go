package routes_test

import (
	"testing"

	"github.com/aosanya/mwanachama-backend-insights/routes"
)

func patterns(rts []routes.Route, prefix string) []string {
	out := make([]string, len(rts))
	for i, rt := range rts {
		out[i] = rt.Pattern(prefix)
	}
	return out
}

func assertPatterns(t *testing.T, got []routes.Route, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d routes, want %d: %v", len(got), len(want), patterns(got, ""))
	}
	for i, p := range patterns(got, "") {
		if p != want[i] {
			t.Fatalf("route %d: got %q, want %q", i, p, want[i])
		}
	}
}

func TestInsightRoutes(t *testing.T) {
	im := newTestManager(t)
	rts := routes.InsightRoutes(im)
	assertPatterns(t, rts, []string{
		"POST /insights",
		"GET /insights",
		"GET /insights/{insightID}",
	})
}

func TestInsightNoteRoutes(t *testing.T) {
	im := newTestManager(t)
	rts := routes.InsightNoteRoutes(im)
	assertPatterns(t, rts, []string{
		"POST /insights/{insightID}/notes",
		"GET /insights/{insightID}/notes",
	})
}

func TestRoutes_ConcatenatesBoth(t *testing.T) {
	im := newTestManager(t)
	all := routes.Routes(im)

	want := 3 + 2 // InsightRoutes + InsightNoteRoutes
	if len(all) != want {
		t.Fatalf("got %d routes, want %d: %v", len(all), want, patterns(all, ""))
	}
}

func TestRoute_PatternWithPrefix(t *testing.T) {
	im := newTestManager(t)
	rts := routes.InsightRoutes(im)

	if got := rts[0].Pattern("/v1/insights-svc"); got != "POST /v1/insights-svc/insights" {
		t.Fatalf("got %q", got)
	}
}
