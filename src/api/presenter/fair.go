package presenter

import (
	"encoding/json"

	"github.com/marcovargas74/m74-fair-api/src/entity"
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
