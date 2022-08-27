package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

const (
	VERSION_PACKAGE = "2022-08-27"
	//DEFAULT VALUES
	DEFAULT_DB_USER     = "root"
	DEFAULT_DB_PASSWORD = "my-secret-pw"
	DEFAULT_DB_DATABASE = "fairAPI"
	DEFAULT_DB_ADDRESS  = "localhost"
	DEFAULT_DB_PORT     = "3307"
	DEFAULT_URL_MYSQL   = "root:my-secret-pw@tcp(localhost:3307)/fairAPI?parseTime=true"

	DEFAULT_SERVER_API_PORT_MEM = ":5000"
	DEFAULT_SERVER_API_PORT_SQL = ":5001"

	DEFAULT_LOG_FILE = "./fairAPI.log"
	TYPE_PROD        = "PROD"
	TYPE_DEV         = "DEV"
	TYPE_TEST        = "TEST"
)

const (
	_DB_USER     = "DB_USER"
	_DB_PASSWORD = "DB_PASSWORD"
	_DB_DATABASE = "DB_DATABASE"
	_DB_ADDRESS  = "DB_IP"
	_DB_PORT     = "DB_PORT"

	_SERVER_API_PORT_MEM = "SERVER_API_PORT_MEM"
	_SERVER_API_PORT_SQL = "SERVER_API_PORT_SQL"
	_TYPE_APP            = "TYPE_APP"
	_LOG_FILE            = "LOG_FILE"
)

//Environs Foi padronizado todas as variaveis de ambiente com underLine no inicio
//Sempre que criar uma variavel de configuracao deve incluir nos ambiente a susa correspondente default
var Environs = map[string]string{

	_SERVER_API_PORT_MEM: DEFAULT_SERVER_API_PORT_MEM,
	_SERVER_API_PORT_SQL: DEFAULT_SERVER_API_PORT_SQL,
	_TYPE_APP:            TYPE_PROD,
	_LOG_FILE:            DEFAULT_LOG_FILE,

	_DB_USER:     DEFAULT_DB_USER,
	_DB_PASSWORD: DEFAULT_DB_PASSWORD,
	_DB_DATABASE: DEFAULT_DB_DATABASE,
	_DB_ADDRESS:  DEFAULT_DB_ADDRESS,
	_DB_PORT:     DEFAULT_DB_PORT,
}

//ConfigAPI è a estrutura de Variaveis de configuracao de todo o sistema
type ConfigAPI struct {
	//DB
	MYSQLUser     string `default:"root"`
	MYSQLPassword string `default:"my-secret-pw"`
	MYSQLDatabase string `default:"fairAPI"`
	MYSQLAddress  string `default:"localhost"`
	MYSQLPortTCP  string `default:"3307"`

	//API
	APIServerPortMem string `default:":5000"`
	APIServerPortSQL string `default:":5001"`
	APITypeApp       string `default:"PROD"`
	APILogFile       string `default:"./fairAPI.log"`
}

func NewConfigAPIDefault() ConfigAPI {
	return ConfigAPI{
		APIServerPortMem: DEFAULT_SERVER_API_PORT_MEM,
		APIServerPortSQL: DEFAULT_SERVER_API_PORT_MEM,
		APITypeApp:       TYPE_PROD,
		APILogFile:       DEFAULT_LOG_FILE,
	}
}

func Getenv(key string) string {
	value, exist := os.LookupEnv(key)
	//fmt.Printf("  Getenv()..key[%s] value[%s]exist[%t]\n", key, value, exist)
	if exist {
		return value
	}
	return Environs[key]
}

func CreateNewConfigAPI() (ConfigAPI, error) {
	config := ConfigAPI{
		MYSQLUser:     Getenv(_DB_USER),
		MYSQLPassword: Getenv(_DB_PASSWORD),
		MYSQLDatabase: Getenv(_DB_DATABASE),
		MYSQLAddress:  Getenv(_DB_ADDRESS),
		MYSQLPortTCP:  Getenv(_DB_PORT),

		APIServerPortMem: Getenv(_SERVER_API_PORT_MEM),
		APIServerPortSQL: Getenv(_SERVER_API_PORT_SQL),
		APITypeApp:       Getenv(_TYPE_APP),
		APILogFile:       Getenv(_LOG_FILE),
	}
	err := config.Validate()
	if err != nil {
		return config, entity.ErrInvalidConfig
	}
	return config, nil
}

//Validate validate book
func (c *ConfigAPI) Validate() error {
	if c.APITypeApp == "" || c.APIServerPortMem == "" || c.APIServerPortSQL == "" || c.APILogFile == "" {
		return entity.ErrInvalidConfig
	}

	if c.MYSQLUser == "" || c.MYSQLPassword == "" || c.MYSQLDatabase == "" || c.MYSQLAddress == "" || c.MYSQLPortTCP == "" {
		return entity.ErrInvalidConfig
	}
	return nil
}

func NewConfigAPI() (ConfigAPI, error) {
	config, err := CreateNewConfigAPI()
	if err != nil {
		logs.Error("Fail to Create NewConfigAPI()-> %v ", err.Error())
		return NewConfigAPIDefault(), entity.ErrDefaultConfig

	}

	return config, nil
}

func ConfigGetMysqlURL() (string, error) {

	mySQLConfig, err := NewConfigAPI()
	if err != nil && err != entity.ErrDefaultConfig {
		logs.Error("Fail to Get MySQL Configurations-> %v ", err.Error())
		return DEFAULT_URL_MYSQL, err
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mySQLConfig.MYSQLUser, mySQLConfig.MYSQLPassword, mySQLConfig.MYSQLAddress, mySQLConfig.MYSQLPortTCP, mySQLConfig.MYSQLDatabase)

	return dataSourceName, err
}

// DataBaseName Retorna nome da Base de dados usada
func DataBaseName() string {
	dataBase, err := ConfigGetMysqlURL()
	if err != nil {
		logs.Warn("Usando Banco DEFAULT err:[%v] ", err.Error())
	}

	i := strings.LastIndex(dataBase, "/")
	if i == -1 {
		return dataBase
	}
	return dataBase[i+1:]

}

// DataBaseURL Retorna a URL da Base de dados usada. Necessario para Abrir o DB
func DataBaseURL() string {
	dataBase, err := ConfigGetMysqlURL()
	if err != nil {
		logs.Warn("Usando Banco DEFAULT err:[%v] ", err.Error())
	}

	i := strings.LastIndex(dataBase, "/")
	if i == -1 {
		return dataBase
	}
	return dataBase[:i+1]

}

//ConfigGetAPIGeneral Return GENERAL Configuration used in API
func ConfigGetAPIGeneral() (ConfigAPI, error) {

	myConfigs, err := NewConfigAPI()

	//fmt.Printf("  ConfigGetAPIGeneralqqq()..isProd[%s]portM[%s]portDB[%s]\n", myConfigs.APITypeApp, myConfigs.APIServerPortMem, myConfigs.APIServerPortSQL)
	if err != nil && err != entity.ErrDefaultConfig {
		logs.Error("Fail to Get APIs GENERAL Configurations-> %v ", err.Error())
		return myConfigs, err
	}

	logs.Debug("   Get APIs GENERAL Configurations-> %v ", myConfigs)
	return myConfigs, nil
}
