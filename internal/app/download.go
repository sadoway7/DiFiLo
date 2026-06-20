package app

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"

	"difilo/internal/content"
	"difilo/internal/ui"
)

// handleDownload renders a page as a self-contained .html download from the
// database: the app stylesheet is inlined, local images/media are embedded as
// base64 data URIs (via inlineAssets), and the page's comments are rendered
// statically. The download is logged against the logged-in user.
func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login?next="+url.QueryEscape(r.URL.String()), http.StatusFound)
		return
	}
	route := r.URL.Query().Get("page")
	if route == "" {
		route = "/"
	}

	// Look up page in DB, falling back to alias resolution.
	page, err := s.DB.GetPageByRoute(route)
	if err != nil || page == nil {
		if canonical, ok := s.DB.ResolveRoute(route, s.aliases); ok {
			route = canonical
			page, err = s.DB.GetPageByRoute(route)
		}
	}
	if err != nil || page == nil {
		http.NotFound(w, r)
		return
	}

	title := page.Title
	if title == "" {
		title = route
	}

	// Render the page content from markdown.
	contentHTML := content.RenderMarkdown(page.BodyMD)

	// Build a standalone HTML document.
	standalone := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8">`+
		`<meta name="viewport" content="width=device-width,initial-scale=1">`+
		`<title>%s</title></head><body>`,
		html.EscapeString(title))
	standalone += ui.DifiCSS()
	standalone += `<div class="df-wiki"><div class="df-wiki-main">`
	standalone += fmt.Sprintf(`<h1>%s</h1>`, html.EscapeString(title))
	standalone += `<div class="df-wiki-content">` + contentHTML + `</div>`
	standalone += `</div></div>`
	// Static comments.
	standalone += s.staticCommentsHTML(route)
	standalone += `</body></html>`

	doc := []byte(standalone)
	// Embed local assets so the single file is self-contained.
	doc = s.inlineAssets(doc)

	_ = s.DB.LogDownload(user.ID, user.Username, route, title)

	fname := sanitizeFilename(title) + ".html"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fname+`"`)
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(doc)
}

// staticCommentsHTML renders a page's comments as plain (non-interactive)
// HTML, so a downloaded file shows the discussion without needing the comment
// API.
func (s *Server) staticCommentsHTML(route string) string {
	comments, _ := s.DB.GetCommentsByRoute(route)
	var b strings.Builder
	b.WriteString(`<div id="df-comments"><h2>Comments</h2>`)
	if len(comments) == 0 {
		b.WriteString(`<p class="df-muted">No comments yet.</p>`)
	} else {
		b.WriteString(`<ul class="df-comment-list">`)
		for _, c := range comments {
			author := c.Username
			if author == "" {
				author = "(unknown)"
			}
			b.WriteString(fmt.Sprintf(`<li class="df-comment"><div class="df-comment-head">`+
				`<span class="df-comment-author">%s</span>`+
				`<span class="df-comment-role">%s</span></div>`+
				`<div class="df-comment-body">%s</div></li>`,
				html.EscapeString(author), html.EscapeString(c.Role),
				strings.ReplaceAll(html.EscapeString(c.Body), "\n", "<br>")))
		}
		b.WriteString(`</ul>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}
