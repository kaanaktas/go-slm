package datafilter

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"log"
	"path/filepath"
)

type ruleSet struct {
	Type  string  `json:"type"`
	Rules []rules `json:"rules"`
}

type rules struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	CustomPath string `json:"custom_path"`
}

var cacheIn = cache.NewInMemory()

//go:embed datafilter_rule_set.json
var dataFilterRuleSet []byte

//go:embed rules/*
var ruleFs embed.FS

func Load(dataFilterRuleSetPath string) {
	var ruleSet, customRuleSet []ruleSet
	err := json.Unmarshal(dataFilterRuleSet, &ruleSet)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of datafilter_rule_set.json. Error: %s", err)
		panic(msg)
	}

	if dataFilterRuleSetPath != "" {
		content, err := config.ReadFile(filepath.Join(config.RootDirectory, dataFilterRuleSetPath))
		if err != nil {
			msg := fmt.Sprintf("Error while reading %s. Error: %s", dataFilterRuleSetPath, err)
			panic(msg)
		}
		err = json.Unmarshal(content, &customRuleSet)
		if err != nil {
			msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", dataFilterRuleSetPath, err)
			panic(msg)
		}

		for i := 0; i < len(customRuleSet); i++ {
			ruleType := customRuleSet[i].Type
			rsIndex := indexOfRuleSet(ruleSet, ruleType)
			if rsIndex == -1 {
				ruleSet = append(ruleSet, customRuleSet[i])
			} else {
				customRules := customRuleSet[i].Rules
				for k := 0; k < len(customRules); k++ {
					index := indexOfRule(ruleSet[rsIndex].Rules, customRules[k].Name)
					if index == -1 {
						ruleSet[rsIndex].Rules = append(ruleSet[rsIndex].Rules, customRules[k])
					} else {
						(ruleSet[rsIndex]).Rules[index] = customRules[k]
					}
				}
			}
		}
	}

	for _, set := range ruleSet {
		for _, rule := range set.Rules {
			content, err := ruleFs.ReadFile(rule.Path)
			if err != nil {
				panic(err)
			}

			var patterns, customPatterns []pattern
			err = json.Unmarshal(content, &patterns)
			if err != nil {
				msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", rule.Path, err)
				panic(msg)
			}

			if rule.CustomPath != "" {
				content, err = config.ReadFile(filepath.Join(config.RootDirectory, rule.CustomPath))
				if err != nil {
					msg := fmt.Sprintf("Error while reading %s. Error: %s", rule.CustomPath, err)
					panic(msg)
				}
				err = json.Unmarshal(content, &customPatterns)
				if err != nil {
					msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", rule.CustomPath, err)
					panic(msg)
				}
			}

			for i := 0; i < len(customPatterns); i++ {
				patternName := customPatterns[i].Name
				patternIndex := indexOfPatterns(patterns, patternName)
				if patternIndex == -1 {
					patterns = append(patterns, customPatterns[i])
				} else {
					patterns[patternIndex] = customPatterns[i]
				}
			}

			validateRule := make([]Validate, len(patterns))
			switch set.Type {
			case PAN:
				for i, v := range patterns {
					validateRule[i] = pan{pattern: v}
				}
			case OWASP:
				for i, v := range patterns {
					validateRule[i] = owasp{pattern: v}
				}
			}

			_ = cacheIn.Set(rule.Name, validateRule, cache.NoExpiration)
		}
	}
	log.Println("datafilter rules have been loaded successfully")
}

func indexOfRuleSet(ruleSet []ruleSet, ruleType string) int {
	for i, set := range ruleSet {
		if ruleType == set.Type {
			return i
		}
	}
	return -1
}

func indexOfRule(rules []rules, ruleName string) int {
	for i, rule := range rules {
		if ruleName == rule.Name {
			return i
		}
	}
	return -1
}

func indexOfPatterns(patterns []pattern, patternName string) int {
	for i, pattern := range patterns {
		if patternName == pattern.Name {
			return i
		}
	}
	return -1
}
