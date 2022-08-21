package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/api/presenter"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/entity"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/logs"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/usecase/fair"
	"github.com/urfave/negroni"
)

func CallbackListFair(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	fmt.Printf("Default data in %v\n", r.URL)
	log.Printf("METHOD[%s] DefaultEndpointFair \n", r.Method)
	w.WriteHeader(http.StatusAccepted)

}

func listFairs(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading Fair"
		var data []*entity.Fair
		var err error
		name := r.URL.Query().Get("id")
		logs.Debug("listFairs name %s \n", name)
		switch {
		case name == "":
			data, err = service.ListFairs()
			logs.Debug("service ListFairs name %s \n", name)
		default:
			data, err = service.SearchFairs(name)
			logs.Debug("service SearchFairs name %s \n", name)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*presenter.Fair
		for _, d := range data {
			toJ = append(toJ, &presenter.Fair{
				ID:           d.ID,
				Name:         d.Name,
				District:     d.District,
				Region5:      d.Region5,
				Neighborhood: d.Neighborhood,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func createFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding Fair"
		var input struct {
			Name         string `json:"name"`
			District     string `json:"district"`
			Region5      string `json:"region5"`
			Neighborhood string `json:"neighborhood"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			logs.Error("DecodeJSON1 %s \n", errorMessage)
			return
		}
		id, err := service.CreateFair(input.Name, input.District, input.Region5, input.Neighborhood)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			logs.Error("InPUT DATA %s \n", errorMessage)
			return
		}
		toJ := &presenter.Fair{
			ID:           id,
			Name:         input.Name,
			District:     input.District,
			Region5:      input.Region5,
			Neighborhood: input.Neighborhood,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			logs.Error("ENCODEJSON1 %s \n", errorMessage)
			return
		}
	})
}

func getFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading fair"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data, err := service.GetFair(id)
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &presenter.Fair{
			ID:           data.ID,
			Name:         data.Name,
			District:     data.District,
			Region5:      data.Region5,
			Neighborhood: data.Neighborhood,
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func deleteFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing Fair"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = service.DeleteFair(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func CallbackCreateFair(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
	fmt.Printf("Default data in %v\n", r.URL)
	log.Printf("METHOD[%s] DefaultEndpointFair \n", r.Method)
	w.WriteHeader(http.StatusAccepted)

}

//CallbackFairByName function Used to handle endpoint /{cpf_or_cnpj_num}
func CallbackFairGet(w http.ResponseWriter, r *http.Request) {

	//var aQueryJSON cpfcnpj.MyQuery
	aFairName := mux.Vars(r)
	logs.Debug("METHOD[%s] NAME in [%s] \n", r.Method, aFairName["name"])

	//cpfcnpj.CreateDB()
	//aQueryJSON.QuerysByNumHTTP(w, r, aCPFNum["cpf_or_cnpj_num"])
}

//CallbackFairByName function Used to handle endpoint /{cpf_or_cnpj_num}
func CallbackFairDelete(w http.ResponseWriter, r *http.Request) {

	//var aQueryJSON cpfcnpj.MyQuery
	aFairName := mux.Vars(r)
	logs.Debug("METHOD[%s] NAME in [%s] \n", r.Method, aFairName["name"])

	//cpfcnpj.CreateDB()
	//aQueryJSON.QuerysByNumHTTP(w, r, aCPFNum["cpf_or_cnpj_num"])
}

//NewServerAPI Create Server
func MakeFairHandlers(r *mux.Router, n negroni.Negroni, service fair.UseCase) {

	//r.HandleFunc("/fair", CallbackListFair).Methods("GET", "OPTIONS")

	//r.HandleFunc("/fair", listFairs(service)).Methods("GET", "OPTIONS")

	r.Handle("/fair", n.With(negroni.Wrap(listFairs(service)))).Methods("GET", "OPTIONS").Name("listFairs")
	r.Handle("/fair", n.With(negroni.Wrap(createFair(service)))).Methods("POST", "OPTIONS").Name("createFair")

	//	r.HandleFunc("/fair", CallbackCreateFair).Methods("POST", "OPTIONS")

	//r.HandleFunc("/fair/{id}", CallbackFairGet).Methods("GET", "OPTIONS")
	//r.HandleFunc("/fair/{name}", CallbackFairDelete).Methods("DELETE", "OPTIONS")
	r.Handle("/fair/{id}", n.With(negroni.Wrap(getFair(service)))).Methods("GET", "OPTIONS").Name("getFair")

	r.Handle("/fair/{id}", n.With(negroni.Wrap(deleteFair(service)))).Methods("DELETE", "OPTIONS").Name("deleteFair")

	//MakeBookHandlers make url handlers
	/*
		func MakeBookHandlers(r *mux.Router, n negroni.Negroni, service book.UseCase) {
		r.Handle("/v1/book", n.With(
			negroni.Wrap(listBooks(service)),
		)).Methods("GET", "OPTIONS").Name("listBooks")

		r.Handle("/v1/book", n.With(
			negroni.Wrap(createBook(service)),
		)).Methods("POST", "OPTIONS").Name("createBook")

		r.Handle("/v1/book/{id}", n.With(
			negroni.Wrap(getBook(service)),
		)).Methods("GET", "OPTIONS").Name("getBook")

		r.Handle("/v1/book/{id}", n.With(
			negroni.Wrap(deleteBook(service)),
		)).Methods("DELETE", "OPTIONS").Name("deleteBook")*/
	//}

}

/*
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

//package m74validatorapi

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"

// 	cpfcnpj "github.com/marcovargas74/m74-val-cpf-cnpj/src/cpf-cnpj"
// )

/*

func listBooks(service book.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading books"
		var data []*entity.Book
		var err error
		title := r.URL.Query().Get("title")
		switch {
		case title == "":
			data, err = service.ListBooks()
		default:
			data, err = service.SearchBooks(title)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*presenter.Book
		for _, d := range data {
			toJ = append(toJ, &presenter.Book{
				ID:       d.ID,
				Title:    d.Title,
				Author:   d.Author,
				Pages:    d.Pages,
				Quantity: d.Quantity,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func createBook(service book.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding book"
		var input struct {
			Title    string `json:"title"`
			Author   string `json:"author"`
			Pages    int    `json:"pages"`
			Quantity int    `json:"quantity"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		id, err := service.CreateBook(input.Title, input.Author, input.Pages, input.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &presenter.Book{
			ID:       id,
			Title:    input.Title,
			Author:   input.Author,
			Pages:    input.Pages,
			Quantity: input.Quantity,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getBook(service book.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading book"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data, err := service.GetBook(id)
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &presenter.Book{
			ID:       data.ID,
			Title:    data.Title,
			Author:   data.Author,
			Pages:    data.Pages,
			Quantity: data.Quantity,
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func deleteBook(service book.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing bookmark"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = service.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

*/
