/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fabianoflorentino/whiterose/internal/domain/entities"
	"github.com/fabianoflorentino/whiterose/update"
	"github.com/spf13/cobra"
)

var (
	updateGoMod       bool
	updateGoVersion   bool
	updateDockerImage bool
	updatePackages    bool
	updateMajor       bool
	updateConfigPath  string
	updateList        bool
	updatePR          bool
	updateDryRun      bool
	updateBase        string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza dependências e versões de projetos",
	Long: `The 'update' command automatically updates dependencies and versions.

Supported updates:
- --go-mod: Update go.mod dependencies
- --go-version: Update Go version in go.mod
- --docker-image: Update base Docker image
- --packages: Update Go packages to latest versions

Update strategies:
- Minor/patch: Automatic (no confirmation)
- Major: Requires confirmation to avoid breaking changes

The command will:
1. Load projects from config file
2. Update specified components
3. Create a new branch
4. Commit changes
5. Push to origin for PR creation
`,
	Run: func(cmd *cobra.Command, args []string) {
		if updateList {
			runListVersions()
			return
		}

		if !updateGoMod && !updateGoVersion && !updateDockerImage && !updatePackages {
			fmt.Println("Error: specify at least one of --go-mod, --go-version, --docker-image, or --packages")
			os.Exit(1)
		}

		if updateMajor {
			if !confirmMajorUpdate() {
				fmt.Println("Update cancelled.")
				os.Exit(0)
			}
		}

		service := update.New()

		projects, err := service.LoadUpdateConfig(updateConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if len(projects) == 0 {
			fmt.Println("No projects found in config.")
			os.Exit(0)
		}

		for _, project := range projects {
			fmt.Printf("\n--- Updating %s ---\n", project.Name)

			var changes []string
			var err error

			if updateGoMod && project.GoMod != nil {
				parsedStrategy := entities.ParseUpdateStrategy(project.GoMod.UpdateStrategy.String())
				err = service.UpdateGoMod(project, parsedStrategy, updateMajor)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating go.mod: %v\n", err)
					continue
				}
				changes = append(changes, "Updated go.mod dependencies")
			}

			if updateGoVersion && project.GoVersion != nil {
				parsedStrategy := entities.ParseUpdateStrategy(project.GoVersion.UpdateStrategy.String())
				err = service.UpdateGoVersion(project, parsedStrategy, updateMajor)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating Go version: %v\n", err)
					continue
				}
				changes = append(changes, "Updated Go version")
			}

			if updateDockerImage && project.DockerImage != nil {
				parsedStrategy := entities.ParseUpdateStrategy(project.DockerImage.UpdateStrategy.String())
				err = service.UpdateDockerImage(project, parsedStrategy, updateMajor)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating Docker image: %v\n", err)
					continue
				}
				changes = append(changes, "Updated Docker base image")
			}

			if updatePackages && project.GoMod != nil {
				checker := update.NewVersionChecker()
				err = checker.UpdatePackages(project.Path, "patch", updateDryRun)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating packages: %v\n", err)
					continue
				}
				changes = append(changes, "Updated Go packages")
			}

			if len(changes) > 0 {
				if updateDryRun {
					fmt.Printf("[DRY RUN] Would create branch and commit: %v\n", changes)
					continue
				}

				branchName, err := service.CreateBranchAndCommit(project, changes)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating branch/commit: %v\n", err)
					continue
				}

				if updatePR {
					if err := service.CreatePRWithBase(project, branchName, changes, updateBase); err != nil {
						fmt.Fprintf(os.Stderr, "Error creating PR: %v\n", err)
						continue
					}
				}

				fmt.Printf("Update completed for %s\n", project.Name)
			}
		}
	},
}

func runListVersions() {
	checker := update.NewVersionChecker()

	if updateGoMod {
		projects, err := update.New().LoadUpdateConfig(updateConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		for _, project := range projects {
			if project.GoMod != nil {
				fmt.Printf("\n=== %s: Library updates ===\n", project.Name)
				if err := checker.ListGoLibUpdates(project.Path); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			}
		}
	}

	if updateGoVersion {
		if err := checker.ListGoVersions(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	if updateDockerImage {
		projects, err := update.New().LoadUpdateConfig(updateConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		for _, project := range projects {
			if project.DockerImage != nil {
				fmt.Printf("\n=== %s: Docker image updates ===\n", project.Name)
				if err := checker.ListDockerUpdates(project.DockerImage.Base); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			}
		}
	}

	if !updateGoMod && !updateGoVersion && !updateDockerImage {
		fmt.Println("Listing all available updates...")
		if err := checker.ListGoVersions(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

func confirmMajorUpdate() bool {
	fmt.Println("WARNING: Major updates may introduce breaking changes.")
	fmt.Print("Are you sure you want to continue? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "yes"
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&updateGoMod, "go-mod", "g", false, "Update/List go.mod dependencies")
	updateCmd.Flags().BoolVarP(&updateGoVersion, "go-version", "v", false, "Update/List Go version")
	updateCmd.Flags().BoolVarP(&updateDockerImage, "docker-image", "d", false, "Update/List Docker base image")
	updateCmd.Flags().BoolVarP(&updatePackages, "packages", "p", false, "Update Go packages to latest versions")
	updateCmd.Flags().BoolVarP(&updateMajor, "major", "m", false, "Update major version (requires confirmation)")
	updateCmd.Flags().BoolVarP(&updateList, "list", "l", false, "List available updates instead of updating")
	updateCmd.Flags().BoolVarP(&updatePR, "pr", "r", false, "Create pull request after pushing (requires gh cli)")
	updateCmd.Flags().BoolVarP(&updateDryRun, "dry-run", "n", false, "Show what would be updated without making changes")
	updateCmd.Flags().StringVarP(&updateConfigPath, "config", "c", "", "Path to update config file")
	updateCmd.Flags().StringVarP(&updateBase, "base", "b", "main", "Base branch for PR")
}
