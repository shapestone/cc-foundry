package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shapestone/claude-code-foundry/embeddata"
	"github.com/shapestone/claude-code-foundry/pkg/doctor"
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
	// Interactive mode - show main menu
	if len(os.Args) < 2 {
		runInteractiveMode()
		return
	}

	// For future: support command-line arguments for scripting
	// For now, always run interactive mode
	runInteractiveMode()
}

func runInteractiveMode() {
	for {
		option, err := installer.ShowMainMenu()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		switch option {
		case installer.MainMenuShow:
			handleShow()
		case installer.MainMenuList:
			handleListInteractive()
		case installer.MainMenuInstall:
			handleInstallInteractive()
		case installer.MainMenuRemove:
			handleRemoveInteractive()
		case installer.MainMenuDoctor:
			handleDoctor()
		case installer.MainMenuVersion:
			fmt.Printf("\nclaude-code-foundry v%s\n", version)
			installer.WaitForKey()
		case installer.MainMenuHelp:
			printUsage()
			installer.WaitForKey()
		case installer.MainMenuExit:
			fmt.Println("\nGoodbye! 👋\n")
			return
		}
	}
}

// handleShow displays the directory structure
func handleShow() {
	if err := installer.ShowDirectoryStructure(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

// handleListInteractive handles the interactive list flow
func handleListInteractive() {
	category, err := installer.ShowCategoryMenu("list")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// User chose to go back
	if category == "" {
		return
	}

	// List all or specific category
	if category == "all" {
		listAll()
	} else {
		listCategory(category)
	}

	installer.WaitForKey()
}

// handleInstallInteractive handles the interactive install flow
func handleInstallInteractive() {
	// Select category
	category, err := installer.ShowCategoryMenu("install")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// User chose to go back
	if category == "" {
		return
	}

	// Prompt for location
	if !installer.PromptForLocation() {
		return
	}

	// Handle install all
	if category == "all" {
		categories, err := embedpkg.ListCategories()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing categories: %v\n", err)
			return
		}

		for _, cat := range categories {
			proceed, err := installer.PreviewInstall(cat, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}

			if !proceed {
				fmt.Println("Installation cancelled.")
				return
			}

			if err := installer.InstallCategory(cat); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
		}
		return
	}

	// Preview and install specific category
	proceed, err := installer.PreviewInstall(category, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if !proceed {
		fmt.Println("Installation cancelled.")
		return
	}

	if err := installer.InstallCategory(category); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}

// handleRemoveInteractive handles the interactive remove flow
func handleRemoveInteractive() {
	// Select category
	category, err := installer.ShowCategoryMenu("remove")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// User chose to go back
	if category == "" {
		return
	}

	// Prompt for location
	if !installer.PromptForLocation() {
		return
	}

	// Handle remove all
	if category == "all" {
		st, err := state.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
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

		if !installer.ConfirmAction("Proceed with removal?") {
			fmt.Println("Removal cancelled.")
			return
		}

		if err := installer.RemoveAll(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		return
	}

	// Preview and remove specific category
	proceed, err := installer.PreviewRemove(category, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if !proceed {
		fmt.Println("Removal cancelled.")
		return
	}

	if err := installer.RemoveCategory(category); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}

// handleDoctor runs the doctor diagnostics
func handleDoctor() {
	report, err := doctor.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running doctor: %v\n", err)
		return
	}

	doctor.PrintReport(report)

	// Offer to fix issues
	if err := doctor.OfferFixes(report, installer.SelectOption); err != nil {
		fmt.Fprintf(os.Stderr, "Error offering fixes: %v\n", err)
		return
	}
}

func printUsage() {
	fmt.Println(`🔧 claude-code-foundry - Manage Claude Code files

Interactive Mode:
  Just run: claude-code-foundry

  The tool will guide you through an interactive menu to:
  - Show directory structure and installed files
  - List available commands, agents, and skills
  - Install files to ~/.claude/ or .claude/
  - Remove installed files
  - Run diagnostics and repair (doctor)

Installation Locations:

User-level (~/.claude/):
  - Available across all projects
  ~/.claude/commands/           Command files (flat .md files)
  ~/.claude/agents/             Agent files (flat .md files)
  ~/.claude/skills/[name]/      Skill subdirectories with SKILL.md

Project-level (.claude/):
  - Specific to current project, can be version-controlled
  .claude/commands/             Command files (flat .md files)
  .claude/agents/               Agent files (flat .md files)
  .claude/skills/[name]/        Skill subdirectories with SKILL.md

File Naming:
  Commands/Agents: ccf-[category]-[filename].md
  Skills: ccf-[category]-[name]/SKILL.md

Note: Non-interactive mode for scripting will be added in a future release.
`)
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

