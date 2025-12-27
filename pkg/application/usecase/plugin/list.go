// Package plugin provides use cases for MCP plugin management.
package plugin

import (
	"sort"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/application/port/output"
)

// PluginInfo represents information about a plugin.
type PluginInfo struct {
	ID        string
	Name      string
	Publisher string
	Enabled   bool
}

// ListUseCase handles listing plugins.
type ListUseCase struct {
	repo output.PluginRepository
}

// NewListUseCase creates a new list use case.
func NewListUseCase(repo output.PluginRepository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

// Execute lists all plugins.
func (uc *ListUseCase) Execute(enabledOnly bool) ([]PluginInfo, error) {
	plugins, err := uc.repo.ListAll()
	if err != nil {
		return nil, err
	}

	var result []PluginInfo
	for id, enabled := range plugins {
		if enabledOnly && !enabled {
			continue
		}

		name, publisher := parsePluginID(id)
		result = append(result, PluginInfo{
			ID:        id,
			Name:      name,
			Publisher: publisher,
			Enabled:   enabled,
		})
	}

	// Sort by ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// parsePluginID splits "name@publisher" into parts.
func parsePluginID(id string) (name, publisher string) {
	for i := len(id) - 1; i >= 0; i-- {
		if id[i] == '@' {
			return id[:i], id[i+1:]
		}
	}
	return id, ""
}
