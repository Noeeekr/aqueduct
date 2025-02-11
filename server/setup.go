package main

import (
	"bufio"
	"os"
	"strings"
)

var config_values *map[string]interface{} = &map[string]interface{}{}

// Wrapper to setup all configs at once if necessary
func setup() {
	parseConfigFile(config_values)
}

// Get user configurations, calls parseConfigFiles() if setup() wasn't called before
func GetConfig(key string) string {
	if *config_values == nil {
		parseConfigFile(config_values)
	}

	val, ok := (*config_values)[key].(string)
	if ok {
		return val
	}

	LogFatal("Incapaz de encontrar o valor esperado.")

	return val
}

// Parse user configuration files to a map
func parseConfigFile(result_map *map[string]interface{}) {
	file, err := os.Open(config_root)
	if err != nil {
		LogFatal(err.Error())
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#") {
			continue
		}
		if value_pair := strings.Split(scanner.Text(), "="); len(value_pair) == 2 {
			if strings.Trim(value_pair[1], " ") != "" && strings.Trim(value_pair[0], " ") != "" {
				(*result_map)[value_pair[0]] = value_pair[1]
			}
		}
	}
}
