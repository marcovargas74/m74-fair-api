package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

//AppLogProg TRUE Log Programado
var AppLogProg bool

//AppFileLog Variavel usado no syslog
var AppFileLog *os.File

//Logger Variavel usado no Log
var Logger zerolog.Logger

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	return original[i+1:]
}

func ThisFunction() string {
	function, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s(%d)", chopPath(runtime.FuncForPC(function).Name()), line)
}

//StartLogger Inicia Login da aplicação
//isProg true to show message
func StartLogger(isProg bool, filepath string) {
	var err error

	AppLogProg = isProg

	/*if !isProg {
		return
	}*/

	if filepath == "" {
		filepath = "/tmp/my-app.log"
	}
	AppFileLog, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	//defer AppFileLog.Close()
	log.SetOutput(AppFileLog)
	// loga data e hora
	log.SetFlags(log.LstdFlags)
	log.Println("começou")
	//AppSyslog(syslog.LOG_INFO, "%s sys/Log  Iniciado\n", ThisFunction())
}

//LoggerClose Finish the Logger
func LoggerClose() {
	AppFileLog.Close()
}

func Start(isProg bool, filepath string) {

	if filepath == "" {
		filepath = "./myApp.log"
	}
	AppFileLog, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	//AppFileLog, err := os.OpenFile("./zero.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	//defer file.Close()
	if !isProg {
		Logger = zerolog.New(os.Stdout)
		return
	}

	//logger := zerolog.New(AppFileLog).With().Timestamp().Logger()

	Logger.Debug().Msg(" ------Log STARTED!-------")
	Logger = zerolog.New(AppFileLog).With().Timestamp().Logger()

}

func Info(format string, a ...interface{}) {
	Logger.Info().Msgf(format, a...)
}

func Warn(format string, a ...interface{}) {
	Logger.Warn().Msgf(format, a...)
}

func Error(format string, a ...interface{}) {
	Logger.Error().Msgf(format, a...)
}

func Debug(format string, a ...interface{}) {
	Logger.Debug().Msgf(format, a...)
}
