// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"strings"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

func newDisableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable <plugin-id>",
		Short: "Disable an MCP plugin",
		Long: `Disable an MCP plugin in Claude Code settings.

The plugin-id should be in the format "name@publisher", for example:
  - context7@claude-plugins-official
  - greptile@claude-plugins-official

Use "mcp-plugin list" to see available plugins.`,
		Args: cobra.ExactArgs(1),
		RunE: runDisable,
	}

	return cmd
}

func runDisable(cmd *cobra.Command, args []string) error {
	pluginID := args[0]

	// Validate plugin ID format
	if !strings.Contains(pluginID, "@") {
		return fmt.Errorf("invalid plugin ID format: expected 'name@publisher', got '%s'", pluginID)
	}

	writer := config.NewWriter()

	// Check current status
	enabled, exists, err := writer.GetPluginStatus(pluginID)
	if err != nil {
		return fmt.Errorf("failed to check plugin status: %w", err)
	}

	if !exists {
		// Show available plugins
		plugins, _ := writer.ListPlugins()
		fmt.Printf("Plugin '%s' not found in settings.\n\n", pluginID)
		fmt.Println("Available plugins:")
		for id := range plugins {
			fmt.Printf("  - %s\n", id)
		}
		return fmt.Errorf("plugin not found")
	}

	if !enabled {
		fmt.Printf("Plugin '%s' is already disabled.\n", pluginID)
		return nil
	}

	// Disable the plugin
	if err := writer.SetPluginEnabled(pluginID, false); err != nil {
		return fmt.Errorf("failed to disable plugin: %w", err)
	}

	fmt.Printf("Plugin '%s' has been disabled.\n", pluginID)
	fmt.Println("Note: Restart Claude Code for changes to take effect.")

	return nil
}
