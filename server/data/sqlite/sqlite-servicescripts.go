package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS servicescripts (
        id TEXT PRIMARY KEY NOT NULL,
        serviceID TEXT NOT NULL,
        name TEXT NOT NULL,
        schedule TEXT NOT NULL,
        source TEXT NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"servicescripts\" failed - %s", err)
		}
		return nil
	})
}

func ServiceScriptQuery(query *schema.ServiceScriptQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceScriptQuery)
	}
	var conditions []string
	if query.ID != nil {
		conditions = append(conditions, "servicescripts.id = ?")
		placeholders = append(placeholders, *query.ID)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "servicescripts.serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if query.Name != nil {
		conditions = append(conditions, "servicescripts.name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if query.Schedule != nil {
		conditions = append(conditions, "servicescripts.schedule = ?")
		placeholders = append(placeholders, *query.Schedule)
	}
	if query.Source != nil {
		conditions = append(conditions, "servicescripts.source = ?")
		placeholders = append(placeholders, *query.Source)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceScripts(query *schema.ServiceScriptQuery) ([]schema.ServiceScript, error) {
	where, wherePlaceholders := ServiceScriptQuery(query)
	rows, err := i.Query(
		`
      SELECT servicescripts.id, servicescripts.serviceID, servicescripts.name, servicescripts.schedule, servicescripts.source
      FROM servicescripts
      INNER JOIN services
      ON servicescripts.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.ServiceScript
	for rows.Next() {
		var record schema.ServiceScript
		if err := rows.Scan(&record.ID, &record.ServiceID, &record.Name, &record.Schedule, &record.Source); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetServiceScript(query schema.ServiceScriptQuery) (schema.ServiceScript, error) {
	where, wherePlaceholders := ServiceScriptQuery(&query)
	var record schema.ServiceScript
	err := i.QueryRow(
		`
      SELECT servicescripts.id, servicescripts.serviceID, servicescripts.name, servicescripts.schedule, servicescripts.source
      FROM servicescripts
      INNER JOIN services
      ON servicescripts.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.ID, &record.ServiceID, &record.Name, &record.Schedule, &record.Source)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.ServiceScript{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateServiceScript(record schema.ServiceScript) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO servicescripts (id, serviceID, name, schedule, source)
      VALUES (?, ?, ?, ?, ?)
    `,
		&record.ID, &record.ServiceID, &record.Name, &record.Schedule, &record.Source,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateServiceScripts(query schema.ServiceScriptQuery, record schema.ServiceScript) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := ServiceScriptQuery(&query)
	result, err := i.Exec(
		`
      UPDATE servicescripts
      SET
        id = ?,
        serviceID = ?,
        name = ?,
        schedule = ?,
        source = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.ID, &record.ServiceID, &record.Name, &record.Schedule, &record.Source},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteServiceScripts(query schema.ServiceScriptQuery) error {
	where, wherePlaceholders := ServiceScriptQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM servicescripts
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
