package plugin

import (
	"testing"
)

// mockRepo is a mock implementation of PluginRepository.
type mockRepo struct {
	plugins map[string]bool
	err     error
}

func (m *mockRepo) ListAll() (map[string]bool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.plugins, nil
}

func (m *mockRepo) ListEnabled() ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	var enabled []string
	for id, e := range m.plugins {
		if e {
			enabled = append(enabled, id)
		}
	}
	return enabled, nil
}

func (m *mockRepo) IsEnabled(pluginID string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.plugins[pluginID], nil
}

func (m *mockRepo) Enable(pluginID string) error {
	if m.err != nil {
		return m.err
	}
	m.plugins[pluginID] = true
	return nil
}

func (m *mockRepo) Disable(pluginID string) error {
	if m.err != nil {
		return m.err
	}
	m.plugins[pluginID] = false
	return nil
}

func TestListUseCase_Execute(t *testing.T) {
	repo := &mockRepo{
		plugins: map[string]bool{
			"context7@claude-plugins-official": true,
			"serena@claude-plugins-official":   false,
			"greptile@claude-plugins-official": true,
		},
	}

	uc := NewListUseCase(repo)

	// Test list all
	result, err := uc.Execute(false)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 plugins, got %d", len(result))
	}

	// Test enabled only
	result, err = uc.Execute(true)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 enabled plugins, got %d", len(result))
	}
}

func TestParsePluginID(t *testing.T) {
	tests := []struct {
		id            string
		wantName      string
		wantPublisher string
	}{
		{"context7@claude-plugins-official", "context7", "claude-plugins-official"},
		{"serena@claude-plugins-official", "serena", "claude-plugins-official"},
		{"no-publisher", "no-publisher", ""},
		{"multi@at@signs", "multi@at", "signs"},
	}

	for _, tt := range tests {
		name, publisher := parsePluginID(tt.id)
		if name != tt.wantName {
			t.Errorf("parsePluginID(%q) name = %q, want %q", tt.id, name, tt.wantName)
		}
		if publisher != tt.wantPublisher {
			t.Errorf("parsePluginID(%q) publisher = %q, want %q", tt.id, publisher, tt.wantPublisher)
		}
	}
}
