package plugin

import (
	"fmt"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/application/port/output"
)

// ToggleResult represents the result of enabling/disabling a plugin.
type ToggleResult struct {
	PluginID   string
	Enabled    bool
	WasAlready bool
	Message    string
}

// ToggleUseCase handles enabling and disabling plugins.
type ToggleUseCase struct {
	repo output.PluginRepository
}

// NewToggleUseCase creates a new toggle use case.
func NewToggleUseCase(repo output.PluginRepository) *ToggleUseCase {
	return &ToggleUseCase{repo: repo}
}

// Enable enables a plugin.
func (uc *ToggleUseCase) Enable(pluginID string) (*ToggleResult, error) {
	// Check current state
	currentlyEnabled, err := uc.repo.IsEnabled(pluginID)
	if err != nil {
		return nil, fmt.Errorf("check plugin status: %w", err)
	}

	if currentlyEnabled {
		return &ToggleResult{
			PluginID:   pluginID,
			Enabled:    true,
			WasAlready: true,
			Message:    fmt.Sprintf("Plugin '%s' is already enabled", pluginID),
		}, nil
	}

	if err := uc.repo.Enable(pluginID); err != nil {
		return nil, fmt.Errorf("enable plugin: %w", err)
	}

	return &ToggleResult{
		PluginID:   pluginID,
		Enabled:    true,
		WasAlready: false,
		Message:    fmt.Sprintf("Plugin '%s' enabled successfully", pluginID),
	}, nil
}

// Disable disables a plugin.
func (uc *ToggleUseCase) Disable(pluginID string) (*ToggleResult, error) {
	// Check current state
	currentlyEnabled, err := uc.repo.IsEnabled(pluginID)
	if err != nil {
		return nil, fmt.Errorf("check plugin status: %w", err)
	}

	if !currentlyEnabled {
		return &ToggleResult{
			PluginID:   pluginID,
			Enabled:    false,
			WasAlready: true,
			Message:    fmt.Sprintf("Plugin '%s' is already disabled", pluginID),
		}, nil
	}

	if err := uc.repo.Disable(pluginID); err != nil {
		return nil, fmt.Errorf("disable plugin: %w", err)
	}

	return &ToggleResult{
		PluginID:   pluginID,
		Enabled:    false,
		WasAlready: false,
		Message:    fmt.Sprintf("Plugin '%s' disabled successfully", pluginID),
	}, nil
}

// Status checks the current status of a plugin.
func (uc *ToggleUseCase) Status(pluginID string) (*ToggleResult, error) {
	enabled, err := uc.repo.IsEnabled(pluginID)
	if err != nil {
		return nil, fmt.Errorf("check plugin status: %w", err)
	}

	status := "disabled"
	if enabled {
		status = "enabled"
	}

	return &ToggleResult{
		PluginID: pluginID,
		Enabled:  enabled,
		Message:  fmt.Sprintf("Plugin '%s' is %s", pluginID, status),
	}, nil
}
