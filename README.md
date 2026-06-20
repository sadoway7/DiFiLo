# DigitalFire-Offline

An offline archive of **Digitalfire** — the ceramic materials, glaze chemistry,
and reference site created by **Tony Hansen** over 35+ years of research.

This data was gathered from the **Internet Archive (Wayback Machine)** and
preserved here so the information survives. The Wayback Machine is not a
guaranteed permanent record — archived pages can be removed when issues arise.
**All content is Tony Hansen's work; original authorship and rights remain with him.**

> This data was **not** scraped from digitalfire.com directly. Strict rate
> limiting was followed and the Internet Archive's guidelines were respected at
> every step — no duplicate data was gathered over and over, connections were
> not hammered. This was done carefully and respectfully.

---

## Quick start

Runs entirely on your own machine — nothing phones home.

**Mac** — double-click `Start DIFI-LOCAL (Mac).command`
**Windows** — double-click `Start DIFI-LOCAL (Windows).bat`

Your browser opens automatically to **http://localhost:8000/** once the server
is ready. (The very first launch builds a search index, which can take a minute
or two.)

**To stop:** double-click the matching `Stop DIFI-LOCAL (...)` script.

## Using the archive

- **Search** — the search bar at the top of every page searches all ~11,400
  pages. Results are ranked by term frequency, with thumbnails and snippets.
- **Browse** — the top menu groups the site's sections (Materials; Recipes &
  Firing; Learn; Media & More). Click a section for its full list, then use the
  filter box to narrow it down.
- **Reading** — click any result or list item to open the captured page.
  In-page links work; external links are blocked (copy and paste the address
  into your browser if you need to open one).
- **Back** — the ← arrow in the top bar returns to the previous page.

### User accounts & comments

- **Register** with an email and a unique public username. The first registered
  user automatically becomes an **admin**.
- **Comments** appear at the bottom of every content page — visible to everyone,
  but only logged-in users can post.
- **Edit / Delete** — users can edit their own comments. Admins and managers can
  delete any comment (moderation).
- **Bookmarks** — click the pin button at the top of any page to bookmark it.
  Your bookmarks appear in the ★ Bookmarks dropdown in the nav bar.
- **Recent comments** — the home page shows the latest comments across all pages
  so you can see what's being discussed.

**Three user roles:**

| Role | Comment | Edit Own | Delete Others | Manage Users |
|------|---------|----------|---------------|--------------|
| admin | yes | yes | yes | yes |
| manager | yes | yes | yes (moderation) | no |
| general | yes | yes | no | no |

---

## Repository contents

| Path | What it is |
|------|------------|
| `cmd/difilo/main.go` | Entry point: flags, DB open, content import, graceful shutdown. |
| `internal/textutil/` | Shared text helpers (prose cleaning, excerpts, snippets, formatting). |
| `internal/db/` | SQLite layer: schema, users, comments, bookmarks, content pages, FTS5 search. |
| `internal/auth/` | Authentication: bcrypt passwords, signed session cookies, role permissions. |
| `internal/content/` | Content pipeline: import from mirror, markdown rendering, wiki page layout. |
| `internal/ui/` | HTML/CSS/JS fragments: nav panel, page shell, comments, bookmarks, CSS. |
| `internal/app/` | Application wiring: Server struct, routing, all HTTP handlers, middleware. |
| `go.mod` | Go module (deps: `modernc.org/sqlite`, `golang.org/x/crypto`, `goldmark`). |
| `scripts/` | Double-click launchers for Mac (`.command`) and Windows (`.bat`). |
| `mirror/html/` | Captured pages (~11,400), stored **gzip-compressed at rest**. |
| `mirror/md/` | Markdown source for the search index (~11,400), also gzip-compressed. |
| `mirror/images/` | All images as **WebP** (originally jpg/png/gif). |
| `mirror/vendor/` | Vendored CDN assets (Bootstrap, jQuery) for offline use. |
| `mirror/pages.json` | Page manifest (URL, title, section) used to build the index. |
| `mirror/difilo.db` | SQLite database (users, comments, bookmarks, content). Auto-created on first run. |

## How it works

A single Go binary serves the `mirror/` directory on `localhost`.

- **Routing** — page routes resolve via the database: direct route lookup,
  alias resolution (token-sorted slug matching), or fuzzy slug search. A
  normalized fallback handles `+`/`_`/`-`/space slug differences.
- **Rendering** — pages are rendered from stored markdown via Goldmark (GFM
  tables, strikethrough, autolinks), with URL rewriting to localize internal
  links and images for offline use. An overlay panel (logo, grouped nav,
  search, user menu) is injected after `<body>`.
- **Search** — SQLite FTS5 full-text search with BM25 ranking, snippet
  highlighting, and a LIKE fallback for special characters. Built from the
  `pages` table on import; stays in sync via triggers.
- **Auth & database** — user accounts, comments, bookmarks, and all content
  are stored in a SQLite database (`mirror/difilo.db`). Passwords are
  bcrypt-hashed; sessions use signed HMAC cookies (no external auth service
  needed). The database is auto-created on first run.
- **Graceful shutdown** — the server handles SIGTERM/SIGINT for clean
  Docker/container shutdowns.

### Storage optimizations (no data loss)

To keep the archive compact without losing anything:

- **Images → WebP.** Every raster image was converted to WebP (longest side
  capped at 1200 px, quality 80). Old `.jpg`/`.png` URLs still resolve — the
  server transparently serves the `.webp` sibling with the correct content type.
- **HTML & Markdown → gzip at rest.** Captured pages and markdown are stored as
  gzip-compressed bytes (the filename stays `.html`/`.md`). The server detects
  the gzip magic header and decompresses transparently before rewriting/serving
  or indexing. Files that are still plain are passed through unchanged, so a mix
  is fine.

Result: the mirror compresses from ~4.6 GB to ~1 GB with zero data loss.

---

## Building from source

Requires [Go](https://go.dev/) 1.25+.

```bash
go build -o DiFiLo ./cmd/difilo
# Windows:
GOOS=windows go build -o DiFiLo.exe ./cmd/difilo
```

## Running

```bash
./DiFiLo --mirror ./mirror --port 8000
```

Flags:

| Flag | Default | Purpose |
|------|---------|---------|
| `--mirror` | `./mirror` | Path to the mirror directory. |
| `--host` | `127.0.0.1` | Interface to bind. Use `0.0.0.0` to expose (e.g. in Docker). |
| `--port` | `8000` | Port to serve on. |
| `--reindex` | off | Force-rebuild the search index from `md/` + `pages.json`. |

To rebuild the search index after editing markdown:

```bash
./DiFiLo --mirror ./mirror --reindex
```

## Docker / self-hosting

A multi-stage `Dockerfile` builds a Linux binary from source and bakes the
mirror into the image, so difilo runs as a container — no volume mounts needed.

```bash
docker build -t difilo .
docker run -d --name difilo --restart unless-stopped -p 8299:8000 difilo
# then visit http://<host>:8299/
```

`.gitlab-ci.yml` automates this: on every push to `main` it rebuilds the image
and redeploys the `difilo` container (host port `8299` → container `8000`) via
the host Docker socket.

---

## About Tony Hansen

Tony Hansen is a ceramic engineer and the author of Digitalfire and
Insight-Live. Everything in this archive is his work. If it has helped you,
the best thank-you is to buy him a coffee:

**https://ko-fi.com/tonyhansen**

Visit the original website: **https://www.digitalfire.com**
