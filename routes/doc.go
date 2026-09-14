// Package routes is mwanachama-backend-insights' own HTTP surface: decode
// a request, call one InsightManager method, encode the response. Mirrors
// mwanachama-backend-assetmanager/routes's shape and conventions on
// purpose: Route/Route.Pattern, one XRoutes function per aggregate
// concatenated by Routes, writeJSON/writeErr/readJSON, and decode-call-
// encode handlers with no caller-identity gate of their own — a route
// built from this package still needs an auth/capability check wrapped
// around it by whatever mounts it (mwanachama-wakala-api's requireCaller,
// today).
//
// This surface exists for a human caller (mwanachama-wakala-studio's
// agency-chat screen, letting an operator add a user-generated insight and
// notes on it) — the existing mcp/ package remains the AI-agent surface,
// unchanged by this package's addition.
package routes
