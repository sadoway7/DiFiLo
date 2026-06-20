package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// DB wraps a *sql.DB connection to the application's SQLite database.
// The underlying handle is unexported so all access goes through the
// package's query methods.
type DB struct {
	db *sql.DB
}

// SQL returns the underlying *sql.DB handle.
// This is intended for the content import pipeline, which needs raw
// transaction access for batch inserts.
func (d *DB) SQL() *sql.DB {
	return d.db
}

// Exec is a convenience wrapper around the underlying db.Exec.
func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// OpenDB opens (or creates) the SQLite database at path, applies the full
// schema (existing tables plus placeholder tables for upcoming wiki
// features), and runs the lightweight migrations required when upgrading
// from older schemas.
func OpenDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite single-writer

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		email        TEXT UNIQUE NOT NULL,
		username     TEXT UNIQUE NOT NULL DEFAULT '',
		password_hash TEXT NOT NULL,
		role         TEXT NOT NULL DEFAULT 'general',
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS comments (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id    INTEGER NOT NULL REFERENCES users(id),
		route      TEXT NOT NULL,
		body       TEXT NOT NULL,
		parent_id  INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_comments_route ON comments(route);
	CREATE TABLE IF NOT EXISTS bookmarks (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id    INTEGER NOT NULL REFERENCES users(id),
		route      TEXT NOT NULL,
		title      TEXT NOT NULL DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, route)
	);
	CREATE TABLE IF NOT EXISTS settings (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL DEFAULT ''
	);
	CREATE TABLE IF NOT EXISTS downloads (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id    INTEGER NOT NULL REFERENCES users(id),
		username   TEXT NOT NULL DEFAULT '',
		route      TEXT NOT NULL,
		title      TEXT NOT NULL DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_downloads_created ON downloads(created_at DESC);

	-- ---- Content tables ----

	CREATE TABLE IF NOT EXISTS pages (
		id               INTEGER PRIMARY KEY,
		section          TEXT NOT NULL,
		title            TEXT NOT NULL DEFAULT '',
		slug             TEXT NOT NULL DEFAULT '',
		source_url       TEXT NOT NULL DEFAULT '',
		route            TEXT NOT NULL DEFAULT '',
		wayback_ts       TEXT NOT NULL DEFAULT '',
		html_sha1        TEXT NOT NULL DEFAULT '',
		body_text        TEXT NOT NULL DEFAULT '',
		body_md          TEXT NOT NULL DEFAULT '',
		thumb            TEXT NOT NULL DEFAULT '',
		meta_description TEXT NOT NULL DEFAULT '',
		author_byline    TEXT NOT NULL DEFAULT '',
		status           TEXT NOT NULL DEFAULT 'published',
		sort_title       TEXT NOT NULL DEFAULT '',
		word_count       INTEGER NOT NULL DEFAULT 0,
		created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(route)
	);
	CREATE INDEX IF NOT EXISTS idx_pages_section ON pages(section);
	CREATE INDEX IF NOT EXISTS idx_pages_status ON pages(status);

	CREATE VIRTUAL TABLE IF NOT EXISTS pages_fts USING fts5(
		title, body_text, section,
		content='pages', content_rowid='id',
		tokenize='porter unicode61'
	);
	CREATE TRIGGER IF NOT EXISTS pages_ai AFTER INSERT ON pages BEGIN
		INSERT INTO pages_fts(rowid, title, body_text, section)
		VALUES (new.id, new.title, new.body_text, new.section);
	END;
	CREATE TRIGGER IF NOT EXISTS pages_ad AFTER DELETE ON pages BEGIN
		INSERT INTO pages_fts(pages_fts, rowid, title, body_text, section)
		VALUES ('delete', old.id, old.title, old.body_text, old.section);
	END;
	CREATE TRIGGER IF NOT EXISTS pages_au AFTER UPDATE ON pages BEGIN
		INSERT INTO pages_fts(pages_fts, rowid, title, body_text, section)
		VALUES ('delete', old.id, old.title, old.body_text, old.section);
		INSERT INTO pages_fts(rowid, title, body_text, section)
		VALUES (new.id, new.title, new.body_text, new.section);
	END;

	CREATE TABLE IF NOT EXISTS page_images (
		id           INTEGER PRIMARY KEY,
		page_id      INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
		image_path   TEXT NOT NULL DEFAULT '',
		original_ref TEXT NOT NULL DEFAULT '',
		caption      TEXT NOT NULL DEFAULT '',
		is_primary   INTEGER NOT NULL DEFAULT 0,
		sort_order   INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_page_images_page ON page_images(page_id);

	CREATE TABLE IF NOT EXISTS page_links (
		id           INTEGER PRIMARY KEY,
		page_id      INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
		target_type  TEXT NOT NULL DEFAULT '',
		target_url   TEXT NOT NULL DEFAULT '',
		target_route TEXT NOT NULL DEFAULT '',
		target_title TEXT NOT NULL DEFAULT '',
		blurb        TEXT NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_page_links_page ON page_links(page_id);
	CREATE INDEX IF NOT EXISTS idx_page_links_target ON page_links(target_route);

	-- ---- Placeholder tables for upcoming wiki features ----

	CREATE TABLE IF NOT EXISTS schema_version (
		version    INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS page_revisions (
		id           INTEGER PRIMARY KEY,
		page_id      INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
		body_md      TEXT NOT NULL,
		edited_by    INTEGER REFERENCES users(id),
		edit_summary TEXT DEFAULT '',
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_revisions_page ON page_revisions(page_id, created_at DESC);
	CREATE TABLE IF NOT EXISTS edit_proposals (
		id               INTEGER PRIMARY KEY,
		page_id          INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
		proposed_body_md TEXT NOT NULL,
		proposed_by      INTEGER NOT NULL REFERENCES users(id),
		status           TEXT DEFAULT 'pending',
		reviewed_by      INTEGER REFERENCES users(id),
		review_note      TEXT DEFAULT '',
		created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
		reviewed_at      DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_proposals_status ON edit_proposals(status);
	CREATE TABLE IF NOT EXISTS page_maintainers (
		page_id INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		PRIMARY KEY(page_id, user_id)
	);
	CREATE TABLE IF NOT EXISTS page_tags (
		id      INTEGER PRIMARY KEY,
		page_id INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
		tag     TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_tags_tag ON page_tags(tag);
	CREATE INDEX IF NOT EXISTS idx_tags_page ON page_tags(page_id);
	CREATE TABLE IF NOT EXISTS notifications (
		id         INTEGER PRIMARY KEY,
		user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type       TEXT NOT NULL,
		message    TEXT NOT NULL,
		route      TEXT DEFAULT '',
		read       INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, read);
	CREATE TABLE IF NOT EXISTS page_uploads (
		id         INTEGER PRIMARY KEY,
		user_id    INTEGER NOT NULL REFERENCES users(id),
		page_id    INTEGER REFERENCES pages(id) ON DELETE SET NULL,
		filename   TEXT NOT NULL,
		image_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
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
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	// Migration: add username column if upgrading from an older schema.
	db.Exec("ALTER TABLE users ADD COLUMN username TEXT NOT NULL DEFAULT ''")
	// Ensure username uniqueness for existing databases.
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE username != ''")
	// Enable foreign keys for cascade deletes.
	db.Exec("PRAGMA foreign_keys = ON")
	return &DB{db: db}, nil
}

// ContentImported reports whether the pages table already has any rows.
// main.go uses this to decide whether to run the initial content import.
func (d *DB) ContentImported() bool {
	var n int
	err := d.db.QueryRow("SELECT COUNT(*) FROM pages LIMIT 1").Scan(&n)
	if err != nil {
		return false
	}
	return n > 0
}
