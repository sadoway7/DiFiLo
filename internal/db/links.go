package db

// LinkRow is one entry from the page_links table.
type LinkRow struct {
	TargetType  string
	TargetURL   string
	TargetRoute string
	TargetTitle string
	Blurb       string
}

// GetPageLinks returns all outbound links for a page, ordered by id.
func (d *DB) GetPageLinks(pageID int64) []LinkRow {
	rows, err := d.db.Query(`
		SELECT target_type, target_url, target_route, target_title, blurb
		FROM page_links
		WHERE page_id = ?
		ORDER BY id`, pageID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []LinkRow
	for rows.Next() {
		var r LinkRow
		if err := rows.Scan(&r.TargetType, &r.TargetURL, &r.TargetRoute, &r.TargetTitle, &r.Blurb); err != nil {
			continue
		}
		out = append(out, r)
	}
	return out
}

// GetInboundLinks returns pages that link TO the given route.
func (d *DB) GetInboundLinks(route string) []ContentPage {
	rows, err := d.db.Query(`
		SELECT DISTINCT p.id, p.section, p.title, p.route
		FROM page_links l
		JOIN pages p ON p.id = l.page_id
		WHERE l.target_route = ?
		ORDER BY p.sort_title`, route)
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
