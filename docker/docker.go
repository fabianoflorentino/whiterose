package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DockerManager struct {
	workDir string
}

func NewDockerManager(workDir string) *DockerManager {
	return &DockerManager{workDir: workDir}
}

func (d *DockerManager) DetectDockerFile() ([]string, error) {
	var dockerfiles []string

	err := filepath.Walk(d.workDir, func(path string, info os.FileInfo, err error) error {
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
	})

	if err != nil {
		return dockerfiles, err
	}
	if len(dockerfiles) == 0 {
		return nil, fmt.Errorf("no Dockerfile found in directory %s", d.workDir)
	}

	return dockerfiles, nil
}

func (d *DockerManager) BuildDockerImage(dockerfilePath, imageName string, buildArgs map[string]string) error {
	fmt.Printf("Building Docker image '%s' from Dockerfile at '%s'\n", imageName, dockerfilePath)

	buildContext := "."

	args := []string{"build"}

	// Adiciona build-args antes do contexto
	for key, value := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, "--progress=plain", "--no-cache")
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
