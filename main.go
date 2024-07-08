package main

import (
	"flag"
	"fmt"
)

func main() {
	config := parseFlags()
	flag.Parse()

	fmt.Println("Starting server with config", *config)

	server := NewServer(config)

	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
