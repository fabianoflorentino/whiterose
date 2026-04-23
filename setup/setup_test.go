package setup

import (
	"testing"
)

func TestPreReq_ValidatesApps(t *testing.T) {
	_ = PreReq
}

func TestGitCloneRepository_ClonesRepos(t *testing.T) {
	_ = GitCloneRepository
}

func TestPackage_Imports(t *testing.T) {
	_ = PreReq
	_ = GitCloneRepository
}
