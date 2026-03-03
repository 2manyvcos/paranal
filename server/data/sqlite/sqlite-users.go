package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS users (
        name TEXT PRIMARY KEY NOT NULL,
        role INTEGER NOT NULL,
        password_hash TEXT NOT NULL,
        display_name TEXT NOT NULL,
        start_page TEXT NOT NULL,
        error_alerts BOOLEAN NOT NULL,
        health_alerts BOOLEAN NOT NULL,
        version_alerts BOOLEAN NOT NULL
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
		conditions = append(conditions, "users.name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if query.Role != nil {
		conditions = append(conditions, "users.role = ?")
		placeholders = append(placeholders, *query.Role)
	}
	if query.HasPassword != nil {
		if *query.HasPassword {
			conditions = append(conditions, "IFNULL(users.password_hash, '') <> ''")
		} else {
			conditions = append(conditions, "IFNULL(users.password_hash, '') = ''")
		}
	}
	if query.DisplayName != nil {
		conditions = append(conditions, "users.display_name = ?")
		placeholders = append(placeholders, *query.DisplayName)
	}
	if query.StartPage != nil {
		conditions = append(conditions, "users.start_page = ?")
		placeholders = append(placeholders, *query.StartPage)
	}
	if query.ErrorAlerts != nil {
		conditions = append(conditions, "IFNULL(users.error_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.ErrorAlerts)
	}
	if query.HealthAlerts != nil {
		conditions = append(conditions, "IFNULL(users.health_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.HealthAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "IFNULL(users.version_alerts, FALSE) = ?")
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
      SELECT name, role, password_hash, display_name, start_page, error_alerts, health_alerts, version_alerts
      FROM users
      WHERE `+where+`
      ORDER BY name
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
		if err := rows.Scan(&record.Name, &record.Role, &record.PasswordHash, &record.DisplayName, &record.StartPage, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts); err != nil {
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
      SELECT name, role, password_hash, display_name, start_page, error_alerts, health_alerts, version_alerts
      FROM users
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.Role, &record.PasswordHash, &record.DisplayName, &record.StartPage, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts)
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
      INSERT INTO users (name, role, password_hash, display_name, start_page, error_alerts, health_alerts, version_alerts)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `,
		&record.Name, &record.Role, &record.PasswordHash, &record.DisplayName, &record.StartPage, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts,
	)
	return requireNoConflict(err)
}

func (i *impl) CreateOrUpdateUser(record schema.User) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO users (name, role, password_hash, display_name, start_page, error_alerts, health_alerts, version_alerts)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?)
      ON CONFLICT (name)
      DO UPDATE
      SET
        role = excluded.role,
        password_hash = excluded.password_hash,
        display_name = excluded.display_name,
        start_page = excluded.start_page,
        error_alerts = excluded.error_alerts,
        health_alerts = excluded.health_alerts,
        version_alerts = excluded.version_alerts
    `,
		&record.Name, &record.Role, &record.PasswordHash, &record.DisplayName, &record.StartPage, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts,
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
        role = ?,
        password_hash = ?,
        display_name = ?,
        start_page = ?,
        error_alerts = ?,
        health_alerts = ?,
        version_alerts = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.Role, &record.PasswordHash, &record.DisplayName, &record.StartPage, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts},
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
