package mwanachamainsights_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-shared/spec"
)

func testSpec(t *testing.T, instance string) *spec.Spec {
	t.Helper()
	s, err := mwanachamainsights.LoadSpec("insights.wakala.json")
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	s.Instance = instance
	return s
}

func newTestManager(t *testing.T) mwanachamainsights.InsightManager {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	s := testSpec(t, "test")
	if err := mwanachamainsights.Provision(db, s); err != nil {
		t.Fatalf("Provision: %v", err)
	}

	mgr, err := mwanachamainsights.NewInsightManager(db, s)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	return mgr
}

func TestNewInsightManager_NilDB(t *testing.T) {
	if _, err := mwanachamainsights.NewInsightManager(nil, testSpec(t, "test")); err == nil {
		t.Fatal("expected error for nil db")
	}
}
