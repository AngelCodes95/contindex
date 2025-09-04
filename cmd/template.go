package cmd

import (
	"fmt"
	"strings"
	"text/template"

	contindexTemplate "github.com/angelcodes95/contindex/internal/template"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage context file templates",
	Long: `Template command provides operations for managing context file templates.

Available subcommands:
  list  - Show all available templates
  show  - Display template details and content
  info  - Get detailed information about a specific template

Templates determine how the main context file is structured and what
reference syntax is used for different AI tools.`,
}

// templateListCmd lists available templates
var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long: `List displays all built-in templates available for context organization.

Each template is optimized for specific AI tools:
- generic: Universal template that can be adapted to any AI tool
- claude: Optimized for Claude Code with @context/ references
- cursor: Designed for Cursor IDE with folder icons  
- copilot: GitHub Copilot compatible with .github placement
- gemini: Optimized for Google Gemini conversational context loading`,
	RunE: runTemplateList,
}

// templateShowCmd shows template content
var templateShowCmd = &cobra.Command{
	Use:   "show <template-name>",
	Short: "Show template content",
	Long: `Show displays the raw template content for a specific template.
This is useful for understanding how templates are structured and
what variables are available for customization.`,
	Args: cobra.ExactArgs(1),
	RunE: runTemplateShow,
}

// templateInfoCmd shows template information
var templateInfoCmd = &cobra.Command{
	Use:   "info <template-name>",
	Short: "Show detailed template information",
	Long: `Info provides comprehensive information about a template including:
- Template description and use case
- Main context file name and location
- Compatible AI tools
- Reference syntax used`,
	Args: cobra.ExactArgs(1),
	RunE: runTemplateInfo,
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateInfoCmd)

	// Flags for template show
	templateShowCmd.Flags().BoolP("raw", "r", false,
		"Show raw template without processing")
}

func runTemplateList(cmd *cobra.Command, args []string) error {
	templateManager := contindexTemplate.New()
	templates := templateManager.ListTemplates()

	fmt.Printf("Available Templates\n\n")

	for _, templateName := range templates {
		info, err := templateManager.GetTemplateInfo(templateName)
		if err != nil {
			logVerbose(cmd, "Warning: could not get info for template %s: %v", templateName, err)
			fmt.Printf("   %s - (no description available)\n", templateName)
			continue
		}

		fmt.Printf("   %s - %s\n", templateName, info.Description)

		// Show main file info
		mainFile := info.MainFile
		if info.SubDir != "" {
			mainFile = fmt.Sprintf("%s/%s", info.SubDir, info.MainFile)
		}
		fmt.Printf("     File: %s\n", mainFile)
	}

	fmt.Println()
	fmt.Printf("Usage: contindex init --template=<name>\n")
	fmt.Printf("       contindex template show <name>\n")

	return nil
}

func runTemplateShow(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	raw, _ := cmd.Flags().GetBool("raw")

	templateManager := contindexTemplate.New()
	info, err := templateManager.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("template not found: %v", err)
	}

	fmt.Printf("Template: %s\n", templateName)
	fmt.Printf("Description: %s\n\n", info.Description)

	if raw {
		fmt.Println("--- Raw Template Content ---")
		fmt.Println(info.Content)
	} else {
		fmt.Println("--- Template Preview ---")
		preview, err := generateTemplatePreview(info)
		if err != nil {
			return fmt.Errorf("failed to generate preview: %v", err)
		}
		fmt.Println(preview)
	}

	fmt.Println("--- End Template ---")

	return nil
}

func runTemplateInfo(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	templateManager := contindexTemplate.New()
	info, err := templateManager.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("template not found: %v", err)
	}

	fmt.Printf("Template Information\n\n")
	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("Description: %s\n", info.Description)
	fmt.Printf("Main file: %s\n", info.MainFile)

	if info.SubDir != "" {
		fmt.Printf("Subdirectory: %s\n", info.SubDir)
		fmt.Printf("Full path: %s/%s\n", info.SubDir, info.MainFile)
	}

	// Show compatible AI tools
	fmt.Printf("\nCompatible AI Tools:\n")
	switch templateName {
	case "claude":
		fmt.Printf("   - Claude Code (primary)\n")
		fmt.Printf("   - Claude web interface\n")
		fmt.Printf("   - Any tool that supports @context/ references\n")
	case "cursor":
		fmt.Printf("   - Cursor IDE (primary)\n")
		fmt.Printf("   - VS Code with appropriate extensions\n")
	case "copilot":
		fmt.Printf("   - GitHub Copilot (primary)\n")
		fmt.Printf("   - GitHub Copilot for VS Code\n")
		fmt.Printf("   - GitHub Copilot CLI\n")
	case "generic":
		fmt.Printf("   - Any AI coding tool\n")
		fmt.Printf("   - Universal compatibility\n")
	}

	// Show reference syntax
	fmt.Printf("\nReference Syntax:\n")
	switch templateName {
	case "claude":
		fmt.Printf("   Individual files are referenced directly\n")
	case "cursor":
		fmt.Printf("   Individual files are referenced directly\n")
	case "copilot":
		fmt.Printf("   Individual files are referenced directly\n")
	case "generic":
		fmt.Printf("   Individual files are referenced directly\n")
	}

	// Show usage example
	fmt.Printf("\nUsage:\n")
	fmt.Printf("   contindex init --template=%s\n", templateName)
	fmt.Printf("   contindex update --template=%s\n", templateName)

	return nil
}

func generateTemplatePreview(info *contindexTemplate.Info) (string, error) {
	// Create sample template data
	sampleData := struct {
		ProjectName      string
		ProjectRoot      string
		ContextDir       string
		Categories       []struct{ Name, Description, Path string }
		Template         string
		GeneratedAt      string
		ContindexVersion string
		ReferenceSyntax  string
	}{
		ProjectName: "sample-project",
		ProjectRoot: "/path/to/project",
		ContextDir:  "/path/to/project/context",
		Categories: []struct{ Name, Description, Path string }{
			{"sample-file-1", "Example descriptively named file", "context/sample-file-1.md"},
			{"sample-file-2", "Another example file with semantic naming", "context/sample-file-2.md"},
			{"sample-file-3", "Third example showing file-based organization", "context/sample-file-3.md"},
		},
		Template:         info.Name,
		GeneratedAt:      "2024-01-01 12:00:00",
		ContindexVersion: "0.0.3",
		ReferenceSyntax:  "@context/%s/",
	}

	// Parse and execute template
	tmpl, err := template.New("preview").Parse(info.Content)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, sampleData); err != nil {
		return "", err
	}

	return result.String(), nil
}
