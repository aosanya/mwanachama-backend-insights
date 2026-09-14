// Package mwanachamainsights records insights — a short summary of a
// working session or a wakala Agency chat, plus any repetitive
// tasks/friction ("challenges") and shortcut/tooling ideas
// ("suggestions") it surfaced — and InsightNotes, append-only follow-ups
// attached to an existing Insight so a later remark on the same topic
// extends it instead of duplicating its Summary. It exposes
// [InsightManager], the single interface for creating and reading both.
//
// Layout:
//   - models/    — the Insight/InsightNote domain types; callers use
//     models.Insight/models.InsightNote directly, no re-export in this
//     package
//   - gormstore/ — GORM row structs, row<->domain conversion, migration
//   - doc.go (this file), tables.go — table-name/migrate wrappers
//   - insight.go, insightnote.go — InsightManager interface, insightManager
//     struct and its methods
//   - errors.go  — sentinel errors
//   - mcp/       — Model Context Protocol tools exposing InsightManager,
//     mirroring mwanachama-backend-agency/mcp's shape
//   - routes/    — plain REST HTTP surface, mirroring
//     mwanachama-backend-assetmanager/routes's shape, for a human caller
//     (e.g. mwanachama-wakala-studio) rather than an AI agent
package mwanachamainsights
