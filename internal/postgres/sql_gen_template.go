package postgres

import (
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/util"
	"strings"
	"text/template"
)

func PostgresTemplateCrudSql() (*template.Template, error) {
	const postgresCrudTemplate = `
-- name: Get{{ .TableName | ToCamel | Singular }} :one
SELECT * FROM {{ .TableName | Quote }}
WHERE id = $1 LIMIT 1;

-- name: Get{{ .TableName | ToCamel | Plural }} :many
SELECT * FROM {{ .TableName | Quote }};

-- name: Count{{ .TableName | ToCamel | Plural }} :one
SELECT count(*) FROM {{ .TableName | Quote }};

-- name: Create{{  .TableName | ToCamel | Singular }} :one
INSERT INTO {{ .TableName | Quote }} (
{{- range $index, $column := .Columns }}
  {{ $column.ColumnName  | Quote }}{{ if not (Last $index (len $.Columns)) }},{{ end }}
{{- end }}
) VALUES (
{{- range $index, $column := .Columns }}
  ${{ Add $index 1 }} {{ if not (Last $index (len $.Columns)) }},{{ end }}
{{- end }}
)
RETURNING *;

-- name: Update{{  .TableName | ToCamel | Singular }} :one
UPDATE {{ .TableName | Quote }}
SET {{ range $index, $column := .Columns }}
    {{ $column.ColumnName  | Quote }} = CASE WHEN @{{ $column.ColumnName }} IS NOT NULL THEN @{{ $column.ColumnName }} ELSE {{ $column.ColumnName  | Quote }} END,   {{ if not (Last $index (len $.Columns)) }},{{ end }}
{{- end }}
WHERE id = ${{ Add (len .Columns) 1 }}
RETURNING *;

-- name: Delete{{  .TableName | ToCamel | Singular }} :exec
DELETE FROM {{ .TableName | Quote }}
WHERE id = $1;
`

	tmpl := template.Must(template.New("postgresCrudTemplate").Funcs(util.TemplateFuncMap()).Parse(postgresCrudTemplate))

	return tmpl, nil
}

func ContentTemplateCrudSql(tableMeta common.TableMeta) ([]byte, error) {
	tmpl, err := PostgresTemplateCrudSql()
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
