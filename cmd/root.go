/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whiterose",
	Short: "A CLI tool to automate cloning and setup of multiple Git repositories.",
	Long: `Whiterose is a command-line tool for automating the cloning and setup of multiple Git repositories.
It streamlines the process of preparing development environments, especially for teams working with several repositories.

Features:
- Clone repositories using HTTPS or SSH
- Automatically checkout the development branch if available
- Create and checkout a user-specific branch if development does not exist
- Load environment variables from a .env file
- Configure repositories via a JSON file

Example usage:
  whiterose setup
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}
