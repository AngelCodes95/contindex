package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Constants for validation limits and thresholds
const (
	MaxPathLength         = 255
	MaxMarkdownFileSize   = 50 * 1024 * 1024 // 50MB
	MaxProjectNameLength  = 100
	MaxTemplateNameLength = 50
	MaxCategoryNameLength = 30
	MaxFileNameLength     = 100
	BinaryThreshold       = 0.3 // 30% control characters indicates binary
	CheckBytesLength      = 512 // Number of bytes to check for binary detection
)

// Regular expression patterns
var (
	validProjectNamePattern  = regexp.MustCompile(`^[a-zA-Z0-9\s_-]+$`)
	validTemplateNamePattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	validCategoryNamePattern = regexp.MustCompile(`^[a-z0-9-]+$`)
	multipleDashPattern      = regexp.MustCompile(`-+`)
)

// Common dangerous characters for path validation
var dangerousChars = []string{";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]"}

// validatePathCommon performs common path validation checks
func validatePathCommon(path, pathType string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("%s path cannot be empty", pathType)
	}

	// Prevent path traversal attacks
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	// Check for dangerous characters
	for _, char := range dangerousChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in path: %s", char, path)
		}
	}

	// Validate path length
	if len(path) > MaxPathLength {
		return fmt.Errorf("path too long (max %d characters): %s", MaxPathLength, path)
	}

	return nil
}

// ValidateFilePath ensures a file path is safe and accessible
func ValidateFilePath(path string) error {
	return validatePathCommon(path, "file")
}

// ValidateFileExists checks if a file exists and is readable
func ValidateFileExists(path string) error {
	if err := ValidateFilePath(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("cannot access file: %s (%w)", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// Check if file is readable
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file is not readable: %s (%w)", path, err)
	}
	file.Close()

	return nil
}

// ValidateDirectoryPath ensures a directory path is safe
func ValidateDirectoryPath(path string) error {
	return validatePathCommon(path, "directory")
}

// ValidateDirectoryWritable checks if a directory exists and is writable
func ValidateDirectoryWritable(path string) error {
	if err := ValidateDirectoryPath(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Try to create the directory to test writability
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("cannot create directory: %s (%w)", path, err)
			}
			return nil
		}
		return fmt.Errorf("cannot access directory: %s (%w)", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	// Test writability by creating a temporary file with proper cleanup
	tempFile := filepath.Join(path, ".contindex_write_test")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("directory is not writable: %s (%w)", path, err)
	}

	// Ensure cleanup happens even if there's an error
	defer func() {
		file.Close()
		os.Remove(tempFile)
	}()

	return nil
}

// ValidateMarkdownFile checks if a file is a valid markdown file
func ValidateMarkdownFile(path string) error {
	if err := ValidateFileExists(path); err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".md" && ext != ".markdown" {
		return fmt.Errorf("file is not a markdown file: %s", path)
	}

	// Check file size (reasonable limit for context files)
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot get file info: %s (%w)", path, err)
	}

	// Size limit for markdown files
	if info.Size() > MaxMarkdownFileSize {
		return fmt.Errorf("markdown file too large (max %dMB): %s", MaxMarkdownFileSize/(1024*1024), path)
	}

	// Basic content validation
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file: %s (%w)", path, err)
	}

	// Check for binary content
	if isBinaryContent(content) {
		return fmt.Errorf("file appears to be binary, not text: %s", path)
	}

	return nil
}

// ValidateProjectName ensures a project name is valid
func ValidateProjectName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check length
	if len(name) > MaxProjectNameLength {
		return fmt.Errorf("project name too long (max %d characters)", MaxProjectNameLength)
	}

	// Check for valid characters (alphanumeric, dash, underscore, space)
	if !validProjectNamePattern.MatchString(name) {
		return fmt.Errorf("invalid project name: must contain only letters, numbers, spaces, dashes, and underscores")
	}

	return nil
}

// ValidateTemplateName ensures a template name is valid
func ValidateTemplateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	// Check for valid characters (alphanumeric, dash only)
	if !validTemplateNamePattern.MatchString(name) {
		return fmt.Errorf("invalid template name: must contain only letters, numbers, and dashes")
	}

	// Check length
	if len(name) > MaxTemplateNameLength {
		return fmt.Errorf("template name too long (max %d characters)", MaxTemplateNameLength)
	}

	return nil
}

// ValidateCategoryName ensures a category name is valid
func ValidateCategoryName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("category name cannot be empty")
	}

	// Check for valid characters (lowercase alphanumeric, dash only)
	if !validCategoryNamePattern.MatchString(name) {
		return fmt.Errorf("invalid category name: must contain only lowercase letters, numbers, and dashes")
	}

	// Check length
	if len(name) > MaxCategoryNameLength {
		return fmt.Errorf("category name too long (max %d characters)", MaxCategoryNameLength)
	}

	return nil
}

// isBinaryContent checks if content appears to be binary
func isBinaryContent(content []byte) bool {
	// Simple heuristic: if more than threshold of first checkLen bytes are null or control characters
	checkLen := CheckBytesLength
	if len(content) < checkLen {
		checkLen = len(content)
	}

	if checkLen == 0 {
		return false
	}

	controlChars := 0
	for i := 0; i < checkLen; i++ {
		c := content[i]
		if c == 0 || (c < 32 && c != '\t' && c != '\n' && c != '\r') {
			controlChars++
		}
	}

	return float64(controlChars)/float64(checkLen) > BinaryThreshold
}

// SanitizeFileName creates a safe filename from user input
func SanitizeFileName(name string) string {
	// Replace dangerous characters with safe alternatives
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "*", "-")
	name = strings.ReplaceAll(name, "?", "-")
	name = strings.ReplaceAll(name, "\"", "-")
	name = strings.ReplaceAll(name, "<", "-")
	name = strings.ReplaceAll(name, ">", "-")
	name = strings.ReplaceAll(name, "|", "-")

	// Remove multiple consecutive dashes
	name = multipleDashPattern.ReplaceAllString(name, "-")

	// Trim dashes from start/end
	name = strings.Trim(name, "-")

	// Ensure reasonable length
	if len(name) > MaxFileNameLength {
		name = name[:MaxFileNameLength]
	}

	// Ensure not empty
	if name == "" {
		name = "unnamed"
	}

	return name
}
