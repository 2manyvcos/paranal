package data

import (
	"fmt"
)

type ScriptTable struct {
	JoinedTableHeader
}

var Scripts CombinedTable[Script] = ScriptTable{
	JoinedTableHeader{
		LeftName:    "scripts",
		LeftFields:  scriptFields,
		RightName:   Services.TableName(),
		RightFields: nil,
	},
}

func (t ScriptTable) TableCorrelations() [][2]string {
	return ScriptCorrelations
}

var ScriptCorrelations = [][2]string{{"serviceID", "id"}}

func (t ScriptTable) NewJoinedRecord() JoinedTableRecord[Script] {
	return new(Script)
}

func (t ScriptTable) TableName() string {
	return t.LeftName
}

func (t ScriptTable) RecordFieldNames() []string {
	return t.LeftRecordFieldNames()
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

func (r *Script) LeftRecordFields() []any {
	return []any{&r.ID, &r.Name, &r.Schedule, &r.Source, &r.ServiceID}
}

func (r *Script) RightRecordFields() []any {
	return nil
}

func (r *Script) RecordFields() []any {
	return r.LeftRecordFields()
}

func (r *Script) Dataset() Script {
	return *r
}

var _ Upsertable[Script] = Script{}

func (r Script) Table() Table[Script] {
	return Scripts
}

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

type ScriptName string

var _ JoinedIdentifier[Script] = ScriptName("")
var _ Identifier[Script] = ScriptName("")

func (i ScriptName) JoinedTable() JoinedTable[Script] {
	return Scripts
}

func (i ScriptName) Table() Table[Script] {
	return Scripts
}

func (i ScriptName) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (i ScriptName) LeftIDNames() []string {
	return ScriptNameIDs
}

var ScriptNameIDs = []string{"name"}

func (i ScriptName) LeftIDs() []any {
	return []any{string(i)}
}
func (i ScriptName) RightIDNames() []string {
	return nil
}

func (i ScriptName) RightIDs() []any {
	return nil
}

func (i ScriptName) IDNames() []string {
	return i.LeftIDNames()
}

func (i ScriptName) IDs() []any {
	return i.LeftIDs()
}
