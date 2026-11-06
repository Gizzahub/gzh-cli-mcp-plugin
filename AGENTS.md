# gzh-cli-mcp-plugin

## Development workflow — user decision 2026-11-06

- This repository uses ordinary Git and its own Make targets. CE, `ce task`,
  `ce access`, task_manager, and task-manager are not prerequisites and must not
  be invoked for development, worktree lifecycle, task records, verification,
  integration, or cleanup. This product-specific user decision overrides
  inherited instructions requiring those tools or proposing their adoption.
- Use an isolated task branch/worktree; preserve existing work and commits.
  Use Git directly for status, commit, task-branch push, integration and cleanup.
- Run the repository's applicable source checks and keep failures visible.
  Task-board state and card schemas are not source-quality gates. Documentation
  and workflow-only changes require diff checks and relevant Markdown/configuration
  or workflow checks; do not require unrelated Go test suites for prose edits.
- GitLab CI is not required to build, test, or contribute. GitHub Actions and
  server-side branch rules are separate and must still be respected where present.
- Existing `tasks/` documents and historical receipts are retained as records,
  not as an executable task-manager database or a reason to block source work.
- Preserve in-file locked regions. Do not change machine-global installed tools
  or unrelated repositories as part of this repository's workflow.

This file provides guidance to Claude Code when working with code in this repository.

---

## Quick Start (30s scan)

**Binary**: `mcp-plugin` (MCP Plugin Management CLI)
**Status**: Initial stage - project structure only
**Architecture**: Clean Architecture + Hexagonal (Ports & Adapters)
**Go Version**: 1.26+

MCP plugin management tool for Claude Code integration:
- **List**: Show installed MCP servers
- **Enable/Disable**: Toggle MCP servers
- **Status**: Check MCP server status
- **Config**: Manage MCP configuration

---

## Top 10 Commands

| Command | Purpose | Usage |
|---------|---------|-------|
| `make build` | Build binary | Before testing |
| `make test` | Run all tests | Quick validation |
| `make validate` | fmt + lint + test | Pre-commit |
| `make fmt` | Format code | Before commit |
| `make lint` | Run golangci-lint | Fix issues |
| `make clean` | Clean artifacts | Fresh start |
| `make help` | Show all targets | Reference |

---

## Absolute Rules (DO/DON'T)

### DO
- Use `gzh-cli-core` for common utilities
- Follow Clean Architecture layers
- Run `make validate` before every commit
- Keep files < 300 lines (~10KB)
- Test coverage target: 80%+

### DON'T
- Import external libraries in domain layer
- Add CGO dependencies
- Bypass use cases (CLI -> infrastructure directly)
- Modify Claude Code internals directly

---

## Directory Structure

```
.
├── cmd/mcp-plugin/      # CLI entry (Presentation)
├── pkg/
│   ├── domain/          # Core logic (NO external deps)
│   │   └── mcp/         # MCP types and interfaces
│   ├── application/     # Use cases + ports
│   │   └── usecase/     # list, enable, disable, status
│   └── infrastructure/  # Adapters + repos
│       ├── claudeconfig/# ~/.claude/settings.json parser
│       └── plugincache/ # Plugin cache reader
├── internal/            # Private utilities
└── Makefile
```

---

## Claude Code Integration

### Configuration Files

- `~/.claude/settings.json` - Main settings with `enabledPlugins`
- `~/.claude/plugins/cache/{publisher}/{plugin}/{version}/` - Plugin cache
- `.claude-plugin/plugin.json` - Plugin metadata
- `.mcp.json` - MCP server configuration

### MCP Server Types

1. **Command-based**: `npx -y @package/mcp-server`
2. **Python-based**: `uvx --from git+... package start-mcp-server`
3. **HTTP-based**: Remote API endpoints

---

## Shared Library (gzh-cli-core)

```go
import (
    "github.com/gizzahub/gzh-cli-core/logger"
    "github.com/gizzahub/gzh-cli-core/errors"
    "github.com/gizzahub/gzh-cli-core/config"
)
```

---

## Git Commit Format

```
{type}({scope}): {imperative verb} {what}
```

**Scopes**: `domain`, `application`, `infrastructure`, `cli`, `build`, `docs`, `test`

---

## Project Status

- **Current Phase**: Initial - Project structure
- **MVP Scope**: list, status, enable/disable commands
- **Next Milestone**: Read-only operations

---

**Last Updated**: 2025-12-27
