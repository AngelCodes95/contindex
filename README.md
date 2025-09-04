# Contindex

A lightweight CLI tool that transforms monolithic AI context files into an index-chapter architecture. 

Similar to a book with table of contents and chapters, it enables selective loading of context for better AI development workflows, with AI deciding the names based on context file context for optimal matching.

## A note from the Dev
This repository was both manually coded (now known as the old fashioned way) and AI pair programmed for revisions and refinements, serving as a showcase for my AI Pair Programming output standards. Read to the end for a fun fact about this project!

## Quick Start

### Installation

**Option 1: Install with Go (Recommended)**
```bash
go install github.com/angelcodes95/contindex@latest
```

**Option 2: Download Binary**
```bash
# Visit GitHub releases page and download for your platform
# https://github.com/angelcodes95/contindex/releases
```

**Option 3: Build from Source**
```bash
git clone https://github.com/angelcodes95/contindex.git
cd contindex
go build -o contindex .
# Move to PATH or run ./contindex
```

### Basic Usage

**For new projects:**
```bash
# Create organized context structure
contindex init --template=claude    # For Claude Code
contindex init --template=cursor    # For Cursor IDE  
contindex init --template=copilot   # For GitHub Copilot
contindex init --template=gemini    # For Google Gemini
contindex init --template=generic   # Universal template
```

**For existing large context files:**
```bash
# Convert monolithic file to organized chapters
contindex convert --source=CLAUDE.md

# Preview what would be created
contindex convert --source=CLAUDE.md --dry-run

# Customize directories to avoid conflicts
contindex convert --source=CLAUDE.md --context-dir=docs-context --backup-dir=backups

# Skip backup creation (not recommended)
contindex convert --source=CLAUDE.md --no-backup

# Overwrite existing context directory
contindex convert --source=CLAUDE.md --force
```

**After adding/removing chapters:**
```bash
# Updates the index file and creates a semantically aligned name that reflects the chapter contents
contindex update
```

## The Problem

AI coding tools suffer from context dilution when context files stay monolithic or unorganized. As projects grow, these files become:
- Too large for AI tools to process effectively (you probably notice it ready 50 lines or so after a while)
- Filled with irrelevant information for specific tasks
- Difficult to maintain and navigate
- A bottleneck for development velocity

## The Index-Chapter Solution

Contindex uses a book-like architecture that creates:

```
example using Claude Code:

Before: CLAUDE.md (monolithic file)
After:  CLAUDE.md (index) + context/chapter-files.md
```

**Architecture Components:**
- **Index file** (CLAUDE.md) - lightweight table of contents with chapter references
- **Chapter files** - individual semantically named context markdown files in context/ directory
- **Selective loading** - AI to read main md file, (CLAUDE.md for the claude code example) to read the file names and load only relevant files
- **Backup system** - preserves original files during conversion

## Performance

Independent studies demonstrate measurable efficiency improvements:

- **Significant token reduction** per development task
- **Extended work sessions** before hitting rate limits
- **Selective context loading** eliminates processing irrelevant content

**Note**: Conversion requires upfront token cost but achieves break-even after several development tasks.

### Performance Studies

**[Automated Testing Study](docs/performance-studies/CONTINDEX-PERFORMANCE-STUDY.md)**
- Up to 55% token reduction using tiktoken-validated measurements
- Authentication task focused testing methodology
- Reproducible results with real AI development workflows

**[Manual Validation Study](docs/performance-studies/MANUAL-VALIDATION-STUDY.md)** 
- Up to 45% token reduction through realistic AI chapter selection
- Real AI tool workflow validation
- Conversion cost analysis for honest assessment

**[Workflow Impact Analysis](docs/performance-studies/WORKFLOW-IMPACT-ANALYSIS.md)**
- Side-by-side development task comparison 
- Break-even analysis (approximately 2 tasks to justify conversion cost)
- Honest evaluation of when contindex provides value

## How It Works

1. **Conversion**: Analyzes monolithic file content and creates descriptively-named chapter files
2. **Index Generation**: Creates lightweight index file that references all chapters
3. **Selective Loading**: AI tools can read index and load only relevant chapters
4. **Maintenance**: Update index when chapters are added, removed, or modified

## Directory Customization

Contindex provides flexible directory options to avoid conflicts with existing project structure. **Both directory names are fully customizable** - you can use any valid directory name you prefer.

### Default Structure
```
your-project/
├── context/                    # Chapter files (default)
├── backup/                     # Original file backup (default)
└── [AGENT].md                  # Index file (varies by template)
```

### Custom Directories
```bash
# Avoid conflicts with existing directories
contindex convert --source=my-context.md --context-dir=docs-context --backup-dir=backups

# Results in:
your-project/
├── docs-context/              # Custom context directory (--context-dir)
├── backups/                   # Custom backup directory (--backup-dir)
│   └── my-context.md          # Original source file backed up
└── CLAUDE.md                  # Index file (using claude template)
```

### Options for Different Scenarios

**Existing `context/` directory:**
```bash
# Option 1: Use different name
contindex convert --context-dir=project-context

# Option 2: Keep simple names, different location
contindex convert --context-dir=chapters

# Option 3: Overwrite existing (careful!)
contindex convert --force
```

**No backup needed:**
```bash
# Skip backup creation entirely
contindex convert --no-backup
```

**Simple directory names:**
```bash
# Use simple, short names
contindex convert --context-dir=docs --backup-dir=old

# Or even shorter
contindex convert --context-dir=ctx --backup-dir=bak
```

**Frontend projects with conflicts:**
```bash
# Avoid common frontend directory names
contindex convert --context-dir=ai-context --backup-dir=.backups
```

### Conflict Detection

Contindex automatically detects directory conflicts and provides helpful suggestions:
- If `context/` exists with files, suggests `--context-dir` or `--force`
- Prevents same name for context and backup directories
- Shows clear error messages with solution options

## Commands

### For Fresh Projects
```bash
# Initialize index-chapter structure
contindex init --template=claude    # Creates CLAUDE.md
contindex init --template=cursor    # Creates AGENTS.md  
contindex init --template=copilot   # Creates .github/copilot-instructions.md
contindex init --template=gemini    # Creates GEMINI.md
contindex init --template=generic   # Creates template.md (universal)
```

### For Existing Monolithic Files
```bash
# Convert monolithic file to index-chapter architecture
contindex convert --source=CLAUDE.md --project="My Project"

# Preview conversion without changes
contindex convert --source=CLAUDE.md --dry-run

# Customize directories to avoid conflicts  
contindex convert --source=CLAUDE.md --context-dir=docs-context --backup-dir=backups

# Advanced options
contindex convert --source=CLAUDE.md --no-backup --force --template=cursor
```

### Maintaining Your Index
```bash
# Update index when you add/remove chapter files (specify your template)
contindex update --template=claude    # For CLAUDE.md
contindex update --template=cursor    # For AGENTS.md
contindex update --template=gemini    # For GEMINI.md

# Force update even if no changes detected
contindex update --force
```

## Project Structure

### Default Structure
After initialization or conversion:
```
your-project/
├── context/                    # Chapter files directory (default)
│   ├── authentication-service.md # Example files generated from monolithic context file
│   ├── database-schema.md
│   └── api-endpoints.md
├── [AGENT].md                  # Index file (varies by template - see below)
└── backup/                     # Original files (default backup location)
    └── [source-file].md        # Your original file backed up here

TIP: Use `--no-backup` during convert command to skip backup completely! 
```

**Index filenames by template:**
- `claude` → `CLAUDE.md`
- `cursor` → `AGENTS.md`
- `copilot` → `copilot-instructions.md`  
- `gemini` → `GEMINI.md`
- `generic` → `template.md`

### Custom Structure
Using `--context-dir` and `--backup-dir` flags:
```bash
contindex convert --source=my-docs.md --context-dir=docs-context --backup-dir=backups --template=cursor
```
Results in:
```
your-project/
├── docs-context/              # Custom chapter directory (--context-dir)
│   ├── authentication-service.md
│   ├── database-schema.md
│   └── api-endpoints.md  
├── AGENTS.md                  # Index file (cursor template)
└── backups/                   # Custom backup directory (--backup-dir)
    └── my-docs.md             # Original source file backed up
```

## Template Examples

Each template creates an index optimized for specific AI tools:

**Claude Code Template (CLAUDE.md):**
```markdown
# Project Context Index

This index provides a table of contents for the context files.
AI tools can reference this index and load specific chapters (context files) as needed.

## Available Chapter Examples: (these will be the filenames of your context documents, semantically renamed by AI for optimal consumption by AI)

1. context/authentication-service.md
2. context/database-schema.md
...
100. context/another-semantically-aligned-filename.md

## How the LLM/Agent will consume

AI tools should:
1. Reference this index to understand available chapters and what their contents will be about, all filenames will be semanticly aligned to their contents.
2. Load specific chapter files from the context/ directory  
3. Process only relevant chapters instead of everything
```

**Cursor IDE Template (AGENTS.md):**
- Optimized for Cursor IDE with `@path/file.md` references
- Simple structure focused on development workflow

**GitHub Copilot Template (.github/copilot-instructions.md):**
- Placed in .github/ directory for GitHub Copilot integration
- Includes context organization instructions

**Google Gemini Template (GEMINI.md):**
- Optimized for Gemini's conversational context loading
- Request-based file loading workflow

**Generic Template (template.md):**
- Universal template that can be adapted to any AI tool
- Tool-agnostic approach with flexible instructions

## Use Cases

- **Large codebases** with extensive context requirements or large amounts of context documents
- **Team projects** needing organized, addressable context
- **AI-assisted development** workflows requiring efficient context management
- **Legacy projects** with monolithic context files needing modernization

## Development

### Contributing

Built as a focused solution to context organization. Contributions should maintain the core philosophy of simplicity, security, and efficiency.

### Building from Source

Prerequisites:
- Go 1.21 or later

```bash
# Clone the repository
git clone https://github.com/angelcodes95/contindex.git
cd contindex

# Build the binary
go build -o contindex .

# Run tests
go test ./...

# Install globally (optional)
go install .
```

### Architecture

The project follows clean architecture principles:

```
contindex/
├── cmd/                     # CLI commands
│   ├── convert.go          # Convert command
│   ├── init.go             # Init command
│   ├── root.go             # Root command and CLI setup
│   ├── template.go         # Template command
│   └── update.go           # Update command
├── docs/                   # Documentation
│   └── performance-studies/
├── internal/               # Internal packages
│   ├── classifier/         # Content analysis and categorization  
│   ├── config/             # Configuration management
│   ├── errors/             # Centralized error types
│   ├── logging/            # Structured logging
│   ├── template/           # Template management
│   │   ├── embed.go        # Embedded file system
│   │   ├── template.go     # Template processing
│   │   └── templates/      # Template files
│   │       ├── claude/
│   │       ├── cursor/  
│   │       ├── copilot/
│   │       ├── gemini/
│   │       └── generic/
│   └── validation/         # Input validation and security
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
└── README.md               # This file
```

## THE FUN FACT
This is my first project using Go. I consumed a series of lectures on Go administered by Matthew Holiday and then used my own approach to leverage CLI Agents to clarify pieces I was still confused on with examples and best practice approaches. I have the ability to deep dive if still confused and I make sure to double check against other online sources, not assume the agent is giving factual information at any time. The following is a simple example doing this method to learn about the Template System:

```
 internal/template/ - Template System
  - embed.go: //go:embed templates/* + var TemplateFS embed.FS
    - Why: Embeds template files into binary at compile time
    - Best Practice: Zero external file dependencies for distribution
    - Go Pattern: Modern embed.FS for asset bundling
  - template.go: Template processing engine
    - Why: Handles template parsing, data injection, file generation
    - Best Practice: Structured data types, comprehensive error handling
    - Go Pattern: text/template for safe template rendering
  - templates/: Actual template files
    - Why: Separate templates for different AI tools (Claude, Cursor, etc.)
    - Best Practice: Tool-specific optimization while maintaining consistency

```

## License

MIT License

Copyright (c) 2025 AngelCodes95

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OR GUARANTEE OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.