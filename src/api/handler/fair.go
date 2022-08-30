package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcovargas74/m74-fair-api/src/api/presenter"
	"github.com/marcovargas74/m74-fair-api/src/config"
	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
	"github.com/marcovargas74/m74-fair-api/src/usecase/fair"
	"github.com/urfave/negroni"
)

func handlerID(w http.ResponseWriter, r *http.Request) (entity.ID, error) {
	vars := mux.Vars(r)
	id, err := entity.StringToID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, entity.ErrInvalidID.Error())
		logs.Error("ERROR: [%s] Failed to convert ID %v \n", entity.ErrInvalidID.Error(), id)
		return entity.ID(id), err
	}

	return entity.ID(id), err
}

func handlerSearchError(w http.ResponseWriter, data *entity.Fair, err error) error {
	if err != nil && err != entity.ErrNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logs.Error("Unknown ERROR: [%s] \n", err.Error())
		return err
	}

	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, entity.ErrNotFound.Error())
		logs.Info("INFO: [%s] No elements recorded in the list", entity.ErrNotFound.Error())
		return entity.ErrNotFound
	}
	return nil
}

func handlerValidateUpdate(w http.ResponseWriter, r *http.Request, newData presenter.Fair, dataToUpdate *entity.Fair) error {

	if r.Method == http.MethodPatch && newData.Name == "" && newData.District == "" && newData.Region5 == "" && newData.Neighborhood == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(w, entity.ErrInvalidEntity.Error())
		logs.Error("Invalid parameter(s) %s ", entity.ErrInvalidEntity.Error())
		return entity.ErrInvalidEntity
	}

	err := newData.Validate()
	if r.Method == http.MethodPut && err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(w, entity.ErrInvalidEntity.Error())
		logs.Error("%sInvalid parameter(s) %s", logs.ThisFunction(), err.Error())
		return err
	}

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

	if err = dataToUpdate.Validate(); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprint(w, err.Error())
		logs.Error("%sInvalid parameter(s) %s", logs.ThisFunction(), err.Error())
		return err
	}

	return nil
}

func listFairs(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var data []*entity.Fair
		var err error

		key, value := presenter.SelectKeySearch(r)
		if key == "" {
			data, err = service.ListFairs()
			logs.Debug("%s key[%s] value[%s] \n", logs.ThisFunction(), key, value)
		}

		if key != "" {
			data, err = service.SearchFairs(key, value)
			logs.Debug("%s key[%s] value[%s] ", logs.ThisFunction(), key, value)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			logs.Error("%s Unknown ERROR: [%s]", logs.ThisFunction(), err.Error())
			return
		}

		//----- Mosta o Retorno da busca
		var toJSONList []*presenter.Fair
		for _, d := range data {
			toJSON := presenter.NewCreateFairPresenter(d)
			toJSONList = append(toJSONList, &toJSON)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(toJSONList); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("%sFailed to convert data to JSON %s", logs.ThisFunction(), err.Error())
			return
		}

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
			logs.Error("%sFailed to retrieve past data from endpoint [%s] ", logs.ThisFunction(), err.Error())
			return
		}
		defer r.Body.Close()

		if err = input.Validate(); err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprint(w, entity.ErrInvalidEntity.Error())
			logs.Error("%sInvalid parameter(s) %s", logs.ThisFunction(), err.Error())
			return
		}

		id, err := service.CreateFair(input.Name, input.District, input.Region5, input.Neighborhood)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, entity.ErrCannotBeCreated.Error())
			logs.Error("%sFailed on Data Create %s", logs.ThisFunction(), err.Error())
			return
		}

		input.ID = id
		json, err := json.Marshal(input)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("%sFailed to convert data to JSON %s", logs.ThisFunction(), err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusCreated)
	})
}

func getFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := handlerID(w, r)
		if err != nil {
			return
		}

		logs.Debug("%sID:[%s]", logs.ThisFunction(), id)

		data, err := service.GetFair(id)
		if err = handlerSearchError(w, data, err); err != nil {
			return
		}

		json, err := presenter.NewCreateFairPresenterJSON(data)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("%sFailed to convert data to JSON %s", logs.ThisFunction(), err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusOK)
	})
}

func updateFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := handlerID(w, r)
		if err != nil {
			return
		}

		dataToUpdate, err := service.GetFair(id)
		if err = handlerSearchError(w, dataToUpdate, err); err != nil {
			logs.Error("%sERROR: [%s] ID [%v] did not find\n", logs.ThisFunction(), entity.ErrInvalidID.Error(), id)
			return
		}

		var newData presenter.Fair
		err = json.NewDecoder(r.Body).Decode(&newData)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("%sFailed to retrieve past data from endpoint %s ", logs.ThisFunction(), err.Error())
			return
		}
		defer r.Body.Close()

		if err = handlerValidateUpdate(w, r, newData, dataToUpdate); err != nil {
			return
		}

		if err = service.UpdateFair(dataToUpdate); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			logs.Error("%sFailed to update data %s", logs.ThisFunction(), err.Error())
			return
		}

		json, err := presenter.NewCreateFairPresenterJSON(dataToUpdate)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, entity.ErrCannotConvertJSON.Error())
			logs.Error("%sFailed to convert data to JSON %s", logs.ThisFunction(), err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusOK)
	})
}

func deleteFair(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := handlerID(w, r)
		if err != nil {
			return
		}

		err = service.DeleteFair(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, entity.ErrCannotBeDeleted.Error())
			logs.Error("%sData deleted - FAILED Err[%s]", logs.ThisFunction(), err.Error())
			return
		}
		fmt.Fprint(w, "Data deleted - SUCCESSFULLY")
		w.WriteHeader(http.StatusOK)

	})
}

func importFileCSV(service fair.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		file := (vars["file"])
		if file == "" {
			file = config.DEFAULT_CSV_FILE
		}

		err := service.ImportFileCSV(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			logs.Error("%s()Data imported FAILED from .CSV file %s", logs.ThisFunction(), err.Error())
			return
		}
		fmt.Fprint(w, "Data imported SUCCESSFULLY from .CSV file")
		w.WriteHeader(http.StatusOK)

	})
}

//MakeFairHandlers Cria Rotas usado para manipular a Feira
func MakeFairHandlers(r *mux.Router, n negroni.Negroni, service fair.UseCase) {

	r.Handle("/fairs", n.With(negroni.Wrap(listFairs(service)))).Methods("GET").Name("listFairs")
	r.Handle("/fairs", n.With(negroni.Wrap(createFair(service)))).Methods("POST").Name("createFair")
	r.Handle("/fairs/{id}", n.With(negroni.Wrap(getFair(service)))).Methods("GET").Name("getFair")
	r.Handle("/fairs/{id}", n.With(negroni.Wrap(updateFair(service)))).Methods("PUT", "PATCH").Name("updateFair")
	r.Handle("/fairs/{id}", n.With(negroni.Wrap(deleteFair(service)))).Methods("DELETE").Name("deleteFair")

	r.Handle("/fairs/import/{file}", n.With(negroni.Wrap(importFileCSV(service)))).Methods("POST").Name("updateFile")
}
