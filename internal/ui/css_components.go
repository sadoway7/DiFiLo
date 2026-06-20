package ui

// cssComponents holds the styles for reusable UI widgets that are shared across
// pages: the nav user-area and role badge, the login/register forms, the
// comments section and its edit affordances, the admin user/role table, the
// in-content pin/download toolbar, the bookmarks dropdown, and the homepage
// recent-comments list.
const cssComponents = `
/* ---- user area in nav ---- */
.difi-user{display:flex;align-items:center;gap:6px;margin-left:8px;font-size:13px;white-space:nowrap;
  flex-shrink:0;flex-grow:0}
.difi-user a{color:#3f3a34;text-decoration:none;padding:4px 8px;border-radius:8px;font-weight:600}
.difi-user a:hover{background:#efe6d8}
.difi-uname{color:#6b7280;font-size:12px;max-width:160px;overflow:hidden;text-overflow:ellipsis}
.difi-badge-role{font-size:10px;font-weight:700;text-transform:uppercase;background:#f4efe5;color:#57534e;
  border-radius:999px;padding:1px 7px;letter-spacing:.03em}
.difi-admin-link{background:#ea580c!important;color:#fff!important}
.difi-admin-link:hover{background:#c2410c!important}
@media(max-width:860px){.difi-user{display:none!important}}

/* ---- auth forms ---- */
.df-auth{max-width:400px;margin:40px auto;padding:0 18px;font-family:system-ui,sans-serif}
.df-auth h1{font-size:24px;margin-bottom:18px}
.df-auth-form{display:flex;flex-direction:column;gap:14px}
.df-auth-form label{display:flex;flex-direction:column;gap:4px;font-size:14px;font-weight:600;color:#374151}
.df-auth-form input{height:42px;border:1px solid #d1d5db;border-radius:9px;padding:0 14px;font-size:15px;outline:none}
.df-auth-form input:focus{border-color:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.15)}
.df-auth-form button{height:44px;border:0;background:#ea580c;color:#fff;border-radius:9px;font-size:15px;
  font-weight:700;cursor:pointer;margin-top:4px}
.df-auth-form button:hover{background:#c2410c}
.df-err{color:#dc2626;font-size:14px;margin-top:12px;padding:10px 14px;background:#fef2f2;border-radius:8px;
  border:1px solid #fecaca}

/* ---- comments ---- */
#df-comments{max-width:940px;margin:40px auto;padding:0 18px 60px;font-family:system-ui,sans-serif;
  border-top:1px solid #e5e7eb;padding-top:30px}
#df-comments h2{font-size:18px;margin-bottom:16px}
.df-comment-list{list-style:none;padding:0;margin:0}
.df-comment{padding:14px 0;border-bottom:1px solid #f0f0f0}
.df-comment-head{display:flex;align-items:center;gap:8px;flex-wrap:wrap;margin-bottom:6px}
.df-comment-author{font-weight:600;font-size:14px;color:#1f2937}
.df-comment-role{font-size:10px;font-weight:700;text-transform:uppercase;background:#f4efe5;color:#57534e;
  border-radius:999px;padding:1px 7px;letter-spacing:.03em}
.df-comment-date{font-size:12px;color:#9ca3af;margin-left:auto}
.df-comment-del{font-size:11px;color:#dc2626;background:none;border:1px solid #fecaca;border-radius:6px;
  padding:2px 8px;cursor:pointer}
.df-comment-del:hover{background:#fef2f2}
.df-comment-edit{font-size:11px;color:#2563eb;background:none;border:1px solid #bfdbfe;border-radius:6px;
  padding:2px 8px;cursor:pointer}
.df-comment-edit:hover{background:#eff6ff}
.df-edit-input{width:100%;border:1px solid #d1d5db;border-radius:8px;padding:8px 12px;font-size:14px;
  font-family:inherit;outline:none;resize:vertical;box-sizing:border-box}
.df-edit-input:focus{border-color:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.15)}
.df-edit-actions{margin-top:6px;display:flex;gap:6px}
.df-edit-save{height:32px;border:0;background:#ea580c;color:#fff;border-radius:7px;padding:0 14px;
  font-size:13px;font-weight:600;cursor:pointer}
.df-edit-save:hover{background:#c2410c}
.df-edit-cancel{height:32px;border:1px solid #d1d5db;background:#fff;color:#6b7280;border-radius:7px;
  padding:0 14px;font-size:13px;cursor:pointer}
.df-edit-cancel:hover{background:#f5f5f5}
.df-comment-body{font-size:14px;color:#374151;line-height:1.6;word-break:break-word}
#df-comment-form{margin-top:18px}
#df-comment-input{width:100%;border:1px solid #d1d5db;border-radius:9px;padding:10px 14px;font-size:14px;
  font-family:inherit;outline:none;resize:vertical;box-sizing:border-box}
#df-comment-input:focus{border-color:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.15)}
#df-comment-form button{margin-top:8px;height:38px;border:0;background:#ea580c;color:#fff;border-radius:9px;
  padding:0 20px;font-size:14px;font-weight:700;cursor:pointer}
#df-comment-form button:hover{background:#c2410c}

/* ---- admin ---- */
.df-admin{max-width:940px;margin:30px auto;padding:0 18px;font-family:system-ui,sans-serif}
.df-admin table{width:100%;border-collapse:collapse;font-size:14px;margin-top:12px}
.df-admin th{text-align:left;padding:10px 12px;background:#f4efe5;color:#57534e;font-size:12px;
  text-transform:uppercase;letter-spacing:.04em;border-bottom:2px solid #e0d6c5}
.df-admin td{padding:8px 12px;border-bottom:1px solid #eee6d8}
.df-role-select{height:32px;border:1px solid #d1d5db;border-radius:6px;font-size:13px;padding:0 8px}
.df-del-user{font-size:12px;color:#dc2626;background:none;border:1px solid #fecaca;border-radius:6px;
  padding:4px 10px;cursor:pointer}
.df-del-user:hover{background:#fef2f2}
.df-user-dl{font-size:12px;font-weight:600;color:#b45309;background:#fff7ed;border:1px solid #fed7aa;border-radius:6px;
  padding:4px 10px;text-decoration:none;cursor:pointer;display:inline-block}
.df-user-dl:hover{background:#fed7aa;color:#9a3412}
.df-setting-row{display:flex;gap:8px;align-items:center;margin-top:12px}
.df-setting-row label{font-weight:600;font-size:14px;white-space:nowrap}
.df-setting-row input{height:38px;border:1px solid #d1d5db;border-radius:8px;padding:0 12px;font-size:14px;outline:none}
.df-setting-row input:focus{border-color:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.15)}
.df-setting-row button{height:38px;border:0;background:#ea580c;color:#fff;border-radius:8px;padding:0 16px;
  font-size:14px;font-weight:600;cursor:pointer;white-space:nowrap}
.df-setting-row button:hover{background:#c2410c}

/* ---- pin/bookmark button (inlaid top-right of content) ---- */
#df-pin-bar{display:flex;justify-content:flex-end;align-items:center;gap:8px;padding:8px 18px 0}
#df-dl-btn{display:inline-flex;align-items:center;gap:6px;border:1.5px solid #e7dfd2;background:#fff;
  color:#1f2937;padding:5px 14px;border-radius:8px;font-size:13px;font-weight:700;text-decoration:none;
  cursor:pointer;transition:.12s}
#df-dl-btn:hover{background:#fdf2e3;border-color:#ea580c;color:#c2410c}
#df-pin-btn{display:inline-flex;align-items:center;gap:6px;border:1.5px solid #ea580c;
  background:#fff7ed;color:#c2410c;padding:5px 14px;border-radius:8px;cursor:pointer;
  font-size:13px;font-weight:700;transition:.12s;box-shadow:0 1px 3px rgba(234,88,12,.15)}
#df-pin-btn:hover{background:#fed7aa;border-color:#c2410c;color:#9a3412}
#df-pin-btn.pinned{background:#ea580c;border-color:#ea580c;color:#fff;box-shadow:none}
#df-pin-btn.pinned svg{fill:#fff}

/* ---- bookmarks dropdown ---- */
.difi-bookmarks-group summary{color:#57534e}
.difi-bookmarks-menu{max-height:400px;overflow-y:auto}
.difi-bookmarks-menu .df-muted{padding:12px;font-size:13px}

/* ---- recent comments (home) ---- */
.df-recent-comments{list-style:none;padding:0;margin:0}
.df-recent-comments li{padding:14px 0;border-bottom:1px solid #f0f0f0}
.df-rc-meta{display:flex;align-items:center;flex-wrap:wrap;gap:6px;font-size:13px;line-height:1.4;color:#6b7280}
.df-rc-author{font-weight:600;color:#1f2937}
.df-rc-badge{font-size:10px;font-weight:700;text-transform:uppercase;background:#f4efe5;color:#57534e;
  border-radius:999px;padding:1px 7px;letter-spacing:.03em}
.df-rc-on{color:#9ca3af}
.df-rc-page{color:#007bff;text-decoration:none;font-weight:500}
.df-rc-page:hover{text-decoration:underline}
.df-rc-snippet{font-size:14px;color:#374151;margin-top:6px;line-height:1.6;word-break:break-word}`
