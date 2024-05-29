package util

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
)

func IsID(columnPair common.ColumnPair) bool {
	return columnPair.IsAutoIncrement || columnPair.IsPrimaryKey
}

func IsNotID(columnPair common.ColumnPair) bool {
	return !columnPair.IsAutoIncrement && !columnPair.IsPrimaryKey
}

func DeterminePK(tableMeta *common.TableMeta) {
	// determine the primary key
	// when pk column is not integer, it will appear in the position[0]
	// or if the columnNames contains "id" , it should be named "id"
	tableMeta.PkColumnName = tableMeta.Columns[0].ColumnName
	for _, column := range tableMeta.Columns {
		if column.ColumnName == "id" {
			tableMeta.PkColumnName = "id"
			break
		}
	}
	for _, columnPair := range tableMeta.Columns {
		if IsID(columnPair) {
			tableMeta.PkColumnName = columnPair.ColumnName
			break
		}
	}
	if common.LogEnabled {
		fmt.Printf(" PK Column for Table %s is %s \n", tableMeta.TableName, tableMeta.PkColumnName)
	}
}
