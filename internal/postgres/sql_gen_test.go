package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/models"
	"testing"
)

func TestDbConn(t *testing.T) {
	db, err := sql.Open("pgx", "postgresql://postgres:pI9TpV5PFw88KfF3UU@localhost:5432/t3-blog")
	defer db.Close()
	if err != nil {
		fmt.Println("Failed to open the db...")
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the db...")
		fmt.Println(err)
	}

}

func TestTablesMetadata(t *testing.T) {
	db, err := sql.Open("pgx", "postgresql://postgres:pI9TpV5PFw88KfF3UU@localhost:5432/t3-blog")
	defer db.Close()
	if err != nil {
		fmt.Println("Failed to open the db...")
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the db...")
		fmt.Println(err)
	}
	rows, err := db.QueryContext(context.Background(), "SELECT TABLE_NAME ,table_schema FROM information_schema.TABLES  WHERE table_schema = 'public'")
	for rows.Next() {
		var tableName, tableSchema string
		err := rows.Scan(&tableName, &tableSchema)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Sprintf("%s : %s", tableName, tableSchema))
	}

	rows, err = db.QueryContext(context.Background(), `SELECT DISTINCT "column_name",data_type FROM information_schema.COLUMNS WHERE table_schema = 'public' AND "table_name" = 'User'`)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var columnName, dataType string
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Sprintf("%s : %s", columnName, dataType))
	}
}

func TestGen(t *testing.T) {
	config := models.Config{
		DbDriver: "pgx",
		DbDsn:    "postgresql://postgres:pI9TpV5PFw88KfF3UU@localhost:5432/t3-blog",
		DbName:   "t3-blog",
		Tables:   []string{"Account", "Comment", "Post", "Session", "User", "VerificationToken"},
	}
	service := NewPostgresGenService()
	service.Gen(config)
}
