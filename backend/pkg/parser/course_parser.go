package parser

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/NeroQue/course-management-backend/internal/models"
	"github.com/google/uuid"
)

// FileInfo holds basic file/directory info
type FileInfo struct {
	Path         string `json:"path"`          // full path
	RelativePath string `json:"relative_path"` // relative to base dir
	Name         string `json:"name"`          // just the filename
	Size         int64  `json:"size"`          // size in bytes
	IsDir        bool   `json:"is_dir"`        // whether it's a directory
	Extension    string `json:"extension"`     // file extension
}

// CourseParser handles reading course files and converting to structured data
type CourseParser struct {
	BasePath string // where course files live
	Debug    bool   // enable extra logging
}

// NewCourseParser creates parser with base directory
func NewCourseParser(basePath string) *CourseParser {
	// Log the base path to help with debugging
	log.Printf("Initializing CourseParser with base path: %s", basePath)

	return &CourseParser{
		BasePath: basePath,
		Debug:    os.Getenv("DEBUG") == "true",
	}
}

// ValidateBasePath checks if the course directory exists and we can read it
func (p *CourseParser) ValidateBasePath() error {
	// check if directory exists
	info, err := os.Stat(p.BasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("courses directory does not exist: %s", p.BasePath)
		}
		return fmt.Errorf("error accessing courses directory: %w", err)
	}

	// make sure it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("courses path is not a directory: %s", p.BasePath)
	}

	// test if we can read it
	f, err := os.Open(p.BasePath)
	if err != nil {
		return fmt.Errorf("cannot open courses directory: %w", err)
	}
	defer f.Close()

	// try reading one entry to verify access
	_, err = f.Readdir(1)
	if err != nil && err != io.EOF {
		return fmt.Errorf("cannot read contents of courses directory: %w", err)
	}

	return nil
}

// ListCourseDirectories shows what course directories are available to import
func (p *CourseParser) ListCourseDirectories() ([]FileInfo, error) {
	var directories []FileInfo

	entries, err := os.ReadDir(p.BasePath)
	if err != nil {
		return nil, fmt.Errorf("error reading courses directory: %w", err)
	}

	// only want directories, not files
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue // skip if we can't get info
			}

			dirPath := filepath.Join(p.BasePath, entry.Name())

			directories = append(directories, FileInfo{
				Path:         dirPath,
				RelativePath: entry.Name(),
				Name:         entry.Name(),
				Size:         info.Size(),
				IsDir:        true,
				Extension:    "",
			})
		}
	}

	return directories, nil
}

// ParseCourseFolder converts a directory into a Course structure
func (p *CourseParser) ParseCourseFolder(folderPath string) (*models.Course, error) {
	// make sure folder exists
	info, err := os.Stat(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing course folder: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("specified path is not a directory: %s", folderPath)
	}

	// scan the folder structure
	modules, err := p.scanCourseFolder(folderPath)
	if err != nil {
		return nil, err
	}

	// figure out relative path
	relativePath, err := filepath.Rel(p.BasePath, folderPath)
	if err != nil {
		// if we can't get relative path, use full path
		relativePath = folderPath
	}

	course := &models.Course{
		ID:           uuid.New(),
		Title:        filepath.Base(folderPath),
		Description:  fmt.Sprintf("Course located at %s", relativePath),
		BasePath:     p.BasePath,
		RelativePath: relativePath,
		Modules:      modules,
	}

	return course, nil
}

// scanCourseFolder recursively scans folder and builds the course structure
func (p *CourseParser) scanCourseFolder(folderPath string) ([]*models.Module, error) {
	var modules []*models.Module

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading course directory: %w", err)
	}

	// look for subdirectories to turn into modules
	moduleCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			modulePath := filepath.Join(folderPath, entry.Name())
			relativePath, err := filepath.Rel(p.BasePath, modulePath)
			if err != nil {
				relativePath = modulePath
			}

			module := &models.Module{
				ID:           uuid.New(),
				Title:        entry.Name(),
				Description:  fmt.Sprintf("Module: %s", entry.Name()),
				RelativePath: relativePath,
				ContentItems: []*models.ContentItem{},
			}

			// scan for content inside this module
			contentItems, err := p.scanModuleForContentRecursive(modulePath, p.BasePath)
			if err != nil {
				log.Printf("Error scanning module %s: %v", entry.Name(), err)
			} else {
				module.ContentItems = contentItems
				log.Printf("Module '%s' found %d content items", entry.Name(), len(contentItems))
			}

			modules = append(modules, module)
			moduleCount++
		}
	}

	// if no subdirectories, treat files in this folder as one module
	if len(modules) == 0 {
		module := &models.Module{
			ID:           uuid.New(),
			Title:        "Main Content",
			Description:  "Default module for course content",
			RelativePath: filepath.Base(folderPath),
			ContentItems: []*models.ContentItem{},
		}

		contentItems, err := p.scanModuleForContentRecursive(folderPath, p.BasePath)
		if err != nil {
			return nil, fmt.Errorf("error scanning for content: %w", err)
		}

		module.ContentItems = contentItems
		modules = append(modules, module)
		log.Printf("Created default module with %d content items", len(contentItems))
	}

	log.Printf("Course parsing completed: found %d modules", len(modules))
	return modules, nil
}

// scanModuleForContentRecursive finds all the actual content files in a module
func (p *CourseParser) scanModuleForContentRecursive(modulePath, basePath string) ([]*models.ContentItem, error) {
	var contentItems []*models.ContentItem

	entries, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, fmt.Errorf("error reading module directory: %w", err)
	}

	// process each file/directory
	for i, entry := range entries {
		entryPath := filepath.Join(modulePath, entry.Name())

		if entry.IsDir() {
			// recursively scan subdirectories
			subContentItems, err := p.scanModuleForContentRecursive(entryPath, basePath)
			if err != nil {
				log.Printf("Error scanning subdirectory %s: %v", entry.Name(), err)
				continue
			}
			contentItems = append(contentItems, subContentItems...)
		} else {
			// process file
			info, err := entry.Info()
			if err != nil {
				log.Printf("Error getting info for %s: %v", entry.Name(), err)
				continue
			}

			relativePath, err := filepath.Rel(basePath, entryPath)
			if err != nil {
				relativePath = entryPath
			}

			// figure out what type of content this is
			contentType := p.determineContentType(entry.Name())

			contentItem := &models.ContentItem{
				ID:           uuid.New(),
				Title:        entry.Name(),
				Description:  fmt.Sprintf("Content file: %s", entry.Name()),
				RelativePath: relativePath,
				Size:         info.Size(),
				ContentType:  contentType,
				Order:        i, // use file order in directory
			}

			contentItems = append(contentItems, contentItem)
		}
	}

	return contentItems, nil
}

// scanModuleForContent scans module for content (kept for compatibility)
func (p *CourseParser) scanModuleForContent(modulePath string) ([]*models.ContentItem, error) {
	// just use the recursive version
	return p.scanModuleForContentRecursive(modulePath, p.BasePath)
}

// determineContentType figures out what kind of file this is based on extension
func (p *CourseParser) determineContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".mp4", ".avi", ".mov", ".mkv", ".wmv":
		return "video"
	case ".pdf":
		return "pdf"
	case ".md", ".txt":
		return "text"
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".ppt", ".pptx":
		return "presentation"
	case ".doc", ".docx":
		return "document"
	case ".xls", ".xlsx":
		return "spreadsheet"
	default:
		return "unknown"
	}
}
