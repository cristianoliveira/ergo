package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var VERSION string

	version, err := ioutil.ReadFile(".version")

	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	VERSION = string(version)

	fmt.Print(VERSION)

}
