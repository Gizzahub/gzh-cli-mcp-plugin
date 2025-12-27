package plugin

import (
	"testing"
)

func TestToggleUseCase_Enable(t *testing.T) {
	tests := []struct {
		name        string
		plugins     map[string]bool
		pluginID    string
		wantEnabled bool
		wantAlready bool
	}{
		{
			name:        "enable disabled plugin",
			plugins:     map[string]bool{"test@pub": false},
			pluginID:    "test@pub",
			wantEnabled: true,
			wantAlready: false,
		},
		{
			name:        "enable already enabled plugin",
			plugins:     map[string]bool{"test@pub": true},
			pluginID:    "test@pub",
			wantEnabled: true,
			wantAlready: true,
		},
		{
			name:        "enable new plugin",
			plugins:     map[string]bool{},
			pluginID:    "new@pub",
			wantEnabled: true,
			wantAlready: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{plugins: tt.plugins}
			uc := NewToggleUseCase(repo)

			result, err := uc.Enable(tt.pluginID)
			if err != nil {
				t.Fatalf("Enable() error = %v", err)
			}

			if result.Enabled != tt.wantEnabled {
				t.Errorf("Enable().Enabled = %v, want %v", result.Enabled, tt.wantEnabled)
			}

			if result.WasAlready != tt.wantAlready {
				t.Errorf("Enable().WasAlready = %v, want %v", result.WasAlready, tt.wantAlready)
			}
		})
	}
}

func TestToggleUseCase_Disable(t *testing.T) {
	tests := []struct {
		name        string
		plugins     map[string]bool
		pluginID    string
		wantEnabled bool
		wantAlready bool
	}{
		{
			name:        "disable enabled plugin",
			plugins:     map[string]bool{"test@pub": true},
			pluginID:    "test@pub",
			wantEnabled: false,
			wantAlready: false,
		},
		{
			name:        "disable already disabled plugin",
			plugins:     map[string]bool{"test@pub": false},
			pluginID:    "test@pub",
			wantEnabled: false,
			wantAlready: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{plugins: tt.plugins}
			uc := NewToggleUseCase(repo)

			result, err := uc.Disable(tt.pluginID)
			if err != nil {
				t.Fatalf("Disable() error = %v", err)
			}

			if result.Enabled != tt.wantEnabled {
				t.Errorf("Disable().Enabled = %v, want %v", result.Enabled, tt.wantEnabled)
			}

			if result.WasAlready != tt.wantAlready {
				t.Errorf("Disable().WasAlready = %v, want %v", result.WasAlready, tt.wantAlready)
			}
		})
	}
}

func TestToggleUseCase_Status(t *testing.T) {
	repo := &mockRepo{
		plugins: map[string]bool{
			"enabled@pub":  true,
			"disabled@pub": false,
		},
	}
	uc := NewToggleUseCase(repo)

	// Test enabled plugin
	result, err := uc.Status("enabled@pub")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if !result.Enabled {
		t.Error("Expected plugin to be enabled")
	}

	// Test disabled plugin
	result, err = uc.Status("disabled@pub")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if result.Enabled {
		t.Error("Expected plugin to be disabled")
	}
}
