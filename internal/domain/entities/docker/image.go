package docker

import (
	"fmt"
	"strings"
	"time"
)

// DockerImage represents a Docker image in domain.
type DockerImage struct {
	ID         string
	Name       string
	Tag        string
	FullName   string
	Size       int64
	Created    time.Time
	BuildArgs  map[string]string
	Target     string
	Context    string
	Dockerfile string
}

// NewDockerImage creates a new DockerImage instance after validating the input data.
func NewDockerImage(name, tag string) (*DockerImage, error) {
	if err := validateDockerImageData(name, tag); err != nil {
		return nil, err
	}

	fullName := buildFullImageName(name, tag)

	return &DockerImage{
		ID:        generateDockerImageIR(fullName),
		Name:      name,
		Tag:       tag,
		FullName:  fullName,
		BuildArgs: make(map[string]string),
		Created:   time.Now(),
	}, nil
}

// AddBuildArg adds a build argument to the Docker image after validation.
func (di *DockerImage) AddBuildArg(key, value string) error {
	if err := validateBuildArgs(map[string]string{key: value}); err != nil {
		return err
	}

	di.BuildArgs[key] = value

	return nil
}

// SetDockerFile sets the Dockerfile path for the Docker image.
func (di *DockerImage) SetDockerFile(path string) error {
	if path == "" {
		return ErrDockerFilePathEmpty
	}

	di.Dockerfile = path

	return nil
}

// SetContext sets the build context for the Docker image.
func (di *DockerImage) SetContext(context string) error {
	if context == "" {
		return ErrContextEmpty
	}

	return nil
}

// SetTarget sets the build target for the Docker image.
func (di *DockerImage) SetTarget(target string) {
	di.Target = target
}

func (di *DockerImage) Validate() error {
	if err := validateDockerImageData(di.Name, di.Tag); err != nil {
		return err
	}

	if di.Dockerfile != "" && !strings.HasSuffix(di.Dockerfile, "Dockerfile") {
		return ErrDockerFileNotFound
	}

	return nil
}

// validateBuildArgs validates the build arguments.
func validateBuildArgs(args map[string]string) error {
	var spcecialChars = "!@#$%^&*()_-+=[]{}|;:'\",.<>?/`~"

	for k := range args {
		if k == "" {
			return ErrArgumentKeyNotBeEmpty
		}

		if k == strings.Trim(k, spcecialChars) {
			return ErrArgumentKeyNotBeSpecialChars
		}
	}
	return nil
}

// validateDockerImageData validates the Docker image name and tag.
func validateDockerImageData(name, tag string) error {
	if name == "" {
		return ErrInvalidDockerImageName
	}

	if tag == "" {
		return ErrInvalidDockerImageTag
	}

	if strings.Contains(name, " ") {
		return ErrInvalidDockerImageName
	}

	return nil
}

// buildFullImageName constructs the full image name from name and tag.
func buildFullImageName(name, tag string) string {
	return fmt.Sprintf("%s:%s", name, tag)
}

// generateDockerImageIR generates a unique identifier for the Docker image.
func generateDockerImageIR(fullName string) string {
	return fmt.Sprintf("img_%d", len(fullName)*int(time.Now().Unix())%10000)
}
