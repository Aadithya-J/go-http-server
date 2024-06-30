package main

import (
	"flag"
	"fmt"
)

func main() {
	config := parseFlags()
	flag.Parse()

	fmt.Println("Starting server with config", config)
}
