package datafilter

import (
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/policy"
	"testing"
)

func TestExecuteDataFilter(t *testing.T) {
	Load("")

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
				data:        "https://testing.com/book.html?default=<script>alert(document.cookie)</script>",
				serviceName: "test",
			}},
		{
			name:  "test_pan_filter",
			panic: true,
			args: args{
				data:        "44044 3360110004012 8888 88881881990139424332 2221111",
				serviceName: "test",
			}},
		{
			name:  "test_no_match",
			panic: false,
			args: args{
				data:        "test data",
				serviceName: "test3",
			}},
	}

	actions := []policy.Action{
		{Name: "sqli", Active: true},
		{Name: "xss", Active: true},
		{Name: "pan_process", Active: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) && tt.panic == false {
					t.Errorf("%s did panic", tt.name)
				} else if (r == nil) && tt.panic == true {
					t.Errorf("%s didn't panic", tt.name)
				}
			}()

			executor := &Executor{
				Actions: actions,
				Data:    &tt.args.data,
			}
			executor.Apply()
		})
	}
}

func TestDataFilterRuleLoad(t *testing.T) {
	cacheIn := cache.NewInMemory()
	cacheIn.Flush()
	Load("/testdata/custom_datafilter_rule_set.yaml")

	type cachedDataFilterRule struct {
		name string
		size int
	}

	tests := []struct {
		name   string
		policy cachedDataFilterRule
	}{
		{
			name: "cached_pan_process", policy: cachedDataFilterRule{
				name: "pan_process",
				size: 1,
			},
		},
		{
			name: "cached_custom_pan_process", policy: cachedDataFilterRule{
				name: "custom_pan_process",
				size: 1,
			},
		},
		{
			name: "cached_sqli", policy: cachedDataFilterRule{
				name: "sqli",
				size: 44,
			},
		},
		{
			name: "cached_xss", policy: cachedDataFilterRule{
				name: "xss",
				size: 27,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if cachedData, ok := cacheIn.Get(test.policy.name); ok {
				if len(cachedData.([]Validate)) != test.policy.size {
					t.Errorf("cached data size doesn't match up. Expected: %d, got:%d", test.policy.size,
						len(cachedData.([]Validate)))
				}
			} else {
				t.Errorf("%s is not in the cache", test.policy.name)
			}
		})
	}
}
