package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/repository"
	"github.com/marcovargas74/m74-fair-api/src/usecase/fair"
	"github.com/urfave/negroni"

	_ "github.com/go-sql-driver/mysql"
)

//ServerAPI is struct to start server
type ServerAPI struct {
	http.Handler
}

//DefaultEndpoint function Used to handle endpoint /- can be a load a page in html to configure
func (s *ServerAPI) DefaultEndpoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	logs.Info("METHOD[%s] DefaultEndpoint \n", r.Method)

	w.WriteHeader(http.StatusAccepted)

}

//NewServerAPIMemory Cria as Rotas Padrão
func NewServerAPIMemory() *ServerAPI {

	logs.Info("Cria as Rotas e salva dados em Memoria %v", "NewServerAPIMemory")
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	n := negroni.New(
		negroni.NewLogger(),
	)

	fairRepo := fair.NewInmem()
	fairService := fair.NewService(fairRepo)
	MakeFairHandlers(routerG, *n, fairService)

	server.Handler = routerG
	return server

}

//NewServerAPI Cria as Rotas Padrão e TRUE é em Memoria
/*func NewServerAPI(mode bool) *ServerAPI {
	if mode {
		return NewServerAPIMemory()
	}
	return NewServerAPIMYSQL()
}*/

//StartAPI http Inicia Servidor que vai prover a API
func StartAPI_Memory(portToAcessAPI string) {
	servidor := NewServerAPIMemory()

	logs.Info("Server Started successfully at port-> %v", portToAcessAPI)
	if err := http.ListenAndServe(portToAcessAPI, servidor); err != nil {
		logs.Error("Fail to conect in port-> %v %v", portToAcessAPI, err)
	}
}

//StartAPI http Inicia Servidor que vai prover a API
func StartAPI_MySQL(portToAcessAPI string) {
	servidor := NewServerAPIMYSQL()

	logs.Info("Server MYSQL Started successfully at port-> %v", portToAcessAPI)
	if err := http.ListenAndServe(portToAcessAPI, servidor); err != nil {
		logs.Error("Fail to conect in port-> %v %v", portToAcessAPI, err)
	}
}

//----------------------------------------------------------------------------------

//TODO refatorar essa parte
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

//AddrOpenDB VAR used to open and to access BD
//var AddrOpenDB string

//AddrDB VAR data source name
//var AddrDB string

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}

//NewServerAPI Cria as Rotas Padrão
func NewServerAPIMYSQL() *ServerAPI {

	logs.Info("Cria as Rotas %v", "NewServerAPIMYSQL")
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	n := negroni.New(
		negroni.NewLogger(),
	)

	//--------------SE CONECTA AO BANCO
	//dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	//db, err := sql.Open("mysql", dataSourceName)

	//AddrOpenDB = DBSourceOpenDocker
	//AddrDB = DBSourceDocker
	repository.AddrOpenDB = DBSourceOpenLocal
	repository.AddrDB = DBSourceLocal

	db, err := sql.Open("mysql", repository.AddrOpenDB)
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
	fairRepo := repository.NewFairMySQL(db)

	fairService := fair.NewService(fairRepo)
	MakeFairHandlers(routerG, *n, fairService)

	server.Handler = routerG
	return server

}
