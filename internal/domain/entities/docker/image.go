package docker

import (
	"fmt"
	"strings"
	"time"
)

// Image represents a Docker image in domain.
type Image struct {
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

// NewImage creates a new Image instance after validating the input data.
func NewImage(name, tag string) (*Image, error) {
	if err := validateImageData(name, tag); err != nil {
		return nil, err
	}

	fullName := buildFullImageName(name, tag)

	return &Image{
		ID:        generateImageIR(fullName),
		Name:      name,
		Tag:       tag,
		FullName:  fullName,
		BuildArgs: make(map[string]string),
		Created:   time.Now(),
	}, nil
}

// AddBuildArg adds a build argument to the Docker image after validation.
func (i *Image) AddBuildArg(key, value string) error {
	if err := validateBuildArgs(map[string]string{key: value}); err != nil {
		return err
	}

	i.BuildArgs[key] = value

	return nil
}

// SetDockerFile sets the Dockerfile path for the Docker image.
func (i *Image) SetDockerFile(path string) error {
	if path == "" {
		return ErrDockerFilePathEmpty
	}

	i.Dockerfile = path

	return nil
}

// SetContext sets the build context for the Docker image.
func (i *Image) SetContext(context string) error {
	if context == "" {
		return ErrContextEmpty
	}

	return nil
}

// SetTarget sets the build target for the Docker image.
func (i *Image) SetTarget(target string) {
	i.Target = target
}

func (i *Image) Validate() error {
	if err := validateImageData(i.Name, i.Tag); err != nil {
		return err
	}

	if i.Dockerfile != "" && !strings.HasSuffix(i.Dockerfile, "Dockerfile") {
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

// validateImageData validates the Docker image name and tag.
func validateImageData(name, tag string) error {
	if name == "" {
		return ErrInvalidImageName
	}

	if tag == "" {
		return ErrInvalidImageTag
	}

	if strings.Contains(name, " ") {
		return ErrInvalidImageName
	}

	return nil
}

// buildFullImageName constructs the full image name from name and tag.
func buildFullImageName(name, tag string) string {
	return fmt.Sprintf("%s:%s", name, tag)
}

// generateImageIR generates a unique identifier for the Docker image.
func generateImageIR(fullName string) string {
	return fmt.Sprintf("img_%d", len(fullName)*int(time.Now().Unix())%10000)
}
