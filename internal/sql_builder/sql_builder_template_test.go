package sql_builder

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/mysql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenMapperCode(t *testing.T) {
	tableMeta := mysql.TableMeta{TableName: "Todo", Columns: []mysql.ColumnPair{
		{ColumnName: "id"},
		{ColumnName: "title"},
		{ColumnName: "created_at"},
	}}
	bytes, err := GenMapperCode(tableMeta)
	require.Nil(t, err)
	fmt.Println(string(bytes))
}
