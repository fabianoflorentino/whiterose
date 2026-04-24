# Architecture

## Project Structure

```
whiterose/
├── cmd/                      # CLI commands (Cobra)
│   ├── root.go              # Root command
│   ├── setup.go             # Setup command
│   ├── preReq.go            # Pre-requisites command
│   ├── docker.go            # Docker command
│   └── update.go            # Update command
├── internal/                # Clean Architecture layers
│   ├── interfaces/          # SOLID interfaces
│   ├── services/            # Service implementations
│   ├── domain/
│   │   ├── entities/        # Domain entities
│   │   └── application/    # Application services
│   └── infrastructure/
├── git/                     # Git operations
├── docker/                  # Docker operations
├── prereq/                  # Pre-requisites validation
├── setup/                   # Setup orchestration
├── update/                  # Update operations
├── utils/                   # Utilities
├── .github/workflows/       # CI/CD
├── Makefile                # Build targets
└── Dockerfile              # Container image
```

## Component Diagram

```mermaid
graph TB
    subgraph CLI
        ROOT[Root Command]
        SETUP[Setup Command]
        PREREQ[Pre-Req Command]
        DOCKER[Docker Command]
        UPDATE[Update Command]
    end

    subgraph Services
        SETUP_SVC[Setup Service]
        PREREQ_SVC[Prereq Service]
        DOCKER_SVC[Docker Service]
        UPDATE_SVC[Update Service]
    end

    subgraph Infrastructure
        GIT[Git Module]
        DOCKER_INFRA[Docker Module]
        UTILS[Utils Module]
    end

    ROOT --> SETUP
    ROOT --> PREREQ
    ROOT --> DOCKER
    ROOT --> UPDATE

    SETUP --> SETUP_SVC
    PREREQ --> PREREQ_SVC
    DOCKER --> DOCKER_SVC
    UPDATE --> UPDATE_SVC

    SETUP_SVC --> GIT
    SETUP_SVC --> PREREQ_SVC
    PREREQ_SVC --> UTILS
    DOCKER_SVC --> DOCKER_INFRA
    UPDATE_SVC --> GIT
    UPDATE_SVC --> UTILS
```

## Layer Responsibilities

| Layer | Purpose |
|-------|---------|
| `cmd/` | CLI entry points, Cobra commands |
| `internal/services/` | Business logic, orchestration |
| `internal/interfaces/` | Contracts for dependency injection |
| `internal/domain/` | Domain entities, business rules |
| `git/`, `docker/`, `utils/` | Infrastructure adapters |