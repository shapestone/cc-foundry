package installer

import (
	"fmt"

	embedpkg "github.com/shapestone/claude-code-foundry/pkg/embed"
)

// MainMenuOption represents a main menu choice
type MainMenuOption string

const (
	MainMenuShow    MainMenuOption = "show"
	MainMenuList    MainMenuOption = "list"
	MainMenuInstall MainMenuOption = "install"
	MainMenuRemove  MainMenuOption = "remove"
	MainMenuDoctor  MainMenuOption = "doctor"
	MainMenuVersion MainMenuOption = "version"
	MainMenuHelp    MainMenuOption = "help"
	MainMenuExit    MainMenuOption = "exit"
)

// ShowMainMenu displays the main interactive menu and returns the selected option
func ShowMainMenu() (MainMenuOption, error) {
	fmt.Println("\n🔧 claude-code-foundry - Manage Claude Code files\n")

	options := []string{
		"Show directory structure",
		"List available files",
		"Install files",
		"Remove files",
		"Doctor (verify & repair)",
		"Version information",
		"Help",
		"Exit",
	}

	selected, err := SelectOption("What would you like to do?", options)
	if err != nil {
		if err.Error() == "cancelled by user" {
			return MainMenuExit, nil
		}
		return "", err
	}

	// Map selection to menu option
	switch selected {
	case 0:
		return MainMenuShow, nil
	case 1:
		return MainMenuList, nil
	case 2:
		return MainMenuInstall, nil
	case 3:
		return MainMenuRemove, nil
	case 4:
		return MainMenuDoctor, nil
	case 5:
		return MainMenuVersion, nil
	case 6:
		return MainMenuHelp, nil
	case 7:
		return MainMenuExit, nil
	default:
		return "", fmt.Errorf("invalid selection")
	}
}

// ShowCategoryMenu displays available categories and returns the selected category
// action parameter is used for display purposes ("list", "install", "remove")
func ShowCategoryMenu(action string) (string, error) {
	categories, err := embedpkg.ListCategories()
	if err != nil {
		return "", fmt.Errorf("failed to list categories: %w", err)
	}

	if len(categories) == 0 {
		fmt.Println("\nNo categories available")
		return "", fmt.Errorf("no categories found")
	}

	// Build display options with category names and file counts
	var options []string
	for _, category := range categories {
		// Get file count for this category
		files, err := embedpkg.ListCategoryFiles(category)
		if err != nil {
			options = append(options, fmt.Sprintf("%s", category))
			continue
		}

		// Count by type
		counts := make(map[string]int)
		for _, file := range files {
			counts[file.Type]++
		}

		// Build count display
		countStr := fmt.Sprintf("%d commands, %d agents, %d skills",
			counts["commands"], counts["agents"], counts["skills"])

		options = append(options, fmt.Sprintf("%s (%s)", category, countStr))
	}

	// Add "All categories" option at the beginning for install/remove
	if action == "install" || action == "remove" {
		options = append([]string{"All categories"}, options...)
	}

	// Add "Back" option
	options = append(options, "← Back to main menu")

	prompt := fmt.Sprintf("Select category to %s", action)
	selected, err := SelectOption(prompt, options)
	if err != nil {
		if err.Error() == "cancelled by user" {
			return "", nil
		}
		return "", err
	}

	// Handle "All categories" selection
	if (action == "install" || action == "remove") && selected == 0 {
		return "all", nil
	}

	// Handle "Back" selection
	backIndex := len(options) - 1
	if selected == backIndex {
		return "", nil
	}

	// Adjust index if "All categories" was added
	categoryIndex := selected
	if action == "install" || action == "remove" {
		categoryIndex = selected - 1
	}

	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return "", fmt.Errorf("invalid category selection")
	}

	return categories[categoryIndex], nil
}

// ShowTypeMenu displays file types (commands, agents, skills) and returns the selected type
func ShowTypeMenu() (string, error) {
	options := []string{
		"Commands",
		"Agents",
		"Skills",
		"All types",
		"← Back",
	}

	selected, err := SelectOption("Select file type", options)
	if err != nil {
		if err.Error() == "cancelled by user" {
			return "", nil
		}
		return "", err
	}

	switch selected {
	case 0:
		return "commands", nil
	case 1:
		return "agents", nil
	case 2:
		return "skills", nil
	case 3:
		return "", nil // Empty string means all types
	case 4:
		return "", nil // Back
	default:
		return "", fmt.Errorf("invalid selection")
	}
}

// ConfirmAction asks for confirmation before proceeding with an action
func ConfirmAction(prompt string) bool {
	options := []string{
		"Yes, proceed",
		"No, cancel",
	}

	selected, err := SelectOption(prompt, options)
	if err != nil {
		return false
	}

	return selected == 0
}
