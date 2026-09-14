package models

// Insight is a record of one working session's friction and suggestions —
// repetitive tasks, challenges hit, or shortcuts/tooling worth building —
// surfaced as a short write-up at (or near) the end of the session.
// Append-only: there is no update or delete, the same as
// mwanachama-backend-assetmanager's Movement ledger.
type Insight struct {
	// ID is the unique identifier for this insight. Set by the backend on
	// creation; callers should leave it empty in create requests.
	ID string `json:"id"`

	// Repo is which repo or project the session concerned, if any (e.g.
	// "mwanachama-backend-agency"). Empty when the session wasn't scoped
	// to one repo.
	Repo string `json:"repo,omitempty"`

	// Summary is a short account of what happened in the session.
	Summary string `json:"summary"`

	// Challenges is free text describing repetitive tasks or friction
	// points observed during the session.
	Challenges string `json:"challenges,omitempty"`

	// Suggestions is free text describing shortcuts, tooling, or
	// enhancements worth building based on this session.
	Suggestions string `json:"suggestions,omitempty"`

	// CreatedAt is the RFC 3339 timestamp when this insight was recorded.
	CreatedAt string `json:"created_at"`
}
