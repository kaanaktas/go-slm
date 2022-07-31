package main

import (
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/executor"
	"log"
	"os"
	"runtime"
	"time"
)

func init() {
	_ = os.Setenv("GO_SLM_COMMON_POLICIES_PATH", "/policy/testdata/common_policies.json")
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/policy/testdata/policy_rule_set.json")
	_ = os.Setenv("GO_SLM_DATA_FILTER_RULE_SET_PATH", "/testdata/datafilter_rule_set.json")
	//pretending to be imported by another project
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")
}

func main() {
	defer config.Elapsed("Execution")()
	defer func() {
		time.Sleep(10 * time.Millisecond)
		log.Println("All Channels were closed successfully. Number of goroutine:", runtime.NumGoroutine())
	}()

	serviceName := "test"
	testData := [...]string{
		"clear data with no match",
		"admin' AND 1=1 --",
		"http://testing.com/book.html?default=<script>alert(document.cookie)</script>",
		"44044 3360110004012 8888 88881881990139424332 2221111"}

	for _, data := range testData {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in Execute", r)
				}
			}()
			log.Println("Filtering data:", data)
			executor.Execute(data, serviceName, config.Request)
		}()
	}
}
