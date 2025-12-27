// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_SetPluginEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create initial settings
	settings := map[string]interface{}{
		"$schema": "test",
		"enabledPlugins": map[string]interface{}{
			"test@pub": false,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	writer := &Writer{homeDir: tmpDir}

	// Enable plugin
	if err := writer.SetPluginEnabled("test@pub", true); err != nil {
		t.Fatalf("SetPluginEnabled() error = %v", err)
	}

	// Verify
	enabled, exists, err := writer.GetPluginStatus("test@pub")
	if err != nil {
		t.Fatalf("GetPluginStatus() error = %v", err)
	}
	if !exists {
		t.Error("plugin should exist")
	}
	if !enabled {
		t.Error("plugin should be enabled")
	}
}

func TestWriter_ListPlugins(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	settings := map[string]interface{}{
		"enabledPlugins": map[string]interface{}{
			"plugin1@pub": true,
			"plugin2@pub": false,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	writer := &Writer{homeDir: tmpDir}
	plugins, err := writer.ListPlugins()
	if err != nil {
		t.Fatalf("ListPlugins() error = %v", err)
	}

	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(plugins))
	}
	if !plugins["plugin1@pub"] {
		t.Error("plugin1 should be enabled")
	}
	if plugins["plugin2@pub"] {
		t.Error("plugin2 should be disabled")
	}
}

func TestWriter_PluginExists(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	settings := map[string]interface{}{
		"enabledPlugins": map[string]interface{}{
			"exists@pub": true,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	writer := &Writer{homeDir: tmpDir}

	exists, err := writer.PluginExists("exists@pub")
	if err != nil {
		t.Fatalf("PluginExists() error = %v", err)
	}
	if !exists {
		t.Error("plugin should exist")
	}

	exists, err = writer.PluginExists("notexists@pub")
	if err != nil {
		t.Fatalf("PluginExists() error = %v", err)
	}
	if exists {
		t.Error("plugin should not exist")
	}
}

func TestWriter_GetPluginStatus(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	settings := map[string]interface{}{
		"enabledPlugins": map[string]interface{}{
			"enabled@pub":  true,
			"disabled@pub": false,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		t.Fatal(err)
	}

	writer := &Writer{homeDir: tmpDir}

	// Test enabled plugin
	enabled, exists, err := writer.GetPluginStatus("enabled@pub")
	if err != nil {
		t.Fatalf("GetPluginStatus() error = %v", err)
	}
	if !exists || !enabled {
		t.Error("enabled@pub should exist and be enabled")
	}

	// Test disabled plugin
	enabled, exists, err = writer.GetPluginStatus("disabled@pub")
	if err != nil {
		t.Fatalf("GetPluginStatus() error = %v", err)
	}
	if !exists || enabled {
		t.Error("disabled@pub should exist and be disabled")
	}

	// Test non-existent plugin
	_, exists, err = writer.GetPluginStatus("nonexistent@pub")
	if err != nil {
		t.Fatalf("GetPluginStatus() error = %v", err)
	}
	if exists {
		t.Error("nonexistent@pub should not exist")
	}
}
