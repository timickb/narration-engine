package config

import (
	"errors"
	"fmt"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/parser"
	"os"
	"strings"
)

func readScenarios(path string) ([]*domain.Scenario, error) {
	entries, _ := os.ReadDir(path)
	handledScenarios := make([]*domain.Scenario, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			children, err := readScenarios(path + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			handledScenarios = append(handledScenarios, children...)
		} else {
			scenario, err := readScenario(path, entry)
			if err != nil {
				return nil, err
			}
			handledScenarios = append(handledScenarios, scenario)
		}
	}

	return handledScenarios, nil
}

func readScenario(path string, entry os.DirEntry) (*domain.Scenario, error) {
	nameParts := strings.Split(entry.Name(), ".")
	if len(nameParts) < 3 {
		return nil, errors.New("invalid file name. correct format is <scenario_name>.<scenario_version>.puml")
	}
	if nameParts[len(nameParts)-1] != "puml" {
		return nil, errors.New("invalid file extension")
	}

	scenario, err := parser.New().Parse(path + "/" + entry.Name())
	if err != nil {
		return nil, fmt.Errorf("parse scenario: %w", err)
	}
	return scenario, nil
}
