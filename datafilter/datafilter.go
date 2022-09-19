package datafilter

import (
	"fmt"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/policy"
	"log"
	"sync"
)

type Executor struct {
	Actions []policy.Action
	Data    *string
}

func (e *Executor) Apply() {
	breaker := make(chan string)
	in := make(chan Validate)
	closeCh := make(chan struct{})

	go processor(e.Actions, in, breaker)
	go validator(e.Data, in, closeCh, breaker)

	select {
	case v := <-breaker:
		panic(v)
	case <-closeCh:
	}

	log.Println("no match with datafilter rules")
}

func processor(actions []policy.Action, in chan<- Validate, breaker <-chan string) {
	defer func() {
		close(in)
	}()

	for _, v := range actions {
		if v.Active {
			if rule, ok := cacheIn.Get(v.Name); ok {
				processRule(rule.([]Validate), in, breaker)
			}
		}
	}
}

func processRule(patterns []Validate, in chan<- Validate, breaker <-chan string) {
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

func validator(data *string, in <-chan Validate, closeCh chan<- struct{}, breaker chan<- string) {
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

func worker(wg *sync.WaitGroup, data *string, in <-chan Validate, breaker chan<- string) {
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
