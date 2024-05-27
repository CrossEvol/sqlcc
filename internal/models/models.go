package models

type Config struct {
	DbDriver   string   `toml:"DB_DRIVER"`
	DbDsn      string   `toml:"DB_DSN"`
	DbName     string   `toml:"DB_NAME"`
	Tables     []string `toml:"TABLES"`
	LogEnabled bool     `toml:"LOG_ENABLED"`
}

type Pragma struct {
	CID       int    `db:"cid" `
	Name      string `db:"name" `
	Type      string `db:"type" `
	Notnull   int8   `db:"notnull" `
	DfltValue any    `db:"dflt_value" `
	PK        int8   `db:"pk" `
}
