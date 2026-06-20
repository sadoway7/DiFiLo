package app

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"difilo/internal/textutil"
)

// handleList shows all pages in a section as an alphabetical A-Z list with a
// client-side filter.
func (s *Server) handleList(w http.ResponseWriter, r *http.Request, section string) {
	docs := s.DB.ListPagesBySection(section)

	var b strings.Builder
	b.WriteString(`<div class="df-listpage">`)
	b.WriteString(`<div class="df-lp-head">`)
	b.WriteString(fmt.Sprintf(`<h1><span class="df-lp-title">%s</span><span class="df-lp-count">%d</span></h1>`,
		html.EscapeString(textutil.PrettySection(section)), len(docs)))
	b.WriteString(`<input class="df-filter" id="df-filter" placeholder="Filter this list…" autocomplete="off">`)
	b.WriteString(`</div>`)
	b.WriteString(`<p class="df-muted" id="df-filtercount"></p>`)

	// A-Z jump bar.
	keys := make([]byte, len(docs))
	firstIdx := map[byte]int{}
	for i, d := range docs {
		k := textutil.AZKey(d.Title)
		keys[i] = k
		if _, ok := firstIdx[k]; !ok {
			firstIdx[k] = i
		}
	}
	azOrder := "#abcdefghijklmnopqrstuvwxyz"
	b.WriteString(`<div class="df-az">`)
	for k := 0; k < len(azOrder); k++ {
		ch := azOrder[k]
		label := strings.ToUpper(string(ch))
		if idx, ok := firstIdx[ch]; ok {
			b.WriteString(fmt.Sprintf(`<a class="df-az-item" href="#df-az-%d">%s</a>`, idx, label))
		} else {
			b.WriteString(fmt.Sprintf(`<span class="df-az-item off">%s</span>`, label))
		}
	}
	b.WriteString(`</div>`)

	b.WriteString(`<ul class="df-list" id="df-list">`)
	for i, d := range docs {
		idAttr := ""
		if firstIdx[keys[i]] == i {
			idAttr = fmt.Sprintf(` id="df-az-%d"`, i)
		}
		b.WriteString(fmt.Sprintf(`<li%s><a href="%s">%s</a></li>`,
			idAttr, html.EscapeString(textutil.OrDefault(d.Route, "/")), html.EscapeString(textutil.OrDefault(d.Title, "(untitled)"))))
	}
	b.WriteString(`</ul></div>`)
	b.WriteString(`<script>(function(){var i=document.getElementById('df-filter'),` +
		`L=document.getElementById('df-list'),li=L.getElementsByTagName('li'),` +
		`c=document.getElementById('df-filtercount');` +
		`function run(){var q=i.value.trim().toLowerCase(),n=0;` +
		`for(var k=0;k<li.length;k++){var t=li[k].textContent.toLowerCase();` +
		`var show=!q||t.indexOf(q)!==-1;li[k].style.display=show?'':'none';if(show)n++;}` +
		`c.textContent=(q&&n<li.length)?(n+' of '+li.length+' shown'):'';}` +
		`i.addEventListener('input',run);i.focus();run();})();</script>`)
	s.renderShell(w, r, textutil.PrettySection(section), "", b.String())
}

// handleRandom redirects to a random content page (excludes the 'url' section).
func (s *Server) handleRandom(w http.ResponseWriter, r *http.Request) {
	pages := s.DB.RandomContentPages(1)
	if len(pages) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	route := pages[0].Route
	if route == "" {
		route = "/"
	}
	http.Redirect(w, r, route, http.StatusFound)
}

// handleExternal shows a "blocked external link" page.
func (s *Server) handleExternal(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	var b strings.Builder
	b.WriteString(`<div class="df-external">`)
	b.WriteString(`<h1>External link — blocked</h1>`)
	b.WriteString(`<p>All outside links are blocked in this offline copy. If you want to open it, ` +
		`<b>copy and paste</b> the address below into your browser's address bar.</p>`)
	b.WriteString(`<p class="df-muted">This resource is for offline use only.</p>`)
	if to != "" {
		b.WriteString(fmt.Sprintf(`<p><code class="df-url">%s</code></p>`, html.EscapeString(to)))
	}
	b.WriteString(`<p style="margin-top:18px"><a href="/">← back to home</a></p>`)
	b.WriteString(`</div>`)
	s.renderShell(w, r, "External link", "", b.String())
}
