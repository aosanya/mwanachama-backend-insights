package models

// Source values for Insight.Source and InsightNote.Source.
const (
	SourceAuto = "auto"
	SourceUser = "user"
)

// Insight is a record of a working session's or a wakala Agency chat's
// friction, suggestions, or observations, surfaced as a short write-up.
// Append-only: there is no update or delete, the same as
// mwanachama-backend-assetmanager's Movement ledger. Notes (see
// [InsightNote]) are how it grows over time without duplicating Summary.
type Insight struct {
	// ID is the unique identifier for this insight. Set by the backend on
	// creation; callers should leave it empty in create requests.
	ID string `json:"id"`

	// Repo is which repo or project the session concerned, if any (e.g.
	// "mwanachama-backend-agency"). Empty when the session wasn't scoped
	// to one repo.
	Repo string `json:"repo,omitempty"`

	// AgencyID is which mwanachama-wakala Agency this insight concerns, if
	// any. Empty for a dev-session insight recorded via MCP.
	AgencyID string `json:"agency_id,omitempty"`

	// DraftID is which Draft of AgencyID was open when this insight was
	// captured, if any.
	DraftID string `json:"draft_id,omitempty"`

	// Source is SourceAuto or SourceUser — whether this insight was
	// recorded by an AI agent or explicitly entered by a person. Empty on
	// create defaults to SourceAuto.
	Source string `json:"source,omitempty"`

	// Summary is a short account of what happened or was observed.
	Summary string `json:"summary"`

	// Tags is optional, comma-separated free text (e.g. "risk, decision"),
	// not a JSON-encoded array — same plain-string convention as
	// mwanachama-backend-assetmanager's Asset.AttributesJSON.
	Tags string `json:"tags,omitempty"`

	// Challenges is free text describing repetitive tasks or friction
	// points observed during the session.
	Challenges string `json:"challenges,omitempty"`

	// Suggestions is free text describing shortcuts, tooling, or
	// enhancements worth building based on this session.
	Suggestions string `json:"suggestions,omitempty"`

	// CreatedAt is the RFC 3339 timestamp when this insight was recorded.
	CreatedAt string `json:"created_at"`
}
