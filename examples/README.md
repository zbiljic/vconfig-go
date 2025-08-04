# vconfig-go Examples

This directory contains example code demonstrating how to use the vconfig-go library.

## Running the Examples

### Basic Example

The basic example demonstrates:
- Creating and saving a configuration
- Loading a configuration from file
- Checking version without loading the entire config
- Configuration validation

```bash
cd examples
go run basic_example.go
```

This will create an `app_config.json` file in the current directory.

### Migration Example (Real-World Pattern)

This example demonstrates a production-ready migration pattern for CLI tools:
- Thread-safe configuration access with mutex protection
- Caching to avoid repeated file reads
- Automatic version detection and migration
- Checkpoint saving during long operations
- Progress tracking with file counters

```bash
cd examples
go run migration_example.go
```

The example will:
1. Create a new v2 config if none exists
2. Automatically migrate v1 configs to v2
3. Demonstrate checkpoint saves during processing
4. Show how to track progress across multiple files

## Clean Up

To remove generated config files:
```bash
rm app_config.json
```
