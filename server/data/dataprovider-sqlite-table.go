package data

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func sqliteSelectDatasets[Dataset any](p *sqliteImpl, table Table[Dataset], where *Conditions) ([]Dataset, error) {
	statement := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(table.RecordFieldNames(), ", "),
		table.TableName(),
	)
	conditions, conditionPlaceholders, ok := sqliteResolveConditions(where)
	if ok {
		statement += " WHERE " + conditions
	}
	rows, err := p.DB.Query(statement, conditionPlaceholders...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Dataset
	for rows.Next() {
		record := table.NewRecord()
		if err := rows.Scan(record.RecordFields()...); err != nil {
			return nil, err
		}
		results = append(results, record.Dataset())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func sqliteSelectDataset[Dataset any](p *sqliteImpl, table Table[Dataset], where *Conditions) (Dataset, error) {
	statement := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(table.RecordFieldNames(), ", "),
		table.TableName(),
	)
	conditions, conditionPlaceholders, ok := sqliteResolveConditions(where)
	if ok {
		statement += " WHERE " + conditions
	}
	record := table.NewRecord()
	err := p.DB.QueryRow(statement, conditionPlaceholders...).Scan(record.RecordFields()...)
	if errors.Is(err, sql.ErrNoRows) {
		var dataset Dataset
		return dataset, ErrNotFound
	}
	return record.Dataset(), err
}

func sqliteCreateDataset[Dataset any](p *sqliteImpl, table Table[Dataset], dataset Insertable[Dataset]) error {
	if err := dataset.Valid(); err != nil {
		return err
	}
	fields := dataset.InsertableNames()
	statement := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES(%s)",
		table.TableName(),
		strings.Join(fields, ", "),
		strings.Join(slices.Repeat([]string{"?"}, len(fields)), ", "),
	)
	_, err := p.DB.Exec(statement, dataset.Insertables()...)
	if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return ErrConflict
	}
	return err
}

func sqliteCreateOrUpdateDataset[Dataset any](p *sqliteImpl, table Table[Dataset], dataset Upsertable[Dataset], idNames []string) error {
	if err := dataset.Valid(); err != nil {
		return err
	}
	insertables := dataset.InsertableNames()
	updatables := dataset.UpdatableNames()
	updateMappings := make([]string, len(updatables))
	for i, name := range updatables {
		updateMappings[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES(%s) ON CONFLICT(%s) DO UPDATE SET %s",
		table.TableName(),
		strings.Join(insertables, ", "),
		strings.Join(slices.Repeat([]string{"?"}, len(insertables)), ", "),
		strings.Join(idNames, ", "),
		strings.Join(updateMappings, ", "),
	)
	_, err := p.DB.Exec(statement, slices.Concat(dataset.Insertables(), dataset.Updatables())...)
	if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return ErrConflict
	}
	return err
}

func sqliteUpdateDataset[Dataset any](p *sqliteImpl, table Table[Dataset], where *Conditions, dataset Updatable[Dataset]) error {
	if err := dataset.Valid(); err != nil {
		return err
	}
	fields := dataset.UpdatableNames()
	fieldMappings := make([]string, len(fields))
	for i, name := range fields {
		fieldMappings[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"UPDATE %s SET %s",
		table.TableName(),
		strings.Join(fieldMappings, ", "),
	)
	conditions, conditionPlaceholders, ok := sqliteResolveConditions(where)
	if ok {
		statement += " WHERE " + conditions
	}
	result, err := p.DB.Exec(statement, slices.Concat(dataset.Updatables(), conditionPlaceholders)...)
	if err != nil {
		return err
	}
	return requireFound(result)
}

func sqliteDeleteDataset[Dataset any](p *sqliteImpl, table Table[Dataset], where *Conditions) error {
	statement := fmt.Sprintf(
		"DELETE FROM %s",
		table.TableName(),
	)
	conditions, conditionPlaceholders, ok := sqliteResolveConditions(where)
	if ok {
		statement += " WHERE " + conditions
	}
	result, err := p.DB.Exec(statement, conditionPlaceholders...)
	if err != nil {
		return err
	}
	return requireFound(result)
}

func sqliteResolveConditions(where *Conditions) (string, []any, bool) {
	if where == nil {
		return "", nil, false
	}
	conditions := make([]string, 0, 1+len(where.Conditions))
	placeholders := make([]any, 0, 1+len(where.Conditions))
	if where.Condition != nil {
		conditions = append(conditions, fmt.Sprintf("%s = ?", where.Condition.Field))
		placeholders = append(placeholders, where.Condition.Value)
	}
	if len(where.Conditions) > 0 {
		for _, sub := range where.Conditions {
			subConditions, subPlaceholders, _ := sqliteResolveConditions(&sub)
			conditions = append(conditions, subConditions)
			placeholders = append(placeholders, subPlaceholders...)
		}
	}
	merge := " AND "
	if where.Or {
		merge = " OR "
	}
	if len(conditions) == 0 {
		return "", nil, false
	}
	if len(conditions) == 1 {
		return conditions[0], placeholders, true
	}
	return "(" + strings.Join(conditions, merge) + ")", placeholders, len(conditions) > 0
}
