package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/angelcodes95/contindex/internal/classifier"
	"github.com/angelcodes95/contindex/internal/config"
	"github.com/angelcodes95/contindex/internal/template"
	"github.com/angelcodes95/contindex/internal/validation"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert monolithic context files to index-chapter architecture",
	Long: `Convert transforms existing monolithic context files into the index-chapter model.

This command:
1. Analyzes your monolithic file and breaks it into semantic chapter files
2. Creates a lightweight index file that references the chapters  
3. Backs up your original file before conversion
4. Enables AI tools to selectively load chapters instead of everything

This enables selective chapter loading instead of processing everything.`,
	RunE: runConvert,
}

var (
	sourceFile   string
	templateType string
	backupDir    string
	contextDir   string
	projectName  string
	noBackup     bool
	force        bool
)

func init() {
	convertCmd.Flags().StringVar(&sourceFile, "source", "CLAUDE.md", "Source monolithic context file")
	convertCmd.Flags().StringVar(&templateType, "template", "claude", "Template type (claude, cursor, copilot, gemini, generic)")
	convertCmd.Flags().StringVar(&backupDir, "backup-dir", "backup", "Backup directory for original file")
	convertCmd.Flags().StringVar(&contextDir, "context-dir", "context", "Context directory name for chapter files")
	convertCmd.Flags().StringVar(&projectName, "project", "Project", "Project name for index generation")
	convertCmd.Flags().BoolVar(&noBackup, "no-backup", false, "Skip creating backup of original file")
	convertCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing context directory if it contains files")
	convertCmd.Flags().BoolP("dry-run", "d", false, "Preview changes without writing files")
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// If project name is still default, use directory name
	if projectName == "Project" {
		if wd, err := os.Getwd(); err == nil {
			projectName = filepath.Base(wd)
		}
	}

	if err := validateConvertInputs(); err != nil {
		return err
	}

	printConversionStatus(dryRun)

	if !dryRun && !noBackup {
		if err := createBackup(); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	contextFiles, err := analyzeAndGenerateFiles()
	if err != nil {
		return err
	}

	if dryRun {
		return previewConversion(contextFiles)
	}

	if err := executeConversion(contextFiles); err != nil {
		return err
	}

	printConversionSuccess(contextFiles)
	return nil
}

func validateConvertInputs() error {
	if err := validation.ValidateMarkdownFile(sourceFile); err != nil {
		return fmt.Errorf("invalid source file: %w", err)
	}

	if err := validation.ValidateTemplateName(templateType); err != nil {
		return fmt.Errorf("invalid template name: %w", err)
	}

	if err := config.ValidateTemplate(templateType); err != nil {
		return fmt.Errorf("unsupported template type: %w", err)
	}

	// Only validate backup directory if backups are enabled
	if !noBackup {
		if err := validation.ValidateDirectoryPath(backupDir); err != nil {
			return fmt.Errorf("invalid backup directory: %w", err)
		}
	}

	if err := validation.ValidateDirectoryPath(contextDir); err != nil {
		return fmt.Errorf("invalid context directory: %w", err)
	}

	// Check for context directory conflicts
	if err := checkDirectoryConflicts(); err != nil {
		return err
	}

	return nil
}

func checkDirectoryConflicts() error {
	// Check if context directory already exists and has files
	if stat, err := os.Stat(contextDir); err == nil && stat.IsDir() {
		files, err := os.ReadDir(contextDir)
		if err != nil {
			return fmt.Errorf("failed to read context directory: %w", err)
		}

		if len(files) > 0 && !force {
			return fmt.Errorf("context directory '%s' already exists and contains %d files - use --context-dir to specify a different name or --force to overwrite", contextDir, len(files))
		}

		if len(files) > 0 && force {
			fmt.Printf("Warning: Overwriting existing files in %s/ directory (--force enabled)\n", contextDir)
		}
	}

	// Check if backup directory conflicts (when backups enabled)
	if !noBackup {
		if contextDir == backupDir {
			return fmt.Errorf("context directory and backup directory cannot be the same ('%s')", contextDir)
		}
	}

	return nil
}

func printConversionStatus(dryRun bool) {
	if dryRun {
		fmt.Printf("DRY RUN: Analyzing %s for file-based structure using %s template...\n",
			sourceFile, templateType)
	} else {
		fmt.Printf("Converting %s to file-based contindex structure using %s template...\n",
			sourceFile, templateType)
	}
}

func analyzeAndGenerateFiles() ([]*classifier.ContextFile, error) {
	analyzer := classifier.NewFileAnalyzer(sourceFile)
	contextFiles, err := analyzer.AnalyzeAndGenerate(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to analyze and generate files: %w", err)
	}

	if len(contextFiles) == 0 {
		return nil, fmt.Errorf("no content sections found in source file")
	}

	return contextFiles, nil
}

func previewConversion(contextFiles []*classifier.ContextFile) error {
	fmt.Printf("\nPREVIEW: Would create %d context files:\n\n", len(contextFiles))

	totalTokens := 0
	for i, file := range contextFiles {
		fmt.Printf("%d. %s\n", i+1, file.FileName)
		fmt.Printf("   Summary: %s\n", file.Summary)
		fmt.Printf("   Size: %d words, ~%d tokens\n", file.WordCount, file.TokenCount)
		if len(file.KeyTerms) > 0 {
			fmt.Printf("   Key terms: %s\n", strings.Join(file.KeyTerms, ", "))
		}
		fmt.Printf("\n")
		totalTokens += file.TokenCount
	}

	fmt.Printf("Total estimated tokens: %d\n", totalTokens)
	fmt.Printf("Average tokens per file: %d\n", totalTokens/len(contextFiles))

	return nil
}

func executeConversion(contextFiles []*classifier.ContextFile) error {
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return fmt.Errorf("failed to create context directory: %w", err)
	}

	if err := writeContextFiles(contextFiles, contextDir); err != nil {
		return fmt.Errorf("failed to write context files: %w", err)
	}

	if err := generateIndexFile(contextFiles); err != nil {
		return fmt.Errorf("failed to generate index file: %w", err)
	}

	return nil
}

func writeContextFiles(contextFiles []*classifier.ContextFile, contextDir string) error {
	for _, file := range contextFiles {
		filePath := filepath.Join(contextDir, file.FileName)

		content := fmt.Sprintf("# %s\n\n%s\n",
			strings.TrimSuffix(file.FileName, ".md"), file.Content)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.FileName, err)
		}
	}

	return nil
}

func generateIndexFile(contextFiles []*classifier.ContextFile) error {
	// Use template system to create index file
	projectConfig := config.DefaultConfig(".")
	projectConfig.ContextDir = contextDir // Use configurable context directory
	if err := projectConfig.UpdateForTemplate(templateType); err != nil {
		return fmt.Errorf("failed to configure template: %w", err)
	}

	// Create template manager and generate index
	templateManager := template.New()
	if err := templateManager.ApplyTemplate(projectConfig); err != nil {
		return fmt.Errorf("failed to apply template: %w", err)
	}

	// Update template with AI-generated chapter information
	return UpdateTemplateWithChapters(projectConfig.MainFile, contextFiles, contextDir)
}

func createBackup() error {
	if err := validation.ValidateDirectoryWritable(backupDir); err != nil {
		return fmt.Errorf("backup directory validation failed: %w", err)
	}

	content, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	originalName := filepath.Base(sourceFile)
	safeName := validation.SanitizeFileName(originalName)
	backupFile := filepath.Join(backupDir, safeName)

	if _, err := os.Stat(backupFile); err == nil {
		ext := filepath.Ext(safeName)
		nameWithoutExt := strings.TrimSuffix(safeName, ext)
		backupFile = filepath.Join(backupDir, fmt.Sprintf("%s_backup_%d%s",
			nameWithoutExt, len(content), ext))
	}

	if err := os.WriteFile(backupFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	fmt.Printf("Created backup: %s\n", backupFile)
	return nil
}

func printConversionSuccess(contextFiles []*classifier.ContextFile) {
	totalWords := 0
	totalTokens := 0

	for _, file := range contextFiles {
		totalWords += file.WordCount
		totalTokens += file.TokenCount
	}

	fmt.Printf("\nSuccessfully converted %s to index-chapter architecture\n", sourceFile)
	fmt.Printf("Created %d chapter files in %s/ directory\n", len(contextFiles), contextDir)
	fmt.Printf("Total content: %d words, ~%d tokens\n", totalWords, totalTokens)
	fmt.Printf("Average per chapter: %d tokens\n", totalTokens/len(contextFiles))
	fmt.Printf("Index file: %s\n", getIndexFileName(templateType))
	if !noBackup {
		fmt.Printf("Backup saved in: %s/\n", backupDir)
	} else {
		fmt.Printf("Backup: skipped (--no-backup)\n")
	}

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Review generated chapter files in %s/ directory\n", contextDir)
	fmt.Printf("2. Check the index file - it references all chapters\n")
	fmt.Printf("3. AI tools can now load specific chapters instead of everything\n")
}

func getIndexFileName(templateType string) string {
	switch templateType {
	case "claude":
		return "CLAUDE.md"
	case "cursor":
		return "AGENTS.md"
	case "copilot":
		return "copilot-instructions.md"
	case "gemini":
		return "GEMINI.md"
	default:
		return "template.md"
	}
}

// UpdateTemplateWithChapters updates the template file with AI-generated chapter names
func UpdateTemplateWithChapters(mainFile string, contextFiles []*classifier.ContextFile, contextDirName string) error {
	// Read the current template file
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Generate simple chapter list - AI already created semantic names
	var chapterList strings.Builder
	for i, file := range contextFiles {
		// Use the AI-generated descriptive filename as the TOC entry
		descriptiveName := strings.TrimSuffix(file.FileName, ".md")
		chapterList.WriteString(fmt.Sprintf("%d. **%s** - `%s/%s`\n", i+1, descriptiveName, contextDirName, file.FileName))
	}

	// Replace placeholder with actual chapter list
	placeholders := []string{
		"(Chapter files will be listed here when you run `contindex update` or `contindex convert`)",
		"(Context files will be listed here when you run `contindex update` or `contindex convert`)",
	}

	updatedContent := string(content)
	for _, placeholder := range placeholders {
		updatedContent = strings.ReplaceAll(updatedContent, placeholder, strings.TrimSpace(chapterList.String()))
	}

	// Write updated content
	return os.WriteFile(mainFile, []byte(updatedContent), 0644)
}
