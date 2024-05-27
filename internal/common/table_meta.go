package common

type ColumnPair struct {
	ColumnName string
	ColumnType string
}

type TableMeta struct {
	TableName string
	Columns   []ColumnPair
}
