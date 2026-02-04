package data

import (
	"fmt"
)

type ScriptTable struct{ TableHeader }

var Scripts Table[Script] = ScriptTable{
	TableHeader{
		Name:   "scripts",
		Fields: scriptFields,
	},
}

func (t ScriptTable) NewRecord() TableRecord[Script] {
	return new(Script)
}

type Script struct {
	ID        int
	Name      string
	Schedule  string
	Source    string
	ServiceID string
}

var scriptFields = []string{"id", "name", "schedule", "source", "serviceID"}

func (r *Script) RecordFields() []any {
	return []any{&r.ID, &r.Name, &r.Schedule, &r.Source, &r.ServiceID}
}

func (r *Script) Dataset() Script {
	return *r
}

var _ Upsertable[Script] = Script{}

func (r Script) TableType() Table[Script] { return nil }

func (r Script) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if r.Schedule == "" {
		return fmt.Errorf("invalid schedule")
	}
	if r.Source == "" {
		return fmt.Errorf("invalid source")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}

func (r Script) InsertableNames() []string {
	return scriptInsertables
}

var scriptInsertables = []string{"name", "schedule", "source", "serviceID"}

func (r Script) Insertables() []any {
	return []any{r.Name, r.Schedule, r.Source, r.ServiceID}
}

func (r Script) UpdatableNames() []string {
	return scriptUpdatables
}

var scriptUpdatables = []string{"name", "schedule", "source", "serviceID"}

func (r Script) Updatables() []any {
	return []any{r.Name, r.Schedule, r.Source, r.ServiceID}
}

type ScriptTableInnerJoinServices struct{ JoinedTableHeader }

var ScriptsInnerJoinServices JoinedTable[Script] = ScriptTableInnerJoinServices{
	JoinedTableHeader{
		Type:        JoinTypeInner,
		LeftName:    Scripts.TableName(),
		LeftFields:  Scripts.RecordFieldNames(),
		RightName:   Services.TableName(),
		RightFields: nil,
	},
}

func (t ScriptTableInnerJoinServices) TableCorrelations() [][2]string {
	return ScriptServiceCorrelations
}

var ScriptServiceCorrelations = [][2]string{{"serviceID", "id"}}

func (t ScriptTableInnerJoinServices) NewJoinedRecord() JoinedTableRecord[Script] {
	return new(Script)
}

func (r *Script) LeftRecordFields() []any {
	return r.RecordFields()
}

func (r *Script) RightRecordFields() []any {
	return nil
}
