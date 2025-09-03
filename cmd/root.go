package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
const Version = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "contindex",
	Short: "Transform monolithic AI context into index-chapter architecture",
	Long: `Contindex uses an index-chapter model to solve AI context dilution.
Like a book with table of contents and chapters, it creates:

• Index file (CLAUDE.md) - lightweight table of contents with chapter references
• Chapter files - individual semantic content files that can be selectively loaded
• Backup system - preserves original monolithic files during conversion

This architecture enables AI tools to reference the index and load only relevant 
chapters instead of processing entire monolithic files.

Commands:
  init     - Set up contindex structure for fresh projects  
  convert  - Transform existing monolithic files to index-chapter format

Examples:
  contindex init --template=claude         # Fresh project setup
  contindex convert --source=CLAUDE.md     # Convert monolithic file
  contindex convert --dry-run              # Preview conversion`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("path", "p", ".", "Project directory path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	// Version flag
	rootCmd.Flags().BoolP("version", "", false, "Show version information")
}

// Helper function to get project path from flags or current directory
func getProjectPath(cmd *cobra.Command) string {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return "."
	}
	return path
}

// Helper function to check verbose flag
func isVerbose(cmd *cobra.Command) bool {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return false
	}
	return verbose
}

// Helper function for verbose logging
func logVerbose(cmd *cobra.Command, format string, args ...interface{}) {
	if isVerbose(cmd) {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}
