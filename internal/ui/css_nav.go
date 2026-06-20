package ui

// cssNav is the fixed overlay panel: the Digitalfire top bar, back button,
// brand/logo, grouped expanding navbar, search box, dice, hamburger, the
// mobile slide-in drawer, and all of their interaction styling.
const cssNav = `
/* ---- overlay panel (warm cream; logo sits directly) ---- */
body{padding-top:74px!important;margin-top:0!important;overflow-x:clip}
#difi-bar{position:fixed;top:0;left:0;right:0;height:58px;display:flex;align-items:center;gap:3px;
  background:#faf7f1;color:#3f3a34;padding:0 12px;z-index:2147483647;border-bottom:1px solid #e7dfd2;
  box-shadow:0 2px 10px rgba(70,45,20,.10);font:14px -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;box-sizing:border-box}
#difi-bar *{box-sizing:border-box}
#difi-bar .difi-ic{height:42px;min-width:42px;padding:0 8px;border:0;border-radius:9px;background:transparent;
  color:#8b8278;font-size:24px;line-height:1;cursor:pointer;display:inline-flex;align-items:center;justify-content:center;transition:.12s}
#difi-bar .difi-ic:hover{background:#efe6d8;color:#1c1917}
#difi-bar .difi-brand{display:flex;align-items:center;text-decoration:none;margin:0 6px 0 2px;padding:4px 8px;border-radius:9px}
#difi-bar .difi-brand img{height:34px;width:auto;display:block;object-fit:contain}
#difi-bar .difi-brand:hover{background:#efe6d8}

/* grouped, expanding sub-menu navbar */
#difi-bar .difi-nav{display:flex;align-items:center;gap:2px;margin-left:4px}
#difi-bar .difi-group{position:relative}
#difi-bar .difi-group summary{list-style:none;cursor:pointer;height:42px;padding:0 14px;border-radius:9px;
  display:inline-flex;align-items:center;color:#44403c;font-weight:600;font-size:14px;white-space:nowrap}
#difi-bar .difi-group summary::-webkit-details-marker{display:none}
#difi-bar .difi-group summary::after{content:"\25BE";font-size:10px;margin-left:7px;opacity:.55}
#difi-bar .difi-group summary:hover,#difi-bar .difi-group[open] summary{background:#efe6d8;color:#1c1917}
#difi-bar .difi-submenu{position:absolute;top:100%;left:0;background:#fff;border:1px solid #e7dfd2;
  border-radius:12px;padding:8px;min-width:210px;box-shadow:0 14px 34px rgba(70,45,20,.20);z-index:1}
#difi-bar .difi-submenu a{display:block;color:#44403c;text-decoration:none;padding:9px 12px;border-radius:8px;
  font-weight:500;font-size:14px;white-space:nowrap}
#difi-bar .difi-submenu a:hover{background:#f4ecdf;color:#1c1917}

/* search: white pill input + flame-orange button */
#difi-bar .difi-search{display:flex;flex:1;min-width:0;max-width:560px;margin-left:auto;height:42px}
/* spacer pushes bookmarks/user right on home/search pages (where search bar is hidden) */
.difi-nav-spacer{flex:0}
body.df-homepage .difi-nav-spacer,body.df-searchpage .difi-nav-spacer{flex:1}
#difi-bar .difi-search input{flex:1;height:42px;min-width:0;border:0;background:#fff;color:#0f172a;
  border-radius:999px 0 0 999px;padding:0 4px 0 18px;font-size:14px;outline:none;box-shadow:inset 0 0 0 1px #e7dfd2}
#difi-bar .difi-search input::placeholder{color:#a8a29e}
#difi-bar .difi-search input:focus{box-shadow:inset 0 0 0 2px #ea580c}
#difi-bar .difi-search button{height:42px;border:0;background:#ea580c;color:#fff;
  border-radius:0 999px 999px 0;padding:0 13px;font-size:13px;font-weight:700;cursor:pointer;
  display:inline-flex;align-items:center;justify-content:center}
#difi-bar .difi-search button:hover{background:#c2410c}
/* dice beside the navbar search (desktop) */
.difi-dice{position:relative;width:34px;height:34px;border-radius:9px;flex:none;margin-left:6px;
  text-decoration:none;cursor:pointer;border:1px solid #e4dac6;-webkit-tap-highlight-color:transparent;
  background:
    radial-gradient(circle 3px at 32% 32%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 68% 32%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 50% 50%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 32% 68%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 68% 68%,#b45309 98%,transparent 100%),
    linear-gradient(145deg,#fff,#f3ece0);
  box-shadow:0 3px 0 #dccfb6,0 4px 8px rgba(70,45,20,.18);transition:transform .08s,box-shadow .08s}
.difi-dice:hover{transform:translateY(-1px);box-shadow:0 4px 0 #dccfb6,0 5px 10px rgba(70,45,20,.22)}
.difi-dice:active{transform:translateY(2px);box-shadow:0 1px 0 #dccfb6,0 2px 4px rgba(70,45,20,.15)}

/* mobile dice (beside the hamburger, shown only on mobile) */
.difi-mdice{display:none;width:34px;height:34px;border-radius:9px;flex:none;margin-left:auto;
  align-items:center;justify-content:center;text-decoration:none;cursor:pointer;
  border:1px solid #e4dac6;-webkit-tap-highlight-color:transparent;
  background:
    radial-gradient(circle 3px at 32% 32%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 68% 32%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 50% 50%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 32% 68%,#b45309 98%,transparent 100%),
    radial-gradient(circle 3px at 68% 68%,#b45309 98%,transparent 100%),
    linear-gradient(145deg,#fff,#f3ece0);
  box-shadow:0 3px 0 #dccfb6,0 4px 8px rgba(70,45,20,.18);transition:transform .08s,box-shadow .08s}
.difi-mdice:active{transform:translateY(2px);box-shadow:0 1px 0 #dccfb6,0 2px 4px rgba(70,45,20,.15)}

/* hamburger (mobile only) */
#difi-bar .difi-burger{display:none;height:44px;width:44px;border:1px solid #e7dfd2;background:#fff;border-radius:11px;
  color:#1c1917;line-height:1;cursor:pointer;align-items:center;justify-content:center;margin-left:auto;
  touch-action:manipulation;box-shadow:0 1px 2px rgba(70,45,20,.06)}
#difi-bar .difi-burger:hover{background:#efe6d8;border-color:#dcd1bf}
#difi-bar .difi-burger svg{display:block}
/* backdrop + side slide-in panel (mobile) */
.difi-backdrop{position:fixed;inset:0;background:rgba(28,25,23,.45);z-index:2147483646;
  opacity:0;visibility:hidden;transition:opacity .2s,visibility .2s}
.difi-backdrop.open{opacity:1;visibility:visible}
#difi-drawer{position:fixed;top:0;right:0;bottom:0;width:min(380px,92vw);background:#faf7f1;z-index:2147483647;
  transform:translateX(100%);transition:transform .24s ease;display:flex;flex-direction:column;
  box-shadow:-12px 0 30px rgba(70,45,20,.22);border-left:1px solid #e7dfd2;
  font:14px -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif}
#difi-drawer.open{transform:translateX(0)}
#difi-drawer *{box-sizing:border-box}
/* drawer head: search + close on one row (no title) */
.difi-dhead{display:flex;align-items:center;gap:8px;padding:12px 12px;border-bottom:1px solid #ece3d4}
.difi-dclose{height:44px;width:44px;flex:0 0 auto;border:0;background:transparent;border-radius:10px;color:#57534e;
  line-height:1;cursor:pointer;display:inline-flex;align-items:center;justify-content:center;
  touch-action:manipulation;-webkit-tap-highlight-color:transparent}
.difi-dclose svg{display:block}
.difi-dclose:hover{background:#efe6d8}
/* drawer search (lives in the head row) */
.difi-dhead .difi-search{flex:1;display:flex;height:44px;min-width:0;border-radius:999px;overflow:hidden;
  background:#fff;box-shadow:inset 0 0 0 1px #e7dfd2;transition:box-shadow .12s}
.difi-dhead .difi-search:focus-within{box-shadow:inset 0 0 0 2px #ea580c}
.difi-dhead .difi-search input{flex:1;min-width:0;border:0;background:transparent;color:#0f172a;padding:0 8px 0 16px;
  font-size:16px;outline:none;-webkit-appearance:none;appearance:none;-webkit-tap-highlight-color:transparent}
.difi-dhead .difi-search input::placeholder{color:#a8a29e}
.difi-dhead .difi-search button{border:0;background:#ea580c;color:#fff;padding:0 18px;font-size:14px;font-weight:700;
  cursor:pointer;-webkit-tap-highlight-color:transparent;touch-action:manipulation}
.difi-dhead .difi-search button:hover{background:#c2410c}
/* scrollable body — styled track signals it scrolls */
.difi-dbody{padding:6px 12px 18px;overflow-y:auto;-webkit-overflow-scrolling:touch;flex:1}
.difi-dbody::-webkit-scrollbar{width:8px}
.difi-dbody::-webkit-scrollbar-thumb{background:#d6cab3;border-radius:8px}
.difi-dbody::-webkit-scrollbar-track{background:transparent}
/* accordion groups mirroring the desktop nav (collapsed = short menu) */
.difi-dgroup{margin:0;border-bottom:1px solid #efe7d8}
.difi-dgroup:last-child{border-bottom:0}
.difi-dgroup>summary{list-style:none;cursor:pointer;display:flex;align-items:center;justify-content:space-between;
  padding:15px 2px;font-size:16px;font-weight:700;color:#1c1917;touch-action:manipulation;-webkit-tap-highlight-color:transparent}
.difi-dgroup>summary::-webkit-details-marker{display:none}
.difi-dgroup>summary::after{content:"";width:8px;height:8px;display:inline-block;
  border-right:2px solid #a8a29e;border-bottom:2px solid #a8a29e;transform:rotate(-45deg);transition:transform .15s}
.difi-dgroup[open]>summary::after{transform:rotate(45deg)}
.difi-dgroup[open]>summary{color:#b45309}
.difi-dlinks{display:flex;flex-direction:column;gap:5px;padding:4px 0 12px}
.difi-dlinks a{display:block;width:100%;text-align:left;text-decoration:none;color:#3f3a34;font-size:16px;font-weight:500;
  padding:12px 14px;border-radius:10px;background:#fff;border:1px solid #ece3d4;
  touch-action:manipulation;-webkit-tap-highlight-color:transparent}
.difi-dlinks a:hover,.difi-dlinks a:active{background:#efe6d8;border-color:#dcd1bf;color:#1c1917}`
