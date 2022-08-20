package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/logs"
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

//NewServerAPI Cria as Rotas Padrão
func NewServerAPI(mode bool) *ServerAPI {

	logs.Info("Cria as Rotas %v", mode)
	server := new(ServerAPI)
	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")

	MakeFairHandlers(routerG)

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
