package claudeconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReader_ReadSettings(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test settings.json
	settings := map[string]interface{}{
		"$schema": "https://json.schemastore.org/claude-code-settings.json",
		"enabledPlugins": map[string]bool{
			"context7@claude-plugins-official": true,
			"serena@claude-plugins-official":   false,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(tmpDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	reader := NewReader(tmpDir)
	result, err := reader.ReadSettings()
	if err != nil {
		t.Fatalf("ReadSettings() error = %v", err)
	}

	if len(result.EnabledPlugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(result.EnabledPlugins))
	}

	if !result.EnabledPlugins["context7@claude-plugins-official"] {
		t.Error("context7 should be enabled")
	}

	if result.EnabledPlugins["serena@claude-plugins-official"] {
		t.Error("serena should be disabled")
	}
}

func TestReader_ReadSettings_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	reader := NewReader(tmpDir)

	result, err := reader.ReadSettings()
	if err != nil {
		t.Fatalf("ReadSettings() error = %v", err)
	}

	if result.EnabledPlugins == nil {
		t.Error("EnabledPlugins should be initialized")
	}
}

func TestReader_EnableDisablePlugin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial settings
	settings := map[string]interface{}{
		"enabledPlugins": map[string]bool{
			"test@publisher": false,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(tmpDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	reader := NewReader(tmpDir)

	// Enable plugin
	if err := reader.EnablePlugin("test@publisher"); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}

	enabled, err := reader.IsPluginEnabled("test@publisher")
	if err != nil {
		t.Fatalf("IsPluginEnabled() error = %v", err)
	}
	if !enabled {
		t.Error("Plugin should be enabled")
	}

	// Disable plugin
	if err := reader.DisablePlugin("test@publisher"); err != nil {
		t.Fatalf("DisablePlugin() error = %v", err)
	}

	enabled, err = reader.IsPluginEnabled("test@publisher")
	if err != nil {
		t.Fatalf("IsPluginEnabled() error = %v", err)
	}
	if enabled {
		t.Error("Plugin should be disabled")
	}
}

func TestReader_ListPlugins(t *testing.T) {
	tmpDir := t.TempDir()

	settings := map[string]interface{}{
		"enabledPlugins": map[string]bool{
			"plugin1@pub": true,
			"plugin2@pub": false,
			"plugin3@pub": true,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(tmpDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	reader := NewReader(tmpDir)

	// Test ListEnabledPlugins
	enabled, err := reader.ListEnabledPlugins()
	if err != nil {
		t.Fatalf("ListEnabledPlugins() error = %v", err)
	}

	if len(enabled) != 2 {
		t.Errorf("Expected 2 enabled plugins, got %d", len(enabled))
	}

	// Test ListAllPlugins
	all, err := reader.ListAllPlugins()
	if err != nil {
		t.Fatalf("ListAllPlugins() error = %v", err)
	}

	if len(all) != 3 {
		t.Errorf("Expected 3 total plugins, got %d", len(all))
	}
}

func TestReader_PreservesUnknownFields(t *testing.T) {
	tmpDir := t.TempDir()

	// Create settings with extra fields
	settings := map[string]interface{}{
		"$schema":             "test-schema",
		"alwaysThinkingEnabled": true,
		"customField":         "preserved",
		"enabledPlugins": map[string]bool{
			"test@pub": false,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(tmpDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	reader := NewReader(tmpDir)

	// Enable a plugin
	if err := reader.EnablePlugin("test@pub"); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}

	// Read raw file and check custom field is preserved
	newData, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(newData, &result); err != nil {
		t.Fatal(err)
	}

	if result["customField"] != "preserved" {
		t.Error("customField should be preserved")
	}

	if result["alwaysThinkingEnabled"] != true {
		t.Error("alwaysThinkingEnabled should be preserved")
	}

	plugins := result["enabledPlugins"].(map[string]interface{})
	if plugins["test@pub"] != true {
		t.Error("plugin should be enabled")
	}
}
