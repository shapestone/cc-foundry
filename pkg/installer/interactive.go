package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	embedpkg "github.com/shapestone/cc-foundry/pkg/embed"
	"github.com/shapestone/cc-foundry/pkg/state"
)

// menuModel represents the state of the menu
type menuModel struct {
	prompt     string
	options    []string
	disabled   []bool // whether each option is disabled
	selected   int
	canceled   bool
	showBanner bool // whether to show the banner at the top
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
		case "ctrl+c", "esc":
			m.canceled = true
			return m, tea.Quit
		case "up", "k":
			// Move up, skipping disabled items
			newSelected := m.selected - 1
			for newSelected >= 0 {
				if len(m.disabled) == 0 || !m.disabled[newSelected] {
					m.selected = newSelected
					break
				}
				newSelected--
			}
		case "down", "j":
			// Move down, skipping disabled items
			newSelected := m.selected + 1
			for newSelected < len(m.options) {
				if len(m.disabled) == 0 || !m.disabled[newSelected] {
					m.selected = newSelected
					break
				}
				newSelected++
			}
		case "enter":
			// Don't allow selecting disabled items
			if len(m.disabled) == 0 || !m.disabled[m.selected] {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m menuModel) View() string {
	var content string

	// ASCII art banner at the top (only for full-screen menus)
	if m.showBanner {
		header := bannerStyle.Render(banner)
		content = header + "\n"
	}

	// Styled prompt/title
	prompt := promptStyle.Render(m.prompt)
	content += prompt

	// Build menu items with styling
	var menuItems string
	for i, option := range m.options {
		var line string
		isDisabled := len(m.disabled) > 0 && m.disabled[i]

		if isDisabled {
			// Disabled item: grayed out, no cursor
			line = "  " + disabledItemStyle.Render(option)
		} else if i == m.selected {
			// Selected item: highlighted with styled cursor
			cursor := cursorStyle.Render("❯ ")
			line = cursor + selectedItemStyle.Render(option)
		} else {
			// Normal item
			line = "  " + normalItemStyle.Render(option)
		}

		menuItems += line + "\n"
	}

	// Help text at bottom
	helpText := helpStyle.Render("Navigate: ↑/↓  Select: Enter (↵)  Back: Esc")

	// Combine all elements
	content += "\n\n" + menuItems + "\n" + helpText

	return content
}

// SelectOption displays an arrow-key navigable menu and returns the selected index
func SelectOption(prompt string, options []string) (int, error) {
	return SelectOptionWithDisabled(prompt, options, nil)
}

// SelectOptionAt displays a menu with a pre-selected cursor position
func SelectOptionAt(prompt string, options []string, initialSelected int) (int, error) {
	if initialSelected < 0 || initialSelected >= len(options) {
		initialSelected = 0
	}

	m := menuModel{
		prompt:     prompt,
		options:    options,
		selected:   initialSelected,
		canceled:   false,
		showBanner: true,
	}

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

// SelectOptionWithDisabled displays a menu with some options disabled
func SelectOptionWithDisabled(prompt string, options []string, disabled []bool) (int, error) {
	// Find first non-disabled item to select initially
	initialSelected := 0
	if disabled != nil {
		for i := range options {
			if !disabled[i] {
				initialSelected = i
				break
			}
		}
	}

	m := menuModel{
		prompt:     prompt,
		options:    options,
		disabled:   disabled,
		selected:   initialSelected,
		canceled:   false,
		showBanner: true, // show banner for full-screen menus
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

// SelectOptionInline displays a menu inline (without alt screen) to preserve context
// Use this when you want to show a menu after printing context that should remain visible
func SelectOptionInline(prompt string, options []string) (int, error) {
	m := menuModel{
		prompt:     prompt,
		options:    options,
		selected:   0,
		canceled:   false,
		showBanner: false, // no banner for inline menus
	}

	// Run inline without alternate screen to preserve printed context
	p := tea.NewProgram(m)
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

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get current directory: %v\n", err)
		return false
	}
	claudeDir := filepath.Join(cwd, ".claude")
	projectExists := true
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		projectExists = false
	}

	projectLabel := fmt.Sprintf("Project (%s/.claude/)", cwd)
	if !projectExists {
		projectLabel += " - No Claude Code project directory found"
	}

	options := []string{
		"User (~/.claude/)",
		projectLabel,
	}

	disabled := []bool{
		false,
		!projectExists,
	}

	selected, err := SelectOptionWithDisabled("Choose location", options, disabled)
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
		CurrentInstallMode = InstallModeUser
		return true
	case 1:
		CurrentInstallMode = InstallModeProject
		return true
	default:
		fmt.Println("Invalid selection. Installation cancelled.")
		return false
	}
}

// PromptForLocationForRemoval intelligently prompts for location based on what's available
// Returns true to proceed, false to cancel or nothing to remove
func PromptForLocationForRemoval(category, fileType string) bool {
	avail, err := CheckLocationAvailability(category, fileType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking locations: %v\n", err)
		return false
	}

	// No files in either location
	if !avail.HasUserLevel && !avail.HasProjectLevel {
		if fileType != "" {
			fmt.Printf("\nNo %s installed from category '%s'\n", fileType, category)
		} else if category != "" {
			fmt.Printf("\nNo files installed from category '%s'\n", category)
		} else {
			fmt.Println("\nNo files installed by foundry")
		}
		return false
	}

	// Always show both locations, disable the ones with 0 files
	fmt.Println()
	cwd, _ := os.Getwd()
	options := []string{
		fmt.Sprintf("User (~/.claude/) - %d files", avail.UserCount),
		fmt.Sprintf("Project (%s/.claude/) - %d files", cwd, avail.ProjectCount),
	}

	disabled := []bool{
		!avail.HasUserLevel,    // Disable if no user files
		!avail.HasProjectLevel, // Disable if no project files
	}

	// Map option indices to install modes
	modeMap := []InstallMode{
		InstallModeUser,
		InstallModeProject,
	}

	prompt := "Confirm location to remove from:"

	selected, err := SelectOptionWithDisabled(prompt, options, disabled)
	if err != nil {
		if err.Error() == "cancelled by user" {
			fmt.Println("Removal cancelled.")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return false
	}

	// Check if cancel was selected
	if selected >= len(modeMap) {
		fmt.Println("Removal cancelled.")
		return false
	}

	// Set the selected mode
	CurrentInstallMode = modeMap[selected]
	return true
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

	if category == "" {
		files, err = embedpkg.ListAllFiles()
		if err != nil {
			return false, fmt.Errorf("failed to list files: %w", err)
		}
	} else if fileType != "" {
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
		if category == "" {
			return false, fmt.Errorf("no installable files found")
		}
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

	// Clear screen and display banner and preview
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
	if category == "" {
		fmt.Printf("Preview: all categories [%s]\n", GetInstallModeDescription())
	} else {
		fmt.Printf("Preview: %s [%s]\n", category, GetInstallModeDescription())
	}
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

	// Ask for confirmation (use inline menu to preserve preview above)
	options := []string{
		"Yes, proceed",
		"No, cancel",
	}

	selected, err := SelectOptionInline("Proceed with installation?", options)
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

	installations := ListInstallationsForCurrentMode(st, category, fileType)
	if len(installations) == 0 {
		// No files to remove - skip preview and return true to continue
		return true, nil
	}

	// Clear screen and display banner and preview
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println(bannerStyle.Render(banner))
	if fileType != "" {
		fmt.Printf("Preview: Remove %s from %s [%s]\n", fileType, category, GetInstallModeDescription())
	} else if category == "" {
		fmt.Printf("Preview: Remove all categories [%s]\n", GetInstallModeDescription())
	} else {
		fmt.Printf("Preview: Remove category %s [%s]\n", category, GetInstallModeDescription())
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

	// Ask for confirmation (use inline menu to preserve preview above)
	options := []string{
		"Yes, remove",
		"No, cancel",
	}

	selected, err := SelectOptionInline("Proceed with removal?", options)
	if err != nil {
		return false, err
	}

	return selected == 0, nil
}

// waitModel is a simple Bubble Tea model that waits for any key press
type waitModel struct{}

func (m waitModel) Init() tea.Cmd { return nil }

func (m waitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m waitModel) View() string {
	return "\n" + helpStyle.Render("Back: Esc") + "\n"
}

// WaitForKey displays a styled prompt and waits for any key press to continue
func WaitForKey() {
	p := tea.NewProgram(waitModel{})
	p.Run()
}
