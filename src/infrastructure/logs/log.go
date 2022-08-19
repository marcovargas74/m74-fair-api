package logs

import (
	"log"
	"os"

	"github.com/rs/zerolog"
)

//AppLogProg TRUE Log Programado
var AppLogProg bool

//AppFileLog Variavel usado no syslog
var AppFileLog *os.File

//Logger Variavel usado no Log
var Logger zerolog.Logger

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
	//AppSyslog(syslog.LOG_INFO, "%s {LOG_FINISH}\n", ThisFunction())
	AppFileLog.Close()
}

func Start() {
	AppFileLog, err := os.OpenFile("./zero.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	//defer file.Close()
	logger := zerolog.New(AppFileLog).With().Timestamp().Logger()
	logger.Debug().Msg("Log STARTED!")
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
