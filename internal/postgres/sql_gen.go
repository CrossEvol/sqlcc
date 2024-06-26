package postgres

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/models"
	"github.com/crossevol/sqlcc/internal/sql_builder"
	"github.com/iancoleman/strcase"
	"io/fs"
	"os"
	"path/filepath"
)

type PostgresGenService struct {
}

func (genService *PostgresGenService) GenMapper(config models.Config) {
	tableMetaDir := common.TableMetaDir
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dbConn(config.DbDriver, config.DbDsn, config.DbName, config.LogEnabled)

	for _, tableMeta := range getTableMetas(config.Tables...) {
		sql, err := sql_builder.GenMapperCode(tableMeta)
		if err != nil {
			fmt.Println(err)
		}
		if _, err := os.Stat(filepath.Join(wd, tableMetaDir)); err != nil {
			os.Mkdir(filepath.Join(wd, tableMetaDir), fs.ModePerm)
		}
		if err := os.WriteFile(filepath.Join(wd, tableMetaDir, strcase.ToSnake(tableMeta.TableName)+".go"), sql, os.ModePerm); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Table Metadata Mapper generated success.")
}

func NewPostgresGenService() *PostgresGenService {
	return &PostgresGenService{}
}

func (genService *PostgresGenService) Gen(config models.Config) {
	gen := common.GenSqlDir
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dbConn(config.DbDriver, config.DbDsn, config.DbName, config.LogEnabled)

	for _, tableMeta := range getTableMetas(config.Tables...) {
		sql, err := ContentTemplateCrudSql(tableMeta)
		if err != nil {
			fmt.Println(err)
		}
		if _, err := os.Stat(filepath.Join(wd, gen)); err != nil {
			os.Mkdir(filepath.Join(wd, gen), fs.ModePerm)
		}
		if err := os.WriteFile(filepath.Join(wd, gen, tableMeta.TableName+".sql"), sql, os.ModePerm); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Basic CRUD stmt generated success.")
}
