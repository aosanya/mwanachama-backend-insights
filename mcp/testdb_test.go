package mcp_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

// newTestManager builds an InsightManager backed by a fresh in-memory
// sqlite database, migrated the same way a real deployment would.
func newTestManager(t *testing.T) mwanachamainsights.InsightManager {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	tables := mwanachamainsights.DefaultTableNames("mcp_test")
	if err := mwanachamainsights.Migrate(db, tables); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	mgr, err := mwanachamainsights.NewInsightManager(db, tables)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	return mgr
}
