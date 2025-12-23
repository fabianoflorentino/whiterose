/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fabianoflorentino/whiterose/setup"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("all"):
			setup.PreReq()
			setup.GitCloneRepository()
		case cmd.Flags().Changed("pre-req"):
			setup.PreReq()
		case cmd.Flags().Changed("repos"):
			setup.GitCloneRepository()
		default:
			if err := cmd.Help(); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")
	setupCmd.PersistentFlags().BoolP("all", "a", false, "Check and install pre-requisites")
	setupCmd.PersistentFlags().BoolP("pre-req", "p", false, "Check and install pre-requisites")
	setupCmd.PersistentFlags().BoolP("repos", "r", false, "Clone Git repository")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
