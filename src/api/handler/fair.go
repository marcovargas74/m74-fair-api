package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-fair-api/src/api/presenter"
	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
	"github.com/marcovargas74/m74-fair-api/src/usecase/fair"
	"github.com/urfave/negroni"
)

func listFairs(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//servercpfcnpj.CreateDB()
		errorMessage := "Error reading Fair"
		var data []*entity.Fair
		var err error

		key := ""
		name := r.URL.Query().Get("name")
		if name != "" {
			key = "name"
		}

		district := r.URL.Query().Get("district")
		if district != "" {
			key = "district"
		}

		region5 := r.URL.Query().Get("region5")
		if region5 != "" {
			key = "region5"
		}

		neighborhood := r.URL.Query().Get("neighborhood")
		if neighborhood != "" {
			key = "neighborhood"
		}

		logs.Debug("listFairs key %s \n", key)
		switch {
		case key == "":
			data, err = service.ListFairs()
			logs.Debug("service ListFairs key %s \n", key)

		case key == "name":
			data, err = service.SearchFairs(key, name)
			logs.Debug("service ListFairs key %s \n", key)

		case key == "district":
			data, err = service.SearchFairs(key, district)
			logs.Debug("service ListFairs key %s \n", key)

		case key == "region5":
			data, err = service.SearchFairs(key, region5)
			logs.Debug("service ListFairs key %s \n", key)

		case key == "neighborhood":
			data, err = service.SearchFairs(key, neighborhood)
			logs.Debug("service ListFairs key %s \n", key)

		default:
			data, err = service.SearchFairs(key, name)
			logs.Debug("service SearchFairs key %s \n", name)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			//w.Write([]byte(errorMessage))
			fmt.Fprint(w, entity.ErrNotFound.Error())
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, entity.ErrNotFound.Error())
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
		fmt.Println("PASSOU AQUI")
		logs.Debug("createFair !!!!! %s \n", errorMessage)
		var input presenter.Fair
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("DecodeJSON1 %s \n", err.Error())
			return
		}
		defer r.Body.Close()

		if errs := input.Validate(); errs != nil {
			log.Printf("INVALIDO %v\n", errs)
			w.WriteHeader(http.StatusBadRequest)
		}

		id, err := service.CreateFair(input.Name, input.District, input.Region5, input.Neighborhood)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			logs.Error("InPUT DATA %s \n", errorMessage)
			return
		}

		input.ID = id
		json, err := json.Marshal(input)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusCreated)
	})
}

func getFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading fair"
		vars := mux.Vars(r)

		logs.Debug("getFair vars %v \n", vars)
		id, err := entity.StringToID(vars["id"])
		logs.Debug("getFair id [%v] \n", id)

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

func updateFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//errorMessage := "Error reading fair"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		//logs.Debug("UPDATE Fair id [%v] \n", id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, entity.ErrInvalidID.Error())
			logs.Error("ERRO: [%s] Nao conseguiu Converte o ID %v \n", entity.ErrInvalidID.Error(), id)
			return
		}

		dataToUpdate, err := service.GetFair(id)
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			logs.Error("ERRO: [%s] DESCONHECIDO ID %v \n", err.Error(), id)
			return
		}

		if dataToUpdate == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, entity.ErrNotFound.Error())
			logs.Error("ERRO: [%s] Nao achou o ID %v \n", entity.ErrInvalidID.Error(), id)
			return
		}

		//------Recuperar os dados Passados no Endpoint
		var newData presenter.Fair
		err = json.NewDecoder(r.Body).Decode(&newData)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("updateFair()Falha ao recuparar os dados passado no endpoint %s \n", err.Error())
			return
		}
		defer r.Body.Close()

		//Se for PUT todos os dados tem que estar aqui
		err = newData.Validate()
		if r.Method == http.MethodPut && err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprint(w, entity.ErrInvalidEntity.Error())
			logs.Error("Parametro(s) invalido(s) %s \n", err.Error())
			return
		}

		// ----atualiza os valores retornados .. TODO criar uma funcao auxiliar pra isso
		if newData.Name != "" {
			dataToUpdate.Name = newData.Name
		}

		if newData.District != "" {
			dataToUpdate.District = newData.District
		}

		if newData.Region5 != "" {
			dataToUpdate.Region5 = newData.Region5
		}

		if newData.Neighborhood != "" {
			dataToUpdate.Neighborhood = newData.Neighborhood
		}

		//---------Atualiza os dados no Banco
		if err = dataToUpdate.Validate(); err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprint(w, entity.ErrInvalidEntity.Error())
			logs.Error("Parametro(s) invalido(s) %s", err.Error())
			return
		}

		if err = service.UpdateFair(dataToUpdate); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, entity.ErrCannotBeUpdated.Error())
			logs.Error("updateFair()Falha ao atualizar dados %s", err.Error())
			return
		}

		json, err := json.Marshal(dataToUpdate)
		if err != nil {
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("updateFair()Falha ao converter o dado pra JSON %s", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusOK)
	})
}

func deleteFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing Fair"
		vars := mux.Vars(r)

		logs.Debug("deleteFair vars %v \n", vars)
		id, err := entity.StringToID(vars["id"])
		logs.Debug("deleteFair id [%v] \n", id)

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
		w.Write([]byte("Sucesso ao deletar Feira"))

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

	r.Handle("/fairs", n.With(negroni.Wrap(listFairs(service)))).Methods("GET", "OPTIONS").Name("listFairs")
	r.Handle("/fairs", n.With(negroni.Wrap(createFair(service)))).Methods("POST", "OPTIONS").Name("createFair")

	r.Handle("/fairs/{id}", n.With(negroni.Wrap(getFair(service)))).Methods("GET", "OPTIONS").Name("getFair")
	r.Handle("/fairs/{id}", n.With(negroni.Wrap(deleteFair(service)))).Methods("DELETE", "OPTIONS").Name("deleteFair")

	r.Handle("/fairs/{id}", n.With(negroni.Wrap(updateFair(service)))).Methods("PUT", "PATCH", "OPTIONS").Name("updateFair")

	/*
		func MakeBookHandlers(r *mux.Router, n negroni.Negroni, service book.UseCase) {
		r.Handle("/v1/book", n.With(
			negroni.Wrap(listBooks(service)),
		)).Methods("GET", "OPTIONS").Name("listBooks")

		r.Handle("/v1/book/{id}", n.With(
			negroni.Wrap(getBook(service)),
		)).Methods("GET", "OPTIONS").Name("getBook")

	//}*/

}

/*

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
