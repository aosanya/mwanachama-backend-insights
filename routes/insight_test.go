package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

func TestCreateInsight(t *testing.T) {
	srv, _ := newTestServer(t)

	body := `{"subject_id":"subject-1","context_id":"context-1","source":"user","summary":"Operator noted a gap.","tags":"risk, decision"}`
	res, raw := do(t, http.MethodPost, srv.URL+"/insights", body)

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
	var out mwanachamainsights.Insight
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID == "" || out.SubjectID != "subject-1" || out.Source != "user" || out.Tags != "risk, decision" {
		t.Fatalf("unexpected insight: %+v", out)
	}
}

func TestCreateInsight_InvalidShape(t *testing.T) {
	srv, _ := newTestServer(t)

	res, raw := do(t, http.MethodPost, srv.URL+"/insights", `{"subject_id":"subject-1"}`)

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
}

func TestGetInsight_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)

	res, raw := do(t, http.MethodGet, srv.URL+"/insights/nope", "")

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
}

func TestListInsights_FilterBySubjectID(t *testing.T) {
	srv, im := newTestServer(t)
	seed := func(subjectID, summary string) {
		t.Helper()
		if _, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{SubjectID: subjectID, Summary: summary}); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	seed("subject-1", "first")
	seed("subject-2", "second")

	res, raw := do(t, http.MethodGet, srv.URL+"/insights?subject_id=subject-1", "")

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
	var out []mwanachamainsights.Insight
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out) != 1 || out[0].Summary != "first" {
		t.Fatalf("expected 1 insight for subject-1, got %+v", out)
	}
}

func TestUpsertInsight_TheAddressOutranksTheBody(t *testing.T) {
	srv, im := newTestServer(t)
	insight, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{Summary: "Original."})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	res, raw := do(t, http.MethodPost, srv.URL+"/insights/"+insight.ID+"/notes",
		`{"insight_id":"some-other-insight","text":"attached by the path, not the body"}`)

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
	var out mwanachamainsights.InsightNote
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.InsightID != insight.ID {
		t.Fatalf("the body's insight_id won: got %q, want %q", out.InsightID, insight.ID)
	}
}
