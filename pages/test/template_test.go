package pages_test

import (
	"strings"
	"testing"

	"github.com/laranatech/electrostatic/pages"
)

func TestReadTemplateFile(t *testing.T) {
	root := "./root"

	tmpl, err := pages.ReadTemplateFile(root)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(tmpl, "%CONTENT%") {
		t.Error("template should contain the CONTENT placeholder")
	}

	if !strings.Contains(tmpl, "%title%") {
		t.Error("template should contain the title placeholder")
	}
}

func TestReadTemplateFile_MissingFile(t *testing.T) {
	_, err := pages.ReadTemplateFile("./nonexistent")

	if err == nil {
		t.Error("expected error for missing template.html, got nil")
	}
}

func TestFormatTemplate_ReplacesContent(t *testing.T) {
	tmpl := "<html><body>%CONTENT%</body></html>"

	page := pages.Page{
		Content: []byte("# Hello"),
		Meta:    map[string]string{},
	}

	result := pages.FormatTemplate(tmpl, page)

	if strings.Contains(result, "%CONTENT%") {
		t.Error("CONTENT placeholder was not replaced")
	}

	if !strings.Contains(result, "<h1") {
		t.Errorf("expected rendered heading in output, got: %s", result)
	}
}

func TestFormatTemplate_ReplacesMeta(t *testing.T) {
	tmpl := "<title>%title%</title><meta name='description' content='%description%'><meta name='keywords' content='%keywords%'>"

	page := pages.Page{
		Content: []byte("text"),
		Meta: map[string]string{
			"title":       "My Title",
			"description": "My Desc",
			"keywords":    "kw1, kw2",
		},
	}

	result := pages.FormatTemplate(tmpl, page)

	if strings.Contains(result, "%title%") {
		t.Error("title placeholder was not replaced")
	}

	if !strings.Contains(result, "My Title") {
		t.Errorf("title not found in result: %s", result)
	}

	if !strings.Contains(result, "My Desc") {
		t.Errorf("description not found in result: %s", result)
	}

	if !strings.Contains(result, "kw1, kw2") {
		t.Errorf("keywords not found in result: %s", result)
	}
}

func TestFormatTemplate_ReplacesDate(t *testing.T) {
	tmpl := "<span>%date%</span>%CONTENT%"

	page := pages.Page{
		Content: []byte("text"),
		Meta: map[string]string{
			"date": "2025-01-01",
		},
	}

	result := pages.FormatTemplate(tmpl, page)

	if strings.Contains(result, "%date%") {
		t.Error("date placeholder was not replaced")
	}

	if !strings.Contains(result, "2025-01-01") {
		t.Errorf("date not found in result: %s", result)
	}
}

func TestFormatTemplate_EmptyMeta(t *testing.T) {
	tmpl := "<title>%title%</title>%CONTENT%"

	page := pages.Page{
		Content: []byte("hello"),
		Meta:    map[string]string{},
	}

	result := pages.FormatTemplate(tmpl, page)

	// placeholder stays if no replacement value — that's the current behavior
	if strings.Contains(result, "%CONTENT%") {
		t.Error("CONTENT should always be replaced")
	}
}

func TestFormatTemplate_WithRealTemplate(t *testing.T) {
	root := "./root"

	tmpl, err := pages.ReadTemplateFile(root)

	if err != nil {
		t.Fatal(err)
	}

	page := pages.Page{
		Content: []byte("# Article\n\nSome content here."),
		Meta: map[string]string{
			"title":       "Test Article | My Site",
			"description": "Test description",
			"keywords":    "test, article",
			"date":        "2025-06-01",
		},
	}

	result := pages.FormatTemplate(tmpl, page)

	if strings.Contains(result, "%CONTENT%") {
		t.Error("CONTENT placeholder not replaced in real template")
	}

	if strings.Contains(result, "%title%") {
		t.Error("title placeholder not replaced in real template")
	}

	if !strings.Contains(result, "Test Article | My Site") {
		t.Errorf("title not found in result")
	}

	if !strings.Contains(result, "<h1") {
		t.Errorf("expected heading in rendered content")
	}
}
