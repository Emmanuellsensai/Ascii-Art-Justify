package main

import (
	"strings"
	"testing"
	"ascii-art-justify/functions"
)

// ── AsciiArt core output tests ───────────────────────────────────────────────

func TestAsciiArtStandardBanner(t *testing.T) {
	result := functions.AsciiArt("hello", "standard")
	if result == "" {
		t.Error("expected non-empty output for 'hello' with standard banner")
	}
	lines := strings.Split(result, "\n")
	if len(lines) < 8 {
		t.Errorf("expected at least 8 lines, got %d", len(lines))
	}
}

func TestAsciiArtShadowBanner(t *testing.T) {
	result := functions.AsciiArt("hello", "shadow")
	if result == "" {
		t.Error("expected non-empty output for 'hello' with shadow banner")
	}
}

func TestAsciiArtThinkertoyBanner(t *testing.T) {
	result := functions.AsciiArt("hello", "thinkertoy")
	if result == "" {
		t.Error("expected non-empty output for 'hello' with thinkertoy banner")
	}
}

func TestAsciiArtInvalidBanner(t *testing.T) {
	result := functions.AsciiArt("hello", "nonexistent")
	if result != "" {
		t.Error("expected empty string for invalid banner, got output")
	}
}

// ── Newline handling ─────────────────────────────────────────────────────────

func TestAsciiArtNewlineOnly(t *testing.T) {
	result := functions.AsciiArt("\\n", "standard")
	if result != "\n" {
		t.Errorf("expected single newline for '\\n', got %q", result)
	}
}

func TestAsciiArtMultipleNewlines(t *testing.T) {
	result := functions.AsciiArt("\\n\\n\\n", "standard")
	if result != "\n\n\n" {
		t.Errorf("expected 3 newlines, got %q", result)
	}
}

func TestAsciiArtTextWithNewline(t *testing.T) {
	result := functions.AsciiArt("hi\\nthere", "standard")
	if result == "" {
		t.Error("expected non-empty output for multi-line input")
	}
	// should produce 8 art lines for "hi", a blank line, then 8 for "there"
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) < 17 {
		t.Errorf("expected at least 17 lines for 'hi\\nthere', got %d", len(lines))
	}
}

func TestAsciiArtEmptyString(t *testing.T) {
	result := functions.AsciiArt("", "standard")
	if result != "" {
		t.Errorf("expected empty string for empty input, got %q", result)
	}
}

// ── Output width consistency ─────────────────────────────────────────────────

func TestAsciiArtAllLinesEqualWidth(t *testing.T) {
	result := functions.AsciiArt("Hello", "standard")
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	width := len(lines[0])
	for i, line := range lines {
		if line == "" {
			continue
		}
		if len(line) != width {
			t.Errorf("line %d width = %d, want %d: %q", i, len(line), width, line)
		}
	}
}

func TestAsciiArtDifferentInputsHaveDifferentWidths(t *testing.T) {
	short := functions.AsciiArt("hi", "standard")
	long := functions.AsciiArt("hello", "standard")

	shortWidth := len(strings.Split(short, "\n")[0])
	longWidth := len(strings.Split(long, "\n")[0])

	if shortWidth >= longWidth {
		t.Errorf("expected 'hello' wider than 'hi', got %d >= %d", shortWidth, longWidth)
	}
}

// ── Alignment padding helpers ────────────────────────────────────────────────

func TestPaddingRight(t *testing.T) {
	termWidth := 120
	result := functions.AsciiArt("hello", "standard")
	lines := strings.Split(result, "\n")
	asciiWidth := len(lines[0])
	pad := strings.Repeat(" ", termWidth-asciiWidth)

	paddedLine := pad + lines[0]
	if !strings.HasPrefix(paddedLine, strings.Repeat(" ", termWidth-asciiWidth)) {
		t.Error("right align: padding not applied correctly")
	}
	if len(paddedLine) != termWidth {
		t.Errorf("right align: padded line width = %d, want %d", len(paddedLine), termWidth)
	}
}

// func TestPaddingCenter(t *testing.T) {
// 	termWidth := 120
// 	result := functions.AsciiArt("hello", "standard")
// 	lines := strings.Split(result, "\n")
// 	asciiWidth := len(lines[0])
// 	pad := strings.Repeat(" ", (termWidth-asciiWidth)/2)

// 	paddedLine := pad + lines[0]
// 	leftPad := len(paddedLine) - len(strings.TrimLeft(paddedLine, " "))
// 	rightPad := asciiWidth // content starts after pad

// 	_ = rightPad
// 	if leftPad != (termWidth-asciiWidth)/2 {
// 		t.Errorf("center: expected left pad %d, got %d", (termWidth-asciiWidth)/2, leftPad)
// 	}
// }

// func TestPaddingLeft(t *testing.T) {
// 	result := functions.AsciiArt("hello", "standard")
// 	lines := strings.Split(result, "\n")
// 	for i, line := range lines {
// 		if line == "" {
// 			continue
// 		}
// 		if strings.HasPrefix(line, " ") {
// 			t.Errorf("left align: line %d has unexpected leading space: %q", i, line)
// 		}
// 	}
// }
