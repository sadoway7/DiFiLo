package app

import (
	"encoding/json"
	"net/http"

	"difilo/internal/db"
)

// handleAPIBookmark adds (POST) or removes (DELETE) a bookmark for the
// logged-in user.
func (s *Server) handleAPIBookmark(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		jsonError(w, "login required", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodDelete {
		var body struct {
			Route string `json:"route"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid request", http.StatusBadRequest)
			return
		}
		if body.Route == "" {
			body.Route = "/"
		}
		_ = s.DB.DeleteBookmark(user.ID, body.Route)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// POST = add bookmark.
	var body struct {
		Route string `json:"route"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.Route == "" {
		body.Route = "/"
	}
	if body.Title == "" {
		body.Title = body.Route
	}
	if err := s.DB.CreateBookmark(user.ID, body.Route, body.Title); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handleAPIBookmarks returns all bookmarks for the logged-in user as JSON.
func (s *Server) handleAPIBookmarks(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		jsonError(w, "login required", http.StatusUnauthorized)
		return
	}
	bookmarks, err := s.DB.GetBookmarksByUser(user.ID)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	if bookmarks == nil {
		bookmarks = []db.Bookmark{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(bookmarks)
}

// handleAPIBookmarkCheck reports whether the logged-in user has bookmarked a
// given route. Unauthenticated users always get {"bookmarked":false}.
func (s *Server) handleAPIBookmarkCheck(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		_ = json.NewEncoder(w).Encode(map[string]bool{"bookmarked": false})
		return
	}
	route := r.URL.Query().Get("route")
	if route == "" {
		route = "/"
	}
	pinned, _ := s.DB.IsBookmarked(user.ID, route)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"bookmarked": pinned})
}
