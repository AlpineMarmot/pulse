package main

import "fmt"

func checkError(err error, context string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("Context: %s.", context))
		fmt.Println(err)
	}
}
