package mcp_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
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
	s.Instance = "mcptest"
	if err := mwanachamainsights.Provision(db, s); err != nil {
		t.Fatalf("Provision: %v", err)
	}
	mgr, err := mwanachamainsights.NewInsightManager(db, s)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	return mgr
}
