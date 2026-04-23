package services

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fabianoflorentino/whiterose/internal/interfaces"
	"github.com/fabianoflorentino/whiterose/utils"
)

type Executor interface {
	Run(cmd string, args ...string) (string, error)
	Which(cmd string) (string, error)
}

type ConfigService struct {
	configPath string
}

func NewConfigService(configPath string) *ConfigService {
	return &ConfigService{configPath: configPath}
}

func (c *ConfigService) LoadRepositories() ([]interfaces.RepoInfo, error) {
	if c.configPath == "" {
		cfg, err := utils.LoadDotConfig()
		if err != nil {
			return nil, fmt.Errorf("config not found: %w", err)
		}
		c.configPath = cfg
	}

	repos, err := utils.FetchRepositories(c.configPath)
	if err != nil {
		return nil, err
	}

	return interfaces.ToRepoInfo(repos), nil
}

func (c *ConfigService) LoadApps() ([]interfaces.AppInfo, error) {
	if c.configPath == "" {
		cfg, err := utils.LoadDotConfig()
		if err != nil {
			cfg = ""
		} else {
			c.configPath = cfg
		}
	}

	var apps []interfaces.AppInfo
	if c.configPath != "" {
		utilsApps, err := utils.FetchAppsInfo(c.configPath)
		if err == nil {
			apps = interfaces.ToAppInfo(utilsApps)
		}
	}

	if len(apps) == 0 {
		apps = defaultApps()
	}

	return apps, nil
}

func defaultApps() []interfaces.AppInfo {
	return []interfaces.AppInfo{
		{
			Name:                "Go",
			Command:             "go",
			VersionFlag:         "version",
			RecommendedVersion:  "1.25.0",
			InstallInstructions: map[string]string{"linux": "sudo apt install golang", "darwin": "brew install go"},
		},
		{
			Name:                "Git",
			Command:             "git",
			VersionFlag:         "--version",
			RecommendedVersion:  "2.40.0",
			InstallInstructions: map[string]string{"linux": "sudo apt install git", "darwin": "brew install git"},
		},
		{
			Name:                "Docker",
			Command:             "docker",
			VersionFlag:         "--version",
			RecommendedVersion:  "20.10.0",
			InstallInstructions: map[string]string{"linux": "curl https://get.docker.com", "darwin": "brew install docker"},
		},
	}
}

type ExecutorService struct{}

func NewExecutorService() *ExecutorService {
	return &ExecutorService{}
}

func (e *ExecutorService) Run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("%s %v failed: %w", cmd, args, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (e *ExecutorService) Which(cmd string) (string, error) {
	c := exec.Command("which", cmd)
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("command not found: %s", cmd)
	}
	return strings.TrimSpace(string(out)), nil
}

type AppValidatorService struct {
	executor Executor
	osName   string
}

func NewAppValidatorService() *AppValidatorService {
	return &AppValidatorService{
		executor: NewExecutorService(),
		osName:   runtime.GOOS,
	}
}

func (a *AppValidatorService) WithExecutor(exec interfaces.Executor) *AppValidatorService {
	a.executor = exec
	return a
}

func (a *AppValidatorService) ValidateOne(app interfaces.AppInfo) interfaces.AppValidation {
	result := interfaces.AppValidation{
		Name:               app.Name,
		Command:            app.Command,
		RecommendedVersion: app.RecommendedVersion,
		OS:                 a.osName,
	}

	out, err := a.executor.Run(app.Command, app.VersionFlag)
	if err != nil {
		result.IsInstalled = false
		result.InstallInstruction = app.InstallInstructions[a.osName]
		return result
	}

	result.IsInstalled = true
	result.CurrentVersion = strings.TrimSpace(out)
	result.IsUpToDate = result.CurrentVersion == app.RecommendedVersion

	return result
}

func (a *AppValidatorService) Validate(apps []interfaces.AppInfo) []interfaces.AppValidation {
	var results []interfaces.AppValidation
	for _, app := range apps {
		results = append(results, a.ValidateOne(app))
	}
	return results
}

func (a *AppValidatorService) List() []interfaces.AppInfo {
	return defaultApps()
}

type FileSystemService struct{}

func NewFileSystemService() *FileSystemService {
	return &FileSystemService{}
}

func (fs *FileSystemService) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *FileSystemService) Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (fs *FileSystemService) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (fs *FileSystemService) MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}