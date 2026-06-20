package db

// GetSetting returns the value stored for key, or "" if unset or on error.
func (d *DB) GetSetting(key string) string {
	var val string
	err := d.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&val)
	if err != nil {
		return ""
	}
	return val
}

// SetSetting upserts a key/value pair into the settings table.
func (d *DB) SetSetting(key, value string) error {
	_, err := d.db.Exec(
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		key, value, value)
	return err
}
