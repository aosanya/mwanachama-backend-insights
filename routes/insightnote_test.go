package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-insights/routes"
)

func TestCreateInsightNote(t *testing.T) {
	im := newTestManager(t)
	insight, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{Summary: "Original observation."})
	if err != nil {
		t.Fatalf("seed CreateInsight: %v", err)
	}
	handler := routes.CreateInsightNote(im)

	body := `{"source":"user","text":"Following up after another discussion."}`
	req := withPathValue(httptest.NewRequest(http.MethodPost, "/insights/"+insight.ID+"/notes", strings.NewReader(body)), "insightID", insight.ID)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out mwanachamainsights.InsightNote
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID == "" || out.InsightID != insight.ID || out.Text == "" {
		t.Fatalf("unexpected note: %+v", out)
	}
}

func TestCreateInsightNote_InsightNotFound(t *testing.T) {
	im := newTestManager(t)
	handler := routes.CreateInsightNote(im)

	req := withPathValue(httptest.NewRequest(http.MethodPost, "/insights/nope/notes", strings.NewReader(`{"text":"hi"}`)), "insightID", "nope")
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestListInsightNotes(t *testing.T) {
	im := newTestManager(t)
	insight, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{Summary: "Original observation."})
	if err != nil {
		t.Fatalf("seed CreateInsight: %v", err)
	}
	if _, err := im.CreateInsightNote(context.Background(), mwanachamainsights.InsightNote{InsightID: insight.ID, Text: "note one"}); err != nil {
		t.Fatalf("seed CreateInsightNote: %v", err)
	}
	handler := routes.ListInsightNotes(im)

	req := withPathValue(httptest.NewRequest(http.MethodGet, "/insights/"+insight.ID+"/notes", nil), "insightID", insight.ID)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out []mwanachamainsights.InsightNote
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if len(out) != 1 || out[0].Text != "note one" {
		t.Fatalf("expected 1 note, got %+v", out)
	}
}
