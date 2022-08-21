package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/logs"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/repository"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/usecase/fair"
	"github.com/urfave/negroni"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serverPort = ":5000"
)

//ServerAPI is struct to start server
type ServerAPI struct {
	http.Handler
}

//DefaultEndpoint function Used to handle endpoint /- can be a load a page in html to configure
func (s *ServerAPI) DefaultEndpoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	fmt.Printf("Default data in %v\n", r.URL)
	log.Printf("METHOD[%s] DefaultEndpoint \n", r.Method)

	w.WriteHeader(http.StatusAccepted)

}

/*
//TODO refatorar essa parte
const (
	//DBSourceOpenLocal Const used to Open Local db
	DBSourceOpenLocal = "root:my-secret-pw@tcp(localhost:3307)/"

	//DBSourceLocal Const used to acces Local db
	DBSourceLocal = "root:my-secret-pw@tcp(localhost:3307)/fairAPI" //root:Mysql#my-secret-pw@/fairAPI"

	//DBSourceOpenDocker Const used to Open Docker db
	DBSourceOpenDocker = "root:my-secret-pw@tcp(mysql-api)/" //mysql-api é o nome do serviço no docker-composer

	//DBSourceDocker Const used to acces Docker db
	DBSourceDocker = "root:my-secret-pw@tcp(mysql-api)/fairAPI" //root:Mysql#my-secret-pw@/fairAPI"

)

//AddrOpenDB VAR used to open and to access BD
var AddrOpenDB string

//AddrDB VAR data source name
var AddrDB string

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}
*/
//NewServerAPI Cria as Rotas Padrão
func NewServerAPI(mode bool) *ServerAPI {

	logs.Info("Cria as Rotas %v", mode)
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	n := negroni.New(
		negroni.NewLogger(),
	)

	/*

	   	//--------------SE CONECTA AO BANCO
	   	//dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	   	//db, err := sql.Open("mysql", dataSourceName)

	   	AddrOpenDB = DBSourceOpenDocker
	   	AddrDB = DBSourceDocker

	   	db, err := sql.Open("mysql", AddrOpenDB)
	   	if err != nil {
	   		logs.Error("FALHA ao conectar ao Banco Mysql do DOcker %v", err)
	   		AddrOpenDB = DBSourceOpenLocal
	   		AddrDB = DBSourceLocal
	   		db, err = sql.Open("mysql", AddrOpenDB)
	   		if err != nil {
	   			logs.Error("FALHA ao conectar ao Banco Mysql Local IP 127.0.0.1 %v", err)
	   		}
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
	   	exec(db, "drop table if exists fair")
	   	//}

	   	exec(db, `create table IF NOT EXISTS fair(
	   	idx integer auto_increment,
	   	id varchar(50) ,
	   	name varchar(50),
	   	district varchar(18),
	   	region5 varchar(6),
	   	neighborhood varchar(20),
	       createAt datetime,
	   	updated_at datetime,
	   	PRIMARY KEY (idx)
	   	)`)
	*/
	var db *sql.DB
	//repository.CreateDB()
	fairRepo := repository.NewFairMySQL(db)

	fairService := fair.NewService(fairRepo)
	MakeFairHandlers(routerG, *n, fairService)

	server.Handler = routerG
	return server

}

//StartAPI http Inicia Servidor que vai prover a API
func StartAPI(mode bool) {
	servidor := NewServerAPI(mode)

	logs.Info("Server Started successfully at port-> %v", serverPort)
	if err := http.ListenAndServe(serverPort, servidor); err != nil {
		logs.Error("Fail to conect in port-> %v %v", serverPort, err)
	}
}
