package app

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"difilo/internal/textutil"
)

// handleExploreCards returns random card data as JSON for the homepage
// refresh button.
func (s *Server) handleExploreCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	articlePicks := s.DB.RandomCardsBySection("article", 2)
	generalPicks := s.DB.RandomCardPages(6)
	picks := append(articlePicks, generalPicks...)
	type card struct {
		Title, Route, Section, Thumb, Desc string
	}
	out := make([]card, 0, len(picks))
	for _, p := range picks {
		route := p.Route
		if route == "" {
			route = "/"
		}
		desc := p.MetaDescription
		if desc == "" {
			desc = textutil.Excerpt(p.BodyText, p.Title, 150)
		}
		out = append(out, card{
			Title:   p.Title,
			Route:   route,
			Section: textutil.PrettySection(p.Section),
			Thumb:   p.Thumb,
			Desc:    desc,
		})
	}
	_ = json.NewEncoder(w).Encode(out)
}

// handleSearch renders the search results page. Uses FTS5 full-text search
// via DB.SearchPages.
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	limit := s.Config.SearchLimit
	if limit <= 0 {
		limit = 100
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	var b strings.Builder
	b.WriteString(`<div class="df-searchpage">`)
	b.WriteString(`<form class="df-bigsearch" action="/search" method="get">`)
	b.WriteString(fmt.Sprintf(`<input type="text" name="q" value="%s" placeholder="Search all pages…" autocomplete="off" autofocus>`,
		html.EscapeString(q)))
	b.WriteString(`<button type="submit">Search</button>`)
	b.WriteString(`</form>`)

	if q == "" {
		b.WriteString(`<p class="df-muted">Type something to search `)
		b.WriteString(fmt.Sprintf("%d pages.", s.DB.PageCount()))
		b.WriteString(`</p>`)
	} else {
		hits := s.DB.SearchPages(q, limit)
		b.WriteString(fmt.Sprintf(`<p class="df-muted">%d result(s) for "%s"</p>`,
			len(hits), html.EscapeString(q)))
		b.WriteString(`<ol class="df-results">`)
		for _, h := range hits {
			route := h.Page.Route
			if route == "" {
				route = "/"
			}
			b.WriteString(`<li>`)
			if h.Page.Thumb != "" {
				b.WriteString(fmt.Sprintf(`<img class="df-thumb" loading="lazy" src="%s" alt="">`,
					html.EscapeString(h.Page.Thumb)))
			} else {
				b.WriteString(`<div class="df-thumb ph"></div>`)
			}
			b.WriteString(`<div class="df-rbody">`)
			b.WriteString(fmt.Sprintf(`<a class="df-title" href="%s">%s</a>`,
				html.EscapeString(route), html.EscapeString(textutil.OrDefault(h.Page.Title, "(untitled)"))))
			b.WriteString(fmt.Sprintf(`<span class="df-badge">%s</span>`, html.EscapeString(textutil.PrettySection(h.Page.Section))))
			if h.Snippet != "" {
				b.WriteString(fmt.Sprintf(`<div class="df-snip">%s</div>`, h.Snippet))
			}
			b.WriteString(`</div></li>`)
		}
		b.WriteString(`</ol>`)
	}

	b.WriteString(`</div>`)
	s.renderShell(w, r, "Search: "+q, "df-searchpage", b.String())
}

// handleSuggest returns up to 8 page-title matches as JSON for search
// autocomplete. Uses DB.SuggestTitles.
func (s *Server) handleSuggest(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if q == "" {
		_, _ = w.Write([]byte("[]"))
		return
	}
	type item struct {
		Title, Route, Section string
	}
	suggestions := s.DB.SuggestTitles(q, 8)
	out := make([]item, 0, len(suggestions))
	for _, sug := range suggestions {
		route := sug.Route
		if route == "" {
			route = "/"
		}
		out = append(out, item{sug.Title, route, sug.Section})
	}
	_ = json.NewEncoder(w).Encode(out)
}
