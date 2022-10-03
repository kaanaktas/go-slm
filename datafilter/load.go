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

//filter types
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

func Load(dataFilterRuleSetPath string) {
	var ruleSet, customRuleSet []ruleSet
	unmarshalYamlAndPopulateDefaultDataFilterRuleSet(&ruleSet)
	checkAndAssignCustomDataFilterRuleSet(dataFilterRuleSetPath, customRuleSet, &ruleSet)

	for _, set := range ruleSet {
		for _, rule := range set.Rules {
			var patterns, customPatterns []pattern
			unmarshalYamlAndPopulatePatterns(rule, &patterns, &customPatterns)
			assignPatterns(customPatterns, &patterns)

			validateRule := make([]Validate, len(patterns))
			switch set.Type {
			case PAN:
				for i, v := range patterns {
					validateRule[i] = &pan{pattern: v}
				}
			case OWASP:
				for i, v := range patterns {
					validateRule[i] = &owasp{pattern: v}
				}
			}
			cacheIn.Set(rule.Name, validateRule, cache.NoExpiration)
		}
	}
	log.Println("datafilter rules have been loaded successfully")
}

func assignPatterns(customPatterns []pattern, patterns *[]pattern) {
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

func checkAndAssignCustomDataFilterRuleSet(dataFilterRuleSetPath string, customRuleSet []ruleSet, ruleSet *[]ruleSet) {
	if dataFilterRuleSetPath == "" {
		return
	}
	unmarshalYamlAndPopulateCustomDataFilterRuleSet(dataFilterRuleSetPath, &customRuleSet)

	copyOfRuleSet := *ruleSet
	for i := 0; i < len(customRuleSet); i++ {
		ruleType := customRuleSet[i].Type
		rsIndex := indexOfRuleSet(copyOfRuleSet, ruleType)
		if rsIndex == -1 {
			copyOfRuleSet = append(copyOfRuleSet, customRuleSet[i])
		} else {
			customRules := customRuleSet[i].Rules
			for k := 0; k < len(customRules); k++ {
				index := indexOfRule(copyOfRuleSet[rsIndex].Rules, customRules[k].Name)
				if index == -1 {
					copyOfRuleSet[rsIndex].Rules = append(copyOfRuleSet[rsIndex].Rules, customRules[k])
				} else {
					(copyOfRuleSet[rsIndex]).Rules[index] = customRules[k]
				}
			}
		}
	}
	ruleSet = &copyOfRuleSet
}

func unmarshalYamlAndPopulatePatterns(rule rules, patterns *[]pattern, customPatterns *[]pattern) {
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

func unmarshalYamlAndPopulateCustomDataFilterRuleSet(dataFilterRuleSetPath string, customRuleSet *[]ruleSet) {
	content := config.MustReadFile(filepath.Join(config.RootDirectory, dataFilterRuleSetPath))
	config.MustUnmarshalYaml(dataFilterRuleSetPath, content, &customRuleSet)
}

func unmarshalYamlAndPopulateDefaultDataFilterRuleSet(ruleSet *[]ruleSet) {
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

func indexOfPatterns(patterns []pattern, patternName string) int {
	for i, pattern := range patterns {
		if patternName == pattern.Name {
			return i
		}
	}
	return -1
}
