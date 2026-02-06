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
      CREATE TABLE IF NOT EXISTS users (
        name TEXT PRIMARY KEY NOT NULL,
        displayName TEXT NOT NULL,
        role INTEGER NOT NULL,
        passwordHash TEXT NOT NULL,
        errorAlerts BOOLEAN NOT NULL,
        uptimeAlerts BOOLEAN NOT NULL,
        versionAlerts BOOLEAN NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"users\" failed - %s", err)
		}
		return nil
	})
}

func UserQuery(query *schema.UserQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.UserQuery)
	}
	var conditions []string
	if query.Name != nil {
		conditions = append(conditions, "name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if query.ErrorAlerts != nil {
		conditions = append(conditions, "errorAlerts = ?")
		placeholders = append(placeholders, *query.ErrorAlerts)
	}
	if query.UptimeAlerts != nil {
		conditions = append(conditions, "uptimeAlerts = ?")
		placeholders = append(placeholders, *query.UptimeAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "versionAlerts = ?")
		placeholders = append(placeholders, *query.VersionAlerts)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListUsers(query *schema.UserQuery) ([]schema.User, error) {
	where, wherePlaceholders := UserQuery(query)
	rows, err := i.Query(
		`
      SELECT name, displayName, role, passwordHash, errorAlerts, uptimeAlerts, versionAlerts
      FROM users
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.User
	for rows.Next() {
		var record schema.User
		if err := rows.Scan(&record.Name, &record.DisplayName, &record.Role, &record.PasswordHash, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetUser(query schema.UserQuery) (schema.User, error) {
	where, wherePlaceholders := UserQuery(&query)
	var record schema.User
	err := i.QueryRow(
		`
      SELECT name, displayName, role, passwordHash, errorAlerts, uptimeAlerts, versionAlerts
      FROM users
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.DisplayName, &record.Role, &record.PasswordHash, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.User{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateUser(record schema.User) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO users (name, displayName, role, passwordHash, errorAlerts, uptimeAlerts, versionAlerts)
      VALUES (?, ?, ?, ?, ?, ?, ?)
    `,
		&record.Name, &record.DisplayName, &record.Role, &record.PasswordHash, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts,
	)
	return requireNoConflict(err)
}

func (i *impl) CreateOrUpdateUser(record schema.User) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO users (name, displayName, role, passwordHash, errorAlerts, uptimeAlerts, versionAlerts)
      VALUES (?, ?, ?, ?, ?, ?, ?)
      ON CONFLICT (name)
      DO UPDATE
      SET
        displayName = excluded.displayName,
        role = excluded.role,
        passwordHash = excluded.passwordHash,
        errorAlerts = excluded.errorAlerts,
        uptimeAlerts = excluded.uptimeAlerts,
        versionAlerts = excluded.versionAlerts
    `,
		&record.Name, &record.DisplayName, &record.Role, &record.PasswordHash, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts,
	)
	return err
}

func (i *impl) UpdateUsers(query schema.UserQuery, record schema.User) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := UserQuery(&query)
	result, err := i.Exec(
		`
      UPDATE users
      SET
        name = ?,
        displayName = ?,
        role = ?,
        passwordHash = ?,
        errorAlerts = ?,
        uptimeAlerts = ?,
        versionAlerts = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.DisplayName, &record.Role, &record.PasswordHash, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteUsers(query schema.UserQuery) error {
	where, wherePlaceholders := UserQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM users
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
