// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Writer writes Claude Code MCP configurations.
type Writer struct {
	homeDir string
}

// NewWriter creates a new configuration writer.
func NewWriter() *Writer {
	home, _ := os.UserHomeDir()
	return &Writer{homeDir: home}
}

// SetPluginEnabled enables or disables a plugin in settings.json.
func (w *Writer) SetPluginEnabled(pluginID string, enabled bool) error {
	path := filepath.Join(w.homeDir, ".claude", "settings.json")

	// Read existing settings
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	// Parse as generic map to preserve all fields
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse settings: %w", err)
	}

	// Get or create enabledPlugins map
	enabledPlugins, ok := settings["enabledPlugins"].(map[string]interface{})
	if !ok {
		enabledPlugins = make(map[string]interface{})
	}

	// Update the plugin state
	enabledPlugins[pluginID] = enabled
	settings["enabledPlugins"] = enabledPlugins

	// Write back with pretty formatting
	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// ListPlugins returns the list of all known plugins with their enabled status.
func (w *Writer) ListPlugins() (map[string]bool, error) {
	path := filepath.Join(w.homeDir, ".claude", "settings.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %w", err)
	}

	result := make(map[string]bool)
	enabledPlugins, ok := settings["enabledPlugins"].(map[string]interface{})
	if !ok {
		return result, nil
	}

	for id, val := range enabledPlugins {
		if enabled, ok := val.(bool); ok {
			result[id] = enabled
		}
	}

	return result, nil
}

// PluginExists checks if a plugin exists in the settings.
func (w *Writer) PluginExists(pluginID string) (bool, error) {
	plugins, err := w.ListPlugins()
	if err != nil {
		return false, err
	}
	_, exists := plugins[pluginID]
	return exists, nil
}

// GetPluginStatus returns the current enabled status of a plugin.
func (w *Writer) GetPluginStatus(pluginID string) (enabled bool, exists bool, err error) {
	plugins, err := w.ListPlugins()
	if err != nil {
		return false, false, err
	}
	status, exists := plugins[pluginID]
	return status, exists, nil
}
