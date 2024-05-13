package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: text2code <input_text_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	err := recreateFilesFromText(inputFile)
	if err != nil {
		fmt.Printf("Failed to recreate files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Files recreated successfully.")
}

func recreateFilesFromText(inputTextFile string) error {
	file, err := os.Open(inputTextFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentPath string
	var currentContent []string
	processingFile := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "-----BEGIN ") {
			currentPath = strings.TrimPrefix(line, "-----BEGIN ")
			currentPath = strings.TrimSuffix(currentPath, "-----")
			currentContent = nil // Initialize new slice for the file's content
			processingFile = true
		} else if strings.HasPrefix(line, "-----END ") {
			if processingFile && currentPath != "" {
				err := writeContentToFile(currentPath, currentContent)
				if err != nil {
					return err
				}
			}
			processingFile = false
		} else if processingFile {
			currentContent = append(currentContent, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func writeContentToFile(path string, content []string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create the file
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the content to the file
	for _, line := range content {
		_, err := outputFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
