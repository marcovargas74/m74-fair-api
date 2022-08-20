package repository

import (
	"database/sql"
	"time"

	"github.com/marcovargas74/m74-val-cpf-cnpj/src/entity"
)

//FairMySQL mysql repo
type FairMySQL struct {
	db *sql.DB
}

//NewFairMySQL create new repository
func NewFairMySQL(db *sql.DB) *FairMySQL {
	return &FairMySQL{
		db: db,
	}
}

//Create a Fair
func (r *FairMySQL) Create(e *entity.Fair) (entity.ID, error) {
	stmt, err := r.db.Prepare(`
		insert into fair (id, name, district, region5, neighborhood, created_at)
		values(?,?,?,?,?,?)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Name,
		e.District,
		e.Region5,
		e.Neighborhood,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

//Get a Fair
func (r *FairMySQL) Get(id entity.ID) (*entity.Fair, error) {
	stmt, err := r.db.Prepare(`id, name, district, region5, neighborhood, created_at from fair where id = ?`)
	if err != nil {
		return nil, err
	}
	var f entity.Fair
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&f.ID, &f.Name, &f.District, &f.Region5, &f.Neighborhood, &f.CreatedAt)
	}

	return &f, err
}

//Update a Fair
func (r *FairMySQL) Update(e *entity.Fair) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec("update fair set name = ?, district = ?, region5 = ?, neighborhood = ?, updated_at = ? where id = ?", e.Name, e.District, e.Region5, e.Neighborhood, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		//TODO colocar log de erro
		return err
	}
	return nil
}

//Search fairs
func (r *FairMySQL) Search(query string) ([]*entity.Fair, error) {
	stmt, err := r.db.Prepare(`select id, name, district, region5, neighborhood, created_at from fair where name like ?`)
	if err != nil {
		return nil, err
	}
	var fairs []*entity.Fair
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f entity.Fair
		err = rows.Scan(&f.ID, &f.Name, &f.District, &f.Region5, &f.Neighborhood, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		fairs = append(fairs, &f)
	}

	return fairs, nil
}

//List Fairs
func (r *FairMySQL) List() ([]*entity.Fair, error) {
	stmt, err := r.db.Prepare(`select id, name, district, region5, neighborhood, created_at from fair`)
	if err != nil {
		return nil, err
	}
	var fairs []*entity.Fair
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f entity.Fair
		err = rows.Scan(&f.ID, &f.Name, &f.District, &f.Region5, &f.Neighborhood, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		fairs = append(fairs, &f)
	}
	return fairs, nil
}

//Delete a Fair
func (r *FairMySQL) Delete(id entity.ID) error {
	_, err := r.db.Exec("delete from fair where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
