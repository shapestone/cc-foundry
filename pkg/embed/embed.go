package embed

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// CategoriesFS is the embedded filesystem containing categories
// This must be set by the main package after embedding
var CategoriesFS fs.FS

// CategoryFile represents a file within a category
type CategoryFile struct {
	Category  string
	Type      string // "commands", "agents", or "skills"
	Filename  string
	SubPath   string // for skill support files: relative subdir (e.g., "references"); empty for top-level files
	SkillName string // for skill support files: the source skill .md filename (e.g., "ddd-skill.md"); empty for main skill files
	Content   []byte
}

// ListCategories returns all available categories
func ListCategories() ([]string, error) {
	entries, err := fs.ReadDir(CategoriesFS, "categories")
	if err != nil {
		return nil, err
	}

	var categories []string
	for _, entry := range entries {
		if entry.IsDir() {
			categories = append(categories, entry.Name())
		}
	}
	return categories, nil
}

// ListCategoryFiles returns all files in a specific category
func ListCategoryFiles(category string) ([]CategoryFile, error) {
	var files []CategoryFile

	categoryPath := filepath.Join("categories", category)

	// Check each type directory (commands, agents, skills)
	for _, fileType := range []string{"commands", "agents", "skills"} {
		typePath := filepath.Join(categoryPath, fileType)

		entries, err := fs.ReadDir(CategoriesFS, typePath)
		if err != nil {
			// Directory doesn't exist for this type, skip
			continue
		}

		if fileType == "skills" {
			// For skills, collect top-level skill files and subdirectory support files separately
			var topLevelSkills []string
			type subDirFile struct {
				subPath  string
				filename string
				content  []byte
			}
			var subDirFiles []subDirFile

			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					content, err := fs.ReadFile(CategoriesFS, filepath.Join(typePath, entry.Name()))
					if err != nil {
						return nil, err
					}
					topLevelSkills = append(topLevelSkills, entry.Name())
					files = append(files, CategoryFile{
						Category: category,
						Type:     fileType,
						Filename: entry.Name(),
						Content:  content,
					})
				} else if entry.IsDir() {
					subPath := entry.Name()
					subEntries, err := fs.ReadDir(CategoriesFS, filepath.Join(typePath, subPath))
					if err != nil {
						continue
					}
					for _, subEntry := range subEntries {
						if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".md") {
							content, err := fs.ReadFile(CategoriesFS, filepath.Join(typePath, subPath, subEntry.Name()))
							if err != nil {
								return nil, err
							}
							subDirFiles = append(subDirFiles, subDirFile{subPath, subEntry.Name(), content})
						}
					}
				}
			}

			// Associate subdirectory files with each skill
			for _, subFile := range subDirFiles {
				for _, skillFilename := range topLevelSkills {
					files = append(files, CategoryFile{
						Category:  category,
						Type:      fileType,
						Filename:  subFile.filename,
						SubPath:   subFile.subPath,
						SkillName: skillFilename,
						Content:   subFile.content,
					})
				}
			}
		} else {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					content, err := fs.ReadFile(CategoriesFS, filepath.Join(typePath, entry.Name()))
					if err != nil {
						return nil, err
					}

					files = append(files, CategoryFile{
						Category: category,
						Type:     fileType,
						Filename: entry.Name(),
						Content:  content,
					})
				}
			}
		}
	}

	return files, nil
}

// ListAllFiles returns all files across all categories
func ListAllFiles() ([]CategoryFile, error) {
	categories, err := ListCategories()
	if err != nil {
		return nil, err
	}

	var allFiles []CategoryFile
	for _, cat := range categories {
		files, err := ListCategoryFiles(cat)
		if err != nil {
			return nil, err
		}
		allFiles = append(allFiles, files...)
	}
	return allFiles, nil
}

// ListTypeFiles returns all files of a specific type in a category
func ListTypeFiles(category, fileType string) ([]CategoryFile, error) {
	var files []CategoryFile

	typePath := filepath.Join("categories", category, fileType)

	entries, err := fs.ReadDir(CategoriesFS, typePath)
	if err != nil {
		return nil, err
	}

	if fileType == "skills" {
		var topLevelSkills []string
		type subDirFile struct {
			subPath  string
			filename string
			content  []byte
		}
		var subDirFiles []subDirFile

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				content, err := fs.ReadFile(CategoriesFS, filepath.Join(typePath, entry.Name()))
				if err != nil {
					return nil, err
				}
				topLevelSkills = append(topLevelSkills, entry.Name())
				files = append(files, CategoryFile{
					Category: category,
					Type:     fileType,
					Filename: entry.Name(),
					Content:  content,
				})
			} else if entry.IsDir() {
				subPath := entry.Name()
				subEntries, err := fs.ReadDir(CategoriesFS, filepath.Join(typePath, subPath))
				if err != nil {
					continue
				}
				for _, subEntry := range subEntries {
					if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".md") {
						content, err := fs.ReadFile(CategoriesFS, filepath.Join(typePath, subPath, subEntry.Name()))
						if err != nil {
							return nil, err
						}
						subDirFiles = append(subDirFiles, subDirFile{subPath, subEntry.Name(), content})
					}
				}
			}
		}

		for _, subFile := range subDirFiles {
			for _, skillFilename := range topLevelSkills {
				files = append(files, CategoryFile{
					Category:  category,
					Type:      fileType,
					Filename:  subFile.filename,
					SubPath:   subFile.subPath,
					SkillName: skillFilename,
					Content:   subFile.content,
				})
			}
		}
	} else {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				content, err := fs.ReadFile(CategoriesFS, filepath.Join(typePath, entry.Name()))
				if err != nil {
					return nil, err
				}

				files = append(files, CategoryFile{
					Category: category,
					Type:     fileType,
					Filename: entry.Name(),
					Content:  content,
				})
			}
		}
	}

	return files, nil
}

// GetFile retrieves a specific file's content
func GetFile(category, fileType, filename string) (*CategoryFile, error) {
	path := filepath.Join("categories", category, fileType, filename)

	content, err := fs.ReadFile(CategoriesFS, path)
	if err != nil {
		return nil, err
	}

	return &CategoryFile{
		Category: category,
		Type:     fileType,
		Filename: filename,
		Content:  content,
	}, nil
}
