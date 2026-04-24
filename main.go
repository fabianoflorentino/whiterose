package main

import (
	"log"
	"os"

	"github.com/fabianoflorentino/whiterose/cmd"
	"github.com/fabianoflorentino/whiterose/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func init() {
	if os.Getenv("SKIP_DOTENV") == "true" {
		return
	}

	if err := utils.LoadDotEnv(); err != nil {
		if os.Getenv("CI") == "true" {
			log.Println("CI mode: skipping .env file")
			return
		}
		log.Printf("Warning: Failed to load .env file: %v", err)
	}
}