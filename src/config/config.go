package config

import (
	"fmt"
	"os"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

const (
	VERSION_PACKAGE = "2022-08-21"
	//DEFAULT VALUES
	DEFAULT_DB_USER     = "root"
	DEFAULT_DB_PASSWORD = "my-secret-pw"
	DEFAULT_DB_DATABASE = "fairAPI"
	DEFAULT_DB_ADDRESS  = "localhost"
	DEFAULT_DB_PORT     = "3307"
	DEFAULT_URL_MYSQL   = "root:my-secret-pw@tcp(localhost:3307)/fairAPI?parseTime=true"

	DEFAULT_SERVER_API_PORT_MEM = ":5000"
	DEFAULT_SERVER_API_PORT_SQL = ":5001"
	DEFAULT_DEV_TEST            = true
	TYPE_PROD                   = "PROD"
	TYPE_DEV                    = "DEV"
	TYPE_TEST                   = "TEST"
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
)

//Environs Foi padronizado todas as variaveis de ambiente com underLine no inicio
//Sempre que criar uma variavel de configuracao deve incluir nos ambiente a susa correspondente default
var Environs = map[string]string{

	_SERVER_API_PORT_MEM: DEFAULT_SERVER_API_PORT_MEM,
	_SERVER_API_PORT_SQL: DEFAULT_SERVER_API_PORT_SQL,
	_TYPE_APP:            TYPE_PROD,

	_DB_USER:     DEFAULT_DB_USER,
	_DB_PASSWORD: DEFAULT_DB_PASSWORD,
	_DB_DATABASE: DEFAULT_DB_DATABASE,
	_DB_ADDRESS:  DEFAULT_DB_ADDRESS,
	_DB_PORT:     DEFAULT_DB_PORT,
}

/*
const (
	VAR_DB_USER           = "DB_USER"
	VAR_DB_PASSWORD       = "my-secret-pw"
	VAR_DB_MYSQL_DATABASE = "fairAPI"
	VAR_DB_HOST           = "127.0.0.1"
	VAR_DB_URL            = "root:my-secret-pw@tcp(localhost:3307)"
	VAR_DB_PORT           = "3306"
	VAR_SERVER_API_PORT   = ":5000"
	VAR_IS_PROD           = "PROD"
	VAR_DEV               = "DEV"
	VAR_DEV_TEST          = true
)*/

//ServerAPI is struct to start server
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
	//API_PORT              string = 8080
}

func NewConfigAPIDefault() ConfigAPI {
	return ConfigAPI{
		APIServerPortMem: ":5000",
		APIServerPortSQL: ":5001",
		APITypeApp:       "PROD",
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
	}
	err := config.Validate()
	if err != nil {
		return config, entity.ErrInvalidConfig
	}
	return config, nil
}

//Validate validate book
func (c *ConfigAPI) Validate() error {
	if c.APITypeApp == "" || c.APIServerPortMem == "" || c.APIServerPortSQL == "" {
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

//ConfigGetAPIGeneral Return GENERAL Configuration used in API
func ConfigGetAPIGeneral() (ConfigAPI, error) {

	myConfigs, err := NewConfigAPI()

	fmt.Printf("  ConfigGetAPIGeneralqqq()..isProd[%s]portM[%s]portDB[%s]\n", myConfigs.APITypeApp, myConfigs.APIServerPortMem, myConfigs.APIServerPortSQL)

	if err != nil && err != entity.ErrDefaultConfig {
		logs.Error("Fail to Get APIs GENERAL Configurations-> %v ", err.Error())
		return myConfigs, err
	}

	logs.Debug("   oooookkGet APIs GENERAL Configurations-> %v ", myConfigs)
	return myConfigs, nil
}

//return "", errors.New("Invalid query")
/*


/*
func ConfigDB(string) {

	//Cria valor padrao..
	url = bancodedados

	//Le variavel de ambientefmt
	//se existir usa
	//pega a confuguracao do banco se nao tiver usa a padrao
	envieValueOfEnvironVar("_DB_URL")

	//monta url

}* /


*/
