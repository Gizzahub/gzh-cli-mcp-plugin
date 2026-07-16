// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

// exportFilePerm is owner-only; exports may include auth headers.
const exportFilePerm = 0o600

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
	ExportedAt string                           `json:"exportedAt"` //nolint:tagliatelle // external protocol wire format
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

	if err := os.WriteFile(outputFile, output, exportFilePerm); err != nil {
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
	// #nosec G304 -- inputFile is an intentional user-provided CLI path
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
	existingServers, err := writer.ListMCPServersGlobal()
	if err != nil {
		existingServers = map[string]config.MCPServerEntry{}
	}

	added, updated, skipped := applyImport(writer, importConfig.Servers, existingServers, merge, dryRun)

	if dryRun {
		fmt.Printf("\nDry run summary: %d to add, %d to update, %d to skip\n", added, updated, skipped)
	} else {
		fmt.Printf("\n✅ Import complete: %d added, %d updated, %d skipped\n", added, updated, skipped)
	}

	return nil
}

func applyImport(
	writer *config.Writer,
	servers map[string]config.MCPServerEntry,
	existing map[string]config.MCPServerEntry,
	merge, dryRun bool,
) (added, updated, skipped int) {
	for name, entry := range servers {
		_, exists := existing[name]
		if dryRun {
			a, u, s := dryRunImportEntry(name, exists, merge)
			added += a
			updated += u
			skipped += s
			continue
		}
		a, u, s := applyImportEntry(writer, name, entry, exists, merge)
		added += a
		updated += u
		skipped += s
	}
	return added, updated, skipped
}

func dryRunImportEntry(name string, exists, merge bool) (added, updated, skipped int) {
	switch {
	case exists && merge:
		fmt.Printf("  [update] %s\n", name)
		return 0, 1, 0
	case exists:
		fmt.Printf("  [skip] %s (already exists)\n", name)
		return 0, 0, 1
	default:
		fmt.Printf("  [add] %s\n", name)
		return 1, 0, 0
	}
}

func applyImportEntry(
	writer *config.Writer,
	name string,
	entry config.MCPServerEntry,
	exists, merge bool,
) (added, updated, skipped int) {
	if exists {
		if !merge {
			fmt.Printf("⚠️  Skipped %s (already exists, use --merge to update)\n", name)
			return 0, 0, 1
		}
		if err := writer.RemoveMCPServer(name); err != nil {
			fmt.Printf("⚠️  Failed to update %s: %v\n", name, err)
			return 0, 0, 0
		}
		if err := writer.AddMCPServer(name, entry); err != nil {
			fmt.Printf("⚠️  Failed to update %s: %v\n", name, err)
			return 0, 0, 0
		}
		return 0, 1, 0
	}
	if err := writer.AddMCPServer(name, entry); err != nil {
		fmt.Printf("⚠️  Failed to add %s: %v\n", name, err)
		return 0, 0, 0
	}
	return 1, 0, 0
}
