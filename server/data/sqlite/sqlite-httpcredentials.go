package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS httpcredentials (
        name TEXT PRIMARY KEY,
        type INTEGER,
        key TEXT,
        value TEXT
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"httpcredentials\" failed - %s", err)
		}
		return nil
	})
}

func HTTPCredentialQuery(query *schema.HTTPCredentialQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.HTTPCredentialQuery)
	}
	var conditions []string
	if query.Name != nil {
		conditions = append(conditions, "name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListHTTPCredentials(query *schema.HTTPCredentialQuery) ([]schema.HTTPCredential, error) {
	where, wherePlaceholders := HTTPCredentialQuery(query)
	rows, err := i.Query(
		`
      SELECT name, type, key, value
      FROM httpcredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.HTTPCredential
	for rows.Next() {
		var record schema.HTTPCredential
		if err := rows.Scan(&record.Name, &record.Type, &record.Key, &record.Value); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetHTTPCredential(query schema.HTTPCredentialQuery) (schema.HTTPCredential, error) {
	where, wherePlaceholders := HTTPCredentialQuery(&query)
	var record schema.HTTPCredential
	err := i.QueryRow(
		`
      SELECT name, type, key, value
      FROM httpcredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.Type, &record.Key, &record.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.HTTPCredential{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateHTTPCredential(record schema.HTTPCredential) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO httpcredentials (name, type, key, value)
      VALUES (?, ?, ?, ?)
    `,
		&record.Name, &record.Type, &record.Key, &record.Value,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateHTTPCredentials(query schema.HTTPCredentialQuery, record schema.HTTPCredential) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := HTTPCredentialQuery(&query)
	result, err := i.Exec(
		`
      UPDATE httpcredentials
      SET
        name = ?,
        type = ?,
        key = ?,
        value = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.Type, &record.Key, &record.Value},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteHTTPCredentials(query schema.HTTPCredentialQuery) error {
	where, wherePlaceholders := HTTPCredentialQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM httpcredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
