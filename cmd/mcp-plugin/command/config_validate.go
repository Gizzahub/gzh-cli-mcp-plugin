// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/gizzahub/gzh-cli-mcp-plugin/pkg/config"
	"github.com/spf13/cobra"
)

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

	results, passCount, warnCount, failCount := validateAllServers(servers)

	if verbose {
		printValidationResults(results)
	}

	fmt.Printf("Validation Summary: %d servers checked\n", len(servers))
	fmt.Printf("  ✅ Pass: %d\n", passCount)
	fmt.Printf("  ⚠️  Warnings: %d\n", warnCount)
	fmt.Printf("  ❌ Failures: %d\n", failCount)

	if failCount > 0 {
		return fmt.Errorf("validation failed with %d errors", failCount)
	}

	return nil
}

func validateAllServers(servers []config.MCPServer) (results []ValidationResult, passCount, warnCount, failCount int) {
	seen := make(map[string]string)
	for _, server := range servers {
		if existingSource, exists := seen[server.Name]; exists {
			results = append(results, ValidationResult{
				Server:  server.Name,
				Check:   "duplicate",
				Status:  checkStatusWarn,
				Message: fmt.Sprintf("Duplicate definition (also in %s)", existingSource),
			})
			warnCount++
		}
		seen[server.Name] = server.Source
	}

	for _, server := range servers {
		serverResults, p, w, f := validateOneServer(server)
		results = append(results, serverResults...)
		passCount += p
		warnCount += w
		failCount += f
	}
	return results, passCount, warnCount, failCount
}

func validateOneServer(server config.MCPServer) (results []ValidationResult, passCount, warnCount, failCount int) {
	if server.Type == "" {
		results = append(results, ValidationResult{
			Server:  server.Name,
			Check:   "type",
			Status:  checkStatusWarn,
			Message: "Server type not specified (inferred)",
		})
		warnCount++
	}

	result, ok := typeSpecificValidation(server)
	if !ok {
		return results, passCount, warnCount, failCount
	}
	results = append(results, result)
	switch result.Status {
	case checkStatusPass:
		passCount++
	case checkStatusWarn:
		warnCount++
	case checkStatusFail:
		failCount++
	}
	return results, passCount, warnCount, failCount
}

func typeSpecificValidation(server config.MCPServer) (ValidationResult, bool) {
	switch server.Type {
	case config.TypeHTTP:
		return validateHTTPServer(server), true
	case config.TypeCommand:
		return validateCommandServer(server), true
	default:
		if server.URL != "" {
			return validateHTTPServer(server), true
		}
		if server.Command != "" {
			return validateCommandServer(server), true
		}
		return ValidationResult{}, false
	}
}

func printValidationResults(results []ValidationResult) {
	fmt.Println("Validation Results:")
	fmt.Println("─────────────────────────────────")
	for _, r := range results {
		var icon string
		switch r.Status {
		case checkStatusWarn:
			icon = "⚠️"
		case checkStatusFail:
			icon = "❌"
		default:
			icon = "✅"
		}
		fmt.Printf("%s %s [%s]: %s\n", icon, r.Server, r.Check, r.Message)
	}
	fmt.Println()
}

func validateHTTPServer(server config.MCPServer) ValidationResult {
	if server.URL == "" {
		return ValidationResult{
			Server:  server.Name,
			Check:   "url",
			Status:  checkStatusFail,
			Message: "HTTP server has no URL configured",
		}
	}

	parsedURL, err := url.Parse(server.URL)
	if err != nil {
		return ValidationResult{
			Server:  server.Name,
			Check:   "url_syntax",
			Status:  checkStatusFail,
			Message: fmt.Sprintf("Invalid URL: %v", err),
		}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ValidationResult{
			Server:  server.Name,
			Check:   "url_scheme",
			Status:  checkStatusFail,
			Message: fmt.Sprintf("Invalid URL scheme: %s (expected http or https)", parsedURL.Scheme),
		}
	}

	return checkHTTPReachability(server)
}

func checkHTTPReachability(server config.MCPServer) ValidationResult {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, server.URL, http.NoBody)
	if err != nil {
		return ValidationResult{
			Server:  server.Name,
			Check:   checkReachability,
			Status:  checkStatusWarn,
			Message: fmt.Sprintf("Cannot verify: %v", err),
		}
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return ValidationResult{
				Server:  server.Name,
				Check:   checkReachability,
				Status:  checkStatusWarn,
				Message: "Server unreachable (may be offline or firewalled)",
			}
		}
		return ValidationResult{
			Server:  server.Name,
			Check:   checkReachability,
			Status:  checkStatusWarn,
			Message: fmt.Sprintf("Cannot verify: %v", err),
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return ValidationResult{
			Server:  server.Name,
			Check:   checkReachability,
			Status:  checkStatusPass,
			Message: fmt.Sprintf("Reachable (HTTP %d - may require auth)", resp.StatusCode),
		}
	}

	return ValidationResult{
		Server:  server.Name,
		Check:   checkReachability,
		Status:  checkStatusPass,
		Message: fmt.Sprintf("Reachable (HTTP %d)", resp.StatusCode),
	}
}

func validateCommandServer(server config.MCPServer) ValidationResult {
	if server.Command == "" {
		return ValidationResult{
			Server:  server.Name,
			Check:   checkCommand,
			Status:  checkStatusFail,
			Message: "Command server has no command configured",
		}
	}

	path, err := exec.LookPath(server.Command)
	if err != nil {
		return ValidationResult{
			Server:  server.Name,
			Check:   checkCommand,
			Status:  checkStatusFail,
			Message: fmt.Sprintf("Command '%s' not found in PATH", server.Command),
		}
	}

	return ValidationResult{
		Server:  server.Name,
		Check:   checkCommand,
		Status:  checkStatusPass,
		Message: fmt.Sprintf("Command available: %s", path),
	}
}
