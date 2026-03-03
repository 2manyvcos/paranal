package sqlite

import (
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func UserServiceQuery(query *schema.UserServiceQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.UserServiceQuery)
	}
	var conditions []string
	serviceConditions, servicePlaceholders := ServiceQuery(&query.ServiceQuery)
	conditions = append(conditions, serviceConditions)
	placeholders = append(placeholders, servicePlaceholders...)
	serviceConfigConditions, serviceConfigPlaceholders := ServiceConfigQuery(&query.ServiceConfigQuery)
	conditions = append(conditions, serviceConfigConditions)
	placeholders = append(placeholders, serviceConfigPlaceholders...)
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListUserServices(userName string, query *schema.UserServiceQuery) ([]schema.UserService, error) {
	where, wherePlaceholders := UserServiceQuery(query)
	rows, err := i.Query(
		`
      SELECT services.id, services.name, services.description, services.logo, services.url, service_configs.favorite, service_configs.hidden, service_configs.health_alerts, service_configs.version_alerts
      FROM services
      LEFT JOIN service_configs
      ON service_configs.service_id = services.id AND service_configs.user_name = ?
      WHERE `+where+`
      ORDER BY services.name
    `,
		slices.Concat(
			[]any{userName},
			wherePlaceholders,
		)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.UserService
	for rows.Next() {
		var record schema.UserService
		var favorite *bool
		var hidden *bool
		var healthAlerts *bool
		var versionAlerts *bool
		if err := rows.Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favorite, &hidden, &healthAlerts, &versionAlerts); err != nil {
			return nil, err
		}
		if favorite != nil {
			record.Favorite = *favorite
		}
		if hidden != nil {
			record.Hidden = *hidden
		}
		if healthAlerts != nil {
			record.HealthAlerts = *healthAlerts
		}
		if versionAlerts != nil {
			record.VersionAlerts = *versionAlerts
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetUserService(userName string, query schema.UserServiceQuery) (schema.UserService, error) {
	where, wherePlaceholders := UserServiceQuery(&query)
	var record schema.UserService
	var favorite *bool
	var hidden *bool
	var healthAlerts *bool
	var versionAlerts *bool
	err := i.QueryRow(
		`
      SELECT services.id, services.name, services.description, services.logo, services.url, service_configs.favorite, service_configs.hidden, service_configs.health_alerts, service_configs.version_alerts
      FROM services
      LEFT JOIN service_configs
      ON services.id = service_configs.service_id AND service_configs.user_name = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{userName},
			wherePlaceholders,
		)...,
	).Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favorite, &hidden, &healthAlerts, &versionAlerts)
	if favorite != nil {
		record.Favorite = *favorite
	}
	if hidden != nil {
		record.Hidden = *hidden
	}
	if healthAlerts != nil {
		record.HealthAlerts = *healthAlerts
	}
	if versionAlerts != nil {
		record.VersionAlerts = *versionAlerts
	}
	if errors.Is(err, sql.ErrNoRows) {
		return schema.UserService{}, schema.ErrNotFound
	}
	return record, err
}
