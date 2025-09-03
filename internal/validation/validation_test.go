package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid file path",
			path:    "test.md",
			wantErr: false,
		},
		{
			name:    "valid nested path",
			path:    "docs/test.md",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "file path cannot be empty",
		},
		{
			name:    "whitespace only path",
			path:    "   ",
			wantErr: true,
			errMsg:  "file path cannot be empty",
		},
		{
			name:    "path traversal attack",
			path:    "../../../etc/passwd",
			wantErr: true,
			errMsg:  "path traversal not allowed",
		},
		{
			name:    "dangerous semicolon",
			path:    "test;rm -rf /.md",
			wantErr: true,
			errMsg:  "dangerous character ';' not allowed",
		},
		{
			name:    "dangerous ampersand",
			path:    "test&whoami.md",
			wantErr: true,
			errMsg:  "dangerous character '&' not allowed",
		},
		{
			name:    "dangerous pipe",
			path:    "test|cat.md",
			wantErr: true,
			errMsg:  "dangerous character '|' not allowed",
		},
		{
			name:    "dangerous backtick",
			path:    "test`id`.md",
			wantErr: true,
			errMsg:  "dangerous character '`' not allowed",
		},
		{
			name:    "dangerous dollar",
			path:    "test$HOME.md",
			wantErr: true,
			errMsg:  "dangerous character '$' not allowed",
		},
		{
			name:    "path too long",
			path:    strings.Repeat("a", 256),
			wantErr: true,
			errMsg:  "path too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateFilePath() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateFileExists(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	testDir := filepath.Join(tmpDir, "testdir")

	// Create test file
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test directory
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid existing file",
			path:    testFile,
			wantErr: false,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.md"),
			wantErr: true,
			errMsg:  "file does not exist",
		},
		{
			name:    "directory instead of file",
			path:    testDir,
			wantErr: true,
			errMsg:  "path is a directory",
		},
		{
			name:    "invalid path characters",
			path:    "test;rm.md",
			wantErr: true,
			errMsg:  "dangerous character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileExists(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateFileExists() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateDirectoryPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid directory path",
			path:    "test/dir",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "directory path cannot be empty",
		},
		{
			name:    "path traversal",
			path:    "../../../etc",
			wantErr: true,
			errMsg:  "path traversal not allowed",
		},
		{
			name:    "dangerous characters",
			path:    "test;rm -rf /",
			wantErr: true,
			errMsg:  "dangerous character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectoryPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectoryPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateDirectoryPath() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateDirectoryWritable(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "writable")
	readOnlyDir := filepath.Join(tmpDir, "readonly")

	// Create writable directory
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create read-only directory (if possible)
	if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "writable directory",
			path:    testDir,
			wantErr: false,
		},
		{
			name:    "non-existent directory (should create)",
			path:    filepath.Join(tmpDir, "newdir"),
			wantErr: false,
		},
		{
			name:    "invalid path",
			path:    "test;rm.dir",
			wantErr: true,
			errMsg:  "dangerous character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectoryWritable(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectoryWritable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateDirectoryWritable() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateMarkdownFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid markdown file
	validMd := filepath.Join(tmpDir, "valid.md")
	if err := os.WriteFile(validMd, []byte("# Test\nSome content"), 0644); err != nil {
		t.Fatalf("Failed to create test markdown: %v", err)
	}

	// Create non-markdown file
	nonMd := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(nonMd, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create binary file with enough null bytes to trigger binary detection (>30% of first 512 bytes)
	binaryFile := filepath.Join(tmpDir, "binary.md")
	binaryContent := make([]byte, 512)
	// Fill with null bytes (more than 30% to trigger binary detection)
	for i := 0; i < 200; i++ {
		binaryContent[i] = 0 // null bytes
	}
	for i := 200; i < 512; i++ {
		binaryContent[i] = byte('a') // some text content
	}
	if err := os.WriteFile(binaryFile, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	// Create large file (over 50MB limit)
	largeFile := filepath.Join(tmpDir, "large.md")
	largeContent := strings.Repeat("# Large File\nContent line with more text to increase size\n", 2000000) // ~100MB
	if err := os.WriteFile(largeFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid markdown file",
			path:    validMd,
			wantErr: false,
		},
		{
			name:    "non-markdown extension",
			path:    nonMd,
			wantErr: true,
			errMsg:  "not a markdown file",
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.md"),
			wantErr: true,
			errMsg:  "file does not exist",
		},
		{
			name:    "binary content",
			path:    binaryFile,
			wantErr: true,
			errMsg:  "appears to be binary",
		},
		{
			name:    "file too large",
			path:    largeFile,
			wantErr: true,
			errMsg:  "file too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMarkdownFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMarkdownFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateMarkdownFile() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid project name",
			input:   "my-project",
			wantErr: false,
		},
		{
			name:    "valid with spaces",
			input:   "My Project Name",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			input:   "my_project_name",
			wantErr: false,
		},
		{
			name:    "valid alphanumeric",
			input:   "Project123",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errMsg:  "project name cannot be empty",
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
			errMsg:  "project name cannot be empty",
		},
		{
			name:    "invalid characters",
			input:   "my@project",
			wantErr: true,
			errMsg:  "invalid project name",
		},
		{
			name:    "too long",
			input:   strings.Repeat("a", 101),
			wantErr: true,
			errMsg:  "too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateProjectName() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid template name",
			input:   "claude",
			wantErr: false,
		},
		{
			name:    "valid with dash",
			input:   "my-template",
			wantErr: false,
		},
		{
			name:    "valid alphanumeric",
			input:   "template123",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errMsg:  "template name cannot be empty",
		},
		{
			name:    "invalid characters",
			input:   "my_template",
			wantErr: true,
			errMsg:  "invalid template name",
		},
		{
			name:    "spaces not allowed",
			input:   "my template",
			wantErr: true,
			errMsg:  "invalid template name",
		},
		{
			name:    "too long",
			input:   strings.Repeat("a", 51),
			wantErr: true,
			errMsg:  "too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateTemplateName() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateCategoryName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid category name",
			input:   "data-layer",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			input:   "api-layer-v2",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errMsg:  "category name cannot be empty",
		},
		{
			name:    "uppercase not allowed",
			input:   "Data-Layer",
			wantErr: true,
			errMsg:  "invalid category name",
		},
		{
			name:    "underscore not allowed",
			input:   "data_layer",
			wantErr: true,
			errMsg:  "invalid category name",
		},
		{
			name:    "spaces not allowed",
			input:   "data layer",
			wantErr: true,
			errMsg:  "invalid category name",
		},
		{
			name:    "too long",
			input:   strings.Repeat("a", 31),
			wantErr: true,
			errMsg:  "too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCategoryName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategoryName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateCategoryName() error = %v, expected to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestIsBinaryContent(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "empty content",
			content:  []byte{},
			expected: false,
		},
		{
			name:     "normal text",
			content:  []byte("This is normal text content"),
			expected: false,
		},
		{
			name:     "markdown content",
			content:  []byte("# Header\n\nSome **bold** text\n"),
			expected: false,
		},
		{
			name:     "content with tabs and newlines",
			content:  []byte("Line 1\n\tIndented line\r\nWindows line ending"),
			expected: false,
		},
		{
			name:     "binary content with many nulls",
			content:  append([]byte{0, 0, 0, 0, 0}, []byte("some text")...),
			expected: true,
		},
		{
			name:     "control characters",
			content:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBinaryContent(tt.content)
			if result != tt.expected {
				t.Errorf("isBinaryContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal filename",
			input:    "test-file",
			expected: "test-file",
		},
		{
			name:     "filename with spaces",
			input:    "my test file",
			expected: "my test file",
		},
		{
			name:     "filename with dangerous characters",
			input:    "test/file\\name:with*dangerous?chars",
			expected: "test-file-name-with-dangerous-chars",
		},
		{
			name:     "filename with quotes",
			input:    "\"quoted filename\"",
			expected: "quoted filename", // The quotes get replaced, then trimmed
		},
		{
			name:     "filename with multiple dashes",
			input:    "test---file---name",
			expected: "test-file-name",
		},
		{
			name:     "empty filename",
			input:    "",
			expected: "unnamed",
		},
		{
			name:     "only dangerous characters",
			input:    "/\\:*?\"<>|",
			expected: "unnamed",
		},
		{
			name:     "very long filename",
			input:    strings.Repeat("a", 150),
			expected: strings.Repeat("a", 100),
		},
		{
			name:     "filename starting and ending with dashes",
			input:    "---test-file---",
			expected: "test-file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFileName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFileName() = %v, want %v", result, tt.expected)
			}
		})
	}
}
