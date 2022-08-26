package handler

import (
	"net/http"

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
// const (
// 	//DBSourceOpenLocal Const used to Open Local db
// 	DBSourceOpenLocal = "root:my-secret-pw@tcp(localhost:3307)/"

// 	//DBSourceLocal Const used to acces Local db
// 	DBSourceLocal = "root:my-secret-pw@tcp(localhost:3307)/fairAPI?parseTime=true"

// 	//DBSourceOpenDocker Const used to Open Docker db
// 	//DBSourceOpenDocker = "root:my-secret-pw@tcp(mysql-api)/" //mysql-api é o nome do serviço no docker-composer

// 	//DBSourceDocker Const used to acces Docker db
// 	//DBSourceDocker = "root:my-secret-pw@tcp(mysql-api)/fairAPI"

// )

//AddrOpenDB VAR used to open and to access BD
//var AddrOpenDB string

//AddrDB VAR data source name
//var AddrDB string

/*

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}

*/

//NewServerAPI Cria as Rotas Padrão
func NewServerAPIMYSQL() *ServerAPI {

	logs.Info("Cria as Rotas %v", "NewServerAPIMYSQL")
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	n := negroni.New(
		negroni.NewLogger(),
	)

	repository.CreateDB()
	db := repository.OpenMysql()
	fairRepo := repository.NewFairMySQL(db)
	//fairRepo := repository.NewFairMySQL()

	fairService := fair.NewService(fairRepo)
	MakeFairHandlers(routerG, *n, fairService)

	server.Handler = routerG
	return server

}
