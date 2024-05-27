package postgres

import (
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
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
	for _, table := range tableNames {
		if common.LogEnabled {
			fmt.Printf("Table: %s\n", table)
		}
		columns, err := getColumns(db, table)
		tableMetas = append(tableMetas, common.TableMeta{TableName: table, Columns: columns})
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
	return columns, nil
}
