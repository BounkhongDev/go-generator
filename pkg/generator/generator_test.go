package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit_InvalidName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "empty"},
		{"spaces only", "   ", "empty"},
		{"contains space", "my project", "spaces"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.input)
			if err == nil {
				t.Fatal("Init expected an error, got nil")
			}
			msg := err.Error()
			if tt.want == "empty" && msg != "project name must not be empty" {
				t.Errorf("Init(%q) error = %q, want something about empty", tt.input, msg)
			}
			if tt.want == "spaces" && !strings.Contains(msg, "space") {
				t.Errorf("Init(%q) error = %q, want something about space", tt.input, msg)
			}
		})
	}
}

func TestToPlural(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"user", "users"},
		{"category", "categories"},
		{"box", "boxes"},
		{"match", "matches"},
		{"dish", "dishes"},
		{"example", "examples"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toPlural(tt.input)
			if got != tt.want {
				t.Errorf("toPlural(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsVowel(t *testing.T) {
	for _, c := range "aeiou" {
		if !isVowel(byte(c)) {
			t.Errorf("isVowel(%c) = false, want true", c)
		}
	}
	for _, c := range "bcdfg" {
		if isVowel(byte(c)) {
			t.Errorf("isVowel(%c) = true, want false", c)
		}
	}
}

func TestGetProjectName(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// No go.mod
	_, err := getProjectName()
	if err == nil {
		t.Error("getProjectName with no go.mod expected error, got nil")
	}

	// Valid go.mod
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/myapp\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := getProjectName()
	if err != nil {
		t.Fatal(err)
	}
	if got != "example.com/myapp" {
		t.Errorf("getProjectName() = %q, want example.com/myapp", got)
	}
}

func TestGenerateModule(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/proj\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	err := GenerateModule("coverage")
	if err != nil {
		t.Fatalf("GenerateModule(\"coverage\") = %v", err)
	}

	// Check that key files were created
	wantFiles := []string{
		"internal/models/coverage.go",
		"internal/requests/coverage_request.go",
		"internal/responses/coverage_response.go",
		"internal/repositories/coverage_repository.go",
		"internal/services/coverage_service.go",
		"internal/controllers/coverage_controller.go",
		"tests/services/coverage_service_test.go",
		"tests/mocks/coverage_repository_mock.go",
		"tests/fixtures/coverage_fixture.go",
		"migrations/migrations.go",
	}
	for _, rel := range wantFiles {
		p := filepath.Join(dir, rel)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist: %v", rel, err)
		}
	}
}

func ExampleInit() {
	// Init validates the project name and returns an error for invalid input.
	// Run from an empty directory to actually create a project.
	err := Init("")
	if err != nil {
		// expected: "project name must not be empty"
	}
	_ = err
	// Output:
}

func ExampleGenerateModule() {
	// GenerateModule creates a new module (model, repo, service, controller, etc.).
	// It reads the project name from go.mod in the current directory.
	dir, _ := os.MkdirTemp("", "go-gen-example-")
	defer os.RemoveAll(dir)
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	_ = os.Chdir(dir)
	_ = os.WriteFile("go.mod", []byte("module example.com/app\n\ngo 1.21\n"), 0644)

	err := GenerateModule("product")
	if err != nil {
		panic(err)
	}
	// Example runs without output check (generator prints to stdout).
}
