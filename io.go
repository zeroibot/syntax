package syntax

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cfgFile struct {
	tokenLines   []string
	grammarLines []string
}

// readCfgLines reads a `.cfg` file and returns a cfgFile
func readCfgLines(path string) (*cfgFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open cfg file: %w", err)
	}
	defer file.Close()

	tokenLines := make([]string, 0)
	grammarLines := make([]string, 0)
	tokenMode, grammarMode := false, false

	scanner := bufio.NewScanner(file) // uses default split function ScanLines
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue // skip comments
		} else if line == "tokens:" {
			tokenMode, grammarMode = true, false
		} else if line == "grammar:" {
			grammarMode, tokenMode = true, false
		} else if tokenMode {
			tokenLines = append(tokenLines, line)
		} else if grammarMode {
			grammarLines = append(grammarLines, line)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("cfg scan error: %w", err)
	}
	return new(cfgFile{tokenLines, grammarLines}), nil
}
