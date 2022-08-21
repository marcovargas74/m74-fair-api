package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/repository"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/usecase/fair"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/api/handler"
	logs "github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/logs"
	"github.com/urfave/negroni"
)

func init() {
	//myquery.CreateStatus()
	//myquery.SetUsingMongoDocker(myquery.SetDockerRun)
	//myquery.CreateDB()
	prod := flag.Bool("prod", false, "write log to file")
	logs.Start(*prod, "./fairAPI.log")

}

const (
	//DBSourceOpenLocal Const used to Open Local db
	DBSourceOpenLocal = "root:my-secret-pw@tcp(127.0.0.1:3307)/"

	//DBSourceLocal Const used to acces Local db
	DBSourceLocal = "root:my-secret-pw@tcp(127.0.0.1:3307)/fairAPI" //root:Mysql#my-secret-pw@/validatorAPP"

	//DBSourceOpenDocker Const used to Open Docker db
	DBSourceOpenDocker = "root:my-secret-pw@tcp(mysql-api)/" //mysql-api é o nome do serviço no docker-composer

	//DBSourceDocker Const used to acces Docker db
	DBSourceDocker = "root:my-secret-pw@tcp(mysql-api)/fairAPI"

	//DBisDropTableSQL Clear TAbles
	//DBisDropTableSQL = false
)

const (
	//DB_USER                = "clean_architecture_go_v2"
	//DB_PASSWORD            = "clean_architecture_go_v2"
	//DB_DATABASE            = "clean_architecture_go_v2"
	//DB_HOST                = "127.0.0.1"
	API_PORT = 8085
)

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}

func main() {
	//logs.Start()
	//prod := flag.Bool("prod", false, "write log to file")
	//debug := flag.Bool("debug", false, "enable debug mode")

	//logs.Debug("======== API FAIR Version %s \n", fair.Version())
	//handler.StartAPI(true)
	// dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	// db, err := sql.Open("mysql", dataSourceName)
	// if err != nil {
	// 	log.Panic(err.Error())
	// }
	// defer db.Close()

	//AddrOpenDB := DBSourceOpenDocker
	// repository.AddrOpenDB = DBSourceDocker
	// repository.AddrDB = DBSourceDocker
	repository.AddrOpenDB = DBSourceOpenLocal
	repository.AddrDB = DBSourceLocal

	db, err := sql.Open("mysql", repository.AddrOpenDB)
	if err != nil {
		log.Printf("Failed to connect to db Docker Mysql...")
		repository.AddrOpenDB = DBSourceOpenLocal
		repository.AddrDB = DBSourceLocal
		db, err = sql.Open("mysql", repository.AddrOpenDB)
		if err != nil {
			log.Printf("Failed to connect to db Local Mysql IP 127.0.0.1")
			log.Print(err)
		}

	}

	//defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	time.Sleep(5 * time.Second)
	//--------------------------------------------------------------
	fmt.Println("Conectado ao Banco")
	exec(db, "create database if not exists fairAPI")
	exec(db, "use fairAPI")
	// if isDropTable {
	// exec(db, "drop table if exists fair")
	// }

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

	fmt.Println("Successfully connected to the DB!")

	//bookRepo := repository.NewBookMySQL(db)
	//bookService := book.NewService(bookRepo)
	fairRepo := repository.NewFairMySQL(db)
	fairService := fair.NewService(fairRepo)

	//MakeFairHandlers(routerG, *n, fairService)

	// metricService, err := metric.NewPrometheusService()
	// if err != nil {
	// 	log.Panic(err.Error())
	// }
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.NewLogger(),
	)
	//book
	//handler.MakeBookHandlers(r, *n, bookService)
	handler.MakeFairHandlers(r, *n, fairService)

	http.Handle("/", r)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(API_PORT),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err.Error())
	}
	logs.LoggerClose()
}
