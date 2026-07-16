// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

// Package config provides configuration reading for Claude Code MCP settings.
package config

// MCP server type wire values (Claude Code / MCP protocol).
const (
	TypeHTTP    = "http"
	TypeCommand = "command"
	TypeStdio   = "stdio"
)

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
	MCPServers map[string]MCPServerConfig `json:"mcpServers"` //nolint:tagliatelle // external protocol wire format
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
	MCPServers map[string]MCPServerConfig `json:"mcpServers"` //nolint:tagliatelle // external protocol wire format
}

// SettingsConfig represents ~/.claude/settings.json structure.
type SettingsConfig struct {
	EnabledPlugins map[string]bool `json:"enabledPlugins"` //nolint:tagliatelle // external protocol wire format
}

// resolveServerType picks the configured type or infers from fields.
func resolveServerType(cfg MCPServerConfig) string {
	switch {
	case cfg.Type != "":
		return cfg.Type
	case cfg.Command != "":
		return TypeCommand
	default:
		return TypeHTTP
	}
}

// serverFromConfig builds an MCPServer from raw config fields.
func serverFromConfig(name, source string, cfg MCPServerConfig) MCPServer {
	return MCPServer{
		Name:    name,
		Type:    resolveServerType(cfg),
		URL:     cfg.URL,
		Command: cfg.Command,
		Args:    cfg.Args,
		Headers: cfg.Headers,
		Source:  source,
	}
}
