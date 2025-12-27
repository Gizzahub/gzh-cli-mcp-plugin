// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

// Package command provides CLI commands for mcp-plugin.
package command

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-plugin",
	Short: "MCP Plugin Manager for Claude Code",
	Long: `MCP Plugin Manager - Manage MCP (Model Context Protocol) servers.

This tool helps manage MCP servers used by Claude Code and similar AI tools.
It provides commands to list, enable, disable, and configure MCP servers.

Examples:
  # List all MCP servers
  mcp-plugin list

  # List enabled servers only
  mcp-plugin list --enabled

  # Show server configuration
  mcp-plugin config show

  # Enable a server
  mcp-plugin enable context7

  # Disable a server
  mcp-plugin disable context7`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newEnableCmd())
	rootCmd.AddCommand(newDisableCmd())
	rootCmd.AddCommand(newInstallCmd())
	rootCmd.AddCommand(newRemoveCmd())
	rootCmd.AddCommand(newSearchCmd())
	rootCmd.AddCommand(newInfoCmd())
	rootCmd.AddCommand(newServerCmd())
	rootCmd.AddCommand(newUpdateCmd())
}
