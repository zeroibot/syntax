package syntax

import "fmt"

type deriveStep struct {
	Sentence
	Tokens []Token
}

func (parser *Parser) CheckSyntax(text string, ignoreTypes []string) error {
	// Tokenize lines
	if parser.lexer == nil {
		return fmt.Errorf("lexer not initialized")
	}

	tokens, err := parser.lexer.Tokenize(text, ignoreTypes)
	if err != nil {
		return err
	}

	// Set first step as last step with tokens
	lastStep := new(deriveStep{
		Sentence: Sentence{parser.start},
		Tokens:   tokens,
	})
	q := make([]*deriveStep, 0)
	q = append(q, lastStep)

	// Invariant: sentence is non-empty and first word is a variable
	// Invariant: replacement rule must start with terminal or is a single variable
	for len(q) > 0 {
		step := q[0]
		q = q[1:]

		if len(step.Tokens) > 0 {
			lastStep = step // set as last step if step has tokens
		}

		if len(step.Sentence) == 0 {
			continue // skip empty sentences
		}

		variable := step.Sentence[0]
		for _, rule := range parser.getReplacements(variable) {
			// align front
			equation := newEquation(rule, step.Sentence)
			result, ok := alignFront(equation, step.Tokens)
			if !ok {
				continue // skip if not aligned
			}
			emptyLeft := len(result.Sentence) == 0
			emptyRight := len(result.Tokens) == 0
			if emptyLeft && emptyRight {
				// both sides are fully consumed = success
				return nil
			}
			// add to queue
			q = append(q, result)
		}
	}

	// Exit loop = queue is empty, failed to parse
	token := lastStep.Tokens[0]
	limit := min(10, len(token.Text))
	return fmt.Errorf("syntax error: unexpected %s at line %d, col %d", token.Text[:limit], token.Row+1, token.ColStart+1)
}

// Get replacement rules for variable
func (parser *Parser) getReplacements(start Variable) []Sentence {
	q := make([]string, 0)
	q = append(q, start)

	replacements := make([]Sentence, 0)
	for len(q) > 0 {
		variable := q[0]
		q = q[1:]

		for _, rule := range parser.rules[variable] {
			first := rule[0]
			if isTerminal(first) {
				// Add to replacement if rule's first word is terminal
				replacements = append(replacements, rule)
			} else {
				// If variable, enqueue so it can be expanded
				q = append(q, first)
			}
		}
	}
	return replacements
}

// Create a new equation
func newEquation(rule, prev Sentence) Sentence {
	equation := make(Sentence, 0)
	equation = append(equation, rule...)     // add replacement rule to front
	equation = append(equation, prev[1:]...) // add the rest of the sentence to the back
	validTokens := make(Sentence, 0, len(equation))
	for _, token := range equation {
		// Only add tokens that are not epsilon
		if token != epsilon {
			validTokens = append(validTokens, token)
		}
	}
	return validTokens
}

// Align sentence and tokens from the front
func alignFront(equation Sentence, tokens []Token) (*deriveStep, bool) {
	limit := min(len(equation), len(tokens))
	for i := range limit {
		left, right := equation[i], tokens[i]
		if isTerminal(left) {
			// equation has terminal in front, try to match with token
			if left != right.Type {
				return nil, false
			}
		} else {
			// equation now has variable in front, stop here
			step := &deriveStep{
				Sentence: listCopy(equation[i:]),
				Tokens:   listCopy(tokens[i:]),
			}
			return step, true
		}
	}
	// Everything matched so far, stop here
	step := &deriveStep{
		Sentence: listCopy(equation[limit:]),
		Tokens:   listCopy(tokens[limit:]),
	}
	return step, true
}

func listCopy[T any](items []T) []T {
	return append([]T{}, items...)
}
