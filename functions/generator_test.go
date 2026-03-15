package functions

import (
	"strings"
	"testing"
)

const testTermWidth = 80

// left output should be identical to raw art — no padding added
func TestAlignArt_Left(t *testing.T) {
	result := AlignArt("hello", "standard", "left", testTermWidth)
	raw := AsciiArt("hello", "standard")

	resultRows := strings.Split(result, "\n")
	rawRows := strings.Split(raw, "\n")

	for i := range rawRows {
		if i >= len(resultRows) {
			break
		}
		if resultRows[i] != rawRows[i] {
			t.Errorf("left: row %d expected %q but got %q", i, rawRows[i], resultRows[i])
		}
	}
}

// every row should be exactly termWidth wide
func TestAlignArt_Right(t *testing.T) {
	result := AlignArt("hello", "standard", "right", testTermWidth)

	for i, row := range strings.Split(result, "\n") {
		if row == "" {
			continue
		}
		if len(row) != testTermWidth {
			t.Errorf("right: row %d expected length %d but got %d", i, testTermWidth, len(row))
		}
	}
}

// center output should have more leading spaces than left but less than right
func TestAlignArt_Center(t *testing.T) {
	leftResult := AlignArt("hello", "standard", "left", testTermWidth)
	centerResult := AlignArt("hello", "standard", "center", testTermWidth)
	rightResult := AlignArt("hello", "standard", "right", testTermWidth)

	leftPad := len(strings.Split(leftResult, "\n")[0]) - len(strings.TrimLeft(strings.Split(leftResult, "\n")[0], " "))
	centerPad := len(strings.Split(centerResult, "\n")[0]) - len(strings.TrimLeft(strings.Split(centerResult, "\n")[0], " "))
	rightPad := len(strings.Split(rightResult, "\n")[0]) - len(strings.TrimLeft(strings.Split(rightResult, "\n")[0], " "))

	if centerPad <= leftPad {
		t.Errorf("center: padding (%d) should be greater than left padding (%d)", centerPad, leftPad)
	}
	if centerPad >= rightPad {
		t.Errorf("center: padding (%d) should be less than right padding (%d)", centerPad, rightPad)
	}
}

// every row should be exactly termWidth wide
func TestAlignArt_Justify(t *testing.T) {
	result := AlignArt("how are you", "standard", "justify", testTermWidth)

	for i, row := range strings.Split(result, "\n") {
		if row == "" {
			continue
		}
		if len(row) != testTermWidth {
			t.Errorf("justify: row %d expected width %d but got %d", i, testTermWidth, len(row))
		}
	}
}

// no alignment should return empty string for valid input
func TestAlignArt_NoneReturnEmpty(t *testing.T) {
	aligns := []string{"left", "right", "center", "justify"}

	for _, align := range aligns {
		t.Run(align, func(t *testing.T) {
			result := AlignArt("hello world", "standard", align, testTermWidth)
			if result == "" {
				t.Errorf("%s: expected output but got empty string", align)
			}
		})
	}
}