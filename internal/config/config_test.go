package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid claude template",
			template: "claude",
			wantErr:  false,
		},
		{
			name:     "valid cursor template",
			template: "cursor",
			wantErr:  false,
		},
		{
			name:     "valid copilot template",
			template: "copilot",
			wantErr:  false,
		},
		{
			name:     "valid generic template",
			template: "generic",
			wantErr:  false,
		},
		{
			name:        "invalid template",
			template:    "invalid",
			wantErr:     true,
			errContains: "unsupported template",
		},
		{
			name:        "empty template",
			template:    "",
			wantErr:     true,
			errContains: "unsupported template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplate(tt.template)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateTemplate() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateTemplate() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTemplate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestGetMainFileForTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		projectRoot string
		want        string
		wantErr     bool
	}{
		{
			name:        "claude template",
			template:    "claude",
			projectRoot: "/test/project",
			want:        "/test/project/CLAUDE.md",
			wantErr:     false,
		},
		{
			name:        "cursor template",
			template:    "cursor",
			projectRoot: "/test/project",
			want:        "/test/project/AGENTS.md",
			wantErr:     false,
		},
		{
			name:        "copilot template with subdir",
			template:    "copilot",
			projectRoot: "/test/project",
			want:        "/test/project/.github/copilot-instructions.md",
			wantErr:     false,
		},
		{
			name:        "generic template",
			template:    "generic",
			projectRoot: "/test/project",
			want:        "/test/project/template.md",
			wantErr:     false,
		},
		{
			name:        "invalid template",
			template:    "invalid",
			projectRoot: "/test/project",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMainFileForTemplate(tt.template, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetMainFileForTemplate() expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("GetMainFileForTemplate() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("GetMainFileForTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	projectRoot := "/test/project"
	config := DefaultConfig(projectRoot)

	if config.ProjectRoot != projectRoot {
		t.Errorf("DefaultConfig() ProjectRoot = %v, want %v", config.ProjectRoot, projectRoot)
	}

	expectedContextDir := filepath.Join(projectRoot, "context")
	if config.ContextDir != expectedContextDir {
		t.Errorf("DefaultConfig() ContextDir = %v, want %v", config.ContextDir, expectedContextDir)
	}

	if config.Template != "generic" {
		t.Errorf("DefaultConfig() Template = %v, want %v", config.Template, "generic")
	}

	expectedMainFile := filepath.Join(projectRoot, "template.md")
	if config.MainFile != expectedMainFile {
		t.Errorf("DefaultConfig() MainFile = %v, want %v", config.MainFile, expectedMainFile)
	}
}

func TestProjectConfig_UpdateForTemplate(t *testing.T) {
	config := DefaultConfig("/test/project")

	err := config.UpdateForTemplate("claude")
	if err != nil {
		t.Errorf("UpdateForTemplate() unexpected error = %v", err)
		return
	}

	if config.Template != "claude" {
		t.Errorf("UpdateForTemplate() Template = %v, want %v", config.Template, "claude")
	}

	expectedMainFile := "/test/project/CLAUDE.md"
	if config.MainFile != expectedMainFile {
		t.Errorf("UpdateForTemplate() MainFile = %v, want %v", config.MainFile, expectedMainFile)
	}
}
