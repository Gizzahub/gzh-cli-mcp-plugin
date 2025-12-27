// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Reader reads Claude Code MCP configurations.
type Reader struct {
	homeDir string
}

// NewReader creates a new configuration reader.
func NewReader() *Reader {
	home, _ := os.UserHomeDir()
	return &Reader{homeDir: home}
}

// GetConfigPaths returns the list of configuration file paths.
func (r *Reader) GetConfigPaths() []string {
	return []string{
		filepath.Join(r.homeDir, ".claude.json"),
		filepath.Join(r.homeDir, ".claude", "settings.json"),
		filepath.Join(r.homeDir, ".claude", "plugins", "cache"),
	}
}

// ListMCPServers lists all configured MCP servers.
func (r *Reader) ListMCPServers() ([]MCPServer, error) {
	var servers []MCPServer

	// Read from ~/.claude.json
	claudeServers, err := r.readClaudeJSON()
	if err == nil {
		servers = append(servers, claudeServers...)
	}

	// Read from plugin cache
	pluginServers, err := r.readPluginConfigs()
	if err == nil {
		servers = append(servers, pluginServers...)
	}

	// Check enabled status from settings.json
	enabledPlugins, _ := r.readSettings()
	for i := range servers {
		if enabled, ok := enabledPlugins[servers[i].Name]; ok {
			servers[i].Enabled = enabled
		}
	}

	return servers, nil
}

func (r *Reader) readClaudeJSON() ([]MCPServer, error) {
	path := filepath.Join(r.homeDir, ".claude.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ClaudeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var servers []MCPServer
	for _, proj := range config.Projects {
		for name, cfg := range proj.MCPServers {
			server := MCPServer{
				Name:    name,
				URL:     cfg.URL,
				Command: cfg.Command,
				Args:    cfg.Args,
				Headers: cfg.Headers,
				Source:  path,
			}
			if cfg.Type != "" {
				server.Type = cfg.Type
			} else if cfg.Command != "" {
				server.Type = "command"
			} else {
				server.Type = "http"
			}
			servers = append(servers, server)
		}
	}

	return servers, nil
}

func (r *Reader) readPluginConfigs() ([]MCPServer, error) {
	var servers []MCPServer

	cacheDir := filepath.Join(r.homeDir, ".claude", "plugins", "cache")
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Look for .mcp.json in plugin subdirectories
		pluginDir := filepath.Join(cacheDir, entry.Name())
		subEntries, err := os.ReadDir(pluginDir)
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() {
				continue
			}

			mcpPath := filepath.Join(pluginDir, subEntry.Name(), ".mcp.json")
			data, err := os.ReadFile(mcpPath)
			if err != nil {
				continue
			}

			// Try parsing as map[string]MCPServerConfig first
			var rawConfig map[string]MCPServerConfig
			if err := json.Unmarshal(data, &rawConfig); err == nil {
				for name, cfg := range rawConfig {
					if name == "mcpServers" {
						continue // Skip if it's wrapped
					}
					server := MCPServer{
						Name:    name,
						URL:     cfg.URL,
						Command: cfg.Command,
						Args:    cfg.Args,
						Headers: cfg.Headers,
						Source:  mcpPath,
					}
					if cfg.Type != "" {
						server.Type = cfg.Type
					} else if cfg.Command != "" {
						server.Type = "command"
					} else {
						server.Type = "http"
					}
					servers = append(servers, server)
				}
			}

			// Try parsing as PluginMCPConfig
			var pluginConfig PluginMCPConfig
			if err := json.Unmarshal(data, &pluginConfig); err == nil && len(pluginConfig.MCPServers) > 0 {
				for name, cfg := range pluginConfig.MCPServers {
					server := MCPServer{
						Name:    name,
						URL:     cfg.URL,
						Command: cfg.Command,
						Args:    cfg.Args,
						Headers: cfg.Headers,
						Source:  mcpPath,
					}
					if cfg.Type != "" {
						server.Type = cfg.Type
					} else if cfg.Command != "" {
						server.Type = "command"
					} else {
						server.Type = "http"
					}
					servers = append(servers, server)
				}
			}
		}
	}

	return servers, nil
}

func (r *Reader) readSettings() (map[string]bool, error) {
	path := filepath.Join(r.homeDir, ".claude", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config SettingsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.EnabledPlugins, nil
}
