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

//NewServerAPIMemory Create Default Routes in Memory
func NewServerAPIMemory() *ServerAPI {

	logs.Info("%s Create Default Routes", logs.ThisFunction())
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

//StartAPI_Memory Start http server which will provide API save data in memory
func StartAPI_Memory(portToAcessAPI string) {
	servidor := NewServerAPIMemory()

	logs.Info("%s Started successfully at port-> %v", logs.ThisFunction(), portToAcessAPI)
	if err := http.ListenAndServe(portToAcessAPI, servidor); err != nil {
		logs.Error("%sFailed to connect to port-> %v %v", logs.ThisFunction(), portToAcessAPI, err)
	}
}

//StartAPI_MySQL Start http server which will provide API save data in MySQL
func StartAPI_MySQL(portToAcessAPI string) {
	servidor := NewServerAPIMySQL()

	logs.Info("%s Started successfully at port-> %v", logs.ThisFunction(), portToAcessAPI)
	if err := http.ListenAndServe(portToAcessAPI, servidor); err != nil {
		logs.Error("%sFailed to connect to port-> %v %v", logs.ThisFunction(), portToAcessAPI, err)
	}
}

//NewServerAPIMySQL Create Default Routes in MySQL
func NewServerAPIMySQL() *ServerAPI {

	logs.Info("%s Create Default Routes", logs.ThisFunction())
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	n := negroni.New(
		negroni.NewLogger(),
	)

	repository.CreateDB(repository.NOT_DROP_DB)
	db := repository.OpenMysql()
	fairRepo := repository.NewFairMySQL(db)

	fairService := fair.NewService(fairRepo)
	MakeFairHandlers(routerG, *n, fairService)

	server.Handler = routerG
	return server

}
