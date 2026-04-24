# Execution Flow

## Command Execution Flow

### Main Entry Point

```mermaid
flowchart TD
    START([main.go]) --> LOAD_ENV[Load Environment]
    LOAD_ENV --> EXECUTE[Execute Root Command]
    EXECUTE --> MATCH{Match Command}
    MATCH -->|setup| CMD_SETUP[Setup Command]
    MATCH -->|pre-req| CMD_PREREQ[Pre-Req Command]
    MATCH -->|docker| CMD_DOCKER[Docker Command]
    MATCH -->|update| CMD_UPDATE[Update Command]
    MATCH -->|completion| CMD_COMPLETION[Completion Command]
    MATCH -->|help| CMD_HELP[Show Help]
```

### Setup Command Flow

```mermaid
flowchart TD
    START([setup --all]) --> CHECK_ALL{Check --all flag}
    CHECK_ALL -->|yes| PREREQ[Validate Prerequisites]
    CHECK_ALL -->|no| CHECK_PREREQ{Check --pre-req}
    CHECK_ALL -->|no| CHECK_REPOS{Check --repos}

    PREREQ --> LOAD_CONFIG[Load Config]
    PREREQ --> VALIDATE_APPS[Validate Apps]

    CHECK_PREREQ -->|yes| PREREQ
    CHECK_PREREQ -->|no| EXIT{Exit}

    CHECK_REPOS -->|yes| LOAD_CONFIG
    CHECK_REPOS -->|no| EXIT

    LOAD_CONFIG --> FETCH_REPOS[Fetch Repositories]
    FETCH_REPOS --> CLONE[Clone Repositories]
    CLONE --> SETUP_BRANCH[Setup Branch]
    SETUP_BRANCH --> EXIT
```

### Update Command Flow

```mermaid
flowchart TD
    START([update --config]) --> LOAD_CONFIG[Load Update Config]
    LOAD_CONFIG --> PARSE_PROJECTS[Parse Projects]

    PARSE_PROJECTS --> FOR_EACH{For each project}
    FOR_EACH --> CHECK_FLAGS{Check Flags}
    CHECK_FLAGS -->|--list| LIST_UPDATES[List Updates]
    CHECK_FLAGS -->|--packages| UPDATE_PACKAGES[Update Packages]
    CHECK_FLAGS -->|--go-version| UPDATE_GO[Update Go Version]
    CHECK_FLAGS -->|--docker-image| UPDATE_DOCKER[Update Docker Image]
    CHECK_FLAGS -->|--report| GENERATE_REPORT[Generate Report]

    LIST_UPDATES --> CREATE_PR{Create PR?}
    UPDATE_PACKAGES --> CREATE_PR
    UPDATE_GO --> CREATE_PR
    UPDATE_DOCKER --> CREATE_PR
    GENERATE_REPORT --> CREATE_PR

    CREATE_PR -->|yes| CREATE_BRANCH[Create Branch]
    CREATE_PR -->|no| COMMIT[Commit Changes]

    CREATE_BRANCH --> COMMIT
    COMMIT --> PUSH[Push to Remote]
    PUSH --> CREATE_PR_PULL{Create PR?}
    CREATE_PR_PULL -->|yes| PULL_REQUEST[Create Pull Request]
    CREATE_PR_PULL -->|no| END([Done])
    PULL_REQUEST --> END
```

### Pre-Req Command Flow

```mermaid
flowchart TD
    START([pre-req]) --> VALIDATE[New AppValidator]
    VALIDATE --> CHECK_FLAGS{Check Flags}
    CHECK_FLAGS -->|check| VALIDATE_ALL[Validate All Apps]
    CHECK_FLAGS -->|list| LIST_ALL[List All Apps]
    CHECK_FLAGS -->|apps| VALIDATE_SPECIFIC[Validate Specific Apps]
    CHECK_FLAGS -->|none| SHOW_HELP[Show Help]

    VALIDATE_ALL --> FOR_EACH{For each app}
    FOR_EACH --> CHECK_INSTALLED[Check Installed?]
    CHECK_INSTALLED --> SHOW_VERSION[Show Version]
    SHOW_VERSION --> CHECK_VERSION[Version OK?]
    CHECK_VERSION -->|yes| END_OK[Installed]
    CHECK_VERSION -->|no| SHOW_INSTALL[Show Install Instructions]
    SHOW_INSTALL --> END

    LIST_ALL --> PRINT_LIST[Print App List]
    PRINT_LIST --> END

    VALIDATE_SPECIFIC --> VALIDATE_ONE[Validate One]
    VALIDATE_ONE --> END
```

## Configuration Flow

```mermaid
flowchart LR
    CONFIG[Config File<br/>.config.json/.yaml]
    ENV[Environment<br/>.env]
    ARGS[Command Args]

    subgraph Load Sequence
        ENV --> LOAD_ENV[Load .env]
        CONFIG --> LOAD_CONFIG[Load Config]
        ARGS --> PARSE_ARGS[Parse Args]
    end

    LOAD_ENV --> MERGE[Merge Configs]
    LOAD_CONFIG --> MERGE
    PARSE_ARGS --> MERGE
```

## GitHub Actions Workflow

```mermaid
flowchart TD
    START([Schedule/Dispatch]) --> PREPARE[Prepare Job]
    PREPARE --> VALIDATE[Validate Version]
    VALIDATE --> CHECK_TAG{Check Tag Exists?}

    CHECK_TAG -->|no| TAG[Create & Push Tag]
    CHECK_TAG -->|yes| ERROR[Error: Tag exists]

    TAG --> BUILD[Build Matrix]
    BUILD --> BUILD_LINUX[Build Linux]
    BUILD --> BUILD_MAC[Build macOS]
    BUILD --> BUILD_WINDOWS[Build Windows]

    BUILD_LINUX --> UPLOAD[Upload Artifacts]
    BUILD_MAC --> UPLOAD
    BUILD_WINDOWS --> UPLOAD

    UPLOAD --> RELEASE[Create Release]
    RELEASE --> CHECKSUMS[Generate Checksums]
    CHECKSUMS --> CREATE_NOTES[Create Release Notes]
    CREATE_NOTES --> PUBLISH[Publish Release]

    TAG --> DOCKER[Build Docker Image]
    DOCKER --> PUSH_DOCKER[Push to Registries]
    PUSH_DOCKER --> PUBLISH
```
