package fair

import (
	"strings"

	"github.com/marcovargas74/m74-fair-api/src/entity"
)

//inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Fair
}

//NewInmem create new repository
func NewInmem() *inmem {
	var m = map[entity.ID]*entity.Fair{}
	return &inmem{
		m: m,
	}
}

//Create a Fair
func (r *inmem) Create(e *entity.Fair) (entity.ID, error) {
	r.m[e.ID] = e
	return e.ID, nil
}

//Get a Fair
func (r *inmem) Get(id entity.ID) (*entity.Fair, error) {
	if r.m[id] == nil {
		return nil, entity.ErrNotFound
	}
	return r.m[id], nil
}

//Update a Fair
func (r *inmem) Update(e *entity.Fair) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = e
	return nil
}

//Search Fair
func (r *inmem) Search(key string, value string) ([]*entity.Fair, error) {
	var d []*entity.Fair
	for _, j := range r.m {
		if strings.Contains(strings.ToLower(j.Name), value) {
			d = append(d, j)
		}
	}
	return d, nil
}

//List Fair
func (r *inmem) List() ([]*entity.Fair, error) {
	var d []*entity.Fair
	for _, j := range r.m {
		if j != nil {
			d = append(d, j)
		}
	}
	return d, nil
}

//Delete a Fair
func (r *inmem) Delete(id entity.ID) error {
	if r.m[id] == nil {
		return entity.ErrNotFound
	}
	r.m[id] = nil
	return nil
}
