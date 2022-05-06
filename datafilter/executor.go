package datafilter

import (
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/policy"
	"log"
	"sync"
)

var cacheIn = cache.NewInMemory()

func Execute(data, serviceName, direction string) {
	policyKey := config.PolicyKey(serviceName, direction)
	cachedRule, ok := cacheIn.Get(access.Key)
	if !ok {
		panic("policyRule doesn't exist")
	}

	policyRules := cachedRule.(access.PolicyRules)[policyKey]
	if len(policyRules) == 0 {
		log.Println("No ruleSet found for", serviceName)
		return
	}

	breaker := make(chan string)
	in := make(chan validate)
	closeCh := make(chan struct{})

	go processor(policyRules, in, breaker)
	go validator(&data, in, closeCh, breaker)

	select {
	case v := <-breaker:
		panic(v)
	case <-closeCh:
	}

	log.Println("no_match")
}

func processor(accessList []access.Rule, in chan<- validate, breaker <-chan string) {
	defer func() {
		close(in)
	}()

	for _, v := range accessList {
		if v.Active {
			if rule, ok := cacheIn.Get(v.Name); ok {
				processRule(rule.([]validate), in, breaker)
			}
		}
	}
}

func processRule(patterns []validate, in chan<- validate, breaker <-chan string) {
	var wg sync.WaitGroup

	for _, pattern := range patterns {
		wg.Add(1)

		pattern := pattern
		go func() {
			defer wg.Done()

			if !pattern.disable() {
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

func validator(data *string, in <-chan validate, closeCh chan<- struct{}, breaker chan<- string) {
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

func worker(wg *sync.WaitGroup, data *string, in <-chan validate, breaker chan<- string) {
	go func() {
		defer wg.Done()

		for v := range in {
			if v.validate(data) {
				breaker <- fmt.Sprint(v.toString())
				return
			}
		}
	}()
}
