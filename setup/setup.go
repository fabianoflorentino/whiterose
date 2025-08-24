package setup

import (
	"log"

	"github.com/fabianoflorentino/whiterose/git"
	"github.com/fabianoflorentino/whiterose/utils"
)

func GitCloneRepository() {
	repos, err := utils.FetchReposFromJSON("./repos.json")
	if err != nil {
		log.Fatalf("failed to fetch repositories: %v", err)
	}

	var opts []git.GitCloneOptions
	for _, r := range repos {
		opts = append(opts, git.GitCloneOptions{
			URL:        r.URL,
			Directory:  r.Directory,
			Username:   utils.GetEnvOrDefault("GIT_USER", ""),
			Password:   utils.GetEnvOrDefault("GIT_TOKEN", ""),
			SSHKeyPath: utils.GetEnvOrDefault("SSH_KEY_PATH", ""),
			SSHKeyName: utils.GetEnvOrDefault("SSH_KEY_NAME", "id_rsa"),
		})
	}

	if err := git.FetchRepositories(opts); err != nil {
		log.Fatalf("failed to fetch repositories: %v", err)
	}
}
