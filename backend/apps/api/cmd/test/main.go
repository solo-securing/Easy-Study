package main

import (
	"api/bootstrap"
	"flag"
	"fmt"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "./config/config.yaml", "Path to the configuration file")
	flag.Parse()

	// Initialize the application
	app := bootstrap.App(*configPath)

	// Print the loaded configuration for verification
	fmt.Println(app.Config)
}
