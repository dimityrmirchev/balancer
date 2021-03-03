package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.Backends)
}
