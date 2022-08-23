package handler

import (
	"encoding/json"
	"fmt"
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
		//errorMessage := "Error reading Fair"
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
			fmt.Fprint(w, err.Error())
			logs.Error("ERRO: [%s] DESCONHECIDO \n", err.Error())
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, entity.ErrNotFound.Error())
			logs.Info("Info: [%s] Nenhum elemento gravado na lista \n", entity.ErrInvalidID.Error())
			return
		}

		var toJSONList []*presenter.Fair
		for _, d := range data {
			toJSON := presenter.NewCreateFairPresenter(d)
			toJSONList = append(toJSONList, &toJSON)
		}
		if err := json.NewEncoder(w).Encode(toJSONList); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("listFairs()Falha ao converter o dado pra JSON %s", err.Error())
			return
		}

		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func createFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input presenter.Fair
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("createFair()Falha ao recuparar os dados passado no endpoint %s \n", err.Error())
			return
		}
		defer r.Body.Close()

		if err = input.Validate(); err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprint(w, entity.ErrInvalidEntity.Error())
			logs.Error("Parametro(s) invalido(s) %s", err.Error())
			return
		}

		id, err := service.CreateFair(input.Name, input.District, input.Region5, input.Neighborhood)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, entity.ErrCannotBeCreated.Error())
			logs.Error("createFair()Falha ao Criar dados %s", err.Error())
			return
		}

		input.ID = id
		json, err := json.Marshal(input)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("createFair()Falha ao converter o dado pra JSON %s", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusCreated)
	})
}

func getFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, entity.ErrInvalidID.Error())
			logs.Error("ERRO: [%s] Nao conseguiu Converte o ID %v \n", entity.ErrInvalidID.Error(), id)
			return
		}

		data, err := service.GetFair(id)
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			logs.Error("ERRO: [%s] DESCONHECIDO ID %v \n", err.Error(), id)
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, entity.ErrNotFound.Error())
			logs.Error("ERRO: [%s] Nao achou o ID %v \n", entity.ErrInvalidID.Error(), id)
			return
		}

		json, err := presenter.NewCreateFairPresenterJSON(data)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("getFair()Falha ao converter o dado pra JSON %s", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusOK)
	})
}

func updateFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//errorMessage := "Error reading fair"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
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

		if r.Method == http.MethodPatch && newData.Name == "" && newData.District == "" && newData.Region5 == "" && newData.Neighborhood == "" {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprint(w, entity.ErrInvalidEntity.Error())
			logs.Error("Parametro(s) invalido(s) %s \n", err.Error())
			return
		}

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

		json, err := presenter.NewCreateFairPresenterJSON(dataToUpdate)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
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
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		logs.Debug("deleteFair id [%v] \n", id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, entity.ErrInvalidID.Error())
			logs.Error("ERRO: [%s] Nao conseguiu Converte o ID %v \n", entity.ErrInvalidID.Error(), id)
			return
		}

		err = service.DeleteFair(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, entity.ErrCannotBeDeleted.Error())
			logs.Error("deleteFair()Falha ao DELETAR dados %s", err.Error())
			return
		}
		w.Write([]byte("Sucesso ao deletar dado"))
		w.WriteHeader(http.StatusOK)

	})
}

//NewServerAPI Create Server
func MakeFairHandlers(r *mux.Router, n negroni.Negroni, service fair.UseCase) {

	r.Handle("/fairs", n.With(negroni.Wrap(listFairs(service)))).Methods("GET", "OPTIONS").Name("listFairs")
	r.Handle("/fairs", n.With(negroni.Wrap(createFair(service)))).Methods("POST", "OPTIONS").Name("createFair")

	r.Handle("/fairs/{id}", n.With(negroni.Wrap(getFair(service)))).Methods("GET", "OPTIONS").Name("getFair")
	r.Handle("/fairs/{id}", n.With(negroni.Wrap(deleteFair(service)))).Methods("DELETE", "OPTIONS").Name("deleteFair")

	r.Handle("/fairs/{id}", n.With(negroni.Wrap(updateFair(service)))).Methods("PUT", "PATCH", "OPTIONS").Name("updateFair")
}
