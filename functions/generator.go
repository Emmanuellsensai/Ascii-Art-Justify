package functions

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func AsciiArt(inputText, bannerType string) string {
	bannerFile, err := os.ReadFile("banners/" + bannerType + ".txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return  ""
	}

	asciiChar := strings.Split(string(bannerFile), "\n")
	words := strings.Split(inputText, "\\n")
	result := ""

	if strings.ReplaceAll(inputText, "\\n", "") == "" {
		return  strings.Repeat("\n", len(words)-1)
	}

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

	return  result
}

func TerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Print("error: ", err)
		return  0
	}
	cols := strings.TrimSpace(string(out))
	output, err := strconv.Atoi(cols)
	if err != nil {
		fmt.Print("error: ", err)
		return  0
	}
	return  output
}

func AlignArt(inputText, bannerType, alignType string, termWidth int) string {

	var result strings.Builder

	art := AsciiArt(inputText, bannerType)
	artChar := strings.Split(art, "\n")
	asciiSize := len(artChar[0])

	var spacesNeeded int
	var spacePerGap int
	var extraSpace int

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

		gaps := len(words) - 1
		totalSpaces := termWidth - totalWordWidth

		spacePerGap = totalSpaces / gaps
		extraSpace = totalSpaces % gaps

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

	for i, line := range artChar {
		if line == "" {
			if i != len(artChar) - 1 {
				result.WriteString("\n")
			}
		}else{
			result.WriteString(pad + line + "\n")
		}
	}
	return  result.String()
}

