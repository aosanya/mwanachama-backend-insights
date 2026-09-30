package mwanachamainsights_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-shared/spec"
)

func legacyTableSet(t *testing.T, db *gorm.DB, instance string) {
	t.Helper()
	stmts := []string{
		`create table ` + instance + `_insights (
			id text primary key, repo text, agency_id text, draft_id text,
			source text, summary text, tags text, challenges text,
			suggestions text, created_at text)`,
		`create table ` + instance + `_insight_notes (
			id text primary key, insight_id text, source text, text text,
			created_at text)`,
		`create index idx_` + instance + `_insights_repo on ` + instance + `_insights (repo)`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("seed legacy schema: %v", err)
		}
	}
	err := db.Exec(`insert into `+instance+`_insights
		(id, repo, agency_id, draft_id, source, summary, created_at)
		values (?, ?, ?, ?, ?, ?, ?)`,
		"row-1", "some-repo", "agy-1", "draft-1", "user", "recorded before the conversion",
		"2026-09-01T00:00:00Z").Error
	if err != nil {
		t.Fatalf("seed legacy row: %v", err)
	}
}

func TestProvisionAdoptsTheLegacyTableSet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	s := testSpec(t, "adopt")
	legacyTableSet(t, db, "adopt")

	if err := mwanachamainsights.Provision(db, s); err != nil {
		t.Fatalf("Provision: %v", err)
	}

	o, ok := s.ByRole(mwanachamainsights.RoleInsight)
	if !ok {
		t.Fatal("no object fills the insight role")
	}
	if db.Migrator().HasTable("adopt_insights") {
		t.Errorf("adopt_insights still exists — the rows were left where no read reaches them")
	}
	if !db.Migrator().HasTable(s.TableFor(o)) {
		t.Fatalf("%s was never created", s.RawNameFor(o))
	}

	mgr, err := mwanachamainsights.NewInsightManager(db, s)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	got, err := mgr.GetInsight(context.Background(), "row-1")
	if err != nil {
		t.Fatalf("the pre-conversion row did not survive adoption: %v", err)
	}
	if got.Summary != "recorded before the conversion" {
		t.Errorf("Summary = %q", got.Summary)
	}
	if got.SubjectID != "agy-1" {
		t.Errorf("SubjectID = %q, want %q — agency_id should have been renamed, not dropped", got.SubjectID, "agy-1")
	}
	if got.ContextID != "draft-1" {
		t.Errorf("ContextID = %q, want %q — draft_id should have been renamed, not dropped", got.ContextID, "draft-1")
	}
}

func TestProvisionIsIdempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	s := testSpec(t, "twice")

	for i := 0; i < 3; i++ {
		if err := mwanachamainsights.Provision(db, s); err != nil {
			t.Fatalf("Provision pass %d: %v", i+1, err)
		}
	}
	adopted, err := spec.Adopted(db, s)
	if err != nil {
		t.Fatalf("Adopted: %v", err)
	}
	if !adopted {
		t.Error("Adopted reports false after three provisioning passes")
	}
}
