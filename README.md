# vconfig-go

A Go library for managing versioned configuration files in JSON format, primarily designed for CLI applications.

## Features

- **Versioned Configuration**: Ensures all configuration structs have a `Version` field
- **Type-Safe Loading**: Uses Go generics for type-safe configuration loading
- **Cross-Platform**: Handles line ending differences between Windows and Unix systems
- **Simple API**: Easy to use functions for loading and saving configurations

## Installation

```bash
go get github.com/zbiljic/vconfig-go
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zbiljic/vconfig-go"
)

// Your configuration struct must have a Version field
type Config struct {
    Version  string
    APIKey   string
    Host     string
    Port     int
    Features []string
}

func main() {
    // Create a new configuration
    config := &Config{
        Version:  "1.0",
        APIKey:   "your-api-key",
        Host:     "localhost",
        Port:     8080,
        Features: []string{"feature1", "feature2"},
    }
    
    // Save configuration to file
    err := vconfig.SaveConfig(config, "config.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Load configuration from file
    loadedConfig, err := vconfig.LoadConfig[Config]("config.json")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Loaded config version: %s\n", loadedConfig.Version)
}
```

### Get Version Without Loading Full Config

```go
version, err := vconfig.GetVersion("config.json")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Config version: %s\n", version)
```

### Custom Validation

```go
// CheckData validates that your struct has the required Version field
err := vconfig.CheckData(myConfig)
if err != nil {
    log.Fatal(err)
}
```

## API Reference

### Functions

- `CheckData(data any) error` - Validates that the config struct has a Version field of type string
- `GetVersion(filename string) (string, error)` - Extracts just the version from a config file
- `LoadConfig[T any](filename string) (*T, error)` - Loads and validates a configuration file
- `SaveConfig(config any, filename string) error` - Saves a configuration to a JSON file

## Requirements

- Go 1.18+ (uses generics)
- Configuration structs must have a `Version` field of type `string`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
