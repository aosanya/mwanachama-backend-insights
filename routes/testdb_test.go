package routes_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-insights/routes"
)

func newTestManager(t *testing.T) mwanachamainsights.InsightManager {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	s, err := mwanachamainsights.LoadSpec("../insights.wakala.json")
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	s.Instance = "routestest"
	if err := mwanachamainsights.Provision(db, s); err != nil {
		t.Fatalf("Provision: %v", err)
	}
	mgr, err := mwanachamainsights.NewInsightManager(db, s)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	return mgr
}

func newTestServer(t *testing.T) (*httptest.Server, mwanachamainsights.InsightManager) {
	t.Helper()
	im := newTestManager(t)
	mux := http.NewServeMux()
	for _, rt := range routes.Routes(im) {
		mux.HandleFunc(rt.Pattern(""), rt.Handler)
	}
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, im
}

func do(t *testing.T, method, url, body string) (*http.Response, string) {
	t.Helper()
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return res, string(raw)
}
