package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/angelcodes95/contindex/internal/config"
)

// Manager handles template operations
type Manager struct{}

// New creates a new template manager
func New() *Manager {
	return &Manager{}
}

// Data holds data for template rendering
type Data struct {
	ProjectName      string
	ProjectRoot      string
	ContextDir       string
	Template         string
	GeneratedAt      string
	ContindexVersion string
	ReferenceSyntax  string
}

// ApplyTemplate creates the main context file using the specified template
func (m *Manager) ApplyTemplate(projectConfig *config.ProjectConfig) error {
	// Prepare template data
	templateData, err := m.prepareTemplateData(projectConfig)
	if err != nil {
		return fmt.Errorf("failed to prepare template data: %v", err)
	}

	// Get template content
	templateContent, err := m.getTemplateContent(projectConfig.Template)
	if err != nil {
		return fmt.Errorf("failed to get template content: %v", err)
	}

	// Parse and execute template
	tmpl, err := template.New("context").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Create main context file
	return m.writeContextFile(projectConfig.MainFile, tmpl, templateData)
}

// prepareTemplateData creates the data structure for template rendering
func (m *Manager) prepareTemplateData(projectConfig *config.ProjectConfig) (*Data, error) {
	projectName := filepath.Base(projectConfig.ProjectRoot)

	// Determine reference syntax based on template
	referenceSyntax := "@context/%s/"
	if projectConfig.Template == "cursor" {
		referenceSyntax = "context/%s/"
	}

	return &Data{
		ProjectName:      projectName,
		ProjectRoot:      projectConfig.ProjectRoot,
		ContextDir:       projectConfig.ContextDir,
		Template:         projectConfig.Template,
		GeneratedAt:      time.Now().Format("2006-01-02 15:04:05"),
		ContindexVersion: "0.0.3", // This should match the version in root.go
		ReferenceSyntax:  referenceSyntax,
	}, nil
}

// getTemplateContent retrieves the template content for the specified template type
func (m *Manager) getTemplateContent(templateType string) (string, error) {
	templatePath := fmt.Sprintf("templates/%s/template.md", templateType)

	content, err := TemplateFS.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("template not found: %s", templateType)
	}

	return string(content), nil
}

// writeContextFile writes the rendered template to the main context file
func (m *Manager) writeContextFile(filePath string, tmpl *template.Template, data *Data) error {
	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %v", err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create context file: %v", err)
	}
	defer file.Close()

	// Execute template and write to file
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}

// ListTemplates returns available template names
func (m *Manager) ListTemplates() []string {
	return config.SupportedTemplates
}

// GetTemplateInfo returns detailed information about a template
func (m *Manager) GetTemplateInfo(templateName string) (*Info, error) {
	if err := config.ValidateTemplate(templateName); err != nil {
		return nil, err
	}

	templateConfig := config.TemplateConfigs[templateName]

	// Read template content for preview
	content, err := m.getTemplateContent(templateName)
	if err != nil {
		return nil, err
	}

	return &Info{
		Name:        templateName,
		Description: getTemplateDescription(templateName),
		MainFile:    templateConfig.MainFile,
		SubDir:      templateConfig.SubDir,
		Content:     content,
	}, nil
}

// Info holds detailed information about a template
type Info struct {
	Name        string
	Description string
	MainFile    string
	SubDir      string
	Content     string
}

// Helper function to get template descriptions
func getTemplateDescription(templateName string) string {
	descriptions := map[string]string{
		"generic": "Universal template that can be adapted to any AI tool",
		"claude":  "Optimized for Claude Code with @context/ references",
		"cursor":  "Designed for Cursor IDE with folder icons",
		"copilot": "GitHub Copilot compatible with .github placement",
		"gemini":  "Optimized for Google Gemini conversational context loading",
	}

	if desc, exists := descriptions[templateName]; exists {
		return desc
	}
	return "No description available"
}
