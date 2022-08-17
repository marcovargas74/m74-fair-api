package m74validatorapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	cpfcnpj "github.com/marcovargas74/m74-val-cpf-cnpj/src/cpf-cnpj"
)

const (
	serverPort = ":5000"
)

//ServerValidator is struct to start server
type ServerValidator struct {
	http.Handler
}

//CallbackStatus function Used to handle endpoint /status
func (s *ServerValidator) CallbackStatus(w http.ResponseWriter, r *http.Request) {

	log.Printf("METHOD[%s] STATUS [%s] \n", r.Method, r.UserAgent())
	cpfcnpj.ShowStatus(w, r)
	w.WriteHeader(http.StatusOK)

}

//CallbackQuerysAll function Used to handle endpoint /all
func (s *ServerValidator) CallbackQuerysAll(w http.ResponseWriter, r *http.Request) {

	var aQueryJSON cpfcnpj.MyQuery
	log.Printf("METHOD[%s] SHOW ALL CPF AND CNPJs \n", r.Method)
	aQueryJSON.QuerysHTTP(w, r)
	cpfcnpj.UpdateStatus()

}

//CallbackQuerysCPFAll function Used to handle endpoint /cpfs/
func (s *ServerValidator) CallbackQuerysCPFAll(w http.ResponseWriter, r *http.Request) {

	var aQueryJSON cpfcnpj.MyQuery
	log.Printf("METHOD[%s] SHOW ALL CPFs \n", r.Method)

	aQueryJSON.QuerysByTypeHTTP(w, r, cpfcnpj.IsCPF)
	cpfcnpj.UpdateStatus()

}

//CallbackQuerysByNum function Used to handle endpoint /{cpf_or_cnpj_num}
func (s *ServerValidator) CallbackQuerysByNum(w http.ResponseWriter, r *http.Request) {

	var aQueryJSON cpfcnpj.MyQuery
	aCPFNum := mux.Vars(r)
	log.Printf("METHOD[%s] CPF in [%s] \n", r.Method, aCPFNum["cpf_or_cnpj_num"])

	cpfcnpj.CreateDB()

	aQueryJSON.QuerysByNumHTTP(w, r, aCPFNum["cpf_or_cnpj_num"])
	cpfcnpj.UpdateStatus()

}

//CallbackQuerysCPF function Used to handle endpoint /cpfs/{cpf}
func (s *ServerValidator) CallbackQuerysCPF(w http.ResponseWriter, r *http.Request) {

	var aQueryJSON cpfcnpj.MyQuery

	aCPFNum := mux.Vars(r)
	log.Printf("METHOD[%s] CPF in [%s] \n", r.Method, aCPFNum["cpf_num"])

	cpfcnpj.CreateDB()

	switch r.Method {
	case http.MethodPost:
		aQueryJSON.SaveQueryHTTP(w, r, aCPFNum["cpf_num"], cpfcnpj.IsCPF)
		log.Printf("WriteHeader %v\n", w)
		cpfcnpj.UpdateStatus()

	case http.MethodGet:
		cpfcnpj.UpdateStatus()
		if len(aCPFNum) == 0 {
			aQueryJSON.QuerysByTypeHTTP(w, r, cpfcnpj.IsCPF)
			return
		}

		aQueryJSON.SaveQueryHTTP(w, r, aCPFNum["cpf_num"], cpfcnpj.IsCPF)

	case http.MethodDelete:
		aQueryJSON.DeleteQuerysByNumHTTP(w, r, aCPFNum["cpf_num"], cpfcnpj.IsCPF)

	}

}

//CallbackQuerysCNPJAll function Used to handle endpoint /cnpj
func (s *ServerValidator) CallbackQuerysCNPJAll(w http.ResponseWriter, r *http.Request) {

	var aQueryJSON cpfcnpj.MyQuery
	log.Printf("METHOD[%s] SHOW ALL CNPJs \n", r.Method)

	aQueryJSON.QuerysByTypeHTTP(w, r, cpfcnpj.IsCNPJ)
	cpfcnpj.UpdateStatus()

}

//CallbackQuerysCNPJ function Used to handle endpoint /cnpjs
func (s *ServerValidator) CallbackQuerysCNPJ(w http.ResponseWriter, r *http.Request) {

	var aQueryJSON cpfcnpj.MyQuery
	aCNPJ := mux.Vars(r)
	argCNPJ := fmt.Sprintf("%s/%s", aCNPJ["cnpj_num"], aCNPJ["cnpj_part2"])
	log.Printf("METHOD[%s] CallbackQuerysCNPJ [%s] \n", r.Method, argCNPJ)

	cpfcnpj.CreateDB()

	switch r.Method {
	case http.MethodPost:
		cpfcnpj.UpdateStatus()
		aQueryJSON.SaveQueryHTTP(w, r, argCNPJ, cpfcnpj.IsCNPJ)

	case http.MethodGet:
		cpfcnpj.UpdateStatus()
		if len(aCNPJ) == 0 {
			aQueryJSON.QuerysByTypeHTTP(w, r, cpfcnpj.IsCNPJ)
			return
		}
		aQueryJSON.SaveQueryHTTP(w, r, argCNPJ, cpfcnpj.IsCNPJ)

	case http.MethodDelete:
		aQueryJSON.DeleteQuerysByNumHTTP(w, r, argCNPJ, cpfcnpj.IsCNPJ)

	}

}

//DefaultEndpoint function Used to handle endpoint /- can be a load a page in html to configure
func (s *ServerValidator) DefaultEndpoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	fmt.Printf("Default data in %v\n", r.URL)
	w.WriteHeader(http.StatusAccepted)

}

//NewServerValidator Create Server
func NewServerValidator(mode string) *ServerValidator {

	server := new(ServerValidator)

	routerG := mux.NewRouter()
	routerG.HandleFunc("/", server.DefaultEndpoint).Methods("GET")
	routerG.HandleFunc("/status", server.CallbackStatus).Methods("GET")
	routerG.HandleFunc("/all", server.CallbackQuerysAll).Methods("GET")
	routerG.HandleFunc("/num/{cpf_or_cnpj_num}", server.CallbackQuerysByNum).Methods("GET")

	//Routes CPF
	routerG.HandleFunc("/cpfs", server.CallbackQuerysCPFAll).Methods("GET")
	routerG.HandleFunc("/cpfs/{cpf_num}", server.CallbackQuerysCPF).Methods("GET", "POST", "DELETE")

	//Routes CNPJ
	routerG.HandleFunc("/cnpjs", server.CallbackQuerysCNPJAll).Methods("GET")
	routerG.HandleFunc("/cnpjs/{cnpj_num}", server.CallbackQuerysCNPJ).Methods("GET", "POST", "DELETE")
	routerG.HandleFunc("/cnpjs/{cnpj_num}/{cnpj_part2}", server.CallbackQuerysCNPJ).Methods("GET", "POST", "DELETE")
	server.Handler = routerG
	return server

}

//StartAPI http starter server
func StartAPI(mode string) {
	servidor := NewServerValidator(mode)

	if err := http.ListenAndServe(serverPort, servidor); err != nil {
		log.Printf("Fail to conect in a port-> 5000 %v", err)
	}
}
