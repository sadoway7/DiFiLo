package db

import (
	"database/sql"
	"fmt"
	"strings"

	"difilo/internal/textutil"
)

// SearchHit is one FTS5 search result.
type SearchHit struct {
	Page    ContentPage
	Snippet string
	Rank    float64
}

// SearchPages runs an FTS5 query and returns ranked results with snippets.
// When the MATCH expression fails (e.g. special characters), it falls back
// to a LIKE-based search that builds snippets via textutil.MakeSnippet.
func (d *DB) SearchPages(q string, limit int) []SearchHit {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	// Build FTS5 MATCH expression: wrap each term in quotes for phrase matching
	terms := strings.Fields(q)
	matchParts := make([]string, len(terms))
	for i, t := range terms {
		matchParts[i] = "\"" + strings.ReplaceAll(t, "\"", "") + "\""
	}
	matchExpr := strings.Join(matchParts, " ")

	query := fmt.Sprintf(`
		SELECT p.id, p.section, p.title, p.slug, p.source_url, p.route,
		       p.wayback_ts, p.html_sha1, p.body_text, p.thumb,
		       p.meta_description, p.author_byline, p.status, p.sort_title,
		       p.word_count,
		       snippet(pages_fts, 1, '<mark>', '</mark>', ' … ', 24) AS snip,
		       bm25(pages_fts) AS rank
		FROM pages_fts
		JOIN pages p ON p.id = pages_fts.rowid
		WHERE pages_fts MATCH ?
		  AND p.status = 'published'
		ORDER BY rank
		LIMIT %d`, limit)

	rows, err := d.db.Query(query, matchExpr)
	if err != nil {
		// FTS5 MATCH syntax can be finicky; fall back to LIKE
		return d.searchFallback(q, limit)
	}
	defer rows.Close()

	var hits []SearchHit
	for rows.Next() {
		var h SearchHit
		var bodyText sql.NullString
		err := rows.Scan(
			&h.Page.ID, &h.Page.Section, &h.Page.Title, &h.Page.Slug,
			&h.Page.SourceURL, &h.Page.Route, &h.Page.WaybackTS,
			&h.Page.HTMLSHA1, &bodyText, &h.Page.Thumb,
			&h.Page.MetaDescription, &h.Page.AuthorByline,
			&h.Page.Status, &h.Page.SortTitle, &h.Page.WordCount,
			&h.Snippet, &h.Rank,
		)
		if err != nil {
			continue
		}
		h.Page.BodyText = bodyText.String
		hits = append(hits, h)
	}
	return hits
}

// searchFallback uses LIKE when FTS5 MATCH fails (e.g. special chars).
func (d *DB) searchFallback(q string, limit int) []SearchHit {
	like := "%" + strings.ToLower(q) + "%"
	query := fmt.Sprintf(`
		SELECT id, section, title, slug, source_url, route,
		       wayback_ts, html_sha1, body_text, thumb,
		       meta_description, author_byline, status, sort_title, word_count
		FROM pages
		WHERE status = 'published'
		  AND (LOWER(title) LIKE ? OR LOWER(body_text) LIKE ?)
		ORDER BY CASE WHEN LOWER(title) LIKE ? THEN 0 ELSE 1 END, sort_title
		LIMIT %d`, limit)

	rows, err := d.db.Query(query, like, like, like)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var hits []SearchHit
	for rows.Next() {
		var h SearchHit
		var bodyText sql.NullString
		err := rows.Scan(
			&h.Page.ID, &h.Page.Section, &h.Page.Title, &h.Page.Slug,
			&h.Page.SourceURL, &h.Page.Route, &h.Page.WaybackTS,
			&h.Page.HTMLSHA1, &bodyText, &h.Page.Thumb,
			&h.Page.MetaDescription, &h.Page.AuthorByline,
			&h.Page.Status, &h.Page.SortTitle, &h.Page.WordCount,
		)
		if err != nil {
			continue
		}
		h.Page.BodyText = bodyText.String
		// Build a simple snippet from body text.
		h.Snippet = textutil.MakeSnippet(h.Page.BodyText, h.Page.Title, strings.Fields(q), 30)
		hits = append(hits, h)
	}
	return hits
}

// SuggestTitles returns up to limit title suggestions for autocomplete.
func (d *DB) SuggestTitles(q string, limit int) []ContentPage {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	like := strings.ToLower(q) + "%"
	query := fmt.Sprintf(`
		SELECT id, section, title, route
		FROM pages
		WHERE status = 'published'
		  AND (LOWER(title) LIKE ? OR LOWER(title) LIKE ?)
		ORDER BY CASE WHEN LOWER(title) LIKE ? THEN 0 ELSE 1 END,
		         sort_title
		LIMIT %d`, limit)
	rows, err := d.db.Query(query, like, "% "+like, like)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ContentPage
	for rows.Next() {
		var p ContentPage
		if err := rows.Scan(&p.ID, &p.Section, &p.Title, &p.Route); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}
