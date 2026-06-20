package db

// ImageRow is one entry from the page_images table.
type ImageRow struct {
	ImagePath   string
	OriginalRef string
	Caption     string
	IsPrimary   bool
	SortOrder   int
}

// GetPageImages returns all images for a page, ordered by is_primary then
// sort_order.
func (d *DB) GetPageImages(pageID int64) []ImageRow {
	rows, err := d.db.Query(`
		SELECT image_path, original_ref, caption, is_primary, sort_order
		FROM page_images
		WHERE page_id = ?
		ORDER BY is_primary DESC, sort_order`, pageID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ImageRow
	for rows.Next() {
		var r ImageRow
		var isPrim int
		if err := rows.Scan(&r.ImagePath, &r.OriginalRef, &r.Caption, &isPrim, &r.SortOrder); err != nil {
			continue
		}
		r.IsPrimary = isPrim != 0
		out = append(out, r)
	}
	return out
}
