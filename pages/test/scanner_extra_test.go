package pages_test

import (
	"testing"

	"github.com/laranatech/electrostatic/pages"
)

func TestIsSkipped_UtilityPages(t *testing.T) {
	utilityPages := []string{"404.md", "403.md", "500.md"}

	for _, name := range utilityPages {
		if !pages.IsSkipped(name) {
			t.Errorf("expected IsSkipped(%q) = true", name)
		}
	}
}

func TestIsSkipped_RegularPages(t *testing.T) {
	regularPages := []string{
		"index.md",
		"about.md",
		"blog.md",
		"larana.md",
		"404-extra.md",
		"my-500.md",
	}

	for _, name := range regularPages {
		if pages.IsSkipped(name) {
			t.Errorf("expected IsSkipped(%q) = false", name)
		}
	}
}

func TestFilterUtilityPages_RemovesErrorPages(t *testing.T) {
	input := []pages.Page{
		{Filepath: "/root/index.md"},
		{Filepath: "/root/about.md"},
		{Filepath: "/root/404.md"},
		{Filepath: "/root/403.md"},
		{Filepath: "/root/500.md"},
		{Filepath: "/root/blog.md"},
	}

	result := pages.FilterUtilityPages(input)

	if len(result) != 3 {
		t.Errorf("expected 3 pages after filtering, got %d", len(result))
	}

	for _, p := range result {
		name := p.Filepath[len(p.Filepath)-len("404.md"):]
		switch name {
		case "404.md", "403.md", "500.md":
			t.Errorf("utility page %q should have been filtered out", p.Filepath)
		}
	}
}

func TestFilterUtilityPages_EmptyInput(t *testing.T) {
	result := pages.FilterUtilityPages([]pages.Page{})

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d pages", len(result))
	}
}

func TestFilterUtilityPages_NoUtilityPages(t *testing.T) {
	input := []pages.Page{
		{Filepath: "/root/index.md"},
		{Filepath: "/root/about.md"},
		{Filepath: "/root/contact.md"},
	}

	result := pages.FilterUtilityPages(input)

	if len(result) != len(input) {
		t.Errorf("expected %d pages, got %d", len(input), len(result))
	}
}

func TestFilterUtilityPages_OnlyUtilityPages(t *testing.T) {
	input := []pages.Page{
		{Filepath: "/root/404.md"},
		{Filepath: "/root/403.md"},
		{Filepath: "/root/500.md"},
	}

	result := pages.FilterUtilityPages(input)

	if len(result) != 0 {
		t.Errorf("expected 0 pages after filtering all utility pages, got %d", len(result))
	}
}

func TestScanAllFilepaths_IncludesUtilityPages(t *testing.T) {
	root := "./root"

	paths, err := pages.ScanAllFilepaths(root)

	if err != nil {
		t.Fatal(err)
	}

	found404 := false
	for _, p := range paths {
		if p == root+"/404.md" {
			found404 = true
		}
	}

	if !found404 {
		t.Error("ScanAllFilepaths should include 404.md")
	}
}
