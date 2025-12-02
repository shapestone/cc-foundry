package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shapestone/claude-code-foundry/pkg/state"
)

// ShowDirectoryStructure displays the Claude Code directory structure and installed files
func ShowDirectoryStructure() error {
	fmt.Println("\n📁 Claude Code Directory Structure\n")

	// Show user-level directory
	if err := showLocation("User-level", "~/.claude/", true); err != nil {
		return err
	}

	fmt.Println()

	// Show project-level directory
	if err := showLocation("Project-level", ".claude/", false); err != nil {
		return err
	}

	// Show installed files by category
	fmt.Println("\n📦 Installed Files (managed by foundry)\n")
	if err := showInstalledFiles(); err != nil {
		return err
	}

	fmt.Println()
	return nil
}

// showLocation displays directory structure for a specific location
func showLocation(label, displayPath string, isUser bool) error {
	var basePath string

	if isUser {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(home, ".claude")
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		basePath = filepath.Join(cwd, ".claude")
	}

	fmt.Printf("%s (%s):\n", label, displayPath)

	// Check if directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		fmt.Printf("  ✗ Directory does not exist\n")
		return nil
	}

	// Show subdirectories with file counts
	for _, subdir := range []string{"commands", "agents", "skills"} {
		subdirPath := filepath.Join(basePath, subdir)
		count, err := countFiles(subdirPath, subdir == "skills")
		if err != nil {
			fmt.Printf("  %s/  (error reading: %v)\n", subdir, err)
			continue
		}

		if count == 0 {
			fmt.Printf("  %s/  (empty)\n", subdir)
		} else {
			itemType := "file"
			if subdir == "skills" {
				itemType = "skill"
			}
			if count != 1 {
				itemType += "s"
			}
			fmt.Printf("  %s/  (%d %s)\n", subdir, count, itemType)
		}
	}

	return nil
}

// countFiles counts files in a directory
// For skills directory, counts subdirectories instead of files
func countFiles(dirPath string, isSkillsDir bool) (int, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return 0, nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if isSkillsDir {
			// For skills, count directories
			if entry.IsDir() {
				count++
			}
		} else {
			// For commands/agents, count .md files
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				count++
			}
		}
	}

	return count, nil
}

// showInstalledFiles displays installed files grouped by category
func showInstalledFiles() error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if len(st.Installations) == 0 {
		fmt.Println("  No files installed by foundry yet")
		return nil
	}

	// Group installations by category
	byCategory := make(map[string][]state.Installation)
	for _, inst := range st.Installations {
		byCategory[inst.Category] = append(byCategory[inst.Category], inst)
	}

	// Display by category
	for category, installations := range byCategory {
		// Count by type
		counts := make(map[string]int)
		for _, inst := range installations {
			counts[inst.Type]++
		}

		// Build count display
		var countParts []string
		if counts["commands"] > 0 {
			countParts = append(countParts, fmt.Sprintf("%d command%s", counts["commands"], plural(counts["commands"])))
		}
		if counts["agents"] > 0 {
			countParts = append(countParts, fmt.Sprintf("%d agent%s", counts["agents"], plural(counts["agents"])))
		}
		if counts["skills"] > 0 {
			countParts = append(countParts, fmt.Sprintf("%d skill%s", counts["skills"], plural(counts["skills"])))
		}

		fmt.Printf("  %s: %s\n", category, strings.Join(countParts, ", "))
	}

	fmt.Printf("\n  Total: %d file%s installed\n", len(st.Installations), plural(len(st.Installations)))

	return nil
}

// plural returns "s" if count is not 1, otherwise ""
func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
