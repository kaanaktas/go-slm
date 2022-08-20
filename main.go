package main

import (
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/executor"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	defer config.Elapsed("Execution")()
	defer func() {
		time.Sleep(10 * time.Millisecond)
		log.Println("All Channels were closed successfully. Number of goroutine:", runtime.NumGoroutine())
	}()

	_ = os.Setenv("GO_SLM_COMMON_POLICIES_PATH", "/testconfig/common_policies.yaml")
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/testconfig/policy_rule_set.yaml")
	_ = os.Setenv("GO_SLM_DATA_FILTER_RULE_SET_PATH", "/testconfig/custom_datafilter_rule_set.yaml")
	//pretending to be imported by another project
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")

	serviceName := "test"
	data := []string{
		"clear data with no match",
		"admin' AND 1=1 --",
		"https://testing.com/book.html?default=<script>alert(document.cookie)</script>",
		"44044 3360110004012 8888 88881881990139424332 2221111"}

	for _, data := range data {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in Execute", r)
				}
				log.Println("--------")
			}()
			log.Println("Filtering data:", data)
			executor.Execute(data, serviceName, config.Request)
		}()
	}
}
