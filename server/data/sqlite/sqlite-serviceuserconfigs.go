package sqlite

import (
	"fmt"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS serviceuserconfigs (
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
			return fmt.Errorf("creating table \"serviceuserconfigs\" failed - %s", err)
		}
		return nil
	})
}

func ServiceUserConfigQuery(query *schema.ServiceUserConfigQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceUserConfigQuery)
	}
	var conditions []string
	if query.UserName != nil {
		conditions = append(conditions, "serviceuserconfigs.userName = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "serviceuserconfigs.serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if query.Favorite != nil {
		conditions = append(conditions, "serviceuserconfigs.favorite = ?")
		placeholders = append(placeholders, *query.Favorite)
	}
	if query.Hidden != nil {
		conditions = append(conditions, "serviceuserconfigs.hidden = ?")
		placeholders = append(placeholders, *query.Hidden)
	}
	if query.UptimeAlerts != nil {
		conditions = append(conditions, "serviceuserconfigs.uptimeAlerts = ?")
		placeholders = append(placeholders, *query.UptimeAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "serviceuserconfigs.versionAlerts = ?")
		placeholders = append(placeholders, *query.VersionAlerts)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceUserConfigs(query *schema.ServiceUserConfigQuery) ([]schema.ServiceUserConfig, error) {
	where, wherePlaceholders := ServiceUserConfigQuery(query)
	rows, err := i.Query(
		`
      SELECT serviceuserconfigs.userName, serviceuserconfigs.serviceID, serviceuserconfigs.favorite, serviceuserconfigs.hidden, serviceuserconfigs.uptimeAlerts, serviceuserconfigs.versionAlerts
      FROM serviceuserconfigs
      INNER JOIN services
      ON serviceuserconfigs.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.ServiceUserConfig
	for rows.Next() {
		var record schema.ServiceUserConfig
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

func (i *impl) CreateOrUpdateServiceUserConfig(record schema.ServiceUserConfig) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO serviceuserconfigs (userName, serviceID, favorite, hidden, uptimeAlerts, versionAlerts)
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
