package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/angelcodes95/contindex/internal/config"
	"github.com/angelcodes95/contindex/internal/template"
	"github.com/angelcodes95/contindex/internal/validation"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize index-chapter structure for fresh projects",
	Long: `Initialize sets up the index-chapter architecture for new projects.

This command creates:
1. Empty context/ directory for chapter files
2. Template index file (CLAUDE.md, AGENTS.md, etc.) with instructions
3. Basic structure ready for AI tools

For existing monolithic files, use 'contindex convert' instead.

Templates available:
  claude   - Optimized for Claude Code (creates CLAUDE.md)
  cursor   - Optimized for Cursor IDE (creates AGENTS.md)  
  copilot  - Optimized for GitHub Copilot (creates .github/copilot-instructions.md)
  generic  - Universal template (creates context-index.md)`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Template selection flag
	initCmd.Flags().StringP("template", "t", "generic",
		"Template type (generic, claude, cursor, copilot)")

	// Force flag for overwriting existing structure
	initCmd.Flags().BoolP("force", "f", false,
		"Force initialization even if structure already exists")
}

func runInit(cmd *cobra.Command, args []string) error {
	projectPath := getProjectPath(cmd)
	templateName, _ := cmd.Flags().GetString("template")
	force, _ := cmd.Flags().GetBool("force")

	logVerbose(cmd, "Initializing contindex in: %s", projectPath)
	logVerbose(cmd, "Using template: %s", templateName)

	// Comprehensive input validation
	if err := validation.ValidateDirectoryPath(projectPath); err != nil {
		return fmt.Errorf("invalid project path: %w", err)
	}

	if err := validation.ValidateTemplateName(templateName); err != nil {
		return fmt.Errorf("invalid template name: %w", err)
	}

	if err := config.ValidateTemplate(templateName); err != nil {
		return fmt.Errorf("unsupported template: %w", err)
	}

	// Validate project directory is writable
	if err := validation.ValidateDirectoryWritable(projectPath); err != nil {
		return fmt.Errorf("project directory not writable: %w", err)
	}

	// Create project configuration
	projectConfig := config.DefaultConfig(projectPath)
	if err := projectConfig.UpdateForTemplate(templateName); err != nil {
		return fmt.Errorf("failed to configure template: %v", err)
	}

	logVerbose(cmd, "Project config created: %+v", projectConfig)

	// Check if structure already exists
	if !force {
		if err := checkExistingStructure(projectConfig); err != nil {
			return err
		}
	}

	// Create directory structure
	if err := createDirectoryStructure(cmd, projectConfig); err != nil {
		return fmt.Errorf("failed to create directory structure: %v", err)
	}

	// Create main context file from template
	templateManager := template.New()
	if err := templateManager.ApplyTemplate(projectConfig); err != nil {
		return fmt.Errorf("failed to create context file from template: %v", err)
	}

	// Success message with next steps
	printInitSuccessMessage(projectConfig)

	return nil
}

func checkExistingStructure(config *config.ProjectConfig) error {
	// Check if context directory exists
	if _, err := os.Stat(config.ContextDir); err == nil {
		return fmt.Errorf("context directory already exists: %s\nUse --force to overwrite",
			config.ContextDir)
	}

	// Check if main context file exists
	if _, err := os.Stat(config.MainFile); err == nil {
		return fmt.Errorf("main context file already exists: %s\nUse --force to overwrite",
			config.MainFile)
	}

	return nil
}

func createDirectoryStructure(cmd *cobra.Command, projectConfig *config.ProjectConfig) error {
	// Create main context directory
	logVerbose(cmd, "Creating context directory: %s", projectConfig.ContextDir)
	if err := os.MkdirAll(projectConfig.ContextDir, 0755); err != nil {
		return fmt.Errorf("failed to create context directory: %v", err)
	}

	// Create .gitkeep file to ensure empty directory is tracked
	gitkeepPath := filepath.Join(projectConfig.ContextDir, ".gitkeep")
	if err := createGitkeepFile(gitkeepPath); err != nil {
		logVerbose(cmd, "Warning: could not create .gitkeep in context directory: %v", err)
	}

	// Create subdirectory for main file if needed (e.g., .github for copilot)
	templateConfig := config.TemplateConfigs[projectConfig.Template]
	if templateConfig.SubDir != "" {
		subDirPath := filepath.Join(projectConfig.ProjectRoot, templateConfig.SubDir)
		logVerbose(cmd, "Creating subdirectory for template: %s", subDirPath)
		if err := os.MkdirAll(subDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create template subdirectory: %v", err)
		}
	}

	return nil
}

func createGitkeepFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func printInitSuccessMessage(projectConfig *config.ProjectConfig) {
	fmt.Printf("✓ Successfully initialized contindex structure\n\n")

	fmt.Printf("Created:\n")
	fmt.Printf("  %s/     # Directory for individual context files\n", projectConfig.ContextDir)
	fmt.Printf("  %s     # Main context index file\n\n", projectConfig.MainFile)

	fmt.Printf("Next steps:\n")
	fmt.Printf("1. Use 'contindex convert --source=YOUR_FILE.md' to convert existing monolithic files\n")
	fmt.Printf("2. Or manually add descriptively-named .md files to the context/ directory\n")
	fmt.Printf("3. Start using your AI tool with selective file loading\n\n")

	fmt.Printf("Template: %s\n", projectConfig.Template)
	fmt.Printf("AI tools can now load specific files instead of processing everything.\n")
}
