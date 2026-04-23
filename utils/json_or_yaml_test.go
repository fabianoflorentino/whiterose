package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileJSON(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{"json file", "config.json", true},
		{"yaml file", "config.yaml", false},
		{"yml file", "config.yml", false},
		{"no extension", "config", false},
		{"wrong extension", "config.txt", false},
		{"nested path", "dir/config.json", true},
		{"path with dots", "config.min.json", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileJSON(tt.file); got != tt.want {
				t.Errorf("IsFileJSON(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}

func TestIsFileYAML(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{"yaml file", "config.yaml", true},
		{"yml file", "config.yml", true},
		{"json file", "config.json", false},
		{"no extension", "config", false},
		{"wrong extension", "config.txt", false},
		{"nested path", "dir/config.yaml", true},
		{"path with dots", "config.min.yaml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileYAML(tt.file); got != tt.want {
				t.Errorf("IsFileYAML(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}

func TestYmlOrYamlExistsInHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("could not determine home directory")
	}

	tests := []struct {
		name     string
		setup    func() error
		teardown func() error
		want     bool
	}{
		{
			name:     "yaml exists",
			setup:    func() error { return os.WriteFile(filepath.Join(homeDir, ".config.yaml"), []byte(""), 0644) },
			teardown: func() error { return os.Remove(filepath.Join(homeDir, ".config.yaml")) },
			want:     true,
		},
		{
			name:     "yml exists",
			setup:    func() error { return os.WriteFile(filepath.Join(homeDir, ".config.yml"), []byte(""), 0644) },
			teardown: func() error { return os.Remove(filepath.Join(homeDir, ".config.yml")) },
			want:     true,
		},
		{
			name:     "no file",
			setup:    func() error { return nil },
			teardown: func() error { return nil },
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}
			defer func() { _ = tt.teardown() }()

			if got := YmlOrYamlExistsInHomeDir(); got != tt.want {
				t.Errorf("YmlOrYamlExistsInHomeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
