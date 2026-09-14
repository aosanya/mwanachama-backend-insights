package gormstore

import "gorm.io/gorm"

// TableNames configures which physical tables an InsightManager reads and
// writes.
type TableNames struct {
	Insights     string
	InsightNotes string
}

// DefaultTableNames builds the conventional table names for one mounted
// instance of this package, e.g. DefaultTableNames("app") yields
// app_insights/app_insight_notes.
func DefaultTableNames(instance string) TableNames {
	return TableNames{
		Insights:     instance + "_insights",
		InsightNotes: instance + "_insight_notes",
	}
}

// Migrate creates or updates the tables t names, via GORM's AutoMigrate
// scoped to those table names. Callers run this once at startup (or in
// test setup) before constructing an InsightManager with the same db and
// t.
func Migrate(db *gorm.DB, t TableNames) error {
	if err := db.Table(t.Insights).AutoMigrate(&InsightRow{}); err != nil {
		return err
	}
	return db.Table(t.InsightNotes).AutoMigrate(&InsightNoteRow{})
}
