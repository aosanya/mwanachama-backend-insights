package mwanachamainsights_test

import (
	"context"
	"errors"
	"testing"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

func TestCreateInsight_RequiresSummary(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	_, err := m.CreateInsight(ctx, mwanachamainsights.Insight{Repo: "some-repo"})
	if !errors.Is(err, mwanachamainsights.ErrInvalidInsight) {
		t.Fatalf("expected ErrInvalidInsight, got %v", err)
	}
}

func TestCreateInsight_ThenGet(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	created, err := m.CreateInsight(ctx, mwanachamainsights.Insight{
		Repo:        "mwanachama-backend-insights",
		Summary:     "Scaffolded a new domain repo end to end.",
		Challenges:  "Had to read four other repos to find the right pattern to mirror.",
		Suggestions: "A repo-scaffolding generator would save this every time.",
	})
	if err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected a generated ID")
	}
	if created.CreatedAt == "" {
		t.Fatal("expected a generated CreatedAt")
	}

	got, err := m.GetInsight(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetInsight: %v", err)
	}
	if got != created {
		t.Fatalf("expected GetInsight to round-trip, got %+v want %+v", got, created)
	}
}

func TestGetInsight_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	if _, err := m.GetInsight(ctx, "does-not-exist"); !errors.Is(err, mwanachamainsights.ErrInsightNotFound) {
		t.Fatalf("expected ErrInsightNotFound, got %v", err)
	}
}

func TestCreateInsight_SourceDefaultsToAuto(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	created, err := m.CreateInsight(ctx, mwanachamainsights.Insight{Summary: "captured by the agent"})
	if err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}
	if created.Source != "auto" {
		t.Fatalf("expected Source to default to %q, got %q", "auto", created.Source)
	}
}

func TestCreateInsight_RejectsUnknownSource(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	_, err := m.CreateInsight(ctx, mwanachamainsights.Insight{Summary: "x", Source: "bogus"})
	if !errors.Is(err, mwanachamainsights.ErrInvalidInsight) {
		t.Fatalf("expected ErrInvalidInsight, got %v", err)
	}
}

func TestListInsights_FiltersByAgencyID(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	if _, err := m.CreateInsight(ctx, mwanachamainsights.Insight{AgencyID: "agency-a", Summary: "first"}); err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}
	if _, err := m.CreateInsight(ctx, mwanachamainsights.Insight{AgencyID: "agency-b", Summary: "second"}); err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}

	scoped, err := m.ListInsights(ctx, mwanachamainsights.InsightFilter{AgencyID: "agency-a"})
	if err != nil {
		t.Fatalf("ListInsights scoped: %v", err)
	}
	if len(scoped) != 1 || scoped[0].Summary != "first" {
		t.Fatalf("expected 1 insight for agency-a, got %+v", scoped)
	}
}

func TestListInsights_FiltersByRepo(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	if _, err := m.CreateInsight(ctx, mwanachamainsights.Insight{Repo: "repo-a", Summary: "first"}); err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}
	if _, err := m.CreateInsight(ctx, mwanachamainsights.Insight{Repo: "repo-b", Summary: "second"}); err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}

	all, err := m.ListInsights(ctx, mwanachamainsights.InsightFilter{})
	if err != nil {
		t.Fatalf("ListInsights: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 insights, got %d", len(all))
	}

	scoped, err := m.ListInsights(ctx, mwanachamainsights.InsightFilter{Repo: "repo-a"})
	if err != nil {
		t.Fatalf("ListInsights scoped: %v", err)
	}
	if len(scoped) != 1 || scoped[0].Summary != "first" {
		t.Fatalf("expected 1 insight for repo-a, got %+v", scoped)
	}
}
