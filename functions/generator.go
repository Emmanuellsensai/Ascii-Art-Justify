package functions

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func AsciiArt(inputText, bannerType string) string {
	// Guard clause for invalid banner type.
	if bannerType != "standard" && bannerType != "shadow" && bannerType != "thinkertoy" {
		fmt.Println("Error: Banner must be shadow, standard, or thinkertoy.")
		os.Exit(0)
	}
	// Read and save the content of the bannerfile.
	bannerFile, err := os.ReadFile("banners/" + bannerType + ".txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	// Saving the ascii characters and words.
	asciiChar := strings.Split(string(bannerFile), "\n")
	words := strings.Split(inputText, "\\n")
	result := ""
	// Error handling for string with  only new line.
	if strings.ReplaceAll(inputText, "\\n", "") == "" {
		return strings.Repeat("\n", len(words)-1)
	}
	// Gennerate ascii art for the words.
	for _, word := range words {
		if word == "" {
			result += "\n"
			continue
		}

		for i := 0; i < 8; i++ {
			for _, char := range word {
				if char >= ' ' && char <= '~' {
					result += asciiChar[i+(int(char-' ')*9)+1]
				}
			}
			result += "\n"
		}
	}

	return result
}

func TerminalWidth() int {
	// Use tput colls to get terminal width.
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Print("error: ", err)
		return 0
	}
	// Get the integar value of the terminal width.
	cols := strings.TrimSpace(string(out))
	output, err := strconv.Atoi(cols)
	if err != nil {
		fmt.Print("error: ", err)
		return 0
	}
	return output
}

func AlignArt(inputText, bannerType, alignType string, termWidth int) string {
	var result strings.Builder

	art := AsciiArt(inputText, bannerType)
	artChar := strings.Split(art, "\n")
	asciiSize := len(artChar[0])

	var spacesNeeded int
	var spacePerGap int
	var extraSpace int
	// Guard clause for aligntype.
	if alignType != "right" && alignType != "center" && alignType != "left" && alignType != "justify" {
		fmt.Println("Error : Option should be right, center, left or justify ")
		os.Exit(0)
	}
	// Switch statement for alighntypes
	switch alignType {
	case "right":
		spacesNeeded = termWidth - asciiSize
	case "center":
		spacesNeeded = (termWidth - asciiSize) / 2
	case "left":
		spacesNeeded = 0
	case "justify":
		words := strings.Fields(inputText)

		if len(words) == 0 {
			return result.String()
		}
		// Returning the default ascii art if there is only 1 word.
		if len(words) == 1 {
			art := AsciiArt(words[0], bannerType)
			result.WriteString(art)
			return result.String()
		}

		totalWordWidth := 0
		for _, word := range words {
			art := AsciiArt(word, bannerType)
			rows := strings.Split(art, "\n")
			totalWordWidth += len(rows[0])
		}
		// Calculate the gaps needed between the justified words.
		gaps := len(words) - 1
		totalSpaces := termWidth - totalWordWidth

		if totalSpaces < 0 {
 		   totalSpaces = 0
		}

		spacePerGap = totalSpaces / gaps
		extraSpace = totalSpaces % gaps
		// Print the rows of the ascii art with the space between each word row.
		for row := 0; row < 8; row++ {
			currentExtra := extraSpace
			for i, word := range words {
				art := AsciiArt(word, bannerType)
				rows := strings.Split(art, "\n")
				result.WriteString(rows[row])
				if i < gaps {
					result.WriteString(strings.Repeat(" ", spacePerGap))
					if currentExtra > 0 {
						result.WriteString(" ")
						currentExtra--
					}
				}
			}
			result.WriteString("\n")
		}
		return result.String()
	}

	pad := strings.Repeat(" ", spacesNeeded)
	// Check for empty string between words and replaces with newline, otherwise prints the pad, word and newline at the end.
	for i, line := range artChar {
		if line == "" {
			if i != len(artChar)-1 {
				result.WriteString("\n")
			}
		} else {
			result.WriteString(pad + line + "\n")
		}
	}
	return result.String()
}
