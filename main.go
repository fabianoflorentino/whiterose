/*
Copyright © 2025 FABIANO SANTOS FLORENTINO <FABIANORATM@GMAIL.COM>
*/
package main

import (
	"github.com/fabianoflorentino/whiterose/cmd"
	"github.com/joho/godotenv"
)

func main() {
	cmd.Execute()
}

func init() {
	godotenv.Load()
}
