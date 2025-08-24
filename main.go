package main

import (
	"log"

	"github.com/fabianoflorentino/whiterose/cmd"
	"github.com/fabianoflorentino/whiterose/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func init() {
	if err := utils.LoadDotEnv(); err != nil {
		log.Println("Failed to load .env file")
	}
}
