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

type ScriptID int

var _ Identifier[Script] = ScriptID(0)
var _ JoinedIdentifier[Script] = ScriptID(0)

func (i ScriptID) Table() Table[Script] {
	return Scripts
}

func (i ScriptID) JoinedTable() JoinedTable[Script] {
	return Scripts
}

func (i ScriptID) Valid() error {
	return nil
}

func (i ScriptID) IDNames() []string {
	return ScriptIDIDs
}

var ScriptIDIDs = []string{"id"}

func (i ScriptID) IDs() []any {
	return []any{int(i)}
}

func (i ScriptID) LeftIDNames() []string {
	return i.IDNames()
}

func (i ScriptID) LeftIDs() []any {
	return i.IDs()
}
func (i ScriptID) RightIDNames() []string {
	return nil
}

func (i ScriptID) RightIDs() []any {
	return nil
}

type ScriptServiceID string

var _ Identifier[Script] = ScriptServiceID("")
var _ JoinedIdentifier[Script] = ScriptServiceID("")

func (i ScriptServiceID) Table() Table[Script] {
	return Scripts
}

func (i ScriptServiceID) JoinedTable() JoinedTable[Script] {
	return Scripts
}

func (i ScriptServiceID) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}

func (i ScriptServiceID) IDNames() []string {
	return ScriptServiceIDIDs
}

var ScriptServiceIDIDs = []string{"serviceID"}

func (i ScriptServiceID) IDs() []any {
	return []any{string(i)}
}

func (i ScriptServiceID) LeftIDNames() []string {
	return i.IDNames()
}

func (i ScriptServiceID) LeftIDs() []any {
	return i.IDs()
}
func (i ScriptServiceID) RightIDNames() []string {
	return nil
}

func (i ScriptServiceID) RightIDs() []any {
	return nil
}

type ScriptIDAndServiceID struct {
	ID        int
	ServiceID string
}

var _ Identifier[Script] = ScriptIDAndServiceID{}
var _ JoinedIdentifier[Script] = ScriptIDAndServiceID{}

func (i ScriptIDAndServiceID) Table() Table[Script] {
	return Scripts
}

func (i ScriptIDAndServiceID) JoinedTable() JoinedTable[Script] {
	return Scripts
}

func (i ScriptIDAndServiceID) Valid() error {
	if i.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}

func (i ScriptIDAndServiceID) IDNames() []string {
	return ScriptIDAndServiceIDIDs
}

var ScriptIDAndServiceIDIDs = []string{"id", "serviceID"}

func (i ScriptIDAndServiceID) IDs() []any {
	return []any{i.ID, i.ServiceID}
}

func (i ScriptIDAndServiceID) LeftIDNames() []string {
	return i.IDNames()
}

func (i ScriptIDAndServiceID) LeftIDs() []any {
	return i.IDs()
}
func (i ScriptIDAndServiceID) RightIDNames() []string {
	return nil
}

func (i ScriptIDAndServiceID) RightIDs() []any {
	return nil
}
