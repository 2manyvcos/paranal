package data

type JoinType int

const (
	JoinTypeInner JoinType = iota + 1
	JoinTypeLeft
)

var JoinTypeVerbs = map[JoinType]string{
	JoinTypeInner: "INNER",
	JoinTypeLeft:  "LEFT",
}

type JoinedTable[Dataset any] interface {
	JoinType() JoinType
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

type JoinedConditions struct {
	LeftCondition  *Condition
	RightCondition *Condition
	Conditions     []JoinedConditions
	Or             bool
}

type JoinedTableHeader struct {
	Type        JoinType
	LeftName    string
	LeftFields  []string
	RightName   string
	RightFields []string
}

func (t JoinedTableHeader) JoinType() JoinType {
	return t.Type
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
