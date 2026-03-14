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
		return ""
	}

	asciiChar := strings.Split(string(bannerFile), "\n")
	words := strings.Split(inputText, "\\n")
	result := ""

	if strings.ReplaceAll(inputText, "\\n", "") == "" {
		return strings.Repeat("\n", len(words)-1)
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

	return result
}

func TerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Print("error: ", err)
		return 0
	}
	cols := strings.TrimSpace(string(out))
	output, err := strconv.Atoi(cols)
	if err != nil {
		fmt.Print("error: ", err)
		return 0
	}
	return output
}
