package main

import (
	"fmt"
	"log"

	"github.com/zbiljic/vconfig-go"
)

// AppConfig represents our application configuration
// Note: Version field is required by vconfig
type AppConfig struct {
	Version      string
	AppName      string
	Debug        bool
	DatabaseURL  string
	MaxRetries   int
	AllowedHosts []string
}

func main() {
	// Create a new configuration
	config := &AppConfig{
		Version:      "1",
		AppName:      "My CLI App",
		Debug:        false,
		DatabaseURL:  "postgres://localhost/myapp",
		MaxRetries:   3,
		AllowedHosts: []string{"localhost", "example.com"},
	}

	// Save configuration to file
	fmt.Println("Saving configuration...")
	err := vconfig.SaveConfig(config, "app_config.json")
	if err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}
	fmt.Println("Configuration saved to app_config.json")

	// Check the version without loading the entire config
	version, err := vconfig.GetVersion("app_config.json")
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("Config file version: %s\n", version)

	// Load configuration from file
	fmt.Println("\nLoading configuration...")
	loadedConfig, err := vconfig.LoadConfig[AppConfig]("app_config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Display loaded configuration
	fmt.Println("Loaded configuration:")
	fmt.Printf("  Version: %s\n", loadedConfig.Version)
	fmt.Printf("  App Name: %s\n", loadedConfig.AppName)
	fmt.Printf("  Debug: %v\n", loadedConfig.Debug)
	fmt.Printf("  Database URL: %s\n", loadedConfig.DatabaseURL)
	fmt.Printf("  Max Retries: %d\n", loadedConfig.MaxRetries)
	fmt.Printf("  Allowed Hosts: %v\n", loadedConfig.AllowedHosts)

	// Demonstrate validation
	fmt.Println("\nValidating configuration struct...")
	err = vconfig.CheckData(loadedConfig)
	if err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}
	fmt.Println("Configuration is valid!")

	// Demonstrate what happens with invalid config (no Version field)
	type InvalidConfig struct {
		Name string
		Port int
	}

	invalidConfig := &InvalidConfig{
		Name: "test",
		Port: 8080,
	}

	fmt.Println("\nTrying to save invalid configuration (no Version field)...")
	err = vconfig.SaveConfig(invalidConfig, "invalid.json")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}
