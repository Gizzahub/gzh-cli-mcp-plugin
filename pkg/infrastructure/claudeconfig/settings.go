// Package claudeconfig provides infrastructure for reading and writing
// Claude Code configuration files.
package claudeconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Settings represents the ~/.claude/settings.json structure.
type Settings struct {
	Schema         string                 `json:"$schema,omitempty"`
	Permissions    *Permissions           `json:"permissions,omitempty"`
	Hooks          map[string]interface{} `json:"hooks,omitempty"`
	EnabledPlugins map[string]bool        `json:"enabledPlugins,omitempty"`
	// Preserve other fields
	Extra map[string]interface{} `json:"-"`
}

// Permissions represents permission settings.
type Permissions struct {
	Allow       []string `json:"allow,omitempty"`
	Deny        []string `json:"deny,omitempty"`
	DefaultMode string   `json:"defaultMode,omitempty"`
}

// Reader reads Claude Code configuration.
type Reader struct {
	configDir string
}

// NewReader creates a new configuration reader.
// If configDir is empty, uses ~/.claude.
func NewReader(configDir string) *Reader {
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".claude")
	}
	return &Reader{configDir: configDir}
}

// SettingsPath returns the path to settings.json.
func (r *Reader) SettingsPath() string {
	return filepath.Join(r.configDir, "settings.json")
}

// ReadSettings reads the settings.json file.
func (r *Reader) ReadSettings() (*Settings, error) {
	path := r.SettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty settings if file doesn't exist
			return &Settings{
				EnabledPlugins: make(map[string]bool),
			}, nil
		}
		return nil, fmt.Errorf("read settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("parse settings: %w", err)
	}

	// Initialize map if nil
	if settings.EnabledPlugins == nil {
		settings.EnabledPlugins = make(map[string]bool)
	}

	return &settings, nil
}

// WriteSettings writes settings back to settings.json.
func (r *Reader) WriteSettings(settings *Settings) error {
	path := r.SettingsPath()

	// Read existing file to preserve unknown fields
	existingData, err := os.ReadFile(path)
	if err == nil {
		var existing map[string]interface{}
		if err := json.Unmarshal(existingData, &existing); err == nil {
			// Update only the fields we manage
			existing["enabledPlugins"] = settings.EnabledPlugins
			data, err := json.MarshalIndent(existing, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal settings: %w", err)
			}
			return os.WriteFile(path, data, 0644)
		}
	}

	// If we can't read existing, write fresh
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// ListEnabledPlugins returns a list of enabled plugin IDs.
func (r *Reader) ListEnabledPlugins() ([]string, error) {
	settings, err := r.ReadSettings()
	if err != nil {
		return nil, err
	}

	var enabled []string
	for id, isEnabled := range settings.EnabledPlugins {
		if isEnabled {
			enabled = append(enabled, id)
		}
	}
	return enabled, nil
}

// ListAllPlugins returns all plugins with their enabled status.
func (r *Reader) ListAllPlugins() (map[string]bool, error) {
	settings, err := r.ReadSettings()
	if err != nil {
		return nil, err
	}
	return settings.EnabledPlugins, nil
}

// EnablePlugin enables a plugin by ID.
func (r *Reader) EnablePlugin(pluginID string) error {
	settings, err := r.ReadSettings()
	if err != nil {
		return err
	}

	settings.EnabledPlugins[pluginID] = true
	return r.WriteSettings(settings)
}

// DisablePlugin disables a plugin by ID.
func (r *Reader) DisablePlugin(pluginID string) error {
	settings, err := r.ReadSettings()
	if err != nil {
		return err
	}

	settings.EnabledPlugins[pluginID] = false
	return r.WriteSettings(settings)
}

// IsPluginEnabled checks if a plugin is enabled.
func (r *Reader) IsPluginEnabled(pluginID string) (bool, error) {
	settings, err := r.ReadSettings()
	if err != nil {
		return false, err
	}

	enabled, exists := settings.EnabledPlugins[pluginID]
	if !exists {
		return false, nil
	}
	return enabled, nil
}
