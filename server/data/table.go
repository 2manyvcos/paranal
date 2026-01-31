package data

type Dataset interface {
	Valid() bool
	IDs() []any
	Fields() []any
}

type Record[DatasetType Dataset] interface {
	IDPointers() []any
	FieldPointers() []any
	Dataset() DatasetType
}

type Table[RecordType Record[DatasetType], DatasetType Dataset] interface {
	Name() string
	IDNames() []string
	IDsValid(id ...any) error
	FieldNames() []string
	NewRecord() RecordType
}

type TableHeader struct {
	Table       string
	IDs, Fields []string
}

func (t TableHeader) Name() string {
	return t.Table
}

func (t TableHeader) IDNames() []string {
	return t.IDs
}

func (t TableHeader) FieldNames() []string {
	return t.Fields
}
