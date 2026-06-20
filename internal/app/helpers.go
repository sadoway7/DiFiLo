package app

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"

	"difilo/internal/auth"
	"difilo/internal/db"
	"difilo/internal/ui"
)

// currentUser reads the signed session cookie and returns the authenticated
// *db.User (or nil). It uses auth.ParseSessionToken to validate the cookie
// and s.DB.GetUserByID to load the row.
func (s *Server) currentUser(r *http.Request) *db.User {
	c, err := r.Cookie(auth.SessionCookieName)
	if err != nil {
		return nil
	}
	userID := auth.ParseSessionToken(c.Value)
	if userID == 0 {
		return nil
	}
	user, err := s.DB.GetUserByID(userID)
	if err != nil || user == nil {
		return nil
	}
	return user
}

// userToViewer converts a *db.User into the lightweight *ui.Viewer DTO that
// the ui package functions accept. A nil user yields a nil viewer, which the
// ui functions treat as "logged out".
func userToViewer(u *db.User) *ui.Viewer {
	if u == nil {
		return nil
	}
	return &ui.Viewer{
		LoggedIn: true,
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
}

// canDeleteComment adapts the app-layer db types to the auth.UserRole and
// auth.CommentOwner interfaces and delegates to auth.CanDeleteComment. Admins
// and managers may delete any comment; general users may delete only their own.
func canDeleteComment(user *db.User, comment *db.Comment) bool {
	if user == nil || comment == nil {
		return false
	}
	return auth.CanDeleteComment(
		&userRoleAdapter{role: user.Role, id: user.ID},
		&commentOwnerAdapter{userID: comment.UserID},
	)
}

// userRoleAdapter adapts *db.User to the auth.UserRole interface.
type userRoleAdapter struct {
	role string
	id   int64
}

func (a *userRoleAdapter) GetRole() string { return a.role }
func (a *userRoleAdapter) GetID() int64    { return a.id }

// commentOwnerAdapter adapts *db.Comment to the auth.CommentOwner interface.
type commentOwnerAdapter struct{ userID int64 }

func (a *commentOwnerAdapter) GetUserID() int64 { return a.userID }

// die prints a formatted error to stderr and exits the process with status 1.
// It is used only during startup, before the server begins serving.
func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "difilo: "+format+"\n", a...)
	os.Exit(1)
}

// escape is a shorthand for html.EscapeString, used by the inline HTML
// builders in the handler layer.
func escape(s string) string {
	return html.EscapeString(s)
}

// jsonError writes a JSON error response with the given HTTP status code.
func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
