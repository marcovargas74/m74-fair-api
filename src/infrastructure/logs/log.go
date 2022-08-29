package logs

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

//AppFileLog File where the logs will be saved
var AppFileLog *os.File

//Logger used to handle the log
var Logger zerolog.Logger

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	return original[i+1:]
}

//ThisFunction return a name of function can be used in Log
func ThisFunction() string {
	function, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s(%d)", chopPath(runtime.FuncForPC(function).Name()), line)
}

//LoggerClose Finish the Logger
func LoggerClose() {
	AppFileLog.Close()
}

//Start Start the Logger
func Start(isProg bool, filepath string) {

	if filepath == "" {
		filepath = "./myApp.log"
	}
	AppFileLog, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Error("Err [%s] Could not Open log FILE", err.Error())
	}

	if !isProg {
		Logger = zerolog.New(os.Stdout)
		return
	}
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
