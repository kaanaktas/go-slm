package policy

import (
	"encoding/json"
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"log"
	"path/filepath"
)

const Key = "access_rule"

var cacheIn = cache.NewInMemory()

type CommonRule struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type CommonRules struct {
	Name  string       `json:"Name"`
	Rules []CommonRule `json:"rule"`
}

type CommonRuleSet struct {
	Rules []CommonRules `json:"Rules"`
}

type Policy struct {
	ServiceName string `json:"serviceName"`
	Request     string `json:"request"`
	Response    string `json:"response"`
}

type Policies struct {
	Policies []Policy `json:"policies"`
}

type Rules map[string][]CommonRule

func Load(policyRuleSetPath, commonRulesPath string) {
	if policyRuleSetPath == "" {
		panic("POLICY_RULE_SET_PATH hasn't been set")
	}

	if commonRulesPath == "" {
		panic("COMMON_RULES_PATH hasn't been set")
	}

	var ps Policies
	content, err := config.ReadFile(filepath.Join(config.RootDirectory, policyRuleSetPath))
	if err != nil {
		msg := fmt.Sprintf("Error while reading %s. Error: %s", policyRuleSetPath, err)
		panic(msg)
	}

	err = json.Unmarshal(content, &ps)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", policyRuleSetPath, err)
		panic(msg)
	}

	var rules CommonRuleSet
	content, err = config.ReadFile(filepath.Join(config.RootDirectory, commonRulesPath))
	if err != nil {
		msg := fmt.Sprintf("Error while reading %s. Error: %s", commonRulesPath, err)
		panic(msg)
	}

	err = json.Unmarshal(content, &rules)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", commonRulesPath, err)
		panic(msg)
	}

	policyRules := make(Rules)

	for _, policy := range ps.Policies {
		for _, rule := range rules.Rules {
			if rule.Name == policy.Request {
				policyRules[config.PolicyKey(policy.ServiceName, config.Request)] = rule.Rules
			}
			if rule.Name == policy.Response {
				policyRules[config.PolicyKey(policy.ServiceName, config.Response)] = rule.Rules
			}
		}
	}

	_ = cacheIn.Set(Key, policyRules, cache.NoExpiration)
	log.Println("policy CommonRules have been loaded successfully")
}
