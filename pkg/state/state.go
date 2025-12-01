package state

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	StateFile = ".claude-code-foundry.json"
	Version   = "1.0.0"
)

// State represents the foundry installation state
type State struct {
	Version       string         `json:"version"`
	Installations []Installation `json:"installations"`
}

// Installation represents a single installed file
type Installation struct {
	Category      string    `json:"category"`
	Type          string    `json:"type"`
	File          string    `json:"file"`
	InstalledPath string    `json:"installed_path"`
	Hash          string    `json:"hash"`
	InstalledAt   time.Time `json:"installed_at"`
}

// Load loads the state file from the user's home directory
func Load() (*State, error) {
	stateFilePath, err := GetStateFilePath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty state
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return &State{
			Version:       Version,
			Installations: []Installation{},
		}, nil
	}

	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// Save saves the state file to the user's home directory
func (s *State) Save() error {
	stateFilePath, err := GetStateFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// AddInstallation adds a new installation to the state
func (s *State) AddInstallation(category, fileType, filename, installedPath string, content []byte) {
	hash := calculateHash(content)

	installation := Installation{
		Category:      category,
		Type:          fileType,
		File:          filename,
		InstalledPath: installedPath,
		Hash:          hash,
		InstalledAt:   time.Now(),
	}

	s.Installations = append(s.Installations, installation)
}

// RemoveInstallation removes an installation from the state
func (s *State) RemoveInstallation(installedPath string) {
	var filtered []Installation
	for _, inst := range s.Installations {
		if inst.InstalledPath != installedPath {
			filtered = append(filtered, inst)
		}
	}
	s.Installations = filtered
}

// FindInstallation finds an installation by its installed path
func (s *State) FindInstallation(installedPath string) *Installation {
	for _, inst := range s.Installations {
		if inst.InstalledPath == installedPath {
			return &inst
		}
	}
	return nil
}

// ListInstallations returns all installations, optionally filtered by category and/or type
func (s *State) ListInstallations(category, fileType string) []Installation {
	var filtered []Installation

	for _, inst := range s.Installations {
		if category != "" && inst.Category != category {
			continue
		}
		if fileType != "" && inst.Type != fileType {
			continue
		}
		filtered = append(filtered, inst)
	}

	return filtered
}

// GetStateFilePath returns the full path to the state file
func GetStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, StateFile), nil
}

// calculateHash calculates SHA-256 hash of content
func calculateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}

// HasContentChanged checks if file content has changed from what was installed
func (i *Installation) HasContentChanged(newContent []byte) bool {
	newHash := calculateHash(newContent)
	return newHash != i.Hash
}
