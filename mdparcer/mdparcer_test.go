package mdparcer

import (
	"bytes"
	"strings"
	"testing"
)

// TestParseCodeBlocks_NoCodeBlocks verifies that markdown with no fenced code
// blocks is returned unchanged and produces an empty blocks slice.
func TestParseCodeBlocks_NoCodeBlocks(t *testing.T) {
	input := []byte("# Heading\n\nSome paragraph text.\n\nAnother paragraph.")

	blocks, out := ParseCodeBlocks(input)

	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(blocks))
	}

	if !bytes.Equal(input, out) {
		t.Errorf("expected output to equal input\ngot:  %q\nwant: %q", out, input)
	}
}

// TestParseCodeBlocks_SingleBlock verifies that a single fenced code block is
// extracted correctly. The closing fence is at line index 3, so the Id must
// be "CODE_BLOCK_3".
func TestParseCodeBlocks_SingleBlock(t *testing.T) {
	input := []byte("Some text\n```go\nfmt.Println(\"hello\")\n```\nAfter")

	blocks, out := ParseCodeBlocks(input)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	b := blocks[0]

	if b.Lang != "go" {
		t.Errorf("expected Lang=%q, got %q", "go", b.Lang)
	}

	if b.Code != "fmt.Println(\"hello\")" {
		t.Errorf("expected Code=%q, got %q", "fmt.Println(\"hello\")", b.Code)
	}

	if b.Id != "CODE_BLOCK_3" {
		t.Errorf("expected Id=%q, got %q", "CODE_BLOCK_3", b.Id)
	}

	outStr := string(out)

	if !strings.Contains(outStr, "CODE_BLOCK_3") {
		t.Errorf("expected modified md to contain placeholder %q, got: %q", "CODE_BLOCK_3", outStr)
	}

	if strings.Contains(outStr, "```") {
		t.Errorf("expected modified md to not contain backtick fences, got: %q", outStr)
	}

	if strings.Contains(outStr, "fmt.Println") {
		t.Errorf("expected code content to be replaced by placeholder, got: %q", outStr)
	}
}

// TestParseCodeBlocks_MultipleBlocks verifies that two separate fenced blocks
// are both extracted.
func TestParseCodeBlocks_MultipleBlocks(t *testing.T) {
	input := []byte("Intro\n```python\nprint('hi')\n```\nMiddle\n```bash\necho hello\n```\nEnd")

	blocks, out := ParseCodeBlocks(input)

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	if blocks[0].Lang != "python" {
		t.Errorf("expected first block Lang=%q, got %q", "python", blocks[0].Lang)
	}

	if blocks[0].Code != "print('hi')" {
		t.Errorf("expected first block Code=%q, got %q", "print('hi')", blocks[0].Code)
	}

	if blocks[1].Lang != "bash" {
		t.Errorf("expected second block Lang=%q, got %q", "bash", blocks[1].Lang)
	}

	if blocks[1].Code != "echo hello" {
		t.Errorf("expected second block Code=%q, got %q", "echo hello", blocks[1].Code)
	}

	outStr := string(out)

	if !strings.Contains(outStr, blocks[0].Id) {
		t.Errorf("expected output to contain placeholder %q", blocks[0].Id)
	}

	if !strings.Contains(outStr, blocks[1].Id) {
		t.Errorf("expected output to contain placeholder %q", blocks[1].Id)
	}
}

// TestParseCodeBlocks_UnknownLang verifies that an unknown language tag is
// preserved exactly as supplied.
func TestParseCodeBlocks_UnknownLang(t *testing.T) {
	input := []byte("```xyz123\nsome code\n```")

	blocks, _ := ParseCodeBlocks(input)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].Lang != "xyz123" {
		t.Errorf("expected Lang=%q, got %q", "xyz123", blocks[0].Lang)
	}
}

// TestParseCodeBlocks_NoLang verifies that a fence opened with just ``` (no
// language) produces a block with an empty Lang field.
func TestParseCodeBlocks_NoLang(t *testing.T) {
	// Opening fence is exactly "```" which has len==3 (>2) so it opens a block.
	input := []byte("``` \nsome code\n```")

	blocks, _ := ParseCodeBlocks(input)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	// strings.Replace("``` ", "```", "", 1) == " "
	// The Lang field holds whatever follows the opening backticks.
	if blocks[0].Lang != " " {
		// A fence opened with just "```" (len 3, no trailing chars) would give
		// Lang="". A fence with a trailing space gives Lang=" ". Accept both
		// empty-ish values to be robust to whitespace variation.
		if strings.TrimSpace(blocks[0].Lang) != "" {
			t.Errorf("expected Lang to be empty or whitespace, got %q", blocks[0].Lang)
		}
	}
}

// TestParseCodeBlocks_NoLangExact verifies a fence opened with bare ``` (no
// trailing characters) yields Lang="".
func TestParseCodeBlocks_NoLangExact(t *testing.T) {
	// "```\ncode\n```" — the opening "```" has len==3, which is > 2, so it
	// opens a block. strings.Replace("```","```","",1) == "".
	//
	// However: the closing fence check (line == "```") runs BEFORE the opening
	// fence check, but only when blockLines != nil. On the first "```" line,
	// blockLines IS nil, so it falls through to the opening fence check and
	// opens a block correctly.
	input := []byte("```\ncode line\n```")

	blocks, _ := ParseCodeBlocks(input)

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].Lang != "" {
		t.Errorf("expected Lang=%q, got %q", "", blocks[0].Lang)
	}

	if blocks[0].Code != "code line" {
		t.Errorf("expected Code=%q, got %q", "code line", blocks[0].Code)
	}
}

// TestFormatCode_GoLang verifies that a valid Go code block is formatted into
// non-empty HTML containing syntax-highlight spans.
func TestFormatCode_GoLang(t *testing.T) {
	block := CodeBlock{
		Lang: "go",
		Code: `package main

import "fmt"

func main() {
	fmt.Println("hello")
}`,
		Id: "CODE_BLOCK_0",
	}

	result, err := FormatCode(block)

	if err != nil {
		t.Fatalf("FormatCode returned unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}

	if !strings.Contains(string(result), "<span") {
		t.Errorf("expected HTML output to contain <span, got: %q", string(result))
	}
}

// TestFormatCode_UnknownLang verifies that an unrecognised language falls back
// to the chroma Fallback lexer without returning an error.
func TestFormatCode_UnknownLang(t *testing.T) {
	block := CodeBlock{
		Lang: "xyz123_notexist",
		Code: "some code here",
		Id:   "CODE_BLOCK_0",
	}

	result, err := FormatCode(block)

	if err != nil {
		t.Fatalf("FormatCode returned unexpected error for unknown lang: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("expected non-empty result even for unknown language")
	}
}

// TestFormatCode_EmptyCode verifies that an empty Code string does not cause
// an error.
func TestFormatCode_EmptyCode(t *testing.T) {
	block := CodeBlock{
		Lang: "go",
		Code: "",
		Id:   "CODE_BLOCK_0",
	}

	_, err := FormatCode(block)

	if err != nil {
		t.Fatalf("FormatCode returned unexpected error for empty code: %v", err)
	}
}

// TestRenderCode_ReplacesPlaceholders verifies that RenderCode replaces each
// block's Id placeholder within the HTML with the formatted code HTML.
func TestRenderCode_ReplacesPlaceholders(t *testing.T) {
	block := CodeBlock{
		Lang: "go",
		Code: `fmt.Println("hi")`,
		Id:   "CODE_BLOCK_42",
	}

	// The rendered markdown would wrap the placeholder in a paragraph.
	input := []byte("<p>CODE_BLOCK_42</p>")

	result := RenderCode(input, []CodeBlock{block})

	resultStr := string(result)

	if strings.Contains(resultStr, "CODE_BLOCK_42") {
		t.Errorf("expected placeholder to be replaced, but it still appears in output: %q", resultStr)
	}

	if !strings.Contains(resultStr, "<span") {
		t.Errorf("expected formatted code HTML containing <span in output: %q", resultStr)
	}
}

// TestRenderCode_EmptyBlocks verifies that passing an empty blocks slice
// returns the input HTML unchanged.
func TestRenderCode_EmptyBlocks(t *testing.T) {
	input := []byte("<p>Hello world</p>")

	result := RenderCode(input, []CodeBlock{})

	if !bytes.Equal(input, result) {
		t.Errorf("expected output to equal input\ngot:  %q\nwant: %q", result, input)
	}
}

// TestMdToHTML_BasicMarkdown verifies that standard markdown headings and
// paragraphs are rendered to HTML.
func TestMdToHTML_BasicMarkdown(t *testing.T) {
	input := []byte("# Hello\n\nWorld")

	result := MdToHTML(input)

	resultStr := string(result)

	if !strings.Contains(resultStr, "<h1") {
		t.Errorf("expected output to contain <h1, got: %q", resultStr)
	}

	if !strings.Contains(resultStr, "World") {
		t.Errorf("expected output to contain 'World', got: %q", resultStr)
	}
}

// TestMdToHTML_WithCodeBlock verifies that a fenced code block in markdown is
// syntax-highlighted (contains <span) and the raw ``` fences are absent.
func TestMdToHTML_WithCodeBlock(t *testing.T) {
	input := []byte("# Title\n\n```go\nfmt.Println(\"hello\")\n```\n")

	result := MdToHTML(input)

	resultStr := string(result)

	if !strings.Contains(resultStr, "<span") {
		t.Errorf("expected syntax-highlighted HTML to contain <span, got: %q", resultStr)
	}

	if strings.Contains(resultStr, "```") {
		t.Errorf("expected output to not contain raw ``` fences, got: %q", resultStr)
	}
}

// TestMdToHTML_EmptyInput verifies that an empty byte slice does not panic and
// returns a byte slice (which may be empty or contain only whitespace).
func TestMdToHTML_EmptyInput(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MdToHTML panicked on empty input: %v", r)
		}
	}()

	result := MdToHTML([]byte{})

	if result == nil {
		t.Error("expected non-nil result for empty input")
	}
}
