package main

import (
	"fmt"
	"os"
)

func main() {
	var config Config
	err := ReadConfig(&config)
	if err != nil {
		processError(err)
	}

	fmt.Printf("%+v", config)
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
