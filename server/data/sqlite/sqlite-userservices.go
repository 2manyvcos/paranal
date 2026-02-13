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
      SELECT services.id, services.name, services.description, services.logo, services.url, serviceconfigs.favorite, serviceconfigs.hidden, serviceconfigs.uptimeAlerts, serviceconfigs.versionAlerts
      FROM services
      LEFT JOIN serviceconfigs
      ON serviceconfigs.serviceID = services.id AND serviceconfigs.userName = ?
      WHERE `+where+`
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
		var uptimeAlerts *bool
		var versionAlerts *bool
		if err := rows.Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favorite, &hidden, &uptimeAlerts, &versionAlerts); err != nil {
			return nil, err
		}
		if favorite != nil {
			record.Favorite = *favorite
		}
		if hidden != nil {
			record.Hidden = *hidden
		}
		if uptimeAlerts != nil {
			record.UptimeAlerts = *uptimeAlerts
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
	var uptimeAlerts *bool
	var versionAlerts *bool
	err := i.QueryRow(
		`
      SELECT services.id, services.name, services.description, services.logo, services.url, serviceconfigs.favorite, serviceconfigs.hidden, serviceconfigs.uptimeAlerts, serviceconfigs.versionAlerts
      FROM services
      LEFT JOIN serviceconfigs
      ON services.id = serviceconfigs.serviceID AND serviceconfigs.userName = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{userName},
			wherePlaceholders,
		)...,
	).Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favorite, &hidden, &uptimeAlerts, &versionAlerts)
	if favorite != nil {
		record.Favorite = *favorite
	}
	if hidden != nil {
		record.Hidden = *hidden
	}
	if uptimeAlerts != nil {
		record.UptimeAlerts = *uptimeAlerts
	}
	if versionAlerts != nil {
		record.VersionAlerts = *versionAlerts
	}
	if errors.Is(err, sql.ErrNoRows) {
		return schema.UserService{}, schema.ErrNotFound
	}
	return record, err
}
