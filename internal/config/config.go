package config

import (
	"fmt"
	"path/filepath"
)

// Error messages constants
const (
	ErrUnsupportedTemplate = "unsupported template"
)

// Note: Category-based organization has been replaced with semantic file naming
// Individual descriptively-named files are created instead of rigid categories

// SupportedTemplates defines available template types
var SupportedTemplates = []string{
	"generic",
	"claude",
	"cursor",
	"copilot",
	"gemini",
}

// TemplateConfigs maps template names to their file configurations
var TemplateConfigs = map[string]TemplateConfig{
	"generic": {
		MainFile: "template.md",
		SubDir:   "",
	},
	"claude": {
		MainFile: "CLAUDE.md",
		SubDir:   "",
	},
	"cursor": {
		MainFile: "AGENTS.md",
		SubDir:   "",
	},
	"copilot": {
		MainFile: "copilot-instructions.md",
		SubDir:   ".github",
	},
	"gemini": {
		MainFile: "GEMINI.md",
		SubDir:   "",
	},
}

// TemplateConfig defines the structure for template configurations
type TemplateConfig struct {
	MainFile string // The main context file name
	SubDir   string // Optional subdirectory (e.g., .github for copilot)
}

// ProjectConfig holds configuration for a contindex project
type ProjectConfig struct {
	ContextDir  string // Directory containing individual context files
	Template    string // Template type being used
	MainFile    string // Main context file path
	ProjectRoot string // Root directory of the project
}

// DefaultConfig creates a default project configuration
func DefaultConfig(projectRoot string) *ProjectConfig {
	return &ProjectConfig{
		ContextDir:  filepath.Join(projectRoot, "context"),
		Template:    "generic",
		MainFile:    filepath.Join(projectRoot, "template.md"),
		ProjectRoot: projectRoot,
	}
}

// ValidateTemplate checks if a template name is supported
func ValidateTemplate(template string) error {
	for _, supported := range SupportedTemplates {
		if template == supported {
			return nil
		}
	}
	return fmt.Errorf("%s: %s", ErrUnsupportedTemplate, template)
}

// ValidateCategory is deprecated - categories are no longer used
// Individual descriptively-named files are created instead

// GetMainFileForTemplate returns the appropriate main file name for a template
func GetMainFileForTemplate(template string, projectRoot string) (string, error) {
	if err := ValidateTemplate(template); err != nil {
		return "", err
	}

	config := TemplateConfigs[template]
	if config.SubDir != "" {
		return filepath.Join(projectRoot, config.SubDir, config.MainFile), nil
	}
	return filepath.Join(projectRoot, config.MainFile), nil
}

// UpdateForTemplate modifies a ProjectConfig to use a specific template
func (pc *ProjectConfig) UpdateForTemplate(template string) error {
	if err := ValidateTemplate(template); err != nil {
		return err
	}

	pc.Template = template
	mainFile, err := GetMainFileForTemplate(template, pc.ProjectRoot)
	if err != nil {
		return err
	}
	pc.MainFile = mainFile

	return nil
}
