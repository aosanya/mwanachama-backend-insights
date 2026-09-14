package mwanachamainsights_test

import (
	"context"
	"errors"
	"testing"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

func TestCreateInsightNote_RequiresInsightIDAndText(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	if _, err := m.CreateInsightNote(ctx, mwanachamainsights.InsightNote{Text: "x"}); !errors.Is(err, mwanachamainsights.ErrInvalidInsightNote) {
		t.Fatalf("expected ErrInvalidInsightNote for missing InsightID, got %v", err)
	}
}

func TestCreateInsightNote_InsightMustExist(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	if _, err := m.CreateInsightNote(ctx, mwanachamainsights.InsightNote{InsightID: "does-not-exist", Text: "x"}); !errors.Is(err, mwanachamainsights.ErrInsightNotFound) {
		t.Fatalf("expected ErrInsightNotFound, got %v", err)
	}
}

func TestCreateInsightNote_ThenList(t *testing.T) {
	ctx := context.Background()
	m := newTestManager(t)

	insight, err := m.CreateInsight(ctx, mwanachamainsights.Insight{Summary: "Original observation."})
	if err != nil {
		t.Fatalf("CreateInsight: %v", err)
	}

	if _, err := m.CreateInsightNote(ctx, mwanachamainsights.InsightNote{InsightID: insight.ID, Text: "first note"}); err != nil {
		t.Fatalf("CreateInsightNote: %v", err)
	}
	if _, err := m.CreateInsightNote(ctx, mwanachamainsights.InsightNote{InsightID: insight.ID, Source: "user", Text: "second note"}); err != nil {
		t.Fatalf("CreateInsightNote: %v", err)
	}

	notes, err := m.ListInsightNotes(ctx, insight.ID)
	if err != nil {
		t.Fatalf("ListInsightNotes: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %+v", notes)
	}
	if notes[0].Text != "first note" || notes[1].Text != "second note" {
		t.Fatalf("expected oldest-first order, got %+v", notes)
	}
	if notes[0].Source != "auto" || notes[1].Source != "user" {
		t.Fatalf("expected Source defaulting/passthrough, got %+v", notes)
	}
}
