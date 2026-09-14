// insight.go — CreateInsight/GetInsight/ListInsights for [insightManager].
package mwanachamainsights

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/aosanya/mwanachama-backend-insights/gormstore"
	"github.com/aosanya/mwanachama-backend-insights/models"
)

// Insight is an alias of [models.Insight], so a caller needs only this
// package's import, never models's directly — mirrors
// mwanachama-backend-assetmanager's identical convenience alias.
type Insight = models.Insight

// InsightFilter scopes a [InsightManager.ListInsights] query. Zero-value
// fields are ignored (no filtering applied for that field).
type InsightFilter struct {
	// Repo restricts results to insights recorded for this exact repo.
	Repo string
	// AgencyID restricts results to insights recorded for this exact
	// wakala Agency.
	AgencyID string
}

// InsightManager is the primary interface for recording and reading
// insights and their notes. This is a single-tenant package — one
// deployment serves one owner, so no method takes a tenant-scoping
// argument, mirroring mwanachama-backend-assetmanager.AssetManager's
// shape. Append-only: there is no Update or Delete, for either Insight or
// InsightNote. InsightNote's methods live on this same interface rather
// than a separate manager — it is subordinate to Insight the same way
// AssetManager carries Movement's methods directly.
//
// Implementations must be safe for concurrent use.
type InsightManager interface {
	CreateInsight(ctx context.Context, in models.Insight) (models.Insight, error)
	GetInsight(ctx context.Context, insightID string) (models.Insight, error)
	ListInsights(ctx context.Context, filter InsightFilter) ([]models.Insight, error)

	CreateInsightNote(ctx context.Context, in models.InsightNote) (models.InsightNote, error)
	ListInsightNotes(ctx context.Context, insightID string) ([]models.InsightNote, error)
}

// insightManager is the GORM-backed implementation of [InsightManager].
type insightManager struct {
	db     *gorm.DB
	tables TableNames
}

// NewInsightManager constructs an [InsightManager] backed by db, reading
// and writing the table named by t (see [DefaultTableNames]). Callers must
// run [Migrate] against the same db and t before use. Returns an error if
// db is nil.
func NewInsightManager(db *gorm.DB, t TableNames) (InsightManager, error) {
	if db == nil {
		return nil, fmt.Errorf("NewInsightManager: db must not be nil")
	}
	return &insightManager{db: db, tables: t}, nil
}

// CreateInsight validates and appends a new Insight record. Summary is the
// only required field. Source defaults to models.SourceAuto when empty,
// so the existing insight_create MCP tool (which never sets it) keeps
// working unchanged.
func (m *insightManager) CreateInsight(ctx context.Context, in models.Insight) (models.Insight, error) {
	if in.Summary == "" {
		return models.Insight{}, fmt.Errorf("%w: Summary is required", ErrInvalidInsight)
	}
	if in.Source == "" {
		in.Source = models.SourceAuto
	}
	if in.Source != models.SourceAuto && in.Source != models.SourceUser {
		return models.Insight{}, fmt.Errorf("%w: Source must be %q or %q", ErrInvalidInsight, models.SourceAuto, models.SourceUser)
	}
	in.ID = ""
	in.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	row := gormstore.InsightToRow(in)
	if err := m.db.WithContext(ctx).Table(m.tables.Insights).Create(&row).Error; err != nil {
		return models.Insight{}, fmt.Errorf("CreateInsight: %w", err)
	}
	return gormstore.InsightFromRow(row), nil
}

// GetInsight reads a single Insight row.
func (m *insightManager) GetInsight(ctx context.Context, insightID string) (models.Insight, error) {
	var row gormstore.InsightRow
	err := m.db.WithContext(ctx).Table(m.tables.Insights).Where("id = ?", insightID).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Insight{}, ErrInsightNotFound
		}
		return models.Insight{}, fmt.Errorf("GetInsight: %w", err)
	}
	return gormstore.InsightFromRow(row), nil
}

// ListInsights returns every Insight row matching filter, newest first.
func (m *insightManager) ListInsights(ctx context.Context, filter InsightFilter) ([]models.Insight, error) {
	q := m.db.WithContext(ctx).Table(m.tables.Insights)
	if filter.Repo != "" {
		q = q.Where("repo = ?", filter.Repo)
	}
	if filter.AgencyID != "" {
		q = q.Where("agency_id = ?", filter.AgencyID)
	}

	var rows []gormstore.InsightRow
	if err := q.Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ListInsights: %w", err)
	}
	out := make([]models.Insight, 0, len(rows))
	for _, r := range rows {
		out = append(out, gormstore.InsightFromRow(r))
	}
	return out, nil
}
