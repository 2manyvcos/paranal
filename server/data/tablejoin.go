package data

type JoinedTable[Dataset any] interface {
	LeftTableName() string
	LeftRecordFieldNames() []string
	RightTableName() string
	RightRecordFieldNames() []string
	TableCorrelations() [][2]string
	NewJoinedRecord() JoinedTableRecord[Dataset]
}

type JoinedTableRecord[Dataset any] interface {
	LeftRecordFields() []any
	RightRecordFields() []any
	Dataset() Dataset
}

type JoinedIdentifier[Dataset any] interface {
	JoinedTable() JoinedTable[Dataset]
	Valid() error
	LeftIDNames() []string
	LeftIDs() []any
	RightIDNames() []string
	RightIDs() []any
}

type JoinedTableHeader struct {
	LeftName    string
	LeftFields  []string
	RightName   string
	RightFields []string
}

func (t JoinedTableHeader) LeftTableName() string {
	return t.LeftName
}

func (t JoinedTableHeader) LeftRecordFieldNames() []string {
	return t.LeftFields
}

func (t JoinedTableHeader) RightTableName() string {
	return t.RightName
}

func (t JoinedTableHeader) RightRecordFieldNames() []string {
	return t.RightFields
}

type CombinedTable[Dataset any] interface {
	JoinedTable[Dataset]
	Table[Dataset]
}
