package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zeroibot/syntax"
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

/*
OLD CODE:

const usage string = "Usage: nlp <tokenize|check> cfg={PATH} <file={PATH} | text={TEXT}> (ignore={TYPE1,TYPE2,...})"

type Config struct {
	path   string
	lines  [][]byte
	ignore *ds.Set[string]
}

func main() {
	var err error
	command, options := io.GetCommandOptions("")
	cfg, err := getConfig(options)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch command {
	case "tokenize":
		err = cmdTokenize(cfg)
	case "check":
		err = cmdCheck(cfg)
	default:
		fmt.Println(usage)
	}
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func getConfig(options dict.StringMap) (*Config, error) {
	// Required: cfg={PATH}
	// Required: file={PATH} | text={TEXT}
	// Options:  ignore={TYPE1,TYPE2,...}

	// Process options
	cfg := &Config{ignore: ds.NewSet[string]()}
	filePath, text := "", ""
	for k, v := range options {
		switch k {
		case "cfg":
			cfg.path = v
		case "file":
			filePath = v
		case "text":
			text = v
		case "ignore":
			cfg.ignore.AddItems(strings.Split(v, ","))
		}
	}

	// Check if cfgPath is set and either filePath or text is set
	if cfg.path == "" || (filePath == "" && text == "") {
		fmt.Println(usage)
		return nil, errors.New("missing params")
	}

	if filePath != "" {
		// Read lines from filePath
		lines, err := nlp.ReadLineBytes(filePath)
		if err != nil {
			return nil, err
		}
		cfg.lines = lines
	} else {
		// Read lines from text
		cfg.lines = list.Map(str.Lines(text), str.ToBytes)
	}

	return cfg, nil
}

// Tokenize command handler
func cmdTokenize(cfg *Config) error {
	// Create lexer from cfgPath
	lexer, err := nlp.LoadLexer(cfg.path)
	if err != nil {
		return err
	}

	// Tokenize
	tokens, err := lexer.Tokenize(cfg.lines, cfg.ignore)
	if err != nil {
		return err
	}
	numTokens := len(tokens)
	fmt.Println("Tokens:", numTokens)

	if numTokens == 0 {
		return nil
	}

	// Display tokens
	maxNum := len(str.Int(len(tokens)))
	maxLength := slices.Max(list.Map(tokens, func(token nlp.Token) int {
		return len(token.Type)
	}))
	maxCoords := slices.Max(list.Map(tokens, func(token nlp.Token) int {
		return len(token.Coords())
	}))
	template := fmt.Sprintf("[%%%dd] %%-%ds : %%-%ds %%s\n", maxNum, maxLength, maxCoords)
	for i, token := range tokens {
		fmt.Printf(template, i+1, token.Type, token.Coords(), token.Text)
	}
	return nil
}

// Check command handler
func cmdCheck(cfg *Config) error {
	// Create parser from cfgPath
	parser, err := nlp.LoadParser(cfg.path)
	if err != nil {
		return err
	}

	err = parser.CheckSyntaxError(cfg.lines, cfg.ignore)
	if err != nil {
		return err
	}
	fmt.Println("Parse: OK")

	return nil
}
*/
