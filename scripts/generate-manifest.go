package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shapestone/shape-yaml/pkg/yaml"
)

// Frontmatter represents the YAML frontmatter in markdown files
type Frontmatter struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tools       []string `yaml:"tools,omitempty"`
}

// FileEntry represents a single file in the manifest
type FileEntry struct {
	Name        string `json:"name"`
	File        string `json:"file"`
	Description string `json:"description"`
	SHA256      string `json:"sha256"`
}

// CategoryFiles represents files grouped by type (commands, agents, skills)
type CategoryFiles struct {
	Commands []FileEntry `json:"commands,omitempty"`
	Agents   []FileEntry `json:"agents,omitempty"`
	Skills   []FileEntry `json:"skills,omitempty"`
}

// Manifest represents the complete manifest structure
type Manifest struct {
	Version    string                   `json:"version"`
	Generated  string                   `json:"generated"`
	Categories map[string]CategoryFiles `json:"categories"`
	Bundle     BundleInfo               `json:"bundle"`
}

// BundleInfo contains information about the bundle archive
type BundleInfo struct {
	File   string `json:"file"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	filesDir := "files"
	manifestPath := filepath.Join(filesDir, "manifest.json")
	bundlePath := filepath.Join(filesDir, "bundle.tar.gz")

	fmt.Println("🔍 Scanning files directory...")

	// Check if files directory exists
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		return fmt.Errorf("files directory does not exist: %s", filesDir)
	}

	// Build manifest by scanning files
	manifest := &Manifest{
		Version:    "1.0",
		Generated:  time.Now().UTC().Format(time.RFC3339),
		Categories: make(map[string]CategoryFiles),
	}

	// Walk through files directory
	err := filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the files directory itself, manifest.json, and bundle.tar.gz
		if path == filesDir || filepath.Base(path) == "manifest.json" || filepath.Base(path) == "bundle.tar.gz" {
			return nil
		}

		// Only process .md files
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Parse the path to extract category and type
		// Expected format: files/category/type/filename.md
		relPath, err := filepath.Rel(filesDir, path)
		if err != nil {
			return err
		}

		parts := strings.Split(filepath.ToSlash(relPath), "/")
		if len(parts) < 3 {
			fmt.Printf("⚠️  Skipping %s (unexpected path structure)\n", relPath)
			return nil
		}

		category := parts[0]
		fileType := parts[1]

		// Validate type
		if fileType != "commands" && fileType != "agents" && fileType != "skills" {
			fmt.Printf("⚠️  Skipping %s (unknown type: %s)\n", relPath, fileType)
			return nil
		}

		// Read and parse file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Calculate SHA-256
		hash := fmt.Sprintf("%x", sha256.Sum256(content))

		// Parse frontmatter
		frontmatter, err := parseFrontmatter(content)
		if err != nil {
			fmt.Printf("⚠️  Warning: failed to parse frontmatter in %s: %v\n", relPath, err)
			// Continue anyway with default values
			frontmatter = &Frontmatter{
				Name:        strings.TrimSuffix(filepath.Base(path), ".md"),
				Description: "No description available",
			}
		}

		// Create file entry
		entry := FileEntry{
			Name:        frontmatter.Name,
			File:        relPath,
			Description: frontmatter.Description,
			SHA256:      hash,
		}

		// Add to manifest
		if _, exists := manifest.Categories[category]; !exists {
			manifest.Categories[category] = CategoryFiles{}
		}

		catFiles := manifest.Categories[category]
		switch fileType {
		case "commands":
			catFiles.Commands = append(catFiles.Commands, entry)
		case "agents":
			catFiles.Agents = append(catFiles.Agents, entry)
		case "skills":
			catFiles.Skills = append(catFiles.Skills, entry)
		}
		manifest.Categories[category] = catFiles

		fmt.Printf("✓ Added %s/%s: %s\n", category, fileType, frontmatter.Name)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan files: %w", err)
	}

	// Create bundle
	fmt.Println("\n📦 Creating bundle archive...")
	bundleInfo, err := createBundle(filesDir, bundlePath)
	if err != nil {
		return fmt.Errorf("failed to create bundle: %w", err)
	}
	manifest.Bundle = *bundleInfo

	// Write manifest
	fmt.Println("\n📝 Writing manifest.json...")
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	// Print summary
	fmt.Println("\n✓ Manifest generation complete!")
	fmt.Printf("  Categories: %d\n", len(manifest.Categories))

	totalFiles := 0
	for _, cat := range manifest.Categories {
		totalFiles += len(cat.Commands) + len(cat.Agents) + len(cat.Skills)
	}
	fmt.Printf("  Total files: %d\n", totalFiles)
	fmt.Printf("  Bundle size: %.2f KB\n", float64(bundleInfo.Size)/1024)
	fmt.Printf("  Output: %s\n", manifestPath)

	return nil
}

// parseFrontmatter extracts and parses YAML frontmatter from markdown content
func parseFrontmatter(content []byte) (*Frontmatter, error) {
	// Frontmatter is between --- delimiters at the start of the file
	lines := strings.Split(string(content), "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return nil, fmt.Errorf("no frontmatter found")
	}

	// Find the closing ---
	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return nil, fmt.Errorf("frontmatter not closed")
	}

	// Extract frontmatter content
	yamlContent := strings.Join(lines[1:endIndex], "\n")

	// Parse YAML
	var frontmatter Frontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &frontmatter); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	return &frontmatter, nil
}

// createBundle creates a tar.gz archive of all files
func createBundle(filesDir, bundlePath string) (*BundleInfo, error) {
	// Create bundle file
	bundleFile, err := os.Create(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundleFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(bundleFile)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk through files and add to archive
	err = filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the files directory itself and generated files
		if path == filesDir || filepath.Base(path) == "manifest.json" || filepath.Base(path) == "bundle.tar.gz" {
			return nil
		}

		// Skip directories and non-.md files
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(filesDir, path)
		if err != nil {
			return err
		}

		// Create tar header
		header := &tar.Header{
			Name:    filepath.ToSlash(relPath),
			Size:    info.Size(),
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(tarWriter, file); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Close writers to flush
	tarWriter.Close()
	gzipWriter.Close()
	bundleFile.Close()

	// Calculate SHA-256 of bundle
	bundleContent, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(bundleContent))

	// Get file info for size
	bundleInfo, err := os.Stat(bundlePath)
	if err != nil {
		return nil, err
	}

	return &BundleInfo{
		File:   "bundle.tar.gz",
		SHA256: hash,
		Size:   bundleInfo.Size(),
	}, nil
}
