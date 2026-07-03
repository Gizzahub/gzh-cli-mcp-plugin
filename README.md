# gzh-cli-mcp-plugin

> MCP (Model Context Protocol) plugin manager for Claude Code.

`mcp-plugin` installs, configures, and manages MCP servers used by Claude Code —
listing installed servers, searching npm for packages, and importing/exporting
MCP configuration.

**Module**: `github.com/gizzahub/gzh-cli-mcp-plugin` · **Binary**: `mcp-plugin` · **Go**: 1.24

> **Status**: Early stage. Architecture is Clean/Hexagonal (Ports & Adapters);
> commands and surface area are still evolving.

## Features

- **Server management** — list, install, remove, enable/disable MCP servers
- **Discovery** — search npm for MCP packages and inspect package info
- **Status & info** — check server status and show detailed server information
- **Configuration** — show, export, import, and validate MCP configuration

## Install

```bash
make build     # build the binary
make install   # install to $GOPATH/bin
make validate  # fmt + lint + test (pre-commit)
```

## Commands

### Servers & plugins

| Command                        | Purpose                              |
|--------------------------------|--------------------------------------|
| `mcp-plugin list`              | List MCP servers                     |
| `mcp-plugin install <name> [package]` | Install an MCP server         |
| `mcp-plugin remove <name>`     | Remove an MCP server                 |
| `mcp-plugin enable <plugin-id>`| Enable an MCP plugin                 |
| `mcp-plugin disable <plugin-id>`| Disable an MCP plugin               |
| `mcp-plugin server status [server]` | Check MCP server status         |
| `mcp-plugin server info <server>`   | Show detailed server information |
| `mcp-plugin server update [server]` | Update servers to latest version |

### Discovery

| Command                     | Purpose                               |
|-----------------------------|---------------------------------------|
| `mcp-plugin search <query>` | Search for MCP packages on npm        |
| `mcp-plugin info <package>` | Show information about an MCP package  |

### Configuration

| Command                        | Purpose                          |
|--------------------------------|----------------------------------|
| `mcp-plugin config show`       | Show current configuration       |
| `mcp-plugin config paths`      | Show configuration file paths    |
| `mcp-plugin config export`     | Export MCP configuration to file |
| `mcp-plugin config import <file>` | Import MCP configuration       |
| `mcp-plugin config validate`   | Validate MCP configuration       |

### Misc

| Command             | Purpose                  |
|---------------------|--------------------------|
| `mcp-plugin version`| Show version information  |

## Architecture

Clean Architecture + Hexagonal (Ports & Adapters):

```
cmd/mcp-plugin/       CLI (Cobra commands)
pkg/
  domain/             core entities & business rules
  application/        use cases / orchestration
  infrastructure/     adapters (config, npm, filesystem)
internal/version/     build version info
```

## Development

```bash
make validate       # fmt + lint + test (pre-commit)
make test           # run all tests
make test-coverage  # coverage report
make fmt            # format
make lint           # golangci-lint
make vet            # go vet
```
