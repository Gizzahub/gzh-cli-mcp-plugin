// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

// Package config provides configuration reading for Claude Code MCP settings.
package config

// MCPServer represents an MCP server configuration.
type MCPServer struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`    // "http" or "command"
	URL     string            `json:"url"`     // For HTTP type
	Command string            `json:"command"` // For command type (npx, uvx)
	Args    []string          `json:"args"`    // Command arguments
	Headers map[string]string `json:"headers"` // HTTP headers
	Enabled bool              `json:"enabled"`
	Source  string            `json:"source"` // Config file source
}

// ClaudeConfig represents the ~/.claude.json structure.
type ClaudeConfig struct {
	Projects map[string]ProjectConfig `json:"projects"`
}

// ProjectConfig represents per-project configuration.
type ProjectConfig struct {
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
}

// MCPServerConfig represents the raw MCP server config from JSON.
type MCPServerConfig struct {
	Type    string            `json:"type"`
	URL     string            `json:"url"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Headers map[string]string `json:"headers"`
}

// PluginMCPConfig represents .mcp.json structure.
type PluginMCPConfig struct {
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
}

// SettingsConfig represents ~/.claude/settings.json structure.
type SettingsConfig struct {
	EnabledPlugins map[string]bool `json:"enabledPlugins"`
}
