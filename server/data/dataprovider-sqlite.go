package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type sqliteImpl struct {
	*sql.DB
}

func SQLite(path string) (DataProvider, error) {
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	provider := &sqliteImpl{db}
	return provider, provider.Setup()
}

func (p *sqliteImpl) Setup() error {
	log.Println("Setting up SQLite data backend")

	// Migration strategy as taken from a stackoverflow comment:

	// > For tracking the database version, I use the built in user-version variable that sqlite provides (sqlite does nothing with this variable, you are free to use it however you please). It starts at 0, and you can get/set this variable with the following sqlite statements:
	// >
	// > > PRAGMA user_version;
	// > > PRAGMA user_version = 1;
	// >
	// > When the app starts, I check the current user-version, apply any changes that are needed to bring the schema up to date, and then update the user-version. I wrap the updates in a transaction so that if anything goes wrong, the changes aren't committed.
	// >
	// > For making schema changes, sqlite supports "ALTER TABLE" syntax for certain operations (renaming the table or adding a column). This is an easy way to update existing tables in-place. See the documentation here: http://www.sqlite.org/lang_altertable.html. For deleting columns or other changes that aren't supported by the "ALTER TABLE" syntax, I create a new table, migrate data into it, drop the old table, and rename the new table to the original name.

	// var dbVersion int
	// err = p.DB.QueryRow("PRAGMA user_version").Scan(&dbVersion)

	if err := p.setupUsers(); err != nil {
		return err
	}

	if err := p.setupHTTPCredentials(); err != nil {
		return err
	}

	if err := p.setupSSHCredentials(); err != nil {
		return err
	}

	if err := p.setupUserCredentials(); err != nil {
		return err
	}

	return nil
}

func requireFound(result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func sqliteListDatasets[TableType Table[RecordType, DatasetType], RecordType Record[DatasetType], DatasetType Dataset](p *sqliteImpl, table TableType) ([]DatasetType, error) {
	tableName := table.Name()
	idNames := table.IDNames()
	fieldNames := table.FieldNames()
	fields := slices.Concat(idNames, fieldNames)

	statement := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(fields, ", "),
		tableName,
	)
	rows, err := p.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DatasetType
	for rows.Next() {
		record := table.NewRecord()
		if err := rows.Scan(slices.Concat(record.IDPointers(), record.FieldPointers())...); err != nil {
			return nil, err
		}
		result = append(result, record.Dataset())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func sqliteGetDataset[TableType Table[RecordType, DatasetType], RecordType Record[DatasetType], DatasetType Dataset](p *sqliteImpl, table TableType, ids ...any) (DatasetType, error) {
	if err := table.IDsValid(ids...); err != nil {
		var dataset DatasetType
		return dataset, err
	}

	tableName := table.Name()
	idNames := table.IDNames()
	fieldNames := table.FieldNames()
	fields := slices.Concat(idNames, fieldNames)

	idConditions := make([]string, len(idNames))
	for i, idName := range idNames {
		idConditions[i] = fmt.Sprintf("%s = ?", idName)
	}
	statement := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ", "),
		tableName,
		strings.Join(idConditions, " AND "),
	)
	record := table.NewRecord()
	err := p.DB.QueryRow(statement, ids...).Scan(slices.Concat(record.IDPointers(), record.FieldPointers())...)
	if errors.Is(err, sql.ErrNoRows) {
		var dataset DatasetType
		return dataset, ErrNotFound
	}
	return record.Dataset(), err
}

func sqliteCreateDataset[TableType Table[RecordType, DatasetType], RecordType Record[DatasetType], DatasetType Dataset](p *sqliteImpl, table TableType, dataset DatasetType, updateExisting bool) error {
	if err := dataset.Valid(); err != nil {
		return err
	}

	tableName := table.Name()
	idNames := table.IDNames()
	fieldNames := table.FieldNames()
	fields := slices.Concat(idNames, fieldNames)

	statement := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES(%s)",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(slices.Repeat([]string{"?"}, len(fields)), ", "),
	)
	if updateExisting && len(fieldNames) > 0 {
		fieldMappings := make([]string, len(fieldNames))
		for i, fieldName := range fieldNames {
			fieldMappings[i] = fmt.Sprintf("%s = excluded.%s", fieldName, fieldName)
		}
		statement += fmt.Sprintf(
			" ON CONFLICT(%s) DO UPDATE SET %s",
			strings.Join(idNames, ", "),
			strings.Join(fieldMappings, ", "),
		)
	}
	_, err := p.DB.Exec(statement, slices.Concat(dataset.IDs(), dataset.Fields())...)
	if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return ErrConflict
	}
	return err
}

func sqliteUpdateDataset[TableType Table[RecordType, DatasetType], RecordType Record[DatasetType], DatasetType Dataset](p *sqliteImpl, table TableType, dataset DatasetType) error {
	if err := dataset.Valid(); err != nil {
		return err
	}

	tableName := table.Name()
	idNames := table.IDNames()
	fieldNames := table.FieldNames()

	if len(fieldNames) == 0 {
		return nil
	}
	fieldMappings := make([]string, len(fieldNames))
	for i, fieldName := range fieldNames {
		fieldMappings[i] = fmt.Sprintf("%s = ?", fieldName)
	}
	idConditions := make([]string, len(idNames))
	for i, idName := range idNames {
		idConditions[i] = fmt.Sprintf("%s = ?", idName)
	}
	statement := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		tableName,
		strings.Join(fieldMappings, ", "),
		strings.Join(idConditions, " AND "),
	)
	result, err := p.DB.Exec(statement, slices.Concat(dataset.Fields(), dataset.IDs()))
	if err != nil {
		return err
	}
	return requireFound(result)
}

func sqliteDeleteDataset[TableType Table[RecordType, DatasetType], RecordType Record[DatasetType], DatasetType Dataset](p *sqliteImpl, table TableType, ids ...any) error {
	if err := table.IDsValid(ids...); err != nil {
		return err
	}

	tableName := table.Name()
	idNames := table.IDNames()

	idConditions := make([]string, len(idNames))
	for i, idName := range idNames {
		idConditions[i] = fmt.Sprintf("%s = ?", idName)
	}
	statement := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		tableName,
		strings.Join(idConditions, " AND "),
	)
	result, err := p.DB.Exec(statement, ids...)
	if err != nil {
		return err
	}
	return requireFound(result)
}
