package main

import (
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/executor"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func init() {
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")
	_ = os.Setenv("GO_SLM_SCHEDULE_POLICY_PATH", "/schedule/testdata/schedule.yaml")
	_ = os.Setenv("GO_SLM_COMMON_POLICIES_PATH", "/policy/testdata/common_policies.yaml")
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/policy/testdata/policy_rule_set.yaml")

	log.Println("Starting with number of worker", config.NumberOfWorker)
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func Benchmark(b *testing.B) {

	content := config.MustReadFile("benchmark_load.json")

	data := string(content)
	serviceName := "test3"

	for i := 0; i < b.N; i++ {
		executor.Apply(data, serviceName, config.Request)
	}
}
