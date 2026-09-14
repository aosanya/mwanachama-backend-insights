// insight.go — HTTP routes over CreateInsight/GetInsight/ListInsights. See
// doc.go for scope.
package routes

import (
	"errors"
	"net/http"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-insights/models"
)

// insightStatusFor maps this package's Insight error sentinels to a status
// code.
func insightStatusFor(err error) int {
	switch {
	case errors.Is(err, mwanachamainsights.ErrInsightNotFound):
		return http.StatusNotFound
	case errors.Is(err, mwanachamainsights.ErrInvalidInsight), errors.Is(err, mwanachamainsights.ErrInvalidInsightNote):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeInsightErr(w http.ResponseWriter, err error) {
	code := insightStatusFor(err)
	if code == http.StatusInternalServerError {
		writeErr(w, code, "internal error")
		return
	}
	writeErr(w, code, err.Error())
}

// CreateInsight handles POST /insights — decode, create, encode.
func CreateInsight(im mwanachamainsights.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in models.Insight
		if err := readJSON(r, &in); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		out, err := im.CreateInsight(r.Context(), in)
		if err != nil {
			writeInsightErr(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, out)
	}
}

// GetInsight handles GET /insights/{insightID}.
func GetInsight(im mwanachamainsights.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := im.GetInsight(r.Context(), r.PathValue("insightID"))
		if err != nil {
			writeInsightErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, out)
	}
}

// ListInsights handles GET /insights — optionally filtered by the query
// param agency_id.
func ListInsights(im mwanachamainsights.InsightManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := mwanachamainsights.InsightFilter{
			Repo:     r.URL.Query().Get("repo"),
			AgencyID: r.URL.Query().Get("agency_id"),
		}
		out, err := im.ListInsights(r.Context(), filter)
		if err != nil {
			writeInsightErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, out)
	}
}
