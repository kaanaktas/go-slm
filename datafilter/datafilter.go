package datafilter

import (
	"fmt"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/policy"
	"log"
	"sync"
)

type Executor[T Validator] struct {
	Actions []policy.Action
	Data    *string
}

func (e *Executor[T]) Apply() {
	breaker := make(chan string)
	in := make(chan T)
	closeCh := make(chan struct{})

	go processor[T](e.Actions, in, breaker)
	go validator[T](e.Data, in, closeCh, breaker)

	select {
	case v := <-breaker:
		panic(v)
	case <-closeCh:
		log.Println("no match with datafilter rules")
	}
}

func processor[T Validator](actions []policy.Action, in chan<- T, breaker <-chan string) {
	defer close(in)

	for _, v := range actions {
		if v.Active {
			if rule, ok := cacheIn.Get(v.Name); ok {
				processRule(rule.([]T), in, breaker)
			}
		}
	}
}

func processRule[T Validator](patterns []T, in chan<- T, breaker <-chan string) {
	var wg sync.WaitGroup

	for _, pattern := range patterns {
		wg.Add(1)

		pattern := pattern
		go func() {
			defer wg.Done()

			if !pattern.IsDisabled() {
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

func validator[T Validator](data *string, in <-chan T, closeCh chan<- struct{}, breaker chan<- string) {
	defer func() {
		close(closeCh)
	}()

	var wg sync.WaitGroup

	//Distribute work to multiple workers
	for i := 0; i < config.NumberOfWorker; i++ {
		wg.Add(1)
		worker[T](&wg, data, in, breaker)
	}

	wg.Wait()
}

func worker[T Validator](wg *sync.WaitGroup, data *string, in <-chan T, breaker chan<- string) {
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
