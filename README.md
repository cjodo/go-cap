# go-cap

A Go client library and CLI tool for the REDCap API.

> [!NOTE]
> Forked from in dev repo [tjrivera/go-cap](https://github.com/tjrivera/go-cap) 

> [!IMPORTANT]
> This repo is still in development. It is not production ready.  The first production ready release will be tagged v1.0.0

## Installation

```bash
go get github.com/cjodo/go-cap
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/cjodo/go-cap"
)

func main() {
    client, err := redcap.NewClient(
        "https://redcap.example.com/api/",
        "your-api-token",
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
    fmt.Printf("Project: %v\n", project["project_title"])
    
    // Export records
    records, err := client.ExportRecords(context.Background(),
        redcap.ExportFields([]string{"record_id", "first_name"}),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    for _, r := range records {
        fmt.Printf("Record: %s\n", r.ID)
    }
}
```

## CLI Usage

```bash
# Set environment variables
export REDCAP_URL="https://redcap.example.com/api/"
export REDCAP_TOKEN="your-token"

# Test connection
cap ping

# Export metadata
cap export metadata --format csv -o metadata.csv

# Export records
cap export records --forms demographics --format csv -o data.csv

# Import records
cap import records data.csv --format csv
```

## Features

- Full REDCap API support
- Context cancellation support
- Automatic retry with exponential backoff
- Rate limiting
- Type-safe Go client
- CLI tool for common operations

## API Endpoints Supported

- Project info
- Metadata/Data dictionary
- Records (export/import)
- Instruments/Forms
- Events (longitudinal)
- Arms
- Users
- Data Access Groups (DAGs)
- Files
- Repeating forms/events
- Form-event mappings
