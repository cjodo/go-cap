# API Reference

## Client

### NewClient

```go
func NewClient(baseURL, token string, opts ...Option) (*Client, error)
```

Creates a new REDCap API client.

**Parameters:**
- `baseURL` - REDCap API URL (e.g., "https://redcap.example.com/api/")
- `token` - API token
- `opts` - Optional configuration options

**Options:**

```go
// Set custom HTTP client
redcap.WithHTTPClient(&http.Client{Timeout: 60 * time.Second})

// Set maximum retry attempts (default: 3)
redcap.WithMaxRetries(5)

// Set retry delay (default: 1 second)
redcap.WithRetryDelay(2 * time.Second)

// Set custom rate limiter
redcap.WithRateLimiter(myRateLimiter)
```

### Ping

```go
func (c *Client) Ping(ctx context.Context) error
```

Verifies API connectivity.

### ExportProject

```go
func (c *Client) ExportProject(ctx context.Context) (map[string]interface{}, error)
```

Returns project information.

## Records

### ExportRecords

```go
func (c *Client) ExportRecords(ctx context.Context, opts ...ExportOption) ([]Record, error)
```

Exports records from the project.

**Options:**

```go
// Filter by record IDs
redcap.ExportRecordsFilter([]string{"1", "2", "3"})

// Limit fields
redcap.ExportFields([]string{"record_id", "first_name", "last_name"})

// Limit forms
redcap.ExportForms([]string{"demographics"})

// Filter by events (longitudinal)
redcap.ExportEvents([]string{"baseline_event", "followup_event"})

// Raw or label values
redcap.ExportRawOrLabel("label")  // "raw", "label", "both"

// Export format
redcap.ExportFormat("json")  // "json", "csv", "odm", "xml"

// Checkbox labels
redcap.ExportCheckboxLabel(true)

// Survey fields
redcap.ExportSurveyFields(true)

// DAG field
redcap.ExportDataAccessGroups(true)

// Filter logic
redcap.ExportFilterLogic("[age] > 18")
```

### ExportRecordsRaw

```go
func (c *Client) ExportRecordsRaw(ctx context.Context, opts ...ExportOption) ([]byte, error)
```

Exports raw format (CSV/JSON) for records.

### ImportRecords

```go
func (c *Client) ImportRecords(ctx context.Context, records []Record, opts ...ImportOption) (*ImportResult, error)
```

Imports records into the project.

**Options:**

```go
// Import format
redcap.ImportFormat("json")  // "json", "csv", "odm", "xml"

// Overwrite behavior
redcap.ImportOverwriteBehavior("overwrite")  // "normal", "overwrite", "Upsert"

// Force auto-numbering
redcap.ImportForceAutoNumber(true)

// Return content
redcap.ImportReturnContent("ids")  // "count", "ids", "auto_ids"
```

### GenerateNextRecordName

```go
func (c *Client) GenerateNextRecordName(ctx context.Context) (string, error)
```

Generates the next sequential record name.

## Metadata

### ExportMetadata

```go
func (c *Client) ExportMetadata(ctx context.Context) ([]Field, error)
```

Returns the data dictionary (metadata) for the project.

### ExportFieldNames

```go
func (c *Client) ExportFieldNames(ctx context.Context) ([]string, error)
```

Returns the list of export field names.

## Instruments

### ExportInstruments

```go
func (c *Client) ExportInstruments(ctx context.Context) ([]Instrument, error)
```

Returns the list of instruments/forms in the project.

## Events (Longitudinal)

### ExportEvents

```go
func (c *Client) ExportEvents(ctx context.Context) ([]Event, error)
```

Returns the list of events for longitudinal projects.

### ExportArms

```go
func (c *Client) ExportArms(ctx context.Context) ([]Arm, error)
```

Returns the list of arms for longitudinal projects.

## Users & Permissions

### ExportUsers

```go
func (c *Client) ExportUsers(ctx context.Context) ([]User, error)
```

Returns the list of users in the project.

### ExportDAGs

```go
func (c *Client) ExportDAGs(ctx context.Context) ([]DAG, error)
```

Returns the list of Data Access Groups.

## Repeating & Mappings

### ExportRepeatingFormsEvents

```go
func (c *Client) ExportRepeatingFormsEvents(ctx context.Context) ([]RepeatingForm, error)
```

Returns repeating form/event information.

### ExportFormEventMapping

```go
func (c *Client) ExportFormEventMapping(ctx context.Context) ([]FormEventMapping, error)
```

Returns form-event mappings.

## Files

### ExportFile

```go
func (c *Client) ExportFile(ctx context.Context, recordID, field, event string) ([]byte, error)
```

Exports a file field from a record.

### ImportFile

```go
func (c *Client) ImportFile(ctx context.Context, recordID, field, event string, data []byte, opts ...ImportOption) error
```

Imports a file into a record field.

### DeleteFile

```go
func (c *Client) DeleteFile(ctx context.Context, recordID, field, event string) error
```

Deletes a file from a record field.

## Types

### Record

```go
type Record struct {
    ID        string
    Fields    map[string]any
    EventName string
}
```

### Field

```go
type Field struct {
    FieldName   string
    FieldLabel  string
    FieldType   string
    FormName    string
    // ... other fields
}
```

### Instrument

```go
type Instrument struct {
    Name  string
    Label string
}
```

### Event

```go
type Event struct {
    Name            string
    ArmNum          int
    DayOffset       string
    UniqueEventName string
}
```

### User

```go
type User struct {
    Username         string
    Email            string
    FirstName        string
    LastName         string
    // ... other fields
}
```
