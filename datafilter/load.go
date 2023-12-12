package datafilter

import (
	"embed"
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"gopkg.in/yaml.v3"
	"log"
	"path/filepath"
)

// filter types
const (
	PAN   = "pan"
	OWASP = "owasp"
)

type ruleSet struct {
	Type  string  `yaml:"type"`
	Rules []rules `yaml:"rules"`
}

type rules struct {
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
	CustomPath string `yaml:"custom_path"`
}

var cacheIn = cache.NewInMemory()

//go:embed datafilter_rule_set.yaml
var dataFilterRuleSet []byte

//go:embed rules/*
var ruleFs embed.FS

func LoadDataFilterRules(dataFilterRuleSetPath string) {
	var defaultRuleSet []ruleSet
	readDefaultDataFilterRuleSet(&defaultRuleSet)
	ruleSet := checkAndAssignCustomDataFilterRuleSet(dataFilterRuleSetPath, defaultRuleSet)

	for _, set := range ruleSet {
		for _, rule := range set.Rules {
			var patterns, customPatterns []patternValidator
			readPatterns(rule, &patterns, &customPatterns)
			assignPatterns(customPatterns, &patterns)

			validateRule := make([]Validator, len(patterns))
			switch set.Type {
			case PAN:
				for i, v := range patterns {
					v.compileRule()
					validateRule[i] = &pan{patternValidator: v}
				}
			case OWASP:
				for i, v := range patterns {
					v.compileRule()
					validateRule[i] = &owasp{patternValidator: v}
				}
			}
			cacheIn.Set(rule.Name, validateRule, cache.NoExpiration)
		}
	}
	log.Println("datafilter rules have been loaded successfully")
}

func assignPatterns(customPatterns []patternValidator, patterns *[]patternValidator) {
	copyOfPatterns := *patterns
	for i := 0; i < len(customPatterns); i++ {
		patternName := customPatterns[i].Name
		patternIndex := indexOfPatterns(copyOfPatterns, patternName)
		if patternIndex == -1 {
			copyOfPatterns = append(copyOfPatterns, customPatterns[i])
		} else {
			copyOfPatterns[patternIndex] = customPatterns[i]
		}
	}
	patterns = &copyOfPatterns
}

func checkAndAssignCustomDataFilterRuleSet(dataFilterRuleSetPath string, defaultRuleSet []ruleSet) []ruleSet {
	if dataFilterRuleSetPath == "" {
		return defaultRuleSet
	}

	var customRuleSet []ruleSet
	readCustomDataFilterRuleSet(dataFilterRuleSetPath, &customRuleSet)

	for i := 0; i < len(customRuleSet); i++ {
		ruleType := customRuleSet[i].Type
		rsIndex := indexOfRuleSet(defaultRuleSet, ruleType)
		if rsIndex == -1 {
			defaultRuleSet = append(defaultRuleSet, customRuleSet[i])
		} else {
			customRules := customRuleSet[i].Rules
			for k := 0; k < len(customRules); k++ {
				index := indexOfRule(defaultRuleSet[rsIndex].Rules, customRules[k].Name)
				if index == -1 {
					defaultRuleSet[rsIndex].Rules = append(defaultRuleSet[rsIndex].Rules, customRules[k])
				} else {
					(defaultRuleSet[rsIndex]).Rules[index] = customRules[k]
				}
			}
		}
	}

	return defaultRuleSet
}

func readPatterns(rule rules, patterns *[]patternValidator, customPatterns *[]patternValidator) {
	content, err := ruleFs.ReadFile(rule.Path)
	if err != nil {
		panic(err)
	}

	config.MustUnmarshalYaml(rule.Path, content, &patterns)
	if rule.CustomPath != "" {
		content = config.MustReadFile(filepath.Join(config.RootDirectory, rule.CustomPath))
		config.MustUnmarshalYaml(rule.CustomPath, content, &customPatterns)
	}
}

func readCustomDataFilterRuleSet(dataFilterRuleSetPath string, customRuleSet *[]ruleSet) {
	content := config.MustReadFile(filepath.Join(config.RootDirectory, dataFilterRuleSetPath))
	config.MustUnmarshalYaml(dataFilterRuleSetPath, content, &customRuleSet)
}

func readDefaultDataFilterRuleSet(ruleSet *[]ruleSet) {
	err := yaml.Unmarshal(dataFilterRuleSet, &ruleSet)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of datafilter_rule_set.yaml. Error: %s", err)
		panic(msg)
	}
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

func indexOfPatterns(patterns []patternValidator, patternName string) int {
	for i, pattern := range patterns {
		if patternName == pattern.Name {
			return i
		}
	}
	return -1
}
