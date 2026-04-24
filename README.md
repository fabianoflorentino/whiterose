# Whiterose

Whiterose is a CLI tool written in Go that automates the setup of multiple Git repositories and validates environment prerequisites. It streamlines onboarding and environment preparation for development teams.

## Features

- Clone repositories using HTTPS or SSH
- Automatically checks out the `development` branch if available, or creates a branch `development/<user>`
- Loads environment variables from a `.env` file in your home directory
- Validates required applications (Go, Git, Docker, jq, yq) and shows installation instructions
- Configurable via a JSON (`.json`) or YAML (`.yaml`/`.yml`) file (default: `$HOME/.config.json`) listing repositories and required applications
- Extensible via command-line flags and environment variables
- Docker automation commands (check, build, list, delete images)

## Requirements

- Go 1.26+
- SSH keys (for cloning via SSH)
- `.env` file in your home directory (required, for environment configuration)
- `.config.json` file in your home directory (required, for repository configuration and applications)
- Linux, macOS, or Windows

## Installation

Clone this repository and build the CLI:

```sh
git clone git@github.com:fabianoflorentino/whiterose.git
cd whiterose
go build -o whiterose main.go
```

### Option 1: Download a Pre-built Release

Go to the [Releases page](https://github.com/fabianoflorentino/whiterose/releases) and download the latest binary for your operating system:

- **Linux**: Download `whiterose-linux-amd64.tar.gz`, extract, and move the binary to a directory in your `$PATH`:
  
  ```sh
  wget https://github.com/fabianoflorentino/whiterose/releases/latest/download/whiterose-linux-amd64.tar.gz
  tar -xzf whiterose-linux-amd64.tar.gz
  sudo mv whiterose /usr/local/bin/
  ```

- **macOS**: Download `whiterose-darwin-amd64.tar.gz`, extract, and move the binary:

  ```sh
  curl -LO https://github.com/fabianoflorentino/whiterose/releases/latest/download/whiterose-darwin-amd64.tar.gz
  tar -xzf whiterose-darwin-amd64.tar.gz
  sudo mv whiterose /usr/local/bin/
  ```

- **Windows**: Download `whiterose-windows-amd64.zip` and extract the `whiterose.exe` file. Add its folder to your `PATH` or run it directly.

### Option 2: Build from Source

```sh
git clone git@github.com:fabianoflorentino/whiterose.git
cd whiterose
go build -o whiterose main.go
```

## Usage

### Configure repositories

Edit your configuration file: `.config.json` **or** `.config.yaml` (YAML is also supported).

- Add your Git repository URLs and directory names.
- Add any required applications for your projects.

#### Example (JSON)

```json
{
  "repositories": [
    {
      "url": "git@github.com:fabianoflorentino/mr-robot.git",
      "directory": "mr_robot"
    }
  ],
  "applications": [
    {
      "name": "Go",
      "command": "go",
      "versionFlag": "version",
      "recommendedVersion": "1.26.0",
      "installInstructions": {
        "linux":   "sudo apt install golang",
        "darwin":  "brew install go",
        "windows": "choco install golang"
      }
    }
  ]
}
```

#### Example (YAML)

```yaml
repositories:
  - url: git@github.com:fabianoflorentino/mr-robot.git
    directory: mr_robot
applications:
  - name: Go
    command: go
    versionFlag: version
    recommendedVersion: "1.26.0"
    installInstructions:
      linux: sudo apt install golang
      darwin: brew install go
      windows: choco install golang
```

### Run setup

```sh
./whiterose setup [flags]
```

This will clone all repositories listed in your config file and check out the appropriate branches.

### Main Commands

- `setup` &mdash; Clone and set up repositories
  - Flags:
    - `--all, -a` &mdash; Check prerequisites and clone repositories
    - `--pre-req, -p` &mdash; Only check prerequisites
    - `--repos, -r` &mdash; Only clone repositories
- `pre-req` &mdash; Validate and list required applications
  - Flags:
    - `--check, -c` &mdash; Check if all required applications are installed
    - `--list, -l` &mdash; List all available applications
    - `--apps, -a` &mdash; Validate specific applications (comma-separated)
- `docker` &mdash; Automate Docker operations (check/build/list/delete images)
  - Flags:
    - `--file, -f` &mdash; Check if Dockerfile exists
    - `--build, -b` &mdash; Build Docker image
    - `--list, -l` &mdash; List Docker images
    - `--delete, -d` &mdash; Delete Docker image
- `update` &mdash; Update dependencies and versions (Go, Docker)
  - Flags:
    - `--go-mod, -g` &mdash; Update go.mod dependencies
    - `--go-version, -v` &mdash; Update Go version in go.mod
    - `--docker-image, -d` &mdash; Update base Docker image
    - `--list, -l` &mdash; List available updates
    - `--major, -m` &mdash; Update major version
    - `--dry-run, -n` &mdash; Show what would be updated
    - `--pr, -p` &mdash; Create pull request after update
    - `--base, -b` &mdash; Base branch for PR (default: main)
    - `--config, -c` &mdash; Path to update config file
- `completion` &mdash; Generate shell autocompletion scripts

Use `whiterose [command] --help` for more information about each command and its flags.

## Environment Variables

Whiterose uses environment variables and/or a config file for configuration. Priority: flags > env vars > config file > defaults.

| Variable | Description | Default |
|----------|-------------|---------|
| `GIT_USER` | Git username for HTTPS | - |
| `GIT_TOKEN` | Git token/password | - |
| `CONFIG_FILE` | Path to config file | `.config.json` |
| `SSH_KEY_PATH` | SSH key directory | `~/.ssh` |
| `SSH_KEY_NAME` | SSH key name | `id_rsa` |
| `IMAGE_NAME` | Docker image name | `my_app` |
| `IMAGE_VERSION` | Docker image version | `latest` |

Optional config file in `$HOME/.config/whiterose.yaml`:

```yaml
git:
  user: "your-user"
  token: "your-token"
  base: "main"

ssh:
  keyPath: "~/.ssh"
  keyName: "id_rsa"

image:
  name: "my_app"
  version: "latest"
```

## Project Structure

- `main.go`: Entry point, loads environment and executes commands
- `cmd/`: CLI commands (`setup`, `pre-req`, `docker`, `update`)
- `git/`: Git operations (clone, checkout)
- `prereq/`: Environment validation utilities
- `setup/`: Setup logic for cloning repositories
- `docker/`: Docker-related utilities
- `update/`: Version update utilities
- `utils/`: Helpers for environment variables and JSON parsing
- `config.json` or `config.yaml`: List of repositories and required applications
- `internal/`: Clean Architecture layers
  - `interfaces/`: SOLID interfaces
  - `services/`: Service implementations
  - `domain/`: Domain entities
  - `application/`: Application services

```sh
.
├── cmd/
├── docker/
├── git/
├── internal/
│   ├── application/
│   ├── domain/
│   ├── interfaces/
│   └── services/
├── prereq/
├── setup/
├── update/
├── utils/
├── config.json
├── Dockerfile
├── Makefile
├── main.go
├── README.md
└── ...
```

## Release

To create a new release:

```bash
# Via GitHub Actions (recommended)
# 1. Go to Actions tab
# 2. Select "Release and Deploy Binary"
# 3. Click "Run workflow"
# 4. Enter version tag (e.g., v1.0.0)

# Via command line
git tag v1.0.0
git push origin v1.0.0
```

The workflow will automatically:
1. Validate the version format
2. Create and push the git tag
3. Build binaries for all platforms (linux, darwin, windows)
4. Generate checksums
5. Create a GitHub Release
6. Build and push Docker images to Docker Hub and GHCR

## License

This project is licensed under the MIT License.
