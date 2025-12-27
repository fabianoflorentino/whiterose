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

- Go 1.25+
- SSH keys (for cloning via SSH)
- `.env` file in your home directory (required, for environment configuration)
- `.config.json` file in your home directory (required, for repository configuration and applications validation)
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
      "recommendedVersion": "1.25.0",
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
    recommendedVersion: "1.25.0"
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
- `completion` &mdash; Generate shell autocompletion scripts

Use `whiterose [command] --help` for more information about each command and its flags.

## Environment Variables

Create a `.env` file in your home directory (required):

- `GIT_USER`: Git username for HTTPS authentication
- `GIT_TOKEN`: Git token/password for HTTPS authentication
- `CONFIG_FILE`: Path to the configuration file (default: `$HOME/.config.json`)
- `SSH_KEY_PATH`: Path to your SSH key directory (default: `$HOME/.ssh`)
- `SSH_KEY_NAME`: Name of your SSH private key (default: `id_rsa`)
- `IMAGE_NAME`: Name of the Docker image to build (default: `my_app:latest`)
- `IMAGE_VERSION`: Version of the Docker image to build (default: `latest`)
- `DOCKERFILE_PATH`: Path to the Dockerfile (default: `$PWD/Dockerfile`)

If not set, default values are used.

## Project Structure

- `main.go`: Entry point, loads environment and executes commands
- `cmd/`: CLI commands (`setup`, `pre-req`, `docker`)
- `git/`: Git operations (clone, checkout)
- `prereq/`: Environment validation utilities
- `setup/`: Setup logic for cloning repositories
- `docker/`: Docker-related utilities
- `utils/`: Helpers for environment variables and JSON parsing
- `config.json` or `config.yaml`: List of repositories and required applications

```sh
.
├── cmd/
├── docker/
├── git/
├── prereq/
├── setup/
├── utils/
├── config.json
├── Dockerfile
├── main.go
├── README.md
└── ...
```

## License

This project is licensed under the MIT License.
