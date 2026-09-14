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

func TestCreateInsight(t *testing.T) {
	im := newTestManager(t)
	handler := routes.CreateInsight(im)

	body := `{"agency_id":"agency-1","draft_id":"draft-1","source":"user","summary":"Operator noted a gap.","tags":"risk, decision"}`
	req := httptest.NewRequest(http.MethodPost, "/insights", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out mwanachamainsights.Insight
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID == "" || out.AgencyID != "agency-1" || out.Source != "user" || out.Tags != "risk, decision" {
		t.Fatalf("unexpected insight: %+v", out)
	}
}

func TestCreateInsight_InvalidShape(t *testing.T) {
	im := newTestManager(t)
	handler := routes.CreateInsight(im)

	req := httptest.NewRequest(http.MethodPost, "/insights", strings.NewReader(`{"agency_id":"agency-1"}`))
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestGetInsight_NotFound(t *testing.T) {
	im := newTestManager(t)
	handler := routes.GetInsight(im)

	req := withPathValue(httptest.NewRequest(http.MethodGet, "/insights/nope", nil), "insightID", "nope")
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestListInsights_FilterByAgencyID(t *testing.T) {
	im := newTestManager(t)
	seed := func(agencyID, summary string) {
		t.Helper()
		if _, err := im.CreateInsight(context.Background(), mwanachamainsights.Insight{AgencyID: agencyID, Summary: summary}); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	seed("agency-1", "first")
	seed("agency-2", "second")
	handler := routes.ListInsights(im)

	req := httptest.NewRequest(http.MethodGet, "/insights?agency_id=agency-1", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out []mwanachamainsights.Insight
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if len(out) != 1 || out[0].Summary != "first" {
		t.Fatalf("expected 1 insight for agency-1, got %+v", out)
	}
}
