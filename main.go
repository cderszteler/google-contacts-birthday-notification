package main

import (
	"fmt"
	"os"
)

var config Config

func main() {
	if err := ReadConfig(&config); err != nil {
		processError(err)
	}
	CreateService()

	if err := SendMail("test"); err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
