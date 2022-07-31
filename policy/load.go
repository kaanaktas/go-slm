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

type CommonPolicy struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type commonPolicies struct {
	CommonPolicyName string         `json:"PolicyName"`
	Policy           []CommonPolicy `json:"Policy"`
}

type commonPolicySet struct {
	CommonPolicies []commonPolicies `json:"commonPolicies"`
}

type policy struct {
	ServiceName string `json:"serviceName"`
	Request     string `json:"request"`
	Response    string `json:"response"`
}

type policies struct {
	Policies []policy `json:"policies"`
}

type CommonPolicyMap map[string][]CommonPolicy

func Load(policyRuleSetPath, commonRulesPath string) {
	if policyRuleSetPath == "" {
		panic("GO_SLM_POLICY_RULE_SET_PATH hasn't been set")
	}

	if commonRulesPath == "" {
		panic("GO_SLM_COMMON_POLICIES_PATH hasn't been set")
	}

	var ps policies
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

	var rules commonPolicySet
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

	policyRules := make(CommonPolicyMap)

	for _, policy := range ps.Policies {
		for _, rule := range rules.CommonPolicies {
			if rule.CommonPolicyName == policy.Request {
				policyRules[config.PolicyKey(policy.ServiceName, config.Request)] = rule.Policy
			}
			if rule.CommonPolicyName == policy.Response {
				policyRules[config.PolicyKey(policy.ServiceName, config.Response)] = rule.Policy
			}
		}
	}

	_ = cacheIn.Set(Key, policyRules, cache.NoExpiration)
	log.Println("policy commonPolicies have been loaded successfully")
}
