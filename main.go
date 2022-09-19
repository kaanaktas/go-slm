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

	_ = os.Setenv("GO_SLM_COMMON_POLICIES_PATH", "/policy/testdata/common_policies.yaml")
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/policy/testdata/policy_rule_set.yaml")
	_ = os.Setenv("GO_SLM_SCHEDULE_POLICY_PATH", "/schedule/testdata/schedule.yaml")
	//pretending to be imported by another project
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")

	serviceData := []struct {
		serviceName string
		data        string
	}{
		{
			serviceName: "test3",
			data:        "clear data with no match",
		},
		{
			serviceName: "test3",
			data:        "admin' AND 1=1 --",
		},
		{
			serviceName: "test3",
			data:        "https://testing.com/book.html?default=<script>alert(document.cookie)</script>",
		},
		{
			serviceName: "test3",
			data:        "44044 3360110004012 8888 88881881990139424332 2221111",
		},
		{
			serviceName: "test",
			data:        "catch_schedule",
		},
	}

	for _, sd := range serviceData {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in Apply", r)
				}
				log.Println("--------")
			}()
			log.Println("Filtering data:", sd.data)
			executor.Apply(sd.data, sd.serviceName, config.Request)
		}()
	}
}
