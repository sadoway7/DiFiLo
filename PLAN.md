# DiFiLo — Refactor & Modernization Plan

## Goal

Transform DiFiLo from 20 flat `.go` files in `package main` into a layered,
package-based structure that:

1. Groups code by **what kind of code it is** (database, auth, content, UI,
   application wiring).
2. Splits files within each package by **specific responsibility** so nothing
   gets huge (no file exceeds ~250 lines).
3. Makes dependencies flow **one direction only** — no circular imports.
4. Makes the `db`, `textutil`, `auth`, and `ui` packages **reusable**.
5. Leaves clean room for the **wiki CMS** (content creation, maintainers,
   revision history, proposals, tags) without restructuring again.
6. Closes the **infrastructure and security gaps** found during review.

## External Dependencies (from awesome-go review)

Added now:
- **`github.com/stretchr/testify`** — test assertions (Gap 12)
- **`github.com/justinas/nosurf`** — CSRF protection middleware (Gap 4)

Deferred to feature build:
- `gosimple/slug` — URL slugification (wiki page creation)
- `microcosm-cc/bluemonday` — HTML sanitizer (wiki content editor)
- `disintegration/imaging` — image resize/convert (upload pipeline, Gap 13)
- `gabriel-vasile/mimetype` — magic-byte MIME detection (uploads)
- `cyruzin/tome` — pagination helper (Gap 15)
- `lmittmann/tint` — colored slog output (local dev polish)

Rejected (stdlib or own code is better):
- Logging frameworks (zap/zerolog/logrus) — `log/slog` is stdlib
- Config frameworks (viper/koanf) — our 30-line Config struct suffices
- Web frameworks (gin/echo/chi) — switch router is 130 lines
- ORM (gorm/ent/sqlx) — raw SQL against SQLite is clean
- Migration frameworks (goose/golang-migrate) — own 50-line runner is simpler

---

## Principles

- **One package per concern.** Each package owns a single area of responsibility.
- **One file per entity or feature** within a package.
- **Dependencies point one way** — never circular.
- **Leaf packages first** (`textutil`, `auth`, `ui`) depend on nothing.
- **No file exceeds ~250 lines.** Most are 50–150.
- **Create schema tables early** (even unused) so no painful migrations later.
- **Never create empty Go stubs** — they rot. Add feature files when the feature
  is built.
- **`doc.go` in every package** — documents purpose for `go doc`.

---

## Target Structure

Legend: `(NOW)` = exists after refactor, `(NEW)` = new placeholder/scaffolding,
`(FUTURE)` = built later as a feature.

```
DiFiLo/
│
├── cmd/                              ← where the program starts
│   └── difilo/
│       └── main.go                   (NOW) entry point: flags, open DB, start server
│                                        + graceful shutdown (SIGTERM/SIGINT handling)
│
├── internal/                         ← only this project can import these
│   │
│   ├── textutil/                     ← SHARED TEXT HELPERS (leaf — no deps)
│   │   ├── doc.go                    (NEW) "Package textutil provides shared text helpers"
│   │   ├── prose.go                  (NOW) CleanProse, StripTablePipes, DecodeEntities,
│   │   │                                Excerpt, MakeSnippet
│   │   ├── format.go                 (NOW) PrettySection, OrDefault, SectionFromRoute, AZKey
│   │   └── prose_test.go             (NEW) tests for cleanProse, excerpt, makeSnippet
│   │
│   ├── db/                           ← DATABASE LAYER (imports textutil only)
│   │   ├── doc.go                    (NEW) "Package db provides all SQLite operations"
│   │   ├── db.go                     (NOW) DB struct, OpenDB, ALL schema tables
│   │   │                                (includes wiki CMS tables + Tier 2, created idle)
│   │   ├── migrations.go             (NEW) versioned schema migration runner
│   │   ├── users.go                  (NOW) User type + user queries
│   │   ├── comments.go               (NOW) Comment type + comment queries
│   │   ├── bookmarks.go              (NOW) Bookmark type + bookmark queries
│   │   ├── settings.go               (NOW) settings get/set
│   │   ├── downloads.go              (NOW) DownloadLog type + download log queries
│   │   ├── pages.go                  (NOW) ContentPage type + page read queries
│   │   ├── pages_write.go            (FUTURE) InsertPage, UpdatePage, PublishPage
│   │   ├── search.go                 (NOW) SearchHit + FTS5 search + LIKE fallback + Suggest
│   │   │                                + pagination support (offset/perPage params)
│   │   ├── images.go                 (NOW) ImageRow + GetPageImages
│   │   ├── links.go                  (NOW) LinkRow + GetPageLinks + GetInboundLinks
│   │   ├── aliases.go                (NOW) aliasKey, normName, LoadAliases, ResolveRoute, SearchBySlug
│   │   ├── maintainers.go            (FUTURE) page_maintainers queries
│   │   ├── revisions.go              (FUTURE) page_revisions queries
│   │   ├── proposals.go              (FUTURE) edit_proposals queries
│   │   ├── tags.go                   (FUTURE) page_tags queries
│   │   ├── notifications.go          (FUTURE) notifications queries
│   │   ├── uploads.go                (FUTURE) page_uploads queries
│   │   ├── typed_data.go             (FUTURE) material_oxides, recipe_ingredients, etc.
│   │   └── db_test.go                (NEW) test OpenDB creates all tables, basic CRUD
│   │
│   ├── auth/                         ← AUTHENTICATION & PERMISSIONS (leaf — no deps)
│   │   ├── doc.go                    (NEW) "Package auth handles passwords, sessions, roles"
│   │   ├── passwords.go              (NOW) HashPassword, CheckPassword (bcrypt)
│   │   ├── sessions.go               (NOW) CreateSessionToken, ParseSessionToken,
│   │   │                                SetSessionCookie, ClearSessionCookie
│   │   │                                + persistent session secret (stored in DB/file)
│   │   ├── roles.go                  (NOW) RoleAdmin/Manager/General + CanDeleteComment
│   │   │                                + CanEditPage, CanApproveProposal, CanAssignMaintainer
│   │   └── auth_test.go              (NEW) test password hash/verify, session create/parse
│   │
│   ├── content/                      ← CONTENT PROCESSING (imports db + textutil)
│   │   ├── doc.go                    (NEW) "Package content handles import, render, validation"
│   │   ├── import.go                 (NOW) ImportContent orchestrator (upsert, NOT destructive)
│   │   ├── import_pages.go           (NOW) per-page parsing: strip chrome, extract body, upsert
│   │   ├── import_images.go          (NOW) image reference extraction + path resolution + insertion
│   │   ├── import_links.go           (NOW) link graph parsing (the "Related Links" tables)
│   │   ├── import_metadata.go        (NOW) frontmatter parsing, meta-description extraction, byline
│   │   ├── import_helpers.go         (NOW) gunzipIfCompressed, extractThumb, routeOf, resolveImagePath
│   │   ├── markdown.go               (NOW) RenderMarkdown + URL rewriting + empty-section stripping
│   │   ├── wiki.go                   (NOW) RenderWikiPage (standalone func, NOT a Server method)
│   │   ├── wiki_css.go               (NOW) wikiCSS constant (wiki-specific styles)
│   │   ├── lightbox.go               (NOW) lightboxHTML constant (fullscreen image modal + JS)
│   │   ├── validate.go               (FUTURE) content validation for user submissions
│   │   ├── diff.go                   (FUTURE) revision diffing logic
│   │   └── types.go                  (FUTURE) per-type field definitions (material, recipe, etc.)
│   │
│   ├── ui/                           ← UI FRAGMENTS (leaf — Viewer DTO keeps it pure)
│   │   ├── doc.go                    (NEW) "Package ui generates HTML/CSS/JS fragments"
│   │   ├── viewer.go                 (NEW) Viewer struct {LoggedIn, ID, Username, Role}
│   │   ├── css.go                    (NOW) difiCSS: concatenates CSS sub-constants
│   │   ├── css_base.go               (NOW) reset, body, typography, table modernization
│   │   ├── css_nav.go                (NOW) overlay panel, nav, search bar, dice, mobile drawer
│   │   ├── css_pages.go              (NOW) search results, list pages, A-Z bar, external page
│   │   ├── css_home.go               (NOW) hero, explore cards, shimmer, credit, animations
│   │   ├── css_components.go         (NOW) auth forms, comments, admin, pin, bookmarks
│   │   ├── panel.go                  (NOW) PanelHTML (takes *Viewer, not *User)
│   │   ├── shell.go                  (NOW) ShellHTML + acJS autocomplete
│   │   ├── comments.go               (NOW) CommentsHTML (takes *Viewer)
│   │   ├── pin.go                    (NOW) PinButtonHTML (takes *Viewer)
│   │   ├── editor.go                 (FUTURE) content editing forms
│   │   ├── revisions.go              (FUTURE) revision list + diff view
│   │   ├── proposals.go              (FUTURE) proposal forms + review queue
│   │   ├── maintainer_bar.go         (FUTURE) maintainer assign/reassign UI
│   │   ├── tag_chips.go              (FUTURE) clickable tag chips
│   │   ├── widgets_material.go       (FUTURE) oxide analysis table editor widget
│   │   └── notifications.go          (FUTURE) notification bell + dropdown
│   │
│   └── app/                          ← THE APPLICATION (imports all layers)
│       ├── doc.go                    (NEW) "Package app wires all layers into HTTP handlers"
│       ├── server.go                 (NOW) Server struct, renderShell, routing switch
│       ├── startup.go                (NOW) buildHeroImages, buildAliases
│       ├── config.go                 (NEW) Config struct: all tunable values centralized
│       ├── middleware.go             (NEW) requireAuth, requireRole, withLogging, rateLimit
│       ├── logging.go                (NEW) structured logger (log/slog), request logging
│       ├── health.go                 (NEW) /health endpoint: {"status":"ok","pages":N}
│       ├── helpers.go                (NOW) currentUser (calls auth + db), small helpers
│       ├── home.go                   (NOW) homepage handler
│       ├── search.go                 (NOW) search + suggest + explore API handlers
│       ├── list.go                   (NOW) browse list + random + external handlers
│       ├── page.go                   (NOW) page resolution + wiki render dispatch
│       ├── auth.go                   (NOW) register, login, logout handlers + form HTML
│       ├── comments.go               (NOW) comment API: list, post, edit, delete
│       ├── bookmarks.go              (NOW) bookmark API: add, remove, list, check
│       ├── admin.go                  (NOW) admin panel + admin API (roles, settings, downloads)
│       ├── download.go               (NOW) download handler + staticCommentsHTML
│       ├── export.go                 (NOW) inlineAssets, mimeByExt, sanitizeFilename
│       ├── static.go                 (NOW) serveStatic + cache headers (max-age for assets)
│       ├── edit.go                   (FUTURE) create/edit page handlers
│       ├── maintainers.go            (FUTURE) assign/reassign maintainer handlers
│       ├── revisions.go              (FUTURE) view history, diff, revert handlers
│       ├── proposals.go              (FUTURE) propose/approve/reject handlers
│       ├── upload.go                 (FUTURE) image upload handler
│       ├── notifications.go          (FUTURE) notification API
│       └── authors.go                (FUTURE) browse by author/maintainer
│
├── scripts/                          ← launcher & install scripts (out of root)
│   ├── start-mac.command
│   ├── start-windows.bat
│   ├── stop-mac.command
│   ├── stop-windows.bat
│   ├── install-mac.command
│   └── install-windows.bat
│
├── mirror/                           ← data directory
│   ├── md/                           ← markdown source (needed for import)
│   ├── images/                       ← WebP images (served at runtime)
│   ├── media/                        ← videos + media (served at runtime)
│   ├── vendor/                       ← Bootstrap + jQuery (served at runtime)
│   ├── pages.json                    ← page manifest (needed for import)
│   ├── difilo.db                     ← SQLite database (auto-created)
│   └── (html/ excluded from Docker — only needed during import)
│
├── go.mod
├── go.sum
├── Dockerfile                        ← updated: build ./cmd/difilo, exclude mirror/html/
├── .dockerignore
├── .gitignore                        ← updated: binaries, logs, PID, .DS_Store, screenshots
├── .gitlab-ci.yml
├── README.md
└── PLAN.md
```

---

## Dependency Graph

No cycles. Every arrow points one direction.

```
textutil   ← leaf, no deps
     ↑
     ├── db        (imports textutil for MakeSnippet in search fallback)
     └── content   (imports textutil for CleanProse, PrettySection, OrDefault, SectionFromRoute)

auth       ← leaf, no deps (pure crypto + role constants)

ui         ← leaf, no deps (Viewer DTO, not *db.User)

content    ← imports db, textutil

app        ← imports db, auth, content, ui, textutil
               (currentUser lives here: calls auth.ParseSessionToken + db.GetUserByID)

cmd/difilo ← imports app
```

| Package | Depends on | Reusable standalone? |
|---------|-----------|---------------------|
| `textutil` | nothing | Yes — pure functions |
| `auth` | nothing | Yes — pure crypto + constants |
| `ui` | nothing | Yes — Viewer DTO, no db coupling |
| `db` | textutil | Yes (with textutil as trivial leaf) |
| `content` | db, textutil | Yes (with db) |
| `app` | everything | No — it's the wiring layer |

---

## Review Fixes — 6 Cascade Issues Resolved

These issues were found by tracing the call graph before implementation.
Each has a definitive fix.

### Fix 1 — Shared helpers need a `textutil` package

**Problem:** Functions in `helpers.go` are called by three different future
packages. Without a shared package, imports cycle:
`db` needs `makeSnippet` → `makeSnippet` needs `cleanProse` → if `cleanProse`
lives in `app`, then `db` imports `app` → but `app` imports `db` → **circular**.

Cross-boundary callers:
- `MakeSnippet` → `db` (search fallback)
- `CleanProse` → `content` (import pipeline)
- `Excerpt` → `app` (home, search cards)
- `PrettySection` → `app` + `content` (wiki)
- `OrDefault` → `app` + `content` (wiki)
- `SectionFromRoute` → `content` (wiki)
- `AZKey` → `app` (list A-Z bar)

**Fix:** New `internal/textutil/` package (leaf, zero deps) holds all these.
`db`, `content`, and `app` all import it.

### Fix 2 — `renderWikiPage` is a vestigial Server method

**Problem:** `wiki.go:175` defines `func (s *Server) renderWikiPage(...)`, but
the function body **never uses `s`**. It was copy-pasted from when it needed
DB access.

**Fix:** Convert to standalone function `content.RenderWikiPage(p, images, links)`.
Zero behavior change.

### Fix 3 — UI functions take `*User` (resolved: Viewer DTO)

**Problem:** `panelHTML`, `commentsHTML`, `pinButtonHTML` all take `*User`.
If `User` is in `db`, then `ui` must import `db` — breaking "ui is standalone."

**Fix (chosen: Viewer DTO):** Define a lightweight struct in `ui/viewer.go`:
```go
type Viewer struct {
    LoggedIn bool
    ID       int64
    Username string
    Role     string
}
```
The app layer converts `*db.User` → `ui.Viewer` before calling UI functions.
UI stays pure (no db import).

### Fix 4 — `currentUser` is a Server method

**Problem:** `auth.go:106` defines `func (s *Server) currentUser(...)`.
If auth is a separate package, it can't be a Server method.

**Fix:** Auth exports `ParseSessionToken` and `SessionCookieName`. The
`currentUser` function lives in `app/helpers.go` — it reads the cookie,
calls `auth.ParseSessionToken`, then `db.GetUserByID`. Auth stays pure
(no db dependency).

### Fix 5 — `notCapturedHTML` is dead code

**Problem:** `helpers.go:46` — defined but never called anywhere.

**Fix:** Delete it. Do not migrate.

### Fix 6 — Role constants in wrong package

**Problem:** `RoleAdmin`/`Manager`/`General` are in `db.go`. The DB doesn't
use the Go constants (schema uses SQL string `'general'`). They're used by
`auth` (permission checks) and `app` (handlers).

**Fix:** Move to `auth/roles.go`. Add future permission functions:
`CanEditPage`, `CanApproveProposal`, `CanAssignMaintainer`.

---

## Gap Analysis — 17 Issues Found

Found by reviewing the entire codebase for structural holes, missing
infrastructure, and security risks.

### Category 1: Will Cause Data Loss (CRITICAL for wiki)

#### Gap 1 — Destructive import wipes user content

`import.go:62-65` does `DROP TABLE IF EXISTS` on all content tables before
re-importing. Once users create pages, `--reimport` **destroys all their work**.

**Fix:** Replace with upsert-based import (`INSERT OR REPLACE` keyed on route).
Add a `--force` flag for destructive re-import. Non-forced imports update
existing rows, preserving user-created content.

#### Gap 2 — No schema migration system

Only migration is a one-off `ALTER TABLE users ADD COLUMN username` at
`db.go:172`. Every schema change is "DROP + recreate" or ad-hoc ALTER.

**Fix:** `db/migrations.go` with versioned migrations. A `schema_version`
table tracks applied migrations. New migrations are additive (ALTER TABLE,
CREATE TABLE) — never destructive.

### Category 2: Security

#### Gap 3 — Download password hardcoded in client-side JavaScript

`handler_home.go:64` — the password `pancake` is literally in the page source:
```javascript
if(p==='pancake'){window.location.href='...'}
```

**Fix:** Move to server-side check. Store password in DB settings. Add a POST
endpoint that verifies the password server-side, then redirects.

#### Gap 4 — No CSRF protection

All POST endpoints (comments, bookmarks, admin) rely only on `SameSite=Lax`.

**Fix:** CSRF tokens in `app/middleware.go`. Generate per session, embed in
forms/meta tag, verify on POST.

#### Gap 5 — Session secret not persistent

`auth.go:27-29` — `sessionSecret` is random per startup. **Every restart logs
out all users.**

**Fix:** Persist the secret in the DB (`settings` table) or a file. Generate
once on first run, reuse forever.

#### Gap 6 — No rate limiting

No protection against brute-force login, comment spam, or API abuse.

**Fix:** Basic IP-based rate limiting middleware in `app/middleware.go`.

### Category 3: Missing Infrastructure

#### Gap 7 — No logging framework

Everything uses `fmt.Fprintf(os.Stderr, ...)` — 7 scattered call sites. No
levels, no structure, no request logging.

**Fix:** `app/logging.go` with `log/slog` (Go stdlib since 1.21). Wire through
middleware for per-request logging.

#### Gap 8 — No middleware layer

Every handler manually calls `s.currentUser(r)` and checks roles — 15+
repetitions. With wiki adding 10+ handlers, this explodes.

**Fix:** `app/middleware.go`:
```go
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc
func (s *Server) requireRole(role string, next http.HandlerFunc) http.HandlerFunc
func (s *Server) withLogging(next http.HandlerFunc) http.HandlerFunc
```

#### Gap 9 — No configuration system

Hardcoded values scattered: `pancake`, port 8000, 30-day sessions, 100-result
search limit, 5000-char comments, 8-image gallery, 25-link sidebar, 80KB hero
threshold.

**Fix:** `app/config.go` with a `Config` struct centralizing all tunables.

#### Gap 10 — No graceful shutdown

`main.go:65` uses `http.ListenAndServe` — no signal handling. Docker sends
SIGTERM and the server dies hard, risking SQLite corruption.

**Fix:** `http.Server` with `Shutdown()` + `signal.Notify(SIGTERM, SIGINT)` in
`cmd/difilo/main.go`.

#### Gap 11 — No health check endpoint

No `/health`. Docker/CI can't verify the server is actually serving.

**Fix:** `app/health.go` — returns `200 {"status":"ok","pages":N}`.

#### Gap 12 — No tests

Zero `_test.go` files.

**Fix (create during refactor):**
- `textutil/prose_test.go` — test CleanProse, Excerpt, MakeSnippet
- `auth/auth_test.go` — test password hash/verify, session create/parse
- `db/db_test.go` — test OpenDB creates all tables, basic CRUD

### Category 4: Wiki CMS Needs

#### Gap 13 — No image upload pipeline

Zero `multipart`/`FormFile` references. Users creating wiki pages need to
upload images.

**Schema placeholder NOW** (in `db/db.go`): `page_uploads` table.
**Future files:** `app/upload.go`, `db/uploads.go`, `ui/uploader.go`.

#### Gap 14 — No notification system

No way to tell a user "your proposal was approved" or "someone commented on
your page."

**Schema placeholder NOW:** `notifications` table.
**Future files:** `db/notifications.go`, `app/notifications.go`, `ui/notifications.go`.

#### Gap 15 — No pagination

All queries use fixed `LIMIT`, no offset. Search caps at 100. List pages
render everything (2,843 materials = 2,843 `<li>` elements).

**Fix:** Add `offset` + `perPage` params to `db/pages.go` and `db/search.go`.
Add page navigation in `ui/`.

### Category 5: Performance / Ops

#### Gap 16 — No cache headers on static assets

`static.go` serves 3,375 images with no explicit cache headers. Browsers
re-request every image on every page load.

**Fix:** `Cache-Control: public, max-age=31536000` for `/images/`, `/vendor/`,
`/media/` in `app/static.go`. Files never change (content-addressed).

#### Gap 17 — mirror/html/ is 351 MB dead weight at runtime

The 11,403 HTML files are only read during import (to extract meta
descriptions). At runtime, pages render from DB.

**Fix:** Exclude `mirror/html/` from Docker image. Keep in repo for re-import
capability.

---

## Wiki CMS Roadmap

The app is evolving from a read-only archive into a **living wiki** where users
create and manage content.

### Content lifecycle

1. **Ingest** — content enters (via import OR user creation)
2. **Store** — written to SQLite (`pages` + typed data tables)
3. **Review** — (optional) proposal workflow
4. **Publish** — becomes visible
5. **Render** — displayed to users
6. **Revise** — edited, new version saved

### Maintainer system

- **On creation:** user submits new page → becomes its maintainer automatically
- **Admin reassign:** admin opens page, selects new maintainer from user list
- **Multiple maintainers:** a page can have several (many-to-many table)
- **Permissions by role:**

| Action | Admin | Manager | General (maintainer) | General (other) |
|--------|-------|---------|---------------------|-----------------|
| Edit own page | direct | direct | direct | — |
| Edit others' page | direct | direct | propose | propose |
| Reassign maintainer | yes | no | no | no |
| Approve proposal | yes | yes | no | no |
| Revert to revision | yes | yes | no | no |

### Feature → file map

| Feature | db/ | app/ | ui/ | content/ |
|---------|-----|------|-----|----------|
| Create new page | `pages_write.go` | `edit.go` | `editor.go` | `validate.go` |
| Edit existing page | `pages_write.go` | `edit.go` | `editor.go` | `validate.go` |
| Maintainer assignment | `maintainers.go` | `maintainers.go` | `maintainer_bar.go` | — |
| Revision history | `revisions.go` | `revisions.go` | `revisions.go` | `diff.go` |
| Edit proposals | `proposals.go` | `proposals.go` | `proposals.go` | — |
| Tags / typecodes | `tags.go` | `tags.go` | `tag_chips.go` | — |
| Structured data per type | `typed_data.go` | — | `widgets_material.go` | `types.go` |
| Draft/published | (existing `status`) | `edit.go` | `editor.go` | — |
| Image upload | `uploads.go` | `upload.go` | (drag-drop in editor) | — |
| Notifications | `notifications.go` | `notifications.go` | `notifications.go` | — |
| Author/maintainer browse | (existing `pages.go`) | `authors.go` | `author_card.go` | — |

---

## Schema Tables (all in db/db.go)

Created during refactor, even if unused. Tables sit idle until code uses them.
All use `CREATE TABLE IF NOT EXISTS` — zero risk.

### Existing tables (NOW)

`users`, `comments`, `bookmarks`, `settings`, `downloads`, `pages`, `pages_fts`,
`page_images`, `page_links`

### New placeholder tables (created during refactor)

```sql
-- Migration tracking
CREATE TABLE IF NOT EXISTS schema_version (
    version    INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Page revision history (every edit saved)
CREATE TABLE IF NOT EXISTS page_revisions (
    id           INTEGER PRIMARY KEY,
    page_id      INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    body_md      TEXT NOT NULL,
    edited_by    INTEGER REFERENCES users(id),
    edit_summary TEXT DEFAULT '',
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_revisions_page ON page_revisions(page_id, created_at DESC);

-- Edit proposals (propose → review → publish)
CREATE TABLE IF NOT EXISTS edit_proposals (
    id               INTEGER PRIMARY KEY,
    page_id          INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    proposed_body_md TEXT NOT NULL,
    proposed_by      INTEGER NOT NULL REFERENCES users(id),
    status           TEXT DEFAULT 'pending',  -- pending/approved/rejected
    reviewed_by      INTEGER REFERENCES users(id),
    review_note      TEXT DEFAULT '',
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    reviewed_at      DATETIME
);
CREATE INDEX IF NOT EXISTS idx_proposals_status ON edit_proposals(status);

-- Page maintainers (many-to-many)
CREATE TABLE IF NOT EXISTS page_maintainers (
    page_id INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(page_id, user_id)
);

-- Tags (typecodes as browsable tags)
CREATE TABLE IF NOT EXISTS page_tags (
    id      INTEGER PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    tag     TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_tags_tag ON page_tags(tag);
CREATE INDEX IF NOT EXISTS idx_tags_page ON page_tags(page_id);

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       TEXT NOT NULL,       -- 'comment', 'proposal', 'maintainer'
    message    TEXT NOT NULL,
    route      TEXT DEFAULT '',
    read       INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, read);

-- User image uploads
CREATE TABLE IF NOT EXISTS page_uploads (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id),
    page_id    INTEGER REFERENCES pages(id) ON DELETE SET NULL,
    filename   TEXT NOT NULL,
    image_path TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tier 2: Material oxide analysis
CREATE TABLE IF NOT EXISTS material_oxides (
    id           INTEGER PRIMARY KEY,
    page_id      INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    oxide        TEXT NOT NULL,
    oxide_url    TEXT DEFAULT '',
    analysis_pct TEXT DEFAULT '',
    formula      TEXT DEFAULT '',
    tolerance    TEXT DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_mat_ox_page ON material_oxides(page_id);

-- Tier 2: Recipe ingredients
CREATE TABLE IF NOT EXISTS recipe_ingredients (
    id            INTEGER PRIMARY KEY,
    page_id       INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    material_name TEXT NOT NULL,
    material_url  TEXT DEFAULT '',
    amount        TEXT DEFAULT '',
    units         TEXT DEFAULT '',
    percent       TEXT DEFAULT '',
    sort_order    INTEGER DEFAULT 0
);

-- Tier 2: Firing schedule steps
CREATE TABLE IF NOT EXISTS schedule_steps (
    id          INTEGER PRIMARY KEY,
    page_id     INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    step_num    INTEGER NOT NULL,
    ramp_c      TEXT DEFAULT '',
    ramp_f      TEXT DEFAULT '',
    hold        TEXT DEFAULT '',
    time        TEXT DEFAULT '',
    description TEXT DEFAULT ''
);
```

---

## Placeholders to Create During Refactor

### Schema tables (in `db/db.go`)

All tables listed in the Schema section above — created with
`CREATE TABLE IF NOT EXISTS`, idle until code uses them.

### Package documentation (`doc.go` per package, 7 files)

```go
// Package db provides all SQLite database operations for DiFiLo:
// users, comments, bookmarks, content pages, full-text search, and
// the wiki lifecycle (revisions, proposals, maintainers, tags).
package db
```

One per package: `textutil`, `db`, `auth`, `content`, `ui`, `app`.

### Infrastructure files (in `app/`)

| File | Purpose |
|------|---------|
| `app/config.go` | `Config` struct centralizing all tunables (ports, limits, thresholds) |
| `app/middleware.go` | `requireAuth`, `requireRole`, `withLogging`, `rateLimit` |
| `app/logging.go` | Structured logger (`log/slog`), request logging |
| `app/health.go` | `/health` endpoint: `{"status":"ok","pages":N}` |

### Migration system (`db/migrations.go`)

Versioned migration runner. Reads `schema_version` table, applies pending
migrations in order. Migrations are additive — never destructive.

### Tests (3 files)

| File | Tests |
|------|-------|
| `textutil/prose_test.go` | CleanProse, Excerpt, MakeSnippet, PrettySection |
| `auth/auth_test.go` | HashPassword/CheckPassword, session create/parse |
| `db/db_test.go` | OpenDB creates all tables, basic CRUD |

### NOT created (rot)

- Empty Go function stubs for future features
- `.gitkeep` files (every directory has real content)
- Empty test files for packages without testable logic yet

---

## File Migration Map

Every current file maps to one or more new files. Nothing is lost; dead code
is deleted.

### Entry point

| Current | New | Notes |
|---------|-----|-------|
| `main.go` | `cmd/difilo/main.go` | Move. Add graceful shutdown (SIGTERM/SIGINT). |

### textutil/ (shared helpers — NEW package)

| Current | New | Notes |
|---------|-----|-------|
| `helpers.go` (cleanProse, stripTablePipes, decodeEntities, excerpt, makeSnippet) | `textutil/prose.go` | Move + export. |
| `helpers.go` (prettySection, orDefault, sectionFromRoute, azKey) | `textutil/format.go` | Move + export. |

### db/ (all SQLite code)

| Current | New | Notes |
|---------|-----|-------|
| `db.go` (DB struct, OpenDB, schema) | `db/db.go` | Keep + add all new schema tables. |
| `db.go` (User type + user queries) | `db/users.go` | Extract. |
| `db.go` (Comment type + comment queries) | `db/comments.go` | Extract. |
| `db.go` (Bookmark type + bookmark queries) | `db/bookmarks.go` | Extract. |
| `db.go` (settings get/set) | `db/settings.go` | Extract. |
| `db.go` (DownloadLog type + download queries) | `db/downloads.go` | Extract. |
| `content_db.go` (ContentPage + page queries) | `db/pages.go` | Extract. |
| `content_db.go` (SearchHit + FTS5 + fallback + Suggest) | `db/search.go` | Extract. Uses textutil.MakeSnippet. |
| `content_db.go` (ImageRow + GetPageImages) | `db/images.go` | Extract. |
| `content_db.go` (LinkRow + GetPageLinks + GetInboundLinks) | `db/links.go` | Extract. |
| `content_db.go` (LoadAliases, ResolveRoute, SearchBySlug) | `db/aliases.go` | Extract. |
| `aliases.go` (aliasKey, normName) | `db/aliases.go` | Merge. |
| (new) | `db/migrations.go` | NEW: versioned migration runner. |

### auth/ (pure crypto + roles — leaf package)

| Current | New | Notes |
|---------|-----|-------|
| `auth.go` (hashPassword, checkPassword) | `auth/passwords.go` | Move + export. |
| `auth.go` (createSessionToken, parseSessionToken, cookies) | `auth/sessions.go` | Move + export. Add persistent secret. |
| `auth.go` (canDeleteComment) + `db.go` (role constants) | `auth/roles.go` | Merge. Add future permission funcs. |

**Note:** `auth.go`'s `currentUser` does NOT move to auth. It moves to
`app/helpers.go` (Fix 4). Auth is a pure leaf package.

### content/ (import + rendering)

| Current | New | Notes |
|---------|-----|-------|
| `import.go` (ImportContent orchestrator) | `content/import.go` | Change to upsert (Fix for Gap 1). |
| `import.go` (per-page parsing) | `content/import_pages.go` | Extract. |
| `import.go` (image extraction) | `content/import_images.go` | Extract. |
| `import.go` (link graph parsing) | `content/import_links.go` | Extract. |
| `import.go` (frontmatter, meta-desc, byline) | `content/import_metadata.go` | Extract. |
| `content_helpers.go` | `content/import_helpers.go` | Rename + move. |
| `markdown.go` | `content/markdown.go` | Move + export. |
| `wiki.go` (renderWikiPage + helpers) | `content/wiki.go` | Convert to standalone func (Fix 2). |
| `wiki.go` (wikiCSS constant) | `content/wiki_css.go` | Extract. |
| `wiki.go` (lightboxHTML constant) | `content/lightbox.go` | Extract. |

### ui/ (HTML/CSS/JS fragments — leaf package via Viewer DTO)

| Current | New | Notes |
|---------|-----|-------|
| (new) | `ui/viewer.go` | NEW: Viewer struct (Fix 3). |
| `assets.go` (difiCSS) | `ui/css.go` | Assembler concatenating sub-constants. |
| `assets.go` (base CSS) | `ui/css_base.go` | Extract. |
| `assets.go` (nav CSS) | `ui/css_nav.go` | Extract. |
| `assets.go` (pages CSS) | `ui/css_pages.go` | Extract. |
| `assets.go` (home CSS) | `ui/css_home.go` | Extract. |
| `assets.go` (components CSS) | `ui/css_components.go` | Extract. |
| `assets.go` (panelHTML, navGroups, bookmarkJS) | `ui/panel.go` | Extract. Takes *Viewer. |
| `assets.go` (shellHTML, acJS) | `ui/shell.go` | Extract. |
| `assets.go` (commentsHTML) | `ui/comments.go` | Extract. Takes *Viewer. |
| `assets.go` (pinButtonHTML) | `ui/pin.go` | Extract. Takes *Viewer. |

### app/ (Server struct + all HTTP handlers)

| Current | New | Notes |
|---------|-----|-------|
| `server.go` | `app/server.go` | Server struct, renderShell, routing. |
| `server.go` (startup helpers) | `app/startup.go` | buildHeroImages, buildAliases. |
| `auth.go` (currentUser) | `app/helpers.go` | Move here (Fix 4). Calls auth + db. |
| (new) | `app/config.go` | NEW: Config struct (Gap 9). |
| (new) | `app/middleware.go` | NEW: auth/role/logging middleware (Gap 8). |
| (new) | `app/logging.go` | NEW: structured logger (Gap 7). |
| (new) | `app/health.go` | NEW: /health endpoint (Gap 11). |
| `handler_home.go` | `app/home.go` | Rename. Fix password (Gap 3). |
| `handler_search.go` | `app/search.go` | Rename. |
| `handler_list.go` | `app/list.go` | Rename. |
| `handler_page.go` | `app/page.go` | Rename. |
| `handlers_auth.go` (register/login/logout) | `app/auth.go` | Split. |
| `handlers_auth.go` (comment API) | `app/comments.go` | Split. |
| `handlers_auth.go` (bookmark API) | `app/bookmarks.go` | Split. |
| `handlers_auth.go` (admin panel + API) | `app/admin.go` | Split. |
| `download.go` (handleDownload + staticCommentsHTML) | `app/download.go` | Split. |
| `download.go` (inlineAssets, mimeByExt, sanitize) | `app/export.go` | Split. |
| `static.go` | `app/static.go` | Move. Add cache headers (Gap 16). |

### Deleted (dead code)

| Current | Reason |
|---------|--------|
| `rewrite.go` (446 lines) | `Rewriter` struct never instantiated. Replaced by DB-driven rendering. |
| `helpers.go` (`notCapturedHTML`) | Defined but never called. |
| `mirror/.difilo-index.bin` (12 MB) | Old gob search index, zero references in code. FTS5 replaced it. |

### Moved (not code)

| Current | New |
|---------|-----|
| `Start DIFI-LOCAL (Mac).command` | `scripts/start-mac.command` |
| `Start DIFI-LOCAL (Windows).bat` | `scripts/start-windows.bat` |
| `Stop DIFI-LOCAL (Mac).command` | `scripts/stop-mac.command` |
| `Stop DIFI-LOCAL (Windows).bat` | `scripts/stop-windows.bat` |
| `Install Dependencies (Mac).command` | `scripts/install-mac.command` |
| `Install Dependencies (Windows).bat` | `scripts/install-windows.bat` |

### Removed from repo (added to .gitignore)

| Current | Reason |
|---------|--------|
| `DiFiLo` (macOS binary) | Build artifact. |
| `DiFiLo.exe` (Windows binary) | Build artifact. |
| `difilo.log` | Runtime log. |
| `difilo.pid` | Runtime PID file. |
| `.DS_Store` | macOS metadata. |
| `Screenshot *.png` (3 files) | Not part of codebase. |

---

## Implementation Steps

### Step 1 — Cleanup (no code changes, zero risk)

- [ ] Delete `rewrite.go` (dead code)
- [ ] Delete `mirror/.difilo-index.bin` (dead gob index)
- [ ] Move launcher scripts to `scripts/` with clean names
- [ ] Remove binaries, logs, PID, screenshots, `.DS_Store` from repo
- [ ] Update `.gitignore`
- [ ] Commit: "cleanup: remove dead code, move scripts, ignore build artifacts"

### Step 2 — Create package directories + doc.go files

- [ ] `mkdir -p cmd/difilo internal/{textutil,db,auth,content,ui,app} scripts`
- [ ] Create `doc.go` in each package (7 files)
- [ ] Commit: "scaffold: create package directories and doc.go files"

### Step 3 — Build textutil package (leaf, no deps)

- [ ] Create `textutil/prose.go` — CleanProse, StripTablePipes, DecodeEntities, Excerpt, MakeSnippet
- [ ] Create `textutil/format.go` — PrettySection, OrDefault, SectionFromRoute, AZKey
- [ ] Create `textutil/prose_test.go` — basic tests
- [ ] Verify: `go build ./internal/textutil && go test ./internal/textutil`
- [ ] Commit: "refactor: extract textutil package (shared text helpers)"

### Step 4 — Build db package (imports textutil)

- [ ] Create `db/db.go` — DB struct, OpenDB, ALL schema (existing + new placeholder tables)
- [ ] Create `db/migrations.go` — versioned migration runner
- [ ] Create `db/users.go` through `db/aliases.go` (11 files)
- [ ] Create `db/db_test.go` — test schema creation, basic CRUD
- [ ] Verify: `go build ./internal/db && go test ./internal/db`
- [ ] Commit: "refactor: extract db package (all SQLite + schema + migrations)"

### Step 5 — Build auth package (leaf, no deps)

- [ ] Create `auth/passwords.go` — HashPassword, CheckPassword
- [ ] Create `auth/sessions.go` — token create/parse, cookie set/clear, persistent secret
- [ ] Create `auth/roles.go` — role constants + CanDeleteComment
- [ ] Create `auth/auth_test.go` — test hash/verify, session create/parse
- [ ] Verify: `go build ./internal/auth && go test ./internal/auth`
- [ ] Commit: "refactor: extract auth package (passwords, sessions, roles)"

### Step 6 — Build content package (imports db + textutil)

- [ ] Create `content/import.go` through `content/import_helpers.go` (6 files)
- [ ] Create `content/markdown.go`, `content/wiki.go`, `content/wiki_css.go`, `content/lightbox.go`
- [ ] Convert renderWikiPage to standalone func (Fix 2)
- [ ] Change import to upsert, not destructive (Gap 1 fix)
- [ ] Verify: `go build ./internal/content`
- [ ] Commit: "refactor: extract content package (import + render)"

### Step 7 — Build ui package (leaf, Viewer DTO)

- [ ] Create `ui/viewer.go` — Viewer struct (Fix 3)
- [ ] Create `ui/css_*.go` (6 CSS files) + `ui/css.go` (assembler)
- [ ] Create `ui/panel.go`, `ui/shell.go`, `ui/comments.go`, `ui/pin.go`
- [ ] Change all functions to take `*Viewer` instead of `*User`
- [ ] Verify: `go build ./internal/ui`
- [ ] Commit: "refactor: extract ui package (CSS, panel, shell, Viewer DTO)"

### Step 8 — Build app package (imports all layers)

- [ ] Create `app/config.go`, `app/middleware.go`, `app/logging.go`, `app/health.go`
- [ ] Create `app/server.go`, `app/startup.go`, `app/helpers.go` (includes currentUser)
- [ ] Create `app/home.go` through `app/static.go` (15 handler files)
- [ ] Update all cross-package references (db.User, auth.HashPassword, etc.)
- [ ] Fix Gap 3 (password server-side), Gap 16 (cache headers)
- [ ] Verify: `go build ./internal/app`
- [ ] Commit: "refactor: extract app package (Server, handlers, middleware)"

### Step 9 — Wire entry point

- [ ] Create `cmd/difilo/main.go` — flags, open DB, import, build server, listen
- [ ] Add graceful shutdown (Gap 10)
- [ ] Verify: `go build ./cmd/difilo`
- [ ] Commit: "refactor: wire entry point with graceful shutdown"

### Step 10 — Delete old root files + verify

- [ ] Remove all old `.go` files from root
- [ ] Delete `helpers.go` `notCapturedHTML` (Fix 5)
- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` clean
- [ ] `go test ./...` passes
- [ ] Run server, test all features
- [ ] Commit: "refactor: remove old source files, restructure complete"

### Step 11 — Update Dockerfile + CI

- [ ] Update Dockerfile: `go build -o /difilo ./cmd/difilo`
- [ ] Exclude `mirror/html/` from Docker (Gap 17)
- [ ] Update `.gitlab-ci.yml`
- [ ] Verify: `docker build -t difilo .`
- [ ] Commit: "ci: update for new structure, exclude mirror/html from Docker"

### Step 12 — Update README

- [ ] Update repository contents table
- [ ] Update build instructions: `go build -o DiFiLo ./cmd/difilo`
- [ ] Document package layers
- [ ] Commit: "docs: update README for new file structure"

---

## Verification Checklist

After the refactor:

- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` clean
- [ ] `go test ./...` passes
- [ ] Server starts: `go run ./cmd/difilo --mirror ./mirror --port 8000`
- [ ] Homepage loads
- [ ] Search works (results + snippets)
- [ ] Browse works (section list + A-Z + filter)
- [ ] Page rendering works (wiki layout + sidebar + lightbox)
- [ ] Auth works (register, login, logout)
- [ ] Comments work (post, edit, delete)
- [ ] Bookmarks work (pin, dropdown)
- [ ] Admin panel works (users, roles, settings)
- [ ] Download works (self-contained HTML export)
- [ ] `/health` returns 200 OK
- [ ] Docker build works
- [ ] No `.go` files in root
- [ ] No build artifacts in repo
- [ ] Session survives restart (persistent secret)
- [ ] Static assets have cache headers

---

## Post-Refactor Priorities

### Immediately after refactor (safety)

1. Fix destructive import (upsert, not DROP) — Gap 1
2. Move password server-side — Gap 3
3. Persistent session secret — Gap 5
4. Graceful shutdown — Gap 10

### Soon after (hardening)

5. CSRF tokens — Gap 4
6. Structured logging — Gap 7
7. Rate limiting — Gap 6
8. Expand tests — Gap 12
9. Pagination — Gap 15

### When building wiki features

10. Image upload pipeline — Gap 13
11. Notification system — Gap 14
12. Content editor (create/edit forms)
13. Maintainer assignment
14. Revision history + diff
15. Edit proposal workflow
16. Tags
17. Structured data widgets (oxide tables, recipe builders)

---

## Expected File Sizes

No file exceeds ~250 lines. Most are 50–150.

| File | Approx lines | Notes |
|------|-------------|-------|
| `db/db.go` | ~250 | Schema is long (all tables) — unavoidable |
| `app/admin.go` | ~200 | Admin panel has 3 sections |
| `content/import_pages.go` | ~200 | Page parsing is detailed |
| `app/export.go` | ~170 | Base64 inlining logic |
| `ui/css_nav.go` | ~150 | Nav panel is the most complex CSS |
| `content/import.go` | ~130 | Orchestrator loop |
| `app/server.go` | ~130 | Routing switch |
| `ui/panel.go` | ~120 | Nav bar assembly |
| `content/wiki.go` | ~120 | Two-column layout assembly |
| `app/home.go` | ~115 | Homepage assembly |
| `db/pages.go` | ~110 | Several query functions |
| `db/search.go` | ~100 | FTS5 + fallback + suggest |
| Everything else | 20–90 | Focused and small |

---

## What This Fixes

| Problem now | After restructure |
|-------------|-------------------|
| `assets.go` is 945 lines mixing CSS, HTML, JS, 5 components | 11 files under `ui/` — 6 CSS + 4 HTML/JS + viewer.go |
| `import.go` is 698 lines doing everything | 5 files under `content/` — orchestrator + pages + images + links + metadata |
| `handlers_auth.go` is 585 lines mixing 4 features | 4 files under `app/` — auth + comments + bookmarks + admin |
| `db.go` is 467 lines with 6 entity types | 8+ files under `db/` — one per entity |
| `content_db.go` is 509 lines mixing 5 concerns | 5 files under `db/` — pages + search + images + links + aliases |
| `rewrite.go` is 446 lines of dead code | Deleted |
| `notCapturedHTML` is dead code | Deleted |
| `.difilo-index.bin` is 12 MB of dead data | Deleted |
| 20 `.go` files in root | 0 `.go` files in root |
| 6 launcher scripts in root | Moved to `scripts/` |
| Binaries, logs, screenshots in repo | Removed, gitignored |
| Nothing reusable — all `package main` | 4 leaf packages reusable (textutil, auth, ui, db) |
| Shared helpers cause circular imports | `textutil` package breaks the cycle |
| `renderWikiPage` is vestigial Server method | Standalone `content.RenderWikiPage()` |
| UI coupled to db via `*User` | `ui.Viewer` DTO — ui stays pure |
| `currentUser` trapped in Server method | Lives in app, calls exported auth functions |
| No schema for wiki features | All tables created idle, ready for code |
| No migrations | Versioned migration runner in `db/migrations.go` |
| No middleware | `app/middleware.go` with auth/role/logging |
| No config | `app/config.go` centralizes all tunables |
| No tests | 3 test files covering the leaf packages |
| No health check | `app/health.go` |
| No graceful shutdown | `cmd/difilo/main.go` handles SIGTERM/SIGINT |
| Destructive import wipes user content | Upsert-based import preserves user data |
| Password in client-side JS | Server-side verification |
| Session lost on restart | Persistent secret in DB |
