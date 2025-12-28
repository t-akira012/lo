package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var targetExtensions = map[string]bool{
	".md":  true,
	".txt": true,
	".mkd": true,
}

func isTargetFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return targetExtensions[ext]
}

func readFirstLine(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func collectFiles(dir string) ([][]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var rows [][]string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !isTargetFile(name) {
			continue
		}
		firstLine := readFirstLine(filepath.Join(dir, name))
		rows = append(rows, []string{name, firstLine})
	}
	return rows, nil
}

func renderTable(rows [][]string) string {
	headerStyle := lipgloss.NewStyle().Bold(true)

	t := table.New().
		Headers("File Name", "1st line").
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	return t.Render()
}

func run(dir string) error {
	rows, err := collectFiles(dir)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		fmt.Println("No .md, .txt, or .mkd files found.")
		return nil
	}

	fmt.Println(renderTable(rows))
	return nil
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	if err := run(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
