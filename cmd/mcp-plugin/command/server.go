// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

func newServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage MCP servers",
		Long: `Server management commands for MCP servers.

Note: Claude Code manages MCP server lifecycle internally.
These commands help verify server configuration and availability.

Examples:
  # Check status of all servers
  mcp-plugin server status

  # Check status of a specific server
  mcp-plugin server status context7

  # Show detailed server information
  mcp-plugin server info kubernetes`,
	}

	cmd.AddCommand(newServerStatusCmd())
	cmd.AddCommand(newServerInfoCmd())

	return cmd
}

func newServerStatusCmd() *cobra.Command {
	var checkHealth bool

	cmd := &cobra.Command{
		Use:   "status [server]",
		Short: "Check MCP server status",
		Long: `Check the status and health of MCP servers.

Without arguments, checks all configured servers.
With a server name, checks only that specific server.

Status checks include:
- Configuration validity
- Enabled/disabled state
- HTTP server reachability (with --health flag)
- Command availability (with --health flag)

Examples:
  # Quick status of all servers
  mcp-plugin server status

  # Check with health probes
  mcp-plugin server status --health

  # Check specific server
  mcp-plugin server status context7 --health`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServerStatus(args, checkHealth)
		},
	}

	cmd.Flags().BoolVar(&checkHealth, "health", false, "Perform health checks (HTTP ping, command verification)")

	return cmd
}

func newServerInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <server>",
		Short: "Show detailed server information",
		Long: `Display detailed configuration information for an MCP server.

Shows:
- Server type (http, command, stdio)
- Configuration source file
- Command and arguments (for command-based servers)
- URL and headers (for HTTP servers)
- Environment variables

Examples:
  mcp-plugin server info context7
  mcp-plugin server info kubernetes`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServerInfo(args[0])
		},
	}

	return cmd
}

func runServerStatus(args []string, checkHealth bool) error {
	reader := config.NewReader()

	servers, err := reader.ListMCPServers()
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No MCP servers configured.")
		return nil
	}

	// Filter by name if provided
	var targetName string
	if len(args) > 0 {
		targetName = args[0]
	}

	var found bool
	for _, server := range servers {
		if targetName != "" && server.Name != targetName {
			continue
		}
		found = true

		// Basic status
		status := "disabled"
		statusIcon := "○"
		if server.Enabled {
			status = "enabled"
			statusIcon = "●"
		}

		fmt.Printf("%s %s (%s)\n", statusIcon, server.Name, status)
		fmt.Printf("  Type: %s\n", server.Type)

		if checkHealth {
			health := checkServerHealth(server)
			fmt.Printf("  Health: %s\n", health)
		}

		fmt.Println()
	}

	if targetName != "" && !found {
		return fmt.Errorf("server '%s' not found", targetName)
	}

	return nil
}

func runServerInfo(name string) error {
	reader := config.NewReader()

	servers, err := reader.ListMCPServers()
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	for _, server := range servers {
		if server.Name != name {
			continue
		}

		fmt.Printf("Server: %s\n", server.Name)
		fmt.Printf("─────────────────────────────────\n")

		// Status
		status := "disabled"
		if server.Enabled {
			status = "enabled"
		}
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Type: %s\n", server.Type)
		fmt.Printf("Source: %s\n", server.Source)

		// Command-based server details
		if server.Command != "" {
			fmt.Printf("\nCommand Configuration:\n")
			fmt.Printf("  Command: %s\n", server.Command)
			if len(server.Args) > 0 {
				fmt.Printf("  Args: %s\n", strings.Join(server.Args, " "))
			}

			// Check if command exists
			path, err := exec.LookPath(server.Command)
			if err != nil {
				fmt.Printf("  ⚠️  Command not found in PATH\n")
			} else {
				fmt.Printf("  Path: %s\n", path)
			}
		}

		// HTTP server details
		if server.URL != "" {
			fmt.Printf("\nHTTP Configuration:\n")
			fmt.Printf("  URL: %s\n", server.URL)
			if len(server.Headers) > 0 {
				fmt.Printf("  Headers:\n")
				for key, value := range server.Headers {
					// Mask sensitive values
					displayValue := value
					if strings.Contains(strings.ToLower(key), "auth") ||
						strings.Contains(strings.ToLower(key), "token") ||
						strings.Contains(strings.ToLower(key), "key") {
						if len(value) > 8 {
							displayValue = value[:4] + "..." + value[len(value)-4:]
						} else {
							displayValue = "****"
						}
					}
					fmt.Printf("    %s: %s\n", key, displayValue)
				}
			}
		}

		// Health check
		fmt.Printf("\nHealth Check:\n")
		health := checkServerHealth(server)
		fmt.Printf("  %s\n", health)

		return nil
	}

	return fmt.Errorf("server '%s' not found", name)
}

func checkServerHealth(server config.MCPServer) string {
	switch {
	case server.URL != "":
		return checkHTTPHealth(server.URL)
	case server.Command != "":
		return checkCommandHealth(server.Command)
	default:
		return "⚠️  Unknown server type"
	}
}

func checkHTTPHealth(url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return fmt.Sprintf("❌ Invalid URL: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Try GET if HEAD fails
		req.Method = http.MethodGet
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Sprintf("❌ Unreachable: %v", err)
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return fmt.Sprintf("✅ Reachable (HTTP %d)", resp.StatusCode)
	}
	return fmt.Sprintf("⚠️  HTTP %d", resp.StatusCode)
}

func checkCommandHealth(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		return fmt.Sprintf("❌ Command not found: %s", command)
	}
	return fmt.Sprintf("✅ Command available: %s", path)
}
