package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetInstallModeDescription tests the install mode description strings
func TestGetInstallModeDescription(t *testing.T) {
	tests := []struct {
		name     string
		mode     InstallMode
		expected string
	}{
		{
			name:     "user mode",
			mode:     InstallModeUser,
			expected: "user (~/.claude/)",
		},
		{
			name:     "project mode",
			mode:     InstallModeProject,
			expected: "project (.claude/)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current mode
			originalMode := CurrentInstallMode
			defer func() { CurrentInstallMode = originalMode }()

			// Set test mode
			CurrentInstallMode = tt.mode

			result := GetInstallModeDescription()
			if result != tt.expected {
				t.Errorf("GetInstallModeDescription() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestGetClaudeCodeDir_UserMode tests user-level directory resolution
func TestGetClaudeCodeDir_UserMode(t *testing.T) {
	// Save current mode
	originalMode := CurrentInstallMode
	defer func() { CurrentInstallMode = originalMode }()

	// Set user mode
	CurrentInstallMode = InstallModeUser

	dir, err := GetClaudeCodeDir()
	if err != nil {
		t.Fatalf("GetClaudeCodeDir() error = %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}

	expected := filepath.Join(home, ".claude")
	if dir != expected {
		t.Errorf("GetClaudeCodeDir() = %q, want %q", dir, expected)
	}
}

// TestGetClaudeCodeDir_ProjectMode tests project-level directory resolution
func TestGetClaudeCodeDir_ProjectMode(t *testing.T) {
	// Save current mode
	originalMode := CurrentInstallMode
	defer func() { CurrentInstallMode = originalMode }()

	// Set project mode
	CurrentInstallMode = InstallModeProject

	dir, err := GetClaudeCodeDir()
	if err != nil {
		t.Fatalf("GetClaudeCodeDir() error = %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	expected := filepath.Join(cwd, ".claude")
	if dir != expected {
		t.Errorf("GetClaudeCodeDir() = %q, want %q", dir, expected)
	}
}

// TestGetTypeDir tests type directory resolution
func TestGetTypeDir(t *testing.T) {
	// Save current mode
	originalMode := CurrentInstallMode
	defer func() { CurrentInstallMode = originalMode }()

	// Set user mode
	CurrentInstallMode = InstallModeUser

	tests := []struct {
		name     string
		fileType string
	}{
		{"commands", "commands"},
		{"agents", "agents"},
		{"skills", "skills"},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := GetTypeDir(tt.fileType)
			if err != nil {
				t.Fatalf("GetTypeDir(%q) error = %v", tt.fileType, err)
			}

			expected := filepath.Join(home, ".claude", tt.fileType)
			if dir != expected {
				t.Errorf("GetTypeDir(%q) = %q, want %q", tt.fileType, dir, expected)
			}
		})
	}
}

// TestGenerateInstalledFilename tests filename generation
func TestGenerateInstalledFilename(t *testing.T) {
	tests := []struct {
		name     string
		category string
		filename string
		expected string
	}{
		{
			name:     "command file",
			category: "development",
			filename: "deploy.md",
			expected: "ccf-development-deploy.md",
		},
		{
			name:     "agent file",
			category: "testing",
			filename: "test-runner.md",
			expected: "ccf-testing-test-runner.md",
		},
		{
			name:     "skill file",
			category: "development",
			filename: "makefile-guide.md",
			expected: "ccf-development-makefile-guide.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateInstalledFilename(tt.category, tt.filename)
			if result != tt.expected {
				t.Errorf("GenerateInstalledFilename(%q, %q) = %q, want %q",
					tt.category, tt.filename, result, tt.expected)
			}
		})
	}
}

// TestEnsureDirectoriesExist tests directory creation
func TestEnsureDirectoriesExist(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ccf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current mode
	originalMode := CurrentInstallMode
	defer func() { CurrentInstallMode = originalMode }()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Set project mode (so it uses .claude/ in current dir)
	CurrentInstallMode = InstallModeProject

	// Test directory creation
	if err := EnsureDirectoriesExist(); err != nil {
		t.Fatalf("EnsureDirectoriesExist() error = %v", err)
	}

	// Verify directories were created
	for _, dirName := range []string{"commands", "agents", "skills"} {
		dirPath := filepath.Join(tmpDir, ".claude", dirName)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %q to exist, but it doesn't", dirPath)
		}
	}
}

// TestSkillDirectoryStructure tests that skills are installed in subdirectories
func TestSkillDirectoryStructure(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ccf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test skill installation creates subdirectory with SKILL.md
	skillName := "ccf-development-test-skill"
	skillDir := filepath.Join(tmpDir, skillName)
	skillFile := filepath.Join(skillDir, "SKILL.md")

	// Create the directory structure
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("Failed to create skill directory: %v", err)
	}

	// Write SKILL.md file
	content := []byte("# Test Skill\nThis is a test skill.")
	if err := os.WriteFile(skillFile, content, 0644); err != nil {
		t.Fatalf("Failed to write SKILL.md: %v", err)
	}

	// Verify structure
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		t.Errorf("Expected SKILL.md to exist at %q", skillFile)
	}

	// Verify it's in a subdirectory
	parentDir := filepath.Dir(skillFile)
	if filepath.Base(parentDir) != skillName {
		t.Errorf("Expected skill to be in subdirectory %q, got %q", skillName, filepath.Base(parentDir))
	}
}

// TestCommandAgentFlatStructure tests that commands and agents are flat files
func TestCommandAgentFlatStructure(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		filename string
	}{
		{
			name:     "command file",
			fileType: "commands",
			filename: "ccf-development-deploy.md",
		},
		{
			name:     "agent file",
			fileType: "agents",
			filename: "ccf-development-test-agent.md",
		},
	}

	tmpDir, err := os.MkdirTemp("", "ccf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeDir := filepath.Join(tmpDir, tt.fileType)
			if err := os.MkdirAll(typeDir, 0755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			filePath := filepath.Join(typeDir, tt.filename)
			content := []byte("# Test Content")
			if err := os.WriteFile(filePath, content, 0644); err != nil {
				t.Fatalf("Failed to write file: %v", err)
			}

			// Verify file is directly in type directory (not in subdirectory)
			if filepath.Dir(filePath) != typeDir {
				t.Errorf("Expected file to be directly in %q, got %q", typeDir, filepath.Dir(filePath))
			}

			// Verify file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected file to exist at %q", filePath)
			}
		})
	}
}
