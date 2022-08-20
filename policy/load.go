package policy

import (
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"gopkg.in/yaml.v3"
	"log"
	"path/filepath"
)

const Key = "policy_rule"

var cacheIn = cache.NewInMemory()

type CommonPolicy struct {
	Name   string `yaml:"name"`
	Active bool   `yaml:"active"`
}

type commonPolicies struct {
	PolicyName string         `yaml:"PolicyName"`
	Policy     []CommonPolicy `yaml:"Policy"`
}

type policy struct {
	ServiceName string `yaml:"serviceName"`
	Request     string `yaml:"request"`
	Response    string `yaml:"response"`
}

type CommonPolicies map[string][]CommonPolicy

func Load(policyRuleSetPath, commonRulesPath string) {
	if policyRuleSetPath == "" {
		panic("GO_SLM_POLICY_RULE_SET_PATH hasn't been set")
	}

	if commonRulesPath == "" {
		panic("GO_SLM_COMMON_POLICIES_PATH hasn't been set")
	}

	var policies []policy
	content, err := config.ReadFile(filepath.Join(config.RootDirectory, policyRuleSetPath))
	if err != nil {
		msg := fmt.Sprintf("Error while reading %s. Error: %s", policyRuleSetPath, err)
		panic(msg)
	}

	err = yaml.Unmarshal(content, &policies)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", policyRuleSetPath, err)
		panic(msg)
	}

	var retrievedCommonPolicies []commonPolicies
	content, err = config.ReadFile(filepath.Join(config.RootDirectory, commonRulesPath))
	if err != nil {
		msg := fmt.Sprintf("Error while reading %s. Error: %s", commonRulesPath, err)
		panic(msg)
	}

	err = yaml.Unmarshal(content, &retrievedCommonPolicies)
	if err != nil {
		msg := fmt.Sprintf("Can't unmarshall the content of %s. Error: %s", commonRulesPath, err)
		panic(msg)
	}

	policyRules := make(CommonPolicies)

	for _, policy := range policies {
		for _, rule := range retrievedCommonPolicies {
			if rule.PolicyName == policy.Request {
				policyRules[config.PolicyKey(policy.ServiceName, config.Request)] = rule.Policy
			}
			if rule.PolicyName == policy.Response {
				policyRules[config.PolicyKey(policy.ServiceName, config.Response)] = rule.Policy
			}
		}
	}

	_ = cacheIn.Set(Key, policyRules, cache.NoExpiration)
	log.Println("policy commonPolicies have been loaded successfully")
}
