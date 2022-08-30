package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

const (
	VERSION_PACKAGE = "2022-08-30"

	DEFAULT_DB_USER     = "root"
	DEFAULT_DB_PASSWORD = "my-secret-pw"
	DEFAULT_DB_DATABASE = "fairAPI"
	DEFAULT_DB_ADDRESS  = "localhost"
	DEFAULT_DB_PORT     = "3307"
	DEFAULT_URL_MYSQL   = "root:my-secret-pw@tcp(localhost:3307)/fairAPI?parseTime=true"

	DEFAULT_SERVER_API_PORT_MEM = ":5000"
	DEFAULT_SERVER_API_PORT_SQL = ":5001"

	DEFAULT_LOG_FILE = "./fairAPI.log"
	DEFAULT_ENV_FILE = "./.env"
	DEFAULT_CSV_FILE = "defaultValues.csv"
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

//Environs All environment variables have been standardized with underLine at the beginning
//Whenever you create a configuration variable, you must include its corresponding default in the environment.
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

//ConfigAPI Is the configuration variables structure of the entire system
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

//NewConfigAPIDefault Returns the default api settings
func NewConfigAPIDefault() ConfigAPI {
	return ConfigAPI{
		APIServerPortMem: DEFAULT_SERVER_API_PORT_MEM,
		APIServerPortSQL: DEFAULT_SERVER_API_PORT_MEM,
		APITypeApp:       TYPE_PROD,
		APILogFile:       DEFAULT_LOG_FILE,
	}
}

func getEnv(key string) string {
	value, exist := os.LookupEnv(key)
	logs.Debug("  Getenv()..key[%s] value[%s]exist[%t]\n", key, value, exist)
	if exist && value != "" {
		return value
	}
	return Environs[key]
}

//CreateNewConfigAPI Creates the application's initial settings by searching for environment variables
func CreateNewConfigAPI() (ConfigAPI, error) {
	config := ConfigAPI{
		MYSQLUser:     getEnv(_DB_USER),
		MYSQLPassword: getEnv(_DB_PASSWORD),
		MYSQLDatabase: getEnv(_DB_DATABASE),
		MYSQLAddress:  getEnv(_DB_ADDRESS),
		MYSQLPortTCP:  getEnv(_DB_PORT),

		APIServerPortMem: getEnv(_SERVER_API_PORT_MEM),
		APIServerPortSQL: getEnv(_SERVER_API_PORT_SQL),
		APITypeApp:       getEnv(_TYPE_APP),
		APILogFile:       getEnv(_LOG_FILE),
	}
	err := config.Validate()
	if err != nil {
		return config, entity.ErrInvalidConfig
	}
	return config, nil
}

//Validate validate Configs Vars
func (c *ConfigAPI) Validate() error {
	if c.APITypeApp == "" || c.APIServerPortMem == "" || c.APIServerPortSQL == "" || c.APILogFile == "" {
		return entity.ErrInvalidConfig
	}

	if c.MYSQLUser == "" || c.MYSQLPassword == "" || c.MYSQLDatabase == "" || c.MYSQLAddress == "" || c.MYSQLPortTCP == "" {
		return entity.ErrInvalidConfig
	}
	return nil
}

//NewConfigAPI Load All Configurations Vars used in API
func NewConfigAPI() (ConfigAPI, error) {
	config, err := CreateNewConfigAPI()
	if err != nil {
		logs.Error("%s Fail to Create NewConfigAPI()-> %v ", logs.ThisFunction(), err.Error())
		return NewConfigAPIDefault(), entity.ErrDefaultConfig

	}

	return config, nil
}

//ConfigGetMysqlURL retorna o endereço URL do banco de dados
func ConfigGetMysqlURL() (string, error) {

	mySQLConfig, err := NewConfigAPI()
	if err != nil && err != entity.ErrDefaultConfig {
		logs.Error("%s Fail to Get MySQL Configurations->[%v] ", logs.ThisFunction(), err.Error())
		return DEFAULT_URL_MYSQL, err
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mySQLConfig.MYSQLUser, mySQLConfig.MYSQLPassword, mySQLConfig.MYSQLAddress, mySQLConfig.MYSQLPortTCP, mySQLConfig.MYSQLDatabase)

	return dataSourceName, err
}

// DataBaseName Returns name of the used database
func DataBaseName() string {
	dataBase, err := ConfigGetMysqlURL()
	if err != nil {
		logs.Warn("Using DEFAULT DB err:[%v] ", err.Error())
	}

	i := strings.LastIndex(dataBase, "/")
	if i == -1 {
		return dataBase
	}

	return dataBase[i+1:]

}

// DataBaseURL Returns the URL of the Database used. Required to Open DB
func DataBaseURL() string {
	dataBase, err := ConfigGetMysqlURL()
	if err != nil {
		logs.Warn("Using DEFAULT DB err:[%v] ", err.Error())
	}

	i := strings.LastIndex(dataBase, "/")
	if i == -1 {
		return dataBase
	}
	return dataBase[:i+1]

}

// DataBaseURLToLog Returns the URL of the Database used. Used in logger. Don't show the user and password
func DataBaseURLLog() string {
	dataBase := DataBaseURL()
	i := strings.LastIndex(dataBase, "@")
	if i == -1 {
		return dataBase
	}
	return dataBase[i+1:]

}

//ConfigGetAPIGeneral Return GENERAL Configuration used in API
func ConfigGetAPIGeneral() (ConfigAPI, error) {

	myConfigs, err := NewConfigAPI()

	if err != nil && err != entity.ErrDefaultConfig {
		logs.Error("Fail to Get APIs GENERAL Configurations-> %v ", err.Error())
		return myConfigs, err
	}

	logs.Debug("   Get APIs GENERAL Configurations-> %v ", myConfigs)
	return myConfigs, nil
}

// IsProdType Check if Prodution Version
func (c *ConfigAPI) IsProdType() bool {
	return strings.Contains(c.APITypeApp, TYPE_PROD)
}

//SetModTest Force use test mod
func (c *ConfigAPI) SetModTest() {
	c.APITypeApp = TYPE_TEST
}

func cpyFilesFromDocker() {
	cmd := exec.Command("cp", "-prf", "../../docker/defaultValues.csv", DEFAULT_CSV_FILE)
	cmd.Run()

	cmd = exec.Command("cp", "-prf", "../../docker/.env", ".env")
	cmd.Run()
}

//LoadFromFileEnv sets the value of the environment variable from a file (.env)
func LoadFromFileEnv(file string) error {

	cpyFilesFromDocker()
	if file == "" {
		file = ".env"
	}
	return godotenv.Load(file)
}
