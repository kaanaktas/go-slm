package policy

import (
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"log"
	"path/filepath"
	"sort"
)

const Key = "policy_rule"

var cacheIn = cache.NewInMemory()

type Action struct {
	Name   string `yaml:"name"`
	Active bool   `yaml:"active"`
	Order  int    `yaml:"order"`
}

type statement struct {
	Type    string   `yaml:"type"`
	Order   int      `yaml:"order"`
	Actions []Action `yaml:"action"`
}

type commonPolicies struct {
	Policy struct {
		Name       string      `yaml:"name"`
		Statements []statement `yaml:"statement"`
	} `json:"policy"`
}

type policy struct {
	ServiceName string `yaml:"serviceName"`
	Request     string `yaml:"request"`
	Response    string `yaml:"response"`
}

type Statements map[string][]statement

func Load(policyRuleSetPath, commonRulesPath string) {
	if policyRuleSetPath == "" {
		panic("GO_SLM_POLICY_RULE_SET_PATH hasn't been set")
	}
	if commonRulesPath == "" {
		panic("GO_SLM_COMMON_POLICIES_PATH hasn't been set")
	}

	var policies []policy
	unmarshalYamlAndPopulatePolicies(policyRuleSetPath, &policies)

	var retrievedCommonPolicies []commonPolicies
	unmarshalYamlAndPopulateCommonPolicies(commonRulesPath, &retrievedCommonPolicies)

	statements := make(Statements)
	for _, policy := range policies {
		for _, rule := range retrievedCommonPolicies {
			if rule.Policy.Name == policy.Request {
				populatePolicyRules(rule, statements, config.PolicyKey(policy.ServiceName, config.Request))
			}
			if rule.Policy.Name == policy.Response {
				populatePolicyRules(rule, statements, config.PolicyKey(policy.ServiceName, config.Response))
			}
		}
	}
	cacheIn.Set(Key, statements, cache.NoExpiration)

	log.Println("common policies have been loaded successfully")
}

func unmarshalYamlAndPopulateCommonPolicies(commonRulesPath string, retrievedCommonPolicies *[]commonPolicies) {
	content := config.MustReadFile(filepath.Join(config.RootDirectory, commonRulesPath))
	config.MustUnmarshalYaml(commonRulesPath, content, retrievedCommonPolicies)
}

func unmarshalYamlAndPopulatePolicies(policyRuleSetPath string, policies *[]policy) {
	content := config.MustReadFile(filepath.Join(config.RootDirectory, policyRuleSetPath))
	config.MustUnmarshalYaml(policyRuleSetPath, content, policies)
}

func populatePolicyRules(rule commonPolicies, policyRules Statements, policyRuleKey string) {
	sort.Slice(rule.Policy.Statements, func(i, j int) bool {
		return rule.Policy.Statements[i].Order < rule.Policy.Statements[j].Order
	})
	policyRules[policyRuleKey] = rule.Policy.Statements
}
