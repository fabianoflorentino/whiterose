/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fabianoflorentino/whiterose/prereq"
	"github.com/spf13/cobra"
)

// preReqCmd represents the preReq command
var preReqCmd = &cobra.Command{
	Use:   "pre-req",
	Short: "Validate and list required applications for the environment.",
	Long: `The pre-req command helps you manage environment prerequisites by listing
all required applications or validating the presence of specific ones.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := prereq.NewAppValidator()

		// validApps receives the list of applications to validate
		validApps, _ := cmd.Flags().GetStringSlice("apps")

		switch {
		case cmd.Flags().Changed("check"):
			app.ValidateApps()
		case cmd.Flags().Changed("list"):
			app.ListAvailableApps()
		case cmd.Flags().Changed("apps"):
			app.ValidateSpecificApps(validApps)
		case len(args) == 0:
			cmd.Help()
		default:
		}
	},
}

func init() {
	rootCmd.AddCommand(preReqCmd)

	preReqCmd.Flags().BoolP("check", "c", false, "Check if all required applications are installed")
	preReqCmd.Flags().BoolP("list", "l", false, "List all available applications")
	preReqCmd.Flags().StringSliceP("apps", "a", []string{}, "Validate specific applications (comma-separated)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// preReqCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// preReqCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
