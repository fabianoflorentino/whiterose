package docker

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrInvalidImageName", ErrInvalidImageName, "invalid docker image name"},
		{"ErrInvalidImageTag", ErrInvalidImageTag, "invalid docker image tag"},
		{"ErrArgumentKeyNotBeEmpty", ErrArgumentKeyNotBeEmpty, "build arg key cannot be empty"},
		{"ErrArgumentKeyNotBeSpecialChars", ErrArgumentKeyNotBeSpecialChars, "build arg key cannot consist solely of special characters"},
		{"ErrDockerFilePathEmpty", ErrDockerFilePathEmpty, "dockerfile path cannot be empty"},
		{"ErrContextEmpty", ErrContextEmpty, "context path cannot be empty"},
		{"ErrDockerFileNotFound", ErrDockerFileNotFound, "dockerfile not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("expected error, got nil")
				return
			}
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorsAreValid(t *testing.T) {
	if !errors.Is(ErrInvalidImageName, ErrInvalidImageName) {
		t.Error("ErrInvalidImageName should match itself")
	}
	if !errors.Is(ErrDockerFileNotFound, ErrDockerFileNotFound) {
		t.Error("ErrDockerFileNotFound should match itself")
	}
}

func TestDifferentErrorsNotEqual(t *testing.T) {
	if ErrInvalidImageName == ErrInvalidImageTag {
		t.Error("different errors should not be equal")
	}
}