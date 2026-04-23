package docker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDockerManager(t *testing.T) {
	dm := NewDockerManager("/tmp")
	if dm == nil {
		t.Error("NewDockerManager() returned nil")
	}
}

func TestNewDockerManager_WithPath(t *testing.T) {
	dm := NewDockerManager("/app")
	if dm.workDir != "/app" {
		t.Errorf("workDir = %v, want /app", dm.workDir)
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
	os.MkdirAll(subDir, 0755)
	
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

func TestDockerManager_DeleteDockerImage(t *testing.T) {
	dm := NewDockerManager("/tmp")
	err := dm.DeleteDockerImage("nonexistent-image")
	if err == nil {
		t.Error("expected error for nonexistent image")
	}
}

func TestDockerManager_ListDockerImages(t *testing.T) {
	dm := NewDockerManager("/tmp")
	dm.ListDockerImages("nonexistent-image")
}