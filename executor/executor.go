package executor

import (
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/datafilter"
	"github.com/kaanaktas/go-slm/policy"
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

var cacheIn = cache.NewInMemory()

type Specification struct {
	PolicyRuleSetPath     string `envconfig:"policy_rule_set_path"`
	CommonRulesPath       string `envconfig:"common_rules_path"`
	DataFilterRuleSetPath string `envconfig:"data_filter_rule_set_path"`
	CurrentModuleName     string `envconfig:"current_module_name"`
}

func loadConfiguration() {
	var spec Specification
	err := envconfig.Process("go_slm", &spec)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !config.IsModuleImported(spec.CurrentModuleName) || config.RootDirectory == "" {
		panic("root directory is empty or module is not imported")
	}

	policy.Load(spec.PolicyRuleSetPath, spec.CommonRulesPath)
	datafilter.Load(spec.DataFilterRuleSetPath)
}

var isConfigurationFlagSet = "isConfigurationFlagSet"

func Execute(data, serviceName, direction string) {
	if _, ok := cacheIn.Get(isConfigurationFlagSet); !ok {
		loadConfiguration()
		_ = cacheIn.Set(isConfigurationFlagSet, true, cache.NoExpiration)
	}

	policyKey := config.PolicyKey(serviceName, direction)
	cachedRule, ok := cacheIn.Get(policy.Key)
	if !ok {
		panic("policyRule doesn't exist")
	}

	policies := cachedRule.(policy.CommonPolicyMap)[policyKey]
	if len(policies) == 0 {
		log.Println("No ruleSet found for", serviceName)
		return
	}

	breaker := make(chan string)
	in := make(chan datafilter.Validate)
	closeCh := make(chan struct{})

	go processor(policies, in, breaker)
	go validator(&data, in, closeCh, breaker)

	select {
	case v := <-breaker:
		panic(v)
	case <-closeCh:
	}

	log.Println("no_match with datafilter rules")
}

func processor(policies []policy.CommonPolicy, in chan<- datafilter.Validate, breaker <-chan string) {
	defer func() {
		close(in)
	}()

	for _, v := range policies {
		if v.Active {
			if rule, ok := cacheIn.Get(v.Name); ok {
				processRule(rule.([]datafilter.Validate), in, breaker)
			}
		}
	}
}

func processRule(patterns []datafilter.Validate, in chan<- datafilter.Validate, breaker <-chan string) {
	var wg sync.WaitGroup

	for _, pattern := range patterns {
		wg.Add(1)

		pattern := pattern
		go func() {
			defer wg.Done()

			if !pattern.Disable() {
				select {
				case <-breaker:
					return
				case in <- pattern:
				}
			}
		}()
	}

	wg.Wait()
}

func validator(data *string, in <-chan datafilter.Validate, closeCh chan<- struct{}, breaker chan<- string) {
	defer func() {
		close(closeCh)
	}()

	var wg sync.WaitGroup

	//Distribute work to multiple workers
	for i := 0; i < config.NumberOfWorker; i++ {
		wg.Add(1)
		worker(&wg, data, in, breaker)
	}

	wg.Wait()
}

func worker(wg *sync.WaitGroup, data *string, in <-chan datafilter.Validate, breaker chan<- string) {
	go func() {
		defer wg.Done()

		for v := range in {
			if v.Validate(data) {
				breaker <- fmt.Sprint(v.ToString())
				return
			}
		}
	}()
}
