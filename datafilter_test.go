package main

import (
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/executor"
	"github.com/labstack/gommon/log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/policy/testdata/policy_rule_set.json")
	_ = os.Setenv("GO_SLM_COMMON_RULES_PATH", "/policy/testdata/common_policy_rules.json")
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")

	os.Exit(m.Run())
}

func TestExecute(t *testing.T) {
	type args struct {
		data        string
		serviceName string
	}
	tests := []struct {
		name  string
		panic bool
		args  args
	}{
		{
			name:  "test_sqli_filter",
			panic: true,
			args: args{
				data:        "admin' AND 1=1 --",
				serviceName: "test",
			}},
		{
			name:  "test_xss_filter",
			panic: true,
			args: args{
				data:        "http://testing.com/book.html?default=<script>alert(document.cookie)</script>",
				serviceName: "test",
			}},
		{
			name:  "test_pan_filter",
			panic: true,
			args: args{
				data:        "44044333322221111swfkjbfjksjkf4444333322221111dedeefefefe",
				serviceName: "test",
			}},
		{
			name:  "test_no_match",
			panic: false,
			args: args{
				data:        "test data",
				serviceName: "test",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) && tt.panic == false {
					log.Error(r)
					t.Errorf("%s did panic", tt.name)
				} else if (r == nil) && tt.panic == true {
					t.Errorf("%s didn't panic", tt.name)
				}
			}()
			executor.Execute(tt.args.data, tt.args.serviceName, config.Request)
		})
	}
}

func TestCache(t *testing.T) {
	_ = os.Setenv("GO_SLM_DATA_FILTER_RULE_SET_PATH", "/datafilter/testdata/datafilter_rule_set.json")

	cacheIn := cache.NewInMemory()
	cacheIn.Flush()

	executor.Execute("test_sqli_filter", "test", config.Request)
	if _, ok := cacheIn.Get("test_pan_process"); !ok {
		t.Error("test_pan_process is not in the cache")
	}
	if _, ok := cacheIn.Get("pan_process"); !ok {
		t.Error("pan_process is not in the cache")
	}
}
