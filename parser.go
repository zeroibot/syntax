package syntax

import (
	"maps"
	"slices"
	"strings"
)

const (
	epsilon string = "EPSILON"
)

type (
	Variable = string
	Terminal = string
	Sentence = []string
)

type Parser struct {
	lexer     *Lexer
	variables []Variable
	terminals []Terminal
	rules     map[Variable][]Sentence
	start     string
}

// NewParser creates a new Parser from the given file path
func NewParser(path string) (*Parser, error) {
	cfg, err := readCfgLines(path)
	if err != nil {
		return nil, err
	}
	return newParserFromCfg(cfg)
}

// NewParserFrom creates a new Parser from the given text
func NewParserFrom(text string) (*Parser, error) {
	cfg, err := createCfg(strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	return newParserFromCfg(cfg)
}

// newParserFromCfg creates a new Parser from the given cfgFile
func newParserFromCfg(cfg *cfgFile) (*Parser, error) {
	lexer, err := newLexerFromLines(cfg.tokenLines)
	if err != nil {
		return nil, err
	}

	parser := &Parser{
		lexer:     lexer,
		variables: make([]Variable, 0),
		terminals: make([]Terminal, 0),
		rules:     make(map[Variable][]Sentence),
	}

	terminals := make(map[string]bool)
	for _, line := range cfg.grammarLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		variable := strings.TrimSpace(parts[0])
		rhs := strings.TrimSpace(parts[1])

		parser.variables = append(parser.variables, variable)
		parser.rules[variable] = make([]Sentence, 0)
		for _, part := range strings.Split(rhs, "|") {
			rule := strings.Fields(strings.TrimSpace(part))
			parser.rules[variable] = append(parser.rules[variable], rule)

			for _, item := range rule {
				if isTerminal(item) {
					terminals[item] = true
				}
			}
		}
	}

	if len(parser.variables) > 0 {
		parser.start = parser.variables[0] // start at the first variable
	}

	if len(terminals) > 0 {
		parser.terminals = slices.Sorted(maps.Keys(terminals))
	}

	return parser, nil
}

// isVariable checks if a token is a variable
func isVariable(token string) bool {
	return strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">")
}

// isTerminal checks if a token is a terminal
func isTerminal(token string) bool {
	return !isVariable(token)
}
