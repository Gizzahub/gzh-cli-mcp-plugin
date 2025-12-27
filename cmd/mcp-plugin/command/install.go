// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"strings"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

var (
	installHTTP    bool
	installURL     string
	installUVX     bool
	installCommand string
	installArgs    []string
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <name> [package]",
		Short: "Install an MCP server",
		Long: `Install an MCP server to Claude Code configuration.

By default, installs an npx-based MCP server:
  mcp-plugin install myserver @package/mcp-server

For HTTP-based servers:
  mcp-plugin install myserver --http --url https://api.example.com/mcp

For uvx-based (Python) servers:
  mcp-plugin install myserver --uvx mypackage

For custom command servers:
  mcp-plugin install myserver --command node --args server.js,--port,8080

Examples:
  # Install an npx MCP server
  mcp-plugin install context7 @upstash/context7-mcp

  # Install an HTTP MCP server
  mcp-plugin install myapi --http --url https://api.example.com/mcp

  # Install a uvx (Python) MCP server
  mcp-plugin install serena --uvx serena-mcp`,
		Args: cobra.RangeArgs(1, 2),
		RunE: runInstall,
	}

	cmd.Flags().BoolVar(&installHTTP, "http", false, "Install as HTTP-based server")
	cmd.Flags().StringVar(&installURL, "url", "", "URL for HTTP server (required with --http)")
	cmd.Flags().BoolVar(&installUVX, "uvx", false, "Install as uvx (Python) server")
	cmd.Flags().StringVar(&installCommand, "command", "", "Custom command (e.g., node, python)")
	cmd.Flags().StringSliceVar(&installArgs, "args", nil, "Custom command arguments")

	return cmd
}

func runInstall(cmd *cobra.Command, args []string) error {
	name := args[0]

	writer := config.NewWriter()

	// Check if server already exists
	exists, err := writer.MCPServerExists(name)
	if err != nil {
		return fmt.Errorf("failed to check server: %w", err)
	}
	if exists {
		return fmt.Errorf("MCP server '%s' already exists. Use 'remove' first to reinstall", name)
	}

	var entry config.MCPServerEntry

	switch {
	case installHTTP:
		// HTTP-based server
		if installURL == "" {
			return fmt.Errorf("--url is required for HTTP servers")
		}
		entry = config.MCPServerEntry{
			Type: "http",
			URL:  installURL,
		}
		fmt.Printf("Installing HTTP MCP server '%s'...\n", name)

	case installUVX:
		// uvx (Python) server
		var pkg string
		if len(args) > 1 {
			pkg = args[1]
		} else {
			pkg = name
		}
		entry = config.MCPServerEntry{
			Type:    "stdio",
			Command: "uvx",
			Args:    []string{pkg},
		}
		fmt.Printf("Installing uvx MCP server '%s' (package: %s)...\n", name, pkg)

	case installCommand != "":
		// Custom command server
		entry = config.MCPServerEntry{
			Type:    "stdio",
			Command: installCommand,
			Args:    installArgs,
		}
		fmt.Printf("Installing custom MCP server '%s' (command: %s)...\n", name, installCommand)

	default:
		// Default: npx server
		var pkg string
		if len(args) > 1 {
			pkg = args[1]
		} else {
			return fmt.Errorf("package name required for npx install (e.g., mcp-plugin install %s @package/name)", name)
		}

		// Add -y flag if not already present
		npxArgs := []string{"-y", pkg}
		entry = config.MCPServerEntry{
			Type:    "stdio",
			Command: "npx",
			Args:    npxArgs,
		}
		fmt.Printf("Installing npx MCP server '%s' (package: %s)...\n", name, pkg)
	}

	// Add the server
	if err := writer.AddMCPServer(name, entry); err != nil {
		return fmt.Errorf("failed to install server: %w", err)
	}

	fmt.Printf("MCP server '%s' has been installed.\n", name)
	printServerConfig(name, entry)
	fmt.Println("\nNote: Restart Claude Code for the new server to be available.")

	return nil
}

func printServerConfig(name string, entry config.MCPServerEntry) {
	fmt.Printf("\nConfiguration:\n")
	fmt.Printf("  Name: %s\n", name)
	if entry.Type != "" {
		fmt.Printf("  Type: %s\n", entry.Type)
	}
	if entry.Command != "" {
		fmt.Printf("  Command: %s\n", entry.Command)
	}
	if len(entry.Args) > 0 {
		fmt.Printf("  Args: %s\n", strings.Join(entry.Args, " "))
	}
	if entry.URL != "" {
		fmt.Printf("  URL: %s\n", entry.URL)
	}
}
