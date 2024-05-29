package common

type ColumnPair struct {
	ColumnName      string
	ColumnType      string
	IsAutoIncrement bool
	IsPrimaryKey    bool
}

type TableMeta struct {
	TableName    string
	Columns      []ColumnPair
	PkColumnName string
}
