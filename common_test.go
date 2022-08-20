package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/testconfig/policy_rule_set.json")
	_ = os.Setenv("GO_SLM_COMMON_POLICIES_PATH", "/testconfig/common_policies.json")
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")

	os.Exit(m.Run())
}
