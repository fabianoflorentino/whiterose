package docker

import (
	"testing"
)

func TestNewMissingStatus(t *testing.T) {
	s := NewMissingStatus()
	if s.Status != StatusMissing {
		t.Errorf("Status = %v, want %v", s.Status, StatusMissing)
	}
}

func TestNewBuildingStatus(t *testing.T) {
	s := NewBuildingStatus()
	if s.Status != StatusBuilding {
		t.Errorf("Status = %v, want %v", s.Status, StatusBuilding)
	}
}

func TestNewReadyStatus(t *testing.T) {
	s := NewReadyStatus()
	if s.Status != StatusReady {
		t.Errorf("Status = %v, want %v", s.Status, StatusReady)
	}
}

func TestNewFailedStatus(t *testing.T) {
	err := &testError{"test error"}
	s := NewFailedStatus(err)
	if s.Status != StatusFailed {
		t.Errorf("Status = %v, want %v", s.Status, StatusFailed)
	}
	if s.Error != "test error" {
		t.Errorf("Error = %v, want 'test error'", s.Error)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}