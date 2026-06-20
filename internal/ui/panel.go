package ui

import (
	"fmt"
	"html"
)

// navGroup is a top-level navbar category that expands into a sub-menu of sections.
type navGroup struct {
	Label string
	Items [][2]string // (display label, section key)
}

// NavGroups organizes the site's sections into a few friendly categories.
var NavGroups = []navGroup{
	{"Materials", [][2]string{
		{"Oxides", "oxide"}, {"Materials", "material"}, {"Minerals", "mineral"},
		{"Properties", "property"}, {"Typecodes", "typecode"},
	}},
	{"Recipes & Firing", [][2]string{
		{"Recipes", "recipe"}, {"Firing Schedules", "schedule"}, {"Temperatures", "temperature"},
	}},
	{"Learn", [][2]string{
		{"Glossary", "glossary"}, {"Articles", "article"}, {"Troubles", "trouble"}, {"Tests", "test"},
	}},
	{"Media & More", [][2]string{
		{"Media", "video"}, {"Pictures", "picture"}, {"Hazards", "hazard"}, {"Projects", "project"},
		{"Consultants", "consultants"}, {"Schools", "schools"}, {"Stores", "stores"}, {"Downloads", "downloads"},
	}},
}

// PanelHTML is the overlay panel injected at the top of every page: back nav,
// the Digitalfire logo, a grouped expanding navbar, and a search box (hidden on
// the home/search pages, which have their own big search).
func PanelHTML(route string, viewer *Viewer) string {
	// Desktop inline nav (details/summary dropdowns).
	nav := ""
	for _, g := range NavGroups {
		items := ""
		for _, it := range g.Items {
			items += fmt.Sprintf(`<a href="/list/%s">%s</a>`, it[1], it[0])
		}
		nav += `<details class="difi-group"><summary>` + g.Label + `</summary>` +
			`<div class="difi-submenu">` + items + `</div></details>`
	}
	// Mobile drawer body: accordion groups mirroring the desktop nav.
	body := ""
	for _, g := range NavGroups {
		links := ""
		for _, it := range g.Items {
			links += fmt.Sprintf(`<a href="/list/%s">%s</a>`, it[1], it[0])
		}
		body += `<details class="difi-dgroup"><summary>` + g.Label + `</summary>` +
			`<div class="difi-dlinks">` + links + `</div></details>`
	}

	// Mobile drawer: user/account section (shown only in drawer, hidden in desktop bar on mobile).
	drawerUser := ""
	if viewer != nil && viewer.LoggedIn {
		displayName := viewer.Username
		if displayName == "" {
			displayName = viewer.Email
		}
		drawerUser = `<details class="difi-dgroup"><summary>Account</summary>` +
			`<div class="difi-dlinks">`
		if viewer.Role == "admin" {
			drawerUser += `<a href="/admin">Admin Panel</a>`
		}
		drawerUser += fmt.Sprintf(`<a href="#" style="color:#9ca3af;cursor:default">%s (%s)</a>`,
				html.EscapeString(displayName), viewer.Role) +
			`<a href="/logout">Logout</a>` +
			`</div></details>`
	} else {
		drawerUser = `<div class="difi-dlinks" style="padding:12px 0;display:flex;gap:8px">` +
			`<a href="/login" style="flex:1;text-align:center">Login</a>` +
			`<a href="/register" style="flex:1;text-align:center">Register</a>` +
			`</div>`
	}

	// User menu (right side of nav bar — desktop only).
	userArea := ""
	if viewer != nil && viewer.LoggedIn {
		displayName := viewer.Username
		if displayName == "" {
			displayName = viewer.Email
		}
		userArea = `<div class="difi-user">`
		if viewer.Role == "admin" {
			userArea += `<a href="/admin" class="difi-admin-link" title="Admin">Admin</a>`
		}
		userArea += fmt.Sprintf(`<span class="difi-uname">%s</span>`, html.EscapeString(displayName)) +
			`<span class="difi-badge-role">` + viewer.Role + `</span>` +
			`<a href="/logout" class="difi-logout">Logout</a></div>`
	} else {
		userArea = `<div class="difi-user"><a href="/login">Login</a> <a href="/register">Register</a></div>`
	}

	js := `<script>(function(){` +
		// Single-open desktop menus.
		`var g=document.querySelectorAll('#difi-bar details.difi-group');` +
		`g.forEach(function(x){x.addEventListener('toggle',function(){if(x.open){g.forEach(function(o){if(o!==x)o.open=false})}})});` +
		// Desktop nav: open submenus on hover, auto-close when the pointer leaves.
		`g.forEach(function(x){x.addEventListener('mouseenter',function(){g.forEach(function(o){o.open=(o===x)})});x.addEventListener('mouseleave',function(){x.open=false})});` +
		// Mobile accordion: single-open + reveal the expanded group.
		`var mg=document.querySelectorAll('#difi-drawer details.difi-dgroup');` +
		`mg.forEach(function(x){x.addEventListener('toggle',function(){if(x.open){mg.forEach(function(o){if(o!==x)o.open=false});var s=x.querySelector('summary');if(s)setTimeout(function(){s.scrollIntoView({behavior:'smooth',block:'nearest'})},60)}})});` +
		// Mobile slide-in panel open/close.
		`var bk=document.querySelector('.difi-burger'),dr=document.getElementById('difi-drawer'),bd=document.getElementById('difi-backdrop');` +
		`function mc(){dr.classList.remove('open');bd.classList.remove('open');bk.setAttribute('aria-expanded','false')}` +
		`function mo(){dr.classList.add('open');bd.classList.add('open');bk.setAttribute('aria-expanded','true')}` +
		`if(bk&&dr&&bd){bk.addEventListener('click',function(){dr.classList.contains('open')?mc():mo()});` +
		`bd.addEventListener('click',mc);var dc=document.querySelector('.difi-dclose');if(dc)dc.addEventListener('click',mc);` +
		`document.addEventListener('keydown',function(e){if(e.key==='Escape')mc()});` +
		`window.addEventListener('resize',function(){if(window.innerWidth>860)mc()})}` +
		// Wrap tables in horizontal-scroll containers (mobile-friendly legacy tables).
		`function wt(){document.querySelectorAll('table').forEach(function(t){` +
		`if(!t.closest('.df-tscroll')&&!t.closest('#difi-bar')&&!t.closest('#difi-drawer')){` +
		`var p=document.createElement('div');p.className='df-tscroll';t.parentNode.insertBefore(p,t);p.appendChild(t)}})};` +
		`if(document.readyState==='loading')document.addEventListener('DOMContentLoaded',wt);else wt();` +
		`})();</script>`
	// Bookmarks dropdown — always visible. Logged-out users get a static login
	// prompt (no fetch); logged-in users get the list loaded via JS on open.
	bookmarkArea := `<details class="difi-group difi-bookmarks-group"><summary>&#9733; Bookmarks</summary>` +
		`<div class="difi-submenu difi-bookmarks-menu" id="difi-bookmarks-dropdown">`
	if viewer != nil && viewer.LoggedIn {
		bookmarkArea += `<p class="df-muted difi-bm-loading">Loading…</p>`
	} else {
		bookmarkArea += `<p class="df-muted">Log in to use bookmarks. <a href="/login">Log in &rarr;</a></p>`
	}
	bookmarkArea += `</div></details>`

	return `<div id="difi-bar">` +
		`<button class="difi-ic" onclick="history.back()" title="Back" aria-label="Back">&#8249;</button>` +
		`<a class="difi-brand" href="/" title="Home">` +
		`<img src="/images/logo.png" alt="Digitalfire"></a>` +
		`<nav class="difi-nav">` + nav + `</nav>` +
		`<form class="difi-search" action="/search" method="get">` +
		`<input type="text" name="q" placeholder="Search all pages…">` +
		`<button type="submit" aria-label="Search"><svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="7"/><path d="M21 21l-4.3-4.3"/></svg></button></form>` +
		`<a class="difi-dice" href="/random" title="Open a random page" aria-label="Open a random page"></a>` +
		`<div class="difi-nav-spacer"></div>` +
		bookmarkArea +
		userArea +
		`<a class="difi-mdice" href="/random" title="Open a random page" aria-label="Open a random page"></a>` +
		`<button class="difi-burger" aria-label="Menu" aria-expanded="false"><svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round"><path d="M4 7h16M4 12h16M4 17h16"/></svg></button>` +
		`</div>` +
		`<div class="difi-backdrop" id="difi-backdrop"></div>` +
		`<div id="difi-drawer">` +
		`<div class="difi-dhead">` +
		`<form class="difi-search" action="/search" method="get">` +
		`<input type="text" name="q" placeholder="Search all pages…">` +
		`<button type="submit">Search</button></form>` +
		`<button class="difi-dclose" aria-label="Close menu"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M6 6l12 12M18 6L6 18"/></svg></button></div>` +
		`<div class="difi-dbody">` + body + drawerUser + `</div>` +
		`</div>` + js +
		bookmarkJS
}

// injectFragment is dropped right after <body> on captured pages: styles + panel.
func injectFragment(route string, viewer *Viewer) string {
	return DifiCSS() + PanelHTML(route, viewer)
}

// bookmarkJS is injected into the nav panel to load bookmarks into the dropdown.
const bookmarkJS = `<script>
(function(){
var bmGroup=document.querySelector('.difi-bookmarks-group');
if(!bmGroup)return;
var loaded=false;
bmGroup.addEventListener('toggle',function(){
  if(!bmGroup.open||loaded)return;
  loaded=true;
  var menu=document.getElementById('difi-bookmarks-dropdown');
  if(!menu)return;
  // Skip the fetch when the menu already holds static content (logged-out prompt).
  if(!menu.querySelector('.difi-bm-loading'))return;
  fetch('/api/bookmarks').then(function(r){
    if(r.status===401){menu.innerHTML='<p class="df-muted">Log in to use bookmarks. <a href="/login">Log in &rarr;</a></p>';return null}
    return r.json()
  }).then(function(bms){
    if(!bms)return;
    if(bms.length===0){menu.innerHTML='<p class="df-muted">Visit a page to bookmark it.</p>';return}
    menu.innerHTML=bms.map(function(b){
      var t=b.Title||b.Route;
      return '<a href="'+b.Route+'">'+t+'</a>';
    }).join('');
  }).catch(function(){menu.innerHTML='<p class="df-muted">Failed to load.</p>'});
});
})();
</script>`
