package sql_builder

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenMapperCode(t *testing.T) {
	tableMeta := common.TableMeta{TableName: "Todo", Columns: []common.ColumnPair{
		{ColumnName: "id"},
		{ColumnName: "title"},
		{ColumnName: "created_at"},
	}}
	bytes, err := GenMapperCode(tableMeta)
	require.Nil(t, err)
	fmt.Println(string(bytes))
}
