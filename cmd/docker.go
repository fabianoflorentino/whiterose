/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fabianoflorentino/whiterose/docker"
	"github.com/fabianoflorentino/whiterose/utils"
	"github.com/spf13/cobra"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Automates Docker operations, such as checking and building images.",
	Long: `The docker command in Whiterose automates common Docker tasks, such as 
checking for the existence of a Dockerfile in the current directory and building 
Docker images using environment variables and custom build arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case cmd.Flags().Changed("file"):
			isDockerFile()
		case cmd.Flags().Changed("build"):
			buildDockerImage()
		case cmd.Flags().Changed("delete"):
			deleteDockerImage()
		case len(args) == 0:
			if err := cmd.Help(); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.Flags().BoolP("file", "f", false, "Check if Dockerfile exists in the current directory")
	dockerCmd.Flags().BoolP("build", "b", false, "Build Docker image from Dockerfile")
	dockerCmd.Flags().BoolP("delete", "d", false, "Delete Docker image")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dockerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dockerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// isDockerFile checks if a Dockerfile exists in the current directory
func isDockerFile() {
	workDir := utils.GetEnvOrDefault("DOCKERFILE_PATH", os.Getenv("PWD"))
	d := docker.NewDockerManager(workDir)

	dockerfilePath, err := d.DetectDockerFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Dockerfile found at: %s\n", dockerfilePath[0])
}

// buildDockerImage builds a Docker image from the Dockerfile
func buildDockerImage() {
	workDir := utils.GetEnvOrDefault("DOCKERFILE_PATH", os.Getenv("PWD"))
	imageName := utils.GetEnvOrDefault("IMAGE_NAME", "my_app:latest")
	buildArgs := map[string]string{
		"IMAGE_VERSION": utils.GetEnvOrDefault("IMAGE_VERSION", "latest"),
	}

	d := docker.NewDockerManager(workDir)

	dockerfilePath, err := d.DetectDockerFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := d.BuildDockerImage(dockerfilePath[0], imageName, buildArgs); err != nil {
		fmt.Println(err)
		return
	}
}

func deleteDockerImage() {
	workDir := utils.GetEnvOrDefault("DOCKERFILE_PATH", os.Getenv("PWD"))
	imageName := utils.GetEnvOrDefault("IMAGE_NAME", "my_app:latest")

	d := docker.NewDockerManager(workDir)

	if err := d.DeleteDockerImage(imageName); err != nil {
		fmt.Println(err)
		return
	}
}
