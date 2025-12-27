// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"

	"github.com/gizzahub/gzh-cli-mcp-plugin/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mcp-plugin version %s\n", version.Version)
		},
	}
}
