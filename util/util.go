package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckError(err error, context string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("Context: %s.", context))
		fmt.Println(err)
	}
}

func FatalError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, ":")
		fmt.Println("\t", err.Error())
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
