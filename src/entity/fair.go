package entity

import (
	"time"

	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
	"gopkg.in/validator.v2"
)

//Fair Entity data
type Fair struct {
	ID           ID
	Name         string `validate:"min=3,max=30,regexp=^[a-z A-Z]*$"`
	District     string `validate:"min=3,max=18,regexp=^[a-z A-Z]*$"`
	Region5      string `validate:"min=3,max=6,regexp=^[a-zA-Z]*$"`
	Neighborhood string `validate:"min=3,max=20,regexp=^[a-z A-Z]*$"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//NewFair create a new Fair
func NewFair(name string, district string, region5 string, neighborhood string) (*Fair, error) {
	f := &Fair{
		ID:           NewID(),
		Name:         name,
		District:     district,
		Region5:      region5,
		Neighborhood: neighborhood,
		CreatedAt:    time.Now(),
	}
	err := f.Validate()
	if err != nil {
		logs.Error("NewFair(*) Err [%s] Could not create a entityFair", err.Error())
		return nil, ErrInvalidEntity
	}
	return f, nil
}

//Validate validate Fair
func (f *Fair) Validate() error {

	if errs := validator.Validate(f); errs != nil {
		logs.Error("NewFair(*) Err [%s] Could not create a entityFair", errs.Error())
		return errs
	}

	return nil
}
