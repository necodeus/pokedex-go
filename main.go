package main

import (
	"bufio"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text) // remove leading and trailing spaces
	text = strings.ToLower(text)   // convert to lowercase

	return strings.Fields(text) // split the text into words
}

func main() {
	// wait for user input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		print("Pokedex > ")

		// read the input
		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) > 0 {
			print("Your command was: " + words[0] + "\n")
		}
	}
}
