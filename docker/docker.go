// DockerManager provides methods to manage Docker operations within a specified working directory.
// It supports detecting Dockerfiles and building Docker images using specified build arguments.
//
// NewDockerManager(workDir string) *DockerManager:
//
//	Creates a new DockerManager instance for the given working directory.
//
// DetectDockerFile() ([]string, error):
//
//	Recursively searches for Dockerfiles in the working directory.
//	Returns a slice of paths to found Dockerfiles or an error if none are found.
//
// BuildDockerImage(dockerfilePath, imageName string, buildArgs map[string]string) error:
//
//	Builds a Docker image using the specified Dockerfile and image name.
//	Accepts build arguments as a map.
//	Outputs build progress and duration to stdout/stderr.
//	Returns an error if the build fails.
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fabianoflorentino/whiterose/utils"
)

type DockerManager struct {
	workDir string
}

// NewDockerManager creates a new DockerManager instance for the given working directory.
func NewDockerManager(workDir string) *DockerManager {
	return &DockerManager{workDir: workDir}
}

// DetectDockerFile searches for Dockerfiles in the working directory.
func (d *DockerManager) DetectDockerFile() ([]string, error) {
	var dockerfiles []string

	// Walk the file tree to find Dockerfiles
	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filename := strings.ToLower(filepath.Base(path))
			if filename == "dockerfile" || strings.HasPrefix(filename, "dockerfile.") {
				dockerfiles = append(dockerfiles, path)
			}
		}

		return nil
	}

	err := filepath.Walk(d.workDir, w)
	if err != nil {
		return dockerfiles, err
	}
	if len(dockerfiles) == 0 {
		return nil, fmt.Errorf("no Dockerfile found in directory %s", d.workDir)
	}

	return dockerfiles, nil
}

// BuildDockerImage builds a Docker image using the specified Dockerfile and image name.
func (d *DockerManager) BuildDockerImage(dockerfilePath, imageName string, buildArgs map[string]string) error {
	fmt.Printf("Building Docker image '%s' from Dockerfile at '%s'\n", imageName, dockerfilePath)

	buildContext := "."

	args := []string{"build"}

	// Adiciona build-args antes do contexto
	for key, value := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	var build_target string = utils.GetEnvOrDefault("BUILD_TARGET", "development")

	args = append(args, "--progress=plain", "--no-cache", "--target", build_target)
	args = append(args, "-t", imageName)
	args = append(args, "-f", dockerfilePath)
	args = append(args, buildContext)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running command: %s\n", cmd.String())

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("Error building Docker image: %v\n", err)
		return err
	}

	fmt.Printf("Docker image '%s' built successfully in %v\n", imageName, duration)
	return nil
}
