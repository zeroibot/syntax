package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/roidaradal/syntax"
)

func main() {
	text, err := readFile("input/all.txt")
	if err != nil {
		log.Fatal(fmt.Errorf("failed reading input file: %w", err))
	}

	lexer, err := syntax.NewLexer("cfg/json.cfg")
	if err != nil {
		log.Fatal(fmt.Errorf("failed creating new lexer: %w", err))
	}

	ignore := []string{"LIT_WS"}
	for line := range strings.Lines(text) {
		tokens, err := lexer.Tokenize(line, ignore)
		if err != nil {
			fmt.Printf("Error in tokenizing %q: %v\n", line, err)
			continue
		}
		fmt.Println(tokens)
	}
}

func readFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
