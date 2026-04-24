package main

import (
	"log"

	"github.com/fabianoflorentino/whiterose/cmd"
	"github.com/fabianoflorentino/whiterose/config"
)

func main() {
	cfg := config.LoadOrDefault()
	log.Printf("Using config: git.base=%s, repo.path=%s", cfg.Git.Base, cfg.Repo.Path)

	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}