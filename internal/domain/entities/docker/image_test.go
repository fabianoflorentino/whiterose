package docker

import (
	"testing"
)

func TestNewImage(t *testing.T) {
	img, err := NewImage("golang", "1.20")
	if err != nil {
		t.Fatalf("NewImage() error = %v", err)
	}
	if img.Name != "golang" {
		t.Errorf("Name = %v, want golang", img.Name)
	}
	if img.Tag != "1.20" {
		t.Errorf("Tag = %v, want 1.20", img.Tag)
	}
	if img.FullName != "golang:1.20" {
		t.Errorf("FullName = %v, want golang:1.20", img.FullName)
	}
}

func TestNewImage_InvalidName(t *testing.T) {
	_, err := NewImage("", "latest")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewImage_InvalidTag(t *testing.T) {
	_, err := NewImage("golang", "")
	if err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestNewImage_NameWithSpaces(t *testing.T) {
	_, err := NewImage("golang alpine", "latest")
	if err == nil {
		t.Error("expected error for name with spaces")
	}
}

func TestImage_AddBuildArg(t *testing.T) {
	img, _ := NewImage("golang", "1.20")
	
	err := img.AddBuildArg("key=value", "value")
	if err == nil {
		t.Error("expected error for key with =")
	}
	
	err = img.AddBuildArg("key-name", "value")
	if err == nil {
		t.Error("expected error for key with -")
	}
}

func TestImage_SetDockerFile(t *testing.T) {
	img, _ := NewImage("golang", "1.20")
	
	err := img.SetDockerFile("")
	if err == nil {
		t.Error("expected error for empty path")
	}
	
	err = img.SetDockerFile("Dockerfile")
	if err != nil {
		t.Errorf("SetDockerFile() error = %v", err)
	}
}

func TestImage_SetContext(t *testing.T) {
	img, _ := NewImage("golang", "1.20")
	
	err := img.SetContext("")
	if err == nil {
		t.Error("expected error for empty context")
	}
	
	err = img.SetContext(".")
	if err != nil {
		t.Errorf("SetContext() error = %v", err)
	}
}

func TestImage_SetTarget(t *testing.T) {
	img, _ := NewImage("golang", "1.20")
	img.SetTarget("development")
	
	if img.Target != "development" {
		t.Errorf("Target = %v, want development", img.Target)
	}
}

func TestImage_Validate(t *testing.T) {
	img := &Image{
		Name:      "golang",
		Tag:       "1.20",
		Dockerfile: "Dockerfile",
	}
	
	err := img.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestImage_Validate_Invalid(t *testing.T) {
	img := &Image{
		Name:      "golang",
		Tag:       "1.20",
		Dockerfile: "dockerfile.txt",
	}
	
	err := img.Validate()
	if err == nil {
		t.Error("expected error for invalid dockerfile")
	}
}

func TestBuildFullImageName(t *testing.T) {
	tests := []struct {
		name, tag, want string
	}{
		{"golang", "1.20", "golang:1.20"},
		{"alpine", "latest", "alpine:latest"},
		{"nginx", "alpine", "nginx:alpine"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name+":"+tt.tag, func(t *testing.T) {
			if got := buildFullImageName(tt.name, tt.tag); got != tt.want {
				t.Errorf("buildFullImageName(%q, %q) = %v, want %v", tt.name, tt.tag, got, tt.want)
			}
		})
	}
}