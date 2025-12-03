package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shapestone/claude-code-foundry/pkg/state"
)

// treeNode represents a node in the directory tree
type treeNode struct {
	label      string
	path       string
	isDir      bool
	expanded   bool
	children   []*treeNode
	fileCount  int
	depth      int
}

// treeModel represents an interactive tree view
type treeModel struct {
	nodes    []*treeNode
	cursor   int
	flatList []*treeNode // Flattened view of visible nodes
}

func (m treeModel) Init() tea.Cmd {
	return nil
}

func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				// Skip leaf files - only stop on directories
				for m.cursor > 0 && !m.flatList[m.cursor].isDir {
					m.cursor--
				}
			}
		case "down", "j":
			if m.cursor < len(m.flatList)-1 {
				m.cursor++
				// Skip leaf files - only stop on directories
				for m.cursor < len(m.flatList)-1 && !m.flatList[m.cursor].isDir {
					m.cursor++
				}
			}
		case "right", "l", "enter":
			// Expand current node (only if it has children)
			if m.cursor < len(m.flatList) {
				node := m.flatList[m.cursor]
				if node.isDir && !node.expanded && len(node.children) > 0 {
					node.expanded = true
					m.rebuildFlatList()
				}
			}
		case "left", "h":
			// Collapse current node
			if m.cursor < len(m.flatList) {
				node := m.flatList[m.cursor]
				if node.isDir && node.expanded {
					node.expanded = false
					m.rebuildFlatList()
				}
			}
		}
	}
	return m, nil
}

func (m treeModel) View() string {
	var sb strings.Builder

	// ASCII art banner at the top
	sb.WriteString(bannerStyle.Render(banner))
	sb.WriteString("\n")

	// Styled title
	title := titleStyle.Render("📁 Claude Code Directory Structure")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	for i, node := range m.flatList {
		// Cursor indicator
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("❯ ")
		}

		// Indentation
		indent := strings.Repeat("  ", node.depth)

		// Expand/collapse indicator
		indicator := ""
		if node.isDir {
			if len(node.children) == 0 {
				// Empty directory - no expand indicator
				indicator = "  "
			} else if node.expanded {
				indicator = "▼ "
			} else {
				indicator = "▶ "
			}
		} else {
			indicator = "  "
		}

		// Label with count (but skip if label already contains parentheses with info)
		label := node.label
		hasCountInfo := strings.Contains(node.label, "(") || strings.Contains(node.label, ":")

		if !hasCountInfo && node.isDir && node.fileCount > 0 {
			itemType := "files"
			if strings.Contains(node.label, "skills") {
				itemType = "skills"
			}
			if node.fileCount == 1 {
				itemType = strings.TrimSuffix(itemType, "s")
			}
			label = fmt.Sprintf("%s (%d %s)", node.label, node.fileCount, itemType)
		} else if !hasCountInfo && node.isDir && node.fileCount == 0 {
			label = fmt.Sprintf("%s (empty)", node.label)
		}

		// Add blank line BEFORE root-level sections (except the first one)
		if node.depth == 0 && i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(fmt.Sprintf("%s%s%s%s\n", cursor, indent, indicator, label))
	}

	// Styled help text at bottom
	help := helpStyle.Render("Navigate: ↑/↓  Expand: →  Collapse: ←  Quit: q")
	sb.WriteString("\n")
	sb.WriteString(help)
	sb.WriteString("\n")

	return sb.String()
}

// rebuildFlatList rebuilds the flattened view of visible nodes
func (m *treeModel) rebuildFlatList() {
	m.flatList = []*treeNode{}
	for _, node := range m.nodes {
		m.addNodeToFlatList(node)
	}
}

// addNodeToFlatList recursively adds nodes to flat list
func (m *treeModel) addNodeToFlatList(node *treeNode) {
	m.flatList = append(m.flatList, node)
	if node.expanded {
		for _, child := range node.children {
			m.addNodeToFlatList(child)
		}
	}
}

// ShowDirectoryStructure displays an interactive directory tree
func ShowDirectoryStructure() error {
	nodes, err := buildTree()
	if err != nil {
		return err
	}

	m := treeModel{
		nodes:  nodes,
		cursor: 0,
	}
	m.rebuildFlatList()

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// buildTree builds the directory tree structure
func buildTree() ([]*treeNode, error) {
	var nodes []*treeNode

	// User-level directory
	userNode, err := buildLocationNode("🏠 User-level (~/.claude/)", true, 0)
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, userNode)

	// Project-level directory
	projectNode, err := buildLocationNode("📂 Project-level (.claude/)", false, 0)
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, projectNode)

	// Installed files section
	installedNode, err := buildInstalledFilesNode()
	if err != nil {
		return nil, err
	}
	if installedNode != nil {
		nodes = append(nodes, installedNode)
	} else {
		// No installed files - show a message node
		noFilesNode := &treeNode{
			label:    "📦 Installed Files Managed by Foundry: No files installed yet",
			isDir:    false,
			expanded: false,
			depth:    0,
		}
		nodes = append(nodes, noFilesNode)
	}

	return nodes, nil
}

// buildLocationNode builds a tree node for a specific location
func buildLocationNode(label string, isUser bool, depth int) (*treeNode, error) {
	var basePath string

	if isUser {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(home, ".claude")
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		basePath = filepath.Join(cwd, ".claude")
	}

	node := &treeNode{
		label:    label,
		path:     basePath,
		isDir:    true,
		expanded: false,
		depth:    depth,
	}

	// Check if directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// Directory doesn't exist, but still create the node
		return node, nil
	}

	// Add subdirectories
	totalFiles := 0
	for _, subdir := range []string{"commands", "agents", "skills"} {
		subdirPath := filepath.Join(basePath, subdir)
		subdirNode, err := buildDirNode(subdir+"/", subdirPath, subdir == "skills", depth+1)
		if err != nil {
			continue
		}
		node.children = append(node.children, subdirNode)
		totalFiles += subdirNode.fileCount
	}
	node.fileCount = totalFiles

	return node, nil
}

// buildDirNode builds a tree node for a directory and its files
func buildDirNode(label, dirPath string, isSkillsDir bool, depth int) (*treeNode, error) {
	node := &treeNode{
		label:    label,
		path:     dirPath,
		isDir:    true,
		expanded: false,
		depth:    depth,
	}

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		node.fileCount = 0
		return node, nil
	}

	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return node, nil
	}

	// Add files/subdirectories as children
	for _, entry := range entries {
		if isSkillsDir {
			// For skills, show directories and their contents
			if entry.IsDir() {
				skillPath := filepath.Join(dirPath, entry.Name())
				skillNode := &treeNode{
					label:    entry.Name() + "/",
					path:     skillPath,
					isDir:    true,
					expanded: false,
					depth:    depth + 1,
				}

				// Count files in the skill directory
				skillEntries, err := os.ReadDir(skillPath)
				if err == nil {
					for _, skillEntry := range skillEntries {
						if !skillEntry.IsDir() {
							fileNode := &treeNode{
								label:    skillEntry.Name(),
								path:     filepath.Join(skillPath, skillEntry.Name()),
								isDir:    false,
								expanded: false,
								depth:    depth + 2,
							}
							skillNode.children = append(skillNode.children, fileNode)
							skillNode.fileCount++
						}
					}
				}

				node.children = append(node.children, skillNode)
				node.fileCount++
			}
		} else {
			// For commands/agents, show .md files
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				fileNode := &treeNode{
					label:    entry.Name(),
					path:     filepath.Join(dirPath, entry.Name()),
					isDir:    false,
					expanded: false,
					depth:    depth + 1,
				}
				node.children = append(node.children, fileNode)
				node.fileCount++
			}
		}
	}

	return node, nil
}

// buildInstalledFilesNode builds a tree node for installed files grouped by category
func buildInstalledFilesNode() (*treeNode, error) {
	st, err := state.Load()
	if err != nil {
		return nil, err
	}

	if len(st.Installations) == 0 {
		return nil, nil
	}

	// Create root node for installed files
	rootNode := &treeNode{
		label:    "📦 Installed Files Managed by Foundry",
		isDir:    true,
		expanded: false,
		depth:    0,
	}

	// Group installations by category
	byCategory := make(map[string][]state.Installation)
	for _, inst := range st.Installations {
		byCategory[inst.Category] = append(byCategory[inst.Category], inst)
	}

	// Create category nodes
	for category, installations := range byCategory {
		// Count by type
		counts := make(map[string]int)
		byType := make(map[string][]state.Installation)
		for _, inst := range installations {
			counts[inst.Type]++
			byType[inst.Type] = append(byType[inst.Type], inst)
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

		categoryLabel := fmt.Sprintf("%s: %s", category, strings.Join(countParts, ", "))
		categoryNode := &treeNode{
			label:     categoryLabel,
			isDir:     true,
			expanded:  false,
			depth:     1,
			fileCount: len(installations),
		}

		// Create type nodes (commands/, agents/, skills/)
		for _, fileType := range []string{"commands", "agents", "skills"} {
			typeInstallations := byType[fileType]

			itemType := "file"
			if fileType == "skills" {
				itemType = "skill"
			}
			if len(typeInstallations) != 1 {
				itemType += "s"
			}

			typeLabel := fileType + "/"
			if len(typeInstallations) == 0 {
				typeLabel += " (empty)"
			} else {
				typeLabel += fmt.Sprintf(" (%d %s)", len(typeInstallations), itemType)
			}

			typeNode := &treeNode{
				label:     typeLabel,
				isDir:     true,
				expanded:  false,
				depth:     2,
				fileCount: len(typeInstallations),
			}

			// Add individual files under type
			for _, inst := range typeInstallations {
				// Get filename from path
				filename := filepath.Base(inst.InstalledPath)
				if inst.Type == "skills" {
					// For skills, show the skill directory name
					filename = filepath.Base(filepath.Dir(inst.InstalledPath))
				}

				// Determine location icon
				locationIcon := "🏠"
				if home, err := os.UserHomeDir(); err == nil {
					if !strings.HasPrefix(inst.InstalledPath, home) {
						locationIcon = "📂"
					}
				}

				fileNode := &treeNode{
					label:    fmt.Sprintf("%s %s", locationIcon, filename),
					path:     inst.InstalledPath,
					isDir:    false,
					expanded: false,
					depth:    3,
				}
				typeNode.children = append(typeNode.children, fileNode)
			}

			categoryNode.children = append(categoryNode.children, typeNode)
		}

		rootNode.children = append(rootNode.children, categoryNode)
		rootNode.fileCount += len(installations)
	}

	return rootNode, nil
}

// appendLocation appends directory structure for a specific location to string builder
func appendLocation(sb *strings.Builder, label, displayPath string, isUser bool) error {
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

	sb.WriteString(fmt.Sprintf("%s (%s):\n", label, displayPath))

	// Check if directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		sb.WriteString("  ✗ Directory does not exist\n")
		return nil
	}

	// Show subdirectories with file counts
	for _, subdir := range []string{"commands", "agents", "skills"} {
		subdirPath := filepath.Join(basePath, subdir)
		count, err := countFiles(subdirPath, subdir == "skills")
		if err != nil {
			sb.WriteString(fmt.Sprintf("  %s/  (error reading: %v)\n", subdir, err))
			continue
		}

		if count == 0 {
			sb.WriteString(fmt.Sprintf("  %s/  (empty)\n", subdir))
		} else {
			itemType := "file"
			if subdir == "skills" {
				itemType = "skill"
			}
			if count != 1 {
				itemType += "s"
			}
			sb.WriteString(fmt.Sprintf("  %s/  (%d %s)\n", subdir, count, itemType))
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

// appendInstalledFiles appends installed files grouped by category to string builder
func appendInstalledFiles(sb *strings.Builder) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if len(st.Installations) == 0 {
		sb.WriteString("  No files installed by foundry yet\n")
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

		sb.WriteString(fmt.Sprintf("  %s: %s\n", category, strings.Join(countParts, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\n  Total: %d file%s installed\n", len(st.Installations), plural(len(st.Installations))))

	return nil
}

// plural returns "s" if count is not 1, otherwise ""
func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
