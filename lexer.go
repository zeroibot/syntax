package syntax

import (
	"fmt"
	"regexp"
	"strings"
)

type Token struct {
	Type     string
	Text     string
	Row      int
	ColStart int
	ColEnd   int
}

// Coords returns the (Row: ColStart-ColEnd) as string
func (t Token) Coords() string {
	return fmt.Sprintf("(%d: %d-%d)", t.Row+1, t.ColStart+1, t.ColEnd)
}

// String returns the string representation of Token
func (t Token) String() string {
	if strings.HasPrefix(t.Type, "LIT_") {
		// Literal Token, no need to display inner text
		return t.Type[4:]
	} else if strings.HasPrefix(t.Type, "KW_") {
		// Keyword Token, no need to display inner text
		return t.Type[3:]
	}
	return fmt.Sprintf("%s(%s)", t.Type, t.Text)
}

type Lexer struct {
	tokenTypes []string
	patterns   []*regexp.Regexp
}

// NewLexer creates a new Lexer with token types and patterns from the given file path
func NewLexer(path string) (*Lexer, error) {
	cfg, err := readCfgLines(path)
	if err != nil {
		return nil, err
	}
	return newLexerFromLines(cfg.tokenLines)
}

// newLexerFromLines creates a new Lexer from the list of token lines
func newLexerFromLines(lines []string) (*Lexer, error) {
	numLines := len(lines)
	tokenTypes := make([]string, 0, numLines)
	patterns := make([]*regexp.Regexp, 0, numLines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		tokenType := strings.TrimSpace(parts[0])
		tokenPattern := strings.TrimSpace(parts[1])
		pattern, err := regexp.Compile("^" + tokenPattern) // prefix with ^ so we match the front, not substrings
		if err != nil {
			return nil, fmt.Errorf("invalid token pattern %s : %w", tokenPattern, err)
		}
		tokenTypes = append(tokenTypes, tokenType)
		patterns = append(patterns, pattern)
	}
	lexer := new(Lexer{tokenTypes: tokenTypes, patterns: patterns})
	return lexer, nil
}

// Tokenize tokenizes the given string, and returns the list of Tokens.
// We can pass in a list of TokenTypes to ignore (pass nil if nothing to ignore).
func (lexer *Lexer) Tokenize(text string, ignoreTypes []string) ([]Token, error) {
	ignore := make(map[string]bool)
	for _, tokenType := range ignoreTypes {
		ignore[tokenType] = true
	}
	tokens := make([]Token, 0)
	row := 0
	for line := range strings.Lines(text) {
		lineTokens, err := lexer.tokenizeLine(row, []byte(line), ignore)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, lineTokens...)
		row += 1
	}
	return tokens, nil
}

// tokenizeLine tokenizes one line
func (lexer *Lexer) tokenizeLine(row int, line []byte, ignoreType map[string]bool) ([]Token, error) {
	// Tokenize line by checking each pattern for a match, until line is fully consumed or no pattern matches remaining line
	tokens := make([]Token, 0)
	col := 0 // keeps track of current column index from original line
	for len(line) > 0 {
		found := false
		for i, pattern := range lexer.patterns {
			match := pattern.FindIndex(line)
			if match == nil {
				continue // move on to next if not a match
			}
			start, end := match[0], match[1]
			chunk := string(line[start:end]) // get chunk of text matched by pattern
			tokenType := lexer.tokenTypes[i] // get corresponding token type
			line = line[end:]                // consume chunk and get remaining line
			if !ignoreType[tokenType] {
				// Add to tokens if token type is not ignored
				tokens = append(tokens, Token{
					Type:     tokenType,
					Text:     chunk,
					Row:      row,
					ColStart: col + start,
					ColEnd:   col + end,
				})
			}
			// Move column index forward
			col += end - start
			found = true
			break
		}
		if !found {
			limit := min(10, len(line))
			return nil, fmt.Errorf("syntax error: unexpected %s at line %d, col %d", string(line[:limit]), row+1, col+1)
		}
	}
	return tokens, nil
}
