package access

import (
	_ "embed"
	"encoding/json"
	"github.com/kaanaktas/go-slm/cache"
	"log"
)

const Key = "access_rule"

func init() {
	loadAccesses()
}

type RuleSet map[string][]Rule

type Rule struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type access []struct {
	ServiceName string `json:"serviceName"`
	Rules       []Rule `json:"rules"`
}

var cacheIn = cache.NewInMemory()

//go:embed access_rule_set.json
var accessRuleSet []byte

func loadAccesses() {
	var acc access
	err := json.Unmarshal(accessRuleSet, &acc)
	if err != nil {
		panic("Can't unmarshall the content of access_rule_set.json")
	}

	ruleSet := make(RuleSet)

	for _, v := range acc {
		ruleSet[v.ServiceName] = v.Rules
	}

	_ = cacheIn.Set(Key, ruleSet, cache.NoExpiration)
	log.Println("access rules have been loaded successfully")
}
