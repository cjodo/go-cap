# go-cap Design Specification

## 1. Project Overview

**Project Name:** go-cap  
**Type:** Go library + CLI tool  
**Purpose:** A robust, production-ready Go client library and CLI for the REDCap (Research Electronic Data Capture) API.

### Goals

- Provide a type-safe, idiomatic Go client for interacting with REDCap projects
- Support all major REDCap API endpoints with proper error handling
- Include a CLI tool for common data export/import operations
- Follow current Go best practices (Go 1.21+)
- Enable concurrent operations with proper context support
- Ensure reliability through retries, timeouts, and rate limiting

### Non-Goals

- Database persistence layer (out of scope; users can use exported data as needed)
- GUI (CLI-only)
- REDCap administration operations (user management, project creation)

---

## 2. Project Structure

```
go-cap/
├── cmd/
│   └── cap/
│       └── main.go              # CLI entrypoint
├── pkg/
│   ├── redcap/
│   │   ├── client.go            # Client struct and constructor
│   │   ├── client_test.go
│   │   ├── option.go           # Functional options
│   │   ├── errors.go           # Custom error types
│   │   ├── retry.go            # Retry logic
│   │   ├── rate_limiter.go     # Rate limiting
│   │   ├── types.go            # Core types (Project, Record, Field, etc.)
│   │   ├── types_test.go
│   │   ├── export.go           # Export operations
│   │   ├── export_test.go
│   │   ├── import.go           # Import operations
│   │   ├── import_test.go
│   │   ├── metadata.go         # Metadata/data dictionary
│   │   ├── file.go             # File operations
│   │   ├── user.go             # User operations
│   │   ├── event.go            # Event operations (longitudinal)
│   │   ├── arm.go              # Arm operations
│   │   ├── dag.go              # Data access groups
│   │   ├── instrument.go       # Instrument/form operations
│   │   ├── project.go          # Project settings
│   │   └── version.go          # Version info
│   └── caplib/                 # Reusable helpers (optional, for external use)
├── internal/
│   └── cli/                    # CLI commands (private implementation)
│       ├── root.go
│       ├── export.go
│       ├── import.go
│       ├── forms.go
│       ├── metadata.go
│       └── config.go
├── configs/
│   └── config.yaml.example     # Example configuration
├── scripts/
│   └── build.sh                # Build script
├── Makefile
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .gitignore
```

### Directory Rationale

- **`cmd/cap/`**: CLI application entry point. Thin main that wires dependencies.
- **`pkg/redcap/`**: Public library code. Safe for external import. Contains all REDCap API client logic.
- **`internal/cli/`**: Private CLI command implementations. Not exposed as part of the public API.
- **`pkg/`**: Since this is both a library and CLI, the main `redcap` package belongs in `pkg/` for clean import paths (`github.com/yourname/go-cap/pkg/redcap`).

---

## 3. Core Types

### 3.1 Client Configuration

```go
// Option configures the Client.
type Option func(*Client) error

// Client represents a REDCap API client.
type Client struct {
    baseURL    string
    token      string
    httpClient *http.Client
    rateLimiter RateLimiter
    maxRetries int
    retryDelay time.Duration
}
```

**Key Configuration Options:**

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `WithTimeout` | `time.Duration` | HTTP request timeout | `30s` |
| `WithMaxRetries` | `int` | Maximum retry attempts | `3` |
| `WithRetryDelay` | `time.Duration` | Initial retry delay (exponential backoff) | `1s` |
| `WithRateLimiter` | `RateLimiter` | Custom rate limiter | `NewDefaultRateLimiter()` |
| `WithHTTPClient` | `*http.Client` | Custom HTTP client | `http.DefaultClient` |
| `WithLogLevel` | `string` | Log verbosity (`debug`, `info`, `warn`, `error`) | `"info"` |

### 3.2 Core Domain Types

```go
// Project represents a REDCap project.
type Project struct {
    ProjectID           int          `json:"project_id"`
    ProjectTitle        string       `json:"project_title"`
    CreationTime        time.Time    `json:"creation_time"`
    ProductionTime      time.Time    `json:"production_time"`
    ProjectLanguage     string       `json:"project_language"`
    Purpose             int          `json:"purpose"`
    PurposeOther        string       `json:"purpose_other"`
    ProjectNotes        string       `json:"project_notes"`
    CustomRecordLabel   string       `json:"custom_record_label"`
    SecondaryUniqueField string     `json:"secondary_unique_field"`
    IsLongitudinal      bool         `json:"is_longitudinal"`
    HasSurveys          bool         `json:"has_surveys"`
    HasRepeatingInstrumentsOrEvents bool `json:"has_repeating_instruments_or_events"`
    ExternalModules     []string     `json:"external_modules"`
    
    // Cached metadata (loaded lazily)
    mu       sync.RWMutex
    metadata []Field
    forms    map[string]*Form
    events   []Event
    arms     []Arm
    users    []User
}

// Field represents a REDCap data dictionary field.
type Field struct {
    FieldName                           string   `json:"field_name"`
    FieldLabel                         string   `json:"field_label"`
    FieldType                          string   `json:"field_type"`
    FormName                           string   `json:"form_name"`
    SectionHeader                      string   `json:"section_header"`
    FieldNote                          string   `json:"field_note"`
    TextValidationTypeOrShowSliderNumber string `json:"text_validation_type_or_show_slider_number"`
    TextValidationMin                   string   `json:"text_validation_min"`
    TextValidationMax                   string   `json:"text_validation_max"`
    Identifier                          string   `json:"identifier"`
    BranchingLogic                     string   `json:"branching_logic"`
    RequiredField                       string   `json:"required_field"`
    CustomAlignment                     string   `json:"custom_alignment"`
    QuestionNumber                      string   `json:"question_number"`
    MatrixGroupName                     string   `json:"matrix_group_name"`
    MatrixRanking                       string   `json:"matrix_ranking"`
    FieldAnnotation                     string   `json:"field_annotation"`
    
    // Choices for select/radio/checkbox fields
    Choices []FieldChoice `json:"-"`
}

// FieldChoice represents a choice in a multiple-choice field.
type FieldChoice struct {
    Code  string // e.g., "1", "0", "1,2"
    Label string
}

// Form represents a REDCap instrument/eCRF.
type Form struct {
    Name        string
    Label       string
    Fields      []Field // Ordered as in data dictionary
}

// Record represents a single REDCap record.
type Record struct {
    ID         string
    Fields     map[string]any // FieldName -> Value
    EventName  string
    Repetition FormRepetition
}

// FormRepetition represents repeating instrument/event info.
type FormRepetition struct {
    FormName      string
    CustomRecordLabel string
}

// Event represents a longitudinal event.
type Event struct {
    EventName         string `json:"event_name"`
    ArmNum            int    `json:"arm_num"`
    DayOffset         int    `json:"day_offset"`
    OffsetMin         int    `json:"offset_min"`
    OffsetMax         int    `json:"offset_max"`
    UniqueEventName   string `json:"unique_event_name"`
}

// Arm represents a study arm in longitudinal projects.
type Arm struct {
    ArmNum int    `json:"arm_num"`
    Name   string `json:"name"`
}

// User represents a REDCap user.
type User struct {
    Username         string `json:"username"`
    Email           string `json:"email"`
    FirstName       string `json:"firstname"`
    LastName        string `json:"lastname"`
    RoleID          int    `json:"role_id"`
    RoleLabel       string `json:"role_label"`
    DataAccessGroup string `json:"data_access_group"`
    Expiration      string `json:"expiration"`
    LastLogin       string `json:"last_login"`
    APIExport       bool   `json:"api_export"`
    APIImport       bool   `json:"api_import"`
    APIProject      bool   `json:"api_project"`
    MobileApp       bool   `json:"mobile_app"`
}
```

---

## 4. API Specification

### 4.1 Client Methods

#### Connection & Project

| Method | Description |
|--------|-------------|
| `NewClient(baseURL, token string, opts ...Option) (*Client, error)` | Create a new client |
| `Ping(ctx context.Context) error` | Verify API connectivity |
| `ExportProject(ctx context.Context) (*Project, error)` | Get project information |
| `ExportProjectSettings(ctx context.Context) (map[string]any, error)` | Get project settings |

#### Metadata / Data Dictionary

| Method | Description |
|--------|-------------|
| `ExportMetadata(ctx context.Context, opts ...MetadataOption) ([]Field, error)` | Export data dictionary |
| `ExportFieldNames(ctx context.Context) ([]string, error)` | Export field names as exported |
| `ExportInstrumentList(ctx context.Context) ([]Instrument, error)` | List all instruments |
| `ExportFormEventMapping(ctx context.Context) ([]FormEventMapping, error)` | Get form-event mappings |

#### Records (Export)

| Method | Description |
|--------|-------------|
| `ExportRecords(ctx context.Context, opts ...ExportOption) ([]Record, error)` | Export records |
| `ExportRecordsRaw(ctx context.Context, opts ...ExportOption) ([]byte, error)` | Export raw format (CSV/JSON) |
| `ExportRepeatingFormsEvents(ctx context.Context) ([]RepeatingForm, error)` | Get repeating form/event info |

#### Records (Import)

| Method | Description |
|--------|-------------|
| `ImportRecords(ctx context.Context, records []Record, opts ...ImportOption) (*ImportResult, error)` | Import records |
| `ImportRecordsRaw(ctx context.Context, data []byte, opts ...ImportOption) (*ImportResult, error)` | Import from raw format |
| `GenerateNextRecordName(ctx context.Context) (string, error)` | Generate next sequential record name |

#### Files

| Method | Description |
|--------|-------------|
| `ExportFile(ctx context.Context, recordID, field, event string) ([]byte, error)` | Export a file field |
| `ImportFile(ctx context.Context, recordID, field, event string, data []byte, opts ...ImportOption) error` | Import a file field |
| `DeleteFile(ctx context.Context, recordID, field, event string) error` | Delete a file field |

#### Users & Permissions

| Method | Description |
|--------|-------------|
| `ExportUsers(ctx context.Context) ([]User, error)` | Export project users |
| `ImportUsers(ctx context.Context, users []User) (*ImportResult, error)` | Import users |
| `ExportDAGs(ctx context.Context) ([]DAG, error)` | Export data access groups |
| `ImportDAGs(ctx context.Context, dags []DAG) (*ImportResult, error)` | Import DAGs |
| `AssignUserToDAG(ctx context.Context, username, dag string) error` | Assign user to DAG |

#### Longitudinal

| Method | Description |
|--------|-------------|
| `ExportEvents(ctx context.Context) ([]Event, error)` | Export events |
| `ImportEvents(ctx context.Context, events []Event) (*ImportResult, error)` | Import events |
| `ExportArms(ctx context.Context) ([]Arm, error)` | Export arms |
| `ImportArms(ctx context.Context, arms []Arm) (*ImportResult, error)` | Import arms |
| `ExportFormEventMapping(ctx context.Context) ([]FormEventMapping, error)` | Export form-event mappings |
| `ImportFormEventMapping(ctx context.Context, mappings []FormEventMapping) (*ImportResult, error)` | Import form-event mappings |

### 4.2 Option Types

#### Export Options

```go
type ExportOptions struct {
    Records                []string // Record IDs
    Fields                []string // Field names
    Forms                 []string // Form names
    Events                []string // Event names (longitudinal)
    RawOrLabel            string   // "raw" | "label" | "both"
    RawOrLabelHeaders    string   // "raw" | "label"
    ExportCheckboxLabel   bool     // Export checkbox values as labels
    ExportSurveyFields    bool     // Include survey fields
    ExportDataAccessGroups bool    // Include DAG field
    Format                string   // "json" | "csv" | "odm" | "xml"
    DateRangeBegin        string   // Date filter begin (ISO 8601)
    DateRangeEnd          string   // Date filter end
    FilterLogic           string   // Filter logic expression
}
```

#### Import Options

```go
type ImportOptions struct {
    Format            string   // "json" | "csv" | "odm" | "xml" | "spss" | "r"
    Type              string   // "flat" | "eav"
    OverwriteBehavior string   // "normal" | "overwrite" | "Upsert"
    ForceAutoNumber   bool
    DateFormat        string   // "YMD" | "MDY" | "DMY"
    ReturnContent     string   // "count" | "ids" | "auto_ids"
}
```

#### Import Result

```go
type ImportResult struct {
    Count     int      `json:"count"`
    IDs       []string `json:"ids,omitempty"`
    Error     string   `json:"error,omitempty"`
}
```

---

## 5. Error Handling

### Error Types

```go
// Error codes for REDCap API errors
const (
    ErrCodeInvalidRequest      = "INVALID_REQUEST"
    ErrCodeUnauthorized        = "UNAUTHORIZED"
    ErrCodeForbidden           = "FORBIDDEN"
    ErrCodeNotFound            = "NOT_FOUND"
    ErrCodeRateLimit           = "RATE_LIMIT"
    ErrCodeServerError         = "SERVER_ERROR"
    ErrCodeUnknown             = "UNKNOWN"
)

// Error represents a REDCap API error.
type Error struct {
    Code      string
    Message   string
    StatusCode int
    Err      error
}

func (e *Error) Error() string {
    return fmt.Sprintf("redcap: %s (%d): %s", e.Code, e.StatusCode, e.Message)
}

func (e *Error) Unwrap() error {
    return e.Err
}

// IsRetryable returns true if the error is transient and worth retrying.
func (e *Error) IsRetryable() bool {
    return e.Code == ErrCodeRateLimit || e.Code == ErrCodeServerError
}
```

### Error Handling Patterns

```go
// Library usage
client, err := redcap.NewClient("https://redcap.example.com/api/", "your-token")
if err != nil {
    return fmt.Errorf("creating client: %w", err)
}

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

records, err := client.ExportRecords(ctx, redcap.ExportRecordsOpt{
    Fields: []string{"record_id", "name"},
    Format: "json",
})
if err != nil {
    var redcapErr *redcap.Error
    if errors.As(err, &redcapErr) {
        switch redcapErr.Code {
        case redcap.ErrCodeUnauthorized:
            return fmt.Errorf("invalid or expired token")
        case redcap.ErrCodeRateLimit:
            return fmt.Errorf("rate limited, retry after backoff")
        default:
            return fmt.Errorf("REDCap error: %w", err)
        }
    }
    return fmt.Errorf("exporting records: %w", err)
}
```

---

## 6. CLI Specification

### Command Structure

```
cap [global options] command [command options] [arguments]
```

### Global Options

| Option | Description |
|--------|-------------|
| `--url, -u` | REDCap API URL (required or `CAP_URL` env) |
| `--token, -t` | REDCap API token (required or `CAP_TOKEN` env) |
| `--timeout` | Request timeout (default: `30s`) |
| `--verbose, -v` | Enable verbose output |
| `--config` | Config file path (default: `~/.cap.yaml`) |

### Commands

#### `cap export records`

Export records from a project.

```
cap export records [options]
```

Options:
| Option | Description |
|--------|-------------|
| `--records` | Comma-separated record IDs |
| `--fields` | Comma-separated field names |
| `--forms` | Comma-separated form names |
| `--format` | Output format: `json`, `csv` (default: `json`) |
| `--out, -o` | Output file (default: stdout) |
| `--raw` | Use raw values (default: label) |
| `--filter` | Filter logic expression |

Example:
```bash
cap export records --forms demographics --format csv -o demographics.csv
```

#### `cap export metadata`

Export the data dictionary.

```
cap export metadata [options]
```

Options:
| Option | Description |
|--------|-------------|
| `--format` | Output format: `json`, `csv` (default: `json`) |
| `--out, -o` | Output file |

Example:
```bash
cap export metadata --format csv -o metadata.csv
```

#### `cap export forms`

List all forms/instruments in the project.

```
cap export forms [options]
```

Options:
| Option | Description |
|--------|-------------|
| `--out, -o` | Output file |

#### `cap export users`

Export project users.

```
cap export users [options]
```

#### `cap export events`

Export longitudinal events (if applicable).

```
cap export events [options]
```

#### `cap import records`

Import records from a file.

```
cap import records [file] [options]
```

Options:
| Option | Description |
|--------|-------------|
| `--format` | Input format: `json`, `csv` (default: auto-detect) |
| `--overwrite` | Overwrite behavior: `normal`, `overwrite` (default: `normal`) |
| `--force-number` | Force auto-numbering of records |
| `--dry-run` | Validate without importing |

Example:
```bash
cap import records data.csv --format csv --overwrite overwrite
```

#### `cap project info`

Display project information.

```
cap project info
```

#### `cap version`

Display version information.

```
cap version
```

### Configuration File

Default location: `~/.cap.yaml`

```yaml
url: "https://redcap.example.com/api/"
token: "your-api-token-here"
timeout: "30s"
log_level: "info"

# Per-project aliases
projects:
  study1:
    url: "https://study1.redcap.example.com/api/"
    token: "token1"
  study2:
    url: "https://study2.redcap.example.com/api/"
    token: "token2"
```

Usage with project alias:
```bash
cap --project study1 export records --forms demographics
```

---

## 7. Reliability Features

### 7.1 Retry Logic

- **Exponential backoff** with jitter
- **Retryable errors**: Rate limiting (429), server errors (5xx), network timeouts
- **Non-retryable errors**: Authentication (401), forbidden (403), not found (404), invalid request (400)
- Configurable max retries and initial delay

```go
type RetryConfig struct {
    MaxRetries    int
    InitialDelay  time.Duration
    MaxDelay      time.Duration
    JitterFactor  float64
}
```

### 7.2 Rate Limiting

- Default: 50 requests per second (configurable)
- Uses token bucket algorithm
- Per-request queuing with context cancellation support

```go
type RateLimiter interface {
    Wait(ctx context.Context) error
    SetRate(rps float64)
}
```

### 7.3 Timeouts

- Default HTTP timeout: 30 seconds
- Per-operation context support for cancellation
- Long-running operations should use `context.WithTimeout`

---

## 8. Testing Strategy

### Unit Tests

- Test each client method with mocked HTTP responses
- Use `httptest` for HTTP server simulation
- Test error handling paths
- Test option configurations

### Integration Tests (Optional)

- Skip by default (`//go:build integration`)
- Requires real REDCap instance and API token
- Environment variable: `REDCAP_URL`, `REDCAP_TOKEN`

### Example Test Structure

```go
// pkg/redcap/client_test.go

func TestClient_ExportRecords(t *testing.T) {
    // Setup mock server
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Verify request params
        assert.Equal(t, "record", r.FormValue("content"))
        assert.Equal(t, "json", r.FormValue("format"))
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`[{"record_id":"1","name":"John"}]`))
    })
    
    server := httptest.NewServer(mux)
    defer server.Close()
    
    client, err := NewClient(server.URL+"/", "test-token")
    require.NoError(t, err)
    
    records, err := client.ExportRecords(context.Background())
    require.NoError(t, err)
    require.Len(t, records, 1)
    assert.Equal(t, "1", records[0].ID)
}
```

---

## 9. Dependencies

### Required

- Go 1.21+
- Standard library only for core functionality

### Recommended (for production)

- `github.com/google/uuid` - For UUID generation
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/cenkalti/backoff/v4` - Backoff/retry logic (or implement custom)
- `golang.org/x/time/rate` - Rate limiting

---

## 10. Implementation Priority

### Phase 1: Core Client (MVP)

1. Client struct with functional options
2. ExportRecords (JSON)
3. ExportMetadata
4. Basic error handling
5. Context support

### Phase 2: Additional Exports

1. ExportUsers
2. ExportEvents / ExportArms
3. ExportForms / ExportInstrumentList
4. ExportFile
5. Field mapping and type handling

### Phase 3: Import Operations

1. ImportRecords (JSON/CSV)
2. ImportFile
3. GenerateNextRecordName
4. ImportUsers / ImportDAGs / ImportEvents

### Phase 4: CLI

1. Basic CLI structure with Cobra
2. Export commands
3. Import commands
4. Configuration file support

### Phase 5: Polish

1. Comprehensive tests
2. Documentation
3. Rate limiting tuning
4. Retry logic refinement

---

## 11. Example Usage

### Library

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/yourname/go-cap/pkg/redcap"
)

func main() {
    client, err := redcap.NewClient(
        "https://redcap.example.com/api/",
        "your-api-token",
        redcap.WithTimeout(30 * time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Test connection
    if err := client.Ping(context.Background()); err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    
    // Get project info
    project, err := client.ExportProject(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Project: %s\n", project.ProjectTitle)
    
    // Export specific fields
    records, err := client.ExportRecords(context.Background(),
        redcap.ExportFields([]string{"record_id", "first_name", "last_name"}),
        redcap.ExportFormat("json"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    for _, r := range records {
        fmt.Printf("Record: %s\n", r.ID)
    }
}
```

### CLI

```bash
# Export metadata
cap export metadata --format csv -o metadata.csv

# Export specific form to CSV
cap export records --forms demographics --format csv -o demographics.csv

# Import records
cap import records data.csv --format csv --overwrite overwrite

# Using config file
export CAP_URL="https://redcap.example.com/api/"
export CAP_TOKEN="your-token"
cap project info
```

---

## 12. Version Compatibility

This library targets REDCap API version 14.x and later. The REDCap API is versioned, and different institutions may run different versions. The client should:

- Use the `version` parameter in API calls to request a specific format
- Handle backward-compatible changes gracefully
- Document minimum supported REDCap version

---

*Last Updated: February 2026*
