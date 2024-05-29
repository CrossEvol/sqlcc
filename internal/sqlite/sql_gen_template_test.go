package sqlite

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"log"
	"testing"
)

func TestContentTemplateCrudSql(t *testing.T) {
	tableMeta := common.TableMeta{TableName: "post", PkColumnName: "id", Columns: []common.ColumnPair{
		{ColumnName: "title"},
	}}
	sql, err := ContentTemplateCrudSql(tableMeta)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(sql))
}
