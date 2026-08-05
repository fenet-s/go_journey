package main

import (
	"fmt"
	"regexp"
	"strings"
)

func WordFrequency(text string) map[string]int {
	text = strings.ToLower(text)

	// Remove punctuation
	re := regexp.MustCompile(`[^\w\s]`)
	text = re.ReplaceAllString(text, "")

	// Split into words
	words := strings.Fields(text)

	// Count frequencies
	frequency := make(map[string]int)

	for _, word := range words {
		frequency[word]++
	}

	return frequency
}

func main() {
	text := "Go is great! Go is fast."

	result := WordFrequency(text)

	fmt.Println(result)
}