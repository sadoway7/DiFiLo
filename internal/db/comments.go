package db

import "database/sql"

// Comment is one row from the comments table, joined with the author's
// username and role for display.
type Comment struct {
	ID        int64  `json:"ID"`
	UserID    int64  `json:"UserID"`
	Username  string `json:"Username"`
	Role      string `json:"Role"`
	Route     string `json:"Route"`
	Body      string `json:"Body"`
	CreatedAt string `json:"CreatedAt"`
	ParentID  sql.NullInt64
}

// CreateComment inserts a comment and returns its new id.
func (d *DB) CreateComment(userID int64, route, body string, parentID sql.NullInt64) (int64, error) {
	res, err := d.db.Exec(
		"INSERT INTO comments (user_id, route, body, parent_id) VALUES (?, ?, ?, ?)",
		userID, route, body, parentID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetCommentsByRoute returns all comments on a route, oldest first.
func (d *DB) GetCommentsByRoute(route string) ([]Comment, error) {
	rows, err := d.db.Query(`
		SELECT c.id, c.user_id, u.username, u.role, c.route, c.body, strftime('%Y-%m-%dT%H:%M:%SZ', c.created_at), c.parent_id
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.route = ?
		ORDER BY c.created_at ASC`, route)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.Username, &c.Role, &c.Route, &c.Body, &c.CreatedAt, &c.ParentID); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GetCommentByID returns the comment with the given id, or nil if none.
func (d *DB) GetCommentByID(id int64) (*Comment, error) {
	c := &Comment{}
	err := d.db.QueryRow(`
		SELECT c.id, c.user_id, u.username, u.role, c.route, c.body, strftime('%Y-%m-%dT%H:%M:%SZ', c.created_at), c.parent_id
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.id = ?`, id).
		Scan(&c.ID, &c.UserID, &c.Username, &c.Role, &c.Route, &c.Body, &c.CreatedAt, &c.ParentID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

// DeleteComment removes the comment with the given id.
func (d *DB) DeleteComment(id int64) error {
	_, err := d.db.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}

// UpdateComment sets the body of the comment with the given id.
func (d *DB) UpdateComment(id int64, body string) error {
	_, err := d.db.Exec("UPDATE comments SET body = ? WHERE id = ?", body, id)
	return err
}

// GetRecentComments returns the most recent comments across all routes,
// newest first.
func (d *DB) GetRecentComments(limit int) ([]Comment, error) {
	rows, err := d.db.Query(`
		SELECT c.id, c.user_id, u.username, u.role, c.route, c.body, strftime('%Y-%m-%dT%H:%M:%SZ', c.created_at), c.parent_id
		FROM comments c
		JOIN users u ON u.id = c.user_id
		ORDER BY c.created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.Username, &c.Role, &c.Route, &c.Body, &c.CreatedAt, &c.ParentID); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
