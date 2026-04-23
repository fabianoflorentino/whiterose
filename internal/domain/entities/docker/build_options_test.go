package docker

import (
	"testing"
)

func TestNewBuildOptions(t *testing.T) {
	bo, err := NewBuildOptions("golang", "1.20")
	if err != nil {
		t.Fatalf("NewBuildOptions() error = %v", err)
	}
	if bo.ImageName != "golang" {
		t.Errorf("ImageName = %v, want golang", bo.ImageName)
	}
	if bo.Tag != "1.20" {
		t.Errorf("Tag = %v, want 1.20", bo.Tag)
	}
	if bo.Context != "." {
		t.Errorf("Context = %v, want .", bo.Context)
	}
}

func TestNewBuildOptions_Invalid(t *testing.T) {
	_, err := NewBuildOptions("", "latest")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestBuildOptions_AddBuildArg(t *testing.T) {
	bo, _ := NewBuildOptions("golang", "1.20")

	err := bo.AddBuildArg("", "value")
	if err == nil {
		t.Error("expected error for empty key")
	}

	err = bo.AddBuildArg("GO_VERSION", "1.20")
	if err != nil {
		t.Errorf("AddBuildArg() error = %v", err)
	}
}

func TestBuildOptions_AddBuildArg_InvalidKey(t *testing.T) {
	bo, _ := NewBuildOptions("golang", "1.20")

	err := bo.AddBuildArg("!!!", "value")
	if err == nil {
		t.Error("expected error for special chars only")
	}
}

func TestBuildOptions_SetDockerfile(t *testing.T) {
	bo, _ := NewBuildOptions("golang", "1.20")

	err := bo.SetDockerfile("")
	if err == nil {
		t.Error("expected error for empty path")
	}

	err = bo.SetDockerfile("Dockerfile")
	if err != nil {
		t.Errorf("SetDockerfile() error = %v", err)
	}
}

func TestBuildOptions_GetFullImageName(t *testing.T) {
	bo, _ := NewBuildOptions("golang", "1.20")
	if got := bo.GetFullImageName(); got != "golang:1.20" {
		t.Errorf("GetFullImageName() = %v, want golang:1.20", got)
	}
}

func TestBuildOptions_Validate(t *testing.T) {
	bo := &BuildOptions{
		ImageName: "golang",
		Tag:       "1.20",
		Context:   ".",
	}

	err := bo.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestBuildOptions_Validate_InvalidContext(t *testing.T) {
	bo := &BuildOptions{
		ImageName: "golang",
		Tag:       "1.20",
		Context:   "",
	}

	err := bo.Validate()
	if err == nil {
		t.Error("expected error for empty context")
	}
}
