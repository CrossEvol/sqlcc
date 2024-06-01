package util

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"strings"
	"text/template"
)

func Quote(string string) string {
	return fmt.Sprintf("`%s`", string)
}

func Quote2(string string) string {
	return fmt.Sprintf("'%s'", string)
}

func LastFunc(index, length int) bool {
	return index == length-1
}

var CreateList = []string{"Create", "create", "CREATE"}

func IsNotCreate(columnPair common.ColumnPair) bool {
	for _, s := range CreateList {
		if strings.Contains(columnPair.ColumnName, s) {
			return false
		}
	}
	return true
}

func TemplateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"ToSnake":          strcase.ToSnake,
		"ToCamel":          strcase.ToCamel,
		"ToLower":          strcase.ToLowerCamel,
		"ToScreamingSnake": strcase.ToScreamingSnake,
		"Plural":           inflection.Plural,
		"Singular":         inflection.Singular,
		"Quote":            Quote,
		"Add":              func(a, b int) int { return a + b },
		"Last":             LastFunc,
		"IsID":             IsID,
		"IsNotID":          IsNotID,
		"IsNotCreate":      IsNotCreate,
	}
	return funcMap
}
