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

	result := functions.AlignArt(inputText, bannerType, alignType, termWidth)
	fmt.Print(result)
}
	