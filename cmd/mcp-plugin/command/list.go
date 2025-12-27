// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

var listEnabledOnly bool

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List MCP servers",
		Long:  `List all configured MCP servers from Claude Code configuration.`,
		RunE:  runList,
	}

	cmd.Flags().BoolVar(&listEnabledOnly, "enabled", false, "Show only enabled servers")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	reader := config.NewReader()
	
	servers, err := reader.ListMCPServers()
	if err != nil {
		return fmt.Errorf("failed to list MCP servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No MCP servers found.")
		return nil
	}

	fmt.Printf("Found %d MCP server(s):\n\n", len(servers))
	
	for _, server := range servers {
		if listEnabledOnly && !server.Enabled {
			continue
		}
		
		status := "disabled"
		if server.Enabled {
			status = "enabled"
		}
		
		fmt.Printf("  %s (%s)\n", server.Name, status)
		fmt.Printf("    Type: %s\n", server.Type)
		if server.URL != "" {
			fmt.Printf("    URL: %s\n", server.URL)
		}
		if server.Command != "" {
			fmt.Printf("    Command: %s %v\n", server.Command, server.Args)
		}
		fmt.Println()
	}

	return nil
}
