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

func sqliteListDatasets[Dataset any](p *sqliteImpl, table Table[Dataset]) ([]Dataset, error) {
	statement := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(table.RecordFieldNames(), ", "),
		table.TableName(),
	)
	rows, err := p.DB.Query(statement)
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

func sqliteGetDatasets[Dataset any](p *sqliteImpl, table Table[Dataset], id Identifier[Dataset]) ([]Dataset, error) {
	if err := id.Valid(); err != nil {
		return nil, err
	}
	ids := id.IDNames()
	conditions := make([]string, len(ids))
	for i, name := range ids {
		conditions[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(table.RecordFieldNames(), ", "),
		table.TableName(),
		strings.Join(conditions, " AND "),
	)
	rows, err := p.DB.Query(statement, id.IDs()...)
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

func sqliteGetDataset[Dataset any](p *sqliteImpl, table Table[Dataset], id Identifier[Dataset]) (Dataset, error) {
	if err := id.Valid(); err != nil {
		var dataset Dataset
		return dataset, err
	}
	ids := id.IDNames()
	conditions := make([]string, len(ids))
	for i, name := range ids {
		conditions[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(table.RecordFieldNames(), ", "),
		table.TableName(),
		strings.Join(conditions, " AND "),
	)
	record := table.NewRecord()
	err := p.DB.QueryRow(statement, id.IDs()...).Scan(record.RecordFields()...)
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

func sqliteUpdateDataset[Dataset any](p *sqliteImpl, table Table[Dataset], id Identifier[Dataset], dataset Updatable[Dataset]) error {
	if err := id.Valid(); err != nil {
		return err
	}
	if err := dataset.Valid(); err != nil {
		return err
	}
	fields := dataset.UpdatableNames()
	fieldMappings := make([]string, len(fields))
	for i, name := range fields {
		fieldMappings[i] = fmt.Sprintf("%s = ?", name)
	}
	ids := id.IDNames()
	conditions := make([]string, len(ids))
	for i, name := range ids {
		conditions[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table.TableName(),
		strings.Join(fieldMappings, ", "),
		strings.Join(conditions, " AND "),
	)
	result, err := p.DB.Exec(statement, slices.Concat(dataset.Updatables(), id.IDs())...)
	if err != nil {
		return err
	}
	return requireFound(result)
}

func sqliteDeleteDataset[Dataset any](p *sqliteImpl, table Table[Dataset], id Identifier[Dataset]) error {
	if err := id.Valid(); err != nil {
		return err
	}
	ids := id.IDNames()
	conditions := make([]string, len(ids))
	for i, name := range ids {
		conditions[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		table.TableName(),
		strings.Join(conditions, " AND "),
	)
	result, err := p.DB.Exec(statement, id.IDs()...)
	if err != nil {
		return err
	}
	return requireFound(result)
}
