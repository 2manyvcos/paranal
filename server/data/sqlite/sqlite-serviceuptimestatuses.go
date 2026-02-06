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
      CREATE TABLE IF NOT EXISTS serviceuptimestatuses (
        name TEXT NOT NULL,
        serviceID TEXT NOT NULL,
        status INTEGER NOT NULL,
        PRIMARY KEY (name, serviceID)
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"serviceuptimestatuses\" failed - %s", err)
		}
		return nil
	})
}

func ServiceUptimeStatusQuery(query *schema.ServiceUptimeStatusQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceUptimeStatusQuery)
	}
	var conditions []string
	if query.Name != nil {
		conditions = append(conditions, "serviceuptimestatuses.name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "serviceuptimestatuses.serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceUptimeStatuses(query *schema.ServiceUptimeStatusQuery) ([]schema.ServiceUptimeStatus, error) {
	where, wherePlaceholders := ServiceUptimeStatusQuery(query)
	rows, err := i.Query(
		`
      SELECT serviceuptimestatuses.name, serviceuptimestatuses.serviceID, serviceuptimestatuses.status
      FROM serviceuptimestatuses
      INNER JOIN services
      ON serviceuptimestatuses.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.ServiceUptimeStatus
	for rows.Next() {
		var record schema.ServiceUptimeStatus
		if err := rows.Scan(&record.Name, &record.ServiceID, &record.Status); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetServiceUptimeStatus(query schema.ServiceUptimeStatusQuery) (schema.ServiceUptimeStatus, error) {
	where, wherePlaceholders := ServiceUptimeStatusQuery(&query)
	var record schema.ServiceUptimeStatus
	err := i.QueryRow(
		`
      SELECT serviceuptimestatuses.name, serviceuptimestatuses.serviceID, serviceuptimestatuses.status
      FROM serviceuptimestatuses
      INNER JOIN services
      ON serviceuptimestatuses.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.ServiceID, &record.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.ServiceUptimeStatus{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateServiceUptimeStatus(record schema.ServiceUptimeStatus) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO serviceuptimestatuses (name, serviceID, status)
      VALUES (?, ?, ?)
    `,
		&record.Name, &record.ServiceID, &record.Status,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateServiceUptimeStatuses(query schema.ServiceUptimeStatusQuery, record schema.ServiceUptimeStatus) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := ServiceUptimeStatusQuery(&query)
	result, err := i.Exec(
		`
      UPDATE serviceuptimestatuses
      SET
        name = ?,
        serviceID = ?,
        status = ?
        WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.ServiceID, &record.Status},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteServiceUptimeStatuses(query schema.ServiceUptimeStatusQuery) error {
	where, wherePlaceholders := ServiceUptimeStatusQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM serviceuptimestatuses
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
