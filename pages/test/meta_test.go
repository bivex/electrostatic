package pages_test

import (
	"testing"

	"github.com/laranatech/electrostatic/pages"
)

func TestReadMetaConfig(t *testing.T) {
	root := "./root"

	config, err := pages.ReadMetaConfig(root)

	if err != nil {
		t.Fatal(err)
	}

	if config.TitleTemplate != "%title% | Your site" {
		t.Errorf("unexpected TitleTemplate: %q", config.TitleTemplate)
	}

	if config.DescriptionTemplate != "%description%" {
		t.Errorf("unexpected DescriptionTemplate: %q", config.DescriptionTemplate)
	}

	if config.KeywordsTemplate != "%keywords%, your keywords" {
		t.Errorf("unexpected KeywordsTemplate: %q", config.KeywordsTemplate)
	}

	if config.FallbackTitle != "Fallback title" {
		t.Errorf("unexpected FallbackTitle: %q", config.FallbackTitle)
	}

	if config.FallbackDescription != "Fallback description" {
		t.Errorf("unexpected FallbackDescription: %q", config.FallbackDescription)
	}

	if config.FallbackKeywords != "Fallback keywords" {
		t.Errorf("unexpected FallbackKeywords: %q", config.FallbackKeywords)
	}
}

func TestReadMetaConfig_MissingFile(t *testing.T) {
	_, err := pages.ReadMetaConfig("./nonexistent")

	if err == nil {
		t.Error("expected error for missing meta.json, got nil")
	}
}

func TestNewMetaMap_Fallbacks(t *testing.T) {
	root := "./root"

	meta, err := pages.NewMetaMap(root, map[string]string{})

	if err != nil {
		t.Fatal(err)
	}

	if meta["title"] != "Fallback title | Your site" {
		t.Errorf("unexpected fallback title: %q", meta["title"])
	}

	if meta["description"] != "Fallback description" {
		t.Errorf("unexpected fallback description: %q", meta["description"])
	}

	if meta["keywords"] != "Fallback keywords, your keywords" {
		t.Errorf("unexpected fallback keywords: %q", meta["keywords"])
	}

	if meta["date"] != "" {
		t.Errorf("expected empty date, got: %q", meta["date"])
	}
}

func TestNewMetaMap_WithParams(t *testing.T) {
	root := "./root"

	params := map[string]string{
		"title":       "What is larana",
		"description": "some description",
		"keywords":    "larana, gorana, framework",
		"date":        "2025-01-01",
	}

	meta, err := pages.NewMetaMap(root, params)

	if err != nil {
		t.Fatal(err)
	}

	if meta["title"] != "What is larana | Your site" {
		t.Errorf("unexpected title: %q", meta["title"])
	}

	if meta["description"] != "some description" {
		t.Errorf("unexpected description: %q", meta["description"])
	}

	if meta["keywords"] != "larana, gorana, framework, your keywords" {
		t.Errorf("unexpected keywords: %q", meta["keywords"])
	}

	if meta["date"] != "2025-01-01" {
		t.Errorf("unexpected date: %q", meta["date"])
	}
}

func TestNewMetaMap_PartialParams(t *testing.T) {
	root := "./root"

	params := map[string]string{
		"title": "Custom title",
	}

	meta, err := pages.NewMetaMap(root, params)

	if err != nil {
		t.Fatal(err)
	}

	if meta["title"] != "Custom title | Your site" {
		t.Errorf("unexpected title: %q", meta["title"])
	}

	// description and keywords should fall back to defaults
	if meta["description"] != "Fallback description" {
		t.Errorf("unexpected description: %q", meta["description"])
	}

	if meta["keywords"] != "Fallback keywords, your keywords" {
		t.Errorf("unexpected keywords: %q", meta["keywords"])
	}
}

func TestNewMetaMap_MissingRoot(t *testing.T) {
	_, err := pages.NewMetaMap("./nonexistent", map[string]string{})

	if err == nil {
		t.Error("expected error for missing root, got nil")
	}
}
