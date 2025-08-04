package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/zbiljic/vconfig-go"
)

// Mutex for thread-safe access to config file
var configMutex = &sync.RWMutex{}

// Cached configs to avoid repeated file reads
var (
	cachedConfigV1 *ConfigV1
	cachedConfigV2 *ConfigV2
)

const (
	configStateVersionV1 = "1"
	configStateVersionV2 = "2"
)

// ConfigV1 represents version 1 of our state configuration
type ConfigV1 struct {
	Version    string   `json:"version"`
	Roots      []string `json:"roots"`
	Checkpoint string   `json:"checkpoint"`
}

// ConfigV2 represents version 2 with additional tracking fields
type ConfigV2 struct {
	Version        string    `json:"version"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
	Roots          []string  `json:"roots"`
	TotalCount     int       `json:"total_count"`
	RemainingCount int       `json:"remaining_count"`
	Paths          []string  `json:"paths"`
}

func main() {
	fmt.Println("State Config Migration Example (Real-World Pattern)")
	fmt.Println("===================================================")

	// Load, create, or migrate config as needed
	config, err := loadCreateMigrateConfig(".", "example-state", "root1", "root2")
	if err != nil {
		log.Fatalf("Failed to load/create/migrate config: %v", err)
	}

	fmt.Printf("\nSuccessfully loaded config version: %s\n", config.Version)
	displayConfigV2(config)

	// Demonstrate updating the config with checkpointing
	fmt.Println("\n--- Simulating Progress Updates ---")
	for i := 0; i < 5; i++ {
		path := fmt.Sprintf("path/to/file%d.txt", i)
		config.Paths = append(config.Paths, path)
		config.TotalCount++

		// Save checkpoint every 2 files (like configStateCheckpointInterval)
		if i%2 == 0 {
			fmt.Printf("Checkpoint at file %d\n", i)
			if err := saveConfigV2(".", "example-state", config); err != nil {
				log.Printf("Failed to save checkpoint: %v", err)
			}
		}
	}

	// Final save
	if err := saveConfigV2(".", "example-state", config); err != nil {
		log.Fatalf("Failed to save final config: %v", err)
	}

	fmt.Println("\n--- Final State ---")
	displayConfigV2(config)

	// Clean up
	fmt.Println("\n--- Cleaning Up ---")
	if err := clearConfig(".", "example-state"); err != nil {
		log.Printf("Failed to clear config: %v", err)
	}
	fmt.Println("Config cleared successfully")
}

// Helper to generate config filename
func configFilename(baseDir, stateName string) string {
	return fmt.Sprintf("%s/.state-%s.json", baseDir, stateName)
}

// kept for completeness of the migration example
//
//nolint:unused
func newConfigV1() *ConfigV1 {
	config := new(ConfigV1)
	config.Version = configStateVersionV1
	config.Roots = make([]string, 0)
	return config
}

// newConfigV2 creates a new v2 config
func newConfigV2() *ConfigV2 {
	config := new(ConfigV2)
	config.Version = configStateVersionV2
	config.CreateTime = time.Now()
	config.UpdateTime = time.Now()
	config.Roots = make([]string, 0)
	config.Paths = make([]string, 0)
	return config
}

// loadConfigV1 loads a v1 config with caching
func loadConfigV1(baseDir, stateName string) (*ConfigV1, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Return cached if available
	if cachedConfigV1 != nil {
		return cachedConfigV1, nil
	}

	filename := configFilename(baseDir, stateName)
	config, err := vconfig.LoadConfig[ConfigV1](filename)
	if err != nil {
		return nil, err
	}

	// Cache the config
	cachedConfigV1 = config
	return config, nil
}

// loadConfigV2 loads a v2 config with caching
func loadConfigV2(baseDir, stateName string) (*ConfigV2, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Return cached if available
	if cachedConfigV2 != nil {
		return cachedConfigV2, nil
	}

	filename := configFilename(baseDir, stateName)
	config, err := vconfig.LoadConfig[ConfigV2](filename)
	if err != nil {
		return nil, err
	}

	// Cache the config
	cachedConfigV2 = config
	return config, nil
}

// saveConfigV2 saves v2 config with mutex and cache update
func saveConfigV2(baseDir, stateName string, config *ConfigV2) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	config.UpdateTime = time.Now()
	config.RemainingCount = len(config.Paths)

	filename := configFilename(baseDir, stateName)
	if err := vconfig.SaveConfig(config, filename); err != nil {
		return err
	}

	// Update cache
	cachedConfigV2 = config
	return nil
}

// createConfig creates a new v2 config with roots
func createConfig(roots ...string) (*ConfigV2, error) {
	config := newConfigV2()
	config.Roots = append(config.Roots, roots...)
	return config, nil
}

// loadCreateMigrateConfig loads existing config or creates new one, handling migrations
func loadCreateMigrateConfig(baseDir, stateName string, roots ...string) (*ConfigV2, error) {
	filename := configFilename(baseDir, stateName)

	version, err := vconfig.GetVersion(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// No config exists, create new v2
			fmt.Println("No existing config found, creating new v2 config...")
			config, err := createConfig(roots...)
			if err != nil {
				return nil, err
			}
			if err := saveConfigV2(baseDir, stateName, config); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}

	// Handle different versions
	switch version {
	case configStateVersionV1:
		fmt.Println("Found v1 config, migrating to v2...")
		currentConfig, err := loadConfigV1(baseDir, stateName)
		if err != nil {
			return nil, fmt.Errorf("unable to load config version '%s': %w", version, err)
		}

		// Migrate v1 to v2
		newConfig := newConfigV2()
		newConfig.Roots = make([]string, len(currentConfig.Roots))
		copy(newConfig.Roots, currentConfig.Roots)

		if err := saveConfigV2(baseDir, stateName, newConfig); err != nil {
			return nil, err
		}

		// Recursively call to load the migrated config
		return loadCreateMigrateConfig(baseDir, stateName, roots...)

	case configStateVersionV2:
		fmt.Println("Found v2 config, loading...")
		currentConfig, err := loadConfigV2(baseDir, stateName)
		if err != nil {
			return nil, fmt.Errorf("unable to load config version '%s': %w", version, err)
		}
		return currentConfig, nil

	default:
		return nil, fmt.Errorf("unknown config version: '%s'", version)
	}
}

// displayConfigV2 displays the v2 config details
func displayConfigV2(config *ConfigV2) {
	fmt.Println("\nConfiguration Details:")
	fmt.Printf("  Version: %s\n", config.Version)
	fmt.Printf("  Created: %s\n", config.CreateTime.Format(time.RFC3339))
	fmt.Printf("  Updated: %s\n", config.UpdateTime.Format(time.RFC3339))
	fmt.Printf("  Roots: %v\n", config.Roots)
	fmt.Printf("  Total Files: %d\n", config.TotalCount)
	fmt.Printf("  Remaining: %d\n", config.RemainingCount)
	fmt.Printf("  Paths Processed: %d\n", len(config.Paths))
}

// clearConfig removes the config file and clears caches
func clearConfig(baseDir, stateName string) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Clear caches
	cachedConfigV1 = nil
	cachedConfigV2 = nil

	filename := configFilename(baseDir, stateName)
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return os.Remove(filename)
}
