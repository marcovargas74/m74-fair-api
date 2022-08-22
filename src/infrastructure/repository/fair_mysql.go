package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

//AddrOpenDB VAR used to open and to access BD
var AddrOpenDB string

//AddrDB VAR data source name
var AddrDB string

//FairMySQL mysql repo
type FairMySQL struct {
	db *sql.DB
}

//NewFairMySQL create new repository
func NewFairMySQL(db *sql.DB) *FairMySQL {
	//repository.CreateDB()

	return &FairMySQL{
		db: db,
	}
}

//Create a Fair
func (r *FairMySQL) Create(e *entity.Fair) (entity.ID, error) {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return e.ID, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	logs.Debug("MySQL CREATE BD %s \n", AddrDB)

	// stmt, err := r.db.Prepare(`
	// 	insert into fair (id, name, district, region5, neighborhood, created_at)
	// 	values(?,?,?,?,?,?)`)

	stmt, err := tx.Prepare("insert into fair (id, name, district, region5, neighborhood, created_at)values(?,?,?,?,?,?)")
	if err != nil {
		logs.Error("Insert DATA in Mysql %v \n", err)
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
		tx.Rollback()
		logs.Error("exec in Mysql %v \n", err)
		return e.ID, err
	}
	tx.Commit()

	err = stmt.Close()
	if err != nil {
		logs.Error("Close stmt in Mysql %v \n", err)
		return e.ID, err
	}
	return e.ID, nil
}

//Get a Fair
func (r *FairMySQL) Get(id entity.ID) (*entity.Fair, error) {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	stmt, err := tx.Prepare(`id, name, district, region5, neighborhood, created_at from fair where id = ?`)
	//stmt, err := r.db.Prepare(`id, name, district, region5, neighborhood, created_at from fair where id = ?`)
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

	if len(fairs) == 0 {
		errEmpty := errors.New("FAIR_TABLE: is Empty")
		return nil, errEmpty
	}

	return fairs, nil
}

//List Fairs
func (r *FairMySQL) List() ([]*entity.Fair, error) {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	logs.Debug("MySQL LIST BD %s \n", AddrDB)

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Println(err)
	// }

	logs.Debug("List MySQL %s \n", "init")

	//stmt, err := tx.Prepare(`select id, name, district, region5, neighborhood, created_at from fair`)
	//stmt, err := r.db.Prepare(`select id, name, district, region5, neighborhood, created_at from fair`)

	rows, err := db.Query("select id, name, district, region5, neighborhood, created_at from fair")
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	var fairs []*entity.Fair
	// rows, err := stmt.Query()
	// if err != nil {
	// 	return nil, err
	// }

	logs.Debug("  List MySQL rows-> %v \n", rows)

	for rows.Next() {
		var f entity.Fair
		err = rows.Scan(&f.ID, &f.Name, &f.District, &f.Region5, &f.Neighborhood, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs.Debug("a fair %s \n", f.Name)
		fairs = append(fairs, &f)
	}

	logs.Debug("List MySQL %s len %d \n", "scan", len(fairs))

	if len(fairs) == 0 {
		errEmpty := errors.New("FAIR_TABLE: is Empty")
		return nil, errEmpty
	}

	logs.Debug("List MySQL %s \n", "OK")
	return fairs, nil
}

//Delete a Fair
func (r *FairMySQL) Delete(id entity.ID) error {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	stmt, err := db.Prepare("delete from air where id = ?")
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	numElement, err := result.RowsAffected()
	if numElement == 0 || err != nil {
		errEmpty := errors.New("FAIR_TABLE: Not Found elements to this Type - Delete")
		return errEmpty
	}

	tx.Commit()

	// _, err := tx.Exec("delete from fair where id = ?", id)
	// //_, err := r.db.Exec("delete from fair where id = ?", id)
	// if err != nil {
	// 	return err
	// }
	return nil
}

/*



const (
	//DBSourceOpenLocal Const used to Open Local db
	DBSourceOpenLocal = "root:my-secret-pw@tcp(localhost:3307)/"

	//DBSourceLocal Const used to acces Local db
	DBSourceLocal = "root:my-secret-pw@tcp(localhost:3307)/validatorAPP" //root:Mysql#my-secret-pw@/validatorAPP"

	//DBSourceOpenDocker Const used to Open Docker db
	DBSourceOpenDocker = "root:my-secret-pw@tcp(mysql-api)/" //mysql-api é o nome do serviço no docker-composer

	//DBSourceDocker Const used to acces Docker db
	DBSourceDocker = "root:my-secret-pw@tcp(mysql-api)/validatorAPP" //root:Mysql#my-secret-pw@/validatorAPP"

	//DBisDropTableSQL Clear TAbles
	DBisDropTableSQL = false
)

//AddrOpenDB VAR used to open and to access BD
var AddrOpenDB string

//AddrDB VAR data source name
var AddrDB string

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}


//InitDBSQL Connect To SQL Database
func InitDBSQL(isDropTable bool) {

	AddrOpenDB = DBSourceOpenDocker
	AddrDB = DBSourceDocker

	db, err := sql.Open("mysql", AddrOpenDB)
	if err != nil {
		log.Printf("Failed to connect to db Local Mysql...")
		AddrOpenDB = DBSourceOpenLocal
		AddrDB = DBSourceLocal
		db, err = sql.Open("mysql", AddrOpenDB)
		if err != nil {
			log.Printf("Failed to connect to db Local Mysql IP 127.0.0.1")
			log.Print(err)
		}

	}

	defer db.Close()

	fmt.Println("Successfully connected to the DB")
	exec(db, "create database if not exists validatorAPP")
	exec(db, "use validatorAPP")
	if isDropTable {
		exec(db, "drop table if exists querys")
	}

	exec(db, `create table IF NOT EXISTS querys(
	idx integer auto_increment,
	id varchar(40) ,
	number varchar(18),
	is_valid boolean,
	is_cpf boolean,
	is_cnpj boolean,
    createAt datetime,
	PRIMARY KEY (idx)
	)`)

	fmt.Println("Successfully connected to the DB!")

}


//CreateDB Create SQL dataBase
func CreateDB() {
	InitDBSQL(DBisDropTableSQL)
}
*/
