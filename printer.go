package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

func printCurrent(currentQuote string, currentWord string, startTime time.Time, currentWordIndex int) {
	// ANSI color codes
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorGreen  = "\033[32m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[34m"
		colorPurple = "\033[35m"
		colorCyan   = "\033[36m"
		colorWhite  = "\033[37m"
	)

	PrintWPMScoreboard()

	words := strings.Split(currentQuote, " ")

	if currentWordIndex != 0 {
		currentWpm = calculateWPM(currentWordIndex, startTime)
	}
	printCentered(0, colorYellow+"=== TYPING SPEED TEST ==="+colorReset)
	printCentered(1, "\nWPM:"+strconv.FormatFloat(currentWpm, 'f', 2, 64))
	//goToStart()

	quoteStr := ""

	// Print each word with the current word highlighted
	currentWordColor := colorBlue
	if lastWordError {
		currentWordColor = colorRed
	}
	for i, word := range words {
		if word == currentWord {
			quoteStr += currentWordColor + word + colorReset
		} else {
			quoteStr += word
		}

		if i < len(words)-1 {
			quoteStr += " "
		}
	}

	printCentered(2, quoteStr)
	printCentered(4, "> ")
}

func goToStart() {
	fmt.Print("\r\033[K")
}

func clear() {
	fmt.Print("\033[H\033[2J")
}

func PrintWPMScoreboard() {
	const (
		colorReset  = "\033[0m"
		colorGreen  = "\033[32m"
		colorYellow = "\033[33m"
		colorBlue   = "\033[34m"
	)
	goToStart()
	fmt.Println("\n" + colorYellow + "=== TYPING SPEED HISTORY ===" + colorReset)
	goToStart()

	if len(wpmHistory) == 0 {
		fmt.Println("No previous scores yet")
		return
	}

	// Calculate average WPM
	var total float64
	for _, wpm := range wpmHistory {
		total += wpm
	}
	avgWPM := total / float64(len(wpmHistory))

	// Print each score with formatting
	goToStart()
	fmt.Println(colorBlue + "Your last 5 scores (WPM):" + colorReset)
	for i, wpm := range wpmHistory {
		goToStart()
		fmt.Printf("  %d. %s%.1f WPM%s\n", i+1, colorGreen, wpm, colorReset)
	}

	goToStart()
	// Print the average
	fmt.Printf("\n%sAverage: %.1f WPM%s\n", colorYellow, avgWPM, colorReset)
	goToStart()
	fmt.Println(colorYellow + "=========================" + colorReset)
}

func printAtPosition(row, col int, text string) {
	// ANSI escape code to move cursor to position
	fmt.Printf("\033[%d;%dH%s", row, col, text)
}

// Function to get terminal width and height
func getTerminalSize() (width, height int, err error) {
	width, height, err = term.GetSize(int(os.Stdin.Fd()))
	return
}

// Function to add a new WPM score to the history
func addWPMToHistory(wpm float64) {
	// Add the new score to the history
	wpmHistory = append(wpmHistory, wpm)

	// If we have more than 5 scores, remove the oldest one
	if len(wpmHistory) > 5 {
		wpmHistory = wpmHistory[1:]
	}
}

func printCentered(row int, text string) {
	width, _, err := getTerminalSize()
	if err != nil {
		width = 80 // Default width if we can't get terminal size
	}

	// Calculate padding to center the text
	padding := (width - len(stripANSI(text))) / 2
	if padding < 0 {
		padding = 0
	}

	// Print spaces for padding, then the text
	printAtPosition(row+10, padding, text)
}

// Helper function to strip ANSI escape codes for length calculations
func stripANSI(str string) string {
	ansiRegex := regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")
	return ansiRegex.ReplaceAllString(str, "")
}
