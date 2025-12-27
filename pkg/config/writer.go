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

// MCPServerEntry represents an MCP server configuration for installation.
type MCPServerEntry struct {
	Type    string            `json:"type,omitempty"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Enabled bool              `json:"enabled,omitempty"`
}

// AddMCPServer adds a new MCP server to claude.json.
func (w *Writer) AddMCPServer(name string, entry MCPServerEntry) error {
	path := filepath.Join(w.homeDir, ".claude.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read claude.json: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse claude.json: %w", err)
	}

	// Get or create mcpServers map
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	// Check if server already exists
	if _, exists := mcpServers[name]; exists {
		return fmt.Errorf("MCP server '%s' already exists", name)
	}

	// Add the new server
	serverConfig := make(map[string]interface{})
	if entry.Type != "" {
		serverConfig["type"] = entry.Type
	}
	if entry.Command != "" {
		serverConfig["command"] = entry.Command
	}
	if len(entry.Args) > 0 {
		serverConfig["args"] = entry.Args
	}
	if entry.URL != "" {
		serverConfig["url"] = entry.URL
	}
	if len(entry.Headers) > 0 {
		serverConfig["headers"] = entry.Headers
	}
	if entry.Enabled {
		serverConfig["enabled"] = entry.Enabled
	}

	mcpServers[name] = serverConfig
	config["mcpServers"] = mcpServers

	// Write back
	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal claude.json: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("failed to write claude.json: %w", err)
	}

	return nil
}

// RemoveMCPServer removes an MCP server from claude.json.
func (w *Writer) RemoveMCPServer(name string) error {
	path := filepath.Join(w.homeDir, ".claude.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read claude.json: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse claude.json: %w", err)
	}

	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	if _, exists := mcpServers[name]; !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	delete(mcpServers, name)
	config["mcpServers"] = mcpServers

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal claude.json: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("failed to write claude.json: %w", err)
	}

	return nil
}

// ListMCPServersGlobal returns global MCP servers from claude.json.
func (w *Writer) ListMCPServersGlobal() (map[string]MCPServerEntry, error) {
	path := filepath.Join(w.homeDir, ".claude.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read claude.json: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse claude.json: %w", err)
	}

	result := make(map[string]MCPServerEntry)
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		return result, nil
	}

	for name, v := range mcpServers {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		entry := MCPServerEntry{}
		if t, ok := cfg["type"].(string); ok {
			entry.Type = t
		}
		if cmd, ok := cfg["command"].(string); ok {
			entry.Command = cmd
		}
		if args, ok := cfg["args"].([]interface{}); ok {
			for _, arg := range args {
				if s, ok := arg.(string); ok {
					entry.Args = append(entry.Args, s)
				}
			}
		}
		if url, ok := cfg["url"].(string); ok {
			entry.URL = url
		}
		if enabled, ok := cfg["enabled"].(bool); ok {
			entry.Enabled = enabled
		}

		result[name] = entry
	}

	return result, nil
}

// MCPServerExists checks if an MCP server exists in claude.json.
func (w *Writer) MCPServerExists(name string) (bool, error) {
	servers, err := w.ListMCPServersGlobal()
	if err != nil {
		return false, err
	}
	_, exists := servers[name]
	return exists, nil
}
