package sqlite

import (
	"fmt"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS service_configs (
        user_name TEXT NOT NULL,
        service_id TEXT NOT NULL,
        favorite BOOLEAN NOT NULL,
        hidden BOOLEAN NOT NULL,
        health_alerts BOOLEAN NOT NULL,
        version_alerts BOOLEAN NOT NULL,
        PRIMARY KEY (user_name, service_id)
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"service_configs\" failed - %s", err)
		}
		return nil
	})

	cleanups = append(cleanups, func(i *impl) error {
		_, err := i.Exec(`
      DELETE FROM service_configs
      WHERE
        ROWID NOT IN (
          SELECT service_configs.ROWID
          FROM service_configs
          INNER JOIN services
          ON service_configs.service_id = services.id
          INNER JOIN users
          ON service_configs.user_name = users.name
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
		conditions = append(conditions, "service_configs.user_name = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "service_configs.service_id = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if query.Favorite != nil {
		conditions = append(conditions, "IFNULL(service_configs.favorite, FALSE) = ?")
		placeholders = append(placeholders, *query.Favorite)
	}
	if query.Hidden != nil {
		conditions = append(conditions, "IFNULL(service_configs.hidden, FALSE) = ?")
		placeholders = append(placeholders, *query.Hidden)
	}
	if query.HealthAlerts != nil {
		conditions = append(conditions, "IFNULL(service_configs.health_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.HealthAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "IFNULL(service_configs.version_alerts, FALSE) = ?")
		placeholders = append(placeholders, *query.VersionAlerts)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) CreateOrUpdateServiceConfig(record schema.ServiceConfig) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO service_configs (user_name, service_id, favorite, hidden, health_alerts, version_alerts)
      VALUES (?, ?, ?, ?, ?, ?)
      ON CONFLICT (user_name, service_id)
      DO UPDATE
      SET
        favorite = excluded.favorite,
        hidden = excluded.hidden,
        health_alerts = excluded.health_alerts,
        version_alerts = excluded.version_alerts
    `,
		&record.UserName, &record.ServiceID, &record.Favorite, &record.Hidden, &record.HealthAlerts, &record.VersionAlerts,
	)
	return err
}
