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
        uptimeAlert BOOLEAN NOT NULL,
        versionAlert BOOLEAN NOT NULL,
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
		conditions = append(conditions, "userName = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if query.Favorite != nil {
		conditions = append(conditions, "favorite = ?")
		placeholders = append(placeholders, *query.Favorite)
	}
	if query.Hidden != nil {
		conditions = append(conditions, "hidden = ?")
		placeholders = append(placeholders, *query.Hidden)
	}
	if query.UptimeAlert != nil {
		conditions = append(conditions, "uptimeAlert = ?")
		placeholders = append(placeholders, *query.UptimeAlert)
	}
	if query.VersionAlert != nil {
		conditions = append(conditions, "versionAlert = ?")
		placeholders = append(placeholders, *query.VersionAlert)
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
      SELECT userName, serviceID, favorite, hidden, uptimeAlert, versionAlert
      FROM serviceuserconfigs
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
		if err := rows.Scan(&record.UserName, &record.ServiceID, &record.Favorite, &record.Hidden, &record.UptimeAlert, &record.VersionAlert); err != nil {
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
      INSERT INTO serviceuserconfigs (userName, serviceID, favorite, hidden, uptimeAlert, versionAlert)
      VALUES (?, ?, ?, ?, ?, ?)
      ON CONFLICT (userName, serviceID)
      DO UPDATE
      SET
        favorite = excluded.favorite,
        hidden = excluded.hidden,
        uptimeAlert = excluded.uptimeAlert,
        versionAlert = excluded.versionAlert
    `,
		&record.UserName, &record.ServiceID, &record.Favorite, &record.Hidden, &record.UptimeAlert, &record.VersionAlert,
	)
	return err
}
