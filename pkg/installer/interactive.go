package installer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	embedpkg "github.com/shapestone/claude-code-foundry/pkg/embed"
	"github.com/shapestone/claude-code-foundry/pkg/state"
)

// menuModel represents the state of the menu
type menuModel struct {
	prompt   string
	options  []string
	selected int
	canceled bool
}

// Init implements tea.Model
func (m menuModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.canceled = true
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.options)-1 {
				m.selected++
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model
func (m menuModel) View() string {
	// ASCII art banner at the top
	header := bannerStyle.Render(banner)

	// Styled prompt/title
	prompt := promptStyle.Render(m.prompt)

	// Build menu items with styling
	var menuItems string
	for i, option := range m.options {
		cursor := "  "
		var line string

		if i == m.selected {
			// Selected item: highlighted with styled cursor
			cursor = cursorStyle.Render("❯")
			line = cursor + " " + selectedItemStyle.Render(option)
		} else {
			// Normal item
			line = cursor + " " + normalItemStyle.Render(option)
		}

		menuItems += line + "\n"
	}

	// Help text at bottom
	helpText := helpStyle.Render("Navigate: ↑/↓  Select: Enter (↵)  Quit: q")

	// Combine all elements with banner at top
	content := header + "\n" + prompt + "\n\n" + menuItems + "\n" + helpText

	return content
}

// SelectOption displays an arrow-key navigable menu and returns the selected index
func SelectOption(prompt string, options []string) (int, error) {
	m := menuModel{
		prompt:   prompt,
		options:  options,
		selected: 0,
		canceled: false,
	}

	// Use alternate screen buffer for clean, full-screen display
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return -1, fmt.Errorf("error running menu: %w", err)
	}

	result := finalModel.(menuModel)
	if result.canceled {
		return -1, fmt.Errorf("cancelled by user")
	}

	return result.selected, nil
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

// WaitForKey waits for the user to press any key to continue
func WaitForKey() {
	fmt.Print("\nPress Enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
