package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

//AddrOpenDB VAR used to open and to access BD
var AddrOpenDB string

//AddrDB VAR data source name
var AddrDB string

// //FairMySQL mysql repo
// type FairMySQL struct {
// 	db *sql.DB
// }

//TODO trazer a criacao do banco pra ca
const (
	//DBSourceOpenLocal Const used to Open Local db
	DBSourceOpenLocal = "root:my-secret-pw@tcp(localhost:3307)/"

	//DBSourceLocal Const used to acces Local db
	DBSourceLocal = "root:my-secret-pw@tcp(localhost:3307)/fairAPI?parseTime=true"

	//DBSourceOpenDocker Const used to Open Docker db
	//DBSourceOpenDocker = "root:my-secret-pw@tcp(mysql-api)/" //mysql-api é o nome do serviço no docker-composer

	//DBSourceDocker Const used to acces Docker db
	//DBSourceDocker = "root:my-secret-pw@tcp(mysql-api)/fairAPI"

)

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}

func CreateDB() {
	//--------------SE CONECTA AO BANCO
	//dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	//db, err := sql.Open("mysql", dataSourceName)

	//AddrOpenDB = DBSourceOpenDocker
	//AddrDB = DBSourceDocker
	AddrOpenDB = DBSourceOpenLocal
	AddrDB = DBSourceLocal

	db, err := sql.Open("mysql", AddrOpenDB)
	if err != nil {
		logs.Error("FALHA ao conectar ao Banco Mysql do DOcker %v", err)
		/*AddrOpenDB = DBSourceOpenLocal
		AddrDB = DBSourceLocal
		db, err = sql.Open("mysql", AddrOpenDB)
		if err != nil {
			logs.Error("FALHA ao conectar ao Banco Mysql Local IP 127.0.0.1 %v", err)
		}*/
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	time.Sleep(5 * time.Second)
	//--------------------------------------------------------------
	fmt.Println("Conectado ao Banco")
	exec(db, "create database if not exists fairAPI")
	exec(db, "use fairAPI")
	//if isDropTable {
	//exec(db, "drop table if exists fair")
	//}

	exec(db, `create table IF NOT EXISTS fair(
		   	idx integer auto_increment,
		   	id varchar(50) ,
		   	name varchar(50),
		   	district varchar(18),
		   	region5 varchar(6),
		   	neighborhood varchar(20),
			created_at datetime,
		   	updated_at datetime,
		   	PRIMARY KEY (idx)
		   	)`)

	//var db *sql.DB
	//repository.CreateDB()
	//db, err := sql.Open("mysql", repository.AddrOpenDB)

	/*db.Close()
	db, err = sql.Open("mysql", repository.AddrDB)
	if err != nil {
		log.Println(err)
		//return nil, err
	}*/

}

/*
func startDB() (*sql.Tx, error) {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err

	}

	return tx, nil
}*/

//OpenMysql create new repository
func OpenMysql() *sql.DB {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

//NewFairMySQL create new repository
/*func NewFairMySQL() *FairMySQL {
	return &FairMySQL{
		db: openMysql(),
	}
}*/
