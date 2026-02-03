package data

type Table[Dataset any] interface {
	TableName() string
	RecordFieldNames() []string
	NewRecord() TableRecord[Dataset]
}

type TableRecord[Dataset any] interface {
	RecordFields() []any
	Dataset() Dataset
}

type Identifier[Dataset any] interface {
	Table() Table[Dataset]
	Valid() error
	IDNames() []string
	IDs() []any
}

type Insertable[Dataset any] interface {
	Table() Table[Dataset]
	Valid() error
	InsertableNames() []string
	Insertables() []any
}

type Updatable[Dataset any] interface {
	Table() Table[Dataset]
	Valid() error
	UpdatableNames() []string
	Updatables() []any
}

type Upsertable[Dataset any] interface {
	Insertable[Dataset]
	Updatable[Dataset]
}

type TableHeader struct {
	Name   string
	Fields []string
}

func (t TableHeader) TableName() string {
	return t.Name
}

func (t TableHeader) RecordFieldNames() []string {
	return t.Fields
}
