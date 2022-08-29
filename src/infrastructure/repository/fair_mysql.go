package repository

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

const (
	//TAB_FAIR_CREATE tabela da referente a entidade
	TAB_FAIR_CREATE = `create table IF NOT EXISTS fair(
		idx integer auto_increment,
		id varchar(50) ,
		name varchar(50),
		district varchar(18),
		region5 varchar(6),
		neighborhood varchar(20),
	    created_at datetime,
		updated_at datetime,
		PRIMARY KEY (idx)
		)`
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

func (r *FairMySQL) beginMysql() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		logs.Error("%s beginMysql() Err [%s]", logs.ThisFunction(), err.Error())
		return nil, err

	}
	return tx, nil
}

//Create Save fair data in DB
func (r *FairMySQL) Create(e *entity.Fair) (entity.ID, error) {

	tx, err := r.beginMysql()
	if err != nil {
		logs.Error("%s Err [%s] Could not access the DB ", logs.ThisFunction(), err.Error())
		return e.ID, err
	}

	stmt, err := tx.Prepare("insert into fair (id, name, district, region5, neighborhood, created_at)values(?,?,?,?,?,?)")
	if err != nil {
		logs.Error("%s Err [%s] Could not prepare transaction on DB ", logs.ThisFunction(), err.Error())
		return e.ID, err
	}

	_, err = stmt.Exec(
		e.ID,
		e.Name,
		e.District,
		e.Region5,
		e.Neighborhood,
		e.CreatedAt.Format("2006-01-02"),
	)
	if err != nil {
		tx.Rollback()
		logs.Error("%s Err [%s] Could not execute transaction on DB ", logs.ThisFunction(), err.Error())
		return e.ID, err
	}
	tx.Commit()

	defer stmt.Close()
	return e.ID, nil
}

//Get a Fair data from DB by ID
func (r *FairMySQL) Get(id entity.ID) (*entity.Fair, error) {

	tx, err := r.beginMysql()
	if err != nil {
		log.Println(err)
		logs.Error("%s Err [%s] Could not access the DB ", logs.ThisFunction(), err.Error())
		return nil, err
	}

	logs.Debug("Get(ID: %s)", id)

	stmt, err := tx.Prepare("select id, name, district, region5, neighborhood, created_at from fair where id = ?")
	if err != nil {
		logs.Error("%s Err [%s] Could not prepare transaction on DB ", logs.ThisFunction(), err.Error())
		return nil, err
	}

	var fair entity.Fair
	rows, err := stmt.Query(id)
	if err != nil {
		logs.Error("%s(ID: %s) Err [%s] Could not execute a Query on DB ", logs.ThisFunction(), id, err.Error())
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&fair.ID, &fair.Name, &fair.District, &fair.Region5, &fair.Neighborhood, &fair.CreatedAt)
		if err != nil {
			logs.Warn("%s(ID: %s) Err [%s] Could not \"SCAN\" on DB ", logs.ThisFunction(), id, err.Error())
		}
	}
	defer stmt.Close()

	return &fair, err
}

//Update Update fair data in DB
func (r *FairMySQL) Update(e *entity.Fair) error {
	e.UpdatedAt = time.Now()

	tx, err := r.beginMysql()
	if err != nil {
		logs.Error("%s Err [%s] Could not access the DB ", logs.ThisFunction(), err.Error())
		return err
	}

	stmt, err := tx.Prepare("update fair set name = ?, district = ?, region5 = ?, neighborhood = ?, updated_at = ? where id = ?")
	if err != nil {
		logs.Error("%s Err [%s] Could not prepare transaction on DB ", logs.ThisFunction(), err.Error())
		return err
	}

	logs.Debug("Update(ID: %v)", e.ID)
	_, err = stmt.Exec(e.Name, e.District, e.Region5, e.Neighborhood, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		tx.Rollback()
		logs.Error("%s Err [%s] Fail to \"UPDATE\" DATA in DB ID[%s]", logs.ThisFunction(), err.Error(), e.ID)
		return err
	}

	defer stmt.Close()
	tx.Commit()

	return nil

}

//Search fairs
func (r *FairMySQL) Search(key string, value string) ([]*entity.Fair, error) {

	logs.Debug("%s (key[%s]value[%s]) IN BD", logs.ThisFunction(), key, value)

	search_param := fmt.Sprintf("select id, name, district, region5, neighborhood, created_at from fair where %s like ?", key)
	stmt, err := r.db.Prepare(search_param)
	if err != nil {
		logs.Error("%s Err [%s] Could not prepare transaction on DB ", logs.ThisFunction(), err.Error())
		return nil, err
	}

	var listFair []*entity.Fair
	rows, err := stmt.Query("%" + value + "%")
	if err != nil {
		logs.Error("%s (value[%s]) Err [%s] Could not execute a Query on DB ", logs.ThisFunction(), value, err.Error())
		return nil, err
	}
	defer stmt.Close()

	for rows.Next() {
		var fair entity.Fair
		err = rows.Scan(&fair.ID, &fair.Name, &fair.District, &fair.Region5, &fair.Neighborhood, &fair.CreatedAt)
		if err != nil {
			logs.Warn("%s(key[%s]value[%s]) err[%s] Could not \"SCAN\" on DB ", logs.ThisFunction(), key, value, err.Error())
			return nil, err
		}

		listFair = append(listFair, &fair)
	}

	if len(listFair) == 0 {
		logs.Warn("%s WARNING: [%s] size[%v] \n", logs.ThisFunction(), entity.ErrEmptyDB.Error(), len(listFair))
		return nil, entity.ErrEmptyDB
	}

	return listFair, nil
}

//List All Fairs Save in DB
func (r *FairMySQL) List() ([]*entity.Fair, error) {

	rows, err := r.db.Query("select id, name, district, region5, neighborhood, created_at from fair")
	if err != nil {
		logs.Error(" %s Err [%s] Could not execute a Query on DB ", logs.ThisFunction(), err.Error())
		return nil, err
	}
	defer rows.Close()

	var listfair []*entity.Fair

	for rows.Next() {
		var fair entity.Fair
		err = rows.Scan(&fair.ID, &fair.Name, &fair.District, &fair.Region5, &fair.Neighborhood, &fair.CreatedAt)
		if err != nil {
			logs.Warn("%s err[%s] Could not \"SCAN\" on DB ", logs.ThisFunction(), err.Error())
			return nil, err
		}
		listfair = append(listfair, &fair)
	}

	if len(listfair) == 0 {
		logs.Warn("%s WARNING: [%s] size[%v] \n", logs.ThisFunction(), entity.ErrEmptyDB.Error(), len(listfair))
		return nil, entity.ErrEmptyDB
	}

	return listfair, nil
}

//Delete a Fair data from DB by ID
func (r *FairMySQL) Delete(id entity.ID) error {
	tx, err := r.beginMysql()
	if err != nil {
		logs.Error("%s Err [%s] Could not access the DB ", logs.ThisFunction(), err.Error())
		return err
	}

	stmt, err := tx.Prepare("delete from fair where id = ?")
	if err != nil {
		logs.Error("%s Err [%s] Could not prepare transaction on DB ", logs.ThisFunction(), err.Error())
		return err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		logs.Error("%s(id[%s]) Err [%s] Fail to \"DELETE\" DATA in DB", logs.ThisFunction(), id, err.Error())
		return err
	}
	defer stmt.Close()

	numElement, err := result.RowsAffected()
	if numElement == 0 || err != nil {
		errEmpty := errors.New("FAIR_TABLE: Not Found elements to this Type - Delete")
		return errEmpty
	}

	tx.Commit()
	return nil
}

//ImportFile Import data From CSV file To MySQL
func (r *FairMySQL) ImportFile(filepath string) error {

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
