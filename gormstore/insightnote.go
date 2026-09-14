package gormstore

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aosanya/mwanachama-backend-insights/models"
)

// InsightNoteRow is the GORM row for a [models.InsightNote]. Append-only,
// same as InsightRow — no UpdatedAt/Deleted columns.
type InsightNoteRow struct {
	ID        string `gorm:"primaryKey"`
	InsightID string `gorm:"index"`
	Source    string
	Text      string
	CreatedAt string
}

func (r *InsightNoteRow) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

// InsightNoteToRow converts a domain InsightNote to its row shape.
func InsightNoteToRow(in models.InsightNote) InsightNoteRow {
	return InsightNoteRow{
		ID:        in.ID,
		InsightID: in.InsightID,
		Source:    in.Source,
		Text:      in.Text,
		CreatedAt: in.CreatedAt,
	}
}

// InsightNoteFromRow converts a row back to the domain InsightNote.
func InsightNoteFromRow(r InsightNoteRow) models.InsightNote {
	return models.InsightNote{
		ID:        r.ID,
		InsightID: r.InsightID,
		Source:    r.Source,
		Text:      r.Text,
		CreatedAt: r.CreatedAt,
	}
}
