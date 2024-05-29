package util

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
)

func DeterminePK(tableMeta *common.TableMeta) {
	// determine the primary key
	// when pk column is not integer, it will appear in the position[0]
	tableMeta.PkColumnName = tableMeta.Columns[0].ColumnName
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
