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
	}
	return fmt.Sprintf("%s(%s)", t.Type, t.Text)
}

type Lexer struct {
	tokenTypes []string
	patterns   []*regexp.Regexp
}

// NewLexerFromLines creates a new Lexer with token types from given lines
func NewLexerFromLines(lines []string) (*Lexer, error) {
	numLines := len(lines)
	tokenTypes := make([]string, 0, numLines)
	patterns := make([]*regexp.Regexp, 0, numLines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
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
