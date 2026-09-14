package models

// InsightNote is an append-only follow-up attached to one Insight — how an
// insight grows over time (a person or agent adding color, confirmation,
// or a related observation) without spawning a near-duplicate Insight that
// repeats its Summary.
type InsightNote struct {
	// ID is the unique identifier for this note. Set by the backend on
	// creation; callers should leave it empty in create requests.
	ID string `json:"id"`

	// InsightID is the Insight this note is attached to. Required.
	InsightID string `json:"insight_id"`

	// Source is SourceAuto or SourceUser. Empty on create defaults to
	// SourceAuto.
	Source string `json:"source,omitempty"`

	// Text is the note's content. Required.
	Text string `json:"text"`

	// CreatedAt is the RFC 3339 timestamp when this note was recorded.
	CreatedAt string `json:"created_at"`
}
