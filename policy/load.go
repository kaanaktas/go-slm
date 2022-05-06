package access

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"log"
)

const Key = "access_rule"

func init() {
	loadAccesses()
}

type Rule struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type Rules struct {
	Name  string `json:"Name"`
	Rules []Rule `json:"rule"`
}

type RuleSet struct {
	Rules []Rules `json:"Rules"`
}

type Policy struct {
	ServiceName string `json:"serviceName"`
	Request     string `json:"request"`
	Response    string `json:"response"`
}

type Policies struct {
	Policies []Policy `json:"policies"`
}

type PolicyRules map[string][]Rule

var cacheIn = cache.NewInMemory()

//go:embed policy_rule_set.json
var policyRuleSetContent []byte

//go:embed common_rules.json
var commonRulesContent []byte

func loadAccesses() {
	var ps Policies
	err := json.Unmarshal(policyRuleSetContent, &ps)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of policy_rule_set.json. Error: %s", err)
		panic(msg)
	}

	var rules RuleSet
	err = json.Unmarshal(commonRulesContent, &rules)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of policy_rule_set.json. Error: %s", err)
		panic(msg)
	}

	policyRules := make(PolicyRules)

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
	log.Println("policy Rules have been loaded successfully")
}
