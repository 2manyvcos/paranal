package sqlite

import (
	"fmt"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS serviceconfigs (
        userName TEXT NOT NULL,
        serviceID TEXT NOT NULL,
        favorite BOOLEAN NOT NULL,
        hidden BOOLEAN NOT NULL,
        uptimeAlerts BOOLEAN NOT NULL,
        versionAlerts BOOLEAN NOT NULL,
        PRIMARY KEY (userName, serviceID)
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"serviceconfigs\" failed - %s", err)
		}
		return nil
	})

	cleanups = append(cleanups, func(i *impl) error {
		_, err := i.Exec(`
      DELETE FROM serviceconfigs
      WHERE
        ROWID NOT IN (
          SELECT serviceconfigs.ROWID
          FROM serviceconfigs
          INNER JOIN services
          ON serviceconfigs.serviceID = services.id
          INNER JOIN users
          ON serviceconfigs.userName = users.name
        )
    `)
		return err
	})
}

func ServiceConfigQuery(query *schema.ServiceConfigQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceConfigQuery)
	}
	var conditions []string
	if query.UserName != nil {
		conditions = append(conditions, "serviceconfigs.userName = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "serviceconfigs.serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if query.Favorite != nil {
		conditions = append(conditions, "IFNULL(serviceconfigs.favorite, FALSE) = ?")
		placeholders = append(placeholders, *query.Favorite)
	}
	if query.Hidden != nil {
		conditions = append(conditions, "IFNULL(serviceconfigs.hidden, FALSE) = ?")
		placeholders = append(placeholders, *query.Hidden)
	}
	if query.UptimeAlerts != nil {
		conditions = append(conditions, "IFNULL(serviceconfigs.uptimeAlerts, FALSE) = ?")
		placeholders = append(placeholders, *query.UptimeAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "IFNULL(serviceconfigs.versionAlerts, FALSE) = ?")
		placeholders = append(placeholders, *query.VersionAlerts)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceConfigs(query *schema.ServiceConfigQuery) ([]schema.ServiceConfig, error) {
	where, wherePlaceholders := ServiceConfigQuery(query)
	rows, err := i.Query(
		`
      SELECT serviceconfigs.userName, serviceconfigs.serviceID, serviceconfigs.favorite, serviceconfigs.hidden, serviceconfigs.uptimeAlerts, serviceconfigs.versionAlerts
      FROM serviceconfigs
      INNER JOIN services
      ON serviceconfigs.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.ServiceConfig
	for rows.Next() {
		var record schema.ServiceConfig
		if err := rows.Scan(&record.UserName, &record.ServiceID, &record.Favorite, &record.Hidden, &record.UptimeAlerts, &record.VersionAlerts); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) CreateOrUpdateServiceConfig(record schema.ServiceConfig) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO serviceconfigs (userName, serviceID, favorite, hidden, uptimeAlerts, versionAlerts)
      VALUES (?, ?, ?, ?, ?, ?)
      ON CONFLICT (userName, serviceID)
      DO UPDATE
      SET
        favorite = excluded.favorite,
        hidden = excluded.hidden,
        uptimeAlerts = excluded.uptimeAlerts,
        versionAlerts = excluded.versionAlerts
    `,
		&record.UserName, &record.ServiceID, &record.Favorite, &record.Hidden, &record.UptimeAlerts, &record.VersionAlerts,
	)
	return err
}
