package main

import (
	"fmt"
	"os"
	"ascii-art-justify/functions"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nExample: go run . --align=right something standard")
		return
	}

	var inputText string
	bannerType := "standard"
	alignType := "left"

	if len(os.Args) == 2 {
		inputText = os.Args[1]
	} else if len(os.Args) == 3 {
		if strings.HasPrefix(os.Args[1], "--align=") {
			alignType = strings.TrimPrefix(os.Args[1], "--align=")
			inputText = os.Args[2]
		} else {
			inputText = os.Args[1]
			bannerType = os.Args[2]
		}
	} else {
		if strings.HasPrefix(os.Args[1], "--align=") {
			alignType = strings.TrimPrefix(os.Args[1], "--align=")
			inputText = os.Args[2]
			bannerType = os.Args[3]
		} else {
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nExample: go run . --align=right something standard")
			return
		}
	}
	
	if inputText == "" {
		return
	}

	termWidth := functions.TerminalWidth()
	art := functions.AsciiArt(inputText, bannerType)

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
		//1.
		// rawWords := strings.Split(inputText, " ")
		// words := []string{}
		// for _, w := range rawWords {
		// 	if w != "" {
		// 		words = append(words, w)
		// 	}
		// }

		//FIX
		words := strings.Fields(inputText)

		if len(words) == 0 {
			return
		}

		if len(words) == 1 {
			art := functions.AsciiArt(words[0], bannerType)
			fmt.Print(art)
			return
		}

		totalWordWidth := 0
		for _, word := range words {
			art := functions.AsciiArt(word, bannerType)
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
				art := functions.AsciiArt(word, bannerType)
				rows := strings.Split(art, "\n")
				fmt.Print(rows[row])
				if i < gaps {
					fmt.Print(strings.Repeat(" ", spacePerGap))
					if currentExtra > 0 {
						fmt.Print(" ")
						currentExtra--
					}
				}
			}
			fmt.Println()
		}
		return
	}

	pad := strings.Repeat(" ", spacesNeeded)

	//2.
	// for _, line := range artChar {
	// 	if line != "" {
	// 		fmt.Print(pad, line, "\n")
	// 	}
	// }

	//FIX
	for i, line := range artChar {
		if line == "" {
			if i != len(artChar) - 1 {
				fmt.Println()
			}
		}else{
			fmt.Print(pad, line, "\n")
		}
	}

}

