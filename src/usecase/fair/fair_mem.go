package fair

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

const (
	INDEX_NAME         = 11
	INDEX_DISTRICT     = 6
	INDEX_REGION5      = 9
	INDEX_NEIGHBORHOOD = 15
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
		logs.Warn("%s(ID: %s) Err: Could not \"SCAN\" on MEMORY ", logs.ThisFunction(), id)
		return nil, entity.ErrNotFound
	}
	return r.m[id], nil
}

//Update a Fair
func (r *inmem) Update(e *entity.Fair) error {
	_, err := r.Get(e.ID)
	if err != nil {
		logs.Error("%s Err [%s] Fail to \"UPDATE\" DATA in MEMERY[%s]", logs.ThisFunction(), err.Error(), e.ID)
		return err
	}
	r.m[e.ID] = e
	return nil
}

//Search Fair
func (r *inmem) Search(key string, value string) ([]*entity.Fair, error) {
	var d []*entity.Fair

	for _, j := range r.m {

		switch {
		case key == "name":
			if strings.Contains(strings.ToLower(j.Name), value) {
				d = append(d, j)
			}

		case key == "district":
			if strings.Contains(strings.ToLower(j.District), value) {
				d = append(d, j)
			}

		case key == "region5":
			if strings.Contains(strings.ToLower(j.Region5), value) {
				d = append(d, j)
			}

		case key == "neighborhood":
			if strings.Contains(strings.ToLower(j.Neighborhood), value) {
				d = append(d, j)
			}

		default:
			if strings.Contains(strings.ToLower(j.Name), value) {
				d = append(d, j)
			}
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
		logs.Error("%s(id[%s]) Err Fail to \"DELETE\" DATA in MEMORY", logs.ThisFunction(), id)
		return entity.ErrNotFound
	}
	r.m[id] = nil
	return nil
}

//ImportFile Import data From CSV file To MySQL
func (r *inmem) ImportFile(filepath string) error {

	csvFile, err := os.Open(filepath)
	if err != nil {
		logs.Error("%s Err [%s] Could not Open log FILE", logs.ThisFunction(), err.Error())
		return err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ','
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			logs.Error("%s Err [%s] Could not Read a Line from CSV FILE", logs.ThisFunction(), err.Error())
			continue
		}

		fair, err := entity.NewFair(line[INDEX_NAME], line[INDEX_DISTRICT], line[INDEX_REGION5], line[INDEX_NEIGHBORHOOD])
		if err != nil {
			logs.Error("%s Err [%s] Could not Read create a entity from CSV FILE", logs.ThisFunction(), err.Error())
			continue
		}
		r.Create(fair)

		logs.Debug("%s entityFAIR [%v] READ from .CSV file\n", logs.ThisFunction(), fair)

	}

	return nil
}
