// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

// Package main provides the entry point for the mcp-plugin CLI.
package main

import (
	"fmt"
	"os"

	"github.com/gizzahub/gzh-cli-mcp-plugin/cmd/mcp-plugin/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
