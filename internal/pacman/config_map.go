package pacman

import "strings"

type (
	sectionMap map[string]string
	configMap  map[string]sectionMap
)

func newConfigMap() configMap {
	return make(configMap)
}

func (cm configMap) Get(section, key string) (string, bool) {
	if sec, ok := cm[section]; ok {
		if val, ok := sec[key]; ok {
			return val, true
		}
	}
	return "", false
}

func (cm configMap) Add(section, key, value string) {
	if _, ok := cm[section]; !ok {
		cm[section] = make(sectionMap)
	}
	if _, ok := cm[section][key]; !ok {
		cm[section][key] = value
	}
}

func createConfigMap(lines []string) configMap {
	cm := newConfigMap()
	var sectionName string
	const expectedParts = 2
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if isSection(line) {
			sectionName = extractSectionName(line)
			continue
		}

		keyValuePair := strings.SplitN(line, "=", expectedParts)
		key := strings.TrimSpace(keyValuePair[0])

		value := ""
		if len(keyValuePair) == expectedParts {
			value = strings.TrimSpace(keyValuePair[1])
		}
		cm.Add(sectionName, key, value)
	}
	return cm
}

func isSection(line string) bool {
	return len(line) > 2 && strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
}

func extractSectionName(line string) string {
	return line[1 : len(line)-1]
}
