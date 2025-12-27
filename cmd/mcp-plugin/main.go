// Package main provides the entry point for the mcp-plugin CLI.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = ""
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mcp-plugin",
		Short: "MCP Plugin Management CLI",
		Long: `mcp-plugin is a CLI tool for managing MCP (Model Context Protocol) plugins
for Claude Code integration.

Features:
  - List installed MCP plugins and servers
  - Enable/disable plugins
  - Check plugin status
  - Manage MCP configuration`,
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", Version, GitCommit, BuildDate),
	}

	// Add subcommands
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(statusCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed MCP plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("MCP Plugins:")
			fmt.Println("  (not yet implemented)")
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show MCP plugin status",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("MCP Status:")
			fmt.Println("  (not yet implemented)")
			return nil
		},
	}
}
