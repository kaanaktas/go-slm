package main

import (
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/datafilter"
	"github.com/kaanaktas/go-slm/executor"
	"github.com/kaanaktas/go-slm/policy"
	"os"
	"testing"
)

func TestDataFilterRuleLoad(t *testing.T) {
	_ = os.Setenv("GO_SLM_DATA_FILTER_RULE_SET_PATH", "/testconfig/custom_datafilter_rule_set.yaml")

	cacheIn := cache.NewInMemory()
	cacheIn.Flush()

	executor.Execute("test_data", "test", config.Request)

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
				if len(cachedData.([]datafilter.Validate)) != test.policy.size {
					t.Errorf("cached data size doesn't match up. Expected: %d, got:%d", test.policy.size,
						len(cachedData.([]datafilter.Validate)))
				}
			} else {
				t.Errorf("%s is not in the cache", test.policy.name)
			}
		})
	}
}

func TestPolicyLoad(t *testing.T) {
	cacheIn := cache.NewInMemory()
	cacheIn.Flush()

	executor.Execute("test_data", "test", config.Request)

	type cachedPolicyRule struct {
		name     string
		size     int
		policies []string
	}

	tests := []struct {
		name   string
		policy cachedPolicyRule
	}{
		{
			name: "cached_policy_rule", policy: cachedPolicyRule{
				name:     "policy_rule",
				size:     4,
				policies: []string{"test_request", "test_response", "test2_request", "test2_response"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if cachedData, ok := cacheIn.Get(test.policy.name); ok {
				if len(cachedData.(policy.CommonPolicies)) != test.policy.size {
					t.Errorf("cached data size doesn't match up. Expected: %d, got:%d", test.policy.size,
						len(cachedData.(policy.CommonPolicies)))
				}
				for _, v := range test.policy.policies {
					cachedPolicies := cachedData.(policy.CommonPolicies)
					if _, exists := cachedPolicies[v]; !exists {
						t.Errorf("%s is not in the policy rule set", v)
					}
				}
			} else {
				t.Errorf("%s is not in the policies", test.policy.name)
			}
		})
	}
}
