# 🏗️ Hexagonal Architecture Diagrams

## 🎯 Architecture Overview

```mermaid
graph TB
    subgraph "🔌 Primary Adapters (Input)"
        CLI[CLI Interface]
        API[REST API]
        GUI[GUI Interface]
    end
    
    subgraph "🏛️ Application Core"
        subgraph "Application Layer"
            UC1[Setup Repositories UseCase]
            UC2[Validate Prerequisites UseCase]
            UC3[Manage Docker Images UseCase]
            
            IP1[Setup Service Port]
            IP2[Validation Service Port]
            IP3[Docker Service Port]
        end
        
        subgraph "Domain Layer"
            E1[Repository Entity]
            E2[Application Entity]
            E3[Docker Image Entity]
            E4[Environment Entity]
            
            DS1[Git Domain Service]
            DS2[Validation Domain Service]
            DS3[Docker Domain Service]
        end
    end
    
    subgraph "🔌 Secondary Adapters (Output)"
        GA[Git Adapter]
        CA[Config Adapter]
        DA[Docker Adapter]
        VA[Validation Adapter]
        LA[Logger Adapter]
    end
    
    subgraph "🌐 External Systems"
        GIT[Git Repositories]
        FS[File System]
        DOCKER[Docker Engine]
        SYS[System Commands]
    end
    
    CLI --> IP1
    CLI --> IP2
    CLI --> IP3
    
    IP1 --> UC1
    IP2 --> UC2
    IP3 --> UC3
    
    UC1 --> E1
    UC1 --> E4
    UC1 --> DS1
    
    UC2 --> E2
    UC2 --> DS2
    
    UC3 --> E3
    UC3 --> DS3
    
    UC1 --> GA
    UC1 --> CA
    UC1 --> LA
    
    UC2 --> VA
    UC2 --> CA
    UC2 --> LA
    
    UC3 --> DA
    UC3 --> CA
    UC3 --> LA
    
    GA --> GIT
    CA --> FS
    DA --> DOCKER
    VA --> SYS
    LA --> FS
```

## 🔄 Data Flow - Repository Setup

```mermaid
sequenceDiagram
    participant CLI as CLI Handler
    participant UC as Setup UseCase
    participant DS as Git Domain Service
    participant E as Repository Entity
    participant GA as Git Adapter
    participant CA as Config Adapter
    participant GIT as Git Repository
    
    CLI->>UC: Execute(SetupRequest)
    UC->>CA: LoadRepositories()
    CA-->>UC: []Repository Config
    
    loop For each repository
        UC->>E: NewRepository(url, directory)
        E-->>UC: Repository Entity
        
        UC->>DS: ValidateRepository(entity)
        DS-->>UC: Validation Result
        
        UC->>GA: Clone(repository)
        GA->>GIT: git clone
        GIT-->>GA: Repository Cloned
        GA-->>UC: Clone Result
        
        UC->>GA: CheckoutBranch(branch)
        GA->>GIT: git checkout
        GIT-->>GA: Branch Checked Out
        GA-->>UC: Checkout Result
    end
    
    UC-->>CLI: SetupResponse
```

## 🧪 Prerequisites Validation Flow

```mermaid
sequenceDiagram
    participant CLI as CLI Handler
    participant UC as Validation UseCase
    participant DS as Validation Domain Service
    participant E as Application Entity
    participant VA as Validation Adapter
    participant CA as Config Adapter
    participant SYS as System
    
    CLI->>UC: Execute(ValidationRequest)
    UC->>CA: LoadApplications()
    CA-->>UC: []Application Config
    
    loop For each application
        UC->>E: NewApplication(config)
        E-->>UC: Application Entity
        
        UC->>DS: ValidateApplication(entity)
        DS->>VA: CheckInstalled(command)
        VA->>SYS: Execute command
        SYS-->>VA: Command Result
        VA-->>DS: Installation Status
        DS-->>UC: Validation Result
    end
    
    UC-->>CLI: ValidationResponse
```

## 🐳 Docker Management Flow

```mermaid
sequenceDiagram
    participant CLI as CLI Handler
    participant UC as Docker UseCase
    participant DS as Docker Domain Service
    participant E as DockerImage Entity
    participant DA as Docker Adapter
    participant CA as Config Adapter
    participant DOCKER as Docker Engine
    
    CLI->>UC: Execute(DockerRequest)
    UC->>CA: LoadDockerConfig()
    CA-->>UC: Docker Config
    
    UC->>E: NewDockerImage(name, tag)
    E-->>UC: DockerImage Entity
    
    alt Build Image
        UC->>DS: BuildImage(entity)
        DS->>DA: Build(options)
        DA->>DOCKER: docker build
        DOCKER-->>DA: Build Result
        DA-->>DS: Build Status
        DS-->>UC: Build Result
    else List Images
        UC->>DA: ListImages(filter)
        DA->>DOCKER: docker images
        DOCKER-->>DA: Image List
        DA-->>UC: Images
    else Delete Image
        UC->>DA: DeleteImage(name)
        DA->>DOCKER: docker rmi
        DOCKER-->>DA: Delete Result
        DA-->>UC: Delete Status
    end
    
    UC-->>CLI: DockerResponse
```

## 🏛️ Architecture Layers

```mermaid
graph LR
    subgraph "Layers"
        subgraph "🔌 Infrastructure"
            IA[Input Adapters]
            OA[Output Adapters]
        end
        
        subgraph "🔄 Application"
            UC[Use Cases]
            IP[Input Ports]
            OP[Output Ports]
            DTO[DTOs]
        end
        
        subgraph "💎 Domain"
            E[Entities]
            DS[Domain Services]
            R[Repository Interfaces]
            DE[Domain Errors]
        end
    end
    
    subgraph "Dependencies"
        IA --> IP
        IP --> UC
        UC --> OP
        UC --> E
        UC --> DS
        OP --> OA
        DS --> R
        E --> DE
    end
```

## 🔌 Ports & Adapters Pattern

```mermaid
graph TB
    subgraph "Primary Ports (Input)"
        PS[Setup Service]
        PV[Validation Service]
        PD[Docker Service]
    end
    
    subgraph "Application Core"
        UC1[Setup UseCase]
        UC2[Validation UseCase]
        UC3[Docker UseCase]
    end
    
    subgraph "Secondary Ports (Output)"
        PG[Git Port]
        PC[Config Port]
        PDO[Docker Port]
        PL[Logger Port]
    end
    
    subgraph "Primary Adapters"
        CLI[CLI Adapter]
        WEB[Web Adapter]
    end
    
    subgraph "Secondary Adapters"
        GA[Git Adapter]
        CA[Config Adapter]
        DA[Docker Adapter]
        LA[Logger Adapter]
    end
    
    CLI --> PS
    CLI --> PV
    CLI --> PD
    
    WEB --> PS
    WEB --> PV
    WEB --> PD
    
    PS --> UC1
    PV --> UC2
    PD --> UC3
    
    UC1 --> PG
    UC1 --> PC
    UC1 --> PL
    
    UC2 --> PC
    UC2 --> PL
    
    UC3 --> PDO
    UC3 --> PC
    UC3 --> PL
    
    PG --> GA
    PC --> CA
    PDO --> DA
    PL --> LA
```

## 📦 Package Structure

```mermaid
graph TB
    subgraph "Package Structure"
        subgraph "cmd/"
            CH[CLI Handlers]
        end
        
        subgraph "internal/"
            subgraph "domain/"
                ENT[entities/]
                REPO[repositories/]
                SERV[services/]
                ERR[errors/]
            end
            
            subgraph "application/"
                UC[usecases/]
                subgraph "ports/"
                    IN[input/]
                    OUT[output/]
                end
                DTO[dto/]
            end
            
            subgraph "infrastructure/"
                subgraph "adapters/"
                    GIT[git/]
                    CONF[config/]
                    DOCK[docker/]
                    VALID[validation/]
                    LOG[logging/]
                end
                CONFIG[config/]
            end
        end
        
        subgraph "pkg/"
            LOGGER[logger/]
            ERRORS[errors/]
            VALID_PKG[validation/]
        end
    end
```

## 🔄 Dependency Flow

```mermaid
graph TD
    subgraph "Dependency Direction"
        CMD[cmd/]
        INFRA[infrastructure/]
        APP[application/]
        DOMAIN[domain/]
        PKG[pkg/]
    end
    
    CMD --> APP
    CMD --> INFRA
    INFRA --> APP
    APP --> DOMAIN
    APP --> PKG
    INFRA --> PKG
    
    style DOMAIN fill:#90EE90, color:#000000
    style APP fill:#87CEEB, color:#000000
    style INFRA fill:#DDA0DD, color:#000000
    style CMD fill:#F0E68C, color:#000000
    style PKG fill:#FFB6C1, color:#000000
```

## 🚀 Deployment Architecture

```mermaid
graph TB
    subgraph "Development Environment"
        DEV[Developer Machine]
        subgraph "Local Services"
            LGIT[Local Git]
            LDOCKER[Local Docker]
            LFS[Local Filesystem]
        end
    end
    
    subgraph "WhiteRose CLI"
        CLI[CLI Application]
        subgraph "Adapters"
            GA[Git Adapter]
            DA[Docker Adapter]
            CA[Config Adapter]
        end
    end
    
    subgraph "External Services"
        GITHUB[GitHub]
        GITLAB[GitLab]
        DOCKER_HUB[Docker Hub]
    end
    
    DEV --> CLI
    CLI --> GA
    CLI --> DA
    CLI --> CA
    
    GA --> LGIT
    GA --> GITHUB
    GA --> GITLAB
    
    DA --> LDOCKER
    DA --> DOCKER_HUB
    
    CA --> LFS
```

## 📊 Component Interaction Overview

```mermaid
graph LR
    subgraph "User Interface"
        U[User]
        CLI[CLI Commands]
    end
    
    subgraph "Application Logic"
        UC[Use Cases]
        PORTS[Ports]
        DOMAIN[Domain Entities]
    end
    
    subgraph "External World"
        GIT[Git Repos]
        DOCKER[Docker]
        CONFIG[Config Files]
        SYSTEM[System]
    end
    
    U --> CLI
    CLI --> UC
    UC --> PORTS
    UC --> DOMAIN
    PORTS --> GIT
    PORTS --> DOCKER
    PORTS --> CONFIG
    PORTS --> SYSTEM
```

These diagrams show how hexagonal architecture keeps the application core isolated from external details, allowing easy testing and maintenance. The dependency direction always points inward, ensuring the domain remains pure and independent of external frameworks.
