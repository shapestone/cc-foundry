package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	embedpkg "github.com/shapestone/cc-foundry/pkg/embed"
	"github.com/shapestone/cc-foundry/pkg/state"
)

// InstallMode determines where files are installed
type InstallMode int

const (
	InstallModeUser    InstallMode = iota // ~/.claude/ (user-level, all projects)
	InstallModeProject                     // .claude/ (project-level, version-controlled)
)

// CurrentInstallMode is the active installation mode (default: user-level)
var CurrentInstallMode = InstallModeUser

// GetClaudeCodeDir returns the Claude Code directory path based on install mode
func GetClaudeCodeDir() (string, error) {
	if CurrentInstallMode == InstallModeProject {
		// Project-level: .claude/ in current directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return filepath.Join(cwd, ".claude"), nil
	}

	// User-level: ~/.claude/
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".claude"), nil
}

// GetTypeDir returns the full path to a specific type directory (commands, agents, skills)
func GetTypeDir(fileType string) (string, error) {
	baseDir, err := GetClaudeCodeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDir, fileType), nil
}

// EnsureDirectoriesExist creates Claude Code directories if they don't exist
func EnsureDirectoriesExist() error {
	for _, fileType := range []string{"commands", "agents", "skills"} {
		dir, err := GetTypeDir(fileType)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// GenerateInstalledFilename creates the ccf-prefixed filename
func GenerateInstalledFilename(category, filename string) string {
	// Remove .md extension
	base := strings.TrimSuffix(filename, ".md")

	// Generate: ccf-[category]-[filename].md
	return fmt.Sprintf("ccf-%s-%s.md", category, base)
}

// InstallFile installs a single file
func InstallFile(file embedpkg.CategoryFile, st *state.State) error {
	// Ensure directories exist
	if err := EnsureDirectoriesExist(); err != nil {
		return err
	}

	// Get target directory
	typeDir, err := GetTypeDir(file.Type)
	if err != nil {
		return err
	}

	// Generate installed filename/path based on type
	var installedPath, displayPath, installedFilename string

	if file.Type == "skills" {
		// Skills: subdirectory with SKILL.md
		skillName := GenerateInstalledFilename(file.Category, file.Filename)
		skillName = strings.TrimSuffix(skillName, ".md") // Remove .md extension
		skillDir := filepath.Join(typeDir, skillName)
		installedPath = filepath.Join(skillDir, "SKILL.md")
		installedFilename = filepath.Join(skillName, "SKILL.md")

		// Create skill subdirectory
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return fmt.Errorf("failed to create skill directory %s: %w", skillDir, err)
		}
	} else {
		// Commands and agents: flat .md files
		installedFilename = GenerateInstalledFilename(file.Category, file.Filename)
		installedPath = filepath.Join(typeDir, installedFilename)
	}

	// Format display path (replace home with ~)
	displayPath = installedPath
	if home, err := os.UserHomeDir(); err == nil {
		displayPath = strings.Replace(installedPath, home, "~", 1)
	}

	// Determine type label (singular form)
	typeLabel := strings.TrimSuffix(file.Type, "s") // "agents" -> "agent"

	// Check if already installed
	existing := st.FindInstallation(installedPath)
	isUpdate := false

	if existing != nil {
		// File already installed, check if content changed
		if !existing.HasContentChanged(file.Content) {
			fmt.Printf("  ✓ %s: %s → %s (unchanged)\n", typeLabel, installedFilename, displayPath)
			return nil
		}

		fmt.Printf("  ⚠ %s: %s → %s (updating)\n", typeLabel, installedFilename, displayPath)
		isUpdate = true
	}

	// Write file
	if err := os.WriteFile(installedPath, file.Content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", installedPath, err)
	}

	// Update state
	st.RemoveInstallation(installedPath) // Remove old entry if exists
	st.AddInstallation(file.Category, file.Type, file.Filename, installedPath, file.Content)

	// Show success with type and path
	if isUpdate {
		fmt.Printf("  ✓ %s: %s → %s (updated)\n", typeLabel, installedFilename, displayPath)
	} else {
		fmt.Printf("  ✓ %s: %s → %s\n", typeLabel, installedFilename, displayPath)
	}
	return nil
}

// InstallCategory installs all files in a category
func InstallCategory(category string) error {
	var files []embedpkg.CategoryFile
	var err error

	if category == "" {
		files, err = embedpkg.ListAllFiles()
	} else {
		files, err = embedpkg.ListCategoryFiles(category)
	}
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		if category == "" {
			return fmt.Errorf("no installable files found")
		}
		return fmt.Errorf("no files found in category '%s'", category)
	}

	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Clear screen and show banner
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
	if category == "" {
		fmt.Printf("Installing all categories [%s]\n", GetInstallModeDescription())
	} else {
		fmt.Printf("Installing category: %s [%s]\n", category, GetInstallModeDescription())
	}

	for _, file := range files {
		if err := InstallFile(file, st); err != nil {
			return err
		}
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	if category == "" {
		fmt.Printf("\n✓ Successfully installed %d files from all categories [%s]\n", len(files), GetInstallModeDescription())
	} else {
		fmt.Printf("\n✓ Successfully installed %d files from category '%s' [%s]\n", len(files), category, GetInstallModeDescription())
	}
	return nil
}

// InstallType installs all files of a specific type in a category
func InstallType(category, fileType string) error {
	files, err := embedpkg.ListTypeFiles(category, fileType)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no %s found in category '%s'", fileType, category)
	}

	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Clear screen and show banner
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
	fmt.Printf("Installing %s from category: %s [%s]\n", fileType, category, GetInstallModeDescription())

	for _, file := range files {
		if err := InstallFile(file, st); err != nil {
			return err
		}
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully installed %d %s from category '%s' [%s]\n", len(files), fileType, category, GetInstallModeDescription())
	return nil
}

// GetInstallModeDescription returns a human-readable description of the current install mode
func GetInstallModeDescription() string {
	if CurrentInstallMode == InstallModeProject {
		return "project (.claude/)"
	}
	return "user (~/.claude/)"
}

// ShowBanner displays the application banner with screen clear
func ShowBanner() {
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
}

// pathMatchesInstallMode checks if an installation path matches the current install mode
func pathMatchesInstallMode(installPath string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	userClaudePath := filepath.Join(home, ".claude")
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	projectClaudePath := filepath.Join(cwd, ".claude")

	if CurrentInstallMode == InstallModeUser {
		// User mode: path must be under ~/.claude/
		return strings.HasPrefix(installPath, userClaudePath)
	}

	// Project mode: path must be under current-dir/.claude/
	return strings.HasPrefix(installPath, projectClaudePath)
}

// ListInstallationsForCurrentMode filters installations by current install mode
func ListInstallationsForCurrentMode(st *state.State, category, fileType string) []state.Installation {
	allInstallations := st.ListInstallations(category, fileType)
	var filtered []state.Installation

	for _, inst := range allInstallations {
		if pathMatchesInstallMode(inst.InstalledPath) {
			filtered = append(filtered, inst)
		}
	}

	return filtered
}

// LocationAvailability indicates which locations have files for a category
type LocationAvailability struct {
	HasUserLevel    bool
	HasProjectLevel bool
	UserCount       int
	ProjectCount    int
}

// CheckLocationAvailability checks which locations have files for a category
func CheckLocationAvailability(category, fileType string) (LocationAvailability, error) {
	st, err := state.Load()
	if err != nil {
		return LocationAvailability{}, fmt.Errorf("failed to load state: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return LocationAvailability{}, fmt.Errorf("failed to get home directory: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return LocationAvailability{}, fmt.Errorf("failed to get working directory: %w", err)
	}

	allInstallations := st.ListInstallations(category, fileType)

	var result LocationAvailability
	userClaudePath := filepath.Join(home, ".claude")
	projectClaudePath := filepath.Join(cwd, ".claude")

	for _, inst := range allInstallations {
		if strings.HasPrefix(inst.InstalledPath, userClaudePath) {
			result.HasUserLevel = true
			result.UserCount++
		} else if strings.HasPrefix(inst.InstalledPath, projectClaudePath) {
			result.HasProjectLevel = true
			result.ProjectCount++
		}
	}

	return result, nil
}

// InstallAll installs all files from all categories
func InstallAll() error {
	categories, err := embedpkg.ListCategories()
	if err != nil {
		return fmt.Errorf("failed to list categories: %w", err)
	}

	if len(categories) == 0 {
		return fmt.Errorf("no categories found")
	}

	fmt.Printf("\nInstalling all categories: %s [%s]\n", strings.Join(categories, ", "), GetInstallModeDescription())

	for _, category := range categories {
		if err := InstallCategory(category); err != nil {
			return err
		}
	}

	return nil
}

// RemoveInstallation removes a single installed file
func RemoveInstallation(installation state.Installation) error {
	// For skills, remove the entire subdirectory
	if installation.Type == "skills" {
		// Path is like: ~/.claude/skills/ccf-development-oss-project-setup/SKILL.md
		// We want to remove: ~/.claude/skills/ccf-development-oss-project-setup/
		skillDir := filepath.Dir(installation.InstalledPath)

		if err := os.RemoveAll(skillDir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove skill directory %s: %w", skillDir, err)
		}
	} else {
		// For commands and agents, just remove the file
		if err := os.Remove(installation.InstalledPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove file %s: %w", installation.InstalledPath, err)
		}
	}

	// Format display path (replace home with ~)
	displayPath := installation.InstalledPath
	if home, err := os.UserHomeDir(); err == nil {
		displayPath = strings.Replace(installation.InstalledPath, home, "~", 1)
	}

	// Determine type label (singular form)
	typeLabel := strings.TrimSuffix(installation.Type, "s")

	fmt.Printf("  ✓ %s: %s (removed from %s)\n", typeLabel, filepath.Base(installation.InstalledPath), displayPath)
	return nil
}

// RemoveCategory removes all files from a category
func RemoveCategory(category string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	installations := ListInstallationsForCurrentMode(st, category, "")
	if len(installations) == 0 {
		// No files to remove - skip silently
		return nil
	}

	// Clear screen and show banner
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
	if category == "" {
		fmt.Printf("Removing %d files from all categories [%s]\n", len(installations), GetInstallModeDescription())
	} else {
		fmt.Printf("Removing %d files from category: %s [%s]\n", len(installations), category, GetInstallModeDescription())
	}

	for _, inst := range installations {
		if err := RemoveInstallation(inst); err != nil {
			return err
		}
		st.RemoveInstallation(inst.InstalledPath)
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	if category == "" {
		fmt.Printf("\n✓ Successfully removed %d files from all categories [%s]\n", len(installations), GetInstallModeDescription())
	} else {
		fmt.Printf("\n✓ Successfully removed %d files from category '%s' [%s]\n", len(installations), category, GetInstallModeDescription())
	}
	return nil
}

// RemoveType removes all files of a specific type from a category
func RemoveType(category, fileType string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	installations := ListInstallationsForCurrentMode(st, category, fileType)
	if len(installations) == 0 {
		// No files to remove - skip silently
		return nil
	}

	// Clear screen and show banner
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
	fmt.Printf("Removing %d %s from category: %s [%s]\n", len(installations), fileType, category, GetInstallModeDescription())

	for _, inst := range installations {
		if err := RemoveInstallation(inst); err != nil {
			return err
		}
		st.RemoveInstallation(inst.InstalledPath)
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully removed %d %s from category '%s' [%s]\n", len(installations), fileType, category, GetInstallModeDescription())
	return nil
}

// RemoveAll removes all installed files
func RemoveAll() error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	installations := ListInstallationsForCurrentMode(st, "", "")

	// Clear screen and show banner
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))

	if len(installations) == 0 {
		fmt.Println("No files installed by foundry")
		return nil
	}

	fmt.Printf("Removing all %d installed files [%s]\n", len(installations), GetInstallModeDescription())

	for _, inst := range installations {
		if err := RemoveInstallation(inst); err != nil {
			// Log error but continue
			fmt.Printf("  ⚠ Error removing %s: %v\n", filepath.Base(inst.InstalledPath), err)
		}
		st.RemoveInstallation(inst.InstalledPath)
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully removed all installed files [%s]\n", GetInstallModeDescription())
	return nil
}
