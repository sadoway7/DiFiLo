package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"difilo/internal/db"
	"difilo/internal/ui"
)

// Server is the main application state. All page content comes from the
// SQLite database (DB); static assets (images, media, vendor files) are
// served from on-disk directories derived from MirrorDir. The Config pointer
// holds tunables and may be shared with sub-packages if needed.
type Server struct {
	DB         *db.DB
	Config     *Config
	MirrorDir  string
	imageDir   string
	vendorDir  string
	mediaDir   string
	heroImages []string          // pool of random hero-background image paths
	aliases    map[string]string // slug-based alias → canonical route (from DB)
}

// New constructs a Server wired to the given database and mirror directory.
// If cfg is nil, DefaultConfig is used. The image/vendor/media directories
// are derived from mirrorDir exactly as the legacy binary did.
func New(cfg *Config, database *db.DB, mirrorDir string) *Server {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Server{
		Config:    cfg,
		DB:        database,
		MirrorDir: mirrorDir,
		imageDir:  filepath.Join(mirrorDir, "images"),
		vendorDir: filepath.Join(mirrorDir, "vendor"),
		mediaDir:  filepath.Join(mirrorDir, "media"),
	}
}

// renderShell wraps inner content in the full page shell (CSS + navbar +
// body) using the ui package. The nav panel always uses "/" as its active
// route hint for shell-rendered app pages; wiki pages build their own shell
// with the page route via ui.PanelHTML directly.
func (s *Server) renderShell(w http.ResponseWriter, r *http.Request, title, bodyClass, inner string) {
	viewer := userToViewer(s.currentUser(r))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, ui.ShellHTML(title, bodyClass, ui.PanelHTML("/", viewer)+inner))
}

// Handler returns the top-level http.Handler for the application: a single
// catch-all multiplexer wrapped in the CSRF middleware. The future cmd/difilo
// binary installs this under "/".
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.CSRFMiddleware(s.handle))
	return mux
}

// handle is the central request router. It matches the URL path against the
// application's routes and dispatches to the appropriate handler method.
func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	switch {
	case p == "/" || p == "":
		s.handleHome(w, r)
		return
	case p == "/health":
		s.handleHealth(w, r)
		return
	case p == "/search":
		s.handleSearch(w, r)
		return
	case p == "/suggest":
		s.handleSuggest(w, r)
		return
	case p == "/download":
		s.handleDownload(w, r)
		return
	case p == "/random":
		s.handleRandom(w, r)
		return
	case p == "/api/explore":
		s.handleExploreCards(w, r)
		return
	case p == "/api/download-check" && r.Method == http.MethodPost:
		s.handleAPIDownloadCheck(w, r)
		return
	case p == "/external":
		s.handleExternal(w, r)
		return
	// Auth routes
	case p == "/login":
		s.handleLogin(w, r)
		return
	case p == "/register":
		s.handleRegister(w, r)
		return
	case p == "/logout":
		s.handleLogout(w, r)
		return
	case p == "/admin":
		s.handleAdmin(w, r)
		return
	case p == "/admin/user-downloads":
		s.handleAdminUserDownloads(w, r)
		return
	// Comment API
	case p == "/api/comments":
		s.handleAPIComments(w, r)
		return
	case p == "/api/comment" && r.Method == http.MethodPost:
		s.handleAPIComment(w, r)
		return
	case strings.HasPrefix(p, "/api/comment/delete/"):
		s.handleAPICommentDelete(w, r)
		return
	case strings.HasPrefix(p, "/api/comment/edit/"):
		s.handleAPICommentEdit(w, r)
		return
	// Bookmark API
	case p == "/api/bookmark" && (r.Method == http.MethodPost || r.Method == http.MethodDelete):
		s.handleAPIBookmark(w, r)
		return
	case p == "/api/bookmarks":
		s.handleAPIBookmarks(w, r)
		return
	case p == "/api/bookmark/check":
		s.handleAPIBookmarkCheck(w, r)
		return
	// Admin API
	case p == "/api/admin/role" && r.Method == http.MethodPost:
		s.handleAPIAdminRole(w, r)
		return
	case p == "/api/admin/delete-user" && r.Method == http.MethodPost:
		s.handleAPIAdminDeleteUser(w, r)
		return
	case p == "/api/admin/settings" && r.Method == http.MethodPost:
		s.handleAPIAdminSettings(w, r)
		return
	// Static assets
	case strings.HasPrefix(p, "/images/"):
		s.serveStatic(w, r, s.imageDir, strings.TrimPrefix(p, "/images/"))
		return
	case strings.HasPrefix(p, "/vendor/"):
		s.serveStatic(w, r, s.vendorDir, strings.TrimPrefix(p, "/vendor/"))
		return
	case strings.HasPrefix(p, "/media/"):
		s.serveStatic(w, r, s.mediaDir, strings.TrimPrefix(p, "/media/"))
		return
	// Section list pages
	case strings.HasPrefix(p, "/list/"):
		s.handleList(w, r, strings.TrimPrefix(p, "/list/"))
		return
	case p == "/home":
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Site list pages (e.g. /material/list, /glossary/list) — route to the
	// generated full list with working A-Z navigation.
	if parts := strings.Split(strings.TrimPrefix(p, "/"), "/"); len(parts) >= 2 && parts[1] == "list" {
		s.handleList(w, r, parts[0])
		return
	}

	// Otherwise: treat as a content page route → render from DB.
	s.handlePage(w, r, p)
}
