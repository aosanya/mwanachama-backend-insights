package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

func TestCreateInsightNote(t *testing.T) {
	srv, im := newTestServer(t)
	insight, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{Summary: "Original observation."})
	if err != nil {
		t.Fatalf("seed CreateInsight: %v", err)
	}

	body := `{"source":"user","text":"Following up after another discussion."}`
	res, raw := do(t, http.MethodPost, srv.URL+"/insights/"+insight.ID+"/notes", body)

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
	var out mwanachamainsights.InsightNote
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID == "" || out.InsightID != insight.ID || out.Text == "" {
		t.Fatalf("unexpected note: %+v", out)
	}
}

func TestCreateInsightNote_InsightNotFound(t *testing.T) {
	srv, _ := newTestServer(t)

	res, raw := do(t, http.MethodPost, srv.URL+"/insights/nope/notes", `{"text":"hi"}`)

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
}

func TestListInsightNotes(t *testing.T) {
	srv, im := newTestServer(t)
	insight, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{Summary: "Original observation."})
	if err != nil {
		t.Fatalf("seed CreateInsight: %v", err)
	}
	if _, err := im.CreateInsightNote(context.Background(), mwanachamainsights.InsightNote{InsightID: insight.ID, Text: "note one"}); err != nil {
		t.Fatalf("seed CreateInsightNote: %v", err)
	}

	res, raw := do(t, http.MethodGet, srv.URL+"/insights/"+insight.ID+"/notes", "")

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}
	var out []mwanachamainsights.InsightNote
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out) != 1 || out[0].Text != "note one" {
		t.Fatalf("expected 1 note, got %+v", out)
	}
}
