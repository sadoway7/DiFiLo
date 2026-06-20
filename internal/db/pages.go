package db

import "database/sql"

// ContentPage is a page row from the SQLite content tables.
type ContentPage struct {
	ID             int64
	Section        string
	Title          string
	Slug           string
	SourceURL      string
	Route          string
	WaybackTS      string
	HTMLSHA1       string
	BodyText       string
	BodyMD         string
	Thumb          string
	MetaDescription string
	AuthorByline   string
	Status         string
	SortTitle      string
	WordCount      int
}

// GetPageByRoute fetches a single published page by its route (e.g.
// "/material/2536"). Returns nil, nil when no such page exists.
func (d *DB) GetPageByRoute(route string) (*ContentPage, error) {
	p := &ContentPage{}
	var bodyMD, bodyText sql.NullString
	err := d.db.QueryRow(`
		SELECT id, section, title, slug, source_url, route,
		       wayback_ts, html_sha1, body_text, body_md, thumb,
		       meta_description, author_byline, status, sort_title, word_count
		FROM pages WHERE route = ? AND status = 'published'`, route).Scan(
		&p.ID, &p.Section, &p.Title, &p.Slug, &p.SourceURL, &p.Route,
		&p.WaybackTS, &p.HTMLSHA1, &bodyText, &bodyMD, &p.Thumb,
		&p.MetaDescription, &p.AuthorByline, &p.Status, &p.SortTitle, &p.WordCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.BodyText = bodyText.String
	p.BodyMD = bodyMD.String
	return p, nil
}

// ListPagesBySection returns all published pages in a section, sorted
// alphabetically by sort_title.
func (d *DB) ListPagesBySection(section string) []ContentPage {
	rows, err := d.db.Query(`
		SELECT id, section, title, route, thumb, word_count
		FROM pages
		WHERE section = ? AND status = 'published'
		ORDER BY sort_title`, section)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ContentPage
	for rows.Next() {
		var p ContentPage
		if err := rows.Scan(&p.ID, &p.Section, &p.Title, &p.Route, &p.Thumb, &p.WordCount); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}

// RandomContentPages returns n random pages (excluding the 'url' section).
func (d *DB) RandomContentPages(n int) []ContentPage {
	rows, err := d.db.Query(`
		SELECT id, section, title, route, thumb
		FROM pages
		WHERE section != 'url' AND status = 'published'
		ORDER BY RANDOM()
		LIMIT ?`, n)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ContentPage
	for rows.Next() {
		var p ContentPage
		if err := rows.Scan(&p.ID, &p.Section, &p.Title, &p.Route, &p.Thumb); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}

// RandomCardPages returns n random pages that have a thumbnail and enough
// body text to make an interesting card.
func (d *DB) RandomCardPages(n int) []ContentPage {
	return d.RandomCardsBySection("", n)
}

// RandomCardsBySection returns n random card-worthy pages from a specific
// section. If section is empty, returns from all sections (excluding 'url').
func (d *DB) RandomCardsBySection(section string, n int) []ContentPage {
	var rows *sql.Rows
	var err error
	if section == "" {
		rows, err = d.db.Query(`
			SELECT id, section, title, route, thumb, meta_description, body_text
			FROM pages
			WHERE section != 'url'
			  AND status = 'published'
			  AND thumb != ''
			  AND word_count >= 12
			ORDER BY RANDOM()
			LIMIT ?`, n)
	} else {
		rows, err = d.db.Query(`
			SELECT id, section, title, route, thumb, meta_description, body_text
			FROM pages
			WHERE section = ?
			  AND status = 'published'
			  AND thumb != ''
			  AND word_count >= 12
			ORDER BY RANDOM()
			LIMIT ?`, section, n)
	}
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ContentPage
	for rows.Next() {
		var p ContentPage
		if err := rows.Scan(&p.ID, &p.Section, &p.Title, &p.Route, &p.Thumb, &p.MetaDescription, &p.BodyText); err != nil {
			continue
		}
		out = append(out, p)
	}
	return out
}

// SectionCounts returns a map of section -> page count (excluding empty
// sections).
func (d *DB) SectionCounts() map[string]int {
	rows, err := d.db.Query(`
		SELECT section, COUNT(*)
		FROM pages
		WHERE status = 'published'
		GROUP BY section`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var sec string
		var n int
		if rows.Scan(&sec, &n) == nil {
			out[sec] = n
		}
	}
	return out
}

// PageCount returns the total number of published pages.
func (d *DB) PageCount() int {
	var n int
	d.db.QueryRow("SELECT COUNT(*) FROM pages WHERE status = 'published'").Scan(&n)
	return n
}
