package installer

import (
	"fmt"

	embedpkg "github.com/shapestone/cc-foundry/pkg/embed"
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
// lastSelected is the index to pre-select (for cursor memory)
func ShowMainMenu(lastSelected int) (MainMenuOption, int, error) {
	fmt.Println("\n🔧 cc-foundry - Manage Claude Code files\n")

	options := []string{
		"Show directory structure",
		"List installable files",
		"Install files",
		"Remove files",
		"Doctor (verify & repair)",
		"Version information",
		"Help",
		"Exit",
	}

	selected, err := SelectOptionAt("What would you like to do?", options, lastSelected)
	if err != nil {
		if err.Error() == "cancelled by user" {
			return MainMenuExit, lastSelected, nil
		}
		return "", 0, err
	}

	// Map selection to menu option
	switch selected {
	case 0:
		return MainMenuShow, selected, nil
	case 1:
		return MainMenuList, selected, nil
	case 2:
		return MainMenuInstall, selected, nil
	case 3:
		return MainMenuRemove, selected, nil
	case 4:
		return MainMenuDoctor, selected, nil
	case 5:
		return MainMenuVersion, selected, nil
	case 6:
		return MainMenuHelp, selected, nil
	case 7:
		return MainMenuExit, selected, nil
	default:
		return "", 0, fmt.Errorf("invalid selection")
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
