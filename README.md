# ascii-art-justify

A command-line tool written in Go that renders text as ASCII art and lets you control where the art appears on the terminal — left, right, center, or justified — using the `--align` flag.

---

## Table of Contents

1. [What Does This Program Do?](#1-what-does-this-program-do)
2. [How to Run It](#2-how-to-run-it)
3. [Project Structure](#3-project-structure)
4. [How the Font Files Work](#4-how-the-font-files-work)
5. [How Each Alignment Works](#5-how-each-alignment-works)
6. [Understanding functions/generator.go — Line by Line](#6-understanding-functionsgeneratorgo--line-by-line)
   - [Package and Imports](#61-package-and-imports)
   - [AsciiArt](#62-asciiart)
   - [TerminalWidth](#63-terminalwidth)
7. [Understanding main.go — Line by Line](#7-understanding-maingo--line-by-line)
   - [Package and Imports](#71-package-and-imports)
   - [Argument Count Check](#72-argument-count-check)
   - [Parsing Arguments](#73-parsing-arguments)
   - [Getting Terminal Width and Art](#74-getting-terminal-width-and-art)
   - [The Alignment Switch](#75-the-alignment-switch)
   - [Printing the Result](#76-printing-the-result)
8. [The Two Bug Fixes](#8-the-two-bug-fixes)
9. [Example Runs](#9-example-runs)
10. [Common Mistakes](#10-common-mistakes)

---

## 1. What Does This Program Do?

This program takes text you give it and renders it as large ASCII art in the terminal. On top of that you can control **where on the terminal** the art appears using the `--align` flag.

Without the flag the art prints at the left as normal. With the flag you can push it to the right edge, center it, or spread words evenly across the full terminal width.

---

## 2. How to Run It

```bash
# No flag — defaults to left
go run . "hello"
go run . "hello" shadow

# Left — same as default
go run . --align=left "hello" standard

# Right — pushed to the right edge of the terminal
go run . --align=right "hello" standard

# Center — middle of the terminal
go run . --align=center "hello" shadow

# Justify — words spread evenly across the full terminal width
go run . --align=justify "how are you" shadow

# Multi-line art
go run . "Hello\nWorld" standard

# Wrong flag format — prints usage message
go run . --align right something standard
```

**Available alignments:** `left` `right` `center` `justify`

**Available fonts:** `standard` `shadow` `thinkertoy`

---

## 3. Project Structure

```
ascii-art-justify/
├── main.go              ← Entry point: reads arguments, validates, handles alignment
├── functions/
│   └── generator.go     ← Core engine: builds ASCII art, gets terminal width
└── banners/
    ├── standard.txt     ← Font file: standard ASCII art style
    ├── shadow.txt       ← Font file: shadow ASCII art style
    └── thinkertoy.txt   ← Font file: thinkertoy ASCII art style
```

> **Important:** The `.txt` font files live inside the `banners/` folder. The program looks for them at `banners/<name>.txt` relative to where you run `go run .`

---

## 4. How the Font Files Work

Every font file stores all printable characters from space (ASCII 32) to `~` (ASCII 126) in order. Each character is exactly **8 lines tall** followed by **1 blank separator line** — making **9 lines per character** in the file.

```
              ← one blank line at the very top of the file
              ← 8 lines for SPACE (all blank)
...
              ← blank separator
 _            ← 8 lines for !
| |
| |
| |
|_|
(_)

              ← blank separator
...continues for every character up to ~
```

The program uses this structure to calculate exactly where in the file any character's art lives using a simple formula — more on this in the `AsciiArt` section.

---

## 5. How Each Alignment Works

Every alignment is just **adding spaces before each row of art**. The only question is how many.

```
left:    [art starts here]
right:   [          spaces          ][art]
center:  [     spaces     ][art]
justify: [word][   spaces   ][word][   spaces   ][word]
```

To calculate spaces you need two numbers:

- **termWidth** — how wide the terminal currently is (e.g. 80 or 120 columns)
- **asciiSize** — how wide the rendered art is (length of the first row of output)

| Alignment | Formula                                           |
| --------- | ------------------------------------------------- |
| `left`    | `0` — no padding                                  |
| `right`   | `termWidth - asciiSize`                           |
| `center`  | `(termWidth - asciiSize) / 2`                     |
| `justify` | spaces split _between_ words, not before the line |

---

## 6. Understanding functions/generator.go — Line by Line

---

### 6.1 Package and Imports

```go
package functions
```

This file belongs to the `functions` package. It cannot run on its own — `main.go` imports it. Because it is a separate package, only names starting with a **capital letter** are visible to `main.go`. That is why `AsciiArt` and `TerminalWidth` are capitalised.

```go
import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)
```

| Package     | What it gives us                                                                         |
| ----------- | ---------------------------------------------------------------------------------------- |
| `"fmt"`     | `fmt.Println` — printing error messages                                                  |
| `"os"`      | `os.ReadFile` — reads the font file from disk, `os.Stdin` — connects terminal to command |
| `"os/exec"` | `exec.Command` — runs a terminal command from inside Go                                  |
| `"strconv"` | `strconv.Atoi` — converts a string like `"80"` into the integer `80`                     |
| `"strings"` | text tools — splitting, trimming, replacing                                              |

---

### 6.2 AsciiArt

This function takes text and a font name and returns the complete ASCII art as a string.

```go
func AsciiArt(inputText, bannerType string) string {
```

Two inputs, one output:

- `inputText` — the text to render, e.g. `"hello"`
- `bannerType` — the font to use, e.g. `"standard"`

Returns a `string` — the complete ASCII art ready to be printed or measured.

---

```go
    bannerFile, err := os.ReadFile("banners/" + bannerType + ".txt")
    if err != nil {
        fmt.Println("Error: ", err)
        return ""
    }
```

`os.ReadFile` opens and reads the entire font file at once. `"banners/" + bannerType + ".txt"` builds the path — so `"standard"` becomes `"banners/standard.txt"`.

If the file cannot be found or read, `err` will not be `nil`. We print the error and return an empty string — the program will not crash but nothing will be drawn.

---

```go
    asciiChar := strings.Split(string(bannerFile), "\n")
```

`string(bannerFile)` converts the raw file bytes into a readable Go string. `strings.Split(..., "\n")` cuts it at every newline, giving us a slice where each element is one line of the font file. So `asciiChar[0]` is the first line, `asciiChar[1]` is the second, and so on.

---

```go
    words := strings.Split(inputText, "\\n")
```

Splits the user's input on the literal two-character sequence backslash + `n`. When a user types `"Hello\nWorld"` in the terminal, Go receives a `\` and an `n` next to each other — not a real newline. `"\\n"` in Go source code represents exactly those two characters, so this correctly splits the input into separate art blocks.

`"Hello\nWorld"` → `["Hello", "World"]` — two separate art blocks.

---

```go
    result := ""
```

An empty string that we will build up piece by piece. Every row of every character gets added to this.

---

```go
    if strings.ReplaceAll(inputText, "\\n", "") == "" {
        return strings.Repeat("\n", len(words)-1)
    }
```

This handles the edge case where the input is **only** newlines — like `"\n"` or `"\n\n"`. After removing all the `\n` sequences, if nothing is left then the entire input is just newlines. In that case return the right number of blank lines and stop.

`strings.ReplaceAll(inputText, "\\n", "")` — removes every `\n` from the input. If the result is `""` (empty), the user typed only newlines.

`strings.Repeat("\n", len(words)-1)` — returns the correct number of newlines. For input `"\n"`, `words` is `["", ""]` so `len(words)-1 = 1` — one newline. For `"\n\n"`, `words` is `["", "", ""]` so `len(words)-1 = 2` — two newlines.

---

```go
    for _, word := range words {
        if word == "" {
            result += "\n"
            continue
        }
```

Loop through each segment of the input. If a segment is empty (from a `\n` in the middle of the input), add a blank line to the result and skip to the next segment with `continue`.

---

```go
        for i := 0; i < 8; i++ {
            for _, char := range word {
                if char >= ' ' && char <= '~' {
                    result += asciiChar[i+(int(char-' ')*9)+1]
                }
            }
            result += "\n"
        }
```

This is the core rendering logic — two nested loops.

**Outer loop `i := 0; i < 8`** — runs 8 times, once for each art row. Remember every character is 8 rows tall.

**Inner loop `for _, char := range word`** — goes through every character in the current word.

**`if char >= ' ' && char <= '~'`** — only process printable ASCII characters (space through tilde). This prevents crashes if the user passes unusual input.

**`asciiChar[i+(int(char-' ')*9)+1]`** — this is the key formula that finds the correct line in the font file for any character at any row. Let's break it down:

- `char - ' '` — finds the position of this character relative to space. Space is ASCII 32, `A` is 65, so `A - ' '` = `65 - 32` = `33`. This tells us `A` is the 33rd character in the font file.
- `* 9` — each character takes 9 lines (8 art rows + 1 blank separator), so multiply by 9 to find where this character's block starts in the file.
- `+ 1` — skip the extra blank line at the very top of the font file.
- `+ i` — add the current row number (0 through 7) to get the specific art row we want.

For example, to get row 2 of letter `A`:

```
i + (int('A' - ' ') * 9) + 1
= 2 + (33 * 9) + 1
= 2 + 297 + 1
= 300   ← line 300 in the font file is row 2 of A
```

After all characters in a row are added, `result += "\n"` moves to the next line.

---

```go
    return result
}
```

Return the completed ASCII art string to whoever called this function.

---

### 6.3 TerminalWidth

This function finds out how wide the terminal window currently is.

```go
func TerminalWidth() int {
```

Returns a single integer — the number of columns in the terminal.

---

```go
    cmd := exec.Command("tput", "cols")
    cmd.Stdin = os.Stdin
```

`tput cols` is a standard terminal command that prints the number of columns in the terminal. For example if your terminal is 80 characters wide, `tput cols` prints `80`.

`exec.Command("tput", "cols")` creates a command object — like writing down what you want to run without running it yet.

`cmd.Stdin = os.Stdin` connects our running terminal to the `tput` command. `tput` needs to read terminal information from standard input to know which terminal to measure. Without this line, `tput` cannot detect the terminal and the command fails.

---

```go
    out, err := cmd.Output()
    if err != nil {
        fmt.Print("error: ", err)
        return 0
    }
```

`.Output()` actually runs the command and captures everything it prints. `out` will contain something like `"80\n"`. If it fails, print the error and return `0`.

---

```go
    cols := strings.TrimSpace(string(out))
    output, err := strconv.Atoi(cols)
    if err != nil {
        fmt.Print("error: ", err)
        return 0
    }
    return output
}
```

`string(out)` converts raw bytes to a string: `"80\n"`.

`strings.TrimSpace(...)` removes the trailing newline: `"80"`.

`strconv.Atoi(cols)` converts the string `"80"` to the integer `80`. You cannot do maths with strings in Go — `Atoi` (ASCII to integer) is the conversion tool. If it fails for any reason, return `0`.

Finally return the width as an integer.

---

## 7. Understanding main.go — Line by Line

---

### 7.1 Package and Imports

```go
package main
```

Every runnable Go program needs exactly one `package main`. This is where Go starts executing.

```go
import (
    "fmt"
    "os"
    "ascii-art-justify/functions"
    "strings"
)
```

| Import                          | Purpose                                                                                        |
| ------------------------------- | ---------------------------------------------------------------------------------------------- |
| `"fmt"`                         | Printing to terminal                                                                           |
| `"os"`                          | `os.Args` — reading command-line arguments                                                     |
| `"ascii-art-justify/functions"` | Our `generator.go` functions (`AsciiArt`, `TerminalWidth`)                                     |
| `"strings"`                     | `strings.HasPrefix`, `strings.TrimPrefix`, `strings.Fields`, `strings.Repeat`, `strings.Split` |

The import path `"ascii-art-justify/functions"` combines the module name from `go.mod` with the subfolder name. After this, all exported functions from `generator.go` are accessed with the `functions.` prefix.

---

### 7.2 Argument Count Check

```go
if len(os.Args) < 2 || len(os.Args) > 4 {
    fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nExample: go run . --align=right something standard")
    return
}
```

`os.Args` is a slice of everything typed on the command line. Index 0 is always the program name itself, so user arguments start at index 1.

Valid argument counts are 2, 3, or 4:

```
go run . "hello"                          → len = 2
go run . --align=right "hello"            → len = 3
go run . "hello" standard                 → len = 3
go run . --align=right "hello" standard   → len = 4
```

If outside this range, print the usage message and stop. `\n` inside the `fmt.Println` string creates a blank line between the two lines of the usage message.

---

### 7.3 Parsing Arguments

```go
var inputText string
bannerType := "standard"
alignType := "left"
```

Declare the three variables with their default values. `var inputText string` gives it the zero value for strings which is `""`. `bannerType` defaults to `"standard"` and `alignType` defaults to `"left"` — so if the user provides no banner or no align flag, these sensible defaults are already in place.

---

```go
if len(os.Args) == 2 {
    inputText = os.Args[1]
```

Only one user argument — it must be the text. No flag, no banner.

---

```go
} else if len(os.Args) == 3 {
    if strings.HasPrefix(os.Args[1], "--align=") {
        alignType = strings.TrimPrefix(os.Args[1], "--align=")
        inputText = os.Args[2]
    } else {
        inputText = os.Args[1]
        bannerType = os.Args[2]
    }
```

Two user arguments — two possible shapes:

- If the first argument starts with `"--align="` → it is a flag + text: extract the align value, second arg is the text
- Otherwise → it is text + banner: first arg is text, second is the banner name

`strings.HasPrefix(os.Args[1], "--align=")` — returns `true` if `os.Args[1]` starts with `"--align="`.

`strings.TrimPrefix(os.Args[1], "--align=")` — removes `"--align="` from the front, leaving just the alignment value. So `"--align=right"` becomes `"right"`.

---

```go
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
```

Three user arguments — only one valid shape: `--align=<value>` then text then banner. If the first argument is not a flag, there is no valid interpretation for three arguments without a flag — print usage and stop.

---

```go
if inputText == "" {
    return
}
```

If after all the parsing `inputText` is still empty, there is nothing to render. Stop silently.

---

### 7.4 Getting Terminal Width and Art

```go
termWidth := functions.TerminalWidth()
art := functions.AsciiArt(inputText, bannerType)
```

Call `TerminalWidth()` to find out how wide the terminal is right now. Call `AsciiArt` to render the full art as a string.

---

```go
artChar := strings.Split(art, "\n")
asciiSize := len(artChar[0])
```

`strings.Split(art, "\n")` cuts the art string at every real newline, giving a slice where each element is one row of art.

`len(artChar[0])` measures the length of the very first row. Since all rows of the art are the same width, row 0 is a reliable way to get the total art width. This is stored as `asciiSize` and used in the padding calculations.

---

```go
var spacesNeeded int
var spacePerGap int
var extraSpace int
```

Declare three variables used in different alignment cases. `var` gives them the zero value for integers (`0`). `spacePerGap` and `extraSpace` are only used in the justify case but are declared here at the top.

---

### 7.5 The Alignment Switch

```go
switch alignType {
case "right":
    spacesNeeded = termWidth - asciiSize
case "center":
    spacesNeeded = (termWidth - asciiSize) / 2
case "left":
    spacesNeeded = 0
```

`switch` checks `alignType` and runs only the matching case.

**right:** push the art to the right edge. The padding = everything to the left of the art = `termWidth - asciiSize`.

**center:** put equal space on both sides. Divide the available space in half. Integer division automatically rounds down if the number is odd — meaning the art sits slightly left of perfect center, which is acceptable.

**left:** no padding needed. `spacesNeeded` stays `0`.

---

```go
case "justify":
    words := strings.Fields(inputText)
```

`strings.Fields` splits the input on any whitespace and returns only the actual words — no empty strings even if there are multiple spaces between words. `"how  are  you"` → `["how", "are", "you"]`.

This replaced the original approach of `strings.Split(inputText, " ")` followed by a manual loop to filter out empty strings — `strings.Fields` does both in one call.

---

```go
    if len(words) == 0 {
        return
    }

    if len(words) == 1 {
        art := functions.AsciiArt(words[0], bannerType)
        fmt.Print(art)
        return
    }
```

Two edge cases handled upfront:

- **Zero words** — nothing to render, stop.
- **One word** — justify with a single word has no gaps to distribute space into. Just print the word as-is (left aligned) and stop.

---

```go
    totalWordWidth := 0
    for _, word := range words {
        art := functions.AsciiArt(word, bannerType)
        rows := strings.Split(art, "\n")
        totalWordWidth += len(rows[0])
    }
```

Measure the total art width of all words combined — with no spaces between them. For each word, render its art, split into rows, measure row 0, and add it to the running total.

---

```go
    gaps := len(words) - 1
    totalSpaces := termWidth - totalWordWidth
    spacePerGap = totalSpaces / gaps
    extraSpace = totalSpaces % gaps
```

`gaps` — the number of spaces between words. Three words have two gaps.

`totalSpaces` — how much space is left over after placing all the word art.

`spacePerGap` — divide the total space evenly between all gaps. Integer division gives the base amount.

`extraSpace` — `totalSpaces % gaps` is the remainder after dividing. For example if `totalSpaces = 10` and `gaps = 3`, then `spacePerGap = 3` and `extraSpace = 1`. That one leftover space needs to go somewhere — it gets distributed one-by-one to the first gaps.

---

```go
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
```

Draw all 8 rows of the justified output.

`currentExtra := extraSpace` — reset the extra space counter at the start of every row so each row distributes the remainder the same way.

For each word at the current row: print the word's art row, then after every word except the last (`i < gaps`) print `spacePerGap` spaces. If there are extra spaces left (`currentExtra > 0`), add one more space and decrement the counter. This distributes the remainder to the first few gaps.

`fmt.Println()` moves to the next terminal line after all words in a row are printed.

`return` — exits the function after justify is done, skipping the padding code below which only applies to left/right/center.

---

### 7.6 Printing the Result

```go
pad := strings.Repeat(" ", spacesNeeded)
```

`strings.Repeat(" ", spacesNeeded)` creates a string of `spacesNeeded` space characters. For `left`, this is `strings.Repeat(" ", 0)` which is `""` — nothing added. For `right` or `center`, it is the appropriate number of spaces.

---

```go
for i, line := range artChar {
    if line == "" {
        if i != len(artChar)-1 {
            fmt.Println()
        }
    } else {
        fmt.Print(pad, line, "\n")
    }
}
```

Loop through every row of the art. Two cases:

**Line is empty:**
Check if this is the very last element in the slice — `i != len(artChar)-1`. The last element is always an empty string left over after splitting on `\n`. We skip it to avoid printing an unwanted extra blank line at the end. Any other empty line (from `\n` in the user's input) does get printed as a blank line.

**Line has content:**
Print the padding first, then the art row, then a newline.

This replaced the original code:

```go
// OLD — skipped ALL blank lines including intentional ones
for _, line := range artChar {
    if line != "" {
        fmt.Print(pad, line, "\n")
    }
}
```

The old version filtered out every empty line — meaning intentional blank lines from `"Hello\nWorld"` were also silently dropped. The fix correctly distinguishes between the trailing empty string (skip it) and intentional blank lines (print them).

---

## 8. The Two Bug Fixes

### Fix 1 — `strings.Fields` instead of manual filter

**Old code:**

```go
rawWords := strings.Split(inputText, " ")
words := []string{}
for _, w := range rawWords {
    if w != "" {
        words = append(words, w)
    }
}
```

`strings.Split(inputText, " ")` produces empty strings when there are multiple spaces between words. The loop then manually filters them out — 6 lines of code.

**Fix:**

```go
words := strings.Fields(inputText)
```

`strings.Fields` splits on any whitespace and automatically skips empty strings — same result in one line.

---

### Fix 2 — Blank line handling

**Old code:**

```go
for _, line := range artChar {
    if line != "" {
        fmt.Print(pad, line, "\n")
    }
}
```

Filtered out every empty line — including intentional blank lines from `\n` in the user's input. So `"Hello\nWorld"` would print both blocks with no gap between them.

**Fix:**

```go
for i, line := range artChar {
    if line == "" {
        if i != len(artChar)-1 {
            fmt.Println()
        }
    } else {
        fmt.Print(pad, line, "\n")
    }
}
```

Now only the very last empty string (always present after the final `\n`) is skipped. All other blank lines print correctly.

---

## 9. Example Runs

```bash
# Default left
go run . "hello"

# Explicitly left
go run . --align=left "hello" standard

# Right — pushed to terminal edge
go run . --align=right "hello" standard

# Center
go run . --align=center "hello" shadow

# Justify — words spread across terminal width
go run . --align=justify "how are you" shadow

# Multi-line
go run . "Hello\nWorld" standard

# Multi-line with double gap
go run . "Hello\n\nWorld" standard

# Wrong format — usage message
go run . --align right something standard

# Too many arguments — usage message
go run . --align=right "hello" standard extra
```

---

## 10. Common Mistakes

**`tput` vs `stty size`**
This code uses `tput cols` to get terminal width. Both work but `tput cols` returns only the column count directly as a single number — simpler to parse than `stty size` which returns `"rows cols"` and requires splitting.

**`cmd.Stdin = os.Stdin` missing**
Without this, `tput` cannot read terminal info and returns an error. Always include it.

**`strconv.Atoi` for string-to-integer**
You cannot do maths directly on a string in Go. Always convert with `strconv.Atoi` before using the number in calculations.

**`strings.Fields` vs `strings.Split`**
Use `strings.Fields` when splitting into words — it handles multiple spaces and never produces empty strings. Use `strings.Split` when you need to split on a specific exact separator like `"\n"`.

**The justify `return`**
The justify case has its own `fmt.Print` loop and ends with `return`. Without that `return`, the code would fall through to the padding loop below and print the art again.

**The trailing empty string**
`strings.Split("hello\n", "\n")` always produces a trailing `""` — splitting `"a\nb\n"` gives `["a", "b", ""]`. Always account for this when looping through split art output.

---

> Built with Go — no external dependencies, just the standard library.
