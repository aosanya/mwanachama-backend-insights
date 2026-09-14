package gormstore

import "gorm.io/gorm"

// TableNames configures which physical table an InsightManager reads and
// writes.
type TableNames struct {
	Insights string
}

// DefaultTableNames builds the conventional table name for one mounted
// instance of this package, e.g. DefaultTableNames("app") yields
// app_insights.
func DefaultTableNames(instance string) TableNames {
	return TableNames{
		Insights: instance + "_insights",
	}
}

// Migrate creates or updates the table t names, via GORM's AutoMigrate
// scoped to that table name. Callers run this once at startup (or in test
// setup) before constructing an InsightManager with the same db and t.
func Migrate(db *gorm.DB, t TableNames) error {
	return db.Table(t.Insights).AutoMigrate(&InsightRow{})
}
