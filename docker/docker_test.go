package docker

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fabianoflorentino/whiterose/mocks"
)

func TestNewDockerManager(t *testing.T) {
	dm := NewDockerManager("/tmp")
	if dm == nil {
		t.Error("NewDockerManager() returned nil")
	}
	if dm.workDir != "/tmp" {
		t.Errorf("workDir = %v, want /tmp", dm.workDir)
	}
}

func TestNewDockerManager_WithClient(t *testing.T) {
	dm := NewDockerManager("/tmp")
	mock := &mocks.MockDockerClient{
		BuildFunc: func(dockerfile, image string, args []string) error {
			return nil
		},
	}

	dm2 := dm.WithClient(mock)
	if dm2 == nil {
		t.Error("WithClient returned nil")
	}
}

func TestDockerManager_DetectDockerFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	dm := NewDockerManager(tmpDir)
	_, err := dm.DetectDockerFile()
	if err == nil {
		t.Error("expected error for no dockerfile")
	}
}

func TestDockerManager_DetectDockerFile_Found(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte("FROM golang:1.20"), 0644); err != nil {
		t.Fatalf("failed to create dockerfile: %v", err)
	}

	dm := NewDockerManager(tmpDir)
	files, err := dm.DetectDockerFile()
	if err != nil {
		t.Errorf("DetectDockerFile() error = %v", err)
	}
	if len(files) != 1 {
		t.Errorf("len(files) = %d, want 1", len(files))
	}
}

func TestDockerManager_DetectDockerFile_Named(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile.app")
	if err := os.WriteFile(dockerfilePath, []byte("FROM node:20"), 0644); err != nil {
		t.Fatalf("failed to create dockerfile: %v", err)
	}

	dm := NewDockerManager(tmpDir)
	files, err := dm.DetectDockerFile()
	if err != nil {
		t.Errorf("DetectDockerFile() error = %v", err)
	}
	if len(files) != 1 {
		t.Errorf("len(files) = %d, want 1", len(files))
	}
}

func TestDockerManager_DetectDockerFile_SubDir(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	dockerfilePath := filepath.Join(subDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte("FROM python:3.11"), 0644); err != nil {
		t.Fatalf("failed to create dockerfile: %v", err)
	}

	dm := NewDockerManager(tmpDir)
	files, err := dm.DetectDockerFile()
	if err != nil {
		t.Errorf("DetectDockerFile() error = %v", err)
	}
	if len(files) != 1 {
		t.Errorf("len(files) = %d, want 1", len(files))
	}
}

func TestDockerManager_BuildDockerImage(t *testing.T) {
	buildCalled := false
	mock := &mocks.MockDockerClient{
		BuildFunc: func(dockerfile, image string, args []string) error {
			buildCalled = true
			if image != "test-image" {
				t.Errorf("image = %v, want test-image", image)
			}
			return nil
		},
	}

	dm := NewDockerManager("/tmp").WithClient(mock)
	err := dm.BuildDockerImage("Dockerfile", "test-image", map[string]string{})
	if err != nil {
		t.Errorf("BuildDockerImage() error = %v", err)
	}
	if !buildCalled {
		t.Error("Build was not called")
	}
}

func TestDockerManager_BuildDockerImage_Error(t *testing.T) {
	mock := &mocks.MockDockerClient{
		BuildFunc: func(dockerfile, image string, args []string) error {
			return errors.New("docker build failed")
		},
	}

	dm := NewDockerManager("/tmp").WithClient(mock)
	err := dm.BuildDockerImage("Dockerfile", "test-image", map[string]string{})
	if err == nil {
		t.Error("expected error from BuildDockerImage")
	}
}

func TestDockerManager_DeleteDockerImage(t *testing.T) {
	deleteCalled := false
	deleteImage := ""
	mock := &mocks.MockDockerClient{
		DeleteFunc: func(image string) error {
			deleteCalled = true
			deleteImage = image
			return nil
		},
	}

	dm := NewDockerManager("/tmp").WithClient(mock)
	err := dm.DeleteDockerImage("test-image")
	if err != nil {
		t.Errorf("DeleteDockerImage() error = %v", err)
	}
	if !deleteCalled {
		t.Error("Delete was not called")
	}
	if deleteImage != "test-image" {
		t.Errorf("deleteImage = %v, want test-image", deleteImage)
	}
}

func TestDockerManager_ListDockerImages(t *testing.T) {
	mock := &mocks.MockDockerClient{
		ListFunc: func(pattern string) ([]string, error) {
			return []string{"image1:latest", "image2:v1.0"}, nil
		},
	}

	dm := NewDockerManager("/tmp").WithClient(mock)
	err := dm.ListDockerImages("test-*")
	if err != nil {
		t.Errorf("ListDockerImages() error = %v", err)
	}
}

func TestDockerManager_ListDockerImages_Error(t *testing.T) {
	mock := &mocks.MockDockerClient{
		ListFunc: func(pattern string) ([]string, error) {
			return nil, errors.New("docker daemon not running")
		},
	}

	dm := NewDockerManager("/tmp").WithClient(mock)
	err := dm.ListDockerImages("test-*")
	if err == nil {
		t.Error("expected error from ListDockerImages")
	}
}