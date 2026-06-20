package app

import (
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strings"

	"difilo/internal/textutil"
	"difilo/internal/ui"
)

// handleHome renders the homepage: hero with search, section buttons, random
// explore cards, recent comments, and the Tony Hansen attribution.
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	counts := s.DB.SectionCounts()
	if counts == nil {
		counts = map[string]int{}
	}

	// Flat list of section buttons, sorted by page count (most content first).
	type sec struct {
		label, key string
		n          int
	}
	var secs []sec
	seen := map[string]bool{}
	for _, g := range ui.NavGroups {
		for _, it := range g.Items {
			if seen[it[1]] {
				continue
			}
			seen[it[1]] = true
			secs = append(secs, sec{it[0], it[1], counts[it[1]]})
		}
	}
	// sort.Slice would normally be used, but a simple insertion keeps the
	// import surface tiny; the slice is small (a few dozen sections).
	for i := 1; i < len(secs); i++ {
		for j := i; j > 0; j-- {
			if secs[j].n > secs[j-1].n || (secs[j].n == secs[j-1].n && secs[j].label < secs[j-1].label) {
				secs[j], secs[j-1] = secs[j-1], secs[j]
			} else {
				break
			}
		}
	}

	var b strings.Builder
	b.WriteString(`<div class="df-home">`)

	// Hero with search bar and random background.
	heroBG := ""
	if len(s.heroImages) > 0 {
		heroBG = s.heroImages[rand.Intn(len(s.heroImages))]
	}
	b.WriteString(`<div class="df-herowrap"`)
	if heroBG != "" {
		b.WriteString(fmt.Sprintf(` style="background-image:url('%s')"`, html.EscapeString(heroBG)))
	}
	b.WriteString(`>`)
	// Top-right hint
	b.WriteString(`<div class="df-hero-login">Bookmark &amp; download pages after logging in</div>`)
	b.WriteString(`<div class="df-searchrow">`)
	b.WriteString(`<form class="df-bigsearch" action="/search" method="get">`)
	b.WriteString(`<input type="text" name="q" placeholder="Search all pages…" autocomplete="off" autofocus>`)
	b.WriteString(`<button type="submit"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/></svg></button></form>`)
	b.WriteString(`<a class="df-dice" href="/random" title="Open a random page" aria-label="Open a random page" data-route="/random"></a>`)
	b.WriteString(`</div>`)

	// Download button — password is verified server-side via /api/download-check
	// (Gap 3: previously the password 'pancake' was hardcoded client-side).
	downloadURL := s.DB.GetSetting("download_url")
	if downloadURL != "" {
		b.WriteString(`<p class="df-download"><a href="#" onclick="dfCheckDownload();return false">&#11015; Download full offline version — v June 1, 2026</a></p>`)
	} else {
		b.WriteString(`<p class="df-download"><a href="#" onclick="alert('Download coming soon');return false">&#11015; Download full offline version — v June 1, 2026</a></p>`)
	}
	b.WriteString(`</div>`)

	// Section buttons.
	b.WriteString(`<div class="df-btns">`)
	for _, sc := range secs {
		b.WriteString(fmt.Sprintf(`<a class="df-btn" href="/list/%s">%s <span class="df-num">%d</span></a>`,
			sc.key, sc.label, sc.n))
	}
	b.WriteString(`</div>`)

	// "Explore" — 2 random articles + 6 random from all sections.
	articlePicks := s.DB.RandomCardsBySection("article", 2)
	generalPicks := s.DB.RandomCardPages(6)
	dbPicks := append(articlePicks, generalPicks...)
	if len(dbPicks) > 0 {
		b.WriteString(`<div class="df-explore-head">`)
		b.WriteString(`<button class="df-explore-refresh" onclick="dfRefreshCards()" title="Show new pages">`)
		b.WriteString(`<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M23 4v6h-6"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>`)
		b.WriteString(`<span>Refresh List</span>`)
		b.WriteString(`</button></div>`)
		b.WriteString(`<div class="df-picks" id="df-picks">`)
		for i, p := range dbPicks {
			route := p.Route
			if route == "" {
				route = "/"
			}
			cls := "df-pick"
			if i >= 4 {
				cls = "df-pick df-pick-extra"
			}
			b.WriteString(fmt.Sprintf(
				`<a class="%s" href="%s">`+
					`<div class="df-pick-img" style="background-image:url('%s')"></div>`+
					`<div class="df-pick-body">`+
					`<div class="df-pick-title">%s</div>`+
					`<div class="df-pick-sec">%s</div>`+
					`<div class="df-pick-text">%s</div>`+
					`</div></a>`,
				cls,
				html.EscapeString(route),
				html.EscapeString(p.Thumb),
				html.EscapeString(textutil.OrDefault(p.Title, "(untitled)")),
				html.EscapeString(textutil.PrettySection(p.Section)),
				html.EscapeString(textutil.OrDefault(p.MetaDescription, textutil.Excerpt(p.BodyText, p.Title, 150)))))
		}
		b.WriteString(`</div>`)
	}

	// Recent comments.
	recentComments, _ := s.DB.GetRecentComments(8)
	if len(recentComments) > 0 {
		b.WriteString(`<h2>Recent Comments</h2>`)
		b.WriteString(`<ul class="df-recent-comments">`)
		for _, c := range recentComments {
			title := c.Route
			if p, _ := s.DB.GetPageByRoute(c.Route); p != nil {
				title = p.Title
			}
			author := c.Username
			if author == "" {
				author = "(unknown)"
			}
			b.WriteString(fmt.Sprintf(`<li>`+
				`<div class="df-rc-meta"><span class="df-rc-author">%s</span>`+
				`<span class="df-rc-badge">%s</span>`+
				`<span class="df-rc-on">commented on</span>`+
				`<a class="df-rc-page" href="%s">%s</a></div>`+
				`<div class="df-rc-snippet">%s</div></li>`,
				html.EscapeString(author), html.EscapeString(c.Role),
				html.EscapeString(c.Route), html.EscapeString(title),
				strings.ReplaceAll(html.EscapeString(c.Body), "\n", "<br>")))
		}
		b.WriteString(`</ul>`)
	}

	// Tony Hansen attribution.
	b.WriteString(`<div class="df-credit"><div class="df-credit-row">`)
	b.WriteString(`<img class="df-tony" src="/images/digitalfire.com/images/SignaturePhotoSmall.jpg" alt="Tony Hansen">`)
	b.WriteString(`<div>`)
	b.WriteString(`<p class="df-byline">Everything here was created by <b>Tony Hansen</b> &#8212; ceramic engineer and ` +
		`the author of Digitalfire and Insight-Live. This is an offline record of his life's work: 35+ years of ` +
		`research into ceramic materials, glaze chemistry, and firing (captured before June 1, 2026). If his work ` +
		`has helped you, the best thank-you is to support him &#8212; buy him a coffee below.</p>`)
	socials := [][3]string{
		{"https://instagram.com/tonyatdigitalfire", "/images/digitalfire.com/images/instagram.jpg", "Instagram"},
		{"https://www.facebook.com/insightlive", "/images/digitalfire.com/images/facebook.jpg", "Facebook"},
		{"https://x.com/digitalfiretony", "/images/digitalfire.com/images/twitter.jpg", "X (Twitter)"},
		{"https://www.linkedin.com/in/tonywhansen/", "/images/digitalfire.com/images/linkedin.jpg", "LinkedIn"},
		{"https://www.pinterest.ca/tonywhansen", "/images/digitalfire.com/images/pinterest.jpg", "Pinterest"},
		{"https://www.threads.net/@tonyatdigitalfire", "/images/digitalfire.com/images/threads.jpg", "Threads"},
	}
	b.WriteString(`<p class="df-socials">`)
	for _, so := range socials {
		b.WriteString(fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener" title="%s"><img src="%s" alt="%s"></a>`,
			html.EscapeString(so[0]), so[2], so[1], so[2]))
	}
	b.WriteString(`</p>`)
	b.WriteString(`<div class="df-credit-btns">`)
	b.WriteString(`<a class="df-coffee" href="https://ko-fi.com/tonyhansen" target="_blank" rel="noopener">` +
		`<img src="/images/digitalfire.com/images/ko-fi.svg" alt=""> Buy Tony a Coffee</a>`)
	b.WriteString(`<a class="df-orig-site" href="https://www.digitalfire.com" target="_blank" rel="noopener">Visit Digitalfire.com &#8594;</a>`)
	b.WriteString(`</div>`)
	b.WriteString(`<p class="df-disclaimer">This offline copy preserves the educational content of Digitalfire.com ` +
		`for archival and study purposes. It was captured for personal use before June 1, 2026, in case the site ` +
		`goes offline. The original site may have newer or updated content. If you find this useful, please ` +
		`support Tony's work through the links above.</p>`)
	b.WriteString(`</div>`)
	b.WriteString(`</div>`)

	b.WriteString(`</div>`)

	// Client-side scripts: explore-card refresh + server-side download password
	// check (Gap 3 — the password is no longer embedded in the page).
	b.WriteString(`<script>
function dfRefreshCards(){
  var grid=document.getElementById('df-picks');
  if(!grid)return;
  var shim='';
  for(var i=0;i<8;i++){
    var cls=i>=4?'df-shimmer df-pick-extra':'df-shimmer';
    shim+='<div class="'+cls+'"><div class="df-shimmer-img"></div>'+
      '<div class="df-shimmer-body">'+
      '<div class="df-shimmer-line s1"></div>'+
      '<div class="df-shimmer-line s2"></div>'+
      '<div class="df-shimmer-line s3"></div>'+
      '</div></div>';
  }
  grid.innerHTML=shim;
  fetch('/api/explore').then(function(r){return r.json()}).then(function(cards){
    if(!cards||!cards.length)return;
    grid.innerHTML=cards.map(function(c,i){
      var cls=i>=4?'df-pick df-pick-extra':'df-pick';
      return '<a class="'+cls+'" href="'+c.Route+'" style="animation:dfCardIn .3s ease '+(i*0.04)+'s both">'+
        '<div class="df-pick-img" style="background-image:url(\''+c.Thumb+'\')"></div>'+
        '<div class="df-pick-body">'+
        '<div class="df-pick-title">'+c.Title+'</div>'+
        '<div class="df-pick-sec">'+c.Section+'</div>'+
        '<div class="df-pick-text">'+c.Desc+'</div>'+
        '</div></a>';
    }).join('');
  });
}
function dfCheckDownload(){
  var p=prompt('Enter password to download:');
  if(p===null)return;
  fetch('/api/download-check',{method:'POST',headers:{'Content-Type':'application/json'},
    body:JSON.stringify({password:p})})
  .then(function(r){return r.json()})
  .then(function(d){
    if(d&&d.ok&&d.url){window.location.href=d.url}
    else{alert('Wrong password')}
  })
  .catch(function(){alert('Download check failed')});
}
</script>`)

	s.renderShell(w, r, "DIFI-LOCAL", "df-homepage", b.String())
}

// handleAPIDownloadCheck verifies a download password against the
// download_password setting and, on success, returns the configured download
// URL as JSON. This is the server-side replacement for the legacy hardcoded
// client-side 'pancake' comparison (Gap 3).
//
// When no password is configured yet, the legacy default "pancake" is used
// (for backward compatibility) and persisted to the settings table so it can
// be managed from the admin panel going forward.
func (s *Server) handleAPIDownloadCheck(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	expected := s.DB.GetSetting("download_password")
	if expected == "" {
		expected = "pancake"
		_ = s.DB.SetSetting("download_password", expected)
	}
	if body.Password != expected {
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": false})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "url": s.DB.GetSetting("download_url")})
}
