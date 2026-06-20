package db

import (
	"database/sql"
	"time"
)

// User is one row from the users table. The Role field is a plain string;
// role constants live in the auth package so this package has no dependency
// on auth.
type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

// UserCount returns the total number of registered users.
func (d *DB) UserCount() (int, error) {
	var n int
	err := d.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&n)
	return n, err
}

// CreateUser inserts a new user and returns the populated User.
func (d *DB) CreateUser(email, username, passwordHash, role string) (*User, error) {
	res, err := d.db.Exec(
		"INSERT INTO users (email, username, password_hash, role) VALUES (?, ?, ?, ?)",
		email, username, passwordHash, role)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &User{ID: id, Email: email, Username: username, PasswordHash: passwordHash, Role: role}, nil
}

// GetUserByEmail returns the user with the given email, or nil if none.
func (d *DB) GetUserByEmail(email string) (*User, error) {
	u := &User{}
	err := d.db.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE email = ?",
		email).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// GetUserByUsername returns the user with the given username, or nil if none.
func (d *DB) GetUserByUsername(username string) (*User, error) {
	u := &User{}
	err := d.db.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE username = ?",
		username).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// GetUserByID returns the user with the given id, or nil if none.
func (d *DB) GetUserByID(id int64) (*User, error) {
	u := &User{}
	err := d.db.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE id = ?",
		id).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// AllUsers returns every user ordered by creation time.
func (d *DB) AllUsers() ([]User, error) {
	rows, err := d.db.Query(
		"SELECT id, email, username, password_hash, role, created_at FROM users ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// UpdateUserRole sets the role of the user with the given id.
func (d *DB) UpdateUserRole(id int64, role string) error {
	_, err := d.db.Exec("UPDATE users SET role = ? WHERE id = ?", role, id)
	return err
}

// DeleteUser removes the user and all their related data (comments,
// bookmarks, downloads, proposals, uploads). Tables with ON DELETE CASCADE
// (page_maintainers, notifications) are handled automatically.
func (d *DB) DeleteUser(id int64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Delete related rows that lack ON DELETE CASCADE
	tx.Exec("DELETE FROM comments WHERE user_id = ?", id)
	tx.Exec("DELETE FROM bookmarks WHERE user_id = ?", id)
	tx.Exec("DELETE FROM downloads WHERE user_id = ?", id)
	tx.Exec("DELETE FROM edit_proposals WHERE proposed_by = ?", id)
	tx.Exec("DELETE FROM edit_proposals WHERE reviewed_by = ?", id)
	tx.Exec("DELETE FROM page_revisions WHERE edited_by = ?", id)
	tx.Exec("DELETE FROM page_uploads WHERE user_id = ?", id)
	// Now safe to delete the user
	if _, err := tx.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
		return err
	}
	return tx.Commit()
}
