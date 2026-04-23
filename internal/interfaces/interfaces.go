package interfaces

import "github.com/fabianoflorentino/whiterose/utils"

type ConfigLoader interface {
	LoadRepositories() ([]RepoInfo, error)
	LoadApps() ([]AppInfo, error)
}

type RepoInfo struct {
	URL       string
	Directory string
}

type AppInfo struct {
	Name                string
	Command             string
	VersionFlag         string
	RecommendedVersion  string
	InstallInstructions map[string]string
}

type GitCloner interface {
	Clone(opts CloneOptions) error
	CloneAll(opts []CloneOptions) error
}

type CloneOptions struct {
	URL        string
	Directory string
	Username  string
	Password  string
	SSHKeyPath string
	SSHKeyName string
	Branch    string
}

type AppChecker interface {
	Validate() []AppValidation
	ValidateOne(name string) AppValidation
	List() []AppInfo
}

type AppValidation struct {
	Name               string
	Command            string
	IsInstalled        bool
	CurrentVersion     string
	RecommendedVersion string
	IsUpToDate         bool
	InstallInstruction string
	OS                string
}

type ImageBuilder interface {
	Build(opts ImageBuildOptions) error
	Delete(image string) error
	ListImages(pattern string) error
	FindDockerfiles(root string) ([]string, error)
}

type ImageBuildOptions struct {
	Dockerfile string
	ImageName  string
	Context   string
	BuildArgs map[string]string
	Target   string
	NoCache  bool
}

type Executor interface {
	Run(cmd string, args ...string) (string, error)
	Which(cmd string) (string, error)
}

type FS interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	Exists(path string) bool
	MkdirAll(path string) error
}

type Logger interface {
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func ToAppInfo(apps []utils.AppInfo) []AppInfo {
	var result []AppInfo
	for _, a := range apps {
		result = append(result, AppInfo{
			Name:                a.Name,
			Command:             a.Command,
			VersionFlag:         a.VersionFlag,
			RecommendedVersion:  a.RecommendedVersion,
			InstallInstructions: a.InstallInstructions,
		})
	}
	return result
}

func ToRepoInfo(repos []utils.RepoInfo) []RepoInfo {
	var result []RepoInfo
	for _, r := range repos {
		result = append(result, RepoInfo{
			URL:       r.URL,
			Directory: r.Directory,
		})
	}
	return result
}