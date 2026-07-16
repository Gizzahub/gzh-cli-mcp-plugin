// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/infrastructure/npm"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var all bool
	var dryRun bool
	var force bool

	cmd := &cobra.Command{
		Use:   "update [server]",
		Short: "Update MCP servers to latest version",
		Long: `Update MCP servers to their latest versions.

This command checks for updates for npm-based MCP servers (npx command)
and updates them to the latest version from the npm registry.

Note: Only servers using 'npx' command can be updated automatically.
HTTP-based and uvx-based servers need manual updates.

Examples:
  # Check for updates on all servers
  mcp-plugin update --all --dry-run

  # Update a specific server
  mcp-plugin update context7

  # Update all servers
  mcp-plugin update --all

  # Force update even if already latest
  mcp-plugin update context7 --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("specify a server name or use --all")
			}

			var serverName string
			if len(args) > 0 {
				serverName = args[0]
			}

			return runUpdate(serverName, all, dryRun, force)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Update all updatable servers")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be updated without making changes")
	cmd.Flags().BoolVar(&force, "force", false, "Force update even if already at latest version")

	return cmd
}

// ServerUpdate represents an update check result.
type ServerUpdate struct {
	Name           string
	PackageName    string
	CurrentVersion string
	LatestVersion  string
	CanUpdate      bool
	Reason         string
}

func runUpdate(serverName string, all, dryRun, force bool) error {
	reader := config.NewReader()
	writer := config.NewWriter()
	npmClient := npm.NewClient()

	servers, err := reader.ListMCPServers()
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No MCP servers configured.")
		return nil
	}

	toCheck := filterServersToUpdate(servers, serverName, all)
	if len(toCheck) == 0 {
		if serverName != "" {
			return fmt.Errorf("server '%s' not found", serverName)
		}
		fmt.Println("No servers to update.")
		return nil
	}

	fmt.Println("Checking for updates...")
	fmt.Println()

	updates, updatable := collectServerUpdates(toCheck, npmClient)

	fmt.Println()

	if updatable == 0 && !force {
		fmt.Println("All servers are up to date.")
		return nil
	}

	if dryRun {
		fmt.Printf("Dry run: %d server(s) would be updated.\n", updatable)
		return nil
	}

	updated, failed := applyServerUpdates(writer, updates, force)
	fmt.Println()
	fmt.Printf("Update complete: %d updated, %d failed\n", updated, failed)
	return nil
}

func filterServersToUpdate(servers []config.MCPServer, serverName string, all bool) []config.MCPServer {
	var toCheck []config.MCPServer
	for _, server := range servers {
		if all || server.Name == serverName {
			toCheck = append(toCheck, server)
		}
	}
	return toCheck
}

func collectServerUpdates(toCheck []config.MCPServer, npmClient *npm.Client) (updates []ServerUpdate, updatable int) {
	updates = make([]ServerUpdate, 0, len(toCheck))

	for _, server := range toCheck {
		update := checkServerUpdate(server, npmClient)
		updates = append(updates, update)
		printUpdateStatus(update)
		if update.CanUpdate {
			updatable++
		}
	}
	return updates, updatable
}

func printUpdateStatus(update ServerUpdate) {
	if update.CanUpdate {
		if update.CurrentVersion != "" {
			fmt.Printf("📦 %s: %s → %s\n", update.Name, update.CurrentVersion, update.LatestVersion)
			return
		}
		fmt.Printf("📦 %s: (unversioned) → %s\n", update.Name, update.LatestVersion)
		return
	}
	if update.Reason == "up-to-date" {
		fmt.Printf("✅ %s: %s (up to date)\n", update.Name, update.CurrentVersion)
		return
	}
	fmt.Printf("⏭️  %s: %s\n", update.Name, update.Reason)
}

func applyServerUpdates(writer *config.Writer, updates []ServerUpdate, force bool) (updated, failed int) {
	for _, update := range updates {
		if !update.CanUpdate && !force {
			continue
		}
		if update.PackageName == "" {
			continue
		}
		if err := applyOneServerUpdate(writer, update); err != nil {
			fmt.Printf("⚠️  %v\n", err)
			failed++
			continue
		}
		fmt.Printf("✅ Updated %s to %s\n", update.Name, update.LatestVersion)
		updated++
	}
	return updated, failed
}

func applyOneServerUpdate(writer *config.Writer, update ServerUpdate) error {
	existingServers, err := writer.ListMCPServersGlobal()
	if err != nil {
		return fmt.Errorf("failed to read config for %s: %w", update.Name, err)
	}

	entry, exists := existingServers[update.Name]
	if !exists {
		return fmt.Errorf("server %s not found in global config", update.Name)
	}

	newArgs := updateArgsToLatest(entry.Args, update.PackageName, update.LatestVersion)
	newEntry := config.MCPServerEntry{
		Type:    entry.Type,
		Command: entry.Command,
		Args:    newArgs,
		URL:     entry.URL,
		Headers: entry.Headers,
		Enabled: entry.Enabled,
	}

	if err := writer.RemoveMCPServer(update.Name); err != nil {
		return fmt.Errorf("failed to update %s: %w", update.Name, err)
	}
	if err := writer.AddMCPServer(update.Name, newEntry); err != nil {
		return fmt.Errorf("failed to update %s: %w", update.Name, err)
	}
	return nil
}

func checkServerUpdate(server config.MCPServer, npmClient *npm.Client) ServerUpdate {
	update := ServerUpdate{
		Name: server.Name,
	}

	// Only npx-based servers can be updated
	if server.Command != "npx" {
		if server.Command != "" {
			update.Reason = fmt.Sprintf("not npm-based (%s)", server.Command)
		} else {
			update.Reason = "HTTP-based server"
		}
		return update
	}

	packageName, currentVersion := extractPackageInfo(server.Args)
	if packageName == "" {
		update.Reason = "cannot determine package name"
		return update
	}

	update.PackageName = packageName
	update.CurrentVersion = currentVersion

	pkgDetail, err := npmClient.GetPackage(packageName)
	if err != nil {
		update.Reason = fmt.Sprintf("npm error: %v", err)
		return update
	}

	latestVersion := pkgDetail.LatestVersion()
	if latestVersion == "" {
		update.Reason = "no latest version found"
		return update
	}

	update.LatestVersion = latestVersion

	if currentVersion == latestVersion {
		update.Reason = "up-to-date"
		return update
	}

	update.CanUpdate = true
	return update
}

// extractPackageInfo extracts package name and version from npx args.
// Examples:
//   - ["-y", "@upstash/context7-mcp"] -> "@upstash/context7-mcp", ""
//   - ["-y", "@package/name@1.0.0"] -> "@package/name", "1.0.0"
//   - ["@modelcontextprotocol/server-sequential-thinking"] -> "@modelcontextprotocol/server-sequential-thinking", ""
func extractPackageInfo(args []string) (packageName, version string) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if strings.Contains(arg, "@") {
			atCount := strings.Count(arg, "@")

			switch atCount {
			case 1:
				if strings.HasPrefix(arg, "@") {
					return arg, ""
				}
				parts := strings.SplitN(arg, "@", 2)
				return parts[0], parts[1]
			case 2:
				lastAt := strings.LastIndex(arg, "@")
				return arg[:lastAt], arg[lastAt+1:]
			}
		}

		return arg, ""
	}

	return "", ""
}

// updateArgsToLatest updates the args to use the latest version.
func updateArgsToLatest(args []string, packageName, latestVersion string) []string {
	newArgs := make([]string, len(args))
	copy(newArgs, args)

	versionPattern := regexp.MustCompile(`@\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?$`)

	for i, arg := range newArgs {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if strings.Contains(arg, packageName) || arg == packageName {
			cleanName := versionPattern.ReplaceAllString(arg, "")
			newArgs[i] = fmt.Sprintf("%s@%s", cleanName, latestVersion)
			break
		}
	}

	return newArgs
}
