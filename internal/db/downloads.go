package db

// DownloadLog is one entry in the per-user page-download audit log.
type DownloadLog struct {
	ID        int64
	UserID    int64
	Username  string
	Route     string
	Title     string
	CreatedAt string
}

// LogDownload records that a user downloaded a page.
func (d *DB) LogDownload(userID int64, username, route, title string) error {
	_, err := d.db.Exec(
		"INSERT INTO downloads (user_id, username, route, title) VALUES (?, ?, ?, ?)",
		userID, username, route, title)
	return err
}

// RecentDownloads returns the most recent download events (newest first).
func (d *DB) RecentDownloads(limit int) ([]DownloadLog, error) {
	rows, err := d.db.Query(`
		SELECT id, user_id, username, route, title, strftime('%Y-%m-%d %H:%M', created_at)
		FROM downloads ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DownloadLog
	for rows.Next() {
		var dl DownloadLog
		if err := rows.Scan(&dl.ID, &dl.UserID, &dl.Username, &dl.Route, &dl.Title, &dl.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, dl)
	}
	return out, rows.Err()
}

// DownloadsByUser returns all download events for one user (newest first).
func (d *DB) DownloadsByUser(userID int64) ([]DownloadLog, error) {
	rows, err := d.db.Query(`
		SELECT id, user_id, username, route, title, strftime('%Y-%m-%d %H:%M', created_at)
		FROM downloads WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DownloadLog
	for rows.Next() {
		var dl DownloadLog
		if err := rows.Scan(&dl.ID, &dl.UserID, &dl.Username, &dl.Route, &dl.Title, &dl.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, dl)
	}
	return out, rows.Err()
}
