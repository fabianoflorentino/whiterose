// GitCloneRepository loads repository configurations from a JSON file,
// sets authentication credentials and SSH key information from environment variables,
// and fetches/clones the repositories. If any error occurs during loading or fetching,
// the function logs the error and terminates the application.
package setup

import (
	"github.com/fabianoflorentino/whiterose/git"
	"github.com/fabianoflorentino/whiterose/prereq"
)

func PreReq() {
	p := prereq.NewAppValidator()
	p.ValidateApps()
}

// GitCloneRepository loads repository configurations from a JSON file,
// sets authentication credentials and SSH key information from environment variables,
// and fetches/clones the repositories. If any error occurs during loading or fetching,
// the function logs the error and terminates the application.
func GitCloneRepository() {
	g := git.NewGitRepository()
	g.Setup()
}
