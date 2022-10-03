package policy

import (
	"github.com/kaanaktas/go-slm/cache"
	"testing"
)

func TestPolicyLoad(t *testing.T) {
	cacheIn := cache.NewInMemory()
	cacheIn.Flush()

	LoadPolicies("/testdata/policy_rule_set.yaml", "/testdata/common_policies.yaml")

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
				name: "policy_rule",
				size: 6,
				policies: []string{"test_request", "test_response",
					"test2_request", "test2_response",
					"test3_request", "test3_response"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if cachedData, ok := cacheIn.Get(test.policy.name); ok {
				if len(cachedData.(Statements)) != test.policy.size {
					t.Errorf("cached data size doesn't match up. Expected: %d, got:%d", test.policy.size,
						len(cachedData.(Statements)))
				}
				for _, v := range test.policy.policies {
					cachedPolicies := cachedData.(Statements)
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
