# `srest` CLI Implementation Plan

This document outlines the implementation approach for the `srest` command-line tool, which provides a modern interface to the Slurm REST API while maintaining familiar Slurm command syntax.

## Architecture Overview

The `srest` CLI tool will be implemented using a modular architecture with the following components:

1. **Command Parser**: Handles command-line arguments and maps them to API calls
2. **API Client**: Communicates with the Slurm REST API
3. **Authentication Manager**: Handles various authentication methods
4. **Output Formatter**: Formats API responses for display
5. **Configuration Manager**: Manages user preferences and settings

## Directory Structure

```
srest/
├── cmd/
│   └── srest/                  # Main executable
│       └── main.go             # Entry point
├── internal/
│   ├── api/                    # API client implementation
│   │   ├── client.go           # Base API client
│   │   ├── auth.go             # Authentication handling
│   │   ├── job.go              # Job-related API calls
│   │   ├── node.go             # Node-related API calls
│   │   └── ...                 # Other API modules
│   ├── cli/                    # CLI command implementation
│   │   ├── root.go             # Root command and global flags
│   │   ├── job.go              # Job-related commands
│   │   ├── node.go             # Node-related commands
│   │   └── ...                 # Other command modules
│   └── formatter/              # Output formatting
│       ├── formatter.go        # Base formatter
│       ├── json.go             # JSON formatter
│       ├── yaml.go             # YAML formatter
│       └── table.go            # Table formatter
├── pkg/
│   ├── config/                 # Configuration management
│   │   └── config.go           # Configuration handling
│   └── models/                 # Data models
│       ├── job.go              # Job-related models
│       ├── node.go             # Node-related models
│       └── ...                 # Other models
└── docs/                       # Documentation
```

## Implementation Approach

### 1. Command-Line Parser

We'll use a robust command-line parsing library like [Cobra](https://github.com/spf13/cobra) to implement the command structure. Cobra provides:

- Nested commands (commands with subcommands)
- Automatic help generation
- Flag support
- Shell completions
- Minimum unique prefix matching

Example structure:

```go
// Root command
var rootCmd = &cobra.Command{
    Use:   "srest",
    Short: "Modern CLI for Slurm REST API",
}

// Job command
var jobCmd = &cobra.Command{
    Use:   "job",
    Short: "Job-related operations",
}

// Job submit subcommand
var jobSubmitCmd = &cobra.Command{
    Use:   "submit [options] [script]",
    Short: "Submit a batch job",
    Run: func(cmd *cobra.Command, args []string) {
        // Handle job submission
    },
}
```

### 2. API Client

The API client will handle communication with the Slurm REST API:

- HTTP requests to API endpoints
- Authentication header management
- Request/response serialization
- Error handling

Example:

```go
type APIClient struct {
    BaseURL    string
    Token      string
    HTTPClient *http.Client
}

func (c *APIClient) SubmitJob(job *models.JobSubmitRequest) (*models.JobSubmitResponse, error) {
    // Make API call to submit job
}
```

### 3. Authentication

Support multiple authentication methods:

- Token-based authentication
- JWT authentication
- User+token authentication
- Environment variables for credentials
- Configuration file for persistent settings

Example:

```go
type AuthManager struct {
    Token     string
    JWT       string
    Username  string
}

func (a *AuthManager) AddAuthHeaders(req *http.Request) {
    if a.JWT != "" {
        req.Header.Set("Authorization", "Bearer "+a.JWT)
    } else if a.Token != "" && a.Username != "" {
        req.Header.Set("X-SLURM-USER-TOKEN", a.Token)
        req.Header.Set("X-SLURM-USER-NAME", a.Username)
    } else if a.Token != "" {
        req.Header.Set("X-SLURM-USER-TOKEN", a.Token)
    }
}
```

### 4. Output Formatting

Support multiple output formats:

- JSON (for machine parsing)
- YAML (for human readability)
- Table (for terminal display)
- Custom formats (similar to `-o` in traditional Slurm commands)

Example:

```go
type Formatter interface {
    Format(data interface{}) (string, error)
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(data interface{}) (string, error) {
    // Format data as JSON
}
```

### 5. Configuration Management

Manage user preferences:

- Default output format
- API endpoint URL
- Authentication credentials
- Other user preferences

Example:

```go
type Config struct {
    APIEndpoint string
    Token       string
    JWT         string
    Username    string
    Format      string
}

func LoadConfig() (*Config, error) {
    // Load configuration from file or environment
}
```

### 6. Command Implementation

Each command will:

1. Parse command-line arguments
2. Convert arguments to API request parameters
3. Make API call
4. Format and display response

Example for job submission:

```go
func submitJob(cmd *cobra.Command, args []string) error {
    // Parse flags
    jobName, _ := cmd.Flags().GetString("job-name")
    nodes, _ := cmd.Flags().GetInt("nodes")
    
    // Create job request
    jobReq := &models.JobSubmitRequest{
        Job: &models.JobDescMsg{
            Name:  jobName,
            Nodes: nodes,
            // Other parameters
        },
    }
    
    // Handle script file if provided
    if len(args) > 0 {
        script, err := ioutil.ReadFile(args[0])
        if err != nil {
            return err
        }
        jobReq.Job.Script = string(script)
    }
    
    // Make API call
    client := api.NewClient()
    resp, err := client.SubmitJob(jobReq)
    if err != nil {
        return err
    }
    
    // Format and display response
    formatter := formatter.NewFormatter(config.GetFormat())
    output, err := formatter.Format(resp)
    if err != nil {
        return err
    }
    
    fmt.Println(output)
    return nil
}
```

### 7. Curl Command Generation

For the `--curl` option, implement a function to generate equivalent curl commands:

```go
func generateCurlCommand(method, url string, headers map[string]string, body interface{}) (string, error) {
    // Generate curl command string
}
```

## Implementation Phases

1. **Phase 1: Core Framework**
   - Set up project structure
   - Implement command parser with Cobra
   - Create basic API client
   - Implement authentication
   - Implement output formatting

2. **Phase 2: Job Management**
   - Implement job submission
   - Implement job listing
   - Implement job control (cancel, update)
   - Implement job history

3. **Phase 3: Node and Partition Management**
   - Implement node listing and control
   - Implement partition listing

4. **Phase 4: Advanced Features**
   - Implement remaining commands
   - Add shell completion
   - Implement configuration file
   - Add validation rules

5. **Phase 5: Testing and Documentation**
   - Write unit tests
   - Write integration tests
   - Complete documentation
   - Create examples

## Testing Strategy

1. **Unit Tests**: Test individual components
2. **Integration Tests**: Test command execution with mock API
3. **End-to-End Tests**: Test against a real Slurm REST API

## Documentation

1. **Command Help**: Detailed help for each command
2. **Man Pages**: Generate man pages
3. **Examples**: Provide example commands
4. **README**: Overview and quick start

## Deployment

1. **Binary Distribution**: Compile for multiple platforms
2. **Package Management**: Create packages for different distributions
   - RPM for RHEL/CentOS
   - DEB for Debian/Ubuntu
   - Homebrew for macOS
   - Chocolatey for Windows
