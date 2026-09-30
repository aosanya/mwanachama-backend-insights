package routes_test

import (
	"sort"
	"testing"

	"github.com/aosanya/mwanachama-backend-insights/routes"
)

func TestRoutesAreTheDeclaredTable(t *testing.T) {
	im := newTestManager(t)

	got := make([]string, 0)
	for _, rt := range routes.Routes(im) {
		got = append(got, rt.Pattern(""))
	}
	sort.Strings(got)

	want := []string{
		"GET /insights",
		"GET /insights/{insightID}",
		"GET /insights/{insightID}/notes",
		"POST /insights",
		"POST /insights/{insightID}/notes",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d routes, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("route %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestShapeNeedsNoManager(t *testing.T) {
	if len(routes.Shape()) != 5 {
		t.Fatalf("got %d routes, want 5", len(routes.Shape()))
	}
	for _, rt := range routes.Shape() {
		if rt.Action == "" {
			t.Errorf("%s %s carries no action id", rt.Method, rt.Path)
		}
	}
}

func TestRoutePatternTakesAPrefix(t *testing.T) {
	im := newTestManager(t)
	for _, rt := range routes.Routes(im) {
		if rt.Method == "POST" && rt.Path == "/insights" {
			if got := rt.Pattern("/v1/insights-svc"); got != "POST /v1/insights-svc/insights" {
				t.Fatalf("got %q", got)
			}
			return
		}
	}
	t.Fatal("POST /insights is not in the table")
}

func TestEverySentinelTheSpecMapsIsSupplied(t *testing.T) {
	im := newTestManager(t)
	if _, err := routes.Build(im); err != nil {
		t.Fatalf("Build: %v", err)
	}
}

func TestNothingIsAnonymous(t *testing.T) {
	if len(routes.AnonymousActions) != 0 {
		t.Fatalf("AnonymousActions is %v — every insights door is the mounting host's to gate", routes.AnonymousActions)
	}
}
