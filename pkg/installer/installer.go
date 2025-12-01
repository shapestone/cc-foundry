package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	embedpkg "github.com/shapestone/claude-code-foundry/pkg/embed"
	"github.com/shapestone/claude-code-foundry/pkg/state"
)

// GetClaudeCodeDir returns the Claude Code directory path based on OS
func GetClaudeCodeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Linux uses ~/.config/claude/
	if runtime.GOOS == "linux" {
		return filepath.Join(home, ".config", "claude"), nil
	}

	// macOS and others use ~/.claudecode/
	return filepath.Join(home, ".claudecode"), nil
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

	// Generate installed filename
	installedFilename := GenerateInstalledFilename(file.Category, file.Filename)
	installedPath := filepath.Join(typeDir, installedFilename)

	// Check if already installed
	if existing := st.FindInstallation(installedPath); existing != nil {
		// File already installed, check if content changed
		if !existing.HasContentChanged(file.Content) {
			fmt.Printf("  ✓ %s (already installed, unchanged)\n", installedFilename)
			return nil
		}

		fmt.Printf("  ⚠ %s (already installed, will update)\n", installedFilename)
	}

	// Write file
	if err := os.WriteFile(installedPath, file.Content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", installedPath, err)
	}

	// Update state
	st.RemoveInstallation(installedPath) // Remove old entry if exists
	st.AddInstallation(file.Category, file.Type, file.Filename, installedPath, file.Content)

	fmt.Printf("  ✓ %s\n", installedFilename)
	return nil
}

// InstallCategory installs all files in a category
func InstallCategory(category string) error {
	files, err := embedpkg.ListCategoryFiles(category)
	if err != nil {
		return fmt.Errorf("failed to list category files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found in category '%s'", category)
	}

	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	fmt.Printf("\nInstalling category: %s\n", category)

	for _, file := range files {
		if err := InstallFile(file, st); err != nil {
			return err
		}
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully installed %d files from category '%s'\n", len(files), category)
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

	fmt.Printf("\nInstalling %s from category: %s\n", fileType, category)

	for _, file := range files {
		if err := InstallFile(file, st); err != nil {
			return err
		}
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully installed %d %s from category '%s'\n", len(files), fileType, category)
	return nil
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

	fmt.Printf("\nInstalling all categories: %s\n", strings.Join(categories, ", "))

	for _, category := range categories {
		if err := InstallCategory(category); err != nil {
			return err
		}
	}

	return nil
}

// RemoveInstallation removes a single installed file
func RemoveInstallation(installation state.Installation) error {
	if err := os.Remove(installation.InstalledPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove file %s: %w", installation.InstalledPath, err)
	}

	fmt.Printf("  ✓ Removed %s\n", filepath.Base(installation.InstalledPath))
	return nil
}

// RemoveCategory removes all files from a category
func RemoveCategory(category string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	installations := st.ListInstallations(category, "")
	if len(installations) == 0 {
		return fmt.Errorf("no files installed from category '%s'", category)
	}

	fmt.Printf("\nRemoving %d files from category: %s\n", len(installations), category)

	for _, inst := range installations {
		if err := RemoveInstallation(inst); err != nil {
			return err
		}
		st.RemoveInstallation(inst.InstalledPath)
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully removed %d files from category '%s'\n", len(installations), category)
	return nil
}

// RemoveType removes all files of a specific type from a category
func RemoveType(category, fileType string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	installations := st.ListInstallations(category, fileType)
	if len(installations) == 0 {
		return fmt.Errorf("no %s installed from category '%s'", fileType, category)
	}

	fmt.Printf("\nRemoving %d %s from category: %s\n", len(installations), fileType, category)

	for _, inst := range installations {
		if err := RemoveInstallation(inst); err != nil {
			return err
		}
		st.RemoveInstallation(inst.InstalledPath)
	}

	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("\n✓ Successfully removed %d %s from category '%s'\n", len(installations), fileType, category)
	return nil
}

// RemoveAll removes all installed files
func RemoveAll() error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	installations := st.ListInstallations("", "")
	if len(installations) == 0 {
		fmt.Println("\nNo files installed by foundry")
		return nil
	}

	fmt.Printf("\nRemoving all %d installed files\n", len(installations))

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

	fmt.Printf("\n✓ Successfully removed all installed files\n")
	return nil
}
