package setup

import (
	"log"

	"github.com/fabianoflorentino/whiterose/git"
	"github.com/fabianoflorentino/whiterose/utils"
)

func GitCloneRepository() {
	repos, err := git.LoadRepositoriesFromFile("./config.json")
	if err != nil {
		log.Fatalf("failed to load repositories: %v", err)
	}

	for i := range repos {
		repos[i].Username = utils.GetEnvOrDefault("GIT_USER", "")
		repos[i].Password = utils.GetEnvOrDefault("GIT_TOKEN", "")
		repos[i].SSHKeyPath = utils.GetEnvOrDefault("SSH_KEY_PATH", "")
		repos[i].SSHKeyName = utils.GetEnvOrDefault("SSH_KEY_NAME", "id_rsa")
	}

	if err := git.FetchRepositories(repos); err != nil {
		log.Fatalf("failed to fetch repositories: %v", err)
	}
}
