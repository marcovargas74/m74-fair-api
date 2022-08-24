package presenter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

//Fair data
type Fair struct {
	ID           entity.ID `json:"id"`
	Name         string    `json:"name"`
	District     string    `json:"district"`
	Region5      string    `json:"region5"`
	Neighborhood string    `json:"neighborhood"`
}

//NewCreateFairPresenter .
func NewCreateFairPresenter(e *entity.Fair) Fair {
	return Fair{
		ID:           e.ID,
		Name:         e.Name,
		District:     e.District,
		Region5:      e.Region5,
		Neighborhood: e.Neighborhood,
	}
}

//NewCreateFairPresenterJSON .
func NewCreateFairPresenterJSON(e *entity.Fair) ([]byte, error) {

	toJSON := NewCreateFairPresenter(e)
	return json.Marshal(toJSON)
}

//Validate validate presenter
func (f *Fair) Validate() error {
	if f.Name == "" || f.District == "" || f.Region5 == "" || f.Neighborhood == "" {
		return entity.ErrInvalidEntity
	}
	return nil
}

//ConvertToFairPresenter presenter
func (f *Fair) ConvertToFairPresenter(e *entity.Fair) Fair {
	return Fair{
		ID:           e.ID,
		Name:         e.Name,
		District:     e.District,
		Region5:      e.Region5,
		Neighborhood: e.Neighborhood,
	}
}

//SelectKeySearch	Retorna parametro de busca. Caso for passado mais de um, retorna o primeiro
func SelectKeySearch(r *http.Request) (string, string) {
	//--seleciona chave de busca---

	logs.Debug("SelectKeySearch url[%v] ", r.URL)
	fmt.Printf("SelectKeySearch url[%v] \n", r.URL)
	value := r.URL.Query().Get("name")
	if value != "" {
		return "name", value
	}

	value = r.URL.Query().Get("district")
	if value != "" {
		return "district", value
	}

	value = r.URL.Query().Get("region5")
	if value != "" {
		return "region5", value
	}

	value = r.URL.Query().Get("neighborhood")
	if value != "" {
		return "neighborhood", value
	}
	//logs.Debug("listFairs key[%s] value[%s] \n", key, value)

	return "", ""

}
