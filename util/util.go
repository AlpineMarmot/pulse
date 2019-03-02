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

func GetCurrentPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
