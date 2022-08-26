package repository

import (
	"database/sql"
	"errors"
	"fmt"
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

//TODO trazer a criacao do banco pra ca
const (
	//DBSourceOpenLocal Const used to Open Local db
	DBSourceOpenLocal = "root:my-secret-pw@tcp(localhost:3307)/"

	//DBSourceLocal Const used to acces Local db
	DBSourceLocal = "root:my-secret-pw@tcp(localhost:3307)/fairAPI?parseTime=true"

	//DBSourceOpenDocker Const used to Open Docker db
	//DBSourceOpenDocker = "root:my-secret-pw@tcp(mysql-api)/" //mysql-api é o nome do serviço no docker-composer

	//DBSourceDocker Const used to acces Docker db
	//DBSourceDocker = "root:my-secret-pw@tcp(mysql-api)/fairAPI"

)

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
	}
	return result
}

func CreateDB() {
	//--------------SE CONECTA AO BANCO
	//dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	//db, err := sql.Open("mysql", dataSourceName)

	//AddrOpenDB = DBSourceOpenDocker
	//AddrDB = DBSourceDocker
	AddrOpenDB = DBSourceOpenLocal
	AddrDB = DBSourceLocal

	db, err := sql.Open("mysql", AddrOpenDB)
	if err != nil {
		logs.Error("FALHA ao conectar ao Banco Mysql do DOcker %v", err)
		/*AddrOpenDB = DBSourceOpenLocal
		AddrDB = DBSourceLocal
		db, err = sql.Open("mysql", AddrOpenDB)
		if err != nil {
			logs.Error("FALHA ao conectar ao Banco Mysql Local IP 127.0.0.1 %v", err)
		}*/
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	time.Sleep(5 * time.Second)
	//--------------------------------------------------------------
	fmt.Println("Conectado ao Banco")
	exec(db, "create database if not exists fairAPI")
	exec(db, "use fairAPI")
	//if isDropTable {
	//exec(db, "drop table if exists fair")
	//}

	exec(db, `create table IF NOT EXISTS fair(
		   	idx integer auto_increment,
		   	id varchar(50) ,
		   	name varchar(50),
		   	district varchar(18),
		   	region5 varchar(6),
		   	neighborhood varchar(20),
			created_at datetime,
		   	updated_at datetime,
		   	PRIMARY KEY (idx)
		   	)`)

	//var db *sql.DB
	//repository.CreateDB()
	//db, err := sql.Open("mysql", repository.AddrOpenDB)

	/*db.Close()
	db, err = sql.Open("mysql", repository.AddrDB)
	if err != nil {
		log.Println(err)
		//return nil, err
	}*/

}

/*
func startDB() (*sql.Tx, error) {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err

	}

	return tx, nil
}*/

//OpenMysql create new repository
func OpenMysql() *sql.DB {

	db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

//NewFairMySQL create new repository
/*func NewFairMySQL() *FairMySQL {
	return &FairMySQL{
		db: openMysql(),
	}
}*/

//NewFairMySQL create new repository
func NewFairMySQL(db *sql.DB) *FairMySQL {
	//repository.CreateDB()

	return &FairMySQL{
		db: db,
	}
}

func (r *FairMySQL) beginMysql() (*sql.Tx, error) {

	/*db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()*/

	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err

	}

	return tx, nil
}

//Create a Fair TODO tratar erros
func (r *FairMySQL) Create(e *entity.Fair) (entity.ID, error) {

	//TODO REFATORAR - separae em uma funcao
	tx, err := r.beginMysql()
	if err != nil {
		log.Println(err)
		return e.ID, err
	}

	logs.Debug("MySQL CREATE BD %s \n", AddrDB)

	// stmt, err := r.db.Prepare(`
	// 	insert into fair (id, name, district, region5, neighborhood, created_at)
	// 	values(?,?,?,?,?,?)`)

	stmt, err := tx.Prepare("insert into fair (id, name, district, region5, neighborhood, created_at)values(?,?,?,?,?,?)")
	//stmt, err := r.db.Prepare("insert into fair (id, name, district, region5, neighborhood, created_at)values(?,?,?,?,?,?)")
	if err != nil {
		logs.Error("Err [%s] Insert DATA in Mysql ", err.Error())
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

	tx, err := r.beginMysql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	logs.Debug("getFair(2.MYSQL)ID %s", id)

	stmt, err := tx.Prepare("select id, name, district, region5, neighborhood, created_at from fair where id = ?")
	//stmt, err := r.db.Prepare(`id, name, district, region5, neighborhood, created_at from fair where id = ?`)
	if err != nil {
		logs.Error("Err [%s] Prepare DATA in Mysql ", err.Error())
		return nil, err
	}

	logs.Debug("MySQL 00 GET BD %v \n", id)

	var fair entity.Fair
	rows, err := stmt.Query(id)
	if err != nil {
		logs.Debug("MySQL  11 GET BD %v \n", id)
		return nil, err
	}

	logs.Debug("MySQL GET 22 BD %v \n", id)

	for rows.Next() {
		err = rows.Scan(&fair.ID, &fair.Name, &fair.District, &fair.Region5, &fair.Neighborhood, &fair.CreatedAt)
	}

	logs.Debug("MySQL 333 GET BD %v \n", id)

	return &fair, err
}

//Update a Fair
func (r *FairMySQL) Update(e *entity.Fair) error {
	e.UpdatedAt = time.Now()

	tx, err := r.beginMysql()
	if err != nil {
		log.Println(err)
		return err
	}

	//_, err := r.db.Exec("update fair set name = ?, district = ?, region5 = ?, neighborhood = ?, updated_at = ? where id = ?", e.Name, e.District, e.Region5, e.Neighborhood, e.UpdatedAt.Format("2006-01-02"), e.ID)
	stmt, err := tx.Prepare("update fair set name = ?, district = ?, region5 = ?, neighborhood = ?, updated_at = ? where id = ?")
	if err != nil {
		log.Println(err)
		return err
	}

	logs.Debug("MySQL UPDATE 22 BD %v \n", e.ID)
	_, err = stmt.Exec(e.Name, e.District, e.Region5, e.Neighborhood, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		logs.Error("Err [%s] FAilt To Update DATA in Mysql  [%s] ", err.Error(), e.Region5)
		return err
	}

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	tx.Commit()
	logs.Debug("MySQL UPDATE 2233 BD %v \n", e.ID)

	return nil

}

//Search fairs
func (r *FairMySQL) Search(key string, value string) ([]*entity.Fair, error) {

	logs.Debug("Search MySQL key[%s]value[%s] GET BD \n", key, value)

	/*tx, err := r.beginMysql()
	if err != nil {
		log.Println(err)
		return nil, err
	}*/

	busca := fmt.Sprintf("select id, name, district, region5, neighborhood, created_at from fair where %s like ?", key)
	stmt, err := r.db.Prepare(busca)

	if err != nil {
		logs.Error("Err [%s] Prepare DATA in Mysql ", err.Error())
		return nil, err
	}
	var fairs []*entity.Fair
	rows, err := stmt.Query("%" + value + "%")

	if err != nil {
		logs.Error("Err [%s] QUERY DATA in Mysql ", err.Error())
		return nil, err
	}
	for rows.Next() {
		var f entity.Fair
		err = rows.Scan(&f.ID, &f.Name, &f.District, &f.Region5, &f.Neighborhood, &f.CreatedAt)
		if err != nil {
			logs.Error("Search key[%s]value[%s] 333 GET BD \n", key, value)
			return nil, err
		}
		fairs = append(fairs, &f)
		logs.Debug("....rows.Next() key[%s]value[%s]\n", f.Name, f.Region5)
	}

	if len(fairs) == 0 {
		errEmpty := errors.New("FAIR_TABLE: is Empty")
		return nil, errEmpty
	}

	return fairs, nil
}

//List Fairs
func (r *FairMySQL) List() ([]*entity.Fair, error) {

	/*db, err := sql.Open("mysql", AddrDB)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()*/

	//logs.Debug("MySQL LIST BD %s \n", AddrDB)

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Println(err)
	// }

	logs.Debug("List MySQL %s \n", "init")

	//stmt, err := tx.Prepare(`select id, name, district, region5, neighborhood, created_at from fair`)
	//stmt, err := r.db.Prepare(`select id, name, district, region5, neighborhood, created_at from fair`)

	//OK rows, err := db.Query("select id, name, district, region5, neighborhood, created_at from fair")
	rows, err := r.db.Query("select id, name, district, region5, neighborhood, created_at from fair")
	if err != nil {
		logs.Error("ERRO: [%s] Nao conseguiu Acessar o DB \n", err.Error())
		return nil, err
	}
	defer rows.Close()

	var fairs []*entity.Fair
	// rows, err := stmt.Query()
	// if err != nil {
	// 	return nil, err
	// }

	logs.Debug("  List MySQL rows-> %v \n", rows)

	for rows.Next() {
		var fair entity.Fair
		err = rows.Scan(&fair.ID, &fair.Name, &fair.District, &fair.Region5, &fair.Neighborhood, &fair.CreatedAt)
		if err != nil {
			logs.Error("   ERRO!!!: [%s] Nao conseguiu Acessar a tabelaDB \n", err.Error())
			return nil, err
		}
		logs.Debug("a fair %s \n", fair.Name)
		fairs = append(fairs, &fair)
	}

	logs.Debug("List MySQL %s len %d \n", "scan", len(fairs))

	if len(fairs) == 0 {
		logs.Warn("WARNING: [%s]  %v \n", entity.ErrEmptyDB.Error(), len(fairs))
		return nil, entity.ErrEmptyDB
	}

	logs.Debug("List MySQL %s \n", "OK")
	return fairs, nil
}

//Delete a Fair
func (r *FairMySQL) Delete(id entity.ID) error {

	tx, err := r.beginMysql()
	if err != nil {
		log.Println(err)
		return err
	}

	stmt, err := tx.Prepare("delete from fair where id = ?")
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

	return nil
}
