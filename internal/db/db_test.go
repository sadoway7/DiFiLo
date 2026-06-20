package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestDB opens a fresh database in a temp directory for testing.
func newTestDB(t *testing.T) *DB {
	t.Helper()
	d, err := OpenDB(filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	return d
}

func TestOpenDB_CreatesSchema(t *testing.T) {
	d := newTestDB(t)

	// All key tables (real tables + the FTS5 virtual table) must exist.
	want := []string{
		"users", "pages", "pages_fts", "comments", "bookmarks",
		"settings", "downloads",
		// Placeholder tables for upcoming wiki features.
		"page_revisions", "edit_proposals", "page_maintainers",
		"page_tags", "notifications", "page_uploads",
	}
	rows, err := d.db.Query("SELECT name FROM sqlite_master WHERE type IN ('table','view')")
	require.NoError(t, err)
	defer rows.Close()

	got := map[string]bool{}
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		got[name] = true
	}
	for _, w := range want {
		assert.True(t, got[w], "missing table %q", w)
	}
}

func TestUserCRUD(t *testing.T) {
	d := newTestDB(t)

	u, err := d.CreateUser("alice@example.com", "alice", "hash", "admin")
	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", u.Email)
	assert.NotZero(t, u.ID)

	// Retrieve by email.
	byEmail, err := d.GetUserByEmail("alice@example.com")
	require.NoError(t, err)
	require.NotNil(t, byEmail)
	assert.Equal(t, u.ID, byEmail.ID)
	assert.Equal(t, "admin", byEmail.Role)

	// Retrieve by id.
	byID, err := d.GetUserByID(u.ID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.Equal(t, "alice", byID.Username)

	// Missing email returns nil, nil.
	missing, err := d.GetUserByEmail("nope@example.com")
	require.NoError(t, err)
	assert.Nil(t, missing)
}

func TestCommentCRUD(t *testing.T) {
	d := newTestDB(t)

	u, err := d.CreateUser("bob@example.com", "bob", "hash", "general")
	require.NoError(t, err)

	id, err := d.CreateComment(u.ID, "/material/1", "nice page", sql.NullInt64{})
	require.NoError(t, err)
	assert.NotZero(t, id)

	got, err := d.GetCommentByID(id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "nice page", got.Body)
	assert.Equal(t, "bob", got.Username)
	assert.Equal(t, "/material/1", got.Route)

	list, err := d.GetCommentsByRoute("/material/1")
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, id, list[0].ID)
}

func TestSearchPages_FTS(t *testing.T) {
	d := newTestDB(t)

	// Insert a published page; the AFTER INSERT trigger populates pages_fts.
	_, err := d.db.Exec(`
		INSERT INTO pages
		    (section, title, slug, source_url, route, wayback_ts,
		     html_sha1, body_text, body_md, thumb, meta_description,
		     author_byline, status, sort_title, word_count)
		VALUES ('material', 'Cobalt Oxide', 'cobalt-oxide', '', '/material/1', '',
		        '', 'Cobalt oxide is a deep blue pigment used in glazes.', '', '', '',
		        '', 'published', 'cobalt oxide', 9)`)
	require.NoError(t, err)

	hits := d.SearchPages("cobalt", 10)
	require.Len(t, hits, 1, "FTS5 search should find the cobalt page")
	assert.Equal(t, "Cobalt Oxide", hits[0].Page.Title)
	assert.Equal(t, "/material/1", hits[0].Page.Route)
	assert.NotEmpty(t, hits[0].Snippet, "FTS snippet should be populated")
}

func TestSettings(t *testing.T) {
	d := newTestDB(t)

	assert.Equal(t, "", d.GetSetting("missing"))

	require.NoError(t, d.SetSetting("site_name", "DiFiLo"))
	assert.Equal(t, "DiFiLo", d.GetSetting("site_name"))

	// Upsert overwrites.
	require.NoError(t, d.SetSetting("site_name", "Updated"))
	assert.Equal(t, "Updated", d.GetSetting("site_name"))
}

func TestContentImported(t *testing.T) {
	d := newTestDB(t)

	assert.False(t, d.ContentImported(), "empty DB has no content")

	_, err := d.db.Exec(`INSERT INTO pages (section, route, title) VALUES ('material', '/material/1', 'Test')`)
	require.NoError(t, err)

	assert.True(t, d.ContentImported(), "should report content after insert")
}

func TestApplyMigrations(t *testing.T) {
	d := newTestDB(t)
	assert.NoError(t, d.ApplyMigrations())
}
