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
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS serviceversions (
        name TEXT NOT NULL,
        serviceID TEXT NOT NULL,
        version TEXT NOT NULL,
        PRIMARY KEY (name, serviceID)
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"serviceversions\" failed - %s", err)
		}
		return nil
	})
}

func ServiceVersionQuery(query *schema.ServiceVersionQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceVersionQuery)
	}
	var conditions []string
	if query.Name != nil {
		conditions = append(conditions, "serviceversions.name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "serviceversions.serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceVersions(query *schema.ServiceVersionQuery) ([]schema.ServiceVersion, error) {
	where, wherePlaceholders := ServiceVersionQuery(query)
	rows, err := i.Query(
		`
      SELECT serviceversions.name, serviceversions.serviceID, serviceversions.version
      FROM serviceversions
      INNER JOIN services
      ON serviceversions.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.ServiceVersion
	for rows.Next() {
		var record schema.ServiceVersion
		if err := rows.Scan(&record.Name, &record.ServiceID, &record.Version); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetServiceVersion(query schema.ServiceVersionQuery) (schema.ServiceVersion, error) {
	where, wherePlaceholders := ServiceVersionQuery(&query)
	var record schema.ServiceVersion
	err := i.QueryRow(
		`
      SELECT serviceversions.name, serviceversions.serviceID, serviceversions.version
      FROM serviceversions
      INNER JOIN services
      ON serviceversions.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.ServiceID, &record.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.ServiceVersion{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateServiceVersion(record schema.ServiceVersion) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO serviceversions (name, serviceID, version)
      VALUES (?, ?, ?)
    `,
		&record.Name, &record.ServiceID, &record.Version,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateServiceVersions(query schema.ServiceVersionQuery, record schema.ServiceVersion) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := ServiceVersionQuery(&query)
	result, err := i.Exec(
		`
      UPDATE serviceversions
      SET
        name = ?,
        serviceID = ?,
        version = ?
        WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.ServiceID, &record.Version},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteServiceVersions(query schema.ServiceVersionQuery) error {
	where, wherePlaceholders := ServiceVersionQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM serviceversions
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
