package docker

import "strings"

// DockerBuildOptions represents options for building a Docker image.
type BuildOptions struct {
	ImageName  string
	Tag        string
	Dockerfile string
	Context    string
	BuildArgs  map[string]string
	Target     string
	NoCache    bool
	Progress   string
}

// NewBuildOptions creates a new BuildOptions instance with default values.
func NewBuildOptions(imageName, tag string) (*BuildOptions, error) {
	if err := validateImageData(imageName, tag); err != nil {
		return nil, err
	}

	return &BuildOptions{
		ImageName: imageName,
		Context:   ".",
		Tag:       tag,
		BuildArgs: make(map[string]string),
		NoCache:   false,
		Progress:  "auto",
	}, nil
}

// AddBuildArg adds a build argument to the BuildOptions after validation.
func (bo *BuildOptions) AddBuildArg(key, value string) error {
	if key == "" {
		return ErrArgumentKeyNotBeEmpty
	}

	if err := validateBuildArgKey(key); err != nil {
		return err
	}

	bo.BuildArgs[key] = value
	return nil
}

// SetDockerfile sets the Dockerfile path for the BuildOptions.
func (bo *BuildOptions) SetDockerfile(path string) error {
	if path == "" {
		return ErrDockerFilePathEmpty
	}

	bo.Dockerfile = path

	return nil
}

// validateBuildArgKey checks if the build argument key is valid.
func validateBuildArgKey(key string) error {
	var specialChars = "!@#$%^&*()_-+=[]{}|;:'\",.<>?/`~"

	isSpecialCharOnly := true
	for _, char := range key {
		if !strings.ContainsRune(specialChars, char) {
			isSpecialCharOnly = false
			break
		}
	}

	if isSpecialCharOnly {
		return ErrArgumentKeyContainsOnlySpecialCharacters
	}

	return nil
}

// GetFullImageName returns the full image name with tag.
func (bo *BuildOptions) GetFullImageName() string {
	return buildFullImageName(bo.ImageName, bo.Tag)
}

// Validate checks if the BuildOptions are valid.
func (bo *BuildOptions) Validate() error {
	if err := validateImageData(bo.ImageName, bo.Tag); err != nil {
		return err
	}

	if bo.Context == "" {
		return ErrContextEmpty
	}

	return nil
}
