package entities

import (
	"strings"
	"testing"
)

func TestParseUpdateConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		wantLen  int
		wantErr  bool
	}{
		{
			"valid yaml",
			`projects:
  - name: test
    path: ./test
    goMod:
      updateStrategy: minor`,
			1,
			false,
		},
		{
			"empty",
			`projects: []`,
			0,
			false,
		},
		{
			"invalid yaml",
			`invalid: [`,
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUpdateConfig([]byte(tt.config))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestGetTimestampedBranchName(t *testing.T) {
	name := GetTimestampedBranchName()

	if len(name) < 21 {
		t.Errorf("len = %d, want at least 21", len(name))
	}

	expectedPrefix := "update-"
	if !strings.HasPrefix(name, expectedPrefix) {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, name)
	}
}

func TestUpdateStrategy_String(t *testing.T) {
	tests := []struct {
		strat  UpdateStrategy
		want   string
	}{
		{StrategyPatch, "patch"},
		{StrategyMinor, "minor"},
		{StrategyMajor, "major"},
	}

	for _, tt := range tests {
		if got := tt.strat.String(); got != tt.want {
			t.Errorf("String() = %v, want %v", got, tt.want)
		}
	}
}

func TestParseUpdateStrategy(t *testing.T) {
	tests := []struct {
		input  string
		want   UpdateStrategy
	}{
		{"major", StrategyMajor},
		{"Major", StrategyMajor},
		{"MAJOR", StrategyMajor},
		{"minor", StrategyMinor},
		{"Minor", StrategyMinor},
		{"patch", StrategyPatch},
		{"Patch", StrategyPatch},
		{"unknown", StrategyPatch},
		{"", StrategyPatch},
	}

	for _, tt := range tests {
		if got := ParseUpdateStrategy(tt.input); got != tt.want {
			t.Errorf("ParseUpdateStrategy(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}