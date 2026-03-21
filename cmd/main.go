package main

import (
	"fmt"
	"log"
	"os"

	"github.com/roidaradal/syntax"
)

func main() {
	text, err := readFile("input/1.json")
	if err != nil {
		log.Fatal(fmt.Errorf("failed reading input file: %w", err))
	}

	lexer, err := syntax.NewLexer("cfg/json.cfg")
	if err != nil {
		log.Fatal(fmt.Errorf("failed creating new lexer: %w", err))
	}

	tokens, err := lexer.Tokenize(text, []string{"LIT_WS"})
	if err != nil {
		log.Fatal(fmt.Errorf("failed tokenizing: %w", err))
	}

	for i, token := range tokens {
		fmt.Println(i+1, token)
	}
}

func readFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
