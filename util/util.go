package util

import (
	"fmt"
	"github.com/AlpineMarmot/pulse/logger"
	"os"
	"path/filepath"
)

func CheckError(err error, context string) {
	if err != nil {
		logger.Println(fmt.Sprintf("Context: %s.", context))
		logger.Println(err)
	}
}

func FatalError(msg string, err error) {
	if err != nil {
		logger.Print(msg, ": ")
		logger.Println(err.Error())
		os.Exit(1)
	}
}

func GetCurrentPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// check config file existence
func FileExists(configFilePath string) bool {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}
