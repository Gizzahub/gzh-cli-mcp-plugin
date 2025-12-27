// Package main provides the entry point for the mcp-plugin CLI.
package main

import (
	"fmt"
	"os"

	"github.com/gizzahub/gzh-cli-mcp-plugin/cmd/mcp-plugin/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
