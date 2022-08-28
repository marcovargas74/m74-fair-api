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
)

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		logs.Error("FALHA ao executar \"o comando\" [%s] no DB Err[%v]", sql, err)
	}
	return result
}

//CreateDB Cria Base de dados e Tabelas caso nao existam ainda
func CreateDB(isDropTable bool) {
	//logs.Info(" CreateDB() Conectado  a URL [%s] do MySQL", config.DataBaseURL())
	db, err := sql.Open("mysql", config.DataBaseURL())
	if err != nil {
		logs.Error("CreateDB() FALHA ao conectar ao Banco Mysql %v", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	time.Sleep(3 * time.Second)
	//--------------------------------------------------------------
	logs.Info(" CreateDB() Conectado  a BASE [%s] do MySQL", config.DataBaseName())

	exec(db, "create database if not exists fairAPI")
	exec(db, "use fairAPI")
	if isDropTable {
		logs.Warn("WARNING: TODOS os DADOS da TABELA[%s] serão excluidos do Banco Mysql", config.DataBaseName())
		exec(db, "drop table if exists fair")
	}

	exec(db, TAB_FAIR_CREATE)
}

//OpenMysql create new repository
func OpenMysql() *sql.DB {

	dataSourceName, err := config.ConfigGetMysqlURL()
	if err != nil {
		logs.Warn("OpenMysql()Usando Banco DEFAULT err:[%v] ", err.Error())
	}

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logs.Warn("OpenMysql() Falha ao Tentar abrir a BASE de dados :[%v] ", err.Error())
		return nil
	}

	logs.Info("OpenMysql() Conectado  a BASE [%s] do MySQL", config.DataBaseName())
	return db
}
