// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"strings"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/infrastructure/npm"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for MCP packages on npm",
		Long: `Search npm registry for MCP-related packages.

This searches npm for packages containing "mcp" along with your query.
Use this to discover MCP servers you can use with Claude Code.

Examples:
  # Search for kubernetes-related MCP packages
  mcp-plugin search kubernetes

  # Search for database MCP packages
  mcp-plugin search database

  # Limit results
  mcp-plugin search postgres --limit 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(args[0], limit)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of results")

	return cmd
}

func runSearch(query string, limit int) error {
	client := npm.NewClient()

	fmt.Printf("Searching npm for MCP packages matching '%s'...\n\n", query)

	results, err := client.Search(query, limit)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results.Objects) == 0 {
		fmt.Println("No packages found.")
		return nil
	}

	fmt.Printf("Found %d package(s):\n\n", len(results.Objects))

	for _, obj := range results.Objects {
		pkg := obj.Package
		
		// Package name and version
		fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
		
		// Description (truncated)
		if pkg.Description != "" {
			desc := pkg.Description
			if len(desc) > 70 {
				desc = desc[:67] + "..."
			}
			fmt.Printf("    %s\n", desc)
		}

		// Score
		fmt.Printf("    Score: %.2f (quality: %.2f, popularity: %.2f)\n",
			obj.Score.Final, obj.Score.Detail.Quality, obj.Score.Detail.Popularity)

		// Install hint
		if strings.Contains(strings.ToLower(pkg.Name), "mcp") {
			fmt.Printf("    Usage: npx %s\n", pkg.Name)
		}

		fmt.Println()
	}

	fmt.Printf("Total: %d packages found\n", results.Total)
	fmt.Println("\nUse 'mcp-plugin info <package>' for more details.")

	return nil
}
