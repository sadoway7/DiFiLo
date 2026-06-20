package app

import (
	"encoding/json"
	"net/http"
)

// handleHealth returns a small JSON status document used for liveness checks.
// It reports "ok" and the current published-page count so a deploy/monitor
// can confirm the database is loaded and reachable.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"pages":  s.DB.PageCount(),
	})
}
