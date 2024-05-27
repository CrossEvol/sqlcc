package sqlite

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"strings"
	"text/template"
)

func SqliteTemplateCrudSql() (*template.Template, error) {
	const postgresCrudTemplate = `
-- name: Get{{ .TableName | ToCamel | Singular }} :one
SELECT * FROM {{ .TableName | Quote }}
WHERE id = ? LIMIT 1;

-- name: Get{{ .TableName | ToCamel | Plural }} :many
SELECT * FROM {{ .TableName | Quote }};

-- name: Count{{ .TableName | ToCamel | Plural }} :one
SELECT count(*) FROM {{ .TableName | Quote }};

-- name: Create{{  .TableName | ToCamel | Singular }} :execresult
INSERT INTO {{ .TableName | Quote }} (
{{- range $index, $column := .Columns }}
  {{ $column.ColumnName  | Quote }}{{ if not (last $index (len $.Columns)) }},{{ end }}
{{- end }}
) VALUES (
{{- range $index, $column := .Columns }}
  ? {{ if not (last $index (len $.Columns)) }},{{ end }}
{{- end }}
);

-- name: Update{{  .TableName | ToCamel | Singular }} :exec
UPDATE {{ .TableName | Quote }}
SET {{ range $index, $column := .Columns }}
  {{ $column.ColumnName  | Quote }} = ? {{ if not (last $index (len $.Columns)) }},{{ end }}
{{- end }}
WHERE id = ?;

-- name: Delete{{  .TableName | ToCamel | Singular }} :exec
DELETE FROM {{ .TableName | Quote }}
WHERE id = ?;
`

	tmpl := template.Must(template.New("postgresCrudTemplate").Funcs(template.FuncMap{
		"ToSnake":  strcase.ToSnake,
		"ToCamel":  strcase.ToCamel,
		"Plural":   inflection.Plural,
		"Singular": inflection.Singular,
		"Quote":    Quote,
		"Add":      func(a, b int) int { return a + b },
		"last":     LastFunc,
	}).Parse(postgresCrudTemplate))

	return tmpl, nil
}

func ContentTemplateCrudSql(tableMeta common.TableMeta) ([]byte, error) {
	tmpl, err := SqliteTemplateCrudSql()
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

func Quote(string string) string {
	return fmt.Sprintf("`%s`", string)
}

func LastFunc(index, length int) bool {
	return index == length-1
}
