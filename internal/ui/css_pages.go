package ui

// cssPages covers the rendered app pages — search results, A-Z list pages, and
// external-link pages — plus the cross-cutting responsive overrides that hide
// the desktop nav and tighten the page layout on small screens.
const cssPages = `
/* ---- app pages (home / search / list / external) ---- */
.df-searchpage,.df-listpage,.df-external{max-width:940px;margin:30px auto;padding:0 18px;font-family:system-ui,sans-serif;color:#1f2937}
.df-searchpage h1,.df-listpage h1{font-size:24px;margin:0 0 6px}
.df-muted{color:#6b7280;font-size:13px;margin:0 0 8px}
.df-filter{width:220px;max-width:100%;height:36px;border:1px solid #ece3d4;border-radius:999px;
  padding:0 16px;font-size:13px;outline:none;background:#fff;margin:0;
  transition:border-color .12s,box-shadow .12s}
.df-filter:focus{border-color:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.18)}
.df-results{list-style:none;padding:0;margin:8px 0}
.df-results li{display:flex;gap:14px;padding:15px 0;border-bottom:1px solid #eef0f3}
.df-thumb{width:96px;height:72px;object-fit:cover;flex:0 0 auto;background:#f3f4f6;border:1px solid #e5e7eb}
.df-thumb.ph{display:flex;align-items:center;justify-content:center;color:#cbd5e1;font-size:10px}
.df-rbody{min-width:0}
.df-title{font-size:16px;font-weight:600;color:#007bff;text-decoration:none}
.df-title:hover{color:#0056b3;text-decoration:underline}
.df-badge{display:inline-block;margin-left:8px;font-size:10px;text-transform:uppercase;letter-spacing:.03em;color:#6b7280;
  border:1px solid #d1d5db;border-radius:999px;padding:1px 8px;vertical-align:middle}
.df-count{display:inline-block;margin-left:6px;font-size:11px;color:#c2410c;font-weight:600;background:#fff7ed;
  border-radius:999px;padding:1px 8px;vertical-align:middle}
.df-snip{color:#4b5563;font-size:13.5px;margin-top:5px;line-height:1.5}
.df-snip mark,.df-results mark{background:#fde68a;color:#111827;padding:0 1px;border-radius:2px}
.df-listpage{max-width:680px;margin:30px auto;background:transparent;border:1px solid #ece3d4;
  border-radius:16px;padding:30px 28px 26px}
.df-lp-head{display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap;margin-bottom:6px}
.df-listpage h1{text-align:left;margin:0;line-height:1.2}
.df-lp-title{font-size:30px;font-weight:800;letter-spacing:-.01em;padding:0 0 4px;
  background:linear-gradient(92deg,#f59e0b,#f97316,#ef4444);
  -webkit-background-clip:text;background-clip:text;color:transparent;-webkit-text-fill-color:transparent}
.df-lp-count{display:inline-block;vertical-align:middle;margin-left:8px;font-size:13px;font-weight:700;
  color:#92400e;background:#fdf2e3;border-radius:999px;padding:3px 11px;letter-spacing:.02em}
.df-list{list-style:none;padding:0;margin-top:0}
.df-list li{border-bottom:1px solid #ece3d4;scroll-margin-top:130px}
.df-list li>a{display:block;padding:8px 10px;color:#374151;text-decoration:none;font-size:14.5px;
  border-radius:7px;transition:background .12s,color .12s}
.df-list li>a:hover{background:#fdf2e3;color:#c2410c}
/* A-Z jump bar */
.df-az{display:grid;grid-template-columns:repeat(9,1fr);gap:6px;position:sticky;top:58px;z-index:50;background:#fff;margin:12px 0 0;padding:8px 0 10px;border-bottom:1px solid #ece3d4}
.df-az-item{display:flex;align-items:center;justify-content:center;
  height:44px;min-width:0;border-radius:9px;font-size:17px;font-weight:700;
  line-height:1;color:#b45309;text-decoration:none;background:#faf7f1;border:1px solid #ece3d4;
  transition:background .1s,color .1s,border-color .1s}
.df-az-item:hover{background:#ea580c;color:#fff;border-color:#ea580c}
.df-az-item.off{color:#d6cab3;background:transparent;border-color:transparent;cursor:default}
@media (max-width:480px){
  .df-az{gap:4px}
  .df-az-item{height:36px;font-size:14px;border-radius:7px}
}
/* list page — mobile */
@media (max-width:640px){
  .df-listpage{padding:20px 16px;border-radius:12px}
  .df-lp-title{font-size:24px}
  .df-filter{width:100%}
}
.df-external h1{font-size:22px}
a[href^="/external?to="]{color:#9ca3af!important;text-decoration:line-through}
.df-noembed{display:inline-block;font-size:12px;color:#6b7280;background:#f9fafb;border:1px dashed #d1d5db;
  padding:4px 8px;border-radius:6px;margin:4px 0}

/* ---- responsive ---- */
@media (max-width:860px){
  #difi-bar .difi-nav,#difi-bar .difi-search,#difi-bar .difi-dice{display:none!important}
  #difi-bar .difi-burger{display:inline-flex!important;margin-left:6px}
  #difi-bar .difi-mdice{display:flex}
}
@media (max-width:640px){
  body{font-size:15px}
  .df-hero{margin:18px 0}
  .df-credit-row{flex-direction:column;align-items:center;text-align:center}
  .df-thumb{width:72px;height:54px}
  .df-results li{gap:10px}
  .df-searchpage,.df-listpage,.df-external,.df-home{margin:18px auto}
}`
