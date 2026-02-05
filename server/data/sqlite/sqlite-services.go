package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec("CREATE TABLE IF NOT EXISTS services (id TEXT PRIMARY KEY, name TEXT, description TEXT, logo TEXT, url TEXT)")
		if err != nil {
			return fmt.Errorf("creating table \"services\" failed - %s", err)
		}
		return nil
	})
}

func ServiceQuery(query *schema.ServiceQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceQuery)
	}
	var conditions []string
	if query.ID != nil {
		conditions = append(conditions, "services.id = ?")
		placeholders = append(placeholders, *query.ID)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServicesWithFavorite(userName string, query *schema.ServiceQuery) ([]schema.ServiceWithFavorite, error) {
	where, wherePlaceholders := ServiceQuery(query)
	rows, err := i.Query(
		`
      SELECT services.id, services.name, services.description, services.logo, services.url, favorites.userName
      FROM services
      LEFT JOIN favorites
      ON services.id = favorites.serviceID AND favorites.userName = ?
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
	var records []schema.ServiceWithFavorite
	for rows.Next() {
		var record schema.ServiceWithFavorite
		var favoriteUserName *string
		if err := rows.Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favoriteUserName); err != nil {
			return nil, err
		}
		record.IsFavorite = favoriteUserName != nil
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetService(query schema.ServiceQuery) (schema.Service, error) {
	where, wherePlaceholders := ServiceQuery(&query)
	var record schema.Service
	err := i.QueryRow(
		`
      SELECT id, name, description, logo, url
      FROM services
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.Service{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) GetServiceWithFavorite(userName string, query schema.ServiceQuery) (schema.ServiceWithFavorite, error) {
	where, wherePlaceholders := ServiceQuery(&query)
	var record schema.ServiceWithFavorite
	var favoriteUserName *string
	err := i.QueryRow(
		`
      SELECT services.id, services.name, services.description, services.logo, services.url, favorites.userName
      FROM services
      LEFT JOIN favorites
      ON services.id = favorites.serviceID AND favorites.userName = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{userName},
			wherePlaceholders,
		)...,
	).Scan(&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL, &favoriteUserName)
	record.IsFavorite = favoriteUserName != nil
	if errors.Is(err, sql.ErrNoRows) {
		return schema.ServiceWithFavorite{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateService(record schema.Service) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO services (id, name, description, logo, url)
      VALUES (?, ?, ?, ?, ?)
    `,
		&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateService(query schema.ServiceQuery, record schema.Service) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := ServiceQuery(&query)
	result, err := i.Exec(
		`
      UPDATE services
      SET
        id = ?,
        name = ?,
        description = ?,
        logo = ?,
        url = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.ID, &record.Name, &record.Description, &record.Logo, &record.URL},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteService(query schema.ServiceQuery) error {
	where, wherePlaceholders := ServiceQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM services
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
