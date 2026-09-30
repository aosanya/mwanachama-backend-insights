//go:build integration

package mwanachamainsights_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-shared/spec"
)

func pgSpec(t *testing.T) (*gorm.DB, *spec.Spec) {
	t.Helper()
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		t.Skip("POSTGRES_URL is unset")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	s, err := mwanachamainsights.LoadSpec("insights.wakala.json")
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	s.Instance = fmt.Sprintf("pg%06x", rand.Intn(1<<24))

	t.Cleanup(func() {
		for _, o := range s.Objects {
			db.Exec("drop table if exists " + s.TableFor(o))
		}
		db.Exec("drop table if exists " + s.NameRegistryTable())
	})
	return db, s
}

func TestPostgresProvisionAndRoundTrip(t *testing.T) {
	db, s := pgSpec(t)
	if err := mwanachamainsights.Provision(db, s); err != nil {
		t.Fatalf("Provision: %v", err)
	}
	mgr, err := mwanachamainsights.NewInsightManager(db, s)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	ctx := context.Background()

	in, err := mgr.CreateInsight(ctx, mwanachamainsights.Insight{
		Repo: "some-repo", SubjectID: "subject-1", Summary: "recorded on postgres",
	})
	if err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}
	if in.Source != "auto" {
		t.Errorf("Source = %q, want the blueprint's default", in.Source)
	}

	got, err := mgr.GetInsight(ctx, in.ID)
	if err != nil {
		t.Fatalf("GetInsight: %v", err)
	}
	if got.Summary != "recorded on postgres" || got.SubjectID != "subject-1" {
		t.Fatalf("round trip lost a value: %+v", got)
	}
}

func TestPostgresEveryDeclaredNameFitsTheIdentifierLimit(t *testing.T) {
	_, s := pgSpec(t)
	for _, o := range s.Objects {
		if n := len(s.TableFor(o)); n > spec.MaxIdentifier {
			t.Errorf("%s needs %d bytes as %q, and postgres truncates at %d",
				s.RawNameFor(o), n, s.TableFor(o), spec.MaxIdentifier)
		}
		for _, idx := range o.Indexes {
			if n := len(s.IndexFor(o, idx)); n > spec.MaxIdentifier {
				t.Errorf("index %s.%s needs %d bytes as %q, and postgres truncates at %d",
					o.Name, idx.Name, n, s.IndexFor(o, idx), spec.MaxIdentifier)
			}
		}
	}
}

func TestPostgresAdoptsTheLegacyTableSet(t *testing.T) {
	db, s := pgSpec(t)
	instance := s.Instance
	stmts := []string{
		`create table ` + instance + `_insights (
			id text primary key, repo text, agency_id text, draft_id text,
			source text, summary text, tags text, challenges text,
			suggestions text, created_at text)`,
		`create table ` + instance + `_insight_notes (
			id text primary key, insight_id text, source text, text text,
			created_at text)`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("seed legacy schema: %v", err)
		}
	}
	t.Cleanup(func() {
		db.Exec("drop table if exists " + instance + "_insights")
		db.Exec("drop table if exists " + instance + "_insight_notes")
	})
	err := db.Exec(`insert into `+instance+`_insights
		(id, repo, agency_id, draft_id, source, summary, created_at)
		values ($1, $2, $3, $4, $5, $6, $7)`,
		"row-1", "some-repo", "agy-1", "draft-1", "user",
		"recorded before the conversion", "2026-09-01T00:00:00Z").Error
	if err != nil {
		t.Fatalf("seed legacy row: %v", err)
	}

	if err := mwanachamainsights.Provision(db, s); err != nil {
		t.Fatalf("Provision: %v", err)
	}

	mgr, err := mwanachamainsights.NewInsightManager(db, s)
	if err != nil {
		t.Fatalf("NewInsightManager: %v", err)
	}
	got, err := mgr.GetInsight(context.Background(), "row-1")
	if err != nil {
		t.Fatalf("the pre-conversion row did not survive adoption: %v", err)
	}
	if got.SubjectID != "agy-1" || got.ContextID != "draft-1" {
		t.Fatalf("the renamed columns lost their values: %+v", got)
	}
}
