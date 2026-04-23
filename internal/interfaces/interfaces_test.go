package interfaces

import (
	"testing"

	"github.com/fabianoflorentino/whiterose/utils"
)

func TestToAppInfo(t *testing.T) {
	utilsApps := []utils.AppInfo{
		{Name: "Go", Command: "go", RecommendedVersion: "1.20"},
		{Name: "Git", Command: "git", RecommendedVersion: "2.40"},
	}

	apps := ToAppInfo(utilsApps)
	if len(apps) != 2 {
		t.Fatalf("len = %d, want 2", len(apps))
	}
	if apps[0].Name != "Go" {
		t.Errorf("apps[0].Name = %v, want Go", apps[0].Name)
	}

	empty := ToAppInfo(nil)
	if len(empty) != 0 {
		t.Errorf("empty apps should return empty slice")
	}
}

func TestToRepoInfo(t *testing.T) {
	utilsRepos := []utils.RepoInfo{
		{URL: "https://github.com/test/repo1", Directory: "repo1"},
		{URL: "https://github.com/test/repo2", Directory: "repo2"},
	}

	repos := ToRepoInfo(utilsRepos)
	if len(repos) != 2 {
		t.Fatalf("len = %d, want 2", len(repos))
	}
	if repos[0].URL != "https://github.com/test/repo1" {
		t.Errorf("repos[0].URL = %v, want https://github.com/test/repo1", repos[0].URL)
	}

	empty := ToRepoInfo(nil)
	if len(empty) != 0 {
		t.Errorf("empty repos should return empty slice")
	}
}