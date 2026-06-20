package ui

// cssBase is the "refresh layer" applied to every page: modern typography and
// link colors layered over captured legacy pages, Bootstrap table restyling,
// the horizontal-scroll wrapper for wide tables, and rules that hide leftover
// site chrome (crisp chat, the original navbar, Google CSE, and the panel's own
// search box on pages that ship their own big search).
const cssBase = `/* ---- refresh layer (captured pages + app pages) ---- */
html{overflow-x:hidden;scroll-behavior:smooth}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif!important;
  color:#1f2937;font-size:16px;line-height:1.6}
body p{line-height:1.65}
body a{color:#007bff}
body a:hover{color:#0056b3;text-decoration:underline}
body img{max-width:100%;height:auto;border-radius:6px}
body h1,body h2,body h3,body h4{letter-spacing:-.01em;color:#111827}
body table{border-collapse:collapse;font-size:14px}
/* modernize the Bootstrap data tables on captured pages */
body .table{background:#fff;margin:14px 0}
body .table>thead>tr>th{background:#f4efe5;color:#57534e;font-size:12px;font-weight:700;
  text-transform:uppercase;letter-spacing:.04em;border-bottom:2px solid #e0d6c5}
body .table td,body .table th{padding:9px 12px;border-color:#eee6d8}
body .table tbody tr:nth-child(even){background:#fbf8f1}
body .table tbody tr:hover{background:#fdf2e3}
.df-tscroll{overflow-x:auto;margin:14px 0;-webkit-overflow-scrolling:touch}
/* hide leftover live-chat widget markup */
#crisp-chatbox,[class^="crisp-"],[id^="crisp-"]{display:none!important}
/* hide the site's original navbar (we replace it with our overlay panel) */
nav.navbar{display:none!important}
/* hide the site's Google CSE box on the home page */
.gcse-search{display:none!important}
/* hide the panel's search box on pages that have their own big search (home + search) */
body.df-homepage #difi-bar .difi-search,body.df-searchpage #difi-bar .difi-search,
body.df-homepage #difi-bar .difi-dice,body.df-searchpage #difi-bar .difi-dice{display:none!important}`
