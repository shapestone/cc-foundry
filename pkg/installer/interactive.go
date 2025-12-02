package installer

import (
	"fmt"
	"os"
	"strings"

	embedpkg "github.com/shapestone/claude-code-foundry/pkg/embed"
	"github.com/shapestone/claude-code-foundry/pkg/state"
	"golang.org/x/term"
)

// SelectOption displays an arrow-key navigable menu and returns the selected index
func SelectOption(prompt string, options []string) (int, error) {
	// Save original terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return -1, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Hide cursor during selection
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor when done

	selected := 0
	buf := make([]byte, 3)
	firstRender := true

	// Initial render
	renderMenu(prompt, options, selected, firstRender)
	firstRender = false

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return -1, err
		}

		// Check for arrow keys and Enter
		if n == 3 && buf[0] == 27 && buf[1] == 91 {
			// Arrow key sequence
			switch buf[2] {
			case 65: // Up arrow
				if selected > 0 {
					selected--
					renderMenu(prompt, options, selected, firstRender)
				}
			case 66: // Down arrow
				if selected < len(options)-1 {
					selected++
					renderMenu(prompt, options, selected, firstRender)
				}
			}
		} else if n == 1 && (buf[0] == 13 || buf[0] == 10) {
			// Enter key
			fmt.Println() // Move to next line after selection
			return selected, nil
		} else if n == 1 && buf[0] == 3 {
			// Ctrl+C
			fmt.Println()
			return -1, fmt.Errorf("cancelled by user")
		}
	}
}

// renderMenu renders the menu with the cursor at the selected option
func renderMenu(prompt string, options []string, selected int, firstRender bool) {
	if !firstRender {
		// Move cursor up to beginning of menu (prompt + number of options)
		linesToMove := len(options) + 1 // +1 for prompt line
		fmt.Printf("\033[%dA", linesToMove)
		// Move cursor to beginning of line
		fmt.Print("\r")
	}

	// Clear from cursor to end of screen
	fmt.Print("\033[J")

	// Print prompt (in raw mode, need \r\n for proper newline)
	fmt.Print(prompt + "\r\n")

	// Print options with cursor (in raw mode, need \r\n for proper newline)
	for i, option := range options {
		if i == selected {
			fmt.Print("❯ " + option + "\r\n")
		} else {
			fmt.Print("  " + option + "\r\n")
		}
	}
}

// PromptForLocation asks the user to choose between project and personal installation
// Returns true to proceed, false to cancel
func PromptForLocation() bool {
	fmt.Println()

	options := []string{
		"1. Project (.claude/)",
		"2. Personal (~/.claude/)",
	}

	selected, err := SelectOption("Choose location", options)
	if err != nil {
		if err.Error() == "cancelled by user" {
			fmt.Println("Installation cancelled.")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return false
	}

	switch selected {
	case 0:
		CurrentInstallMode = InstallModeProject
		fmt.Println("→ Installing to project location (.claude/)")
		return true
	case 1:
		CurrentInstallMode = InstallModeUser
		fmt.Println("→ Installing to personal location (~/.claude/)")
		return true
	default:
		fmt.Println("Invalid selection. Installation cancelled.")
		return false
	}
}

// PreviewChange represents a single file change to preview
type PreviewChange struct {
	Action      string // "install", "update", "skip"
	Type        string // "command", "agent", "skill"
	Name        string // Display name
	Path        string // Installation path
	IsUnchanged bool   // True if file content hasn't changed
}

// PreviewInstall shows what will be installed and asks for confirmation
func PreviewInstall(category string, fileType string) (bool, error) {
	var files []embedpkg.CategoryFile
	var err error

	if fileType != "" {
		files, err = embedpkg.ListTypeFiles(category, fileType)
		if err != nil {
			return false, fmt.Errorf("failed to list files: %w", err)
		}
	} else {
		files, err = embedpkg.ListCategoryFiles(category)
		if err != nil {
			return false, fmt.Errorf("failed to list files: %w", err)
		}
	}

	if len(files) == 0 {
		if fileType != "" {
			return false, fmt.Errorf("no %s found in category '%s'", fileType, category)
		}
		return false, fmt.Errorf("no files found in category '%s'", category)
	}

	// Load state to check for existing installations
	st, err := state.Load()
	if err != nil {
		return false, fmt.Errorf("failed to load state: %w", err)
	}

	// Build preview
	var changes []PreviewChange
	for _, file := range files {
		var installPath string
		var displayName string

		if file.Type == "skills" {
			skillName := GenerateInstalledFilename(file.Category, file.Filename)
			skillName = strings.TrimSuffix(skillName, ".md")
			typeDir, _ := GetTypeDir(file.Type)
			installPath = fmt.Sprintf("%s/%s/SKILL.md", typeDir, skillName)
			displayName = fmt.Sprintf("%s/SKILL.md", skillName)
		} else {
			fileName := GenerateInstalledFilename(file.Category, file.Filename)
			typeDir, _ := GetTypeDir(file.Type)
			installPath = fmt.Sprintf("%s/%s", typeDir, fileName)
			displayName = fileName
		}

		// Replace home with ~ for display
		displayPath := installPath
		if home, err := os.UserHomeDir(); err == nil {
			displayPath = strings.Replace(installPath, home, "~", 1)
		}

		// Check if already installed
		existing := st.FindInstallation(installPath)
		action := "install"
		isUnchanged := false

		if existing != nil {
			if existing.HasContentChanged(file.Content) {
				action = "update"
			} else {
				action = "skip"
				isUnchanged = true
			}
		}

		typeLabel := strings.TrimSuffix(file.Type, "s")
		changes = append(changes, PreviewChange{
			Action:      action,
			Type:        typeLabel,
			Name:        displayName,
			Path:        displayPath,
			IsUnchanged: isUnchanged,
		})
	}

	// Display preview
	fmt.Printf("\nPreview: %s [%s]\n", category, GetInstallModeDescription())
	fmt.Println()

	installCount := 0
	updateCount := 0
	skipCount := 0

	for _, change := range changes {
		switch change.Action {
		case "install":
			fmt.Printf("  + %s: %s → %s\n", change.Type, change.Name, change.Path)
			installCount++
		case "update":
			fmt.Printf("  ↻ %s: %s → %s (will update)\n", change.Type, change.Name, change.Path)
			updateCount++
		case "skip":
			fmt.Printf("  · %s: %s → %s (unchanged)\n", change.Type, change.Name, change.Path)
			skipCount++
		}
	}

	fmt.Println()
	fmt.Printf("Summary: %d to install, %d to update, %d unchanged\n", installCount, updateCount, skipCount)
	fmt.Println()

	// Ask for confirmation
	options := []string{
		"Yes, proceed",
		"No, cancel",
	}

	selected, err := SelectOption("Proceed with installation?", options)
	if err != nil {
		return false, err
	}

	return selected == 0, nil
}

// PreviewRemove shows what will be removed and asks for confirmation
func PreviewRemove(category string, fileType string) (bool, error) {
	st, err := state.Load()
	if err != nil {
		return false, fmt.Errorf("failed to load state: %w", err)
	}

	installations := st.ListInstallations(category, fileType)
	if len(installations) == 0 {
		if fileType != "" {
			fmt.Printf("\nNo %s installed from category '%s'\n", fileType, category)
		} else {
			fmt.Printf("\nNo files installed from category '%s'\n", category)
		}
		return false, nil
	}

	// Display preview
	if fileType != "" {
		fmt.Printf("\nPreview: Remove %s from %s [%s]\n", fileType, category, GetInstallModeDescription())
	} else {
		fmt.Printf("\nPreview: Remove category %s [%s]\n", category, GetInstallModeDescription())
	}
	fmt.Println()

	for _, inst := range installations {
		displayPath := inst.InstalledPath
		if home, err := os.UserHomeDir(); err == nil {
			displayPath = strings.Replace(inst.InstalledPath, home, "~", 1)
		}
		typeLabel := strings.TrimSuffix(inst.Type, "s")
		fmt.Printf("  - %s: %s\n", typeLabel, displayPath)
	}

	fmt.Println()
	fmt.Printf("Summary: %d files will be removed\n", len(installations))
	fmt.Println()

	// Ask for confirmation
	options := []string{
		"Yes, remove",
		"No, cancel",
	}

	selected, err := SelectOption("Proceed with removal?", options)
	if err != nil {
		return false, err
	}

	return selected == 0, nil
}
