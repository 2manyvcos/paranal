package sqlite

import (
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func UserServiceQuery(query *schema.UserServiceQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.UserServiceQuery)
	}
	var conditions []string
	if query.ID != nil {
		conditions = append(conditions, "services.id = ?")
		placeholders = append(placeholders, *query.ID)
	}
	if query.Favorite != nil {
		conditions = append(conditions, "serviceuserconfigs.favorite = ?")
		placeholders = append(placeholders, *query.Favorite)
	}
	if query.UptimeAlert != nil {
		conditions = append(conditions, "serviceuserconfigs.uptimeAlert = ?")
		placeholders = append(placeholders, *query.UptimeAlert)
	}
	if query.VersionAlert != nil {
		conditions = append(conditions, "serviceuserconfigs.versionAlert = ?")
		placeholders = append(placeholders, *query.VersionAlert)
	}
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
      SELECT services.id, services.name, services.description, services.logo, services.url, serviceuserconfigs.favorite, serviceuserconfigs.uptimeAlert, serviceuserconfigs.versionAlert
      FROM services
      LEFT JOIN serviceuserconfigs
      ON services.id = serviceuserconfigs.serviceID AND serviceuserconfigs.userName = ?
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
		var uptimeAlert *bool
		var versionAlert *bool
		if err := rows.Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favorite, &uptimeAlert, &versionAlert); err != nil {
			return nil, err
		}
		if favorite != nil {
			record.Favorite = *favorite
		}
		if uptimeAlert != nil {
			record.UptimeAlert = *uptimeAlert
		}
		if versionAlert != nil {
			record.VersionAlert = *versionAlert
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
	var uptimeAlert *bool
	var versionAlert *bool
	err := i.QueryRow(
		`
      SELECT services.id, services.name, services.description, services.logo, services.url, serviceuserconfigs.favorite, serviceuserconfigs.uptimeAlert, serviceuserconfigs.versionAlert
      FROM services
      LEFT JOIN serviceuserconfigs
      ON services.id = serviceuserconfigs.serviceID AND serviceuserconfigs.userName = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{userName},
			wherePlaceholders,
		)...,
	).Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favorite, &uptimeAlert, &versionAlert)
	if favorite != nil {
		record.Favorite = *favorite
	}
	if uptimeAlert != nil {
		record.UptimeAlert = *uptimeAlert
	}
	if versionAlert != nil {
		record.VersionAlert = *versionAlert
	}
	if errors.Is(err, sql.ErrNoRows) {
		return schema.UserService{}, schema.ErrNotFound
	}
	return record, err
}
