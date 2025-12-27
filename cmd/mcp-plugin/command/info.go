// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"strings"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/infrastructure/npm"
	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <package>",
		Short: "Show information about an MCP package",
		Long: `Display detailed information about an npm package.

This fetches package metadata from npm registry including
version, description, repository, and usage hints.

Examples:
  # Get info about context7 MCP server
  mcp-plugin info @upstash/context7-mcp

  # Get info about playwright MCP server
  mcp-plugin info @playwright/mcp`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(args[0])
		},
	}

	return cmd
}

func runInfo(packageName string) error {
	client := npm.NewClient()

	fmt.Printf("Fetching information for '%s'...\n\n", packageName)

	pkg, err := client.GetPackage(packageName)
	if err != nil {
		return fmt.Errorf("failed to get package info: %w", err)
	}

	// Package name
	fmt.Printf("Package: %s\n", pkg.Name)
	fmt.Printf("Version: %s\n", pkg.LatestVersion())
	
	if pkg.License != "" {
		fmt.Printf("License: %s\n", pkg.License)
	}

	// Description
	if pkg.Description != "" {
		fmt.Printf("\nDescription:\n  %s\n", pkg.Description)
	}

	// Author
	if pkg.Author != nil && pkg.Author.Name != "" {
		fmt.Printf("\nAuthor: %s", pkg.Author.Name)
		if pkg.Author.Email != "" {
			fmt.Printf(" <%s>", pkg.Author.Email)
		}
		fmt.Println()
	}

	// Links
	fmt.Println("\nLinks:")
	fmt.Printf("  npm: https://www.npmjs.com/package/%s\n", pkg.Name)
	if pkg.Homepage != "" {
		fmt.Printf("  homepage: %s\n", pkg.Homepage)
	}
	if pkg.Repository.URL != "" {
		repoURL := pkg.Repository.URL
		repoURL = strings.TrimPrefix(repoURL, "git+")
		repoURL = strings.TrimSuffix(repoURL, ".git")
		fmt.Printf("  repository: %s\n", repoURL)
	}

	// Usage hint for MCP packages
	if strings.Contains(strings.ToLower(pkg.Name), "mcp") {
		fmt.Println("\nUsage with Claude Code:")
		fmt.Printf("  Add to MCP config: npx -y %s\n", pkg.Name)
	}

	// Show available versions count
	fmt.Printf("\nVersions: %d available\n", len(pkg.Versions))

	return nil
}
