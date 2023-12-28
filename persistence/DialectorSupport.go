package persistence

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"bytes"
	"os"
	"strings"
	"text/template"
)

type DBParameters struct {
	Host     string
	Port     string
	Database string
	User     string
	Pass     string
}

// OpenDialectorFromEnv returns a gorm.Dialector based on the values of the "DB_***" environment variables.
//
// The "DB_TYPE" environment variable specifies the type of database to connect to.
// The "DB_DSN" environment variable specifies the data source name (DSN) for the database connection.
//
// It panics if either "DB_TYPE" or "DB_DSN" environment variables are not configured.
// It returns a gorm.Dialector based on the provided "DB_TYPE" and "DB_DSN".
func OpenDialectorFromEnv() gorm.Dialector {
	dbType, dbTypeFound := os.LookupEnv("DB_TYPE")
	if !dbTypeFound {
		panic("DB_TYPE environment variable not configured")
	}

	dsnTemplate, dsnTemplateFound := os.LookupEnv("DB_DSN")
	if !dsnTemplateFound {
		panic("DB_DSN environment variable not configured")
	}

	return OpenDialector(dbType, dsnTemplate)
}

// OpenDialector creates a GORM dialector based on the given database type and DSN template.
//
// Parameters:
// - dbType: The type of the database (e.g., "SQLITE", "POSTGRES", "MYSQL").
// - dsnTemplate: The DSN template used to create the connection string.
//
// Returns:
// - gorm.Dialector: The GORM dialector based on the given parameters.
func OpenDialector(dbType string, dsnTemplate string) gorm.Dialector {

	switch strings.ToUpper(dbType) {
	case "SQLITE":
		return sqlite.Open(makeDSN(dsnTemplate))
	case "POSTGRES":
		return postgres.Open(makeDSN(dsnTemplate))
	case "MYSQL":
		return mysql.Open(makeDSN(dsnTemplate))
	}

	return nil
}

// makeDSN generates a Data Source Name (DSN) using a template string.
//
// It takes in a dsnTemplate string and returns a string representing the generated DSN.
// The dsnTemplate is parsed using the "dsn" template of the Go template package.
// The parsed template is then executed using a DBParameters struct to populate the template variables
// with the values from the environment variables.
// The function assumes that the necessary environment variables for the DBParameters are set.
func makeDSN(dsnTemplate string) string {
	tmpl := template.New("dsn")
	tmpl, err := tmpl.Parse(dsnTemplate)
	if err != nil {
		panic(err)
	}

	params := DBParameters{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		User:     os.Getenv("DB_USER"),
		Pass:     os.Getenv("DB_PASS"),
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, params)
	if err != nil {
		panic(err)
	}

	return doc.String()
}
