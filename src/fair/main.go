package main

import (
	"flag"

	fair "github.com/marcovargas74/m74-val-cpf-cnpj/src/api-validator"
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/api/handler"
	logs "github.com/marcovargas74/m74-val-cpf-cnpj/src/infrastructure/logs"
)

func init() {
	//myquery.CreateStatus()
	//myquery.SetUsingMongoDocker(myquery.SetDockerRun)
	//myquery.CreateDB()
	prod := flag.Bool("prod", false, "write log to file")
	logs.Start(*prod, "./fairAPI.log")

}

func main() {
	//logs.Start()
	//prod := flag.Bool("prod", false, "write log to file")
	//debug := flag.Bool("debug", false, "enable debug mode")

	logs.Debug("======== API FAIR Version %s \n", fair.Version())
	handler.StartAPI(true)
	logs.LoggerClose()
}
