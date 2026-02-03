package data

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	_ "modernc.org/sqlite"
)

func sqliteListJoinedDatasets[Dataset any](p *sqliteImpl, table JoinedTable[Dataset]) ([]Dataset, error) {
	leftFields := table.LeftRecordFieldNames()
	rightFields := table.RightRecordFieldNames()
	fields := make([]string, 0, len(leftFields)+len(rightFields))
	for _, name := range leftFields {
		fields = append(fields, "l."+name)
	}
	for _, name := range rightFields {
		fields = append(fields, "r."+name)
	}
	correlations := table.TableCorrelations()
	correlationConditions := make([]string, len(correlations))
	for i, names := range correlations {
		correlationConditions[i] = fmt.Sprintf("l.%s = r.%s", names[0], names[1])
	}
	statement := fmt.Sprintf(
		"SELECT %s FROM %s l INNER JOIN %s r ON %s",
		strings.Join(fields, ", "),
		table.LeftTableName(),
		table.RightTableName(),
		strings.Join(correlationConditions, " AND "),
	)
	rows, err := p.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Dataset
	for rows.Next() {
		record := table.NewJoinedRecord()
		if err := rows.Scan(slices.Concat(record.LeftRecordFields(), record.RightRecordFields())...); err != nil {
			return nil, err
		}
		results = append(results, record.Dataset())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func sqliteGetJoinedDatasets[Dataset any](p *sqliteImpl, table JoinedTable[Dataset], id JoinedIdentifier[Dataset]) ([]Dataset, error) {
	if err := id.Valid(); err != nil {
		return nil, err
	}
	leftFields := table.LeftRecordFieldNames()
	rightFields := table.RightRecordFieldNames()
	fields := make([]string, 0, len(leftFields)+len(rightFields))
	for _, name := range leftFields {
		fields = append(fields, "l."+name)
	}
	for _, name := range rightFields {
		fields = append(fields, "r."+name)
	}
	correlations := table.TableCorrelations()
	correlationConditions := make([]string, len(correlations))
	for i, names := range correlations {
		correlationConditions[i] = fmt.Sprintf("l.%s = r.%s", names[0], names[1])
	}
	leftIDs := id.LeftIDNames()
	rightIDs := id.RightIDNames()
	ids := make([]string, 0, len(leftIDs)+len(rightIDs))
	for _, name := range leftIDs {
		ids = append(ids, "l."+name)
	}
	for _, name := range rightIDs {
		ids = append(ids, "r."+name)
	}
	conditions := make([]string, len(ids))
	for i, name := range ids {
		conditions[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"SELECT %s FROM %s l INNER JOIN %s r ON %s WHERE %s",
		strings.Join(fields, ", "),
		table.LeftTableName(),
		table.RightTableName(),
		strings.Join(correlationConditions, " AND "),
		strings.Join(conditions, " AND "),
	)
	rows, err := p.DB.Query(statement, slices.Concat(id.LeftIDs(), id.RightIDs())...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Dataset
	for rows.Next() {
		record := table.NewJoinedRecord()
		if err := rows.Scan(slices.Concat(record.LeftRecordFields(), record.RightRecordFields())...); err != nil {
			return nil, err
		}
		results = append(results, record.Dataset())
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func sqliteGetJoinedDataset[Dataset any](p *sqliteImpl, table JoinedTable[Dataset], id JoinedIdentifier[Dataset]) (Dataset, error) {
	if err := id.Valid(); err != nil {
		var dataset Dataset
		return dataset, err
	}
	leftFields := table.LeftRecordFieldNames()
	rightFields := table.RightRecordFieldNames()
	fields := make([]string, 0, len(leftFields)+len(rightFields))
	for _, name := range leftFields {
		fields = append(fields, "l."+name)
	}
	for _, name := range rightFields {
		fields = append(fields, "r."+name)
	}
	correlations := table.TableCorrelations()
	correlationConditions := make([]string, len(correlations))
	for i, names := range correlations {
		correlationConditions[i] = fmt.Sprintf("l.%s = r.%s", names[0], names[1])
	}
	leftIDs := id.LeftIDNames()
	rightIDs := id.RightIDNames()
	ids := make([]string, 0, len(leftIDs)+len(rightIDs))
	for _, name := range leftIDs {
		ids = append(ids, "l."+name)
	}
	for _, name := range rightIDs {
		ids = append(ids, "r."+name)
	}
	conditions := make([]string, len(ids))
	for i, name := range ids {
		conditions[i] = fmt.Sprintf("%s = ?", name)
	}
	statement := fmt.Sprintf(
		"SELECT %s FROM %s l INNER JOIN %s r ON %s WHERE %s",
		strings.Join(fields, ", "),
		table.LeftTableName(),
		table.RightTableName(),
		strings.Join(correlationConditions, " AND "),
		strings.Join(conditions, " AND "),
	)
	record := table.NewJoinedRecord()
	err := p.DB.QueryRow(statement, slices.Concat(id.LeftIDs(), id.RightIDs())...).Scan(slices.Concat(record.LeftRecordFields(), record.RightRecordFields())...)
	if errors.Is(err, sql.ErrNoRows) {
		var dataset Dataset
		return dataset, ErrNotFound
	}
	return record.Dataset(), err
}
