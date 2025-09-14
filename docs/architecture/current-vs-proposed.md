# 🏗️ Current vs. Proposed Architecture

## 📊 Current Architecture Analysis

### Current Structure

```bash
whiterose/
├── main.go                 # Entry point with mixed logic
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command
│   ├── setup.go           # Setup command
│   ├── preReq.go          # Prereq command
│   └── docker.go          # Docker command
├── git/                    # Git operations directly coupled
│   └── git.go
├── docker/                 # Docker operations directly coupled
│   └── docker.go
├── prereq/                 # Prerequisites validation
│   └── pre_req.go
├── setup/                  # Setup orchestration
│   └── setup.go
└── utils/                  # Various utilities
    ├── get_env_or_default.go
    ├── json.go
    └── load_env_config.go
```

### Identified Problems

#### 1. **High Coupling** 🔗

- CLI commands directly coupled to implementations
- Direct dependencies between layers
- Difficulty replacing implementations

#### 2. **Mixed Responsibilities** 🔄

- Business logic mixed with infrastructure
- Configuration scattered throughout code
- Validation mixed with execution

#### 3. **Low Testability** 🧪

- External dependencies difficult to mock
- Tests too integrated
- Limited coverage of business code

#### 4. **Limited Extensibility** 📈

- Adding features requires modification in multiple points
- No clear separation of responsibilities
- Difficulty adding new adapters

## 🎯 Proposed Architecture (Hexagonal/Clean)

### New Structure

```bash
whiterose/
├── cmd/                           # 🔌 Adapters - CLI Interface
│   ├── handlers/                  # Command handlers
│   │   ├── setup_handler.go
│   │   ├── prereq_handler.go
│   │   └── docker_handler.go
│   └── root.go
├── internal/                      # 🏛️ Application Domain
│   ├── domain/                    # 💎 Domain Layer
│   │   ├── entities/              # Business entities
│   │   │   ├── repository.go
│   │   │   ├── application.go
│   │   │   ├── docker_image.go
│   │   │   └── environment.go
│   │   ├── repositories/          # Repository interfaces (ports)
│   │   │   ├── git_repository.go
│   │   │   ├── config_repository.go
│   │   │   ├── docker_repository.go
│   │   │   └── environment_repository.go
│   │   ├── services/              # Domain services
│   │   │   ├── git_service.go
│   │   │   ├── validation_service.go
│   │   │   └── docker_service.go
│   │   └── errors/                # Domain errors
│   │       └── domain_errors.go
│   ├── application/               # 🔄 Application Layer
│   │   ├── usecases/              # Use cases (interactors)
│   │   │   ├── setup_repositories.go
│   │   │   ├── validate_prerequisites.go
│   │   │   ├── manage_docker_images.go
│   │   │   └── manage_environment.go
│   │   ├── ports/                 # Application ports
│   │   │   ├── input/             # Primary ports
│   │   │   │   ├── setup_service.go
│   │   │   │   ├── validation_service.go
│   │   │   │   └── docker_service.go
│   │   │   └── output/            # Secondary ports
│   │   │       ├── git_port.go
│   │   │       ├── config_port.go
│   │   │       ├── docker_port.go
│   │   │       └── logger_port.go
│   │   └── dto/                   # Data Transfer Objects
│   │       ├── setup_dto.go
│   │       ├── validation_dto.go
│   │       └── docker_dto.go
│   └── infrastructure/            # 🔌 Infrastructure Layer
│       ├── adapters/              # Secondary adapters
│       │   ├── git/               # Git adapter
│       │   │   ├── go_git_adapter.go
│       │   │   └── git_auth.go
│       │   ├── config/            # Configuration adapter
│       │   │   ├── json_config.go
│       │   │   └── env_config.go
│       │   ├── docker/            # Docker adapter
│       │   │   └── docker_client.go
│       │   ├── validation/        # Validation adapter
│       │   │   └── system_validator.go
│       │   └── logging/           # Logging adapter
│       │       └── logger.go
│       ├── external/              # External dependencies
│       │   ├── filesystem/
│       │   ├── network/
│       │   └── system/
│       └── config/                # Infrastructure config
│           ├── container.go       # DI container
│           └── wire.go           # Dependency injection
├── pkg/                          # 📦 Shared packages
│   ├── logger/                   # Logging utilities
│   ├── errors/                   # Error handling
│   └── validation/               # Common validations
├── configs/                      # ⚙️ Configuration files
│   ├── app.yaml
│   └── environments/
└── main.go                       # 🚀 Application entry point
```

### Improvements of the New Architecture

#### 1. **Low Coupling** 🔗

- Well-defined interfaces between layers
- Dependency Injection for flexibility
- Easy to replace implementations

#### 2. **Clear Separation of Responsibilities** 🎯

- **Domain**: Pure business rules
- **Application**: Use case orchestration
- **Infrastructure**: Technical details and external systems

#### 3. **High Testability** 🧪

- Interfaces allow easy mock creation
- Isolated unit tests per layer
- Complete coverage of business code

#### 4. **Extensibility** 📈

- New use cases without modifying core
- Multiple adapters for different sources
- Support for different interfaces (CLI, API, GUI)

## 🔄 Migration Mapping

### Current Architecture → New Architecture

| Current File | New Location | Responsibility |
|--------------|--------------|----------------|
| `cmd/setup.go` | `cmd/handlers/setup_handler.go` | CLI Interface |
| `setup/setup.go` | `internal/application/usecases/setup_repositories.go` | Use case |
| `git/git.go` | `internal/infrastructure/adapters/git/go_git_adapter.go` | Git Adapter |
| `utils/json.go` | `internal/infrastructure/adapters/config/json_config.go` | Config Adapter |
| `prereq/pre_req.go` | `internal/application/usecases/validate_prerequisites.go` | Use case |
| `docker/docker.go` | `internal/infrastructure/adapters/docker/docker_client.go` | Docker Adapter |

### New Components

| Component | Location | Responsibility |
|-----------|----------|----------------|
| Entities | `internal/domain/entities/` | Business models |
| Ports | `internal/application/ports/` | Communication interfaces |
| DTOs | `internal/application/dto/` | Data transfer |
| Domain Services | `internal/domain/services/` | Complex business rules |
| DI Container | `internal/infrastructure/config/container.go` | Dependency injection |

## 📈 Quantifiable Benefits

### Before Migration

- **Testability**: ~30% possible coverage
- **Coupling**: High (cascading modifications)
- **Time for new feature**: 2-3 days
- **Test complexity**: High

### After Migration

- **Testability**: ~90% possible coverage
- **Coupling**: Low (isolated modifications)
- **Time for new feature**: 4-6 hours
- **Test complexity**: Low

## 🚀 Next Step

Continue to the [Execution Plan](../migration/execution-plan.md) to see how to implement this migration step by step.
