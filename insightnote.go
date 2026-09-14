// insightnote.go — CreateInsightNote/ListInsightNotes for [insightManager].
package mwanachamainsights

import (
	"context"
	"fmt"
	"time"

	"github.com/aosanya/mwanachama-backend-insights/gormstore"
	"github.com/aosanya/mwanachama-backend-insights/models"
)

// InsightNote is an alias of [models.InsightNote], same convenience alias
// as Insight.
type InsightNote = models.InsightNote

// CreateInsightNote validates and appends a new InsightNote, attaching it
// to an existing Insight. Text is required; InsightID must reference a
// real Insight. Source defaults to models.SourceAuto when empty.
func (m *insightManager) CreateInsightNote(ctx context.Context, in models.InsightNote) (models.InsightNote, error) {
	if in.InsightID == "" || in.Text == "" {
		return models.InsightNote{}, fmt.Errorf("%w: InsightID and Text are required", ErrInvalidInsightNote)
	}
	if in.Source == "" {
		in.Source = models.SourceAuto
	}
	if in.Source != models.SourceAuto && in.Source != models.SourceUser {
		return models.InsightNote{}, fmt.Errorf("%w: Source must be %q or %q", ErrInvalidInsightNote, models.SourceAuto, models.SourceUser)
	}
	if _, err := m.GetInsight(ctx, in.InsightID); err != nil {
		return models.InsightNote{}, err
	}

	in.ID = ""
	in.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	row := gormstore.InsightNoteToRow(in)
	if err := m.db.WithContext(ctx).Table(m.tables.InsightNotes).Create(&row).Error; err != nil {
		return models.InsightNote{}, fmt.Errorf("CreateInsightNote: %w", err)
	}
	return gormstore.InsightNoteFromRow(row), nil
}

// ListInsightNotes returns every InsightNote for insightID, oldest first
// (a comment-thread reads top-to-bottom, unlike ListInsights' newest-first
// feed).
func (m *insightManager) ListInsightNotes(ctx context.Context, insightID string) ([]models.InsightNote, error) {
	var rows []gormstore.InsightNoteRow
	q := m.db.WithContext(ctx).Table(m.tables.InsightNotes).Where("insight_id = ?", insightID)
	if err := q.Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ListInsightNotes: %w", err)
	}
	out := make([]models.InsightNote, 0, len(rows))
	for _, r := range rows {
		out = append(out, gormstore.InsightNoteFromRow(r))
	}
	return out, nil
}
