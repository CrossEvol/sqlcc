package sql_builder

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"strings"
	"text/template"
)

// TODO: need modify or not ? how do ? concat it later ?
const mapperTemplate = `
package table_meta

import "fmt"

type {{ .TableName | ToLower }}_ struct {
	TABLE      string
	ALIAS      string
{{- range $index, $column := .Columns }}
    {{ $column.ColumnName | ToScreamingSnake }} string
{{- end }}
}

func New{{ .TableName | ToCamel }}_() {{ .TableName | ToLower }}_ {
	return {{ .TableName | ToLower }}_{
		TABLE:      "{{ .TableName | ToCamel }}",
{{- range $index, $column := .Columns }}
        {{ $column.ColumnName | ToScreamingSnake }}: "{{ $column.ColumnName }}", 
{{- end }}
	}
}

func NewAlias{{ .TableName | ToCamel }}_(alias string) {{ .TableName | ToLower }}_ {
	withAlias := func(field string) string {
		return fmt.Sprintf("%s.%s", alias, field)
	}

	return {{ .TableName | ToLower }}_{
		TABLE:      "{{ .TableName | ToCamel }}",
		ALIAS:      alias,
{{- range $index, $column := .Columns }}
        {{ $column.ColumnName | ToScreamingSnake }}: withAlias("{{ $column.ColumnName }}"), 
{{- end }}
	}
}
`

func newTmpl() (*template.Template, error) {
	tmpl := template.Must(template.New("mapperTemplate").Funcs(template.FuncMap{
		"ToSnake":          strcase.ToSnake,
		"ToCamel":          strcase.ToCamel,
		"ToLower":          strcase.ToLowerCamel,
		"ToScreamingSnake": strcase.ToScreamingSnake,
		"Plural":           inflection.Plural,
		"Singular":         inflection.Singular,
		"Quote":            Quote,
		"Add":              func(a, b int) int { return a + b },
		"Last":             LastFunc,
	}).Parse(mapperTemplate))

	return tmpl, nil
}

func Quote(string string) string {
	return fmt.Sprintf("`%s`", string)
}

func LastFunc(index, length int) bool {
	return index == length-1
}

func GenMapperCode(tableMeta common.TableMeta) ([]byte, error) {
	tmpl, err := newTmpl()
	if err != nil {
		return nil, err
	}

	var content strings.Builder
	err = tmpl.Execute(&content, tableMeta)
	if err != nil {
		return nil, err
	}

	return []byte(content.String()), nil
}
