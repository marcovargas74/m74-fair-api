package config

import (
	"fmt"
	"os"
)

const (
	VERSION_PACKAGE = "2022-08-21"

	DB_USER             = "root"
	DB_PASSWORD         = "my-secret-pw"
	DB_MYSQL_DATABASE   = "fairAPI"
	DB_HOST             = "127.0.0.1"
	DB_URL              = "root:my-secret-pw@tcp(localhost:3307)"
	DB_PORT             = "3306"
	SERVER_API_PORT_MEM = ":5000"
	SERVER_API_PORT_SQL = ":5001"
	DEV_TEST            = true
)

const (
	VAR_DB_USER           = "DB_USER"
	VAR_DB_PASSWORD       = "my-secret-pw"
	VAR_DB_MYSQL_DATABASE = "fairAPI"
	VAR_DB_HOST           = "127.0.0.1"
	VAR_DB_URL            = "root:my-secret-pw@tcp(localhost:3307)"
	VAR_DB_PORT           = "3306"
	VAR_SERVER_API_PORT   = ":5000"
	VAR_IS_PROD           = "PROD"
	VAR_DEV               = "DEV"
	VAR_DEV_TEST          = true
)

var Environs = map[string]string{
	"_DB_USER":     "root",
	"_DB_PASSWORD": "my-secret-pw",
	"_DB_URL":      "root:my-secret-pw@tcp(localhost:3307)",
	/*
		VAR_DB_USER           = "DB_USER"
		VAR_DB_PASSWORD       = "my-secret-pw"
		VAR_DB_MYSQL_DATABASE = "fairAPI"
		VAR_DB_HOST           = "127.0.0.1"
		VAR_DB_URL            = "root:my-secret-pw@tcp(localhost:3307)"
		VAR_DB_PORT           = "3306"
		VAR_SERVER_API_PORT   = ":5000"
		VAR_PROD              = "PROD"
		VAR_DEV_TEST          = true,*/

}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func ValueOfEnvironVar(varName string) string {
	return Environs[varName]
}

// Main function
func SetEngGet() {

	os.Setenv("_DB_USER", ValueOfEnvironVar("_DB_USER"))
	os.Setenv("_DB_URL", ValueOfEnvironVar("_DB_URL"))

	// returns value of GEEKS
	fmt.Println("_DB_USER:", os.Getenv("_DB_USER"))
	fmt.Println("_DB_URL:", os.Getenv("_DB_URL"))

	// Unset environment variable GEEKS
	//os.Unsetenv("GEEKS")

	// returns empty string and false,
	// because we removed the GEEKS variable
	//value, ok := os.LookupEnv("GEEKS")

	//fmt.Println("GEEKS:", value, " Is present:", ok)

}
