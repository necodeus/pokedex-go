package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)

	words := strings.Fields(trimmed)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return words
}

func main() {
	fmt.Println("Hello, World!")
}
