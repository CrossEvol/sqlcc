package misc

import (
	"encoding/json"
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"log"
	"strconv"
	"testing"
)

func TestColumnChangeValueWithoutPointer(t *testing.T) {
	var columns []common.ColumnPair
	for i := 0; i < 5; i++ {
		columns = append(columns, common.ColumnPair{ColumnName: "item" + strconv.Itoa(i)})
	}
	s := toJson(columns)
	fmt.Println(s)

	fmt.Println("-------------------------------------------")

	for _, column := range columns {
		column.ColumnType = "string"
		column.IsAutoIncrement = true
		column.IsPrimaryKey = true
	}
	s = toJson(columns)
	fmt.Println(s)
}

func TestColumnChangeValueWithPointer(t *testing.T) {
	var columns []*common.ColumnPair
	for i := 0; i < 5; i++ {
		columns = append(columns, &common.ColumnPair{ColumnName: "item" + strconv.Itoa(i)})
	}
	s := toJson(columns)
	fmt.Println(s)

	fmt.Println("-------------------------------------------")

	for _, column := range columns {
		column.ColumnType = "string"
		column.IsAutoIncrement = true
		column.IsPrimaryKey = true
	}
	s = toJson(columns)
	fmt.Println(s)
}

func toJson(v any) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}
