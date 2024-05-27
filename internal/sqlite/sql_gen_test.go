package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func TestSqliteConn(t *testing.T) {
	db, err := sql.Open("sqlite3", "dev.db")
	if err != nil {
		log.Fatal("Failed to open the db ...")
		fmt.Println(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to Ping the db ...")
	}
}

func TestSqliteTableMetadata(t *testing.T) {
	db, err := sql.Open("sqlite3", "D:\\GOLANG_SAMPLES\\navidrome_0.50.2_windows_amd64\\navidrome.db")
	if err != nil {
		log.Fatal("Failed to open the db ...")
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT tbl_name FROM sqlite_master  WHERE type='table';`)
	if err != nil {
		fmt.Println(err)
	}
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			fmt.Println(err)
		} else {
			tables = append(tables, table)
		}
	}
	for _, table := range tables {
		fmt.Printf("[%s]\n", table)
		rows, err := db.Query(fmt.Sprintf(`PRAGMA table_info('%s');`, table))
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			var pragma models.Pragma
			if err := rows.Scan(&pragma.CID, &pragma.Name, &pragma.Type, &pragma.Notnull, &pragma.DfltValue, &pragma.PK); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s : %s \n", pragma.Name, pragma.Type)
		}
	}
}

func TestName(t *testing.T) {
	config := models.Config{
		DbDriver: "sqlite3",
		DbDsn:    "D:\\GOLANG_SAMPLES\\navidrome_0.50.2_windows_amd64\\navidrome.db",
		Tables:   []string{"album", "artist"},
	}
	service := NewSqliteGenService()
	service.Gen(config)
}
