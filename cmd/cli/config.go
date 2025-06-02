package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// Config represents the CLI configuration
type Config struct {
    DefaultServer string            `json:"default_server"`
    Timeout       int              `json:"timeout_seconds"`
    Aliases       map[string]string `json:"aliases"`
    History       bool             `json:"save_history"`
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
    home, err := os.UserHomeDir()
    if err != nil {
        return ".rushkv-cli.json"
    }
    return filepath.Join(home, ".rushkv-cli.json")
}

// loadConfig loads configuration from file or returns default config
func loadConfig() *Config {
    config := &Config{
        DefaultServer: "localhost:8080",
        Timeout:       5,
        Aliases:       make(map[string]string),
        History:       true,
    }
    
    configPath := getConfigPath()
    data, err := os.ReadFile(configPath)
    if err != nil {
        return config // Use default configuration
    }
    
    if err := json.Unmarshal(data, config); err != nil {
        fmt.Printf("Warning: Invalid configuration file format, using defaults\n")
        return config
    }
    
    return config
}

// save saves the configuration to file
func (c *Config) save() error {
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(getConfigPath(), data, 0644)
}