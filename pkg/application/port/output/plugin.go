// Package output defines output ports for the application layer.
package output

// PluginRepository defines the interface for plugin storage operations.
type PluginRepository interface {
	// ListAll returns all plugins with their enabled status.
	ListAll() (map[string]bool, error)

	// ListEnabled returns only enabled plugin IDs.
	ListEnabled() ([]string, error)

	// IsEnabled checks if a plugin is enabled.
	IsEnabled(pluginID string) (bool, error)

	// Enable enables a plugin.
	Enable(pluginID string) error

	// Disable disables a plugin.
	Disable(pluginID string) error
}
