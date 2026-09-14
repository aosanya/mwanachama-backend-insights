package gormstore

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aosanya/mwanachama-backend-insights/models"
)

// InsightRow is the GORM row for a [models.Insight]. No UpdatedAt/Deleted
// columns — an Insight is append-only, never updated or deleted (see
// models.Insight's doc).
type InsightRow struct {
	ID          string `gorm:"primaryKey"`
	Repo        string `gorm:"index"`
	Summary     string
	Challenges  string
	Suggestions string
	CreatedAt   string
}

func (r *InsightRow) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

// InsightToRow converts a domain Insight to its row shape.
func InsightToRow(in models.Insight) InsightRow {
	return InsightRow{
		ID:          in.ID,
		Repo:        in.Repo,
		Summary:     in.Summary,
		Challenges:  in.Challenges,
		Suggestions: in.Suggestions,
		CreatedAt:   in.CreatedAt,
	}
}

// InsightFromRow converts a row back to the domain Insight.
func InsightFromRow(r InsightRow) models.Insight {
	return models.Insight{
		ID:          r.ID,
		Repo:        r.Repo,
		Summary:     r.Summary,
		Challenges:  r.Challenges,
		Suggestions: r.Suggestions,
		CreatedAt:   r.CreatedAt,
	}
}
