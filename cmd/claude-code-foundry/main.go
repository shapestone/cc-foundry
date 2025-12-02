package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shapestone/claude-code-foundry/embeddata"
	embedpkg "github.com/shapestone/claude-code-foundry/pkg/embed"
	"github.com/shapestone/claude-code-foundry/pkg/installer"
	"github.com/shapestone/claude-code-foundry/pkg/state"
)

const version = "1.0.0"

func init() {
	// Set the embedded filesystem for the embed package to use
	embedpkg.CategoriesFS = embeddata.Categories
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		handleList()
	case "install":
		handleInstall()
	case "remove":
		handleRemove()
	case "version", "--version", "-v":
		fmt.Printf("claude-code-foundry v%s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`claude-code-foundry - Manage Claude Code files

Usage:
  claude-code-foundry <command> [arguments]

Commands:
  list <target>              List available categories and their contents
  install <target> [type]    Install or update files from categories
  remove <target> [type]     Remove installed files
  version                    Show version information
  help                       Show this help message

List targets:
  list all                   Show all categories and their files
  list <category>            Show files in a specific category

Install/Remove targets:
  <command> all              All categories
  <command> <category>       Specific category
  <command> <category> <type>  Specific type (commands|agents|skills)

Examples:
  claude-code-foundry list all
  claude-code-foundry list development
  claude-code-foundry install all
  claude-code-foundry install development
  claude-code-foundry install development agents
  claude-code-foundry remove development skills

Note:
  The install and remove commands prompt you to choose the installation location:
  - User-level (~/.claude/) - Available across all projects
  - Project-level (.claude/) - Specific to current project

  The install command automatically updates existing files if they've changed.
  Files that are already installed and unchanged will be skipped.

Installation Locations:

User-level:
  ~/.claude/commands/           - Command files (flat .md files)
  ~/.claude/agents/             - Agent files (flat .md files)
  ~/.claude/skills/[name]/      - Skill subdirectories with SKILL.md

Project-level:
  .claude/commands/             - Command files (flat .md files)
  .claude/agents/               - Agent files (flat .md files)
  .claude/skills/[name]/        - Skill subdirectories with SKILL.md

File Naming:
  Commands/Agents: ccf-[category]-[filename].md
  Skills: ccf-[category]-[name]/SKILL.md
`)
}

func handleList() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: list command requires a target")
		fmt.Fprintln(os.Stderr, "Usage: claude-code-foundry list <all|category>")
		os.Exit(1)
	}

	target := os.Args[2]

	if target == "all" {
		listAll()
	} else {
		listCategory(target)
	}
}

func listAll() {
	categories, err := embedpkg.ListCategories()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing categories: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAvailable Categories:\n")

	for _, category := range categories {
		fmt.Printf("📁 %s/\n", category)

		files, err := embedpkg.ListCategoryFiles(category)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error listing files: %v\n", err)
			continue
		}

		// Group by type
		byType := make(map[string][]string)
		for _, file := range files {
			byType[file.Type] = append(byType[file.Type], file.Filename)
		}

		// Display by type
		for _, fileType := range []string{"commands", "agents", "skills"} {
			if files, ok := byType[fileType]; ok {
				typeLabel := strings.Title(fileType)
				fmt.Printf("  %s:\n", typeLabel)
				for _, filename := range files {
					fmt.Printf("    - %s\n", filename)
				}
			}
		}
		fmt.Println()
	}
}

func listCategory(category string) {
	files, err := embedpkg.ListCategoryFiles(category)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing category '%s': %v\n", category, err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Printf("No files found in category '%s'\n", category)
		return
	}

	fmt.Printf("\nCategory: %s\n\n", category)

	// Group by type
	byType := make(map[string][]string)
	for _, file := range files {
		byType[file.Type] = append(byType[file.Type], file.Filename)
	}

	// Display by type
	for _, fileType := range []string{"commands", "agents", "skills"} {
		if files, ok := byType[fileType]; ok {
			typeLabel := strings.Title(fileType)
			fmt.Printf("%s:\n", typeLabel)
			for _, filename := range files {
				fmt.Printf("  - %s\n", filename)
			}
			fmt.Println()
		}
	}
}

func handleInstall() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: install command requires a target")
		fmt.Fprintln(os.Stderr, "Usage: claude-code-foundry install <all|category> [type]")
		os.Exit(1)
	}

	target := os.Args[2]

	// Prompt for location
	if !installer.PromptForLocation() {
		os.Exit(0)
	}

	// Handle install all
	if target == "all" {
		// For "all", we need to preview and confirm each category
		categories, err := embedpkg.ListCategories()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing categories: %v\n", err)
			os.Exit(1)
		}

		for _, category := range categories {
			// Preview changes
			proceed, err := installer.PreviewInstall(category, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if !proceed {
				fmt.Println("Installation cancelled.")
				os.Exit(0)
			}

			// Proceed with installation
			if err := installer.InstallCategory(category); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
		return
	}

	// Handle category with optional type
	category := target
	var fileType string

	if len(os.Args) >= 4 {
		fileType = os.Args[3]

		// Validate type
		if fileType != "commands" && fileType != "agents" && fileType != "skills" {
			fmt.Fprintf(os.Stderr, "Error: invalid type '%s'. Must be: commands, agents, or skills\n", fileType)
			os.Exit(1)
		}
	}

	// Preview changes
	proceed, err := installer.PreviewInstall(category, fileType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !proceed {
		fmt.Println("Installation cancelled.")
		os.Exit(0)
	}

	// Proceed with installation
	if fileType != "" {
		// Install specific type
		if err := installer.InstallType(category, fileType); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Install entire category
		if err := installer.InstallCategory(category); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func handleRemove() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: remove command requires a target")
		fmt.Fprintln(os.Stderr, "Usage: claude-code-foundry remove <all|category> [type]")
		os.Exit(1)
	}

	target := os.Args[2]

	// Prompt for location
	if !installer.PromptForLocation() {
		os.Exit(0)
	}

	// Handle remove all
	if target == "all" {
		// Preview all removals
		st, err := state.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		installations := st.ListInstallations("", "")
		if len(installations) == 0 {
			fmt.Println("\nNo files installed by foundry")
			return
		}

		fmt.Printf("\nPreview: Remove all installed files [%s]\n", installer.GetInstallModeDescription())
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

		// Ask for confirmation with arrow-key menu
		options := []string{
			"Yes, remove",
			"No, cancel",
		}

		selected, err := installer.SelectOption("Proceed with removal?", options)
		if err != nil {
			if err.Error() == "cancelled by user" {
				fmt.Println("Removal cancelled.")
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			os.Exit(1)
		}

		if selected != 0 {
			fmt.Println("Removal cancelled.")
			os.Exit(0)
		}

		if err := installer.RemoveAll(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle category with optional type
	category := target
	var fileType string

	if len(os.Args) >= 4 {
		fileType = os.Args[3]

		// Validate type
		if fileType != "commands" && fileType != "agents" && fileType != "skills" {
			fmt.Fprintf(os.Stderr, "Error: invalid type '%s'. Must be: commands, agents, or skills\n", fileType)
			os.Exit(1)
		}
	}

	// Preview changes
	proceed, err := installer.PreviewRemove(category, fileType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !proceed {
		fmt.Println("Removal cancelled.")
		os.Exit(0)
	}

	// Proceed with removal
	if fileType != "" {
		// Remove specific type
		if err := installer.RemoveType(category, fileType); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Remove entire category
		if err := installer.RemoveCategory(category); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

