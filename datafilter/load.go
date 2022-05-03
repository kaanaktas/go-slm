package datafilter

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"log"
)

func init() {
	loadRules()
}

type ruleSet struct {
	Type  string  `json:"type"`
	Rules []rules `json:"rules"`
}

type rules struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

//go:embed datafilter_rule_set.json
var dataFilterRuleSet []byte

//go:embed rules/*
var ruleFs embed.FS

func loadRules() {
	var ruleSet []ruleSet
	err := json.Unmarshal(dataFilterRuleSet, &ruleSet)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of datafilter_rule_set.json. Error: %s", err)
		panic(msg)
	}

	for _, set := range ruleSet {
		for _, rule := range set.Rules {
			content, err := ruleFs.ReadFile(rule.Path)
			if err != nil {
				panic(err)
			}

			var patterns []pattern
			err = json.Unmarshal(content, &patterns)
			if err != nil {
				msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", rule.Path, err)
				panic(msg)
			}

			validateRule := make([]validate, len(patterns))
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
