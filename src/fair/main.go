package main

import (
	"flag"
	"fmt"

	"github.com/marcovargas74/m74-fair-api/src/api/handler"
	"github.com/marcovargas74/m74-fair-api/src/config"
	logs "github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

//LogInFile write log to file
var LogInFile *bool

//GeneralConfigs  GENERAL Configuration
var GeneralConfigs config.ConfigAPI

func init() {
	LogInFile = flag.Bool("log", false, "write log to file")
	configs, err := config.ConfigGetAPIGeneral()
	if err != nil || GeneralConfigs.APIServerPortSQL == "" {
		logs.Error("Fail to Get Configurations-> %v ", GeneralConfigs)
	}

	GeneralConfigs = configs

	logs.Start(GeneralConfigs.IsProdType(), GeneralConfigs.APILogFile)

	err = config.LoadFromFileEnv(config.DEFAULT_ENV_FILE)
	if !GeneralConfigs.IsProdType() {
		err = config.LoadFromFileEnv("../../docker/.env")
	}

	if err != nil {
		logs.Error("Error[%s]loading .env file", err.Error())
	}
	//fmt.Printf("VARIAVEL port[%s] LIDO DO ARQUIVO .env", os.Getenv("SERVER_API_PORT_SQL"))

}

func main() {
	fmt.Printf("======== SAVE API FAIR Version %s TYPE:%v\n", config.VERSION_PACKAGE, GeneralConfigs.APITypeApp)
	logs.Debug("Get Configurations-> [%s] [%s] ", GeneralConfigs.APIServerPortSQL, GeneralConfigs.APIServerPortMem)

	go handler.StartAPI_MySQL(GeneralConfigs.APIServerPortSQL)
	go handler.StartAPI_Memory(GeneralConfigs.APIServerPortMem)

	var input string
	fmt.Scanln(&input)
	fmt.Println("APP FINALIZADO!")
	logs.LoggerClose()
}
