package routes

import (
	"net/http"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
)

// Route is one address this package answers, relative to wherever the
// mounting process prefixes it. Path uses net/http's ServeMux pattern
// syntax ("{insightID}" etc.), so the mounting process only ever needs
// prefix+rt.Path, never its own copy of the path text.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Pattern returns the http.ServeMux registration pattern for this route
// once mounted under prefix — r.Method+" "+prefix+r.Path, net/http's own
// "METHOD /path" syntax (Go 1.22+ mux patterns).
func (r Route) Pattern(prefix string) string {
	return r.Method + " " + prefix + r.Path
}

// InsightRoutes is CreateInsight/GetInsight/ListInsights, addressed under
// "/insights".
func InsightRoutes(im mwanachamainsights.InsightManager) []Route {
	base := "/insights"
	return []Route{
		{Method: http.MethodPost, Path: base, Handler: CreateInsight(im)},
		{Method: http.MethodGet, Path: base, Handler: ListInsights(im)},
		{Method: http.MethodGet, Path: base + "/{insightID}", Handler: GetInsight(im)},
	}
}

// InsightNoteRoutes is CreateInsightNote/ListInsightNotes, addressed under
// "/insights/{insightID}/notes".
func InsightNoteRoutes(im mwanachamainsights.InsightManager) []Route {
	base := "/insights/{insightID}/notes"
	return []Route{
		{Method: http.MethodPost, Path: base, Handler: CreateInsightNote(im)},
		{Method: http.MethodGet, Path: base, Handler: ListInsightNotes(im)},
	}
}

// Routes is every address this package answers today: InsightRoutes and
// InsightNoteRoutes concatenated.
func Routes(im mwanachamainsights.InsightManager) []Route {
	out := InsightRoutes(im)
	out = append(out, InsightNoteRoutes(im)...)
	return out
}
