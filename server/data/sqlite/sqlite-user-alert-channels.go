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
      CREATE TABLE IF NOT EXISTS user_alert_channels (
        id TEXT PRIMARY KEY NOT NULL,
        user_name TEXT NOT NULL,
        url TEXT NOT NULL,
        error_alerts BOOLEAN NOT NULL,
        health_alerts BOOLEAN NOT NULL,
        version_alerts BOOLEAN NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"user_alert_channels\" failed - %s", err)
		}
		return nil
	})

	cleanups = append(cleanups, func(i *impl) error {
		_, err := i.Exec(`
      DELETE FROM user_alert_channels
      WHERE
        ROWID NOT IN (
          SELECT user_alert_channels.ROWID
          FROM user_alert_channels
          INNER JOIN users
          ON user_alert_channels.user_name = users.name
        )
    `)
		return err
	})
}

func UserAlertChannelQuery(query *schema.UserAlertChannelQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.UserAlertChannelQuery)
	}
	var conditions []string
	if query.ID != nil {
		conditions = append(conditions, "user_alert_channels.id = ?")
		placeholders = append(placeholders, *query.ID)
	}
	if query.UserName != nil {
		conditions = append(conditions, "user_alert_channels.user_name = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.HasURL != nil {
		if *query.HasURL {
			conditions = append(conditions, "IFNULL(user_alert_channels.url, '') <> ''")
		} else {
			conditions = append(conditions, "IFNULL(user_alert_channels.url, '') = ''")
		}
	}
	if query.ErrorAlerts != nil {
		conditions = append(conditions, "IFNULL(user_alert_channels.error_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.ErrorAlerts)
	}
	if query.HealthAlerts != nil {
		conditions = append(conditions, "IFNULL(user_alert_channels.health_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.HealthAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "IFNULL(user_alert_channels.version_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.VersionAlerts)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListUserAlertChannels(query *schema.UserAlertChannelQuery) ([]schema.UserAlertChannel, error) {
	where, wherePlaceholders := UserAlertChannelQuery(query)
	rows, err := i.Query(
		`
      SELECT user_alert_channels.id, user_alert_channels.user_name, user_alert_channels.url, user_alert_channels.error_alerts, user_alert_channels.health_alerts, user_alert_channels.version_alerts
      FROM user_alert_channels
      INNER JOIN users
      ON user_alert_channels.user_name = users.name
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.UserAlertChannel
	for rows.Next() {
		var record schema.UserAlertChannel
		if err := rows.Scan(&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetUserAlertChannel(query schema.UserAlertChannelQuery) (schema.UserAlertChannel, error) {
	where, wherePlaceholders := UserAlertChannelQuery(&query)
	var record schema.UserAlertChannel
	err := i.QueryRow(
		`
      SELECT user_alert_channels.id, user_alert_channels.user_name, user_alert_channels.url, user_alert_channels.error_alerts, user_alert_channels.health_alerts, user_alert_channels.version_alerts
      FROM user_alert_channels
      INNER JOIN users
      ON user_alert_channels.user_name = users.name
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.UserAlertChannel{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateUserAlertChannel(record schema.UserAlertChannel) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO user_alert_channels (id, user_name, url, error_alerts, health_alerts, version_alerts)
      VALUES (?, ?, ?, ?, ?, ?)
    `,
		&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateUserAlertChannels(query schema.UserAlertChannelQuery, record schema.UserAlertChannel) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := UserAlertChannelQuery(&query)
	result, err := i.Exec(
		`
      UPDATE user_alert_channels
      SET
        id = ?,
        user_name = ?,
        url = ?,
        error_alerts = ?,
        health_alerts = ?,
        version_alerts = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteUserAlertChannels(query schema.UserAlertChannelQuery) error {
	where, wherePlaceholders := UserAlertChannelQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM user_alert_channels
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
