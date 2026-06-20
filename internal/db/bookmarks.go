package db

import "time"

// Bookmark is one row from the bookmarks table.
type Bookmark struct {
	ID        int64
	UserID    int64
	Route     string
	Title     string
	CreatedAt time.Time
}

// CreateBookmark records a bookmark for the given user and route. Duplicate
// (user, route) pairs are silently ignored.
func (d *DB) CreateBookmark(userID int64, route, title string) error {
	_, err := d.db.Exec(
		"INSERT OR IGNORE INTO bookmarks (user_id, route, title) VALUES (?, ?, ?)",
		userID, route, title)
	return err
}

// DeleteBookmark removes a user's bookmark for the given route.
func (d *DB) DeleteBookmark(userID int64, route string) error {
	_, err := d.db.Exec("DELETE FROM bookmarks WHERE user_id = ? AND route = ?", userID, route)
	return err
}

// GetBookmarksByUser returns all bookmarks for a user, newest first.
func (d *DB) GetBookmarksByUser(userID int64) ([]Bookmark, error) {
	rows, err := d.db.Query(
		"SELECT id, user_id, route, title, created_at FROM bookmarks WHERE user_id = ? ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Bookmark
	for rows.Next() {
		var b Bookmark
		if err := rows.Scan(&b.ID, &b.UserID, &b.Route, &b.Title, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// IsBookmarked reports whether the user has bookmarked the given route.
func (d *DB) IsBookmarked(userID int64, route string) (bool, error) {
	var exists int
	err := d.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = ? AND route = ?)",
		userID, route).Scan(&exists)
	return exists == 1, err
}
