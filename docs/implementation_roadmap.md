# Slurm REST CLI Implementation Roadmap

This document outlines the detailed implementation plan for the Slurm REST CLI tool (`srest`). The implementation is divided into phases to ensure a structured and manageable development process.

## Phase 1: Core Framework and Basic Commands

### 1.1 Project Setup and Structure
- [x] Initialize Go project with modules
- [x] Set up directory structure (cmd, internal, pkg)
- [x] Add essential dependencies (Cobra, Viper, etc.)

### 1.2 Configuration Management
- [x] Implement configuration file handling with Viper
- [x] Create configuration structure and YAML serialization
- [x] Implement environment variable support
- [x] Add config command for managing configuration

### 1.3 Authentication Framework
- [x] Design authentication interfaces for different auth methods
- [x] Implement token-based authentication
- [x] Implement JWT-based authentication
- [x] Implement user+token authentication
- [x] Create auth login command for Keycloak integration

### 1.4 API Client
- [x] Design API client interface
- [x] Implement HTTP request/response handling
- [x] Add authentication header support
- [x] Implement curl command generation
- [ ] Add error handling and retries

### 1.5 Output Formatting
- [x] Design formatter interfaces
- [x] Implement JSON formatter
- [x] Implement YAML formatter
- [x] Implement table formatter
- [ ] Add support for custom formats

### 1.6 Basic Commands
- [x] Implement root command with global flags
- [x] Implement job command structure
- [x] Implement node command structure
- [ ] Connect commands to API client
- [ ] Add proper error handling and validation

## Phase 2: Job Management Commands

### 2.1 Job Submission
- [ ] Implement job submit command with script input
- [ ] Add support for all job parameters
- [ ] Implement job script validation
- [ ] Add support for job templates

### 2.2 Job Query and Control
- [ ] Implement job list command with filters
- [ ] Implement job show command with detailed output
- [ ] Implement job cancel command
- [ ] Add support for job array operations

### 2.3 Job Monitoring
- [ ] Implement job wait command
- [ ] Add support for job notifications
- [ ] Implement job stats command

## Phase 3: Node and Partition Management

### 3.1 Node Management
- [ ] Implement node list command with filters
- [ ] Implement node show command with detailed output
- [ ] Implement node update command
- [ ] Add support for node state changes

### 3.2 Partition Management
- [ ] Implement partition list command
- [ ] Implement partition show command
- [ ] Implement partition update command

## Phase 4: Advanced Features

### 4.1 Reservation Management
- [ ] Implement reservation list command
- [ ] Implement reservation create command
- [ ] Implement reservation update command
- [ ] Implement reservation delete command

### 4.2 QOS Management
- [ ] Implement QOS list command
- [ ] Implement QOS show command

### 4.3 Advanced Authentication
- [ ] Implement token refresh
- [ ] Add support for multiple authentication profiles
- [ ] Implement secure credential storage

### 4.4 Advanced Output
- [ ] Add support for custom output templates
- [ ] Implement output filtering
- [ ] Add support for output redirection

## Phase 5: Testing and Documentation

### 5.1 Unit Testing
- [ ] Add unit tests for API client
- [ ] Add unit tests for formatters
- [ ] Add unit tests for commands
- [ ] Add unit tests for authentication

### 5.2 Integration Testing
- [ ] Set up integration test environment
- [ ] Add integration tests for job commands
- [ ] Add integration tests for node commands
- [ ] Add integration tests for authentication

### 5.3 Documentation
- [ ] Create user documentation
- [ ] Create developer documentation
- [ ] Add examples and tutorials
- [ ] Create man pages

## Phase 6: Deployment and Distribution

### 6.1 Packaging
- [ ] Create release process
- [ ] Add version information
- [ ] Create installation packages

### 6.2 CI/CD
- [ ] Set up continuous integration
- [ ] Set up continuous deployment
- [ ] Add automated testing

## Minimum Viable Product (MVP)

The MVP for the Slurm REST CLI tool will include:

1. **Core Framework**
   - Configuration management
   - Authentication with JWT and token
   - Basic API client
   - Output formatting (JSON, YAML, table)

2. **Job Commands**
   - `job submit`: Submit a job script
   - `job list`: List jobs with basic filtering
   - `job show`: Show job details
   - `job cancel`: Cancel a job

3. **Node Commands**
   - `node list`: List nodes with basic filtering
   - `node show`: Show node details

4. **Authentication**
   - `auth login`: Authenticate with Keycloak
   - Support for token and JWT authentication

5. **Configuration**
   - `config get`: Get configuration values
   - `config set`: Set configuration values
   - `config list`: List all configuration values

The MVP will focus on providing the essential functionality needed for basic Slurm cluster interaction via the REST API, with a clean and user-friendly interface.
