package cmd

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	if rootCmd.Use != "whiterose" {
		t.Errorf("Use = %v, want whiterose", rootCmd.Use)
	}
	if rootCmd.Short == "" {
		t.Error("Short should not be empty")
	}
	commands := rootCmd.Commands()
	if len(commands) < 4 {
		t.Errorf("len(Commands) = %d, want at least 4", len(commands))
	}
}

func TestSetupCmd(t *testing.T) {
	if setupCmd.Use != "setup" {
		t.Errorf("Use = %v, want setup", setupCmd.Use)
	}
	if setupCmd.Short == "" {
		t.Error("Short should not be empty")
	}
	flags := setupCmd.PersistentFlags()
	if flags.Lookup("all") == nil {
		t.Error("all flag should exist")
	}
	if flags.Lookup("pre-req") == nil {
		t.Error("pre-req flag should exist")
	}
	if flags.Lookup("repos") == nil {
		t.Error("repos flag should exist")
	}
}

func TestPreReqCmd(t *testing.T) {
	if preReqCmd.Use != "pre-req" {
		t.Errorf("Use = %v, want pre-req", preReqCmd.Use)
	}
	if preReqCmd.Short == "" {
		t.Error("Short should not be empty")
	}
	flags := preReqCmd.Flags()
	if flags.Lookup("check") == nil {
		t.Error("check flag should exist")
	}
	if flags.Lookup("list") == nil {
		t.Error("list flag should exist")
	}
	if flags.Lookup("apps") == nil {
		t.Error("apps flag should exist")
	}
}

func TestDockerCmd(t *testing.T) {
	if dockerCmd.Use != "docker" {
		t.Errorf("Use = %v, want docker", dockerCmd.Use)
	}
	if dockerCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestUpdateCmd(t *testing.T) {
	if updateCmd.Use != "update" {
		t.Errorf("Use = %v, want update", updateCmd.Use)
	}
	if updateCmd.Short == "" {
		t.Error("Short should not be empty")
	}
	flags := updateCmd.Flags()
	if flags.Lookup("go-mod") == nil {
		t.Error("go-mod flag should exist")
	}
	if flags.Lookup("go-version") == nil {
		t.Error("go-version flag should exist")
	}
	if flags.Lookup("docker-image") == nil {
		t.Error("docker-image flag should exist")
	}
	if flags.Lookup("list") == nil {
		t.Error("list flag should exist")
	}
	if flags.Lookup("major") == nil {
		t.Error("major flag should exist")
	}
	if flags.Lookup("dry-run") == nil {
		t.Error("dry-run flag should exist")
	}
	if flags.Lookup("pr") == nil {
		t.Error("pr flag should exist")
	}
	if flags.Lookup("base") == nil {
		t.Error("base flag should exist")
	}
	if flags.Lookup("config") == nil {
		t.Error("config flag should exist")
	}
}

func TestAllCommands_HaveShortDescription(t *testing.T) {
	commands := rootCmd.Commands()
	for _, c := range commands {
		if c.Short == "" {
			t.Errorf("Command %s has empty Short description", c.Use)
		}
	}
}