package persistence

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"github.com/kelseyhightower/envconfig"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBParameters struct {
	// One of "sqlite", "mysql", "postgres"
	Type string `required:"true"`

	// Text template for use in building the complete DSN
	Dsn      string `required:"true"`
	Host     string `required:"true"`
	Port     int    `required:"true"`
	Database string `required:"true"`
	User     string `required:"true"`
	Pass     string `required:"true"`
}

// GetParametersFromEnv retrieves the DBParameters from environment variables.
//
// Returns *DBParameters and error.
func GetParametersFromEnv() (*DBParameters, error) {
	var params DBParameters
	err := envconfig.Process("db", &params)

	return &params, err
}

// OpenDialectorFromEnv returns a gorm.Dialector based on the values of the "DB_***" environment variables.
//
// The "DB_TYPE" environment variable specifies the type of database to connect to.
// The "DB_DSN" environment variable specifies the data source name (DSN) for the database connection.
//
// It panics if either "DB_TYPE" or "DB_DSN" environment variables are not configured.
// It returns a gorm.Dialector based on the provided "DB_TYPE" and "DB_DSN".
func OpenDialectorFromEnv() (gorm.Dialector, error) {
	var params DBParameters

	err := envconfig.Process("db", &params)
	if err != nil {
		return nil, err
	}

	return OpenDialector(&params)
}

// OpenDialector creates a GORM dialector based on the given database type and DSN template.
//
// Parameters:
// - dbType: The type of the database (e.g., "SQLITE", "POSTGRES", "MYSQL").
// - dsnTemplate: The DSN template used to create the connection string.
//
// Returns:
// - gorm.Dialector: The GORM dialector based on the given parameters.
func OpenDialector(params *DBParameters) (gorm.Dialector, error) {

	dsn, err := makeDSN(params)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(params.Type) {
	case "SQLITE":
		return sqlite.Open(dsn), nil
	case "POSTGRES":
		return postgres.Open(dsn), nil
	case "MYSQL":
		return mysql.Open(dsn), nil
	}

	return nil, errors.New("unknown database type")
}

// makeDSN generates a Data Source Name (DSN) using a template string.
//
// It takes in a dsnTemplate string and returns a string representing the generated DSN.
// The dsnTemplate is parsed using the "dsn" template of the Go template package.
// The parsed template is then executed using a DBParameters struct to populate the template variables
// with the values from the environment variables.
// The function assumes that the necessary environment variables for the DBParameters are set.
func makeDSN(params *DBParameters) (string, error) {
	tmpl := template.New("dsn")

	tmpl, err := tmpl.Parse(params.Dsn)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, params)
	if err != nil {
		panic(err)
	}

	return doc.String(), nil
}
