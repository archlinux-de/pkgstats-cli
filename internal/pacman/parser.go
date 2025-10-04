package pacman

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configMaxRecursion = 10
	repository         = "core"
	defaultDBPath      = "/var/lib/pacman"
)

type config struct {
	dbPath string
	server string
}

type Pacman struct {
	config *config
}

type parser struct {
	includeFileCache map[string][]string
}

func Parse(pacmanConfPath string) (*Pacman, error) {
	conf, err := newParser().parseConfigFile(pacmanConfPath, defaultDBPath)
	if err != nil {
		return nil, err
	}

	return &Pacman{config: conf}, nil
}

func newParser() *parser {
	return &parser{
		includeFileCache: make(map[string][]string),
	}
}

func (parser *parser) parseConfigFile(configPath string, defaultDBPath string) (*config, error) {
	lines, err := parser.readConfigFile(configPath, 0)
	if err != nil {
		return nil, err
	}

	confMap := createConfigMap(lines)

	dbPath, ok := confMap.Get("options", "DBPath")
	if !ok {
		dbPath = defaultDBPath
	}
	server, ok := confMap.Get(repository, "Server")
	if !ok {
		return nil, errors.New("server not found")
	}
	return &config{
		dbPath: dbPath,
		server: server,
	}, nil
}

func (parser *parser) readConfigFile(fileName string, depth int) ([]string, error) {
	if depth >= configMaxRecursion {
		return nil, fmt.Errorf("reached maximum Include depth of %d", depth)
	}
	readFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	var fileLines []string
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if isCommentOrEmpty(line) {
			continue
		}
		if isIncludeLine(line) {
			includeLines, err := parser.processIncludeLine(line, depth+1)
			if err != nil {
				return nil, err
			}
			fileLines = append(fileLines, includeLines...)
		} else {
			fileLines = append(fileLines, line)
		}
	}
	if err := fileScanner.Err(); err != nil {
		return nil, err
	}
	return fileLines, nil
}

func (parser *parser) processIncludeLine(line string, depth int) ([]string, error) {
	const expectedParts = 2
	includeFileName := strings.TrimSpace(strings.SplitN(line, "=", expectedParts)[1])
	includeFileNames, err := filepath.Glob(includeFileName)
	if err != nil {
		return nil, err
	}

	var allIncludedLines []string
	for _, fileName := range includeFileNames {
		if cachedLines, ok := parser.includeFileCache[fileName]; ok {
			allIncludedLines = append(allIncludedLines, cachedLines...)
		} else {
			includedLines, err := parser.readConfigFile(fileName, depth)
			if err != nil {
				return nil, err
			}
			parser.includeFileCache[fileName] = includedLines
			allIncludedLines = append(allIncludedLines, includedLines...)
		}
	}
	return allIncludedLines, nil
}

func isCommentOrEmpty(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func isIncludeLine(line string) bool {
	return strings.HasPrefix(line, "Include")
}
