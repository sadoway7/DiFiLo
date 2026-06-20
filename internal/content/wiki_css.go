package content

// WikiCSS contains styles for wiki-rendered pages and the image lightbox.
// Aesthetic: warm cream + stone neutrals + orange accent — wiki meets IMDB.
const WikiCSS = `
<style>
/* ---- wiki page layout ---- */
.df-wiki{display:flex;gap:36px;max-width:1200px;margin:0 auto;padding:28px 24px 60px}
.df-wiki-main{flex:1;min-width:0}
.df-wiki-sidebar{width:280px;flex-shrink:0;font-size:13px}

/* wiki header */
.df-wiki-header{margin-bottom:32px;padding-bottom:20px;border-bottom:1px solid #e7dfd2}
.df-wiki-badge{display:inline-block;background:linear-gradient(135deg,#fef3e2,#fde8c8);color:#92400e;
  font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:.07em;
  padding:4px 12px;border-radius:20px;margin-bottom:10px;border:1px solid #f5d9a8}
.df-wiki-header h1{font-size:32px;margin:0 0 8px;line-height:1.15;color:#1c1917;
  letter-spacing:-.025em;font-weight:800}
.df-wiki-desc{color:#57534e;font-size:16px;line-height:1.55;margin:10px 0 0;max-width:680px}
.df-wiki-byline{color:#a8a29e;font-size:12.5px;margin-top:10px;
  display:inline-flex;align-items:center;gap:8px;
  background:#fafaf9;padding:4px 12px;border-radius:20px}
.df-wiki-byline::before{content:"";width:6px;height:6px;border-radius:50%;background:#ea580c}

/* wiki content — rendered markdown */
.df-wiki-content{font-size:16px;line-height:1.75;color:#292524;word-wrap:break-word}

/* headings — section hierarchy with accent treatment */
.df-wiki-content h1{font-size:22px;margin:36px 0 16px;padding-bottom:8px;
  border-bottom:2px solid #f4efe5;color:#1c1917;font-weight:700}
.df-wiki-content h2{font-size:19px;margin:32px 0 12px;color:#1c1917;font-weight:700;
  padding-left:12px;border-left:3px solid #ea580c}
.df-wiki-content h3{font-size:16px;margin:24px 0 8px;color:#44403c;font-weight:700}
.df-wiki-content h4{font-size:12px;margin:20px 0 8px;color:#78716c;
  text-transform:uppercase;letter-spacing:.06em;font-weight:700;
  padding-bottom:4px;border-bottom:1px solid #f5f0e8}

/* paragraphs */
.df-wiki-content p{margin:14px 0}

/* tables — polished data cards */
.df-wiki-content table{width:100%;border-collapse:separate;border-spacing:0;margin:20px 0;
  font-size:14px;background:#fff;border:1px solid #e7dfd2;border-radius:12px;
  overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.04)}
.df-wiki-content thead{background:linear-gradient(180deg,#f8f2e6,#f4efe5)}
.df-wiki-content th{padding:11px 16px;text-align:left;
  font-weight:700;font-size:11px;text-transform:uppercase;letter-spacing:.05em;
  color:#78716c;border-bottom:2px solid #e0d6c5;white-space:nowrap}
.df-wiki-content td{padding:10px 16px;border-bottom:1px solid #f5f0e8;vertical-align:top;
  line-height:1.5}
.df-wiki-content tbody tr:last-child td{border-bottom:none}
.df-wiki-content tbody tr:nth-child(even){background:#fdfbf6}
.df-wiki-content tbody tr:hover{background:#fef9f0}
/* key-value style tables: bold first column (labels) */
.df-wiki-content td:first-child{font-weight:600;color:#44403c}
/* numeric data alignment */
.df-wiki-content td:nth-child(2),.df-wiki-content th:nth-child(2){text-align:left}

/* images — clickable cards */
.df-wiki-content img{max-width:100%;height:auto;border-radius:12px;margin:14px 0;
  box-shadow:0 2px 10px rgba(0,0,0,.08);cursor:pointer;
  transition:transform .18s,box-shadow .18s}
.df-wiki-content img:hover{transform:scale(1.015);box-shadow:0 6px 20px rgba(0,0,0,.12)}
.df-wiki-content a img{display:block;margin:14px 0}

/* links — clean, subtle */
.df-wiki-content a{color:#0369a1;text-decoration:none;font-weight:500;
  border-bottom:1px solid rgba(3,105,161,.25);transition:.12s}
.df-wiki-content a:hover{color:#075985;border-bottom-color:#075985}
.df-wiki-content a{cursor:pointer}

/* lists */
.df-wiki-content ul,.df-wiki-content ol{padding-left:24px;margin:14px 0}
.df-wiki-content li{margin:6px 0}
.df-wiki-content ul li::marker{color:#ea580c}

/* blockquotes — callout style */
.df-wiki-content blockquote{border-left:4px solid #ea580c;padding:10px 0 10px 18px;
  color:#57534e;margin:18px 0;background:linear-gradient(90deg,#fef9f0,transparent);
  border-radius:0 8px 8px 0;font-style:italic}

/* inline code */
.df-wiki-content code{background:#f4efe5;padding:2px 7px;border-radius:5px;
  font-size:13px;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;color:#92400e;
  font-weight:600}

/* code blocks */
.df-wiki-content pre{background:#1c1917;color:#e7e5e4;padding:18px 20px;border-radius:12px;
  overflow-x:auto;margin:18px 0;font-size:13px;line-height:1.6}
.df-wiki-content pre code{background:none;padding:0;color:inherit;font-weight:400}

/* horizontal rules — subtle dividers */
.df-wiki-content hr{border:none;height:1px;
  background:linear-gradient(90deg,transparent,#e7dfd2,transparent);margin:32px 0}

/* strong/bold text */
.df-wiki-content strong{color:#1c1917;font-weight:700}

/* ---- wiki sidebar ---- */
.df-wiki-sidebox{background:#fff;border:1px solid #e7dfd2;border-radius:14px;
  padding:18px 20px;margin-bottom:14px;
  box-shadow:0 1px 3px rgba(0,0,0,.03)}
.df-wiki-sidebox h3{font-size:11px;text-transform:uppercase;letter-spacing:.07em;
  color:#a8a29e;margin:0 0 14px;font-weight:700;
  padding-bottom:8px;border-bottom:1px solid #f5f0e8}
.df-wiki-meta-row{display:flex;justify-content:space-between;align-items:center;
  padding:6px 0;border-bottom:1px solid #f5f0e8}
.df-wiki-meta-row:last-child{border-bottom:none}
.df-wiki-meta-row dt{color:#a8a29e;font-weight:500;font-size:12px}
.df-wiki-meta-row dd{margin:0;color:#292524;font-weight:600;font-size:13px}
.df-wiki-meta-row dd a{color:#0369a1;text-decoration:none}
.df-wiki-meta-row dd a:hover{text-decoration:underline}

/* image gallery thumbs */
.df-wiki-thumbs{display:grid;grid-template-columns:1fr 1fr;gap:8px}
.df-wiki-gal{width:100%;height:88px;object-fit:cover;border-radius:10px;cursor:pointer;
  margin:0;box-shadow:none;border:1px solid #f0ebe0;
  transition:transform .14s,box-shadow .14s}
.df-wiki-gal:hover{transform:scale(1.04);box-shadow:0 6px 16px rgba(0,0,0,.14);border-color:#ea580c}

/* related links list */
.df-wiki-links a{display:flex;align-items:center;justify-content:space-between;
  padding:8px 0;color:#292524;text-decoration:none;font-size:13px;font-weight:500;
  border-bottom:1px solid #f5f0e8;transition:padding .1s}
.df-wiki-links a:last-child{border-bottom:none}
.df-wiki-links a:hover{color:#ea580c;padding-left:4px}
.df-wiki-links a span:first-child{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;
  margin-right:8px}
.df-wiki-link-sec{font-size:10px;color:#a8a29e;text-transform:uppercase;
  letter-spacing:.04em;background:#f4efe5;padding:2px 8px;border-radius:10px;
  white-space:nowrap;flex-shrink:0;font-weight:600}

/* lightbox modal */
#df-lightbox{position:fixed;inset:0;background:rgba(15,12,8,.92);z-index:999999;
  display:none;align-items:center;justify-content:center;cursor:zoom-out;
  animation:dfFade .15s ease;backdrop-filter:blur(4px)}
#df-lightbox.open{display:flex}
#df-lightbox img{max-width:88%;max-height:82vh;border-radius:12px;
  box-shadow:0 12px 48px rgba(0,0,0,.5);cursor:default;position:relative;z-index:2}
#df-lightbox-close{position:fixed;top:16px;right:16px;z-index:10;color:#fff;font-size:28px;
  width:44px;height:44px;border-radius:50%;background:rgba(255,255,255,.15);
  backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;
  cursor:pointer;line-height:1;user-select:none;font-weight:300;
  transition:background .15s,transform .15s}
#df-lightbox-close:hover{background:rgba(255,255,255,.3);transform:scale(1.1)}
#df-lightbox-cap{position:fixed;bottom:24px;left:50%;transform:translateX(-50%);
  color:rgba(255,255,255,.8);font-size:14px;max-width:600px;text-align:center;
  line-height:1.4;padding:0 20px;z-index:3}
@keyframes dfFade{from{opacity:0}to{opacity:1}}

/* responsive */
@media(max-width:900px){
  .df-wiki{flex-direction:column;padding:20px 16px;gap:20px}
  .df-wiki-sidebar{width:100%;display:grid;grid-template-columns:1fr 1fr;gap:12px}
  .df-wiki-sidebox{margin-bottom:0}
  .df-wiki-header h1{font-size:24px}
}
@media(max-width:600px){
  .df-wiki-sidebar{grid-template-columns:1fr}
  .df-wiki-content{font-size:15px}
  .df-wiki-content table{font-size:13px}
  .df-wiki-content td,.df-wiki-content th{padding:8px 10px}
}
</style>
`
