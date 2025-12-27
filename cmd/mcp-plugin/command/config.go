// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage MCP configuration",
		Long:  `Manage MCP server configuration and settings.`,
	}

	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigPathsCmd())

	return cmd
}

func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := config.NewReader()
			
			servers, err := reader.ListMCPServers()
			if err != nil {
				return fmt.Errorf("failed to read configuration: %w", err)
			}

			enabledCount := 0
			for _, s := range servers {
				if s.Enabled {
					enabledCount++
				}
			}

			fmt.Printf("MCP Configuration Summary:\n")
			fmt.Printf("  Total servers: %d\n", len(servers))
			fmt.Printf("  Enabled: %d\n", enabledCount)
			fmt.Printf("  Disabled: %d\n", len(servers)-enabledCount)
			
			return nil
		},
	}
}

func newConfigPathsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "paths",
		Short: "Show configuration file paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := config.NewReader()
			paths := reader.GetConfigPaths()

			fmt.Println("Configuration file paths:")
			for _, p := range paths {
				fmt.Printf("  %s\n", p)
			}
			
			return nil
		},
	}
}
