package ui

import (
	"fmt"
	"html"
	"net/url"
)

// PinButtonHTML returns an inlaid toolbar at the top of content pages with a
// pin/bookmark button. Logged-out viewers get nothing (same as the old
// user==nil behavior).
func PinButtonHTML(route string, viewer *Viewer) string {
	if viewer == nil || !viewer.LoggedIn {
		return ""
	}
	routeEsc := html.EscapeString(route)
	dlHref := "/download?page=" + url.QueryEscape(route)
	return fmt.Sprintf(`
<div id="df-pin-bar" data-route="%s">
  <a id="df-dl-btn" href="%s" title="Download this page as a self-contained HTML file">&#8623; Download</a>
  <button id="df-pin-btn" onclick="dfTogglePin()" title="Bookmark this page">
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 2l2.4 7.4H22l-6 4.6 2.3 7.4-6.3-4.6-6.3 4.6 2.3-7.4-6-4.6h7.6z"/>
    </svg>
    <span id="df-pin-label">Bookmark</span>
  </button>
</div>
<script>
(function(){
var btn=document.getElementById('df-pin-btn');
var label=document.getElementById('df-pin-label');
if(!btn)return;
var route=document.getElementById('df-pin-bar').dataset.route;
fetch('/api/bookmark/check?route='+encodeURIComponent(route))
.then(function(r){return r.json()})
.then(function(d){if(d&&d.bookmarked){btn.classList.add('pinned');label.textContent='Bookmarked'}});
window.dfTogglePin=function(){
  if(btn.classList.contains('pinned')){
    fetch('/api/bookmark',{method:'DELETE',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({route:route})})
    .then(function(){btn.classList.remove('pinned');label.textContent='Bookmark'});
  } else {
    var title=document.title||route;
    fetch('/api/bookmark',{method:'POST',headers:{'Content-Type':'application/json'},
      body:JSON.stringify({route:route,title:title})})
    .then(function(r){if(r.ok){btn.classList.add('pinned');label.textContent='Bookmarked'}});
  }
};
})();
</script>
`, routeEsc, html.EscapeString(dlHref))
}
