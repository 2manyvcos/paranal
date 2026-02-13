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
      CREATE TABLE IF NOT EXISTS useralertchannels (
        id TEXT PRIMARY KEY NOT NULL,
        userName TEXT NOT NULL,
        url TEXT NOT NULL,
        errorAlerts BOOLEAN NOT NULL,
        uptimeAlerts BOOLEAN NOT NULL,
        versionAlerts BOOLEAN NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"useralertchannels\" failed - %s", err)
		}
		return nil
	})
}

func UserAlertChannelQuery(query *schema.UserAlertChannelQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.UserAlertChannelQuery)
	}
	var conditions []string
	if query.ID != nil {
		conditions = append(conditions, "useralertchannels.id = ?")
		placeholders = append(placeholders, *query.ID)
	}
	if query.UserName != nil {
		conditions = append(conditions, "useralertchannels.userName = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.HasURL != nil {
		if *query.HasURL {
			conditions = append(conditions, "IFNULL(useralertchannels.url, '') <> ''")
		} else {
			conditions = append(conditions, "IFNULL(useralertchannels.url, '') = ''")
		}
	}
	if query.ErrorAlerts != nil {
		conditions = append(conditions, "IFNULL(useralertchannels.errorAlerts, FALSE) = ?")
		placeholders = append(placeholders, *query.ErrorAlerts)
	}
	if query.UptimeAlerts != nil {
		conditions = append(conditions, "IFNULL(useralertchannels.uptimeAlerts, FALSE) = ?")
		placeholders = append(placeholders, *query.UptimeAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "IFNULL(useralertchannels.versionAlerts, FALSE) = ?")
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
      SELECT useralertchannels.id, useralertchannels.userName, useralertchannels.url, useralertchannels.errorAlerts, useralertchannels.uptimeAlerts, useralertchannels.versionAlerts
      FROM useralertchannels
      INNER JOIN users
      ON useralertchannels.userName = users.name
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
		if err := rows.Scan(&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts); err != nil {
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
      SELECT useralertchannels.id, useralertchannels.userName, useralertchannels.url, useralertchannels.errorAlerts, useralertchannels.uptimeAlerts, useralertchannels.versionAlerts
      FROM useralertchannels
      INNER JOIN users
      ON useralertchannels.userName = users.name
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts)
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
      INSERT INTO useralertchannels (id, userName, url, errorAlerts, uptimeAlerts, versionAlerts)
      VALUES (?, ?, ?, ?, ?, ?)
    `,
		&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts,
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
      UPDATE useralertchannels
      SET
        id = ?,
        userName = ?,
        url = ?,
        errorAlerts = ?,
        uptimeAlerts = ?,
        versionAlerts = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteUserAlertChannels(query schema.UserAlertChannelQuery) error {
	where, wherePlaceholders := UserAlertChannelQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM useralertchannels
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
