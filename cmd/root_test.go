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
}

func TestRootCmd_Execute(t *testing.T) {
	t.Skip("Execute requires full setup")
}

func TestExecute(t *testing.T) {
	t.Skip("Execute requires full setup")
}

func TestSetupCmd(t *testing.T) {
	t.Skip("SetupCmd requires config")
}

func TestPreReqCmd(t *testing.T) {
	t.Skip("PreReqCmd requires config")
}

func TestDockerCmd(t *testing.T) {
	t.Skip("DockerCmd requires config")
}
