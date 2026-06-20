package db

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reNormSep  = regexp.MustCompile(`[-_+\s\.]+`)
	reAliasTok = regexp.MustCompile(`[0-9a-z]+`)
)

// normName normalises a file path for fuzzy matching: lower-case, strip
// extension, replace separators with underscores.
func normName(rel string) string {
	s := strings.ToLower(rel)
	s = strings.TrimSuffix(s, ".html")
	s = reNormSep.ReplaceAllString(s, "_")
	return strings.Trim(s, "_")
}

// aliasKey produces a word-order-independent key from a string by extracting
// its alphanumeric tokens, lower-casing, sorting, and joining with "_".
// This lets slug URLs (e.g. /material/cobalt+oxide) match page titles
// (e.g. "Cobalt Oxide") regardless of token order.
func aliasKey(s string) string {
	toks := reAliasTok.FindAllString(strings.ToLower(s), -1)
	if len(toks) == 0 {
		return ""
	}
	sort.Strings(toks)
	return strings.Join(toks, "_")
}

// LoadAliases builds the alias map from the pages table. The key is
// aliasKey(section+" "+title) — a word-order-independent token hash — so
// slug URLs like /material/cobalt+oxide resolve to /material/230.
// Only unambiguous keys (one page per key) are included.
func (d *DB) LoadAliases() map[string]string {
	rows, err := d.db.Query(`
		SELECT section, title, route FROM pages
		WHERE status = 'published' AND route != '' AND title != ''`)
	if err != nil {
		return map[string]string{}
	}
	defer rows.Close()

	type pageInfo struct{ route string }
	keyPages := map[string][]pageInfo{}
	for rows.Next() {
		var section, title, route string
		if err := rows.Scan(&section, &title, &route); err != nil {
			continue
		}
		k := aliasKey(section + " " + title)
		if k == "" {
			continue
		}
		keyPages[k] = append(keyPages[k], pageInfo{route})
	}

	out := map[string]string{}
	for k, pages := range keyPages {
		if len(pages) == 1 {
			out[k] = pages[0].route
		}
	}
	return out
}

// ResolveRoute checks if a route exists directly in the DB, or via alias.
// Returns the canonical route and whether the page exists.
func (d *DB) ResolveRoute(route string, aliases map[string]string) (string, bool) {
	if i := strings.IndexByte(route, '?'); i >= 0 {
		route = route[:i]
	}
	// Direct lookup.
	var count int
	d.db.QueryRow(
		"SELECT COUNT(*) FROM pages WHERE route = ? AND status = 'published'",
		route,
	).Scan(&count)
	if count > 0 {
		return route, true
	}
	// Alias lookup — slug-based resolution.
	k := aliasKey(route)
	if k != "" {
		if canonical, ok := aliases[k]; ok && canonical != "" {
			return canonical, true
		}
	}
	return "", false
}

// SearchBySlug tries to find a page by extracting the last path segment from
// a route and searching for a title match in the same section. This is the
// last-resort fallback for old numeric IDs and slug mismatches.
//
// e.g. /test/14 -> search section="test" for title LIKE "%14%"
//      /glossary/analysis -> search section="glossary" for title="Analysis"
func (d *DB) SearchBySlug(route string) (*ContentPage, bool) {
	if i := strings.IndexByte(route, '?'); i >= 0 {
		route = route[:i]
	}
	route = strings.TrimPrefix(route, "/")
	parts := strings.SplitN(route, "/", 2)
	if len(parts) < 2 || parts[1] == "" {
		return nil, false
	}
	section := parts[0]
	slug := strings.ReplaceAll(parts[1], "+", " ")
	slug = strings.ReplaceAll(slug, "%20", " ")

	// Try exact title match in this section first
	p := &ContentPage{}
	err := d.db.QueryRow(`
		SELECT id, section, title, route FROM pages
		WHERE section = ? AND status = 'published'
		  AND (LOWER(title) = LOWER(?) OR LOWER(title) = LOWER(?))
		LIMIT 1`, section, slug, section+" "+slug).Scan(
		&p.ID, &p.Section, &p.Title, &p.Route)
	if err == nil {
		return p, true
	}

	// Try title containing the slug
	rows, err := d.db.Query(`
		SELECT id, section, title, route FROM pages
		WHERE section = ? AND status = 'published'
		  AND LOWER(title) LIKE ?
		ORDER BY sort_title LIMIT 1`, section, "%"+strings.ToLower(slug)+"%")
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Section, &p.Title, &p.Route)
		if err == nil {
			return p, true
		}
	}
	return nil, false
}
