// Package repository provides repository implementations.
package repository

import (
	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/application/port/output"
	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/infrastructure/claudeconfig"
)

// Ensure PluginRepository implements the interface.
var _ output.PluginRepository = (*PluginRepository)(nil)

// PluginRepository implements plugin storage using Claude config.
type PluginRepository struct {
	reader *claudeconfig.Reader
}

// NewPluginRepository creates a new plugin repository.
func NewPluginRepository(configDir string) *PluginRepository {
	return &PluginRepository{
		reader: claudeconfig.NewReader(configDir),
	}
}

// ListAll returns all plugins with their enabled status.
func (r *PluginRepository) ListAll() (map[string]bool, error) {
	return r.reader.ListAllPlugins()
}

// ListEnabled returns only enabled plugin IDs.
func (r *PluginRepository) ListEnabled() ([]string, error) {
	return r.reader.ListEnabledPlugins()
}

// IsEnabled checks if a plugin is enabled.
func (r *PluginRepository) IsEnabled(pluginID string) (bool, error) {
	return r.reader.IsPluginEnabled(pluginID)
}

// Enable enables a plugin.
func (r *PluginRepository) Enable(pluginID string) error {
	return r.reader.EnablePlugin(pluginID)
}

// Disable disables a plugin.
func (r *PluginRepository) Disable(pluginID string) error {
	return r.reader.DisablePlugin(pluginID)
}
