package app

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"difilo/internal/content"
	"difilo/internal/db"
	"difilo/internal/ui"
)

// handlePage renders a wiki-style page from the database. Resolution order:
//  1. Direct route lookup in pages table
//  2. Alias resolution (slug-based → canonical)
//  3. Fuzzy slug search (handles old numeric IDs, mismatched slugs)
//  4. If all fail: show a helpful "not found" with search links
func (s *Server) handlePage(w http.ResponseWriter, r *http.Request, p string) {
	// 1. Direct DB lookup.
	page, err := s.DB.GetPageByRoute(p)
	if err == nil && page != nil {
		s.renderWikiPageFromDB(w, r, page)
		return
	}

	// 2. Alias resolution for slug-based routes.
	if canonical, ok := s.DB.ResolveRoute(p, s.aliases); ok && canonical != p {
		http.Redirect(w, r, canonical, http.StatusFound)
		return
	}

	// 3. Fallback: search by slug segment.
	if found, ok := s.DB.SearchBySlug(p); ok && found.Route != "" && found.Route != p {
		http.Redirect(w, r, found.Route, http.StatusFound)
		return
	}

	// 4. Not found — show search results instead of a dead end.
	searchQ := p
	if parts := strings.Split(strings.TrimPrefix(p, "/"), "/"); len(parts) >= 2 {
		searchQ = strings.ReplaceAll(parts[len(parts)-1], "+", " ")
	}
	s.renderShell(w, r, "Not found", "", notFoundWithSearch(html.EscapeString(p), html.EscapeString(searchQ)))
}

// renderWikiPageFromDB renders a full wiki page from a ContentPage. It pulls
// the page's images and links, hands them to content.RenderWikiPage, and
// wraps the result in the page shell together with the per-page pin/bookmark
// toolbar and the interactive comment section.
func (s *Server) renderWikiPageFromDB(w http.ResponseWriter, r *http.Request, p *db.ContentPage) {
	viewer := userToViewer(s.currentUser(r))

	images := s.DB.GetPageImages(p.ID)
	links := s.DB.GetPageLinks(p.ID)
	wikiBody := content.RenderWikiPage(p, images, links)
	commentSection := ui.CommentsHTML(p.Route, viewer)

	inner := ui.PinButtonHTML(p.Route, viewer) + content.WikiCSS + wikiBody + commentSection

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprint(w, ui.ShellHTML(p.Title, "df-wikipage", ui.PanelHTML(p.Route, viewer)+inner))
}

// notFoundWithSearch renders a "not found" page with a search link.
func notFoundWithSearch(route, query string) string {
	return fmt.Sprintf(`<div class="df-external">
		<h1>Page not found</h1>
		<p>This page isn't in the database.</p>
		<p class="df-muted">Requested: <code>%s</code></p>
		<div style="margin:20px 0">
			<form class="df-bigsearch" action="/search" method="get" style="max-width:500px">
				<input type="text" name="q" value="%s" placeholder="Search for this topic…">
				<button type="submit">Search</button>
			</form>
		</div>
		<p><a href="/">← back to home</a></p>
	</div>`, route, query)
}
