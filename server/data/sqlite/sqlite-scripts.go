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
		_, err := i.Exec("CREATE TABLE IF NOT EXISTS scripts (id INTEGER PRIMARY KEY, name TEXT, schedule TEXT, source TEXT, serviceID INTEGER)")
		if err != nil {
			return fmt.Errorf("creating table \"scripts\" failed - %s", err)
		}
		return nil
	})
}

func ScriptQuery(query *schema.ScriptQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ScriptQuery)
	}
	var conditions []string
	if query.ID != nil {
		conditions = append(conditions, "scripts.id = ?")
		placeholders = append(placeholders, *query.ID)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListScripts(query *schema.ScriptQuery) ([]schema.Script, error) {
	where, wherePlaceholders := ScriptQuery(query)
	rows, err := i.Query(
		`
      SELECT scripts.id, scripts.name, scripts.schedule, scripts.source, scripts.serviceID
      FROM scripts
      INNER JOIN services
      ON scripts.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.Script
	for rows.Next() {
		var record schema.Script
		if err := rows.Scan(&record.ID, &record.Name, &record.Schedule, &record.Source, &record.ServiceID); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetScript(query schema.ScriptQuery) (schema.Script, error) {
	where, wherePlaceholders := ScriptQuery(&query)
	var record schema.Script
	err := i.QueryRow(
		`
      SELECT scripts.id, scripts.name, scripts.schedule, scripts.source, scripts.serviceID
      FROM scripts
      INNER JOIN services
      ON scripts.serviceID = services.id
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.ID, &record.Name, &record.Schedule, &record.Source, &record.ServiceID)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.Script{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateScript(record schema.Script) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO scripts (id, name, schedule, source, serviceID)
      VALUES (?, ?, ?, ?, ?)
    `,
		&record.ID, &record.Name, &record.Schedule, &record.Source, &record.ServiceID,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateScript(query schema.ScriptQuery, record schema.Script) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := ScriptQuery(&query)
	result, err := i.Exec(
		`
      UPDATE scripts
      SET
        id = ?,
        name = ?,
        schedule = ?,
        source = ?,
        serviceID = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.ID, &record.Name, &record.Schedule, &record.Source, &record.ServiceID},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteScript(query schema.ScriptQuery) error {
	where, wherePlaceholders := ScriptQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM scripts
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
