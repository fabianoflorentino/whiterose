package docker

import "errors"

var (
	ErrInvalidDockerImageName                   error = errors.New("invalid docker image name")
	ErrInvalidDockerImageTag                    error = errors.New("invalid docker image tag")
	ErrArgumentKeyNotBeEmpty                    error = errors.New("build arg key cannot be empty")
	ErrArgumentKeyNotBeSpecialChars             error = errors.New("build arg key cannot consist solely of special characters")
	ErrDockerFilePathEmpty                      error = errors.New("dockerfile path cannot be empty")
	ErrContextEmpty                             error = errors.New("context path cannot be empty")
	ErrDockerFileNotFound                       error = errors.New("dockerfile not found")
	ErrArgumentKeyContainsOnlySpecialCharacters error = errors.New("build arg key cannot contain only special characters")
)
