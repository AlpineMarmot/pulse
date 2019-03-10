package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type logger struct {
	logger      *log.Logger
	fileHandler *os.File
	file        interface{}
	created     bool
	quiet       bool
}

var loggerInstance logger

func New(file interface{}, quiet bool) {
	if loggerInstance.created == true {
		return
	}
	loggerInstance = logger{
		file:    file,
		created: true,
		quiet:   quiet,
	}

	if file != nil {
		createLoggerWithFile(file.(string))
	} else {
		createLogger()
	}
}

func createLoggerWithFile(file string) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Opening/Creating file " + file + " as log output fail")
	}
	multi := io.MultiWriter(f, os.Stdout)
	loggerInstance.logger = log.New(multi, "", log.Ldate|log.Ltime)
}

func createLogger() {
	loggerInstance.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func CloseLogFile() {
	if loggerInstance.file != nil {
		loggerInstance.fileHandler.Close()
	}
}

func Println(v ...interface{}) {
	if loggerInstance.quiet == false {
		loggerInstance.logger.Println(v)
	}
}

func Print(v ...interface{}) {
	if loggerInstance.quiet == false {
		loggerInstance.logger.Print(v)
	}
}

func Printf(format string, v ...interface{}) {
	if loggerInstance.quiet == false {
		loggerInstance.logger.Printf(format, v)
	}
}
