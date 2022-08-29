package presenter

import (
	"encoding/json"
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

//NewCreateFairPresenter create fair to be used in the return of the endpoint
func NewCreateFairPresenter(e *entity.Fair) Fair {
	return Fair{
		ID:           e.ID,
		Name:         e.Name,
		District:     e.District,
		Region5:      e.Region5,
		Neighborhood: e.Neighborhood,
	}
}

//NewCreateFairPresenterJSON create fair(in JSON format) to be used in the return of the endpoint
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

//SelectKeySearch	Returns search parameter. If more than one is passed, returns the first
func SelectKeySearch(r *http.Request) (string, string) {
	logs.Debug("SelectKeySearch url[%v] ", r.URL)
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
	logs.Debug("SelectKeySearch value[%s] \n", value)

	return "", ""

}
