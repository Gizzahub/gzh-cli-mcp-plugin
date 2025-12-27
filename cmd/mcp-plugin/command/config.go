// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

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
	cmd.AddCommand(newConfigExportCmd())
	cmd.AddCommand(newConfigImportCmd())
	cmd.AddCommand(newConfigValidateCmd())

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

// ExportConfig represents the export file format.
type ExportConfig struct {
	Version    string                           `json:"version"`
	ExportedAt string                           `json:"exportedAt"`
	Servers    map[string]config.MCPServerEntry `json:"servers"`
}

func newConfigExportCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export MCP configuration to file",
		Long: `Export all MCP server configurations to a JSON file.

The exported file can be used to:
- Backup your configuration
- Share configuration with teammates
- Migrate to a new machine

Examples:
  # Export to stdout
  mcp-plugin config export

  # Export to file
  mcp-plugin config export -o mcp-backup.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigExport(outputFile)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")

	return cmd
}

func runConfigExport(outputFile string) error {
	writer := config.NewWriter()

	servers, err := writer.ListMCPServersGlobal()
	if err != nil {
		return fmt.Errorf("failed to read servers: %w", err)
	}

	export := ExportConfig{
		Version:    "1.0",
		ExportedAt: time.Now().Format(time.RFC3339),
		Servers:    servers,
	}

	output, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if outputFile == "" {
		fmt.Println(string(output))
		return nil
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Exported %d servers to %s\n", len(servers), outputFile)
	return nil
}

func newConfigImportCmd() *cobra.Command {
	var merge bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import MCP configuration from file",
		Long: `Import MCP server configurations from a JSON file.

By default, import will fail if servers already exist.
Use --merge to update existing servers.

Examples:
  # Import from file
  mcp-plugin config import mcp-backup.json

  # Dry run to preview changes
  mcp-plugin config import mcp-backup.json --dry-run

  # Merge with existing config
  mcp-plugin config import mcp-backup.json --merge`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigImport(args[0], merge, dryRun)
		},
	}

	cmd.Flags().BoolVar(&merge, "merge", false, "Merge with existing configuration (update existing servers)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be imported without making changes")

	return cmd
}

func runConfigImport(inputFile string, merge, dryRun bool) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var importConfig ExportConfig
	if err := json.Unmarshal(data, &importConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if len(importConfig.Servers) == 0 {
		fmt.Println("No servers found in import file.")
		return nil
	}

	writer := config.NewWriter()
	existingServers, _ := writer.ListMCPServersGlobal()

	var added, updated, skipped int

	for name, entry := range importConfig.Servers {
		_, exists := existingServers[name]

		if dryRun {
			if exists {
				if merge {
					fmt.Printf("  [update] %s\n", name)
					updated++
				} else {
					fmt.Printf("  [skip] %s (already exists)\n", name)
					skipped++
				}
			} else {
				fmt.Printf("  [add] %s\n", name)
				added++
			}
			continue
		}

		if exists {
			if merge {
				// Remove then add to update
				if err := writer.RemoveMCPServer(name); err != nil {
					fmt.Printf("⚠️  Failed to update %s: %v\n", name, err)
					continue
				}
				if err := writer.AddMCPServer(name, entry); err != nil {
					fmt.Printf("⚠️  Failed to update %s: %v\n", name, err)
					continue
				}
				updated++
			} else {
				fmt.Printf("⚠️  Skipped %s (already exists, use --merge to update)\n", name)
				skipped++
			}
		} else {
			if err := writer.AddMCPServer(name, entry); err != nil {
				fmt.Printf("⚠️  Failed to add %s: %v\n", name, err)
				continue
			}
			added++
		}
	}

	if dryRun {
		fmt.Printf("\nDry run summary: %d to add, %d to update, %d to skip\n", added, updated, skipped)
	} else {
		fmt.Printf("\n✅ Import complete: %d added, %d updated, %d skipped\n", added, updated, skipped)
	}

	return nil
}

func newConfigValidateCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate MCP configuration",
		Long: `Validate all MCP server configurations for common issues.

Checks performed:
- URL syntax and reachability (for HTTP servers)
- Command availability (for command-based servers)
- Required fields presence
- Duplicate server detection

Examples:
  # Quick validation
  mcp-plugin config validate

  # Verbose validation with details
  mcp-plugin config validate --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigValidate(verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed validation results")

	return cmd
}

// ValidationResult represents a validation check result.
type ValidationResult struct {
	Server  string
	Check   string
	Status  string // "pass", "warn", "fail"
	Message string
}

func runConfigValidate(verbose bool) error {
	reader := config.NewReader()

	servers, err := reader.ListMCPServers()
	if err != nil {
		return fmt.Errorf("failed to read servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No MCP servers configured.")
		return nil
	}

	var results []ValidationResult
	var passCount, warnCount, failCount int

	// Check for duplicates
	seen := make(map[string]string)
	for _, server := range servers {
		if existingSource, exists := seen[server.Name]; exists {
			results = append(results, ValidationResult{
				Server:  server.Name,
				Check:   "duplicate",
				Status:  "warn",
				Message: fmt.Sprintf("Duplicate definition (also in %s)", existingSource),
			})
			warnCount++
		}
		seen[server.Name] = server.Source
	}

	for _, server := range servers {
		// Check required fields
		if server.Type == "" {
			results = append(results, ValidationResult{
				Server:  server.Name,
				Check:   "type",
				Status:  "warn",
				Message: "Server type not specified (inferred)",
			})
			warnCount++
		}

		// Type-specific validation
		switch server.Type {
		case "http":
			result := validateHTTPServer(server)
			results = append(results, result)
			switch result.Status {
			case "pass":
				passCount++
			case "warn":
				warnCount++
			case "fail":
				failCount++
			}

		case "command":
			result := validateCommandServer(server)
			results = append(results, result)
			switch result.Status {
			case "pass":
				passCount++
			case "warn":
				warnCount++
			case "fail":
				failCount++
			}

		default:
			if server.URL != "" {
				result := validateHTTPServer(server)
				results = append(results, result)
				switch result.Status {
				case "pass":
					passCount++
				case "warn":
					warnCount++
				case "fail":
					failCount++
				}
			} else if server.Command != "" {
				result := validateCommandServer(server)
				results = append(results, result)
				switch result.Status {
				case "pass":
					passCount++
				case "warn":
					warnCount++
				case "fail":
					failCount++
				}
			}
		}
	}

	// Print results
	if verbose {
		fmt.Println("Validation Results:")
		fmt.Println("─────────────────────────────────")
		for _, r := range results {
			var icon string
			switch r.Status {
			case "warn":
				icon = "⚠️"
			case "fail":
				icon = "❌"
			default:
				icon = "✅"
			}
			fmt.Printf("%s %s [%s]: %s\n", icon, r.Server, r.Check, r.Message)
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("Validation Summary: %d servers checked\n", len(servers))
	fmt.Printf("  ✅ Pass: %d\n", passCount)
	fmt.Printf("  ⚠️  Warnings: %d\n", warnCount)
	fmt.Printf("  ❌ Failures: %d\n", failCount)

	if failCount > 0 {
		return fmt.Errorf("validation failed with %d errors", failCount)
	}

	return nil
}

func validateHTTPServer(server config.MCPServer) ValidationResult {
	if server.URL == "" {
		return ValidationResult{
			Server:  server.Name,
			Check:   "url",
			Status:  "fail",
			Message: "HTTP server has no URL configured",
		}
	}

	// Validate URL syntax
	parsedURL, err := url.Parse(server.URL)
	if err != nil {
		return ValidationResult{
			Server:  server.Name,
			Check:   "url_syntax",
			Status:  "fail",
			Message: fmt.Sprintf("Invalid URL: %v", err),
		}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ValidationResult{
			Server:  server.Name,
			Check:   "url_scheme",
			Status:  "fail",
			Message: fmt.Sprintf("Invalid URL scheme: %s (expected http or https)", parsedURL.Scheme),
		}
	}

	// Check reachability (with timeout)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Head(server.URL)
	if err != nil {
		// Check if it's a connection issue or auth required
		if strings.Contains(err.Error(), "connection refused") {
			return ValidationResult{
				Server:  server.Name,
				Check:   "reachability",
				Status:  "warn",
				Message: "Server unreachable (may be offline or firewalled)",
			}
		}
		return ValidationResult{
			Server:  server.Name,
			Check:   "reachability",
			Status:  "warn",
			Message: fmt.Sprintf("Cannot verify: %v", err),
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return ValidationResult{
			Server:  server.Name,
			Check:   "reachability",
			Status:  "pass",
			Message: fmt.Sprintf("Reachable (HTTP %d - may require auth)", resp.StatusCode),
		}
	}

	return ValidationResult{
		Server:  server.Name,
		Check:   "reachability",
		Status:  "pass",
		Message: fmt.Sprintf("Reachable (HTTP %d)", resp.StatusCode),
	}
}

func validateCommandServer(server config.MCPServer) ValidationResult {
	if server.Command == "" {
		return ValidationResult{
			Server:  server.Name,
			Check:   "command",
			Status:  "fail",
			Message: "Command server has no command configured",
		}
	}

	// Check if command exists in PATH
	path, err := exec.LookPath(server.Command)
	if err != nil {
		return ValidationResult{
			Server:  server.Name,
			Check:   "command",
			Status:  "fail",
			Message: fmt.Sprintf("Command '%s' not found in PATH", server.Command),
		}
	}

	return ValidationResult{
		Server:  server.Name,
		Check:   "command",
		Status:  "pass",
		Message: fmt.Sprintf("Command available: %s", path),
	}
}
