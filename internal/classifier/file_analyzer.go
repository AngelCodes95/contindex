package classifier

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/angelcodes95/contindex/internal/validation"
)

// Constants for the new file-based architecture
const (
	MinWordCountForFile  = 10 // Minimum words for a standalone file
	MaxDescriptiveLength = 60 // Maximum length for descriptive filenames
	TokenEstimationRatio = 4  // Estimated characters per token
)

// ContentSection represents a section of content with metadata
type ContentSection struct {
	Title     string // Section title extracted from headers
	Content   string // The actual content text
	StartLine int    // Starting line number in source file
	EndLine   int    // Ending line number in source file
	WordCount int    // Word count for this section
}

// ContextFile represents a single context file with descriptive naming
type ContextFile struct {
	FileName   string   // Descriptive filename based on content
	Content    string   // The actual file content
	WordCount  int      // Word count for the file
	TokenCount int      // Estimated token count
	Summary    string   // Brief content summary for indexing
	KeyTerms   []string // Key terms extracted from content
}

// FileAnalyzer processes monolithic files and generates descriptive individual files
type FileAnalyzer struct {
	SourceFile   string            // Path to source monolithic file
	content      string            // Cached source content
	sections     []*ContentSection // Parsed sections from source
	contextFiles []*ContextFile    // Generated context files
}

// New creates a new FileAnalyzer instance
func NewFileAnalyzer(sourceFile string) *FileAnalyzer {
	return &FileAnalyzer{
		SourceFile: sourceFile,
	}
}

// AnalyzeAndGenerate performs complete analysis and generates individual context files
func (fa *FileAnalyzer) AnalyzeAndGenerate(ctx context.Context) ([]*ContextFile, error) {
	// Validate source file
	if err := validation.ValidateMarkdownFile(fa.SourceFile); err != nil {
		return nil, fmt.Errorf("invalid source file: %w", err)
	}

	// Parse the source file into sections
	if err := fa.parseSourceFile(); err != nil {
		return nil, fmt.Errorf("failed to parse source file: %w", err)
	}

	// Generate individual context files
	if err := fa.generateContextFiles(); err != nil {
		return nil, fmt.Errorf("failed to generate context files: %w", err)
	}

	return fa.contextFiles, nil
}

// parseSourceFile reads and parses the monolithic file into content sections
func (fa *FileAnalyzer) parseSourceFile() error {
	content, err := os.ReadFile(fa.SourceFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fa.content = string(content)
	scanner := bufio.NewScanner(strings.NewReader(fa.content))

	var sections []*ContentSection
	var currentSection *ContentSection
	var contentBuffer strings.Builder
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Detect markdown headers (## or ###)
		if strings.HasPrefix(strings.TrimSpace(line), "##") {
			// Save previous section
			if currentSection != nil {
				currentSection.Content = strings.TrimSpace(contentBuffer.String())
				currentSection.EndLine = lineNum - 1
				currentSection.WordCount = len(strings.Fields(currentSection.Content))

				// Only keep sections with meaningful content
				if currentSection.WordCount >= MinWordCountForFile {
					sections = append(sections, currentSection)
				}
			}

			// Start new section
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			currentSection = &ContentSection{
				Title:     title,
				StartLine: lineNum,
			}
			contentBuffer.Reset()
		} else if currentSection != nil {
			contentBuffer.WriteString(line + "\n")
		}
	}

	// Handle final section
	if currentSection != nil {
		currentSection.Content = strings.TrimSpace(contentBuffer.String())
		currentSection.EndLine = lineNum
		currentSection.WordCount = len(strings.Fields(currentSection.Content))

		if currentSection.WordCount >= MinWordCountForFile {
			sections = append(sections, currentSection)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	fa.sections = sections
	return nil
}

// generateContextFiles creates individual context files with descriptive names
func (fa *FileAnalyzer) generateContextFiles() error {
	var contextFiles []*ContextFile

	for _, section := range fa.sections {
		// Generate descriptive filename based on content analysis
		fileName := fa.generateDescriptiveFileName(section)

		// Extract key terms for indexing
		keyTerms := fa.extractKeyTerms(section)

		// Generate content summary
		summary := fa.generateContentSummary(section)

		// Estimate token count
		tokenCount := len(section.Content) / TokenEstimationRatio

		contextFile := &ContextFile{
			FileName:   fileName,
			Content:    section.Content,
			WordCount:  section.WordCount,
			TokenCount: tokenCount,
			Summary:    summary,
			KeyTerms:   keyTerms,
		}

		contextFiles = append(contextFiles, contextFile)
	}

	fa.contextFiles = contextFiles
	return nil
}

// generateDescriptiveFileName creates meaningful filenames based on content analysis
func (fa *FileAnalyzer) generateDescriptiveFileName(section *ContentSection) string {
	content := strings.ToLower(section.Title + " " + section.Content)

	// Extract domain-specific terms
	var descriptors []string

	// Add title-based descriptor
	titleDesc := fa.extractTitleDescriptor(section.Title)
	if titleDesc != "" {
		descriptors = append(descriptors, titleDesc)
	}

	// Add technology-based descriptors
	techDesc := fa.extractTechnologyDescriptor(content)
	if techDesc != "" {
		descriptors = append(descriptors, techDesc)
	}

	// Add function-based descriptors
	funcDesc := fa.extractFunctionDescriptor(content)
	if funcDesc != "" {
		descriptors = append(descriptors, funcDesc)
	}

	// Combine descriptors into filename
	fileName := strings.Join(descriptors, "-")

	// Sanitize and validate filename
	fileName = validation.SanitizeFileName(fileName)

	// Ensure reasonable length
	if len(fileName) > MaxDescriptiveLength {
		fileName = fileName[:MaxDescriptiveLength]
	}

	// Default if empty
	if fileName == "" {
		fileName = "general-context"
	}

	return fileName + ".md"
}

// extractTitleDescriptor extracts meaningful terms from section title
func (fa *FileAnalyzer) extractTitleDescriptor(title string) string {
	title = strings.ToLower(title)

	// Remove common stop words
	stopWords := []string{"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"}
	words := strings.Fields(title)

	var meaningful []string
	for _, word := range words {
		isStopWord := false
		for _, stopWord := range stopWords {
			if word == stopWord {
				isStopWord = true
				break
			}
		}
		if !isStopWord && len(word) > 2 {
			meaningful = append(meaningful, word)
		}
	}

	// Take first 2-3 meaningful words
	if len(meaningful) > 3 {
		meaningful = meaningful[:3]
	}

	return strings.Join(meaningful, "-")
}

// extractTechnologyDescriptor identifies technology-specific terms
func (fa *FileAnalyzer) extractTechnologyDescriptor(content string) string {
	techPatterns := map[string]string{
		"postgresql|postgres|pg": "postgresql",
		"mongodb|mongo":          "mongodb",
		"redis":                  "redis",
		"kubernetes|k8s":         "kubernetes",
		"docker":                 "docker",
		"jwt|oauth":              "oauth",
		"stripe|payment":         "payments",
		"webhook":                "webhooks",
		"graphql|gql":            "graphql",
		"rest|api":               "rest-api",
	}

	for pattern, descriptor := range techPatterns {
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			return descriptor
		}
	}

	return ""
}

// extractFunctionDescriptor identifies functional aspects
func (fa *FileAnalyzer) extractFunctionDescriptor(content string) string {
	funcPatterns := map[string]string{
		"authentication|auth|login":      "authentication",
		"authorization|permission":       "authorization",
		"database|schema|model":          "database",
		"deployment|deploy|production":   "deployment",
		"monitoring|metrics|logging":     "monitoring",
		"security|encryption|compliance": "security",
		"testing|test|spec":              "testing",
		"configuration|config|setup":     "configuration",
	}

	for pattern, descriptor := range funcPatterns {
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			return descriptor
		}
	}

	return ""
}

// extractKeyTerms identifies important terms for indexing
func (fa *FileAnalyzer) extractKeyTerms(section *ContentSection) []string {
	content := strings.ToLower(section.Content)

	// Define important term patterns
	termPatterns := []string{
		"api", "endpoint", "database", "schema", "authentication", "authorization",
		"security", "deployment", "monitoring", "testing", "configuration",
		"jwt", "oauth", "postgresql", "mongodb", "redis", "kubernetes", "docker",
	}

	var foundTerms []string
	for _, term := range termPatterns {
		if strings.Contains(content, term) {
			foundTerms = append(foundTerms, term)
		}
	}

	return foundTerms
}

// generateContentSummary creates a brief summary for indexing
func (fa *FileAnalyzer) generateContentSummary(section *ContentSection) string {
	content := section.Content

	// Take first sentence or first 100 characters
	sentences := strings.Split(content, ".")
	if len(sentences) > 0 && len(sentences[0]) > 10 {
		summary := strings.TrimSpace(sentences[0])
		if len(summary) > 100 {
			summary = summary[:100] + "..."
		}
		return summary
	}

	// Fallback to first 100 characters
	if len(content) > 100 {
		return content[:100] + "..."
	}

	return content
}
