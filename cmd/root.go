/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/crossevol/sqlcc/internal/models"
	mysqlgen "github.com/crossevol/sqlcc/internal/mysql"
	postgresgen "github.com/crossevol/sqlcc/internal/postgres"
	"github.com/crossevol/sqlcc/internal/service"
	sqlitegen "github.com/crossevol/sqlcc/internal/sqlite"
	_ "github.com/go-sql-driver/mysql" // Import the Mysql driver
	_ "github.com/jackc/pgx/v5/stdlib" // Import the Postgres driver
	_ "github.com/mattn/go-sqlite3"    // Import the Sqlite driver
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	tomlFile = "sqlcc.toml"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sqlcc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "This command can generate the sqlcc.toml and then validate it, or it can generate the sqlc.json for sqlc generation.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate the sqlcc.toml config file.",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		blob := `
# Choose the dialect among ['pgx', 'mysql', 'sqlite3']
DB_DRIVER = "<?>"

# mysql =>  user:password@tcp(host:port)/db?k=v
# pgx   =>  user:password@host:port/db?k=v
DB_DSN = "user:password@tcp(host:port)/test_db?charset=utf8"

# target database name you choose , not needed for sqlite and pgx, only for mysql
DB_NAME = "test_db"

# target tables used for stmt generation
TABLES = ["XXXX", "YYYY"]

# Enable logger or not
LOG_ENABLED = false
	`
		if err := os.WriteFile("sqlcc.toml", []byte(blob), fs.ModePerm); err != nil {
			fmt.Println("Failed to initialize the sqlcc.toml...")
			os.Exit(1)
		}
		fmt.Println("Generate sqlcc.toml success.")
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the sqlcc.toml , whether it can connect to the database or not.",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		if len(args) > 1 {
			return errors.New("Invalid Params")
		}
		if len(args) == 1 {
			if args[0] == "." {
				return nil
			}
			if _, err := os.Stat(args[0]); err != nil {
				return errors.New("Config File Not Found")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var config models.Config
		if len(args) == 0 || args[0] == "." {
			toml.DecodeFile(tomlFile, &config)
		}
		db, err := sql.Open(config.DbDriver, config.DbDsn)
		if err != nil {
			fmt.Println("Failed to open the DB_DSN...")
		}
		err = db.Ping()
		if err != nil {
			fmt.Println("Failed to connect to the DB_DSN...")
		}
		fmt.Println("Succeed to verify configuration.")
		db.Close()
	},
}

var configSqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Generate the sqlc.json for sqlite dialect.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		const sqlite3SqlcJson = `
{
  "version": "1",
  "cloud": {
    "project": "01HAQMMECEYQYKFJN8MP16QC41"
  },
  "packages": [
    {
      "name": "dao",
      "schema": "schemas",
      "queries": "queries",
      "path": "internal/database/dao",
      "engine": "sqlite",
      "database": {
        "uri": "file:ondeck?mode=memory&cache=shared"
      },
      "rules": [
        "sqlc/db-prepare"
      ],
      "emit_json_tags": true,
      "emit_db_tags": true,
      "emit_exported_queries": true,
      "emit_prepared_queries": true,
      "emit_interface": true,
      "overrides": [
        {
          "db_type": "int",
          "go_type": {
            "type": "int"
          }
        }
      ]
    }
  ]
}
`
		if _, err := os.Stat("sqlc.json"); err != nil {
			os.WriteFile("sqlc.json", []byte(sqlite3SqlcJson), fs.ModePerm)
			fmt.Println("sqlc.json created.")
		} else {
			os.WriteFile(fmt.Sprintf("sqlc.%s.json", strings.ReplaceAll(strconv.Itoa(int(time.Time{}.Unix())), "-", "")), []byte(sqlite3SqlcJson), fs.ModePerm)
			fmt.Println("sqlc.json has already been created.")
		}

	},
}

var configMysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Generate the sqlc.json for mysql dialect.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		const mysqlSqlcJson = `
{
  "version": "1",
  "cloud": {
    "project": "01HAQMMECEYQYKFJN8MP16QC41"
  },
  "packages": [
    {
      "name": "dao",
      "schema": "schemas",
      "queries": "queries",
      "path": "internal/database/dao",
      "engine": "mysql",
      "database": {
        "uri": "db_dsn?parseTime=true"
      },
      "rules": [
        "sqlc/db-prepare"
      ],
      "emit_json_tags": true,
      "emit_db_tags": true,
      "emit_exported_queries": true,
      "emit_prepared_queries": true,
      "emit_interface": true,
      "overrides": [
        {
          "db_type": "int",
          "go_type": {
            "type": "int"
          }
        }
      ]
    }
  ]
}
`
		if _, err := os.Stat("sqlc.json"); err != nil {
			os.WriteFile("sqlc.json", []byte(mysqlSqlcJson), fs.ModePerm)
			fmt.Println("sqlc.json created.")
		} else {
			os.WriteFile(fmt.Sprintf("sqlc.%s.json", strings.ReplaceAll(strconv.Itoa(int(time.Time{}.Unix())), "-", "")), []byte(mysqlSqlcJson), fs.ModePerm)
			fmt.Println("sqlc.json has already been created.")
		}

	},
}

var configPgxCmd = &cobra.Command{
	Use:   "pgx",
	Short: "Generate the sqlc.json for postgres dialect.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		const pgxSqlcJson = `
{
  "version": "1",
  "cloud": {
    "project": "01HAQMMECEYQYKFJN8MP16QC41"
  },
  "packages": [
    {
      "name": "dao",
	  "path": "internal/database/dao",
      "schema": "postgresql/schema",
      "queries": "postgresql/query",
      "engine": "postgresql",
      "sql_package": "database/sql",
      "database": {
        "uri": "${VET_TEST_EXAMPLES_POSTGRES_ONDECK}"
      },
      "analyzer": {
        "database": false
      },
      "rules": [
        "sqlc/db-prepare"
      ],
      "emit_json_tags": true,
      "emit_prepared_queries": true,
      "emit_interface": true,
      "emit_prepared_queries": true,
      "emit_interface": true,
      "overrides": [
        {
          "db_type": "pg_catalog.int4",
          "go_type": {
            "type": "int"
          }
        }
      ]
    }
  ]
}
`

		if _, err := os.Stat("sqlc.json"); err != nil {
			os.WriteFile("sqlc.json", []byte(pgxSqlcJson), fs.ModePerm)
			fmt.Println("sqlc.json created.")
		} else {
			os.WriteFile(fmt.Sprintf("sqlc.%s.json", strings.ReplaceAll(strconv.Itoa(int(time.Time{}.Unix())), "-", "")), []byte(pgxSqlcJson), fs.ModePerm)
			fmt.Println("sqlc.json has already been created.")
		}

	},
}

var genCrudSqlCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate the basic curd sql statement for sqlc",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var config models.Config
		var genService service.SqlGenService
		toml.DecodeFile(tomlFile, &config)
		if config.LogEnabled {
			marshalIndent, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				fmt.Println("Failed to marshal config...")
			}
			fmt.Println(string(marshalIndent))
		}
		genService = newGenService(config.DbDriver)
		genService.Gen(config)
	},
}

var genSqlBuilderMapperCmd = &cobra.Command{
	Use:   "map",
	Short: "map the table metadata to an object.",
	Long:  `when using the sql builder for golang, use this command to map the column name to static constants by organize them into an object using the column names as values, which can improve the type safety.`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var config models.Config
		toml.DecodeFile(tomlFile, &config)
		if config.LogEnabled {
			marshalIndent, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				fmt.Println("Failed to marshal config...")
			}
			fmt.Println(string(marshalIndent))
		}
		genService := newGenService(config.DbDriver)
		genService.GenMapper(config)
	},
}

func newGenService(driver string) service.SqlGenService {
	var genService service.SqlGenService
	switch strings.ToLower(driver) {
	case "mysql":
		genService = mysqlgen.NewMysqlGenService()
	case "pgx":
		genService = postgresgen.NewPostgresGenService()
	case "sqlite3":
		genService = sqlitegen.NewSqliteGenService()
	}

	return genService
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(genCrudSqlCmd)
	rootCmd.AddCommand(genSqlBuilderMapperCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configMysqlCmd)
	configCmd.AddCommand(configPgxCmd)
	configCmd.AddCommand(configSqliteCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sqlcc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
