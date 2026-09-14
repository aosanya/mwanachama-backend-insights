// insightnote.go — HTTP routes over CreateInsightNote/ListInsightNotes.
// See doc.go for scope.
package routes

import (
	"net/http"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

// createInsightNoteBody is the wire shape for POST
// /insights/{insightID}/notes — InsightID comes from the path, not the
// body.
type createInsightNoteBody struct {
	Source string `json:"source,omitempty"`
	Text   string `json:"text"`
}

// CreateInsightNote handles POST /insights/{insightID}/notes — decode,
// create, encode.
func CreateInsightNote(im mwanachamainsights.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body createInsightNoteBody
		if err := readJSON(r, &body); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		out, err := im.CreateInsightNote(r.Context(), mwanachamainsights.InsightNote{
			InsightID: r.PathValue("insightID"),
			Source:    body.Source,
			Text:      body.Text,
		})
		if err != nil {
			writeInsightErr(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, out)
	}
}

// ListInsightNotes handles GET /insights/{insightID}/notes.
func ListInsightNotes(im mwanachamainsights.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := im.ListInsightNotes(r.Context(), r.PathValue("insightID"))
		if err != nil {
			writeInsightErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, out)
	}
}
