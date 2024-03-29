package executor

import (
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/datafilter"
	"github.com/kaanaktas/go-slm/policy"
	"github.com/kaanaktas/go-slm/schedule"
	"github.com/kelseyhightower/envconfig"
	"log"
)

var cacheIn = cache.NewInMemory()

type Executor interface {
	Apply()
}

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

	schedule.LoadSchedules(spec.SchedulePolicyPath)
	datafilter.LoadDataFilterRules(spec.DataFilterRuleSetPath)
	policy.LoadPolicies(spec.PolicyRuleSetPath, spec.CommonRulesPath)
}

const isConfigurationFlagSet = "isConfigurationFlagSet"

func Apply(data, serviceName, direction string) {
	if _, ok := cacheIn.Get(isConfigurationFlagSet); !ok {
		loadSpecification()
		cacheIn.Set(isConfigurationFlagSet, true, cache.NoExpiration)
	}

	cachedPolicy, ok := cacheIn.Get(policy.CacheKey)
	if !ok {
		panic("policy doesn't exist in the cache")
	}

	dynamicPolicyKey := config.PolicyKey(serviceName, direction)
	policyStatements := cachedPolicy.(policy.Statements)[dynamicPolicyKey]
	if len(policyStatements) == 0 {
		log.Printf("No policy statements found for %s", serviceName)
		return
	}

	for _, statement := range policyStatements {
		statementExecutor := createExecutor(data, statement.Type, statement.Actions)
		statementExecutor.Apply()
	}
}

func createExecutor(data string, statementType string, actions []policy.Action) Executor {
	switch statementType {
	case config.StatementSchedule:
		return &schedule.Executor{Actions: actions}
	case config.StatementData:
		return &datafilter.Executor[datafilter.Validator]{Actions: actions, Data: &data}
	default:
		panic(fmt.Sprintf("StatementType: %s doesn't exist", statementType))
	}
}
