package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

var (
	startTime     time.Time
	wordCount     int
	currentWpm    float64 = 0
	lastWordError bool    = false
)

var wpmHistory = make([]float64, 0, 5)

func calculateWPM(quoteWordCount int, startTime time.Time) float64 {
	// Calculate elapsed time in minutes
	elapsedTime := time.Since(startTime)
	minutes := elapsedTime.Minutes()

	// Calculate words per minute
	wpm := float64(quoteWordCount) / minutes

	return wpm
}

func getNextWord(quote string, currentWordIndex int) string {
	words := strings.Split(quote, " ")
	if currentWordIndex < len(words) {
		return words[currentWordIndex]
	}
	return "" // Return empty if index is out of bounds
}

func getRandomQuote() string {
	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}

func checkEscape(buffer []byte) {
	if buffer[0] == 27 {
		fmt.Println("\nExiting...")
		os.Exit(0)
	}
}

func checkQuit(currentUserWord string) {
	// Check if this is a word completion
	if currentUserWord == "quit" {
		clear()
		printCentered(4, "Goodbye!")
		// close the program immediatly
		os.Exit(0)
	}
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting terminal to raw mode: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState) // Restore on exit

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	// Init
	currentQuote := getRandomQuote()
	currentWordIndex := 0
	currentWord := getNextWord(currentQuote, currentWordIndex)
	currentUserWord := ""
	startTime = time.Now()

	// Clear screen and show initial instructions
	clear()
	printCurrent(currentQuote, currentWord, startTime, currentWordIndex)

	// Buffer for reading characters
	buffer := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			fmt.Println("\nError reading input:", err)
			break
		}

		checkEscape(buffer)

		// Check for backspace (ASCII 8 or 127)
		if buffer[0] == 8 || buffer[0] == 127 {
			if len(currentUserWord) > 0 {
				currentUserWord = currentUserWord[:len(currentUserWord)-1]
				fmt.Print("\b \b")
			}
			continue
		}

		// Check for space (ASCII 32) or enter (ASCII 13)
		if buffer[0] == 32 || buffer[0] == 13 {
			checkQuit(currentUserWord)

			// Check if the word is correct
			if currentUserWord == currentWord {
				lastWordError = false
				currentWordIndex++

				// Check if we've completed the quote WIN
				if currentWordIndex >= len(strings.Split(currentQuote, " ")) {
					quoteWordCount := len(strings.Split(currentQuote, " "))
					addWPMToHistory(calculateWPM(quoteWordCount, startTime))
					goToStart()
					currentQuote = getRandomQuote()
					startTime = time.Now()
					currentWordIndex = 0
				}

				// Get the next word to type
				currentWord = getNextWord(currentQuote, currentWordIndex)
				if currentWord == "" {
					goToStart()
					currentQuote = getRandomQuote()
					currentWordIndex = 0
					currentWord = getNextWord(currentQuote, currentWordIndex)
				}

			} else {
				lastWordError = true
			}

			currentUserWord = ""
			clear()
			printCurrent(currentQuote, currentWord, startTime, currentWordIndex)
			continue
		}

		// Regular character - add to current word and echo it
		if buffer[0] >= 32 && buffer[0] <= 126 { // Printable ASCII
			currentUserWord += string(buffer[0])
			fmt.Print(string(buffer[0]))
		}
	}
}
