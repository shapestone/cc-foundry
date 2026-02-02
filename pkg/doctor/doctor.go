package doctor

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shapestone/cc-foundry/pkg/state"
)

// Issue represents a detected problem
type Issue struct {
	Type        string // "error", "warning", "info"
	Category    string
	Description string
	CanFix      bool
	FixFunc     func() error
}

// HealthReport contains the results of the health check
type HealthReport struct {
	Issues           []Issue
	Errors           int
	Warnings         int
	FilesChecked     int
	CorruptedFiles   int
	MissingFiles     int
	ModifiedFiles    int
	OrphanedFiles    int
}

// Run performs a comprehensive health check and returns a report
func Run() (*HealthReport, error) {
	report := &HealthReport{}

	fmt.Println("🏥 Running doctor diagnostics...")
	fmt.Println()

	// 1. Verify ~/.claude.json
	if err := checkClaudeConfig(report); err != nil {
		fmt.Println("✗ Checking Claude Code configuration (~/.claude.json)")
		return report, err
	}
	fmt.Println("✓ Checking Claude Code configuration (~/.claude.json)")

	// 2. Check file integrity
	if err := checkFileIntegrity(report); err != nil {
		fmt.Println("✗ Checking foundry-managed files")
		return report, err
	}
	fmt.Printf("✓ Checking foundry-managed files (%d files)\n", report.FilesChecked)

	// 3. Detect conflicts
	if err := detectConflicts(report); err != nil {
		fmt.Println("✗ Detecting orphaned and conflicting files")
		return report, err
	}
	fmt.Println("✓ Detecting orphaned and conflicting files")

	fmt.Println()
	return report, nil
}

// checkClaudeConfig verifies ~/.claude.json exists and is valid
func checkClaudeConfig(report *HealthReport) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeConfigPath := filepath.Join(home, ".claude.json")

	// Check if file exists
	if _, err := os.Stat(claudeConfigPath); os.IsNotExist(err) {
		report.Warnings++
		report.Issues = append(report.Issues, Issue{
			Type:        "warning",
			Category:    "config",
			Description: "~/.claude.json does not exist (Claude Code may not be installed)",
			CanFix:      false,
		})
		return nil
	}

	// Check if file is valid JSON
	data, err := os.ReadFile(claudeConfigPath)
	if err != nil {
		report.Errors++
		report.Issues = append(report.Issues, Issue{
			Type:        "error",
			Category:    "config",
			Description: fmt.Sprintf("Cannot read ~/.claude.json: %v", err),
			CanFix:      false,
		})
		return nil
	}

	var configData interface{}
	if err := json.Unmarshal(data, &configData); err != nil {
		report.Errors++
		report.Issues = append(report.Issues, Issue{
			Type:        "error",
			Category:    "config",
			Description: fmt.Sprintf("~/.claude.json is not valid JSON: %v", err),
			CanFix:      false,
		})
		return nil
	}

	// Check file size (warn if > 50MB as per performance issue #5024)
	fileInfo, _ := os.Stat(claudeConfigPath)
	sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	if sizeMB > 50 {
		report.Warnings++
		report.Issues = append(report.Issues, Issue{
			Type:        "warning",
			Category:    "config",
			Description: fmt.Sprintf("~/.claude.json is large (%.1fMB) - may cause performance issues", sizeMB),
			CanFix:      false,
		})
	}

	return nil
}

// checkFileIntegrity verifies installed files match expected hashes
func checkFileIntegrity(report *HealthReport) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if len(st.Installations) == 0 {
		return nil
	}

	for _, inst := range st.Installations {
		report.FilesChecked++

		// Check if file exists
		if _, err := os.Stat(inst.InstalledPath); os.IsNotExist(err) {
			report.MissingFiles++
			report.Errors++
			report.Issues = append(report.Issues, Issue{
				Type:        "error",
				Category:    inst.Category,
				Description: fmt.Sprintf("Missing file: %s", inst.InstalledPath),
				CanFix:      true,
				FixFunc:     createFixMissingFileFunc(inst),
			})
			continue
		}

		// Read file and check hash
		content, err := os.ReadFile(inst.InstalledPath)
		if err != nil {
			report.Errors++
			report.Issues = append(report.Issues, Issue{
				Type:        "error",
				Category:    inst.Category,
				Description: fmt.Sprintf("Cannot read file %s: %v", inst.InstalledPath, err),
				CanFix:      false,
			})
			continue
		}

		currentHash := fmt.Sprintf("%x", sha256.Sum256(content))
		if currentHash != inst.Hash {
			report.ModifiedFiles++
			report.Warnings++
			report.Issues = append(report.Issues, Issue{
				Type:        "warning",
				Category:    inst.Category,
				Description: fmt.Sprintf("Modified file detected: %s (hash mismatch)", inst.InstalledPath),
				CanFix:      false,
			})
		}
	}

	return nil
}

// detectConflicts finds duplicate files or naming issues
func detectConflicts(report *HealthReport) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Load state to know which files are managed by foundry
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	managedPaths := make(map[string]bool)
	for _, inst := range st.Installations {
		managedPaths[inst.InstalledPath] = true
	}

	// Check user-level directories
	claudeDir := filepath.Join(home, ".claude")
	if err := detectConflictsInDir(claudeDir, managedPaths, report); err != nil {
		return err
	}

	// Check project-level directories
	cwd, err := os.Getwd()
	if err == nil {
		projectClaudeDir := filepath.Join(cwd, ".claude")
		if _, err := os.Stat(projectClaudeDir); err == nil {
			if err := detectConflictsInDir(projectClaudeDir, managedPaths, report); err != nil {
				return err
			}
		}
	}

	return nil
}

// detectConflictsInDir checks a directory for conflicts
func detectConflictsInDir(baseDir string, managedPaths map[string]bool, report *HealthReport) error {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil
	}

	for _, subdir := range []string{"commands", "agents", "skills"} {
		subdirPath := filepath.Join(baseDir, subdir)
		if _, err := os.Stat(subdirPath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(subdirPath)
		if err != nil {
			continue
		}

		// Check for orphaned ccf- files (not in state)
		for _, entry := range entries {
			fullPath := filepath.Join(subdirPath, entry.Name())

			// Skip directories for commands/agents
			if subdir != "skills" && entry.IsDir() {
				continue
			}

			// Check if this is a ccf- file but not managed
			if strings.HasPrefix(entry.Name(), "ccf-") && !managedPaths[fullPath] {
				// For skills, check if it's a directory
				if subdir == "skills" && entry.IsDir() {
					skillFile := filepath.Join(fullPath, "SKILL.md")
					if !managedPaths[skillFile] {
						report.OrphanedFiles++
						report.Warnings++
						report.Issues = append(report.Issues, Issue{
							Type:        "warning",
							Category:    "orphaned",
							Description: fmt.Sprintf("Orphaned foundry file: %s (not tracked in state)", fullPath),
							CanFix:      true,
							FixFunc:     createRemoveOrphanedFunc(fullPath, true),
						})
					}
				} else if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					report.OrphanedFiles++
					report.Warnings++
					report.Issues = append(report.Issues, Issue{
						Type:        "warning",
						Category:    "orphaned",
						Description: fmt.Sprintf("Orphaned foundry file: %s (not tracked in state)", fullPath),
						CanFix:      true,
						FixFunc:     createRemoveOrphanedFunc(fullPath, false),
					})
				}
			}
		}
	}

	return nil
}

// createFixMissingFileFunc creates a fix function for missing files
func createFixMissingFileFunc(inst state.Installation) func() error {
	return func() error {
		// For now, just remove from state
		// Future: could reinstall from embedded files
		st, err := state.Load()
		if err != nil {
			return err
		}
		st.RemoveInstallation(inst.InstalledPath)
		return st.Save()
	}
}

// createRemoveOrphanedFunc creates a fix function for orphaned files
func createRemoveOrphanedFunc(path string, isDir bool) func() error {
	return func() error {
		if isDir {
			return os.RemoveAll(path)
		}
		return os.Remove(path)
	}
}

// PrintReport displays the health report
func PrintReport(report *HealthReport) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 Health Report")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	if len(report.Issues) == 0 {
		fmt.Println("✓ No issues found - everything looks healthy!")
		fmt.Printf("\nFiles checked: %d\n", report.FilesChecked)
		return
	}

	// Print summary
	fmt.Printf("Files checked: %d\n", report.FilesChecked)
	if report.Errors > 0 {
		fmt.Printf("❌ Errors: %d\n", report.Errors)
	}
	if report.Warnings > 0 {
		fmt.Printf("⚠️  Warnings: %d\n", report.Warnings)
	}
	if report.MissingFiles > 0 {
		fmt.Printf("Missing files: %d\n", report.MissingFiles)
	}
	if report.ModifiedFiles > 0 {
		fmt.Printf("Modified files: %d\n", report.ModifiedFiles)
	}
	if report.OrphanedFiles > 0 {
		fmt.Printf("Orphaned files: %d\n", report.OrphanedFiles)
	}

	// Print issues by type
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Issues Found:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	for _, issue := range report.Issues {
		icon := "ℹ️ "
		if issue.Type == "error" {
			icon = "❌"
		} else if issue.Type == "warning" {
			icon = "⚠️ "
		}

		fmt.Printf("%s [%s] %s\n", icon, issue.Category, issue.Description)
		if issue.CanFix {
			fmt.Println("   (can be fixed)")
		}
		fmt.Println()
	}
}

// OfferFixes prompts the user to fix issues that can be fixed
func OfferFixes(report *HealthReport, selectOptionFunc func(string, []string) (int, error)) error {
	fixableIssues := []Issue{}
	for _, issue := range report.Issues {
		if issue.CanFix && issue.FixFunc != nil {
			fixableIssues = append(fixableIssues, issue)
		}
	}

	if len(fixableIssues) == 0 {
		return nil
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n%d issue(s) can be automatically fixed.\n\n", len(fixableIssues))

	options := []string{
		"Yes, fix all issues",
		"No, leave as is",
	}

	selected, err := selectOptionFunc("Would you like to fix these issues?", options)
	if err != nil || selected != 0 {
		return nil
	}

	fmt.Println("\nFixing issues...")
	fixed := 0
	failed := 0

	for _, issue := range fixableIssues {
		if err := issue.FixFunc(); err != nil {
			fmt.Printf("❌ Failed to fix: %s (%v)\n", issue.Description, err)
			failed++
		} else {
			fmt.Printf("✓ Fixed: %s\n", issue.Description)
			fixed++
		}
	}

	fmt.Printf("\nFixed: %d, Failed: %d\n", fixed, failed)
	return nil
}
