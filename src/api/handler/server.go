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

//NewServerAPIMemory Cria as Rotas Padrão em Memoria
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

//StartAPI_Memory Inicia Servidor http que vai prover a API salva dados em memoria
func StartAPI_Memory(portToAcessAPI string) {
	servidor := NewServerAPIMemory()

	logs.Info("Server Started successfully at port-> %v", portToAcessAPI)
	if err := http.ListenAndServe(portToAcessAPI, servidor); err != nil {
		logs.Error("Fail to conect in port-> %v %v", portToAcessAPI, err)
	}
}

//StartAPI_MySQL Inicia Servidor http que vai prover a API salva os dados no MySQL
func StartAPI_MySQL(portToAcessAPI string) {
	servidor := NewServerAPIMYSQL()

	logs.Info("Server MYSQL Started successfully at port-> %v", portToAcessAPI)
	if err := http.ListenAndServe(portToAcessAPI, servidor); err != nil {
		logs.Error("Fail to conect in port-> %v %v", portToAcessAPI, err)
	}
}

//----------------------------------------------------------------------------------

//NewServerAPIMYSQL Cria as Rotas Padrão
func NewServerAPIMYSQL() *ServerAPI {

	logs.Info("Cria as Rotas %v", "NewServerAPIMYSQL")
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	n := negroni.New(
		negroni.NewLogger(),
	)

	repository.CreateDB(repository.NOT_DROP_DB)
	//repository.CreateDB(repository.DROP_DB)
	db := repository.OpenMysql()
	fairRepo := repository.NewFairMySQL(db)

	fairService := fair.NewService(fairRepo)
	MakeFairHandlers(routerG, *n, fairService)

	server.Handler = routerG
	return server

}
