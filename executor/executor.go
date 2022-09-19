package executor

import (
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/datafilter"
	"github.com/kaanaktas/go-slm/policy"
	"github.com/kaanaktas/go-slm/schedule"
	"github.com/kelseyhightower/envconfig"
	"log"
)

var cacheIn = cache.NewInMemory()

type specification struct {
	PolicyRuleSetPath     string `envconfig:"policy_rule_set_path"`
	CommonRulesPath       string `envconfig:"common_policies_path"`
	DataFilterRuleSetPath string `envconfig:"data_filter_rule_set_path"`
	SchedulePolicyPath    string `envconfig:"schedule_policy_path"`
	CurrentModuleName     string `envconfig:"current_module_name"`
}

func loadSpecification() {
	var spec specification
	if err := envconfig.Process("go_slm", &spec); err != nil {
		log.Fatal(err.Error())
	}

	if !config.IsModuleImported(spec.CurrentModuleName) || config.RootDirectory == "" {
		panic("root directory is empty or module is not imported")
	}

	schedule.Load(spec.SchedulePolicyPath)
	datafilter.Load(spec.DataFilterRuleSetPath)
	policy.Load(spec.PolicyRuleSetPath, spec.CommonRulesPath)
}

var isConfigurationFlagSet = "isConfigurationFlagSet"

func Apply(data, serviceName, direction string) {
	if _, ok := cacheIn.Get(isConfigurationFlagSet); !ok {
		loadSpecification()
		cacheIn.Set(isConfigurationFlagSet, true, cache.NoExpiration)
	}

	cachedPolicy, ok := cacheIn.Get(policy.Key)
	if !ok {
		panic("policy doesn't exist in the cache")
	}

	dynamicPolicyKey := config.PolicyKey(serviceName, direction)
	policyStatements := cachedPolicy.(policy.Statements)[dynamicPolicyKey]
	if len(policyStatements) == 0 {
		log.Println("No policy statements found for", serviceName)
		return
	}

	for _, statement := range policyStatements {
		switch statement.Type {
		case config.StatementSchedule:
			schedule.Apply(statement.Actions)
		case config.StatementData:
			datafilter.Apply(statement.Actions, &data)
		}
	}
}
