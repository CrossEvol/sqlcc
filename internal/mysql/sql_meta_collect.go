package mysql

import (
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/util"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
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
	args = append(args, common.DbName)
	if len(targetTables) == 0 {
		stmt = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?"
	} else {
		stmt = fmt.Sprintf("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME IN (%s)", strings.Join(targetTables, " , "))
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

		tableMeta := common.TableMeta{TableName: table, Columns: columns}

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

func getColumns(db *sql.DB, tableName string) ([]*common.ColumnPair, error) {
	// Use information_schema to get column information
	rows, err := db.Query("SELECT DISTINCT COLUMN_NAME,COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_NAME = ?  AND TABLE_SCHEMA = ? ", tableName, common.DbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*common.ColumnPair
	for rows.Next() {
		//var colName string
		var columnPair common.ColumnPair
		if err := rows.Scan(&columnPair.ColumnName, &columnPair.ColumnType); err != nil {
			return nil, err
		}
		columns = append(columns, &columnPair)
	}

	// find primary key for mysql
	pkColumnName, err := findPkName(db, tableName)
	if err != nil {
		return columns, err
	}

	for _, column := range columns {
		if column.ColumnName == pkColumnName {
			column.IsAutoIncrement = true
			column.IsPrimaryKey = true
			break
		}
	}

	return columns, nil
}

func findPkName(db *sql.DB, tableName string) (string, error) {
	var pkColumnName string
	rows, err := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = ?  AND TABLE_NAME = ?  AND CONSTRAINT_NAME = 'PRIMARY'", common.DbName, tableName)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&pkColumnName); err != nil {
			return "", err
		}
	}
	return pkColumnName, nil
}
