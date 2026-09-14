// Package mwanachamainsights records session insights — a short summary of
// a working session plus any repetitive tasks/friction ("challenges") and
// shortcut/tooling ideas ("suggestions") it surfaced. It exposes
// [InsightManager], the single interface for creating and reading Insight
// records.
//
// Layout:
//   - models/    — the Insight domain type; callers use models.Insight
//     directly, no re-export in this package
//   - gormstore/ — GORM row struct, row<->domain conversion, migration
//   - doc.go (this file), tables.go — table-name/migrate wrappers
//   - insight.go — InsightManager interface, insightManager struct
//   - errors.go  — sentinel errors
//   - mcp/       — Model Context Protocol tools exposing InsightManager,
//     mirroring mwanachama-backend-agency/mcp's shape
package mwanachamainsights
