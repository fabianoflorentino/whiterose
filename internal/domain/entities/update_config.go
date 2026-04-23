package entities

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type UpdateStrategy string

const (
	StrategyPatch UpdateStrategy = "patch"
	StrategyMinor UpdateStrategy = "minor"
	StrategyMajor UpdateStrategy = "major"
)

type GoModConfig struct {
	UpdateStrategy UpdateStrategy `json:"updateStrategy" yaml:"updateStrategy"`
}

type DockerImageConfig struct {
	Base        string         `json:"base" yaml:"base"`
	UpdateStrategy UpdateStrategy `json:"updateStrategy" yaml:"updateStrategy"`
}

type GoVersionConfig struct {
	Version     string        `json:"version" yaml:"version"`
	UpdateStrategy UpdateStrategy `json:"updateStrategy" yaml:"updateStrategy"`
}

type UpdateProject struct {
	Name            string           `json:"name" yaml:"name"`
	Path            string           `json:"path" yaml:"path"`
	GoMod           *GoModConfig     `json:"goMod" yaml:"goMod"`
	GoVersion       *GoVersionConfig `json:"goVersion" yaml:"goVersion"`
	DockerImage     *DockerImageConfig `json:"dockerImage" yaml:"dockerImage"`
}

type UpdateConfig struct {
	Projects []UpdateProject `json:"projects" yaml:"projects"`
}

func ParseUpdateConfig(data []byte) ([]UpdateProject, error) {
	config := UpdateConfig{}
	if err := unmarshalYAML(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return config.Projects, nil
}

func unmarshalYAML(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func GetTimestampedBranchName() string {
	ts := time.Now().Format("20060102-150405")
	return fmt.Sprintf("update-%s", ts)
}

func (s UpdateStrategy) String() string {
	return string(s)
}

func ParseUpdateStrategy(s string) UpdateStrategy {
	switch strings.ToLower(s) {
	case "major":
		return StrategyMajor
	case "minor":
		return StrategyMinor
	default:
		return StrategyPatch
	}
}