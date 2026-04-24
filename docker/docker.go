// Package docker provides Docker operations with SOLID principles.
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

type DockerClient interface {
	Build(dockerfile, image string, args []string) error
	Delete(image string) error
	ListImages(pattern string) ([]string, error)
}

type RealDockerClient struct{}

func (c *RealDockerClient) Build(dockerfile, image string, args []string) error {
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *RealDockerClient) Delete(image string) error {
	cmd := exec.Command("docker", "rmi", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_, err := cmd.CombinedOutput()
	return err
}

func (c *RealDockerClient) ListImages(pattern string) ([]string, error) {
	cmd := exec.Command("docker", "images", "--filter", "reference="+pattern, "--format", "{{.Repository}}:{{.Tag}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, line := range lines {
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

type DockerManager struct {
	workDir      string
	dockerClient DockerClient
}

func NewDockerManager(workDir string) *DockerManager {
	return &DockerManager{
		workDir:      workDir,
		dockerClient: &RealDockerClient{},
	}
}

func (dm *DockerManager) WithClient(client DockerClient) *DockerManager {
	return &DockerManager{
		workDir:      dm.workDir,
		dockerClient: client,
	}
}

func (dm *DockerManager) DetectDockerFile() ([]string, error) {
	var dockerfiles []string
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
	err := filepath.Walk(dm.workDir, w)
	if err != nil {
		return dockerfiles, err
	}
	if len(dockerfiles) == 0 {
		return nil, fmt.Errorf("no Dockerfile found in directory %s", dm.workDir)
	}
	return dockerfiles, nil
}

func (dm *DockerManager) BuildDockerImage(dockerfilePath, imageName string, buildArgs map[string]string) error {
	fmt.Printf("Building Docker image '%s' from Dockerfile at '%s'\n", imageName, dockerfilePath)

	args := []string{"build"}
	for key, value := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}
	args = append(args, "--progress=plain", "--no-cache", "--target", utils.GetEnvOrDefault("BUILD_TARGET", "development"))
	args = append(args, "-t", imageName, "-f", dockerfilePath, ".")

	startTime := time.Now()
	err := dm.dockerClient.Build(dockerfilePath, imageName, args)
	duration := time.Since(startTime)
	if err != nil {
		fmt.Printf("Error building Docker image: %v\n", err)
		return err
	}
	fmt.Printf("Docker image '%s' built successfully in %v\n", imageName, duration)
	return nil
}

func (dm *DockerManager) DeleteDockerImage(imageName string) error {
	fmt.Printf("Deleting Docker image '%s'\n", imageName)
	return dm.dockerClient.Delete(imageName)
}

func (dm *DockerManager) ListDockerImages(pattern string) error {
	fmt.Printf("Listing Docker images matching '%s'\n", pattern)
	images, err := dm.dockerClient.ListImages(pattern)
	if err != nil {
		fmt.Printf("Error listing images: %v\n", err)
		return err
	}
	if len(images) == 0 {
		fmt.Println("No images found matching the pattern.")
		return nil
	}
	fmt.Println("Found images:")
	for _, img := range images {
		fmt.Printf("  %s\n", img)
	}
	return nil
}