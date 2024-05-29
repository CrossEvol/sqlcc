package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/crossevol/sqlcc/internal/models"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestDbConn(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
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

func TestPostgresGenService_Gen(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logEnabled, err := strconv.ParseBool(os.Getenv("LOG_ENABLED"))
	if err != nil {
		log.Fatal(err)
	}
	var tables []string
	for _, table := range strings.Split(os.Getenv("TABLES"), ",") {
		tables = append(tables, table)
	}
	config := models.Config{DbDriver: os.Getenv("DB_DRIVER"), DbDsn: os.Getenv("DB_DSN"), DbName: os.Getenv("DB_NAME"), Tables: tables, LogEnabled: logEnabled}
	service := NewPostgresGenService()
	service.Gen(config)
}
