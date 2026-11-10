package db

func (db *sqliteDB) setup() error {
	db.Lock()
	defer db.Unlock()

	_, err := db.Exec(
		`
      CREATE TABLE IF NOT EXISTS states (
        registry TEXT NOT NULL,
        id TEXT NOT NULL,
        rev TEXT NOT NULL,
        collector TEXT NOT NULL,
        data BLOB NOT NULL CHECK (json_valid(data, 8)),
        PRIMARY KEY (registry, id)
      )
      WITHOUT ROWID;

      CREATE INDEX IF NOT EXISTS idx_states_collector
      ON states(collector)
      WHERE collector != '';
    `,
	)
	return err
}
