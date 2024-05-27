package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/models"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Import the Sqlite driver
)

var (
	dbDriver   string
	dbDsn      string
	logEnabled bool
)

func init() {

}

type ColumnPair struct {
	ColumnName string
	ColumnType string
}

type TableMeta struct {
	TableName string
	Columns   []ColumnPair
}

func dbConn(_dbDriver string, _dbDsn string, _logEnabled bool) {
	dbDriver = _dbDriver
	dbDsn = _dbDsn
	logEnabled = _logEnabled
}

func getTableMetas(selectedTables ...string) []TableMeta {
	var targetTables []string
	for _, selectedTable := range selectedTables {
		targetTables = append(targetTables, fmt.Sprintf("'%s'", selectedTable))
	}
	db, err := sql.Open(dbDriver, dbDsn)
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
	var tableMetas []TableMeta
	for _, table := range tableNames {
		if logEnabled {
			fmt.Printf("Table: %s\n", table)
		}
		columns, err := getColumns(db, table)
		tableMetas = append(tableMetas, TableMeta{TableName: table, Columns: columns})
		if err != nil {
			fmt.Println("Error getting columns:", err)
			continue
		}
		if logEnabled {
			for _, column := range columns {
				fmt.Printf("  - Column Name: %s , Column Type: %s\n", column.ColumnName, column.ColumnType)
			}
		}
	}

	return tableMetas
}

func getColumns(db *sql.DB, tableName string) ([]ColumnPair, error) {
	// Use information_schema to get column information
	rows, err := db.Query(fmt.Sprintf(`PRAGMA table_info('%s');`, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnPair
	for rows.Next() {
		var pragma models.Pragma
		if err := rows.Scan(&pragma.CID, &pragma.Name, &pragma.Type, &pragma.Notnull, &pragma.DfltValue, &pragma.PK); err != nil {
			fmt.Println(err)
		}

		columns = append(columns, ColumnPair{ColumnName: pragma.Name, ColumnType: pragma.Type})
	}
	return columns, nil
}
