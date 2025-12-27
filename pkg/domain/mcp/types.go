// Package mcp defines domain types for MCP (Model Context Protocol) management.
package mcp

// ServerType represents the type of MCP server.
type ServerType string

const (
	// ServerTypeCommand is a command-based MCP server (npx, uvx).
	ServerTypeCommand ServerType = "command"
	// ServerTypeHTTP is an HTTP-based MCP server.
	ServerTypeHTTP ServerType = "http"
)

// Server represents an MCP server configuration.
type Server struct {
	// Name is the server identifier.
	Name string
	// Type is the server type (command or http).
	Type ServerType
	// Command is the executable command (for command type).
	Command string
	// Args are the command arguments (for command type).
	Args []string
	// URL is the server URL (for http type).
	URL string
	// Headers are HTTP headers (for http type).
	Headers map[string]string
	// Enabled indicates if the server is enabled.
	Enabled bool
}

// Plugin represents a Claude Code plugin.
type Plugin struct {
	// ID is the unique plugin identifier (e.g., "context7@claude-plugins-official").
	ID string
	// Name is the plugin name.
	Name string
	// Publisher is the plugin publisher.
	Publisher string
	// Description is the plugin description.
	Description string
	// Version is the installed version.
	Version string
	// Enabled indicates if the plugin is enabled.
	Enabled bool
	// MCPServers are the MCP servers provided by this plugin.
	MCPServers []Server
}

// PluginID represents a plugin identifier with name and publisher.
type PluginID struct {
	Name      string
	Publisher string
}

// String returns the plugin ID string (name@publisher).
func (p PluginID) String() string {
	return p.Name + "@" + p.Publisher
}
