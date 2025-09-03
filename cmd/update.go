package cmd

import (
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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update index file to reflect current chapter files",
	Long: `Update scans the context/ directory and regenerates the index file
to reflect the current state of chapter files.

This command:
1. Scans context/ directory for all chapter files
2. Analyzes file content to generate descriptions
3. Updates the index file with current chapter references
4. Maintains the lightweight table of contents

Use this command when you've added, removed, or modified chapter files
and need the index to reflect the current state.`,
	RunE: runUpdate,
}

var (
	updateTemplate string
	forceUpdate    bool
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateTemplate, "template", "claude",
		"Template type for index file (claude, cursor, copilot, generic)")
	updateCmd.Flags().BoolVar(&forceUpdate, "force", false,
		"Force update even if no changes detected")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	projectPath := getProjectPath(cmd)

	logVerbose(cmd, "Updating index in: %s", projectPath)
	logVerbose(cmd, "Using template: %s", updateTemplate)

	// Validate inputs
	if err := validation.ValidateDirectoryPath(projectPath); err != nil {
		return fmt.Errorf("invalid project path: %w", err)
	}

	if err := config.ValidateTemplate(updateTemplate); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Check if context directory exists
	contextDir := filepath.Join(projectPath, "context")
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return fmt.Errorf("context directory not found: %s\nRun 'contindex init' to set up the structure", contextDir)
	}

	// Scan for chapter files
	chapterFiles, err := scanContextDirectory(contextDir)
	if err != nil {
		return fmt.Errorf("failed to scan chapter files: %w", err)
	}

	if len(chapterFiles) == 0 {
		fmt.Printf("No chapter files found in %s\n", contextDir)
		fmt.Printf("Add .md files to the context/ directory and run update again\n")
		return nil
	}

	logVerbose(cmd, "Found %d chapter files", len(chapterFiles))

	// Get current index file path
	indexFile, err := config.GetMainFileForTemplate(updateTemplate, projectPath)
	if err != nil {
		return fmt.Errorf("failed to determine index file path: %w", err)
	}

	// Check if update is needed (unless forced)
	if !forceUpdate {
		if needsUpdate, err := checkIfUpdateNeeded(indexFile, chapterFiles); err != nil {
			logVerbose(cmd, "Warning: could not check update status: %v", err)
		} else if !needsUpdate {
			fmt.Printf("Index file is up to date. Use --force to regenerate anyway.\n")
			return nil
		}
	}

	// Generate updated index using template system
	projectConfig := config.DefaultConfig(projectPath)
	if err := projectConfig.UpdateForTemplate(updateTemplate); err != nil {
		return fmt.Errorf("failed to configure template: %w", err)
	}

	// Use template system to regenerate the index
	templateManager := template.New()
	if err := templateManager.ApplyTemplate(projectConfig); err != nil {
		return fmt.Errorf("failed to apply template: %w", err)
	}

	// Update template with chapter filenames (already semantic from AI)
	if err := UpdateTemplateWithChapters(indexFile, chapterFiles); err != nil {
		return fmt.Errorf("failed to update template with chapters: %w", err)
	}

	// Success message
	printUpdateSuccess(indexFile, chapterFiles)

	return nil
}

func checkIfUpdateNeeded(indexFile string, chapterFiles []*classifier.ContextFile) (bool, error) {
	// Check if index file exists
	indexStat, err := os.Stat(indexFile)
	if os.IsNotExist(err) {
		return true, nil // Index doesn't exist, update needed
	}
	if err != nil {
		return true, err // Can't check, assume update needed
	}

	// Check if any chapter file is newer than index
	for _, file := range chapterFiles {
		filePath := filepath.Join("context", file.FileName)
		fileStat, err := os.Stat(filePath)
		if err != nil {
			continue // Skip if file doesn't exist
		}
		if fileStat.ModTime().After(indexStat.ModTime()) {
			return true, nil // Chapter file is newer
		}
	}

	return false, nil // No update needed
}

func printUpdateSuccess(indexFile string, chapterFiles []*classifier.ContextFile) {
	fmt.Printf("✓ Successfully updated index file: %s\n\n", indexFile)

	fmt.Printf("Chapter files referenced:\n")
	for i, file := range chapterFiles {
		fmt.Printf("%d. context/%s\n", i+1, file.FileName)
	}

	fmt.Printf("\nTotal chapters: %d\n", len(chapterFiles))
	fmt.Printf("\nAI tools can now reference the updated index to load specific chapters.\n")
}

// scanContextDirectory scans the context directory for .md files
func scanContextDirectory(contextDir string) ([]*classifier.ContextFile, error) {
	entries, err := os.ReadDir(contextDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read context directory: %w", err)
	}

	var contextFiles []*classifier.ContextFile

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Simple ContextFile with just filename - no analysis needed
		contextFile := &classifier.ContextFile{
			FileName: entry.Name(),
		}

		contextFiles = append(contextFiles, contextFile)
	}

	return contextFiles, nil
}

