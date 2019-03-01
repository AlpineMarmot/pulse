package util

import "fmt"

func CheckError(err error, context string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("Context: %s.", context))
		fmt.Println(err)
	}
}
