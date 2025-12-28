package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsTargetFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"md file", "README.md", true},
		{"txt file", "notes.txt", true},
		{"mkd file", "doc.mkd", true},
		{"MD uppercase", "README.MD", true},
		{"go file", "main.go", false},
		{"no extension", "Makefile", false},
		{"hidden md", ".hidden.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTargetFile(tt.filename)
			if got != tt.want {
				t.Errorf("isTargetFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestReadFirstLine(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"single line", "Hello World", "Hello World"},
		{"multi line", "First\nSecond\nThird", "First"},
		{"empty file", "", ""},
		{"japanese", "日本語テスト\n2行目", "日本語テスト"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name+".txt")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			got := readFirstLine(path)
			if got != tt.want {
				t.Errorf("readFirstLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadFirstLineNotExist(t *testing.T) {
	got := readFirstLine("/nonexistent/file.txt")
	if got != "" {
		t.Errorf("readFirstLine(nonexistent) = %q, want empty", got)
	}
}

func TestCollectFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	files := map[string]string{
		"readme.md":  "# README",
		"notes.txt":  "My notes",
		"doc.mkd":    "Documentation",
		"main.go":    "package main",
		"data.json":  `{"key": "value"}`,
	}

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a subdirectory (should be ignored)
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	rows, err := collectFiles(dir)
	if err != nil {
		t.Fatalf("collectFiles() error = %v", err)
	}

	if len(rows) != 3 {
		t.Errorf("collectFiles() returned %d rows, want 3", len(rows))
	}

	// Check that only target files are included
	foundFiles := make(map[string]string)
	for _, row := range rows {
		foundFiles[row[0]] = row[1]
	}

	expected := map[string]string{
		"doc.mkd":   "Documentation",
		"notes.txt": "My notes",
		"readme.md": "# README",
	}

	for name, content := range expected {
		if got, ok := foundFiles[name]; !ok {
			t.Errorf("missing file %q", name)
		} else if got != content {
			t.Errorf("file %q: got %q, want %q", name, got, content)
		}
	}
}

func TestCollectFilesEmpty(t *testing.T) {
	dir := t.TempDir()

	rows, err := collectFiles(dir)
	if err != nil {
		t.Fatalf("collectFiles() error = %v", err)
	}

	if len(rows) != 0 {
		t.Errorf("collectFiles() returned %d rows, want 0", len(rows))
	}
}

func TestCollectFilesNotExist(t *testing.T) {
	_, err := collectFiles("/nonexistent/directory")
	if err == nil {
		t.Error("collectFiles(nonexistent) should return error")
	}
}

func TestRenderTable(t *testing.T) {
	rows := [][]string{
		{"test.md", "# Test"},
		{"日本語.txt", "日本語コンテンツ"},
	}

	output := renderTable(rows)

	// Check that output contains expected content
	if output == "" {
		t.Error("renderTable() returned empty string")
	}

	// Check headers exist
	if !contains(output, "File Name") {
		t.Error("output missing 'File Name' header")
	}
	if !contains(output, "1st line") {
		t.Error("output missing '1st line' header")
	}

	// Check data exists
	if !contains(output, "test.md") {
		t.Error("output missing 'test.md'")
	}
	if !contains(output, "# Test") {
		t.Error("output missing '# Test'")
	}
	if !contains(output, "日本語.txt") {
		t.Error("output missing '日本語.txt'")
	}
}

func TestRenderSimple(t *testing.T) {
	rows := [][]string{
		{"test.md", "# Test"},
		{"日本語.txt", "日本語コンテンツ"},
	}

	output := renderSimple(rows)

	expected := "\"test.md\"\t# Test\n\"日本語.txt\"\t日本語コンテンツ\n"
	if output != expected {
		t.Errorf("renderSimple() = %q, want %q", output, expected)
	}
}

func TestRenderSimpleEmpty(t *testing.T) {
	rows := [][]string{}
	output := renderSimple(rows)

	if output != "" {
		t.Errorf("renderSimple([]) = %q, want empty", output)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
