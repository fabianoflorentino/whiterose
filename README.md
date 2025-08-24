# Whiterose

Whiterose is a CLI tool written in Go to automate the setup of multiple Git repositories and validate environment prerequisites. It streamlines onboarding and environment preparation for development teams.

## Features

- Clone repositories using HTTPS or SSH
- Automatically checks out the `development` branch if available, or creates a branch `development/<user>`
- Loads environment variables from a `.env` file in your home directory
- Validates required applications (Go, Git, Docker, jq, yq) and shows installation instructions
- Configurable via a JSON file (`repos.json`) listing repositories
- Extensible via flags and environment variables

## Requirements

- Go 1.25+
- SSH keys (for cloning via SSH)
- `.env` file in your home directory (optional, for credentials)
- Linux, macOS, or Windows

## Installation

Clone this repository and build the CLI:

```sh
git clone git@github.com:fabianoflorentino/whiterose.git
cd whiterose
go build -o whiterose main.go
```

## Usage

### Configure repositories

Edit `repos.json`:

```json
[
  {
    "url": "git@github.com:fabianoflorentino/mr-robot.git",
    "directory": "mr_robot"
  }
]
```

### Clone repositories

```sh
./whiterose setup
```

### Validate environment prerequisites

```sh
./whiterose pre-req --check
```

List all available applications for validation:

```sh
./whiterose pre-req --list
```

Validate specific applications:

```sh
./whiterose pre-req --apps go,git
```

## Environment Variables

Create a `.env` file in your home directory (required):

- `GIT_USER`: Git username for HTTPS authentication
- `GIT_TOKEN`: Git token/password for HTTPS authentication
- `SSH_KEY_PATH`: Path to your SSH key directory (default: `$HOME/.ssh`)
- `SSH_KEY_NAME`: Name of your SSH private key (default: `id_rsa`)

If not set, default values are used.

## Project Structure

- `main.go`: Entry point, loads environment and executes commands
- `cmd/`: CLI commands (`setup`, `pre-req`)
- `git/`: Git operations (clone, checkout)
- `prereq/`: Environment validation utilities
- `setup/`: Setup logic for cloning repositories
- `utils/`: Helpers for environment variables and JSON parsing
- `repos.json`: List of repositories to clone

## License

This project is licensed under the MIT License.
