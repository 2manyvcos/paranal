package data

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	_ "modernc.org/sqlite"
)

func sqliteSelectJoinedDatasets[Dataset any](p *sqliteImpl, table JoinedTable[Dataset], on *JoinedConditions, where *JoinedConditions) ([]Dataset, error) {
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
		"SELECT %s FROM %s l %s JOIN %s r ON %s",
		strings.Join(fields, ", "),
		table.LeftTableName(),
		JoinTypeVerbs[table.JoinType()],
		table.RightTableName(),
		strings.Join(correlationConditions, " AND "),
	)
	onConditions, onConditionPlaceholders, ok := sqliteResolveJoinedConditions(on)
	if ok {
		statement += " AND " + onConditions
	}
	conditions, conditionPlaceholders, ok := sqliteResolveJoinedConditions(where)
	if ok {
		statement += " WHERE " + conditions
	}
	rows, err := p.DB.Query(statement, slices.Concat(onConditionPlaceholders, conditionPlaceholders)...)
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

func sqliteSelectJoinedDataset[Dataset any](p *sqliteImpl, table JoinedTable[Dataset], on *JoinedConditions, where *JoinedConditions) (Dataset, error) {
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
		"SELECT %s FROM %s l %s JOIN %s r ON %s",
		strings.Join(fields, ", "),
		table.LeftTableName(),
		JoinTypeVerbs[table.JoinType()],
		table.RightTableName(),
		strings.Join(correlationConditions, " AND "),
	)
	onConditions, onConditionPlaceholders, ok := sqliteResolveJoinedConditions(on)
	if ok {
		statement += " AND " + onConditions
	}
	conditions, conditionPlaceholders, ok := sqliteResolveJoinedConditions(where)
	if ok {
		statement += " WHERE " + conditions
	}
	record := table.NewJoinedRecord()
	err := p.DB.QueryRow(statement, slices.Concat(onConditionPlaceholders, conditionPlaceholders)...).Scan(slices.Concat(record.LeftRecordFields(), record.RightRecordFields())...)
	if errors.Is(err, sql.ErrNoRows) {
		var dataset Dataset
		return dataset, ErrNotFound
	}
	return record.Dataset(), err
}

func sqliteResolveJoinedConditions(where *JoinedConditions) (string, []any, bool) {
	if where == nil {
		return "", nil, false
	}
	conditions := make([]string, 0, 2+len(where.Conditions))
	placeholders := make([]any, 0, 2+len(where.Conditions))
	if where.LeftCondition != nil {
		conditions = append(conditions, fmt.Sprintf("l.%s = ?", where.LeftCondition.Field))
		placeholders = append(placeholders, where.LeftCondition.Value)
	}
	if where.RightCondition != nil {
		conditions = append(conditions, fmt.Sprintf("r.%s = ?", where.RightCondition.Field))
		placeholders = append(placeholders, where.RightCondition.Value)
	}
	if len(where.Conditions) > 0 {
		for _, sub := range where.Conditions {
			subConditions, subPlaceholders, _ := sqliteResolveJoinedConditions(&sub)
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
