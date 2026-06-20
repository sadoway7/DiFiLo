# AGENTS.md — DiFiLo Developer Guide

This document explains the file structure, conventions, and Go best practices
used in DiFiLo. Read this before adding code.

---

## Quick Reference

```bash
go build -o DiFiLo ./cmd/difilo          # build the binary
go run ./cmd/difilo --mirror ./mirror    # run without building
go test ./...                            # run all tests
go vet ./...                             # lint
```

---

## Project Overview

DiFiLo is a self-hosted, offline-first reference library. A single Go binary
serves ~11,400 archived web pages from a local SQLite database with full-text
search, user accounts, comments, bookmarks, and role-based moderation. No
external services, no phone-home, no network calls at runtime.

---

## Package Structure

The project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
conventions. Code is organized by **what kind of code it is** — each package
owns a single concern.

```
DiFiLo/
├── cmd/difilo/          Entry point. Parses flags, opens DB, starts server.
│                        NO business logic lives here.
│
├── internal/            Private application code. Go enforces: only this
│   │                    module can import these packages.
│   │
│   ├── textutil/        SHARED TEXT HELPERS (leaf package — zero deps)
│   │                    CleanProse, Excerpt, MakeSnippet, PrettySection,
│   │                    OrDefault, SectionFromRoute, AZKey.
│   │                    Used by db (search fallback), content (rendering),
│   │                    and app (handlers).
│   │
│   ├── db/              DATABASE LAYER (imports textutil only)
│   │                    DB struct, OpenDB, schema, all SQLite queries.
│   │                    One file per entity: users.go, comments.go,
│   │                    bookmarks.go, pages.go, search.go, etc.
│   │                    Includes placeholder tables for future wiki features
│   │                    (revisions, proposals, maintainers, tags, etc.).
│   │
│   ├── auth/            AUTHENTICATION (leaf package — zero deps)
│   │                    bcrypt passwords, HMAC session tokens, cookies,
│   │                    role constants, permission helpers.
│   │                    Uses interfaces (UserRole, CommentOwner) so it never
│   │                    imports db types. Fully reusable standalone.
│   │
│   ├── content/         CONTENT PROCESSING (imports db + textutil)
│   │                    Import pipeline (mirror files → SQLite),
│   │                    markdown rendering (Goldmark), wiki page layout,
│   │                    image lightbox.
│   │                    ImportContent is non-destructive: uses DELETE on
│   │                    imported pages (source_url != ''), preserving
│   │                    user-created content.
│   │
│   ├── ui/              UI FRAGMENTS (leaf package — zero deps)
│   │                    CSS, HTML, JS string constants and builders.
│   │                    Uses a Viewer DTO (not *db.User) so the package
│   │                    has zero coupling to the database layer.
│   │
│   └── app/             APPLICATION WIRING (imports all packages above)
│                        Server struct, HTTP routing, all request handlers,
│                        middleware (auth, CSRF, logging), config, health
│                        check. This is the ONLY package that touches all
│                        layers — that's its job.
│
├── start/             Double-click launcher scripts (Mac/Windows).
├── mirror/              Data directory (md/, html/, images/, media/, pages.json).
├── go.mod / go.sum      Module definition and dependency locks.
└── Dockerfile           Container build (builds from ./cmd/difilo).
```

---

## Dependency Graph

Dependencies flow **one direction only**. No cycles. Ever.

```
cmd/difilo  →  app  →  db, auth, content, ui, textutil
                       content  →  db, textutil
                       db       →  textutil
                       auth     →  (nothing)
                       ui       →  (nothing)
                       textutil →  (nothing)
```

**Leaf packages** (auth, ui, textutil) depend on nothing — they can be copied
into another project as-is. The `db` package depends on textutil only for
snippet generation in the search fallback path.

If you find yourself wanting to add an import that would create a cycle,
you've put code in the wrong package. Stop and restructure.

---

## Conventions

### One file per entity

Inside each package, files are split by entity or feature:
- `db/users.go` — User type + all user queries
- `db/comments.go` — Comment type + all comment queries
- `db/search.go` — SearchHit type + FTS5 search + fallback

Don't put two unrelated entities in the same file.

### Export everything that crosses package boundaries

Functions called from another package must be exported (capitalized first
letter). Internal helpers within a package stay unexported (lowercase).

### UI functions take *Viewer, not *db.User

The `ui` package is a pure leaf with zero dependencies. It cannot import `db`.
Instead of passing `*db.User`, pass a `*ui.Viewer`:

```go
type Viewer struct {
    LoggedIn bool
    ID       int64
    Username string
    Email    string
    Role     string
}
```

The app layer converts before calling UI functions:
```go
viewer := &ui.Viewer{LoggedIn: true, ID: u.ID, Username: u.Username, Role: u.Role}
html := ui.PanelHTML(route, viewer)
```

### Auth uses interfaces for permission checks

`auth.CanDeleteComment` takes interfaces, not concrete types:
```go
auth.CanDeleteComment(userRoleAdapter, commentOwnerAdapter)
```
This keeps auth decoupled from db. The app layer provides adapter types.

### Handlers are methods on Server

All HTTP handlers are methods on `*app.Server`:
```go
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request)
```

The Server struct holds the database, config, and cached indexes.

### Schema tables are created eagerly

All tables — including placeholders for future features — are created in
`OpenDB()` via `CREATE TABLE IF NOT EXISTS`. This avoids migration headaches
when features are built later. The tables sit idle until code uses them.

### No dead code

If a function or file isn't called anywhere, delete it. Don't leave commented-out
code or unused helpers. The codebase should be clean enough to navigate by grep.

---

## How to Add New Features

| Feature type | Where it goes | Example |
|---|---|---|
| New database table | `internal/db/db.go` schema (CREATE TABLE IF NOT EXISTS) | `page_revisions` |
| New DB query | `internal/db/<entity>.go` | `GetRevisionsByPage(pageID)` in `revisions.go` |
| New HTTP endpoint | `internal/app/<feature>.go` + add route in `server.go` | `handleRevisions` in `revisions.go` |
| New UI component | `internal/ui/<component>.go` | `RevisionList` in `revisions.go` |
| New content rendering | `internal/content/<feature>.go` | `RenderDiff` in `diff.go` |
| New shared helper | `internal/textutil/<file>.go` | If used by 2+ packages |
| New test | Same package as the code, `<name>_test.go` | `db/revisions_test.go` |

---

## Build Commands

```bash
# Build binary
go build -o DiFiLo ./cmd/difilo

# Build for Windows
GOOS=windows go build -o DiFiLo.exe ./cmd/difilo

# Run directly
go run ./cmd/difilo --mirror ./mirror --port 8000

# Test everything
go test ./...

# Test one package
go test ./internal/db -v

# Vet
go vet ./...
```

---

## Key Files to Know

| File | What to read it for |
|---|---|
| `cmd/difilo/main.go` | Startup flow: flags → DB → import → server → shutdown |
| `internal/app/server.go` | Server struct + routing switch (all URL routes) |
| `internal/app/config.go` | All tunable values (limits, thresholds) |
| `internal/db/db.go` | Full database schema — all 19 tables |
| `internal/content/import.go` | How content enters the system from mirror files |
| `internal/ui/viewer.go` | The Viewer DTO pattern |
| `PLAN.md` | Full refactor plan + gap analysis + wiki CMS roadmap |

---

## Dependencies (External)

Minimal by design. The app philosophy is "single binary, no phone-home."

| Dependency | Purpose |
|---|---|
| `modernc.org/sqlite` | Pure-Go SQLite driver (no CGO) |
| `golang.org/x/crypto/bcrypt` | Password hashing |
| `github.com/yuin/goldmark` | Markdown rendering (GFM) |
| `github.com/stretchr/testify` | Test assertions |
| `github.com/justinas/nosurf` | CSRF middleware (available, currently using built-in) |

No web framework. No ORM. No config library. No logging framework (uses stdlib
`log/slog`). The routing is a `switch` statement. This is intentional — every
dependency is a liability for an offline tool.
