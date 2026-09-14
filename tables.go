package mwanachamainsights

import (
	"gorm.io/gorm"

	"github.com/aosanya/mwanachama-backend-insights/gormstore"
)

// TableNames configures which physical table an [InsightManager] reads and
// writes. See [gormstore.TableNames].
type TableNames = gormstore.TableNames

// DefaultTableNames builds the conventional table name for one mounted
// instance of this package. See [gormstore.DefaultTableNames].
func DefaultTableNames(instance string) TableNames {
	return gormstore.DefaultTableNames(instance)
}

// Migrate creates or updates the table t names. See [gormstore.Migrate].
func Migrate(db *gorm.DB, t TableNames) error {
	return gormstore.Migrate(db, t)
}
