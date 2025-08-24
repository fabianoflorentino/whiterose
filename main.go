/*
Copyright © 2025 FABIANO SANTOS FLORENTINO <FABIANORATM@GMAIL.COM>
*/
package main

import (
	"github.com/fabianoflorentino/whiterose/cmd"
	"github.com/joho/godotenv"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic("Failed to execute command")
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}
}
