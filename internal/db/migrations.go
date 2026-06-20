package db

// ApplyMigrations runs any pending schema migrations.
// Currently a no-op placeholder — the schema is created by OpenDB.
// Future migrations (ALTER TABLE) go here, tracked via schema_version table.
func (d *DB) ApplyMigrations() error {
	return nil
}
