// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

var removeForce bool

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm", "uninstall"},
		Short:   "Remove an MCP server",
		Long: `Remove an MCP server from Claude Code configuration.

This removes the server entry from ~/.claude.json but does not uninstall
any npm or Python packages that may have been installed.

Examples:
  # Remove an MCP server
  mcp-plugin remove context7

  # Remove without confirmation
  mcp-plugin remove context7 --force`,
		Args: cobra.ExactArgs(1),
		RunE: runRemove,
	}

	cmd.Flags().BoolVarP(&removeForce, "force", "f", false, "Skip confirmation")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	writer := config.NewWriter()

	// Check if server exists
	exists, err := writer.MCPServerExists(name)
	if err != nil {
		return fmt.Errorf("failed to check server: %w", err)
	}

	if !exists {
		// Show available servers
		servers, _ := writer.ListMCPServersGlobal()
		fmt.Printf("MCP server '%s' not found.\n\n", name)
		if len(servers) > 0 {
			fmt.Println("Available servers:")
			for serverName := range servers {
				fmt.Printf("  - %s\n", serverName)
			}
		} else {
			fmt.Println("No MCP servers installed.")
		}
		return fmt.Errorf("server not found")
	}

	// Get server info before removing
	servers, _ := writer.ListMCPServersGlobal()
	serverInfo := servers[name]

	// Remove the server
	if err := writer.RemoveMCPServer(name); err != nil {
		return fmt.Errorf("failed to remove server: %w", err)
	}

	fmt.Printf("MCP server '%s' has been removed.\n", name)

	// Show what was removed
	if serverInfo.Command != "" {
		fmt.Printf("  (was: %s)\n", serverInfo.Command)
	} else if serverInfo.URL != "" {
		fmt.Printf("  (was: %s)\n", serverInfo.URL)
	}

	fmt.Println("\nNote: Restart Claude Code for changes to take effect.")

	return nil
}
