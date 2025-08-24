# Whiterose

Whiterose is a CLI tool written in Go for automating the cloning and setup of multiple Git repositories. It is designed to streamline the process of preparing development environments, especially for teams working with several repositories.

## Features

- Clone repositories using HTTPS or SSH
- Automatically checks out the `development` branch if available
- Creates and checks out a user-specific branch if `development` does not exist
- Loads environment variables from a `.env` file
- Configurable via a JSON file listing repositories

## Requirements

- Go 1.18+
- SSH keys (for cloning via SSH)
- `.env` file in your home directory with credentials (optional)

## Installation

Clone this repository and build the CLI:

```sh
git clone git@github.com:fabianoflorentino/whiterose.git
cd whiterose
go build -o whiterose main.go
```

## Usage

Configure your repositories in `repos.json`:

```json
[
  {
    "url": "git@github.com:fabianoflorentino/mr-robot.git",
    "directory": "mr_robot"
  }
]
```

Run the setup command to clone all repositories:

```sh
./whiterose setup
```

## Environment Variables

Whiterose can use the following environment variables (set in your `.env` file):

- `GIT_USER`: Git username for HTTPS authentication
- `GIT_TOKEN`: Git token/password for HTTPS authentication
- `SSH_KEY_PATH`: Path to your SSH key directory
- `SSH_KEY_NAME`: Name of your SSH private key (default: `id_rsa`)

## License

This project is licensed under the MIT License.
