package app

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"difilo/internal/db"
)

// handleAPIComments returns all comments for a route as JSON.
func (s *Server) handleAPIComments(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Query().Get("route")
	if route == "" {
		route = "/"
	}
	comments, err := s.DB.GetCommentsByRoute(route)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	if comments == nil {
		comments = []db.Comment{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comments)
}

// handleAPIComment creates a new comment (POST).
func (s *Server) handleAPIComment(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		jsonError(w, "you must be logged in to comment", http.StatusUnauthorized)
		return
	}
	var body struct {
		Route    string `json:"route"`
		Body     string `json:"body"`
		ParentID *int64 `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	body.Body = strings.TrimSpace(body.Body)
	if body.Body == "" {
		jsonError(w, "comment cannot be empty", http.StatusBadRequest)
		return
	}
	if len(body.Body) > s.commentMaxLen() {
		jsonError(w, "comment too long", http.StatusBadRequest)
		return
	}
	if body.Route == "" {
		body.Route = "/"
	}
	var parentID sql.NullInt64
	if body.ParentID != nil {
		parentID = sql.NullInt64{Int64: *body.ParentID, Valid: true}
	}
	id, err := s.DB.CreateComment(user.ID, body.Route, body.Body, parentID)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	comment, _ := s.DB.GetCommentByID(id)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comment)
}

// handleAPICommentDelete deletes a comment (POST). Admins and managers may
// delete any comment; general users may delete only their own.
func (s *Server) handleAPICommentDelete(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		jsonError(w, "login required", http.StatusUnauthorized)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/comment/delete/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "invalid comment id", http.StatusBadRequest)
		return
	}
	comment, err := s.DB.GetCommentByID(id)
	if err != nil || comment == nil {
		jsonError(w, "comment not found", http.StatusNotFound)
		return
	}
	if !canDeleteComment(user, comment) {
		jsonError(w, "permission denied", http.StatusForbidden)
		return
	}
	if err := s.DB.DeleteComment(id); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleAPICommentEdit edits a comment's body (POST). Only the comment owner
// may edit (not even admins/managers edit others' comments).
func (s *Server) handleAPICommentEdit(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		jsonError(w, "login required", http.StatusUnauthorized)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/comment/edit/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "invalid comment id", http.StatusBadRequest)
		return
	}
	var body struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	body.Body = strings.TrimSpace(body.Body)
	if body.Body == "" {
		jsonError(w, "comment cannot be empty", http.StatusBadRequest)
		return
	}
	if len(body.Body) > s.commentMaxLen() {
		jsonError(w, "comment too long", http.StatusBadRequest)
		return
	}
	comment, err := s.DB.GetCommentByID(id)
	if err != nil || comment == nil {
		jsonError(w, "comment not found", http.StatusNotFound)
		return
	}
	if user.ID != comment.UserID {
		jsonError(w, "you can only edit your own comments", http.StatusForbidden)
		return
	}
	if err := s.DB.UpdateComment(id, body.Body); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	comment.Body = body.Body
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comment)
}

// commentMaxLen returns the configured maximum comment length, falling back to
// 5000 when no config is present.
func (s *Server) commentMaxLen() int {
	if s.Config != nil && s.Config.CommentMaxLen > 0 {
		return s.Config.CommentMaxLen
	}
	return 5000
}
