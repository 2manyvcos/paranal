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
      CREATE TABLE IF NOT EXISTS usercredentials (
        name TEXT PRIMARY KEY NOT NULL,
        description TEXT NOT NULL,
        value TEXT NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"usercredentials\" failed - %s", err)
		}
		return nil
	})
}

func UserCredentialQuery(query *schema.UserCredentialQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.UserCredentialQuery)
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

func (i *impl) ListUserCredentials(query *schema.UserCredentialQuery) ([]schema.UserCredential, error) {
	where, wherePlaceholders := UserCredentialQuery(query)
	rows, err := i.Query(
		`
      SELECT name, description, value
      FROM usercredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.UserCredential
	for rows.Next() {
		var record schema.UserCredential
		if err := rows.Scan(&record.Name, &record.Description, &record.Value); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetUserCredential(query schema.UserCredentialQuery) (schema.UserCredential, error) {
	where, wherePlaceholders := UserCredentialQuery(&query)
	var record schema.UserCredential
	err := i.QueryRow(
		`
      SELECT name, description, value
      FROM usercredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.Description, &record.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.UserCredential{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateUserCredential(record schema.UserCredential) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO usercredentials (name, description, value)
      VALUES (?, ?, ?)
    `,
		&record.Name, &record.Description, &record.Value,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateUserCredentials(query schema.UserCredentialQuery, record schema.UserCredential) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := UserCredentialQuery(&query)
	result, err := i.Exec(
		`
      UPDATE usercredentials
      SET
        name = ?,
        description = ?,
        value = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.Description, &record.Value},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteUserCredentials(query schema.UserCredentialQuery) error {
	where, wherePlaceholders := UserCredentialQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM usercredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
