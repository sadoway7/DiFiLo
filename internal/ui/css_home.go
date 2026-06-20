package ui

// cssHome styles the homepage hero: the full-bleed hero wrapper, big search
// row with random-page dice, the random "Explore" picks grid and its shimmer
// loading placeholders, the homepage entrance animations, the credit/byline
// block, and the search autocomplete dropdown.
const cssHome = `
/* home */
/* Homepage hero tucks up under the fixed navbar, so drop the global body offset here. */
body.df-homepage{padding-top:0!important}
.df-home{max-width:940px;margin:0 auto;padding:0 18px;font-family:system-ui,sans-serif;color:#1f2937;
  min-height:calc(100vh - 104px);display:flex;flex-direction:column}
.df-disclaimer{color:#9ca3af;font-size:12.5px;line-height:1.6;text-align:center;margin:48px auto 0;max-width:780px}
.df-home h2{font-size:15px;text-transform:uppercase;letter-spacing:.06em;color:#b45309;margin:30px 0 10px;font-weight:800}
.df-home h2::before{content:"";display:inline-block;width:22px;height:4px;margin-right:10px;vertical-align:middle;background:linear-gradient(90deg,#fbbf24,#f97316,#ef4444);border-radius:2px}
.df-explore-head{display:flex;align-items:center;justify-content:flex-start;margin:30px 0 10px}
.df-explore-refresh{border:none;background:transparent;width:auto;height:auto;
  display:inline-flex;align-items:center;gap:6px;color:#1c1917;cursor:pointer;
  font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:.05em;
  padding:4px 8px;border-radius:8px;transition:.15s}
.df-explore-refresh:hover{color:#ea580c;background:#fef3e2}
.df-explore-refresh svg{transition:transform .4s}
.df-explore-refresh:hover svg{transform:rotate(180deg)}
.df-herowrap{position:relative;z-index:2;text-align:center;margin:0 0 6px;
  width:100vw;margin-left:calc(50% - 50vw);
  margin-top:40px;
  padding:170px 18px 80px;
  background-color:#faf7f1;
  background-size:cover;background-position:center;background-repeat:no-repeat}
.df-herowrap::before{content:"";position:absolute;inset:0;z-index:0;
  background:radial-gradient(ellipse 62% 70% at 50% 45%,rgba(255,255,255,.72),rgba(255,255,255,.18) 42%,rgba(255,255,255,0) 70%)}
.df-herowrap>*{position:relative;z-index:1}
.df-herowrap .df-searchrow{z-index:3}
.df-searchrow{display:flex;align-items:center;justify-content:center;gap:12px;
  margin:0 auto 20px;max-width:800px;width:100%;box-sizing:border-box}
.df-herowrap .df-bigsearch{margin:0;flex:1 1 auto;max-width:none;min-width:0;box-sizing:border-box}
.df-searchrow .df-dice{flex:none}
/* mobile + tablet: search bar fits screen, dice hidden on homepage */
@media(max-width:768px){
  .df-searchrow{gap:8px;padding:0;width:100%;max-width:100%}
  .df-searchrow .df-bigsearch{min-width:0;width:100%;height:44px;flex:1 1 100%}
  .df-searchrow .df-bigsearch input{height:44px;font-size:16px;padding-left:16px;min-width:0}
  .df-searchrow .df-bigsearch button{height:44px;padding:0 14px;font-size:14px;flex-shrink:0}
  body.df-homepage .df-searchrow .df-dice{display:none}
  .df-hero-login{font-size:11px;top:16px;right:16px;padding:5px 10px}
}
/* hero mobile fixes */
@media(max-width:768px){
  /* override the full-bleed trick — on mobile the container is already full width */
  .df-herowrap{width:100%;margin-left:0;margin-top:0;padding:60px 12px 40px}
  .df-hero-login{display:none}
  .df-download{margin:10px 0 0}
  .df-download a{font-size:12px;padding:7px 14px}
  .df-download-hint{font-size:11px}
  .df-btns{gap:6px;margin:16px auto 0}
  .df-btn{font-size:12px;padding:6px 10px}
  .df-dice{width:40px;height:40px}
}
.df-dice{position:relative;width:46px;height:46px;border-radius:12px;flex:none;text-decoration:none;cursor:pointer;
  border:1px solid #e4dac6;-webkit-tap-highlight-color:transparent;
  background:
    radial-gradient(circle 4px at 32% 32%,#b45309 98%,transparent 100%),
    radial-gradient(circle 4px at 68% 32%,#b45309 98%,transparent 100%),
    radial-gradient(circle 4px at 50% 50%,#b45309 98%,transparent 100%),
    radial-gradient(circle 4px at 32% 68%,#b45309 98%,transparent 100%),
    radial-gradient(circle 4px at 68% 68%,#b45309 98%,transparent 100%),
    linear-gradient(145deg,#fff,#f3ece0);
  box-shadow:0 4px 0 #dccfb6,0 7px 12px rgba(70,45,20,.22);transition:transform .08s,box-shadow .08s}
.df-dice:hover{transform:translateY(-1px);box-shadow:0 5px 0 #dccfb6,0 9px 15px rgba(70,45,20,.26)}
.df-dice:active{transform:translateY(3px);box-shadow:0 1px 0 #dccfb6,0 2px 5px rgba(70,45,20,.18)}
.df-herowrap .df-tagline{margin:0 auto 18px}
.df-download{margin:14px 0 0;font-size:14px}
.df-download a{display:inline-block;background:#fff;color:#ea580c;font-weight:700;
  text-decoration:none;padding:8px 18px;border-radius:24px;
  box-shadow:0 2px 8px rgba(0,0,0,.12);transition:.15s}
.df-download a:hover{background:#ea580c;color:#fff;text-decoration:none;
  box-shadow:0 4px 12px rgba(234,88,12,.35)}
.df-download-sub{margin:4px 0 0;font-size:12px;color:#9ca3af;font-style:italic}
.df-download-hint{margin:8px 0 0;font-size:12.5px;color:#57534e;
  display:inline-block;background:rgba(255,255,255,.55);
  backdrop-filter:blur(6px);-webkit-backdrop-filter:blur(6px);
  padding:5px 14px;border-radius:20px}
.df-hero-login{position:absolute;top:24px;right:56px;z-index:2;
  font-size:12.5px;color:#44403c;
  background:rgba(255,255,255,.55);backdrop-filter:blur(6px);-webkit-backdrop-filter:blur(6px);
  padding:6px 14px;border-radius:20px}
.df-orig-site{display:inline-block;color:#007bff;text-decoration:none;font-size:14px;font-weight:600}
.df-orig-site:hover{color:#0056b3;text-decoration:underline}
.df-herowrap .df-disclaimer{margin:18px auto 0}
.df-tagline{color:#6b7280;font-size:16px;margin:2px 0 18px}
.df-btns{display:flex;flex-wrap:wrap;justify-content:center;gap:10px;max-width:880px;margin:24px auto 0}
.df-btn{display:inline-flex;align-items:center;gap:7px;padding:7px 12px;border-radius:9px;
  background:#fff;border:1px solid #e7dfd2;color:#1f2937;text-decoration:none;font-weight:600;font-size:13px;
  box-shadow:0 1px 2px rgba(0,0,0,.04);transition:.12s}
.df-btn:hover{background:#ea580c;border-color:#ea580c;color:#fff;box-shadow:0 2px 6px rgba(234,88,12,.3)}
.df-btn .df-num{font-size:11px;color:#9ca3af;font-weight:500}
.df-btn:hover .df-num{color:#ffe4d1}
/* random "Explore" picks grid */
.df-picks{display:grid;grid-template-columns:repeat(4,1fr);gap:16px;margin:8px 0 10px}

/* shimmer loading cards */
.df-shimmer{border-radius:14px;overflow:hidden;background:#fff;border:1px solid #e7dfd2}
.df-shimmer-img{width:100%;aspect-ratio:16/9;background:linear-gradient(90deg,#f0ebe0 25%,#f8f4ed 50%,#f0ebe0 75%);
  background-size:200% 100%;animation:dfShimmer 1.2s infinite}
.df-shimmer-body{padding:11px 13px 13px;display:flex;flex-direction:column;gap:6px;min-height:96px}
.df-shimmer-line{height:12px;border-radius:6px;
  background:linear-gradient(90deg,#f0ebe0 25%,#f8f4ed 50%,#f0ebe0 75%);
  background-size:200% 100%;animation:dfShimmer 1.2s infinite}
.df-shimmer-line.s1{width:80%}.df-shimmer-line.s2{width:50%;height:9px}.df-shimmer-line.s3{width:100%;height:9px}
@keyframes dfShimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}
@keyframes dfCardIn{from{opacity:0;transform:translateY(8px)}to{opacity:1;transform:translateY(0)}}

/* ---- homepage entrance animations ---- */
/* snappy, subtle stagger — each element type gets slightly different treatment */
@keyframes dfHeroBg{from{opacity:0}to{opacity:1}}
@keyframes dfHeroItem{from{opacity:0;transform:translateY(8px)}to{opacity:1;transform:translateY(0)}}
@keyframes dfHeroFade{from{opacity:0}to{opacity:1}}

.df-herowrap{animation:dfHeroBg .4s ease-out 0s both}
.df-searchrow{animation:dfHeroItem .3s ease-out .08s both}
.df-download{animation:dfHeroFade .25s ease-out .15s both}
.df-hero-login{animation:dfHeroFade .25s ease-out .2s both}
.df-btns{animation:dfHeroFade .25s ease-out .22s both}
.df-explore-head{animation:dfHeroItem .3s ease-out .28s both}
.df-picks{animation:dfHeroItem .35s ease-out .34s both}
.df-credit{animation:dfHeroFade .3s ease-out .42s both}
.df-home>h2:not(.df-explore-head h2),.df-recent-comments{animation:dfHeroFade .25s ease-out .38s both}

@media(prefers-reduced-motion:reduce){
  .df-herowrap,.df-searchrow,.df-download,.df-hero-login,
  .df-btns,.df-explore-head,.df-picks,.df-credit{animation:none!important}
}
@media(min-width:1100px){.df-picks{width:calc(100% + 240px);margin-left:-120px;margin-right:-120px;gap:20px}}
.df-pick{display:flex;flex-direction:column;background:#fff;border:1px solid #e7dfd2;border-radius:14px;
  overflow:hidden;text-decoration:none;color:#1f2937;box-shadow:0 1px 2px rgba(0,0,0,.04);transition:.18s}
.df-pick:hover{transform:translateY(-4px) scale(1.02);border-color:#ea580c;
  box-shadow:0 12px 28px rgba(234,88,12,.16);text-decoration:none}
.df-pick-img{width:100%;aspect-ratio:16/9;background-size:cover;background-position:center;background-color:#f3eee2}
.df-pick-body{padding:11px 13px 13px;display:flex;flex-direction:column;gap:5px;min-height:96px}
.df-pick-title{font-size:14.5px;font-weight:700;color:#1f2937;line-height:1.3;overflow:hidden;text-overflow:ellipsis;
  display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical}
.df-pick-sec{font-size:10.5px;font-weight:700;text-transform:uppercase;letter-spacing:.05em;color:#b45309}
.df-pick-text{font-size:12.5px;color:#6b7280;line-height:1.45;overflow:hidden;text-overflow:ellipsis;
  display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical}
@media(max-width:860px){.df-picks{grid-template-columns:repeat(2,1fr)}.df-pick-extra{display:none}}
@media(max-width:480px){.df-picks{grid-template-columns:1fr}}
.df-credit{margin:34px 0 6px;padding-top:22px;border-top:1px solid #eef0f3}
.df-credit-row{display:flex;gap:16px;align-items:flex-start}
.df-tony{width:96px;height:96px;border-radius:50%;object-fit:cover;flex:0 0 auto;background:#f3f4f6}
.df-byline{margin:0 0 12px;color:#374151;font-size:14.5px;line-height:1.55}
.df-byline b{color:#b45309}
.df-socials{display:flex;flex-wrap:wrap;gap:10px;align-items:center}
.df-socials a{display:inline-block;line-height:0;border-radius:50%;transition:.12s}
.df-socials a:hover{transform:translateY(-2px)}
.df-socials img{width:34px;height:34px;border-radius:50%;object-fit:cover;display:block;background:#f3f4f6}
.df-coffee{display:inline-flex;align-items:center;gap:8px;background:#ff5f5f;color:#fff!important;text-decoration:none;
  font-weight:700;font-size:13px;padding:9px 16px;border-radius:999px;transition:.12s}
.df-coffee:hover{background:#e04b4b;transform:translateY(-1px)}
.df-coffee img{width:20px;height:20px;border-radius:0;background:transparent}
.df-credit-btns{display:flex;gap:10px;flex-wrap:wrap;margin-top:14px}
.df-orig-site{display:inline-flex;align-items:center;gap:6px;background:#007bff;color:#fff!important;text-decoration:none;
  font-weight:700;font-size:13px;padding:9px 16px;border-radius:999px;transition:.12s}
.df-orig-site:hover{background:#0056b3;transform:translateY(-1px)}
.df-url{display:inline-block;background:#f3f4f6;border:1px solid #e5e7eb;border-radius:8px;padding:8px 12px;
  font-size:14px;word-break:break-all;max-width:100%}
.df-bigsearch{display:flex;gap:0;margin:8px 0 18px;height:46px;border:none;
  border-radius:999px;overflow:visible;position:relative;background:#fff;
  box-shadow:0 4px 14px rgba(15,23,42,.08);transition:transform .18s,box-shadow .18s}
.df-bigsearch:hover{transform:translateY(-1px);box-shadow:0 8px 24px rgba(15,23,42,.14)}
.df-bigsearch:focus-within{transform:translateY(-1px);
  box-shadow:0 10px 28px rgba(15,23,42,.16),0 0 0 3px rgba(234,88,12,.18)}
.df-bigsearch input{flex:1 1 0;width:0;height:46px;min-width:0;border:0;background:#fff;color:#0f172a;
  padding:0 4px 0 22px;font-size:16px;outline:none;border-radius:999px 0 0 999px}
.df-bigsearch input::placeholder{color:#a8a29e}
.df-bigsearch button{height:46px;border:0;background:#ea580c;color:#fff;
  padding:0 18px;font-size:15px;font-weight:700;cursor:pointer;display:inline-flex;align-items:center;justify-content:center;border-radius:0 999px 999px 0;flex-shrink:0}
.df-bigsearch button:hover{background:#c2410c}
/* search autocomplete dropdown */
.df-ac{position:absolute;left:0;right:0;top:calc(100% + 10px);background:#fff;border:1px solid #e7dfd2;
  border-radius:16px;box-shadow:0 12px 30px rgba(15,23,42,.18);overflow:hidden;z-index:80;display:none;text-align:left;max-width:100%}
.df-ac.open{display:block}
.df-ac-head{padding:9px 18px 4px;font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:.07em;color:#b6a98e}
.df-ac-item{display:flex;align-items:baseline;gap:8px;padding:9px 18px;cursor:pointer;font-size:14px;color:#374151;
  border-bottom:1px solid #f3ede3;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.df-ac-item:last-child{border-bottom:0}
.df-ac-item .df-ac-t{overflow:hidden;text-overflow:ellipsis}
.df-ac-item .df-ac-sec{margin-left:auto;flex:none;font-size:11px;color:#b6a98e;text-transform:uppercase;letter-spacing:.04em}
.df-ac-item.sel{background:#fdf2e3;color:#c2410c}
.df-bigsearch input::selection{background:rgba(234,88,12,.22)}`
