package repository

import (
	"database/sql"
	"time"

	"github.com/marcovargas74/m74-fair-api/src/config"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

const (
	//DROP_DB CLEAR DATABASE
	DROP_DB = true

	//NOT_DROP_DB CLEAR DATABASE
	NOT_DROP_DB = false

	INDEX_NAME         = 11
	INDEX_DISTRICT     = 6
	INDEX_REGION5      = 9
	INDEX_NEIGHBORHOOD = 15
)

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		logs.Error("FAILED to execute \"command\" [%s] in DB Err[%v]", sql, err)
	}
	return result
}

//CreateDB Create Database and Tables if they don't exist yet
func CreateDB(isDropTable bool) error {

	logs.Info(" %s Starting the connection to MySQL URL:[%s], WAIT...", logs.ThisFunction(), config.DataBaseURL())

	db, err := sql.Open("mysql", config.DataBaseURL())
	if err != nil {
		logs.Error("%s FAILURE to CONNECT to MySQL Err[%v]", logs.ThisFunction(), err.Error())
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	time.Sleep(3 * time.Second)

	logs.Info(" %s Connected To BASE [%s] do MySQL", logs.ThisFunction(), config.DataBaseURL())
	exec(db, "create database if not exists fairAPI")
	exec(db, "use fairAPI")
	if isDropTable {
		logs.Warn("%s WARNING: ALL DATA in TABLE [%s] will be deleted from MySQL", logs.ThisFunction(), config.DataBaseName())
		exec(db, "drop table if exists fair")
	}

	exec(db, TAB_FAIR_CREATE)
	return err
}

//OpenMysql create new repository
func OpenMysql() *sql.DB {

	dataSourceName, err := config.ConfigGetMysqlURL()
	if err != nil {
		logs.Warn("%s Using DEFAULT DATABSE err:[%v] ", logs.ThisFunction(), err.Error())
	}

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logs.Error("%s FAILURE to OPEN to MySQL Err[%v]", logs.ThisFunction(), err.Error())
		return nil
	}

	logs.Info(" %s Connected To BASE [%s] do MySQL", logs.ThisFunction(), config.DataBaseURL())
	return db
}
