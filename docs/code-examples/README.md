# Code Examples - Hexagonal Architecture

This directory contains complete code examples for implementing the hexagonal architecture migration for the WhiteRose project.

## 📁 Structure

```text
code-examples/
├── domain/                    # Domain Layer Examples
│   ├── entities/             # Business entities
│   ├── repositories/         # Repository interfaces (ports)
│   ├── services/             # Domain services
│   └── errors/              # Domain-specific errors
├── application/              # Application Layer Examples
│   ├── usecases/            # Use case implementations
│   ├── dtos/                # Data Transfer Objects
│   └── ports/               # Primary and secondary ports
└── infrastructure/          # Infrastructure Layer Examples
    ├── adapters/            # External system adapters
    ├── config/              # Configuration management
    └── di/                  # Dependency injection
```

## 🎯 How to Use

1. **Study the examples** to understand the hexagonal architecture patterns
2. **Copy and adapt** the code for your specific use cases
3. **Follow the naming conventions** established in the examples
4. **Maintain the layer separation** as demonstrated
5. **Use the error handling patterns** shown in the examples

## 📋 Examples Included

### Domain Layer

- Repository entity with validation
- Git repository interface
- Docker image entity
- Environment configuration entity
- Domain services for business rules
- Custom error types

### Application Layer

- Setup repositories use case
- Validate prerequisites use case
- Manage Docker images use case
- DTOs for data transfer
- Port definitions

### Infrastructure Layer

- Git adapter using go-git
- Docker adapter using Docker SDK
- Configuration adapter
- Validation adapter
- Dependency injection container

## ⚠️ Important Notes

- These examples are **templates** - adapt them to your specific needs
- Follow the **dependency rule**: inner layers should not depend on outer layers
- Use **interfaces** to maintain loose coupling between layers
- Implement **proper error handling** at each layer
- Add **comprehensive tests** for each component

## 🔗 Related Documentation

- [Migration Execution Plan](../migration/execution-plan.md)
- [Architecture Overview](../architecture/current-vs-proposed.md)
- [Phase 1: Domain Layer Guide](../migration/phases/phase-1-domain.md)
- [Testing Guide](../how-to/testing-guide.md)
