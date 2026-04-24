package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Git   GitConfig
	SSH   SSHConfig
	Image ImageConfig
	Repo  RepoConfig
}

type GitConfig struct {
	User string
	Token string
	Base  string
}

type SSHConfig struct {
	KeyPath string
	KeyName string
}

type ImageConfig struct {
	Name    string
	Version string
}

type RepoConfig struct {
	Path string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("WHITEROSE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("git.base", "main")
	v.SetDefault("ssh.keyName", "id_rsa")
	v.SetDefault("image.name", "my_app")
	v.SetDefault("image.version", "latest")
	v.SetDefault("repo.path", ".config.json")

	_ = v.BindEnv("git.user", "GIT_USER")
	_ = v.BindEnv("git.token", "GIT_TOKEN")
	_ = v.BindEnv("ssh.keyPath", "SSH_KEY_PATH")
	_ = v.BindEnv("ssh.keyName", "SSH_KEY_NAME")
	_ = v.BindEnv("image.name", "IMAGE_NAME")
	_ = v.BindEnv("image.version", "IMAGE_VERSION")
	_ = v.BindEnv("repo.path", "CONFIG_FILE", "WHITEROSE_REPO_PATH")

	v.SetConfigName("whiterose")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	cfg := &Config{
		Git: GitConfig{
			User: v.GetString("git.user"),
			Token: v.GetString("git.token"),
			Base: v.GetString("git.base"),
		},
		SSH: SSHConfig{
			KeyPath: v.GetString("ssh.keyPath"),
			KeyName: v.GetString("ssh.keyName"),
		},
		Image: ImageConfig{
			Name:    v.GetString("image.name"),
			Version: v.GetString("image.version"),
		},
		Repo: RepoConfig{
			Path: v.GetString("repo.path"),
		},
	}

	return cfg, nil
}

func LoadOrDefault() *Config {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{
			Git: GitConfig{Base: "main"},
		}
	}
	return cfg
}

func GetConfigPath() string {
	if path := os.Getenv("CONFIG_FILE"); path != "" {
		return path
	}
	if path := viper.GetString("repo.path"); path != "" {
		return path
	}
	return ".config.json"
}