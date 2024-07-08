package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := parseFlags()
	flag.Parse()

	fmt.Println("Starting server with config", *config)

	server := NewServer(config)
	//server.static("./public") can be used to serve static files as alternative to flags

	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	server.logger.Printf("Received signal %v, initiating shutdown...\n", sig)

	if err := server.Stop(); err != nil {
		server.logger.Printf("Server failed to stop gracefully: %v\n", err)
	} else {
		server.logger.Println("Server stopped gracefully")
	}
}
