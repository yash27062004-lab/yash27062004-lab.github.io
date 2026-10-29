package render

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func RenderPage(
	files []string,
	destination string,
	data any,
) error {
	if len(files) == 0 {
		return fmt.Errorf("at least one template file is required")
	}

	if strings.TrimSpace(destination) == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	executeErr := tmpl.ExecuteTemplate(file, "layout", data)
	closeErr := file.Close()

	if executeErr != nil {
		return fmt.Errorf("render template: %w", executeErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close output file: %w", closeErr)
	}

	return nil
}

func CopyDir(sourceDir, destinationDir string) error {
	if strings.TrimSpace(sourceDir) == "" {
		return fmt.Errorf("source directory cannot be empty")
	}

	if strings.TrimSpace(destinationDir) == "" {
		return fmt.Errorf("destination directory cannot be empty")
	}

	sourceAbs, err := filepath.Abs(filepath.Clean(sourceDir))
	if err != nil {
		return fmt.Errorf("resolve source directory: %w", err)
	}

	destinationAbs, err := filepath.Abs(filepath.Clean(destinationDir))
	if err != nil {
		return fmt.Errorf("resolve destination directory: %w", err)
	}

	if pathsOverlap(sourceAbs, destinationAbs) {
		return fmt.Errorf("source and destination directories cannot overlap")
	}

	info, err := os.Stat(sourceAbs)
	if err != nil {
		return fmt.Errorf("read source directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("source path is not a directory")
	}

	return filepath.Walk(
		sourceAbs,
		func(
			sourcePath string,
			info os.FileInfo,
			walkErr error,
		) error {
			if walkErr != nil {
				return walkErr
			}

			relativePath, err := filepath.Rel(sourceAbs, sourcePath)
			if err != nil {
				return err
			}

			targetPath := filepath.Join(destinationAbs, relativePath)

			if info.IsDir() {
				return os.MkdirAll(targetPath, info.Mode().Perm())
			}

			return copyFile(sourcePath, targetPath)
		},
	)
}

func copyFile(sourcePath, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer source.Close()

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return fmt.Errorf("create destination directory: %w", err)
	}

	destination, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}

	_, copyErr := io.Copy(destination, source)
	closeErr := destination.Close()

	if copyErr != nil {
		return fmt.Errorf("copy file contents: %w", copyErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close destination file: %w", closeErr)
	}

	return nil
}

func pathsOverlap(first, second string) bool {
	return directoryContains(first, second) ||
		directoryContains(second, first)
}

func directoryContains(parent, child string) bool {
	relative, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	if relative == "." {
		return true
	}

	if relative == ".." || filepath.IsAbs(relative) {
		return false
	}

	return !strings.HasPrefix(
		relative,
		".."+string(os.PathSeparator),
	)
}
