package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS settings (
        name TEXT PRIMARY KEY NOT NULL,
        value TEXT NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"settings\" failed - %s", err)
		}
		return nil
	})
}

func (i *impl) GetSetting(name string) (string, error) {
	var value string
	err := i.QueryRow(
		`
      SELECT value
      FROM settings
      WHERE
        name = ?
    `,
		name,
	).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", schema.ErrNotFound
	}
	return value, err
}

func (i *impl) CreateOrUpdateSetting(name string, value string) error {
	_, err := i.Exec(
		`
      INSERT INTO settings (name, value)
      VALUES (?, ?)
      ON CONFLICT (name)
      DO UPDATE
      SET
        value = excluded.value
    `,
		&name, &value,
	)
	return err
}
