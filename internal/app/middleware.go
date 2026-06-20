package app

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// csrfCookieName is the name of the double-submit CSRF cookie. It is NOT
// HttpOnly so client-side JavaScript can read it and echo it back in the
// X-CSRF-Token header on state-changing fetch() calls.
const csrfCookieName = "difilo_csrf"

// RequireAuth wraps a handler so that only authenticated users may reach it.
// Unauthenticated requests are redirected to /login when they look like a
// browser navigation, or receive a 401 JSON error otherwise.
func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.currentUser(r) == nil {
			if wantsJSON(r) {
				jsonError(w, "login required", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login?next="+url.QueryEscape(r.URL.String()), http.StatusFound)
			return
		}
		next(w, r)
	}
}

// RequireRole wraps a handler so that only users whose role matches the given
// role may reach it. Any other user (including unauthenticated visitors)
// receives a 403.
func (s *Server) RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.currentUser(r)
		if user == nil || user.Role != role {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// CSRFMiddleware sets a double-submit CSRF cookie on every response and, for
// unsafe (state-changing) HTTP methods, verifies that the request echoes the
// cookie value back via the X-CSRF-Token header or _csrf form field.
//
// The bootstrap authentication endpoints (/login, /register) and JSON-bodied
// API requests are exempt: the former must work before any cookie exists and
// the latter require a CORS preflight, which makes them naturally
// CSRF-resistant.
//
// This is a self-contained replacement for the nosurf package so the app has
// no external CSRF dependency; it can be swapped for nosurf with no caller
// changes.
func (s *Server) CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := ensureCSRFCookie(w, r)
		if isUnsafeMethod(r.Method) && !csrfExempt(r) {
			sent := r.Header.Get("X-CSRF-Token")
			if sent == "" {
				sent = r.FormValue("_csrf")
			}
			if sent == "" || subtle.ConstantTimeCompare([]byte(sent), []byte(cookie.Value)) != 1 {
				http.Error(w, "invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		next(w, r)
	}
}

// ensureCSRFCookie returns the request's existing CSRF cookie, or sets a fresh
// one on the response and returns it when none is present yet.
func ensureCSRFCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	if c, err := r.Cookie(csrfCookieName); err == nil && c.Value != "" {
		return c
	}
	token := randomCSRFToken()
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // readable by client JS so it can be echoed back
		SameSite: http.SameSiteLaxMode,
	})
	return &http.Cookie{Name: csrfCookieName, Value: token}
}

func isUnsafeMethod(m string) bool {
	switch m {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

// csrfExempt reports whether a state-changing request is exempt from CSRF
// verification: the login/register bootstrap endpoints and JSON API requests.
func csrfExempt(r *http.Request) bool {
	switch r.URL.Path {
	case "/login", "/register":
		return true
	}
	if ct := r.Header.Get("Content-Type"); strings.HasPrefix(ct, "application/json") {
		return true
	}
	return false
}

// wantsJSON reports whether the client expects a JSON response (used to decide
// between a JSON error and an HTML redirect on auth failure).
func wantsJSON(r *http.Request) bool {
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		return true
	}
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

// randomCSRFToken returns a fresh 32-byte random hex token, falling back to a
// time-based value if the system RNG is unavailable.
func randomCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(b)
}
