package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/common"
	"github.com/crossevol/sqlcc/internal/models"
	"github.com/crossevol/sqlcc/internal/util"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Import the Sqlite driver
)

func init() {

}

func dbConn(_dbDriver string, _dbDsn string, _logEnabled bool) {
	common.DbDriver = _dbDriver
	common.DbDsn = _dbDsn
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
		stmt = "SELECT tbl_name FROM sqlite_master  WHERE type='table'"
	} else {
		stmt = fmt.Sprintf("SELECT tbl_name FROM sqlite_master  WHERE type='table' AND tbl_name in (%s);", strings.Join(targetTables, " , "))
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
	rows, err := db.Query(fmt.Sprintf(`PRAGMA table_info('%s');`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*common.ColumnPair
	for rows.Next() {
		var pragma models.Pragma
		if err := rows.Scan(&pragma.CID, &pragma.Name, &pragma.Type, &pragma.Notnull, &pragma.DfltValue, &pragma.PK); err != nil {
			fmt.Println(err)
		}

		columnPair := common.ColumnPair{ColumnName: pragma.Name, ColumnType: pragma.Type}
		if pragma.PK == 1 {
			columnPair.IsAutoIncrement = true
			columnPair.IsPrimaryKey = true
		}
		columns = append(columns, &columnPair)
	}
	return columns, nil
}
