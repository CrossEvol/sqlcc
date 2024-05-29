package postgres

import (
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/util"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the Postgres driver
)

func init() {

}

func dbConn(_dbDriver string, _dbDsn string, _dbName string, _logEnabled bool) {
	common.DbDriver = _dbDriver
	common.DbDsn = _dbDsn
	common.DbName = _dbName
	common.LogEnabled = _logEnabled
}

func getTableMetas(selectedTables ...string) []common.TableMeta {
	var targetTables []string
	for _, selectedTable := range selectedTables {
		targetTables = append(targetTables, fmt.Sprintf("'%s'", selectedTable))
	}
	db, err := sql.Open(common.DbDriver, common.DbDsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var stmt string
	var args []any
	if len(targetTables) == 0 {
		stmt = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'public'"
	} else {
		stmt = fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'public' AND TABLE_NAME IN (%s)", strings.Join(targetTables, " , "))
	}

	// Get all tableNames in the database
	rows, err := db.Query(stmt, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		tableNames = append(tableNames, tableName)
	}

	// Loop through each table and get its columns
	var tableMetas []common.TableMeta
	for _, tableName := range tableNames {
		if common.LogEnabled {
			fmt.Printf("Table: %s\n", tableName)
		}
		columns, err := getColumns(db, tableName)
		tableMeta := common.TableMeta{TableName: tableName, Columns: columns}

		util.DeterminePK(&tableMeta)

		tableMetas = append(tableMetas, tableMeta)
		if err != nil {
			fmt.Println("Error getting columns:", err)
			continue
		}
		if common.LogEnabled {
			for _, column := range columns {
				fmt.Printf("  - Column Name: %s , Column Type: %s\n", column.ColumnName, column.ColumnType)
			}
		}
	}

	return tableMetas
}

func getColumns(db *sql.DB, tableName string) ([]common.ColumnPair, error) {
	// Use information_schema to get column information
	rows, err := db.Query(fmt.Sprintf(`SELECT DISTINCT "column_name",data_type FROM information_schema.COLUMNS WHERE table_schema = 'public' AND "table_name" = '%s'`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []common.ColumnPair
	for rows.Next() {
		//var colName string
		var columnPair common.ColumnPair
		if err := rows.Scan(&columnPair.ColumnName, &columnPair.ColumnType); err != nil {
			return nil, err
		}
		columns = append(columns, columnPair)
	}

	// TODO: after uncomment here, will miss the columns data , why ?
	// find primary key for postgres
	//pkMeta, err := findPkColumn(db, tableName)
	//if err != nil {
	//	return nil, err
	//}
	//for _, column := range columns {
	//	if column.ColumnName == pkMeta.ColumnName {
	//		column.IsAutoIncrement = true
	//		column.IsPrimaryKey = true
	//		break
	//	}
	//}
	return columns, nil
}

func findPkColumn(db *sql.DB, tableName string) (*findPkMeta, error) {
	var meta findPkMeta
	rows, err := db.Query(fmt.Sprintf(`SELECT column_name, is_identity, column_default FROM information_schema.columns WHERE table_schema = 'public' AND table_name = %s AND column_default LIKE 'nextval(%%'`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&meta.ColumnName, &meta.IsIdentity, &meta.ColumnDefault); err != nil {
			return nil, err
		}
	}
	return &meta, err
}

type findPkMeta struct {
	ColumnName    string
	IsIdentity    string
	ColumnDefault string
}
