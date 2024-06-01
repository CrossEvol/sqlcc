package sqlite

import (
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/util"
	"strings"
	"text/template"
)

func SqliteTemplateCrudSql() (*template.Template, error) {
	const postgresCrudTemplate = `
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
  {{- if IsNotID $column }}{{ $column.ColumnName  | Quote }}{{ if not (Last $index (len $.Columns)) }},{{ end }}{{end}}
{{- end }}
) VALUES (
{{ range $index, $column := .Columns }}
  {{- if IsNotID $column }}? {{ if not (Last $index (len $.Columns)) }},{{ end }}{{ end }}
{{- end }}
);

-- name: Update{{  .TableName | ToCamel | Singular }} :execresult
UPDATE {{ .TableName | Quote }}
SET {{ range $index, $column := .Columns }}
  {{ if IsNotID $column }}{{ $column.ColumnName  | Quote }} = CASE WHEN @{{ $column.ColumnName }} IS NOT NULL THEN @{{ $column.ColumnName }} ELSE {{ $column.ColumnName  | Quote }} END{{ if not (Last $index (len $.Columns)) }},{{ end }}{{end}}
{{- end }}
WHERE {{ .PkColumnName }} = ?;

-- name: Delete{{  .TableName | ToCamel | Singular }} :exec
DELETE FROM {{ .TableName | Quote }}
WHERE {{ .PkColumnName }} = ?;
`

	tmpl := template.Must(template.New("postgresCrudTemplate").Funcs(util.TemplateFuncMap()).Parse(postgresCrudTemplate))

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
