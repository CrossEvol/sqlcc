package mysql

import (
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/util"
	"strings"
	"text/template"
)

func MySQLTemplateCrudSql() (*template.Template, error) {
	const mysqlCrudTemplate = `
-- name: Get{{ .TableName | ToCamel | Singular }} :one
SELECT * FROM {{ .TableName | Quote }}
WHERE {{ .PkColumnName }} = ? LIMIT 1;

-- name: Get{{ .TableName | ToCamel | Plural }} :many
SELECT * FROM {{ .TableName | Quote }};

-- name: Get{{ .TableName | ToCamel | Plural }}By{{ .PkColumnName | ToCamel }}s :many
SELECT * FROM {{ .TableName | Quote }} WHERE {{ .PkColumnName }} IN (sqlc.slice('{{ .PkColumnName }}s'));

-- name: Count{{ .TableName | ToCamel | Plural }} :one
SELECT count(*) FROM {{ .TableName | Quote }};

-- name: Create{{  .TableName | ToCamel | Singular }} :execresult
INSERT INTO {{ .TableName | Quote }} (
{{ range $index, $column := .Columns }}
  {{- if IsNotID $column }}{{ $column.ColumnName  | Quote }}{{ if not (Last $index (len $.Columns)) }},{{ end }}{{ end }}
{{- end }}
) VALUES (
{{ range $index, $column := .Columns }}
  {{- if IsNotID $column }}? {{ if not (Last $index (len $.Columns)) }},{{ end }}{{ end }}
{{- end }}
);

-- name: Update{{  .TableName | ToCamel | Singular }} :execresult
UPDATE {{ .TableName | Quote }}
SET {{ range $index, $column := .Columns }}
  {{ if and  (IsNotID $column) (IsNotCreate $column) }}{{ $column.ColumnName  | Quote }} = CASE WHEN sqlc.arg('{{ $column.ColumnName }}') IS NOT NULL THEN sqlc.arg('{{ $column.ColumnName }}') ELSE {{ $column.ColumnName  | Quote }} END{{ if not (Last $index (len $.Columns)) }},{{ end }}{{ end }}
{{- end }}
WHERE {{ .PkColumnName }} = ?;

-- name: Delete{{  .TableName | ToCamel | Singular }} :exec
DELETE FROM {{ .TableName | Quote }}
WHERE {{ .PkColumnName }} = ?;
`

	tmpl := template.Must(template.New("mysqlCrudTemplate").Funcs(util.TemplateFuncMap()).Parse(mysqlCrudTemplate))

	return tmpl, nil
}

func ContentTemplateCrudSql(tableMeta common.TableMeta) ([]byte, error) {
	tmpl, err := MySQLTemplateCrudSql()
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
